package cae

import (
	"crypto/sha512"
	"golang.org/x/crypto/hkdf"
)

func kdf(key []byte, cipher_key_len, mac_key_len int) (cipher_key, mac_key []byte) {
	kdf := hkdf.New(sha512.New, key, nil, nil)
	cipher_key = make([]byte, cipher_key_len)
	mac_key = make([]byte, mac_key_len)
	kdf.Read(cipher_key)
	kdf.Read(mac_key)
	return
}
