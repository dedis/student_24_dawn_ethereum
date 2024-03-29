package cae

type nullNokdf struct {}

var NullNokdf = nullNokdf{}

func (nullNokdf) Name() string {
	return "null-nokdf"
}

func (nullNokdf) TagLen() int {
	return 0
}

func (nullNokdf) Encrypt(ciphertext, tag, key, plaintext []byte) error {
	copy(ciphertext, plaintext)

	return nil
}

func (nullNokdf) Decrypt(plaintext, key, ciphertext, tag []byte) error {
	copy(plaintext, ciphertext)

	return nil
}
