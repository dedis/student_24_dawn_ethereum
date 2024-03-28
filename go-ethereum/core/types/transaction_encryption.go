// Copyright EPFL DEDIS

package types

import (
	"errors"
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/cae"
	"github.com/ethereum/go-ethereum/f3b"
	"github.com/ethereum/go-ethereum/log"
)

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

	dkgcli := f3b.NewDkgCli()

	label := binary.BigEndian.AppendUint64(from.Bytes(), tx.Nonce)
	key, err := dkgcli.Decrypt(label, tx.EncKey)
	if err != nil {
		return nil, err
	}

	// TODO: if the ciphertext is too short, penalize the sender
	plaintext := make([]byte, len(tx.Ciphertext))
	err = cae.Selected.Decrypt(plaintext, key, tx.Ciphertext, tx.Tag)
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
		Key:     key,
		EncKey:     tx.EncKey,

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

	plaintext := append(tx.To.Bytes(), tx.Data...)

	ciphertext := make([]byte, len(plaintext))
	tag := make([]byte, cae.Selected.TagLen())
	err := cae.Selected.Encrypt(ciphertext, tag, tx.Key, plaintext)
	if err != nil {
		panic(err)
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
		EncKey:     tx.EncKey,

		V: tx.V,
		R: tx.R,
		S: tx.S,
	}), nil
}
