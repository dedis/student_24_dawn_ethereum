package main

import (
	"log"

	"github.com/drand/drand/chain"
	drand_crypto "github.com/drand/drand/crypto"

	"github.com/drand/kyber"
	"github.com/drand/kyber/sign"
	"github.com/drand/kyber/sign/bls"
	"github.com/drand/kyber/util/random"
	_ "github.com/drand/tlock"
)

// note: default scheme is not compatible with timelock encryption
const SchemeID = drand_crypto.ShortSigSchemeID

// Network simulates a threshold network by holding the whole private key
type Network struct {
	*drand_crypto.Scheme
	sigScheme sign.Scheme
	priv      kyber.Scalar
	pub       kyber.Point
}

/*
	func NewNetwork() (*Network, error) {
		scheme, err := drand_crypto.SchemeFromName(SchemeID)
		if err != nil {
			return nil, err
		}
		random, err := chacha20.NewUnauthenticatedCipher(make([]byte, 32), []byte("realrandomiv"))
		if err != nil {
			return nil, err
		}
		// ref: https://github.com/drand/kyber/blob/master/sign/test/threshold.go#L14
		const threshold = 1
		priPoly := share.NewPriPoly(scheme.KeyGroup, threshold, nil, random)
		pubPoly := priPoly.Commit(scheme.KeyGroup.Point().Base())
		priv := priPoly.Shares(1)[0]
		pub := pubPoly.Commit()
		return &Network{priv, pub}, nil
	}
*/
func NewNetwork() (*Network, error) {
	scheme, err := drand_crypto.SchemeFromName(SchemeID)
	if err != nil {
		return nil, err
	}
	sigScheme := bls.NewSchemeOnG1(scheme.Pairing)
	priv, pub := sigScheme.NewKeyPair(random.New())
	return &Network{scheme, sigScheme, priv, pub}, nil
}

func (b *Network) PublicKey() kyber.Point {
	return b.pub
}

func (b *Network) SignRound(rn uint64) (*chain.Beacon, error) {
	sig, err := b.sigScheme.Sign(b.priv, b.DigestBeacon(&chain.Beacon{Round: rn}))
	if err != nil {
		return nil, err
	}
	return &chain.Beacon{Round: rn, Signature: sig}, nil
}

func main() {
	network, err := NewNetwork()
	if err != nil {
		log.Fatal(err)
	}

	const rn uint64 = 1337
	b, err := network.SignRound(rn)
	if err != nil {
		log.Fatal(err)
	}

	err = network.VerifyBeacon(b, network.PublicKey())
	if err != nil {
		log.Fatal(err)
	}
}
