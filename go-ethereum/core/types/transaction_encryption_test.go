// Copyright EPFL DEDIS

package types_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/f3b/vdf"
)

func TestHash(t *testing.T) {
	acct, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	label := append(crypto.PubkeyToAddress(acct.PublicKey).Bytes(), 0,0,0,0,0,0,0,0,)
	
	_, n := vdf.ShareSecret(label, types.Log2t)
	l, π := vdf.Proof(label, n, types.Log2t)
	dec_tx := types.NewTx(&types.DecryptedTx{
		ChainID:    big.NewInt(1),
		Nonce:      0,
		GasFeeCap:  big.NewInt(1),
		Gas:        100000,
		Value:      big.NewInt(0),
		To:         &common.Address{},
		Data:       []byte("hello world"),
		N:          n.Bytes(),
		L:          l.Bytes(),
		Π:          π.Bytes(),
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
	t.Skip("Skipping test")
	dec_tx := types.NewTx(&types.DecryptedTx{
		ChainID:    big.NewInt(1),
		Nonce:      0,
		GasFeeCap:  big.NewInt(1),
		Gas:        100000,
		Value:      big.NewInt(0),
		To:         &common.Address{},
		Data:       []byte("hello world"),
		N:          []byte("n"),
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
