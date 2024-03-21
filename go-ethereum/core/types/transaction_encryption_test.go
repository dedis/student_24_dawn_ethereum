// Copyright EPFL DEDIS

package types

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestHash(t *testing.T) {
	enc_key := []byte{1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,}
	key :=  []byte{2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,}
	dec_tx := types.NewTx(&types.DecryptedTx{
		ChainID:    big.NewInt(1),
		Nonce:      0,
		GasFeeCap:  big.NewInt(1),
		Gas:        100000,
		Value:      big.NewInt(0),
		To:         &common.Address{},
		Data:       []byte("hello world"),
		EncKey:     enc_key,
		Key:        key,
	})
	dec_hash := dec_tx.Hash().String()
	enc_tx, err := dec_tx.Reencrypt()
	if err != nil {
		t.Fatal(err)
	}
	enc_hash := enc_tx.Hash().String()
	if dec_hash != enc_hash {
		t.Fatalf("Expected %s to equal %s", dec_hash, enc_hash)
	}
}

func TestSig(t *testing.T) {
	enc_key := []byte{1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,}
	key :=  []byte{2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,}
	dec_tx := types.NewTx(&types.DecryptedTx{
		ChainID:    big.NewInt(1),
		Nonce:      0,
		GasFeeCap:  big.NewInt(1),
		Gas:        100000,
		Value:      big.NewInt(0),
		To:         &common.Address{},
		Data:       []byte("hello world"),
		EncKey:     enc_key,
		Key:        key,
	})
	signer := types.NewLausanneSigner(big.NewInt(1))
	acct, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	dec_tx, err = types.SignTx(dec_tx, signer, acct)
	if err != nil {
		t.Fatal(err)
	}
	enc_tx, err := dec_tx.Reencrypt()
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
	fmt.Println(enc_sender, dec_sender)
	if enc_sender != dec_sender {
		t.Fatal("Expected sender to be the same")
	}
}
