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

func TestRecoverSecretWithProof(t *testing.T) {
	const steps uint64 = 1000
	label := []byte("test")
	n, secret := ShareSecret(label, steps)
	y, π, l := recoverSecretWithProof(label, n, steps)
	g := deriveInitial(label, n)
	if !checkProof(g, y, π, l, n, steps) {
		t.Fatal("bad proof")
	
	}
	if secret.Cmp(y) != 0 {
		t.Fatal("bad recovered secret")
	}
}

func BenchmarkVdfRecoverSecret(b *testing.B) {
	label := []byte("test")

	for steps := uint64(100); steps <= 1_000_000; steps *= 10 {
		b.Run(fmt.Sprintf("Steps=%d", steps), func(b *testing.B) {
			n, _ := ShareSecret(label, steps)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				RecoverSecret(label, n, steps)
			}
		})
	}
}
