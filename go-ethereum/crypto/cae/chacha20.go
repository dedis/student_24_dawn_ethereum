package cae

import (
	"golang.org/x/crypto/chacha20"
)


type _chacha20 struct {}

var Chacha20 = _chacha20{}

func (_chacha20) Name() string {
	return "chacha20"
}

func (_chacha20) TagLen() int {
	return 0
}

func (_chacha20) Encrypt(ciphertext, tag, key, plaintext []byte) error {
	cipher_key, _ := kdf(key, chacha20.KeySize, 0)

	//NOTE: nul iv! this is ok because the key is single use
	var iv [chacha20.NonceSize]byte
	cipher, err := chacha20.NewUnauthenticatedCipher(cipher_key, iv[:])
	if err != nil {
		return err
	}
	cipher.XORKeyStream(ciphertext, plaintext)

	return nil
}

func (_chacha20) Decrypt(plaintext, key, ciphertext, tag []byte) error {
	cipher_key, _ := kdf(key, chacha20.KeySize, 0)

	var iv [chacha20.NonceSize]byte
	cipher, err := chacha20.NewUnauthenticatedCipher(cipher_key, iv[:])
	if err != nil {
		return err
	}

	cipher.XORKeyStream(plaintext, ciphertext)
	return nil
}
