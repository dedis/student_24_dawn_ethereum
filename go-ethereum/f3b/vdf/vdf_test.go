// Copyright EPFL DEDIS

package vdf

import (
	"fmt"
	"testing"
)

func TestRecoverSecret(t *testing.T) {
	label := []byte("test")
	n, secret := ShareSecret(label)
	rsecret := RecoverSecret(n, label)
	if secret.Cmp(rsecret) != 0 {
		t.Fatal("bad recovered secret")
	}
}

func TestRecoverSecretWithProof(t *testing.T) {
	label := []byte("test")
	var steps uint64 = 1000
	n, secret := shareSecret(label, steps)
	y, π, l := recoverSecretWithProof(n, label, uint64(steps))
	g := deriveInitial(n, label)
	if !checkProof(g, y, π, l, n, uint64(steps)) {
		t.Fatal("bad proof")
	
	}
	fmt.Println("secret", secret)
	fmt.Println("y", y)
	if secret.Cmp(y) != 0 {
		t.Fatal("bad recovered secret")
	}
}

func BenchmarkVdfRecoverSecret(b *testing.B) {
	label := []byte("test")
	n, _ := ShareSecret(label)

	for steps := uint64(100); steps <= 1_000_000; steps *= 10 {
		b.Run(fmt.Sprintf("Steps=%d", steps), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				recoverSecret(n, label, steps)
			}
		})
	}
}
