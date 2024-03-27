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

func (ChaCha20HmacSha256) TagSize() int {
	return 32
}

func (ChaCha20HmacSha256) Encrypt(ciphertext, tag, key, plaintext []byte) error {
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

func (ChaCha20HmacSha256) Decrypt(plaintext, key, ciphertext, tag []byte) error {
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
