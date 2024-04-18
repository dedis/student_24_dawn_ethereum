// Copyright EPFL DEDIS

package vdf

import (
	"fmt"
	"testing"
)

func TestRecoverSecret(t *testing.T) {
	const log2t = 10
	label := []byte("test")
	secret, _, _, n := ShareSecret(label, log2t)
	rsecret := RecoverSecret(label, n, log2t)
	if secret.Cmp(rsecret) != 0 {
		t.Fatal("bad recovered secret")
	}
}

func TestRecoverSecretFromProof(t *testing.T) {
	const log2t = 10
	label := []byte("test")
	secret, l, π, n := ShareSecret(label, log2t)
	y, ok := RecoverSecretFromProof(label, l, π, n, log2t) 
	if !ok {
		t.Fatal("bad proof")
	
	}
	if secret.Cmp(y) != 0 {
		t.Fatal("bad recovered secret")
	}
}

func BenchmarkRecoverSecret(b *testing.B) {
	for log2t := 5; log2t <= 20; log2t++ {
		b.Run(fmt.Sprintf("log2t=%d", log2t), func(b *testing.B) {
			label := []byte("test")
			n, _, _, _ := ShareSecret(label, log2t)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				RecoverSecret(label, n, log2t)
			}
		})
	}
}

func BenchmarkRecoverSecretFromProof(b *testing.B) {
	for log2t := 5; log2t <= 50; log2t += 5 {
		b.Run(fmt.Sprintf("log2t=%d", log2t), func(b *testing.B) {
			label := []byte("test")
			n, l, π, _ := ShareSecret(label, log2t)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				RecoverSecretFromProof(label, l, π, n, log2t)
			}
		})
	}
}

func BenchmarkProof(b *testing.B) {
	for log2t := 5; log2t <= 20; log2t += 5 {
		b.Run(fmt.Sprintf("log2t=%d", log2t), func(b *testing.B) {
			label := []byte("test")
			n, _, _, _ := ShareSecret(label, log2t)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = Proof(label, n, log2t)
			}
		})
	}
}
