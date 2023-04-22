// This code is (c) by DEDIS/EPFL 2017 under the MPL v2 or later version.

package main

import (
	"errors"

	kyber "go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing"
	"go.dedis.ch/kyber/v3/util/random"
	"go.dedis.ch/kyber/v3/xof/keccak"
)

// Based on https://github.com/drand/kyber/blob/master/encrypt/ibe/ibe.go
// with G1 and G2 flipped and Keccak instead of Blake2s

type CiphertextCPA struct {
	U kyber.Point
	V []byte
}

type hashablePoint interface {
	Hash([]byte) kyber.Point
}

func gtToXOF(GidT kyber.Point) (kyber.XOF, error) {
	seed, err := GidT.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return keccak.New(seed), nil
}

func EncryptCPAonG2(suite pairing.Suite, X kyber.Point, label, msg []byte) (*CiphertextCPA, error) {
	hashable, ok := suite.G1().Point().(hashablePoint)
	if !ok {
		return nil, errors.New("point needs to implement hashablePoint")
	}
	P := suite.Pair(hashable.Hash(label), X)
	r := suite.G2().Scalar().Pick(random.New())
	U := suite.G2().Point().Mul(r, suite.G2().Point().Base())
	xof, err := gtToXOF(suite.GT().Point().Mul(r, P))
	if err != nil {
		return nil, err
	}
	V := make([]byte, len(msg))
	xof.XORKeyStream(V, msg)

	return &CiphertextCPA{U, V}, nil
}

func DecryptCPAonG2(suite pairing.Suite, private kyber.Point, ct *CiphertextCPA) ([]byte, error) {
	xof, err := gtToXOF(suite.Pair(private, ct.U))
	if err != nil {
		return nil, err
	}
	msg := make([]byte, len(ct.V))
	xof.XORKeyStream(msg, ct.V)
	return msg, nil
}
