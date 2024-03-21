// Copyright EPFL DEDIS

package cae

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto/cae"
)

func TestHappyPath(t *testing.T) {
	scheme := cae.Selected
	key := []byte{1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,}
	ciphertext, err := scheme.Encrypt(key, []byte("hello"))
	if err != nil {
		t.Fatal(err)
	}

	plaintext, err := scheme.Decrypt(key, ciphertext)
	if err != nil {
		t.Fatal(err)
	}

	if string(plaintext) != "hello" {
		t.Fatalf("%x != hello", plaintext)
	}
}
