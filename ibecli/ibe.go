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
	RP kyber.Point
	C  []byte
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

func EncryptCPAonG2(s pairing.Suite, pubKey kyber.Point, label, msg []byte) (*CiphertextCPA, error) {
	hashable, ok := s.G1().Point().(hashablePoint)
	if !ok {
		return nil, errors.New("point needs to implement hashablePoint")
	}
	Qid := hashable.Hash(label)
	r := s.G2().Scalar().Pick(random.New())
	rP := s.G2().Point().Mul(r, s.G2().Point().Base())

	// e(Qid, Ppub) = e( H(round), s*P) where s is dist secret key
	Ppub := pubKey
	rQid := s.G1().Point().Mul(r, Qid)
	GidT := s.Pair(rQid, Ppub)

	xof, err := gtToXOF(GidT)
	if err != nil {
		return nil, err
	}
	xored := make([]byte, len(msg))
	xof.XORKeyStream(xored, msg)

	return &CiphertextCPA{rP, xored}, nil
}

func DecryptCPAonG2(s pairing.Suite, private kyber.Point, c *CiphertextCPA) ([]byte, error) {
	GidT := s.Pair(private, c.RP)
	xof, err := gtToXOF(GidT)
	if err != nil {
		return nil, err
	}
	xored := make([]byte, len(c.C))
	xof.XORKeyStream(xored, c.C)
	return xored, nil
}
