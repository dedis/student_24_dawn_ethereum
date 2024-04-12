// Copyright EPFL DEDIS

package vdf

import (
	"fmt"
	"math/big"
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
	n, secret := ShareSecret(label)
	steps := 1000
	y, π, l := recoverSecretWithProof(n, label, uint64(steps))
	init := deriveInitial(n, label)
	two := new(big.Int).SetInt64(2)
	y2 := new(big.Int).Exp(two, new(big.Int).SetInt64(int64(steps)), l)
	y2.Exp(init, y2, n)
	y2.Mul(y2, π.Exp(π, l, n))
	if y2.Cmp(y) != 0 {
		t.Fatal("bad proof")
	}
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
