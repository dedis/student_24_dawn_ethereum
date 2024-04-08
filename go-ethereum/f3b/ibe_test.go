// Copyright EPFL DEDIS

package f3b_test

import (
	"fmt"
	"testing"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/sign/bls"
	"go.dedis.ch/kyber/v3/util/random"

	"github.com/ethereum/go-ethereum/f3b"
)

func TestVerifyIdentity(t *testing.T) {
	sk, pk := bls.NewKeyPair(f3b.Suite, random.New())
	label := []byte("test")
	sig, err := bls.Sign(f3b.Suite, sk, label)
	if err != nil {
		t.Fatal("error", err)
	}
	s := f3b.Suite.G1().Point()
	err = s.UnmarshalBinary(sig)
	if err != nil {
		t.Fatal("error", err)
	}
	ok := f3b.VerifyIdentity(pk, s, label)
	if err != nil {
		t.Fatal("error", err)
	}
	if !ok {
		t.Fatal("bad identity")
	}
}

func BenchmarkVerifyIdentity(b *testing.B) {
	for _, variant := range []struct{name string; f func(pk, sigma kyber.Point, label []byte) bool}{{"Slow", f3b.VerifyIdentitySlow}, {"Fast", f3b.VerifyIdentityFast}} {
		verify := variant.f
		b.Run(fmt.Sprintf("Variant=%s", variant.name), func(b *testing.B) {
			sk, pk := bls.NewKeyPair(f3b.Suite, random.New())
			label := []byte("test")
			sig, err := bls.Sign(f3b.Suite, sk, label)
			if err != nil {
				b.Fatal("error", err)
			}
			s := f3b.Suite.G1().Point()
			err = s.UnmarshalBinary(sig)
			if err != nil {
				b.Fatal("error", err)
			}
			var ok bool
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ok = verify(pk, s, label)
			}
			b.StopTimer()
			if !ok {
				b.Fatal("bad identity")
			}
		})
	}
}

func BenchmarkRecoverSecret(b *testing.B) {
	sk, pk := bls.NewKeyPair(f3b.Suite, random.New())
	label := []byte("test")
	U, secret := f3b.ShareSecret(pk, label)
	sig, err := bls.Sign(f3b.Suite, sk, label)
	if err != nil {
		b.Fatal("error", err)
	}
	s := f3b.Suite.G1().Point()
	err = s.UnmarshalBinary(sig)
	if err != nil {
		b.Fatal("error", err)
	}
	b.ResetTimer()
	var rsecret kyber.Point
	for i := 0; i < b.N; i++ {
		rsecret = f3b.RecoverSecret(s, U)
	}
	b.StopTimer()
	if !secret.Equal(rsecret) {
		b.Fatal("secret mismatch")
	}
}
