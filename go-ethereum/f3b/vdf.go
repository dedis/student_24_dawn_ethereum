package f3b

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/f3b/vdf"
)

const Log2t = 15

type VDF struct { }

func (_ *VDF) ShareSecret(label []byte) (seed, encKey []byte, err error) {
	secret, _, _, n := vdf.ShareSecret(label, Log2t)

	encKey = n.Bytes()
	seed = secret.Bytes()
	return
}

const lBytes = 32
const πBytes = 512

func (_ *VDF) RevealSecret(label, encKey []byte) (reveal []byte, err error) {
	n := new(big.Int).SetBytes(encKey)
	l, π := vdf.Proof(label, n, Log2t)

	reveal = make([]byte, lBytes+πBytes)
	l.FillBytes(reveal[:lBytes])
	π.FillBytes(reveal[lBytes:])
	return
}

func (_ *VDF) RecoverSecret(label, encKey, reveal []byte) (seed []byte, err error) {
	n := new(big.Int).SetBytes(encKey)
	l := new(big.Int).SetBytes(reveal[:lBytes])
	π := new(big.Int).SetBytes(reveal[lBytes:])
	secret, ok := vdf.RecoverSecretFromProof(label, l, π, n, Log2t)
	if !ok {
		return nil, errors.New("bad VDF proof")
	}

	seed = secret.Bytes()
	return
}

func (_ *VDF) IsVdf() bool {
	return true
}
