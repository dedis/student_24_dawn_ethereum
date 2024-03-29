package cae

import (
	"crypto/sha256"
	"crypto/hmac"
)


type nullHmacSha256 struct {}

var NullHmacSha256 = nullHmacSha256{}

func (nullHmacSha256) Name() string {
	return "null-hmac-sha256"
}

func (nullHmacSha256) TagLen() int {
	return 32
}

func (nullHmacSha256) Encrypt(ciphertext, tag, key, plaintext []byte) error {
	_, mac_key := kdf(key, 0, 32)

	copy(ciphertext, plaintext)

	mac := hmac.New(sha256.New, mac_key)
	mac.Write(ciphertext)
	mac.Sum(tag[:0])
	return nil
}

func (nullHmacSha256) Decrypt(plaintext, key, ciphertext, tag []byte) error {
	_, mac_key := kdf(key, 0, 32)

	var buf [32]byte
	mac := hmac.New(sha256.New, mac_key)
	mac.Write(ciphertext)
	mac.Sum(buf[:0])
	if !hmac.Equal(tag, buf[:]) {
		return AuthenticationError
	}

	copy(plaintext, ciphertext)
	return nil
}
