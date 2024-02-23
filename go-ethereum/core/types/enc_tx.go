// Copyright 2021 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type EncryptedTx struct {
	ChainID    *big.Int
	Nonce      uint64
	GasTipCap  *big.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap  *big.Int // a.k.a. maxFeePerGas
	Gas        uint64
	Value      *big.Int
	Payload    []byte // Enc_k(to | data) symmetric encryption
	Key        []byte // hash(k) | Enc_smc(k)
	AccessList AccessList

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`
}

// copy creates a deep copy of the transaction data and initializes all fields.
func (tx *EncryptedTx) copy() TxData {
	cpy := &EncryptedTx{
		Nonce: tx.Nonce,
		Payload:  common.CopyBytes(tx.Payload),
		Key:   common.CopyBytes(tx.Key),
		Gas:   tx.Gas,
		// These are copied below.
		AccessList: make(AccessList, len(tx.AccessList)),
		Value:      new(big.Int),
		ChainID:    new(big.Int),
		GasTipCap:  new(big.Int),
		GasFeeCap:  new(big.Int),
		V:          new(big.Int),
		R:          new(big.Int),
		S:          new(big.Int),
	}
	copy(cpy.AccessList, tx.AccessList)
	if tx.Value != nil {
		cpy.Value.Set(tx.Value)
	}
	if tx.ChainID != nil {
		cpy.ChainID.Set(tx.ChainID)
	}
	if tx.GasTipCap != nil {
		cpy.GasTipCap.Set(tx.GasTipCap)
	}
	if tx.GasFeeCap != nil {
		cpy.GasFeeCap.Set(tx.GasFeeCap)
	}
	if tx.V != nil {
		cpy.V.Set(tx.V)
	}
	if tx.R != nil {
		cpy.R.Set(tx.R)
	}
	if tx.S != nil {
		cpy.S.Set(tx.S)
	}
	return cpy
}

// accessors for innerTx.
func (tx *EncryptedTx) txType() byte           { return EncryptedTxType }
func (tx *EncryptedTx) chainID() *big.Int      { return tx.ChainID }
func (tx *EncryptedTx) accessList() AccessList { return tx.AccessList }
func (tx *EncryptedTx) data() []byte           { return tx.Payload }
func (tx *EncryptedTx) gas() uint64            { return tx.Gas }
func (tx *EncryptedTx) gasFeeCap() *big.Int    { return tx.GasFeeCap }
func (tx *EncryptedTx) gasTipCap() *big.Int    { return tx.GasTipCap }
func (tx *EncryptedTx) gasPrice() *big.Int     { return tx.GasFeeCap }
func (tx *EncryptedTx) value() *big.Int        { return tx.Value }
func (tx *EncryptedTx) nonce() uint64          { return tx.Nonce }
func (tx *EncryptedTx) to() *common.Address    { return nil }
func (tx *EncryptedTx) key() []byte            { return tx.Key }

func (tx *EncryptedTx) rawSignatureValues() (v, r, s *big.Int) {
	return tx.V, tx.R, tx.S
}

func (tx *EncryptedTx) setSignatureValues(chainID, v, r, s *big.Int) {
	tx.ChainID, tx.V, tx.R, tx.S = chainID, v, r, s
}