// Copyright EPFL DEDIS

package types_test

import (
	"math/big"
	"encoding/binary"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/f3b"
)

func TestHashEncryptedTx(t *testing.T) {
	vdf := &f3b.VDF{Log2t: 5}
	f3b.ForceSelectedProtocol(t, vdf)

	acct, err := crypto.GenerateKey()
	from := crypto.PubkeyToAddress(acct.PublicKey)
	label := binary.BigEndian.AppendUint64(from.Bytes(), 0)
	_, encKey, err := vdf.ShareSecret(label)
	if err != nil {
		t.Error(err)
	}
	reveal, err := vdf.RevealSecret(label, encKey)
	if err != nil {
		t.Fatal(err)
	}

	dec_tx := types.NewTx(&types.DecryptedTx{
		ChainID:    big.NewInt(1),
		Nonce:      0,
		GasFeeCap:  big.NewInt(1),
		Gas:        100000,
		Value:      big.NewInt(0),
		To:         &common.Address{},
		Data:       []byte("hello world"),
		EncKey:     encKey,
		Reveal:     reveal,
		From:       from,
	})
	enc_tx, err := dec_tx.Reencrypt(vdf)
	if err != nil {
		t.Fatal(err)
	}
	dec_hash := dec_tx.Hash().String()
	enc_hash := enc_tx.Hash().String()
	if dec_hash != enc_hash {
		t.Fatalf("Expected %s to equal %s", dec_hash, enc_hash)
	}
}

func TestSignEncryptedTx(t *testing.T) {
	vdf := &f3b.VDF{Log2t: 5}
	f3b.ForceSelectedProtocol(t, vdf)

	acct, err := crypto.GenerateKey()
	from := crypto.PubkeyToAddress(acct.PublicKey)
	label := binary.BigEndian.AppendUint64(from.Bytes(), 0)
	_, encKey, err := vdf.ShareSecret(label)
	if err != nil {
		t.Error(err)
	}
	reveal, err := vdf.RevealSecret(label, encKey)
	if err != nil {
		t.Fatal(err)
	}

	dec_tx := types.NewTx(&types.DecryptedTx{
		ChainID:    big.NewInt(1),
		Nonce:      0,
		GasFeeCap:  big.NewInt(1),
		Gas:        100000,
		Value:      big.NewInt(0),
		To:         &common.Address{},
		Data:       []byte("hello world"),
		EncKey:     encKey,
		Reveal:     reveal,
		From:       from,
	})
	signer := types.NewLausanneSigner(big.NewInt(1))
	dec_tx, err = types.SignTx(dec_tx, signer, acct)
	if err != nil {
		t.Fatal(err)
	}
	enc_tx, err := dec_tx.Reencrypt(vdf)
	if err != nil {
		t.Fatal(err)
	}
	enc_sender, err := types.Sender(signer, enc_tx)
	if err != nil {
		t.Fatal(err)
	}
	dec_sender, err := types.Sender(signer, dec_tx)
	if err != nil {
		t.Fatal(err)
	}
	if enc_sender != dec_sender {
		t.Fatal("Expected sender to be the same")
	}
}
