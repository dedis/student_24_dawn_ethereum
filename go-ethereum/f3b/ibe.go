// Copyright EPFL DEDIS

package f3b

import (
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing/bn256"
	"go.dedis.ch/kyber/v3/util/random"
)

var Suite = bn256.NewSuite()

type hashablePoint interface {
	Hash([]byte) kyber.Point
}

var hashable = Suite.G1().Point().(hashablePoint)

func HashToG1(label []byte) kyber.Point {
	return hashable.Hash(label)
}

// Share a secret with the SMC pk by computing
// r = random scalar
// U = rg2
// secret = rP
func ShareSecret(pk kyber.Point, label []byte) (U, secret kyber.Point) {
	P := Suite.Pair(HashToG1(label), pk)
	r := Suite.G2().Scalar().Pick(random.New())
	U = Suite.G2().Point().Mul(r, Suite.G2().Point().Base())
	secret = Suite.GT().Point().Mul(r, P)
	return
}

func RecoverSecret(sigma kyber.Point, U kyber.Point) kyber.Point {
	return Suite.Pair(sigma, U)
}

func VerifyIdentity(pk, sigma kyber.Point, label []byte) (bool, error) {
	// note: this could be optimized
	nbase := Suite.G2().Point().Neg(Suite.G2().Point().Base())
	return Suite.PairingCheck([]kyber.Point{sigma, HashToG1(label)}, []kyber.Point{nbase, pk}), nil
}
