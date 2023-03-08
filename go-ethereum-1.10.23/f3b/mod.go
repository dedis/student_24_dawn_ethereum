// Copyright EPFL DEDIS

package f3b

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"

	"golang.org/x/crypto/hkdf"
)

// EncryptCompact encrypts the plaintext using AES-CTR
// The same secret should NEVER be used to encrypt two different messages.
//
// secret can be of any length, and should ideally represent at least 128 bits of entropy.
//
// plaintext will have the same length as ciphertext.
func EncryptCompact(secret []byte, plaintext []byte) (ciphertext []byte) {
	ciphertext = xorWithKeyStream(secret, plaintext)
	return
}

func DecryptCompact(secret []byte, ciphertext []byte) (plaintext []byte) {
	plaintext = xorWithKeyStream(secret, ciphertext)
	return
}

// HKDF personalisation string
const compactInfo = "DEDIS-F3B-Compact-AES256CTR"

func xorWithKeyStream(secret []byte, input []byte) (output []byte) {
	var key [32]byte
	var iv [16]byte

	kdf := hkdf.New(sha512.New, secret, nil, []byte(compactInfo))
	kdf.Read(key[:])
	kdf.Read(iv[:])
	aes256, err := aes.NewCipher(key[:])
	if err != nil {
		panic(err)
	}
	stream := cipher.NewCTR(aes256, iv[:])
	output = make([]byte, len(input))
	stream.XORKeyStream(output, input)
	return
}
