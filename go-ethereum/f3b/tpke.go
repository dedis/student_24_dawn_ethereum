// Copyright EPFL DEDIS

package f3b

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/cae"

	"go.dedis.ch/kyber/v3"
)

type EncryptFn func(from common.Address, tx *types.Transaction) (*types.Transaction, error)

type TPKE struct {
	cae cae.Scheme
	pk kyber.Point
}

func NewTPKE(cae cae.Scheme) (*TPKE, error) {
	smccli := NewSmcCli()

	pk, err := smccli.GetPublicKey()
	if err != nil {
		return nil, err
	}

	return &TPKE{cae, pk}, nil
}

func (e *TPKE) EncryptTx(from common.Address, tx *types.Transaction) (*types.Transaction, error) {
	label := binary.BigEndian.AppendUint64(from.Bytes(), tx.Nonce())
	U, secret := ShareSecret(e.pk, label)

	plaintext := append(tx.To().Bytes(), tx.Data()...)
	ciphertext := make([]byte, len(plaintext))
	tag := make([]byte, e.cae.TagLen())
	seed, err := secret.MarshalBinary()
	if err != nil {
		return nil, err
	}
	err = e.cae.Encrypt(ciphertext, tag, seed, plaintext)
	if err != nil {
		return nil, err
	}
	encKey, err := U.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return types.NewTx(&types.EncryptedTx{
		ChainID:    tx.ChainId(),
		Nonce:      tx.Nonce(),
		GasFeeCap:  tx.GasPrice(),
		Gas:        tx.Gas(),
		Value:      tx.Value(),
		Ciphertext: ciphertext,
		Tag:        tag,
		EncKey:     encKey,
	}), nil
}
