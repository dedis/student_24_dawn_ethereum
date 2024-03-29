package cae

import (
	"crypto/aes"
	"crypto/cipher"
)


type aes256Gcm struct {}

var Aes256Gcm = aes256Gcm{}

func (aes256Gcm) Name() string {
	return "aes256-gcm"
}

func (aes256Gcm) TagLen() int {
	return 16
}

func (aes256Gcm) Encrypt(ciphertext, tag, key, plaintext []byte) error {
	cipher_key, _ := kdf(key, 32, 0)

	//NOTE: nul iv! this is ok because the key is single use
	var iv [12]byte

	block, err := aes.NewCipher(cipher_key)
	if err != nil {
		return err
	}
	cipher, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	buf := make([]byte, 0, len(plaintext)+Aes256Gcm.TagLen())
	buf = cipher.Seal(buf, iv[:], plaintext, nil)
	copy(ciphertext, buf[:len(plaintext)])
	copy(tag, buf[len(plaintext):])

	return nil
}

func (aes256Gcm) Decrypt(plaintext, key, ciphertext, tag []byte) error {
	cipher_key, _ := kdf(key, 32, 0)

	//NOTE: nul iv! this is ok because the key is single use
	var iv [12]byte
	block, err := aes.NewCipher(cipher_key)
	if err != nil {
		return err
	}
	cipher, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	buf := make([]byte, len(plaintext)+Aes256Gcm.TagLen())
	copy(buf[:len(plaintext)], ciphertext)
	copy(buf[len(plaintext):], tag)
	_, err = cipher.Open(plaintext[:0], iv[:], buf, nil)
	if err != nil {
		if err.Error() == "cipher: message authentication failed" {
			return AuthenticationError
		}
		return err
	}

	return nil
}
