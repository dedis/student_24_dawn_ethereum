// Copyright EPFL DEDIS

package f3b

import (
	"errors"

	"github.com/ethereum/go-ethereum/f3b/ibe"
	"github.com/ethereum/go-ethereum/log"

	"go.dedis.ch/kyber/v3"
)

type TPKE struct {
	pk kyber.Point
	smccli *ibe.SmcCli
}

func NewTPKE() (Protocol, error) {
	smccli := ibe.NewSmcCli()

	pk, err := smccli.GetPublicKey()
	if err != nil {
		return nil, err
	}

	return &TPKE{pk, smccli}, nil
}

func (e *TPKE) ShareSecret(labelBytes []byte) (seed, encKey []byte, err error) {
	var label ibe.Label
	copy(label[:], labelBytes)

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

func (e *TPKE) RevealSecret(labelBytes []byte, encKey []byte) (reveal []byte, err error) {
	var label ibe.Label
	copy(label[:], labelBytes)

	return e.smccli.Extract(label)
}

func (e *TPKE) RecoverSecret(labelBytes []byte, encKey, reveal []byte) (seed []byte, err error) {
	var label ibe.Label
	copy(label[:], labelBytes)

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

	log.Info("Identity", "label", label, "reveal", reveal, "identity", identity)
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
