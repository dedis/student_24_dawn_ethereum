// Copyright EPFL DEDIS

package f3b_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/f3b"
)

func TestRecoverSecret(t *testing.T) {
	label := []byte("test")
	n, secret := f3b.VdfShareSecret(label)
	rsecret := f3b.VdfRecoverSecret(n, label)
	if secret.Cmp(rsecret) != 0 {
		t.Fatal("bad recovered secret")
	}
}

func BenchmarkVdfRecoverSecret(b *testing.B) {
	label := []byte("test")
	n, _ := f3b.VdfShareSecret(label)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f3b.VdfRecoverSecret(n, label)
	}
	b.StopTimer()
}
