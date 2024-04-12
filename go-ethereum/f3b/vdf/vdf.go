// Copyright EPFL DEDIS

package vdf

import (
	"crypto/rand"
	"crypto/rsa"
	"math/big"

	"go.dedis.ch/kyber/v3/suites"
)

var Suite = suites.MustFind("ed25519")

const RsaBits = 2048
const SquaringSteps uint64 = 1_000_000

// Share a secret with the SMC pk by computing
// r = random scalar
// U = rg2
// secret = rP
func ShareSecret(label []byte) (n, secret *big.Int) {
	priv, err := rsa.GenerateKey(rand.Reader, RsaBits)
	if err != nil {
		panic(err.Error())
	}
	one := new(big.Int).SetInt64(1)
	p_ := new(big.Int).Sub(priv.Primes[0], one)
	q_ := new(big.Int).Sub(priv.Primes[1], one)
	φ := new(big.Int).Mul(p_, q_)
	init := deriveInitial(priv.N, label)
	two := new(big.Int).SetInt64(2)
	t := new(big.Int).Exp(two, new(big.Int).SetUint64(SquaringSteps), φ)
	secret = new(big.Int).Exp(init, t, priv.N)
	return priv.N, secret
}

func deriveInitial(n *big.Int, label []byte) *big.Int {
	init, err := rand.Int(Suite.XOF(label), n)
	if err != nil {
		panic(err.Error())
	}
	return init
}


func RecoverSecret(n *big.Int, label []byte) *big.Int {
	return recoverSecret(n, label, SquaringSteps)
}

func recoverSecret(n *big.Int, label []byte, steps uint64) *big.Int {
	x := deriveInitial(n, label)
	for i := uint64(0); i < steps; i++ {
		x.Mul(x, x)
		x.Mod(x, n)
	}
	return x
}

func recoverSecretWithProof(n *big.Int, label []byte, steps uint64) (*big.Int, *big.Int, *big.Int) {
	x := deriveInitial(n, label)
	κ := uint64(16) // TODO: set based on steps?
	// TODO: γ?
	memo := make([]*big.Int, 0, steps >> κ)
	for i := uint64(0); i < steps; i++ {
		if i % (1<<κ) == 0 {
			memo = append(memo, new(big.Int).Set(x))
		}
		x.Mul(x, x)
		x.Mod(x, n)
	}
	/*
	for b := uint64(0): b < 1 << κ; b++ {
		r.Mul(s.Mul(b,s))
	*/
	one := new(big.Int).SetInt64(1)
	two := new(big.Int).SetInt64(2)
	l := sampleL()
	π := new(big.Int).Set(one)
	for i := uint64(0); i < steps >> κ + 1; i++ {
		b := new(big.Int).SetUint64(steps - κ*(i+1))
		b.Exp(two, b, l)
		b.Div(b, l)
		c := memo[i] // 2**(κ*i) mod n
		b.Exp(c, b, n)
		π.Mul(π, b)
	}
	return x, π, l
}

func sampleL() *big.Int {
	// FIXME: fiat shamir
	kBits := 128
	l, err := rand.Prime(rand.Reader, kBits+1)
	if err != nil {
		panic(err.Error())
	}
	return l
}

/*
func prove(n, g, y, t *big.Int) (*big.Int, *big.Int) {
	one := new(big.Int).SetInt64(1)
	p_ := new(big.Int).Sub(priv.Primes[0], one)
	q_ := new(big.Int).Sub(priv.Primes[1], one)
	φ := new(big.Int).Mul(p_, q_)
	init := deriveInitial(priv.N, label)
	two := new(big.Int).SetInt64(2)
	t := new(big.Int).Exp(two, new(big.Int).SetUint64(SquaringSteps), φ)
	y := new(big.Int).Exp(init, t, priv.N)

	// proof
	l := sampleL()
	π := new(big.Int).Lsh(one, uint(steps))
	π.Div(π, l)
	π.Exp(init, π, priv.N)

	return y, π
}
*/
