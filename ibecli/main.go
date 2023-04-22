package main

import (
	"encoding/binary"
	"log"

	kyber "go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing"
	bn256 "go.dedis.ch/kyber/v3/pairing/bn256"
	"go.dedis.ch/kyber/v3/share"
	"go.dedis.ch/kyber/v3/sign/bls"
	"go.dedis.ch/kyber/v3/sign/tbls"
	"go.dedis.ch/kyber/v3/util/random"
)

const (
	THRESHOLD = 3
	N_SHARES  = 5
)

// Network simulates a threshold network by holding the whole private key
type Network struct {
	pairing.Suite
	priShares []*share.PriShare
	pubPoly   *share.PubPoly
}

func NewNetwork() (*Network, error) {
	suite := bn256.NewSuite()
	keyGroup := suite.G2()
	// ref: https://github.com/drand/kyber/blob/master/sign/test/threshold.go#L14
	priPoly := share.NewPriPoly(keyGroup, THRESHOLD, nil, random.New())
	pubPoly := priPoly.Commit(keyGroup.Point().Base())
	shares := priPoly.Shares(N_SHARES)

	return &Network{suite, shares, pubPoly}, nil
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
		sigShares[i], err = tbls.Sign(network.Suite, network.priShares[i], msg)
		if err != nil {
			return nil, err
		}
	}
	sig, err := tbls.Recover(network.Suite, network.pubPoly, msg, sigShares, THRESHOLD, N_SHARES)
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

	err = bls.Verify(network.Suite, network.PublicKey(), id, sig)
	if err != nil {
		log.Fatal(err)
	}
	ct, err := EncryptCPAonG2(network.Suite, network.PublicKey(), id, simKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("encrypted key: %v", ct)

	// decrypt
	signature := network.Suite.G1().Point()
	if err := signature.UnmarshalBinary(sig); err != nil {
		log.Fatal(err)
	}
	simKey, err = DecryptCPAonG2(network.Suite, signature, ct)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("decrypted key %x", simKey)
}
