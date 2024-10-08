// Copyright EPFL DEDIS

package f3b

import (
	"errors"

	"github.com/ethereum/go-ethereum/f3b/ibe"

	"go.dedis.ch/kyber/v3"
)

type TPKE struct {
	pk kyber.Point
	smccli SmcCli
}

func NewTPKE(smccli SmcCli) (Protocol, error) {
	pk, err := smccli.GetPublicKey()
	if err != nil {
		return nil, err
	}

	return &TPKE{pk, smccli}, nil
}

func (e *TPKE) ShareSecret(label []byte) (seed, encKey []byte, err error) {
	U, secret := ibe.ShareSecret(e.pk, label)

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

func (e *TPKE) RevealSecret(label []byte, encKey []byte) (reveal []byte, err error) {
	return e.smccli.Extract(label)
}

func (e *TPKE) RecoverSecret(label []byte, encKey, reveal []byte) (seed []byte, err error) {
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

	if !ibe.VerifyIdentity(e.pk, identity, label) {
		return nil, errors.New("bad identity")
	}

	secret := ibe.RecoverSecret(identity, U)
	seed, err = secret.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return
}

func (e *TPKE) IsVdf() bool {
	return false
}

func (e *TPKE) IsTibe() bool {
	return false
}

func (_ *TPKE) IsTpke() bool {
	return true
}
