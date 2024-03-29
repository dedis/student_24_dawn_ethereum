package cae

type null struct {}

var Null = null{}

func (null) Name() string {
	return "null"
}

func (null) TagLen() int {
	return 0
}

func (null) Encrypt(ciphertext, tag, key, plaintext []byte) error {
	_, _ = kdf(key, 0, 0)

	copy(ciphertext, plaintext)

	return nil
}

func (null) Decrypt(plaintext, key, ciphertext, tag []byte) error {
	_, _ = kdf(key, 0, 0)

	copy(plaintext, ciphertext)

	return nil
}
