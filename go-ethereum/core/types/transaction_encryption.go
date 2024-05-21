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

func (t *Transaction) Encrypt(from common.Address, f3bProtocol f3b.Protocol) (*Transaction, error) {
	label := binary.BigEndian.AppendUint64(from.Bytes(), t.Nonce())
	seed, encKey, err := f3bProtocol.ShareSecret(label)
	if err != nil {
		return nil, err
	}

	log.Info("Encrypting", "label", label, "encKey", encKey, "seed", seed)

	plaintext := append(t.To().Bytes(), t.Data()...)
	ciphertext := make([]byte, len(plaintext))
	tag := make([]byte, cae.Selected.TagLen())
	err = cae.Selected.Encrypt(ciphertext, tag, seed, plaintext)
	if err != nil {
		return nil, err
	}

	return NewTx(&EncryptedTx{
		ChainID:    t.ChainId(),
		Nonce:      t.Nonce(),
		GasTipCap:  t.GasTipCap(),
		GasFeeCap:  t.GasFeeCap(),
		Gas:        t.Gas(),
		Value:      t.Value(),
		Ciphertext: ciphertext,
		Tag:        tag,
		EncKey:     encKey,
	}), nil
}
func (t *Transaction) Decrypt(f3bProtocol f3b.Protocol) (*Transaction, error) {
	// Minimal signer for an encrypted transaction
	signer := NewLausanneSigner(t.ChainId())

	from, err := Sender(signer, t)
	if err != nil {
		return nil, err
	}

	tx, ok := t.inner.(*EncryptedTx)
	if !ok {
		return nil, errors.New("cannot decrypt a non-encrypted transaction")
	}


	label := binary.BigEndian.AppendUint64(from.Bytes(), tx.Nonce)
	reveal, err := f3bProtocol.RevealSecret(label, tx.EncKey)
	if err != nil {
		return nil, err
	}
	seed, err := f3bProtocol.RecoverSecret(label, tx.EncKey, reveal)
	if err != nil {
		return nil, err
	}
	log.Info("Reencrypting", "label", label, "enckey", tx.EncKey, "reveal", reveal, "seed", seed)

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
		EncKey:     tx.EncKey,
		Reveal:     reveal,
		From:       from,

		V: tx.V,
		R: tx.R,
		S: tx.S,
	}), nil
}

func (t *Transaction) Reencrypt(protocol f3b.Protocol) (*Transaction, error) {
	tx, ok := t.inner.(*DecryptedTx)
	if !ok {
		return nil, errors.New("cannot reencrypt a non-decrypted transaction")
	}

	label := binary.BigEndian.AppendUint64(tx.From.Bytes(), tx.Nonce)
	seed, err := protocol.RecoverSecret(label, tx.EncKey, tx.Reveal)
	if err != nil {
		return nil, err
	}
	log.Info("Reencrypting", "label", label, "enckey", tx.EncKey, "reveal", tx.Reveal, "seed", seed)

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
		EncKey:     tx.EncKey,

		V: tx.V,
		R: tx.R,
		S: tx.S,
	}), nil
}
