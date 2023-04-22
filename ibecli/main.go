package main

import (
	"encoding/binary"
	"log"

	"github.com/drand/kyber"
	bls "github.com/drand/kyber-bls12381"
	"github.com/drand/kyber/encrypt/ibe"
	"github.com/drand/kyber/pairing"
	"github.com/drand/kyber/share"
	"github.com/drand/kyber/sign"
	"github.com/drand/kyber/sign/tbls"
	"github.com/drand/kyber/util/random"
)

const (
	THRESHOLD = 3
	N_SHARES  = 5
)

// Network simulates a threshold network by holding the whole private key
type Network struct {
	pairing.Suite
	KeyGroup kyber.Group
	sign.ThresholdScheme
	priShares []*share.PriShare
	pubPoly   *share.PubPoly
}

func NewNetwork() (*Network, error) {
	suite := bls.NewBLS12381Suite()
	thresholdScheme := tbls.NewThresholdSchemeOnG2(suite)
	keyGroup := suite.G1()
	// ref: https://github.com/drand/kyber/blob/master/sign/test/threshold.go#L14
	priPoly := share.NewPriPoly(keyGroup, THRESHOLD, nil, random.New())
	pubPoly := priPoly.Commit(keyGroup.Point().Base())
	shares := priPoly.Shares(N_SHARES)

	return &Network{suite, keyGroup, thresholdScheme, shares, pubPoly}, nil
}

func (network *Network) PublicKey() kyber.Point {
	return network.pubPoly.Commit()
}

func (network *Network) LabelForRound(rn uint64) []byte {
	buf := []byte("my cool blockchain")
	binary.BigEndian.AppendUint64(buf, rn)
	return buf
}
func (network *Network) SignRound(rn uint64) ([]byte, error) {
	msg := network.LabelForRound(rn)
	sigShares := make([][]byte, THRESHOLD)
	for i := range sigShares {
		var err error
		sigShares[i], err = network.ThresholdScheme.Sign(network.priShares[i], msg)
		if err != nil {
			return nil, err
		}
	}
	sig, err := network.ThresholdScheme.Recover(network.pubPoly, msg, sigShares, THRESHOLD, N_SHARES)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

func main() {
	network, err := NewNetwork()
	if err != nil {
		log.Fatal(err)
	}

	const rn uint64 = 1337
	sig, err := network.SignRound(rn)
	if err != nil {
		log.Fatal(err)
	}

	simKey := make([]byte, 32)
	random.Bytes(simKey, random.New())
	log.Printf("generated key %x", simKey)
	id := network.LabelForRound(rn)
	ct, err := ibe.EncryptCPAonG1(network.Suite, network.KeyGroup.Point().Base(), network.PublicKey(), id, simKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("encrypted key: %v", ct)

	// decrypt
	var signature bls.KyberG2
	if err := signature.UnmarshalBinary(sig); err != nil {
		log.Fatal(err)
	}
	simKey, err = ibe.DecryptCPAonG1(bls.NewBLS12381Suite(), &signature, ct)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("decrypted key %x", simKey)
}
