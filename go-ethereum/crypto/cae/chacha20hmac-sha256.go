package cae

import (
	"crypto/sha256"
	"crypto/hmac"

	"golang.org/x/crypto/chacha20"
)


type chacha20HmacSha256 struct {}

var Chacha20HmacSha256 = chacha20HmacSha256{}

func (chacha20HmacSha256) Name() string {
	return "chacha20-hmac-sha256"
}

func (chacha20HmacSha256) TagLen() int {
	return 32
}

func (chacha20HmacSha256) Encrypt(ciphertext, tag, key, plaintext []byte) error {
	cipher_key, mac_key := kdf(key, chacha20.KeySize, 32)

	//NOTE: nul iv! this is ok because the key is single use
	var iv [chacha20.NonceSize]byte
	cipher, err := chacha20.NewUnauthenticatedCipher(cipher_key, iv[:])
	if err != nil {
		return err
	}
	cipher.XORKeyStream(ciphertext, plaintext)

	mac := hmac.New(sha256.New, mac_key)
	mac.Write(ciphertext)
	mac.Sum(tag[:0])
	return nil
}

func (chacha20HmacSha256) Decrypt(plaintext, key, ciphertext, tag []byte) error {
	cipher_key, mac_key := kdf(key, chacha20.KeySize, 32)

	var buf [32]byte
	mac := hmac.New(sha256.New, mac_key)
	mac.Write(ciphertext)
	mac.Sum(buf[:0])
	if !hmac.Equal(tag, buf[:]) {
		return AuthenticationError
	}

	var iv [chacha20.NonceSize]byte
	cipher, err := chacha20.NewUnauthenticatedCipher(cipher_key, iv[:])
	if err != nil {
		return err
	}

	cipher.XORKeyStream(plaintext, ciphertext)
	return nil
}
