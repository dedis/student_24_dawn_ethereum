// Copyright 2015 The go-ethereum Authors
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

package core

import (
	"fmt"
	"encoding/binary"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/misc"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/f3b"
)

// StateProcessor is a basic Processor, which takes care of transitioning
// state from one point to another.
//
// StateProcessor implements Processor.
type StateProcessor struct {
	config *params.ChainConfig // Chain configuration options
	bc     *BlockChain         // Canonical block chain
	engine consensus.Engine    // Consensus engine used for block rewards
}

// NewStateProcessor initialises a new StateProcessor.
func NewStateProcessor(config *params.ChainConfig, bc *BlockChain, engine consensus.Engine) *StateProcessor {
	return &StateProcessor{
		config: config,
		bc:     bc,
		engine: engine,
	}
}

// Process processes the state changes according to the Ethereum rules by running
// the transaction messages using the statedb and applying any rewards to both
// the processor (coinbase) and any included uncles.
//
// Process returns the receipts and logs accumulated during the process and
// returns the amount of gas that was used in the process. If any of the
// transactions failed to execute due to insufficient gas it will return an error.
func (p *StateProcessor) Process(block *types.Block, statedb *state.StateDB, cfg vm.Config) (types.Receipts, []*types.Log, uint64, error) {
	log.Info("## Process validation of a received block")
	var (
		receipts    types.Receipts
		usedGas     = new(uint64)
		header      = block.Header()
		blockHash   = block.Hash()
		blockNumber = block.Number()
		allLogs     []*types.Log
		gp          = new(GasPool).AddGas(block.GasLimit())
	)
	// Mutate the block and state according to any hard-fork specs
	if p.config.DAOForkSupport && p.config.DAOForkBlock != nil && p.config.DAOForkBlock.Cmp(block.Number()) == 0 {
		misc.ApplyDAOHardFork(statedb)
	}
	// Retrieve previous ordered enc txs
	// pendingEncryptedTxs := RetrievePendingEncryptedTransactions(p.bc, types.EncryptedBlockDelay)
	// Append with the plaintext tx
	// allTx := append(pendingEncryptedTxs, block.Transactions()...)
	// Iterate over and process the individual transactions
	var (
		beneficiary     common.Address
		receipt         *types.Receipt
		isExecEncrypted bool
		rcAuth          common.Address
		err             error
	)
	// tmp_rcs := p.bc.GetReceiptsByHash(block.Hash())
	// if len(tmp_rcs) != block.Transactions().Len() {
	// 	panic("unequal transaction length and receipts")
	// }
	orderBlock := RetrieveOrderBlock(p.bc, types.EncryptedBlockDelay)
	if orderBlock != nil {
		rcAuth, err = p.bc.engine.Author(orderBlock.Header())
		if err != nil {
			panic(fmt.Sprintf("fail to retrieve author: %s", err))
		}
		log.Error(fmt.Sprintf("last block signer: %s", rcAuth))
	}
	for i, tx := range block.Transactions() {
		signer := types.MakeSigner(p.config, header.Number)
		msg, err := tx.AsMessage(signer, header.BaseFee)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("could not apply tx %d [%v]: %w", i, tx.Hash().Hex(), err)
		}

		statedb.Prepare(tx.Hash(), i)

		if isExecEncrypted, _ = isExecuteEncryptedTx(statedb, signer, p.config, tx); isExecEncrypted {
			beneficiary = rcAuth
			log.Info(fmt.Sprintf("[VERIFY][ENC EXE][beneficiary]: %s", beneficiary))
		} else {
			beneficiary, _ = p.bc.engine.Author(block.Header())
			log.Info(fmt.Sprintf("[VERIFY][Normal][beneficiary]: %s", beneficiary))
		}

		blockContext := NewEVMBlockContext(header, p.bc, &beneficiary)
		vmenv := vm.NewEVM(blockContext, vm.TxContext{}, statedb, p.config, cfg)

		receipt, err = applyTransaction(msg, p.config, &beneficiary, gp, statedb, blockNumber, blockHash, tx, usedGas, vmenv, isExecEncrypted)

		log.Info(fmt.Sprintf("[VERIFY][ENC][RC]] receipt key appended: %v", receipt.Key))

		if err != nil {
			return nil, nil, 0, fmt.Errorf("could not apply tx %d [%v]: %w", i, tx.Hash().Hex(), err)
		}
		receipts = append(receipts, receipt)
		allLogs = append(allLogs, receipt.Logs...)
	}
	// Finalize the block, applying any consensus engine specific extras (e.g. block rewards)
	p.engine.Finalize(p.bc, header, statedb, block.Transactions(), block.Uncles())

	return receipts, allLogs, *usedGas, nil
}

func isExecuteEncryptedTx(statedb *state.StateDB, signer types.Signer, cf *params.ChainConfig, tx *types.Transaction) (bool, error) {
	var (
		from common.Address
		err  error
	)
	if from, err = types.Sender(signer, tx); err != nil {
		return false, err
	}
	stNonce := statedb.GetNonce(from)
	// executing encrypted tx from last finalty block
	if tx.Type() == types.EncryptedTxType && stNonce > tx.Nonce() {
		return true, nil
	}

	return false, nil
}

func RetrieveOrderBlock(wc *BlockChain, numbersBack uint64) *types.Block {
	currentNumber := wc.CurrentHeader().Number.Uint64()
	// regardless of the genesis block
	if currentNumber <= numbersBack {
		return nil
	}

	previousNumber := currentNumber - numbersBack
	preBlock := wc.GetBlockByNumber(previousNumber)
	return preBlock
}

func RetrievePendingEncryptedTransactions(wc *BlockChain, numbersBack uint64) types.Transactions {
	encryptedTxs := make(types.Transactions, 0)
	currentNumber := wc.CurrentHeader().Number.Uint64()
	// regardless of the genesis block
	if currentNumber <= numbersBack {
		return encryptedTxs
	}

	previousNumber := currentNumber - numbersBack
	preBlock := wc.GetBlockByNumber(previousNumber)

	txs := preBlock.Transactions()
	receipts := wc.GetReceiptsByHash(preBlock.Hash())

	if len(txs) != len(receipts) {
		panic("unequal length of txs and receipts")
	}
	var rc *types.Receipt

	// retrieve all the pending encrypted txs that should be executed in this block
	for i, tx := range txs {
		if tx.Type() == types.EncryptedTxType {
			rc = receipts[i]

			if rc.Key == nil || len(rc.Key) == 0 { // if there is no key attached to the receipt, this encrypted tx must have not been executed
				log.Info(fmt.Sprintf("[ENC][RETRIEVE] tx hash: %v", tx.Hash().String()))
				encryptedTxs = append(encryptedTxs, tx)
			}
		}
	}

	return encryptedTxs

}

func applyTransaction(msg types.Message, config *params.ChainConfig, author *common.Address, gp *GasPool, statedb *state.StateDB, blockNumber *big.Int, blockHash common.Hash, tx *types.Transaction, usedGas *uint64, evm *vm.EVM, isExecEncrypted bool) (*types.Receipt, error) {
	// Create a new context to be used in the EVM environment.
	txContext := NewEVMTxContext(msg)
	evm.Reset(txContext, statedb)

	var calldata []byte
	var to *common.Address

	if isExecEncrypted {
		dkgCli := f3b.NewDkgCli()
		label := msg.From().Bytes()
		label = binary.BigEndian.AppendUint64(label, msg.Nonce())
		plaintext, err := dkgCli.Decrypt(label, msg.Data())
		to = new(common.Address)
		*to = common.BytesToAddress(plaintext[:common.AddressLength])
		calldata = plaintext[common.AddressLength:]
		if err != nil {
			panic("decryptMsgData: fail on decryption")
		}
		calldata = plaintext
	} else {
		calldata = msg.Data()
		to = msg.To()
	}

	// Apply the transaction to the current state (included in the env).
	result, err := ApplyMessage(evm, msg, gp, calldata)
	if err != nil {
		return nil, err
	}

	// Update the state with pending changes.
	var root []byte
	if config.IsByzantium(blockNumber) {
		statedb.Finalise(true)
	} else {
		root = statedb.IntermediateRoot(config.IsEIP158(blockNumber)).Bytes()
	}
	*usedGas += result.UsedGas

	// Create a new receipt for the transaction, storing the intermediate root and gas used
	// by the tx.
	receipt := &types.Receipt{Type: tx.Type(), PostState: root, CumulativeGasUsed: *usedGas}
	if result.Failed() {
		receipt.Status = types.ReceiptStatusFailed
	} else {
		receipt.Status = types.ReceiptStatusSuccessful
	}
	receipt.TxHash = tx.Hash()
	receipt.GasUsed = result.UsedGas

	// If the transaction created a contract, store the creation address in the receipt.
	if to == nil {
		receipt.ContractAddress = crypto.CreateAddress(evm.TxContext.Origin, tx.Nonce())
	}

	// If this is the execution of an encrypted tx, then add the key to the receipt
	if isExecEncrypted {
		receipt.Key = []byte("TODO")

	}

	// Set the receipt logs and create the bloom filter.
	receipt.Logs = statedb.GetLogs(tx.Hash(), blockHash)
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	receipt.BlockHash = blockHash
	receipt.BlockNumber = blockNumber
	receipt.TransactionIndex = uint(statedb.TxIndex())
	return receipt, err
}

// ApplyTransaction attempts to apply a transaction to the given state database
// and uses the input parameters for its environment. It returns the receipt
// for the transaction, gas used and an error if the transaction failed,
// indicating the block was invalid.
func ApplyTransaction(config *params.ChainConfig, bc ChainContext, author *common.Address, gp *GasPool, statedb *state.StateDB, header *types.Header, tx *types.Transaction, usedGas *uint64, cfg vm.Config, isExecEncrypted bool) (*types.Receipt, error) {
	msg, err := tx.AsMessage(types.MakeSigner(config, header.Number), header.BaseFee)
	if err != nil {
		return nil, err
	}
	// Create a new context to be used in the EVM environment
	blockContext := NewEVMBlockContext(header, bc, author)
	vmenv := vm.NewEVM(blockContext, vm.TxContext{}, statedb, config, cfg)
	return applyTransaction(msg, config, author, gp, statedb, header.Number, header.Hash(), tx, usedGas, vmenv, isExecEncrypted)
}
