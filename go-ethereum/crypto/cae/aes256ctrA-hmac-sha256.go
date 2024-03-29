package cae

import (
	"crypto/sha256"
	"crypto/hmac"
	"crypto/aes"
	"crypto/cipher"
)


type aes256CtrHmacSha256 struct {}

var Aes256CtrHmacSha256 = aes256CtrHmacSha256{}

func (aes256CtrHmacSha256) Name() string {
	return "aes256ctr-hmac-sha256"
}

func (aes256CtrHmacSha256) TagLen() int {
	return 32
}

func (aes256CtrHmacSha256) Encrypt(ciphertext, tag, key, plaintext []byte) error {
	cipher_key, mac_key := kdf(key, 32, 32)

	//NOTE: nul iv! this is ok because the key is single use
	var iv [16]byte

	block, err := aes.NewCipher(cipher_key)
	if err != nil {
		return err
	}
	cipher := cipher.NewCTR(block, iv[:])
	cipher.XORKeyStream(ciphertext, plaintext)

	mac := hmac.New(sha256.New, mac_key)
	mac.Write(ciphertext)
	mac.Sum(tag[:0])
	return nil
}

func (aes256CtrHmacSha256) Decrypt(plaintext, key, ciphertext, tag []byte) error {
	cipher_key, mac_key := kdf(key, 32, 32)

	var buf [32]byte
	mac := hmac.New(sha256.New, mac_key)
	mac.Write(ciphertext)
	mac.Sum(buf[:0])
	if !hmac.Equal(tag, buf[:]) {
		return AuthenticationError
	}

	var iv [16]byte
	block, err := aes.NewCipher(cipher_key)
	if err != nil {
		return err
	}
	cipher := cipher.NewCTR(block, iv[:])
	cipher.XORKeyStream(plaintext, ciphertext)
	return nil
}
