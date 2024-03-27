// Copyright EPFL DEDIS

// Package cae implements committing authenticated encryption schemes.
//
// We define a committing authenticated encryption scheme (CAE) in the sense of [CR22](https://eprint.iacr.org/2022/1260.pdf).
//
// We do not use nonces, instead keys are single-use.


package cae

type CAE interface {
	// Encrypt encrypts and commits to the plaintext using the key.
	// The key is consumed and must not be used for a different plaintext.
	//
	// ciphertext must be as long as plaintext and tag must have size TagSize()
	Encrypt(ciphertext, tag, key, plaintext []byte) error
	// Encrypt decrypts the cipherext using the key.
	// It returns AuthenticationError if the ciphertext is not valid.
	//
	// plaintext must be as long as ciphertext and tag must have size TagSize()
	Decrypt(plaintext, key, ciphertext, tag []byte) error
	// Return the size of the authentication tag in bytes.
	TagSize() int
}

type authenticationError struct {}
func (authenticationError) Error() string {
	return "cae: authentication error"
}
var AuthenticationError error = authenticationError{}

// For development convenience, this is used to select the CAE to use.
var Selected CAE = ChaCha20HmacSha256{}
