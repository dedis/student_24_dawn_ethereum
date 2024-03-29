
package cae

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
)


type rkAes256Gcm struct {}

var RkAes256Gcm = rkAes256Gcm{}

func (rkAes256Gcm) Name() string {
	return "rk-aes256-gcm"
}

func (rkAes256Gcm) TagLen() int {
	return 32 // 16 GCM tag + 16 RK tag
}

func (rkAes256Gcm) Encrypt(ciphertext, tag, key, plaintext []byte) error {
	cipher_key, key_commit := kdf(key, 32, 16)

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
	buf = append(buf, key_commit...)
	copy(ciphertext, buf[:len(plaintext)])
	copy(tag, buf[len(plaintext):])

	return nil
}

func (rkAes256Gcm) Decrypt(plaintext, key, ciphertext, tag []byte) error {
	cipher_key, key_commit := kdf(key, 32, 16)

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
	copy(buf[len(plaintext):], tag[:16])
	_, err = cipher.Open(plaintext[:0], iv[:], buf, nil)
	if err != nil {
		if err.Error() == "cipher: message authentication failed" {
			return AuthenticationError
		}
		return err
	}

	if !hmac.Equal(tag[16:], key_commit) {
		return AuthenticationError
	}

	return nil
}
