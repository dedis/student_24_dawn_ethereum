// Copyright EPFL DEDIS

package vdf

import (
	"crypto/rand"
	"crypto/rsa"
	"math/big"
	"fmt"

	"go.dedis.ch/kyber/v3/suites"
)

var Suite = suites.MustFind("ed25519")

const RsaBits = 2048
const SquaringSteps uint64 = 1_000_000

// Share a secret with the future
func ShareSecret(label []byte) (n, secret *big.Int) {
	return shareSecret(label, SquaringSteps)
}

func shareSecret(label []byte, steps uint64) (n, secret *big.Int) {
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
	t := new(big.Int).Exp(two, new(big.Int).SetUint64(steps), φ)
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
	var tmp big.Int
	g := deriveInitial(n, label)
	x := new(big.Int).Set(g)
	y := new(big.Int)
	/*
	κ := uint64(16) // TODO: set based on steps?
	κ = 1
	// TODO: γ?
	memo := make([]big.Int, (steps + 1<<κ - 1) >> κ)
	for i := uint64(0); i < steps; i++ {
		if i % (1<<κ) == 0 {
			memo[i >> κ].Set(x)
		}
		x.Mul(x, x)
		x.Mod(x, n)
	}
	for b := uint64(0): b < 1 << κ; b++ {
		r.Mul(s.Mul(b,s))
	one := new(big.Int).SetInt64(1)
	two := new(big.Int).SetInt64(2)
	l := sampleL()
	π := new(big.Int).Set(one)
	fmt.Println((steps >> κ) + 1)
	for i := uint64(0); i < (steps >> κ); i++ {
		fmt.Println("i", i)
		b := new(big.Int)
		b.Exp(two, tmp.SetUint64(steps  >> (κ*(i+1))), l)
		b.Mul(b, tmp.SetUint64(1 << κ))
		b.Div(b, l)
		c := &memo[i] // g**(2**(κ*i)) mod n
		b.Exp(c, b, n)
		π.Mul(π, b)
	}
	*/
	for i := uint64(0); i < steps; i++ {
		x.Mul(x, x)
		x.Mod(x, n)
	}

	// Long-division slow way based on https://eprint.iacr.org/2018/712
	y.Set(x)
	r := new(big.Int).SetUint64(1)
	two := new(big.Int).SetInt64(2)
	l := sampleL()
	x.SetUint64(1)
	for i := uint64(0); i < steps; i++ {
		b := new(big.Int)
		b.Mul(two, r).Div(b,l)
		r.Mul(r, two).Mod(r, l)
		x.Mul(x, x).Mul(x, tmp.Exp(g, b, n))
		x.Mod(x, n)
		fmt.Println("x", x)
		fmt.Println("b", b)
	}
	π := x
	return y, π, l
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

func checkProof(g, y, π, l, n *big.Int, steps uint64) bool {
	fmt.Println("g", g)
	fmt.Println("y", y)
	fmt.Println("π", π)
	fmt.Println("l", l)
	fmt.Println("n", n)
	two := new(big.Int).SetInt64(2)
	t := new(big.Int).SetUint64(steps)
	// r = 2**t mod l
	r := new(big.Int).Exp(two, t, l)
	// π**l * g**r == y
	y2 := new(big.Int).Exp(π, l, n)
	y2.Mul(y2, new(big.Int).Exp(g, r, n))
	y2.Mod(y2, n)
	return y2.Cmp(y) == 0
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
