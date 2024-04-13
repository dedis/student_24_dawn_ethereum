// Copyright EPFL DEDIS

package vdf

import (
	"crypto/rand"
	"crypto/rsa"
	"math/big"

	"go.dedis.ch/kyber/v3/suites"

	"github.com/ethereum/go-ethereum/common"
)

var Suite = suites.MustFind("ed25519")

const RsaBits = 2048

// Share a secret with the future
func ShareSecret(label []byte, log2t int) (secret, n *big.Int) {
	priv, err := rsa.GenerateKey(rand.Reader, RsaBits)
	if err != nil {
		// unreachable if rand.Reader is well-behaved
		panic(err)
	}
	var tmp big.Int
	pMinusOne := new(big.Int).Sub(priv.Primes[0], common.Big1)
	qMinusOne := new(big.Int).Sub(priv.Primes[1], common.Big1)
	φ := new(big.Int).Mul(pMinusOne, qMinusOne)
	g := deriveInitial(label, priv.N)
	secret = new(big.Int).Exp(g, tmp.Exp(common.Big2, tmp.SetUint64(1 << log2t), φ) , priv.N)
	return priv.N, secret
}

func deriveInitial(label []byte, n *big.Int) *big.Int {
	init, err := rand.Int(Suite.XOF(label), n)
	if err != nil {
		panic(err)
	}
	return init
}


func RecoverSecret(label []byte, n *big.Int, log2t int) *big.Int {
	x := deriveInitial(label, n)
	t := 1 << log2t
	for i := 0; i < t; i++ {
		x.Mul(x, x)
		x.Mod(x, n)
	}
	return x
}

func Proof(label []byte, n *big.Int, log2t int) (l *big.Int, π *big.Int) {
	var tmp big.Int
	g := deriveInitial(label, n)
	x := new(big.Int).Set(g)
	κ := uint64(log2t / 2) // TODO: set based on t?
	κ = 32
	t := uint64(1 << log2t)
	// FIXME: assumes κ divides t
	// TODO: γ? O(sqrt(t)) memory
	memo := make([]big.Int, (t+κ-1) / κ)
	for i := uint64(0); i < t; i++ {
		if i % κ == 0 {
			memo[i / κ].Set(x)
		}
		x.Mul(x, x)
		x.Mod(x, n)
	}
	y := new(big.Int).Set(x)
	l = sampleL(g, y)
	x.SetUint64(1)
	b := new(big.Int)
	for i := uint64(0); i < (t+κ-1) / κ; i++ {
		b.Exp(common.Big2, tmp.SetUint64(t  - (κ*(i+1))), l)
		b.Mul(b, tmp.SetUint64(1 << κ))
		b.Div(b, l)
		c := &memo[i] // g**(2**(κ*i)) mod n
		b.Exp(c, b, n)
		x.Mul(x, b)
		x.Mod(x, n)
	}
	if t % κ != 0 {
		// FIXME: no workee
		b.Exp(common.Big2, tmp.SetUint64(t % κ), l)
		b.Mul(b, tmp.SetUint64(1 << (t % κ)))
		b.Div(b, l)
		c := y
		b.Exp(c, b, n)
		x.Mul(x, b)
		x.Mod(x, n)
	}

	return l, x
}

// fiat shamir prime generation
func sampleL(g, y *big.Int) *big.Int {
	kBits := 128
	xof := Suite.XOF(g.Bytes())
	xof.Write(y.Bytes())
	l, err := DeterministicPrime(xof, kBits+1)
	if err != nil {
		panic(err)
	}
	return l
}

func RecoverSecretFromProof(label []byte, l, π, n *big.Int, log2t int) (y *big.Int, ok bool) {
	g := deriveInitial(label, n)
	t := new(big.Int).SetUint64(1 << log2t)
	// r = 2**t mod l
	r := new(big.Int).Exp(common.Big2, t, l)
	// π**l * g**r == y
	y = new(big.Int).Exp(π, l, n)
	y.Mul(y, new(big.Int).Exp(g, r, n))
	y.Mod(y, n)
	if sampleL(g, y).Cmp(l) != 0 {
		return nil, false
	}
	return y, true
}
