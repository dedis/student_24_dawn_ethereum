// Copyright EPFL DEDIS

package cae

import (
	"testing"
)

func TestHappyPath(t *testing.T) {
	scheme := ChaCha20HmacSha256{}
	plaintext := []byte("hello")
	ciphertext := make([]byte, len(plaintext))
	tag := make([]byte, scheme.TagLen())
	key := []byte{1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,}
	err := scheme.Encrypt(ciphertext, tag, key, plaintext)
	if err != nil {
		t.Fatal(err)
	}

	plaintext = make([]byte, len(ciphertext))
	err = scheme.Decrypt(plaintext, key, ciphertext, tag)
	if err != nil {
		t.Fatal(err)
	}

	if string(plaintext) != "hello" {
		t.Fatalf("%x != hello", plaintext)
	}
}

func TestMac(t *testing.T) {
	scheme := ChaCha20HmacSha256{}
	plaintext := []byte("hello")
	ciphertext := make([]byte, len(plaintext))
	tag := make([]byte, scheme.TagLen())
	key := []byte{1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,}
	err := scheme.Encrypt(ciphertext, tag, key, plaintext)
	if err != nil {
		t.Fatal(err)
	}
	tag[2] ^= 1

	plaintext = make([]byte, len(ciphertext))
	err = scheme.Decrypt(plaintext, key, ciphertext, tag)
	if err != AuthenticationError {
		t.Fatalf("expected cae.AuthenticationError, got %v", err)
	}
}

func BenchmarkDecrypt(b *testing.B) {
	scheme := ChaCha20HmacSha256{}
	plaintext := []byte("hellohellohellohello")
	ciphertext := make([]byte, len(plaintext))
	tag := make([]byte, scheme.TagLen())
	key := []byte{1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,}
	err := scheme.Encrypt(ciphertext, tag, key, plaintext)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scheme.Decrypt(plaintext, key, ciphertext, tag)
	}
}
