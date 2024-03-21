package cae

import (
	"crypto/sha512"
	"crypto/sha256"
	"crypto/hmac"

	"golang.org/x/crypto/chacha20"
	"golang.org/x/crypto/hkdf"
)


type ChaCha20HmacSha256 struct {}

func kdf(key []byte, cipher_key_len, mac_key_len int) (cipher_key, mac_key []byte) {
	kdf := hkdf.New(sha512.New, key, nil, nil)
	cipher_key = make([]byte, cipher_key_len)
	mac_key = make([]byte, mac_key_len)
	kdf.Read(cipher_key)
	kdf.Read(mac_key)
	return
}

func (ChaCha20HmacSha256) Encrypt(key, plaintext []byte) ([]byte, error) {
	const tagSize = 32
	cipher_key, mac_key := kdf(key, chacha20.KeySize, 32)
	iv := make([]byte, chacha20.NonceSize) //NOTE: this is ok because the key is used only once
	cipher, err := chacha20.NewUnauthenticatedCipher(cipher_key, iv)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, len(plaintext), len(plaintext)+tagSize)
	cipher.XORKeyStream(ciphertext, plaintext)
	mac := hmac.New(sha256.New, mac_key)
	mac.Write(ciphertext)
	ciphertext = mac.Sum(ciphertext) // appends
	return ciphertext, nil
}

func (ChaCha20HmacSha256) Decrypt(key, ciphertext []byte) ([]byte, error) {
	const tagSize = 32
	cipher_key, mac_key := kdf(key, chacha20.KeySize, 32)
	ciphertext, tag := ciphertext[:len(ciphertext)-tagSize], ciphertext[len(ciphertext)-tagSize:]
	mac := hmac.New(sha256.New, mac_key)
	mac.Write(ciphertext)
	if !hmac.Equal(tag, mac.Sum(nil)) {
		return nil, AuthenticationError
	}

	iv := make([]byte, chacha20.NonceSize)
	cipher, err := chacha20.NewUnauthenticatedCipher(cipher_key, iv)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(ciphertext))
	cipher.XORKeyStream(plaintext, ciphertext)
	return plaintext, nil
}
