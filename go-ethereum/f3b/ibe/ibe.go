// Copyright EPFL DEDIS

package ibe

import (
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing/bn256"
	"go.dedis.ch/kyber/v3/util/random"
)

var Suite = bn256.NewSuite()

type hashablePoint interface {
	Hash([]byte) kyber.Point
}

func HashToG1(label []byte) kyber.Point {
	return Suite.G1().Point().(hashablePoint).Hash(label)
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

func VerifyIdentitySlow(pk, sigma kyber.Point, label []byte) bool {
       lhs := Suite.Pair(sigma, Suite.G2().Point().Base())
       rhs := Suite.Pair(HashToG1(label), pk)
       return lhs.Equal(rhs)
}

func VerifyIdentityFast(pk, sigma kyber.Point, label []byte) bool {
	nbase := Suite.G2().Point().Neg(Suite.G2().Point().Base())
	return Suite.PairingCheck([]kyber.Point{sigma, HashToG1(label)}, []kyber.Point{nbase, pk})
}

var VerifyIdentity = VerifyIdentityFast
