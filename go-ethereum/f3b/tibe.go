// Copyright EPFL DEDIS

package f3b

import (
	"errors"

	"github.com/ethereum/go-ethereum/f3b/ibe"

	"go.dedis.ch/kyber/v3"
)

// fixed length array for map[]
const LabelLength = 8
type Label [LabelLength]byte

type TIBE struct {
	pk kyber.Point
	smccli SmcCli
	cache map[Label][]byte
}

func NewTIBE(smccli SmcCli) (Protocol, error) {
	pk, err := smccli.GetPublicKey()
	if err != nil {
		return nil, err
	}

	return &TIBE{pk, smccli, make(map[Label][]byte)}, nil
}

func (e *TIBE) ShareSecret(labelBytes []byte) (seed, encKey []byte, err error) {
	U, secret := ibe.ShareSecret(e.pk, labelBytes)

	seed, err = secret.MarshalBinary()
	if err != nil {
		return nil, nil, err
	}
	encKey, err = U.MarshalBinary()
	if err != nil {
		return nil, nil, err
	}
	return
}

func (e *TIBE) RevealSecret(labelBytes []byte, encKey []byte) (reveal []byte, err error) {
	var label Label
	copy(label[:], labelBytes)
	reveal, ok := e.cache[label]
	if !ok {
		reveal, err = e.smccli.Extract(labelBytes)
		if err != nil {
			return nil, err
		}
		e.cache[label] = reveal
	}
	return reveal, nil

}

func (e *TIBE) RecoverSecret(labelBytes []byte, encKey, reveal []byte) (seed []byte, err error) {
	U := ibe.Suite.G2().Point()
	err = U.UnmarshalBinary(encKey)
	if err != nil {
		return nil, err
	}

	identity := ibe.Suite.G1().Point()
	err = identity.UnmarshalBinary(reveal)
	if err != nil {
		return nil, err
	}

	if !ibe.VerifyIdentity(e.pk, identity, labelBytes) {
		return nil, errors.New("bad identity")
	}

	secret := ibe.RecoverSecret(identity, U)
	seed, err = secret.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return
}

func (e *TIBE) IsVdf() bool {
	return false
}

func (e *TIBE) IsTibe() bool {
	return true
}

func (_ *TIBE) IsTpke() bool {
	return false
}
