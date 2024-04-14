// Copyright EPFL DEDIS

package types

import (
	"errors"
	"encoding/binary"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/cae"
	"github.com/ethereum/go-ethereum/f3b/vdf"
	"github.com/ethereum/go-ethereum/log"
)

const Log2t = 15

func (t *Transaction) Decrypt() (*Transaction, error) {
	// Minimal signer for an encrypted transaction
	signer := NewLausanneSigner(t.ChainId())

	log.Info("Decrypt()", "tx", t.inner, "signer", signer, "sigcache", t.from.Load().(sigCache).signer)

	from, err := Sender(signer, t)
	if err != nil {
		return nil, err
	}

	tx, ok := t.inner.(*EncryptedTx)
	if !ok {
		return nil, errors.New("cannot decrypt a non-encrypted transaction")
	}

	label := binary.BigEndian.AppendUint64(from.Bytes(), tx.Nonce)
	n := new(big.Int).SetBytes(tx.N)
	l, π := vdf.Proof(label, n, Log2t)
	secret, ok := vdf.RecoverSecretFromProof(label, l, π, n, Log2t)
	if !ok {
		// NOTE: should not happen since it's our proof
		return nil, errors.New("bad VDF proof")
	}

	seed := secret.Bytes()
	log.Info("Decrypt()", "label", label, "secret", secret, "seed", seed)

	// TODO: if the ciphertext is too short, penalize the sender
	plaintext := make([]byte, len(tx.Ciphertext))
	err = cae.Selected.Decrypt(plaintext, seed, tx.Ciphertext, tx.Tag)
	// TODO: if this is an authentication error, penalize the sender
	if err != nil {
		return nil, err
	}

	to := common.BytesToAddress(plaintext[:common.AddressLength])
	data := plaintext[common.AddressLength:]

	return NewTx(&DecryptedTx{
		ChainID:    tx.ChainID,
		Nonce:      tx.Nonce,
		GasTipCap:  tx.GasTipCap,
		GasFeeCap:  tx.GasFeeCap,
		Gas:        tx.Gas,
		Value:      tx.Value,
		To:        &to,
		Data:       data,
		N:          tx.N,
		L:          l.Bytes(),
		Π:          π.Bytes(),

		V: tx.V,
		R: tx.R,
		S: tx.S,
	}), nil
}

func (t *Transaction) Reencrypt() (*Transaction, error) {
	tx, ok := t.inner.(*DecryptedTx)
	if !ok {
		return nil, errors.New("cannot reencrypt a non-decrypted transaction")
	}

	signer := NewLausanneSigner(t.ChainId())
	from, err := Sender(signer, t)
	if err != nil {
		return nil, err
	}

	label := binary.BigEndian.AppendUint64(from.Bytes(), tx.Nonce)
	n := new(big.Int).SetBytes(tx.N)
	l := new(big.Int).SetBytes(tx.L)
	π := new(big.Int).SetBytes(tx.Π)
	secret, ok := vdf.RecoverSecretFromProof(label, l, π, n, Log2t)
	if !ok {
		return nil, errors.New("bad VDF proof")
	}

	seed := secret.Bytes()

	plaintext := append(tx.To.Bytes(), tx.Data...)

	ciphertext := make([]byte, len(plaintext))
	tag := make([]byte, cae.Selected.TagLen())
	err = cae.Selected.Encrypt(ciphertext, tag, seed, plaintext)
	if err != nil {
		return nil, err
	}

	return NewTx(&EncryptedTx{
		ChainID:    tx.ChainID,
		Nonce:      tx.Nonce,
		GasTipCap:  tx.GasTipCap,
		GasFeeCap:  tx.GasFeeCap,
		Gas:        tx.Gas,
		Value:      tx.Value,
		Ciphertext: ciphertext,
		Tag:        tag,
		N: 	    tx.N,

		V: tx.V,
		R: tx.R,
		S: tx.S,
	}), nil
}
