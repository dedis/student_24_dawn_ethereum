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
	"github.com/ethereum/go-ethereum/rlp"
)

const Log2t = 15

func (t *Transaction) Decrypt() (*Transaction, error) {
	// Minimal signer for an encrypted transaction
	signer := NewLausanneSigner(t.ChainId())

	from, err := Sender(signer, t)
	if err != nil {
		return nil, err
	}

	 _ = from

	tx, ok := t.inner.(*EncryptedTx)
	if !ok {
		return nil, errors.New("cannot decrypt a non-encrypted transaction")
	}

	label := binary.BigEndian.AppendUint64(from.Bytes(), tx.Nonce)
	n := new(big.Int).SetBytes(tx.EncKey)
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

	reveal, err := rlp.EncodeToBytes([][]byte{from.Bytes(), l.Bytes(), π.Bytes()})
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

	reveal := [][]byte{}
	rlp.DecodeBytes(tx.Reveal, &reveal)
	from := reveal[0]
	l := new(big.Int).SetBytes(reveal[1])
	π := new(big.Int).SetBytes(reveal[2])

	label := binary.BigEndian.AppendUint64(from, tx.Nonce)

	n := new(big.Int).SetBytes(tx.EncKey)
	secret, ok := vdf.RecoverSecretFromProof(label, l, π, n, Log2t)
	if !ok {
		log.Error("bad vdf proof", "label", label, "l", l, "π", π, "n", n, "Log2t", Log2t)
		return nil, errors.New("bad VDF proof")
	}

	seed := secret.Bytes()

	plaintext := append(tx.To.Bytes(), tx.Data...)

	ciphertext := make([]byte, len(plaintext))
	tag := make([]byte, cae.Selected.TagLen())
	err := cae.Selected.Encrypt(ciphertext, tag, seed, plaintext)
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
