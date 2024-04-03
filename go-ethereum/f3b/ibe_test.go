// Copyright EPFL DEDIS

package f3b_test

import (
	"testing"
	"fmt"

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
	fmt.Println(sig)
	s := f3b.Suite.G1().Point()
	err = s.UnmarshalBinary(sig)
	if err != nil {
		t.Fatal("error", err)
	}
	ok, err := f3b.VerifyIdentity(pk, s, label)
	if err != nil {
		t.Fatal("error", err)
	}
	if !ok {
		t.Fatal("bad identity")
	}
}
