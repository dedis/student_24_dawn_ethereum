// Copyright EPFL DEDIS

package f3b

type Protocol interface {
	ShareSecret(label []byte) (seed, encKey []byte, err error)
	RevealSecret(label, encKey []byte) (reveal []byte, err error)
	RecoverSecret(encKey, reveal []byte) (seed []byte, err error)
	IsVdf() bool
}
