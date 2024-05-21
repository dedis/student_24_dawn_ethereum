// Copyright EPFL DEDIS

package f3b

import (
	"github.com/ethereum/go-ethereum/f3b/ibe"

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

func (e *TPKE) ShareSecret(label []byte) (seed, encKey []byte, err error) {
	//label := binary.BigEndian.AppendUint64(from.Bytes(), tx.Nonce())
	U, secret := ibe.ShareSecret(e.pk, label)

	//plaintext := append(tx.To().Bytes(), tx.Data()...)
	//ciphertext = make([]byte, len(plaintext))
	//tag = make([]byte, e.cae.TagLen())
	seed, err = secret.MarshalBinary()
	if err != nil {
		return nil, nil, err
	}
	//err = e.cae.Encrypt(ciphertext, tag, seed, plaintext)
	//if err != nil {
	//	return nil, err
	//}
	encKey, err = U.MarshalBinary()
	if err != nil {
		return nil, nil, err
	}
	return
}

func (e *TPKE) RevealSecret(label, encKey []byte) (reveal []byte, err error) {
	return e.smccli.Extract(label)
}

func (e *TPKE) RecoverSecret(encKey, reveal []byte) (seed []byte, err error) {
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

	/* FIXME: not possible to have the label due to circular dependency
	if !ibe.VerifyIdentity(e.pk, identity, label) {
		return nil, errors.New("bad identity")
	}
	*/

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
