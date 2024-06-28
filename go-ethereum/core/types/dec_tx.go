// Copyright EPFL DEDIS

package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type DecryptedTx struct {
	ChainID    *big.Int
	Nonce      uint64
	GasTipCap  *big.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap  *big.Int // a.k.a. maxFeePerGas
	Gas        uint64
	Value      *big.Int
	To	 *common.Address
	Data    []byte
	AccessList AccessList
	EncKey     []byte // The symmetric key k encrypted for the SMC
	Reveal     []byte // witness data to help decrypt EncKey
	TargetBlock uint64 // F3B-TIBE: block number the transaction is encrypted for
	From       common.Address // the sender address, needed to properly reencrypt

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`
}

// copy creates a deep copy of the transaction data and initializes all fields.
func (tx *DecryptedTx) copy() TxData {
	cpy := &DecryptedTx{
		Nonce: tx.Nonce,
		To:    copyAddressPtr(tx.To),
		Data:  common.CopyBytes(tx.Data),
		EncKey:   common.CopyBytes(tx.EncKey),
		Reveal:   common.CopyBytes(tx.Reveal),
		TargetBlock: tx.TargetBlock,
		From:     tx.From,
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
func (tx *DecryptedTx) txType() byte           { return DecryptedTxType }
func (tx *DecryptedTx) chainID() *big.Int      { return tx.ChainID }
func (tx *DecryptedTx) accessList() AccessList { return tx.AccessList }
func (tx *DecryptedTx) data() []byte           { return tx.Data }
func (tx *DecryptedTx) gas() uint64            { return tx.Gas }
func (tx *DecryptedTx) gasFeeCap() *big.Int    { return tx.GasFeeCap }
func (tx *DecryptedTx) gasTipCap() *big.Int    { return tx.GasTipCap }
func (tx *DecryptedTx) gasPrice() *big.Int     { return tx.GasFeeCap }
func (tx *DecryptedTx) value() *big.Int        { return tx.Value }
func (tx *DecryptedTx) nonce() uint64          { return tx.Nonce }
func (tx *DecryptedTx) to() *common.Address    { return tx.To }
func (tx *DecryptedTx) ciphertext() []byte     { return nil }
func (tx *DecryptedTx) tag() []byte            { return nil }
func (tx *DecryptedTx) encKey() []byte         { return tx.EncKey }
func (tx *DecryptedTx) reveal() []byte         { return tx.Reveal }
func (tx *DecryptedTx) targetBlock() uint64    { return tx.TargetBlock }

func (tx *DecryptedTx) rawSignatureValues() (v, r, s *big.Int) {
	return tx.V, tx.R, tx.S
}

func (tx *DecryptedTx) setSignatureValues(chainID, v, r, s *big.Int) {
	tx.ChainID, tx.V, tx.R, tx.S = chainID, v, r, s
}
