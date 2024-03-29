package cae

import (
	"fmt"

	"github.com/aead/chacha20/chacha"
)


type _chacha struct {
	rounds int
}

var Chacha20 = _chacha{20}
var Chacha12 = _chacha{12}
var Chacha8 = _chacha{8}

func (cc _chacha) Name() string {
	return fmt.Sprintf("chacha%d", cc.rounds)
}

func (cc _chacha) TagLen() int {
	return 0
}

func (cc _chacha) Encrypt(ciphertext, tag, key, plaintext []byte) error {
	cipher_key, _ := kdf(key, chacha.KeySize, 0)

	//NOTE: nul iv! this is ok because the key is single use
	var iv [chacha.NonceSize]byte
	chacha.XORKeyStream(ciphertext, plaintext, iv[:], cipher_key, cc.rounds)

	return nil
}

func (cc _chacha) Decrypt(plaintext, key, ciphertext, tag []byte) error {
	return cc.Encrypt(plaintext, tag, key, ciphertext)
}
