package main

import (
	"log"

	"github.com/drand/drand/chain"
	drand_crypto "github.com/drand/drand/crypto"

	"github.com/drand/kyber"
	bls "github.com/drand/kyber-bls12381"
	"github.com/drand/kyber/encrypt/ibe"
	"github.com/drand/kyber/share"
	"github.com/drand/kyber/util/random"
)

// note: default scheme is not compatible with timelock encryption
const SchemeID = drand_crypto.UnchainedSchemeID

const (
	THRESHOLD = 3
	N_SHARES  = 5
)

// Network simulates a threshold network by holding the whole private key
type Network struct {
	*drand_crypto.Scheme
	priShares []*share.PriShare
	pubPoly   *share.PubPoly
}

func NewNetwork() (*Network, error) {
	scheme, err := drand_crypto.SchemeFromName(SchemeID)
	if err != nil {
		return nil, err
	}

	// ref: https://github.com/drand/kyber/blob/master/sign/test/threshold.go#L14
	priPoly := share.NewPriPoly(scheme.KeyGroup, THRESHOLD, nil, random.New())
	pubPoly := priPoly.Commit(scheme.KeyGroup.Point().Base())
	shares := priPoly.Shares(N_SHARES)

	return &Network{scheme, shares, pubPoly}, nil
}

func (b *Network) PublicKey() kyber.Point {
	return b.pubPoly.Commit()
}

func (b *Network) SignRound(rn uint64) (*chain.Beacon, error) {
	msg := b.DigestBeacon(&chain.Beacon{Round: rn})
	sigShares := make([][]byte, THRESHOLD)
	for i := range sigShares {
		var err error
		sigShares[i], err = b.ThresholdScheme.Sign(b.priShares[i], msg)
		if err != nil {
			return nil, err
		}
	}
	sig, err := b.ThresholdScheme.Recover(b.pubPoly, msg, sigShares, THRESHOLD, N_SHARES)
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

	simKey := make([]byte, 32)
	random.Bytes(simKey, random.New())
	log.Printf("generated key %x", simKey)
	id := network.DigestBeacon(&chain.Beacon{Round: rn})
	ct, err := ibe.EncryptCPAonG1(bls.NewBLS12381Suite(), network.Scheme.KeyGroup.Point().Base(), network.PublicKey(), id, simKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("encrypted key: %v", ct)

	// decrypt
	var signature bls.KyberG2
	if err := signature.UnmarshalBinary(b.Signature); err != nil {
		log.Fatal(err)
	}
	simKey, err = ibe.DecryptCPAonG1(bls.NewBLS12381Suite(), &signature, ct)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("decrypted key %x", simKey)
}
