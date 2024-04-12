// Copyright EPFL DEDIS

package f3b

import (
	"crypto/rand"
	"crypto/rsa"
	"math/big"
)

const RsaBits = 2048
const SquaringSteps uint64 = 1_000_000

// Share a secret with the SMC pk by computing
// r = random scalar
// U = rg2
// secret = rP
func VdfShareSecret(label []byte) (n, secret *big.Int) {
	priv, err := rsa.GenerateKey(rand.Reader, RsaBits)
	if err != nil {
		panic(err.Error())
	}
	one := new(big.Int).SetInt64(1)
	p_ := new(big.Int).Sub(priv.Primes[0], one)
	q_ := new(big.Int).Sub(priv.Primes[1], one)
	φ := new(big.Int).Mul(p_, q_)
	init := vdfDeriveInitial(priv.N, label)
	two := new(big.Int).SetInt64(2)
	t := new(big.Int).Exp(two, new(big.Int).SetUint64(SquaringSteps), φ)
	secret = new(big.Int).Exp(init, t, priv.N)
	return priv.N, secret
}

func vdfDeriveInitial(n *big.Int, label []byte) *big.Int {
	init, err := rand.Int(Suite.XOF(label), n)
	if err != nil {
		panic(err.Error())
	}
	return init
}


func VdfRecoverSecret(n *big.Int, label []byte) *big.Int {
	x := vdfDeriveInitial(n, label)
	for i := uint64(0); i < SquaringSteps; i++ {
		x.Mul(x, x)
		x.Mod(x, n)
	}
	return x
}
