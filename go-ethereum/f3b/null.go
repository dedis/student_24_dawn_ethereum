// Copyright EPFL DEDIS

package f3b

import "bytes"

type Null struct {
}

func NewNull() Protocol {
	return &Null{}
}

func (_ *Null) ShareSecret(label []byte) (seed, encKey []byte, err error) {
	seed = bytes.Repeat([]byte{1}, 32)
	encKey = bytes.Repeat([]byte{0xfe}, 32)
	return
}

func (_ *Null) RevealSecret(labelBytes []byte, encKey []byte) (reveal []byte, err error) {
	reveal = bytes.Repeat([]byte{0xaa}, 32)
	return
}

func (_ *Null) RecoverSecret(labelBytes []byte, encKey, reveal []byte) (seed []byte, err error) {
	seed = bytes.Repeat([]byte{1}, 32)
	return
}

func (_ *Null) IsVdf() bool {
	return false
}

func (_ *Null) IsTibe() bool {
	return false
}

func (_ *Null) IsTpke() bool {
	return false
}
