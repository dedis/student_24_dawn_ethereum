// Copyright EPFL DEDIS

package types

import (
	"bytes"
)

type ShadowTransaction struct {}
type ShadowTransactions []*ShadowTransaction

func (s ShadowTransactions) Len() int { return len(s) }

// EncodeIndex encodes the i'th transaction to w. Note that this does not check for errors
// because we assume that *Transaction will only ever contain valid txs that were either
// constructed by decoding or via public API in this package.
func (s ShadowTransactions) EncodeIndex(i int, w *bytes.Buffer) {
	/*
	tx := s[i]
	if tx.Type() == LegacyTxType {
		rlp.Encode(w, tx.inner)
	} else {
		tx.encodeTyped(w)
	}
	*/
	panic("fixme")
}
