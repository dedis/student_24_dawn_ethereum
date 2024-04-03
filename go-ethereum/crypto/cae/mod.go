// Copyright EPFL DEDIS

// Package cae implements committing authenticated encryption schemes.
//
// We define a committing authenticated encryption scheme (CAE) in the sense of [CR22](https://eprint.iacr.org/2022/1260.pdf).
//
// We do not use nonces, instead keys are single-use.


package cae

type Scheme interface {
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
	// Return the unique, human-readable name of the Scheme.
	Name() string
	// Return the size of the authentication tag in bytes.
	TagLen() int
}

type authenticationError struct {}
func (authenticationError) Error() string {
	return "cae: authentication error"
}
var AuthenticationError error = authenticationError{}

// For development convenience, this is used to select the scheme to use.
var Selected Scheme = RkChacha20Poly1305

var AllSchemes = []Scheme{
	Aes256Gcm,
	Aes256CtrHmacSha256,
	Chacha20HmacSha256,
	Chacha20Poly1305,
	Chacha20Alt,
	Chacha20,
	Chacha12,
	Chacha8,
	Null,
	NullNokdf,
	NullHmacSha256,
	RkAes256Gcm,
	RkChacha20Poly1305,
}
