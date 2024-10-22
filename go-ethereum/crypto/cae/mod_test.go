// Copyright EPFL DEDIS

package cae_test

import (
	"testing"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto/cae"
)

func TestHappyPath(t *testing.T) {
	for _, scheme := range cae.AllSchemes {
		t.Run("Scheme=" + scheme.Name(), func(t *testing.T) {
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
		})
	}
}

func TestMac(t *testing.T) {
	for _, scheme := range cae.AllSchemes {
		t.Run("Scheme=" + scheme.Name(), func(t *testing.T) {
			if scheme.TagLen() == 0 {
				// no MAC
				t.SkipNow()
			}
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
		})
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
				nsPerByte := float64(b.Elapsed().Nanoseconds() / int64(b.N)) / float64(l)
				b.ReportMetric(nsPerByte, "ns/B")
			})
		}
	}
}
