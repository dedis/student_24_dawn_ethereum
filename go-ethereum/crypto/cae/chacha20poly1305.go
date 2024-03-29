package cae

import (
	"golang.org/x/crypto/chacha20poly1305"
)


type chacha20Poly1305 struct {}

var Chacha20Poly1305 = chacha20Poly1305{}

func (chacha20Poly1305) Name() string {
	return "chacha20-poly1305"
}

func (chacha20Poly1305) TagLen() int {
	return chacha20poly1305.Overhead
}

func (chacha20Poly1305) Encrypt(ciphertext, tag, key, plaintext []byte) error {
	cipher_key, _ := kdf(key, chacha20poly1305.KeySize, 0)

	//NOTE: nul iv! this is ok because the key is single use
	var iv [chacha20poly1305.NonceSize]byte
	cipher, err := chacha20poly1305.New(cipher_key)
	if err != nil {
		return err
	}
	buf := make([]byte, 0, len(plaintext)+chacha20poly1305.Overhead)
	buf = cipher.Seal(buf, iv[:], plaintext, nil)
	copy(ciphertext, buf[:len(plaintext)])
	copy(tag, buf[len(plaintext):])

	return nil
}

func (chacha20Poly1305) Decrypt(plaintext, key, ciphertext, tag []byte) error {
	cipher_key, _ := kdf(key, chacha20poly1305.KeySize, 0)

	//NOTE: nul iv! this is ok because the key is single use
	var iv [chacha20poly1305.NonceSize]byte
	cipher, err := chacha20poly1305.New(cipher_key)
	if err != nil {
		return err
	}
	buf := make([]byte, len(plaintext)+chacha20poly1305.Overhead)
	copy(buf[:len(plaintext)], ciphertext)
	copy(buf[len(plaintext):], tag)
	_, err = cipher.Open(plaintext[:0], iv[:], buf, nil)
	if err != nil {
		if err.Error() == "chacha20poly1305: message authentication failed" {
			return AuthenticationError
		}
		return err
	}

	return nil
}
