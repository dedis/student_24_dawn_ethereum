// Copyright EPFL DEDIS

package f3b

import (
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/util/random"
)

const pointMarshalledSize = 128

func HashToG1(label []byte) kyber.Point {
	// sync with kyber/sign/bls
	h := Suite.Hash()
	h.Write(label)
	x := Suite.G1().Scalar().SetBytes(h.Sum(nil))
	return Suite.G1().Point().Mul(x, nil)
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
	// FIXME do pair(...) == 0 for performance
	lhs := Suite.Pair(sigma, Suite.G2().Point())
	rhs := Suite.Pair(HashToG1(label), pk)
	return lhs.Equal(rhs), nil
}
