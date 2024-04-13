// Copyright EPFL DEDIS

package vdf

import (
	"fmt"
	"testing"
)

func TestRecoverSecret(t *testing.T) {
	const steps uint64 = 1000
	label := []byte("test")
	n, secret := ShareSecret(label, steps)
	rsecret := RecoverSecret(label, n, steps)
	if secret.Cmp(rsecret) != 0 {
		t.Fatal("bad recovered secret")
	}
}

func TestRecoverSecretFromProof(t *testing.T) {
	const steps uint64 = 1000
	label := []byte("test")
	n, secret := ShareSecret(label, steps)
	l, π := Proof(label, n, steps)
	y, ok := RecoverSecretFromProof(label, l, π, n, steps) 
	if !ok {
		t.Fatal("bad proof")
	
	}
	if secret.Cmp(y) != 0 {
		t.Fatal("bad recovered secret")
	}
}

func BenchmarkRecoverSecret(b *testing.B) {
	for steps := uint64(100); steps <= 1_000_000; steps *= 10 {
		b.Run(fmt.Sprintf("Steps=%d", steps), func(b *testing.B) {
			label := []byte("test")
			n, _ := ShareSecret(label, steps)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				RecoverSecret(label, n, steps)
			}
		})
	}
}

func BenchmarkRecoverSecretFromProof(b *testing.B) {
	for steps := uint64(100); steps <= 1_000_000; steps *= 10 {
		b.Run(fmt.Sprintf("Steps=%d", steps), func(b *testing.B) {
			label := []byte("test")
			n, _ := ShareSecret(label, steps)
			l, π := Proof(label, n, steps)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				RecoverSecretFromProof(label, l, π, n, steps)
			}
		})
	}
}

func BenchmarkProof(b *testing.B) {
	for steps := uint64(100); steps <= 1_000_000; steps *= 10 {
		b.Run(fmt.Sprintf("Steps=%d", steps), func(b *testing.B) {
			label := []byte("test")
			n, _ := ShareSecret(label, steps)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = Proof(label, n, steps)
			}
		})
	}
}
