// This code is (c) by DEDIS/EPFL 2017 under the MPL v2 or later version.

package ibe

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"errors"
	"fmt"

	kyber "go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing"
	"go.dedis.ch/kyber/v3/util/random"

	"golang.org/x/crypto/hkdf"
)

// Based on https://github.com/drand/kyber/blob/master/encrypt/ibe/ibe.go
// with G1 and G2 flipped and AES-CTR instead of Blake2s

type CiphertextCPA struct {
	U kyber.Point
	V []byte
}

const pointMarshalledSize = 128

func (ct *CiphertextCPA) Serialize(suite pairing.Suite) ([]byte, error) {
	marshalledU, err := ct.U.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if len(marshalledU) != pointMarshalledSize {
		return nil, fmt.Errorf("unexpected point marshalled size: %v", len(marshalledU))
	}
	buf := make([]byte, 0, pointMarshalledSize+len(ct.V))
	buf =  append(buf, marshalledU...)
	buf =  append(buf, ct.V...)
	return buf, nil
}
func (ct *CiphertextCPA) Deserialize(suite pairing.Suite, data []byte) error {
	marshalledU := data[:pointMarshalledSize]
	U := suite.G2().Point()
	U.UnmarshalBinary(marshalledU)
	ct.U = U
	ct.V = bytes.Clone(data[pointMarshalledSize:])
	return nil
}

type hashablePoint interface {
	Hash([]byte) kyber.Point
}

func gtToStream(GidT kyber.Point) (cipher.Stream, error) {
	seed, err := GidT.MarshalBinary()
	if err != nil {
		return nil, err
	}
	key := make([]byte, 32)
	kdf := hkdf.New(sha512.New, seed, nil, nil)
	kdf.Read(key)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := make([]byte, 16) //NOTE: this is allowed because the key is used only once
	return cipher.NewCTR(block, iv), nil
}

func DeriveEncryptionKeyOnG2(suite pairing.Suite, X kyber.Point, label []byte) (kyber.Point, error) {
	hashable, ok := suite.G1().Point().(hashablePoint)
	if !ok {
		return nil, errors.New("point needs to implement hashablePoint")
	}
	P := suite.Pair(hashable.Hash(label), X)
	return P, nil
}

func EncryptCPAonG2(suite pairing.Suite, P kyber.Point, msg []byte) (*CiphertextCPA, error) {
	r := suite.G2().Scalar().Pick(random.New())
	U := suite.G2().Point().Mul(r, suite.G2().Point().Base())
	xof, err := gtToStream(suite.GT().Point().Mul(r, P))
	if err != nil {
		return nil, err
	}
	V := make([]byte, len(msg))
	xof.XORKeyStream(V, msg)

	return &CiphertextCPA{U, V}, nil
}

func DecryptCPAonG2(suite pairing.Suite, decKey kyber.Point, ct *CiphertextCPA) ([]byte, error) {
	xof, err := gtToStream(suite.Pair(decKey, ct.U))
	if err != nil {
		return nil, err
	}
	msg := make([]byte, len(ct.V))
	xof.XORKeyStream(msg, ct.V)
	return msg, nil
}

func VerifyIdentityOnG2(suite pairing.Suite, pk, identity kyber.Point, label []byte) (bool, error) {
	hashable, ok := suite.G1().Point().(hashablePoint)
	if !ok {
		return false, errors.New("point needs to implement hashablePoint")
	}

	lhs := suite.Pair(identity, suite.G2().Point())
	rhs := suite.Pair(hashable.Hash(label), pk)
	return lhs.Equal(rhs), nil
}
