// Copyright EPFL DEDIS

package cae_test

import (
	"testing"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto/cae"
)

func TestHappyPath(t *testing.T) {
	scheme := cae.Selected
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
	scheme := cae.Selected
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
	if err != cae.AuthenticationError {
		t.Fatalf("expected cae.AuthenticationError, got %v", err)
	}
}

func BenchmarkDecrypt(b *testing.B) {
	for _, scheme := range cae.AllSchemes {
		for l := 20; l <= 2_000_000; l *= 10 {
			b.Run(fmt.Sprintf("Scheme=%s/l=%d", scheme.Name(), l), func(b *testing.B) {
				plaintext := make([]byte, l)
				ciphertext := make([]byte, l)
				tag := make([]byte, scheme.TagLen())
				key := []byte{1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,}
				err := scheme.Encrypt(ciphertext, tag, key, plaintext)
				if err != nil {
					b.Fatal(err)
				}

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					err = scheme.Decrypt(plaintext, key, ciphertext, tag)
				}
				b.StopTimer()
				if err != nil {
					b.Fatal(err)
				}
				for i := range plaintext {
					if plaintext[i] != 0 {
						b.Fatalf("%x != 0", plaintext)
					}
				}
				// gascost for 10MGas/s
				// 10MGas = cost/elapsed
				// cost = 10MGas/s * elpased
				// cost = 10MGas/s * elpased / 10e9 μs/s
				suggestedGas := b.Elapsed().Nanoseconds() / 25 / int64(b.N)
				b.ReportMetric(float64(suggestedGas), "gas/op")
				b.ReportMetric(float64(suggestedGas)/float64(l), "gas/B")
			})
		}
	}
}
