// Copyright EPFL DEDIS

package f3b

import (
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/sign/bdn"
	"go.dedis.ch/kyber/v3/util/random"

	"github.com/ethereum/go-ethereum/f3b/ibe"
)

type fakeSmcCli struct {
	pubkey kyber.Point
	privkey kyber.Scalar
}

// Create a fake SMC that is just a local keypair, for testing purposes
func NewFakeSmcCli() SmcCli {
	privkey, pubkey := bdn.NewKeyPair(ibe.Suite, random.New())
	return &fakeSmcCli{ pubkey, privkey }
}

func (c *fakeSmcCli) GetPublicKey() (kyber.Point, error) {
	return c.pubkey.Clone(), nil
}

func (c *fakeSmcCli) Extract(label []byte) (v []byte, err error) {
	return bdn.Sign(ibe.Suite, c.privkey, label)
}
