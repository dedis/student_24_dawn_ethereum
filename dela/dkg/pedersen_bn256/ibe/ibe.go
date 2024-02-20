// This code is (c) by DEDIS/EPFL 2017 under the MPL v2 or later version.

package ibe

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"

	kyber "go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing"
	"go.dedis.ch/kyber/v3/util/random"

	"golang.org/x/crypto/hkdf"
)

const tagSize = 32

// Based on https://github.com/drand/kyber/blob/master/encrypt/ibe/ibe.go
// with G1 and G2 flipped and AES-CTR instead of Blake2s
// and appended HMAC-SHA256 tag

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
	buf = append(buf, marshalledU...)
	buf = append(buf, ct.V...)
	return buf, nil
}
func (ct *CiphertextCPA) Deserialize(suite pairing.Suite, data []byte) error {
	if len(data) < pointMarshalledSize+tagSize {
		return errors.New("data too short")
	}
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

type sessionKeys struct {
	encKey [32]byte
	macKey [32]byte
}

func (keys *sessionKeys) fromGt(GidT kyber.Point) error {
	seed, err := GidT.MarshalBinary()
	if err != nil {
		return err
	}
	kdf := hkdf.New(sha512.New, seed, nil, nil)
	kdf.Read(keys.encKey[:])
	kdf.Read(keys.macKey[:])
	return nil
}

func (keys *sessionKeys) keyStream() (cipher.Stream, error) {
	block, err := aes.NewCipher(keys.encKey[:])
	if err != nil {
		return nil, err
	}
	iv := make([]byte, 16) //NOTE: this is allowed because the key is used only once
	return cipher.NewCTR(block, iv), nil
}

func (keys *sessionKeys) mac() (hash.Hash, error) {
	return hmac.New(sha256.New, keys.macKey[:]), nil
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
	var keys sessionKeys
	r := suite.G2().Scalar().Pick(random.New())
	U := suite.G2().Point().Mul(r, suite.G2().Point().Base())
	err := keys.fromGt(suite.GT().Point().Mul(r, P))
	if err != nil {
		return nil, err
	}
	xof, err := keys.keyStream()
	if err != nil {
		return nil, err
	}
	V := make([]byte, len(msg), len(msg)+tagSize)
	xof.XORKeyStream(V, msg)

	mac, err := keys.mac()
	mac.Write(V)
	V = mac.Sum(V)

	return &CiphertextCPA{U, V}, nil
}

func DecryptCPAonG2(suite pairing.Suite, decKey kyber.Point, ct *CiphertextCPA) ([]byte, error) {
	var keys sessionKeys
	err := keys.fromGt(suite.Pair(decKey, ct.U))
	if err != nil {
		return nil, err
	}

	mac, err := keys.mac()
	mac.Write(ct.V[:len(ct.V)-tagSize])
	tag := mac.Sum(nil)
	if !hmac.Equal(tag, ct.V[len(ct.V)-tagSize:]) {
		return nil, errors.New("MAC mismatch")
	}

	xof, err := keys.keyStream()
	if err != nil {
		return nil, err
	}
	msg := make([]byte, len(ct.V[:len(ct.V)-tagSize]))
	xof.XORKeyStream(msg, ct.V[:len(ct.V)-tagSize])

	return msg, nil
}
