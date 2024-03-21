// Copyright EPFL DEDIS

package types

import (
	"errors"
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/cae"
	"github.com/ethereum/go-ethereum/f3b"
)

func (t *Transaction) Decrypt() (*Transaction, error) {
	// Minimal signer for an encrypted transaction
	signer := NewLondonSigner(t.ChainId())

	from, err := signer.Sender(t)
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

	ciphertext := tx.Payload
	// TODO: if the ciphertext is too short, penalize the sender

	plaintext, err := cae.Selected.Decrypt(key, ciphertext)
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

		V: tx.V,
		R: tx.R,
		S: tx.S,
	}), nil
}

func (t *Transaction) Reencrypt() ([]byte, error) {
	// Minimal signer for an encrypted transaction
	signer := NewLondonSigner(t.ChainId())

	from, err := signer.Sender(t)
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

	ciphertext := tx.Payload
	// TODO: if the ciphertext is too short, penalize the sender

	plaintext, err := cae.Selected.Decrypt(key, ciphertext)
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

		V: tx.V,
		R: tx.R,
		S: tx.S,
	}), nil
}
