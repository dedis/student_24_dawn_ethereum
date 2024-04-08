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
	U := f3b.Suite.G2().Point()
	err = U.UnmarshalBinary(tx.EncKey)
	if err != nil {
		return nil, err
	}
	identityBytes, err := dkgcli.Extract(label)
	if err != nil {
		return nil, err
	}

	identity := f3b.Suite.G1().Point()
	err = identity.UnmarshalBinary(identityBytes)
	if err != nil {
		return nil, err
	}

	pk, err := dkgcli.GetPublicKey()
	if err != nil {
		return nil, err
	}

	ok = f3b.VerifyIdentity(pk, identity, label)
	if !ok {
		return nil, errors.New("bad identity")
	}

	secret := f3b.RecoverSecret(identity, U)
	seed, err := secret.MarshalBinary()
	if err != nil {
		return nil, err
	}
	log.Info("Decrypt()", "label", label, "U", U, "secret", secret, "seed", seed)

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
		Reveal:     identityBytes,

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

	U := f3b.Suite.G2().Point()
	err := U.UnmarshalBinary(tx.EncKey)
	if err != nil {
		return nil, err
	}

	identity := f3b.Suite.G1().Point()
	err = identity.UnmarshalBinary(tx.Reveal)
	if err != nil {
		return nil, err
	}

	secret := f3b.RecoverSecret(identity, U)
	seed, err := secret.MarshalBinary()
	if err != nil {
		return nil, err
	}

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
