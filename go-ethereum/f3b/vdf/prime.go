// Copyright EPFL DEDIS

// based on go 1.22.2 crypto/rand package:
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package vdf

import (
	"errors"
	"io"
	"math/big"
)

// DeterministicPrime returns a number of the given bit length that is prime with high probability.
// DeterministicPrime will return error for any error returned by [rand.Read] or if bits < 2.
// unlike rand.Prime, DeterministicPrime is deterministic and will always return the same prime given the same randomness source.
func DeterministicPrime(rand io.Reader, bits int) (*big.Int, error) {
	if bits < 2 {
		return nil, errors.New("crypto/rand: prime size must be at least 2-bit")
	}

	b := uint(bits % 8)
	if b == 0 {
		b = 8
	}

	bytes := make([]byte, (bits+7)/8)
	p := new(big.Int)

	for {
		if _, err := io.ReadFull(rand, bytes); err != nil {
			return nil, err
		}

		// Clear bits in the first byte to make sure the candidate has a size <= bits.
		bytes[0] &= uint8(int(1<<b) - 1)
		// Don't let the value be too small, i.e, set the most significant two bits.
		// Setting the top two bits, rather than just the top bit,
		// means that when two of these values are multiplied together,
		// the result isn't ever one bit short.
		if b >= 2 {
			bytes[0] |= 3 << (b - 2)
		} else {
			// Here b==1, because b cannot be zero.
			bytes[0] |= 1
			if len(bytes) > 1 {
				bytes[1] |= 0x80
			}
		}
		// Make the value odd since an even number this large certainly isn't prime.
		bytes[len(bytes)-1] |= 1

		p.SetBytes(bytes)
		if p.ProbablyPrime(20) {
			return p, nil
		}
	}
}
