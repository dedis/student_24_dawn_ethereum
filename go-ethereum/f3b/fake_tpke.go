// Copyright EPFL DEDIS

package f3b

import "golang.org/x/crypto/chacha20"

func NewFakeTPKE() (Protocol, error) {
	var key [chacha20.KeySize]byte
	var iv [chacha20.NonceSize]byte
	// deterministic randomness for a deterministic fake key
	rand, err := chacha20.NewUnauthenticatedCipher(key[:], iv[:])
	if err != nil {
		return nil, err
	}

	return NewTPKE(NewFakeSmcCli(rand))
}
