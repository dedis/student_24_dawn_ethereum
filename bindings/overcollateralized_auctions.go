// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// OvercollateralizedAuctionsMetaData contains all meta data concerning the OvercollateralizedAuctions contract.
var OvercollateralizedAuctionsMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"blockDelay_\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"auctions\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"collection\",\"type\":\"address\",\"internalType\":\"contractIERC721\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"bidToken\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"proceedsReceiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"opening\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commitDeadline\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revealDeadline\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxBid\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"highestAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"highestBidder\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"commitBid\",\"inputs\":[{\"name\":\"auctionId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"commit\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"computeCommitment\",\"inputs\":[{\"name\":\"blinding\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"commit\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"revealBid\",\"inputs\":[{\"name\":\"auctionId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blinding\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"settle\",\"inputs\":[{\"name\":\"auctionId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"startAuction\",\"inputs\":[{\"name\":\"collection\",\"type\":\"address\",\"internalType\":\"contractIERC721\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"bidToken\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"proceedsReceiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"auctionId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AuctionStarted\",\"inputs\":[{\"name\":\"auctionId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Commit\",\"inputs\":[{\"name\":\"auctionId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Reveal\",\"inputs\":[{\"name\":\"auctionId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
}

// OvercollateralizedAuctionsABI is the input ABI used to generate the binding from.
// Deprecated: Use OvercollateralizedAuctionsMetaData.ABI instead.
var OvercollateralizedAuctionsABI = OvercollateralizedAuctionsMetaData.ABI

// OvercollateralizedAuctions is an auto generated Go binding around an Ethereum contract.
type OvercollateralizedAuctions struct {
	OvercollateralizedAuctionsCaller     // Read-only binding to the contract
	OvercollateralizedAuctionsTransactor // Write-only binding to the contract
	OvercollateralizedAuctionsFilterer   // Log filterer for contract events
}

// OvercollateralizedAuctionsCaller is an auto generated read-only Go binding around an Ethereum contract.
type OvercollateralizedAuctionsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OvercollateralizedAuctionsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OvercollateralizedAuctionsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OvercollateralizedAuctionsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OvercollateralizedAuctionsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OvercollateralizedAuctionsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OvercollateralizedAuctionsSession struct {
	Contract     *OvercollateralizedAuctions // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// OvercollateralizedAuctionsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OvercollateralizedAuctionsCallerSession struct {
	Contract *OvercollateralizedAuctionsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// OvercollateralizedAuctionsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OvercollateralizedAuctionsTransactorSession struct {
	Contract     *OvercollateralizedAuctionsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                     // Transaction auth options to use throughout this session
}

// OvercollateralizedAuctionsRaw is an auto generated low-level Go binding around an Ethereum contract.
type OvercollateralizedAuctionsRaw struct {
	Contract *OvercollateralizedAuctions // Generic contract binding to access the raw methods on
}

// OvercollateralizedAuctionsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OvercollateralizedAuctionsCallerRaw struct {
	Contract *OvercollateralizedAuctionsCaller // Generic read-only contract binding to access the raw methods on
}

// OvercollateralizedAuctionsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OvercollateralizedAuctionsTransactorRaw struct {
	Contract *OvercollateralizedAuctionsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOvercollateralizedAuctions creates a new instance of OvercollateralizedAuctions, bound to a specific deployed contract.
func NewOvercollateralizedAuctions(address common.Address, backend bind.ContractBackend) (*OvercollateralizedAuctions, error) {
	contract, err := bindOvercollateralizedAuctions(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OvercollateralizedAuctions{OvercollateralizedAuctionsCaller: OvercollateralizedAuctionsCaller{contract: contract}, OvercollateralizedAuctionsTransactor: OvercollateralizedAuctionsTransactor{contract: contract}, OvercollateralizedAuctionsFilterer: OvercollateralizedAuctionsFilterer{contract: contract}}, nil
}

// NewOvercollateralizedAuctionsCaller creates a new read-only instance of OvercollateralizedAuctions, bound to a specific deployed contract.
func NewOvercollateralizedAuctionsCaller(address common.Address, caller bind.ContractCaller) (*OvercollateralizedAuctionsCaller, error) {
	contract, err := bindOvercollateralizedAuctions(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OvercollateralizedAuctionsCaller{contract: contract}, nil
}

// NewOvercollateralizedAuctionsTransactor creates a new write-only instance of OvercollateralizedAuctions, bound to a specific deployed contract.
func NewOvercollateralizedAuctionsTransactor(address common.Address, transactor bind.ContractTransactor) (*OvercollateralizedAuctionsTransactor, error) {
	contract, err := bindOvercollateralizedAuctions(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OvercollateralizedAuctionsTransactor{contract: contract}, nil
}

// NewOvercollateralizedAuctionsFilterer creates a new log filterer instance of OvercollateralizedAuctions, bound to a specific deployed contract.
func NewOvercollateralizedAuctionsFilterer(address common.Address, filterer bind.ContractFilterer) (*OvercollateralizedAuctionsFilterer, error) {
	contract, err := bindOvercollateralizedAuctions(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OvercollateralizedAuctionsFilterer{contract: contract}, nil
}

// bindOvercollateralizedAuctions binds a generic wrapper to an already deployed contract.
func bindOvercollateralizedAuctions(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OvercollateralizedAuctionsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OvercollateralizedAuctions.Contract.OvercollateralizedAuctionsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OvercollateralizedAuctions.Contract.OvercollateralizedAuctionsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OvercollateralizedAuctions.Contract.OvercollateralizedAuctionsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OvercollateralizedAuctions.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OvercollateralizedAuctions.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OvercollateralizedAuctions.Contract.contract.Transact(opts, method, params...)
}

// Auctions is a free data retrieval call binding the contract method 0x571a26a0.
//
// Solidity: function auctions(uint256 ) view returns(address collection, uint256 tokenId, address bidToken, address proceedsReceiver, uint64 opening, uint64 commitDeadline, uint64 revealDeadline, uint256 maxBid, uint256 highestAmount, address highestBidder)
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsCaller) Auctions(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Collection       common.Address
	TokenId          *big.Int
	BidToken         common.Address
	ProceedsReceiver common.Address
	Opening          uint64
	CommitDeadline   uint64
	RevealDeadline   uint64
	MaxBid           *big.Int
	HighestAmount    *big.Int
	HighestBidder    common.Address
}, error) {
	var out []interface{}
	err := _OvercollateralizedAuctions.contract.Call(opts, &out, "auctions", arg0)

	outstruct := new(struct {
		Collection       common.Address
		TokenId          *big.Int
		BidToken         common.Address
		ProceedsReceiver common.Address
		Opening          uint64
		CommitDeadline   uint64
		RevealDeadline   uint64
		MaxBid           *big.Int
		HighestAmount    *big.Int
		HighestBidder    common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Collection = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.TokenId = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.BidToken = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.ProceedsReceiver = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.Opening = *abi.ConvertType(out[4], new(uint64)).(*uint64)
	outstruct.CommitDeadline = *abi.ConvertType(out[5], new(uint64)).(*uint64)
	outstruct.RevealDeadline = *abi.ConvertType(out[6], new(uint64)).(*uint64)
	outstruct.MaxBid = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.HighestAmount = *abi.ConvertType(out[8], new(*big.Int)).(**big.Int)
	outstruct.HighestBidder = *abi.ConvertType(out[9], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// Auctions is a free data retrieval call binding the contract method 0x571a26a0.
//
// Solidity: function auctions(uint256 ) view returns(address collection, uint256 tokenId, address bidToken, address proceedsReceiver, uint64 opening, uint64 commitDeadline, uint64 revealDeadline, uint256 maxBid, uint256 highestAmount, address highestBidder)
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsSession) Auctions(arg0 *big.Int) (struct {
	Collection       common.Address
	TokenId          *big.Int
	BidToken         common.Address
	ProceedsReceiver common.Address
	Opening          uint64
	CommitDeadline   uint64
	RevealDeadline   uint64
	MaxBid           *big.Int
	HighestAmount    *big.Int
	HighestBidder    common.Address
}, error) {
	return _OvercollateralizedAuctions.Contract.Auctions(&_OvercollateralizedAuctions.CallOpts, arg0)
}

// Auctions is a free data retrieval call binding the contract method 0x571a26a0.
//
// Solidity: function auctions(uint256 ) view returns(address collection, uint256 tokenId, address bidToken, address proceedsReceiver, uint64 opening, uint64 commitDeadline, uint64 revealDeadline, uint256 maxBid, uint256 highestAmount, address highestBidder)
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsCallerSession) Auctions(arg0 *big.Int) (struct {
	Collection       common.Address
	TokenId          *big.Int
	BidToken         common.Address
	ProceedsReceiver common.Address
	Opening          uint64
	CommitDeadline   uint64
	RevealDeadline   uint64
	MaxBid           *big.Int
	HighestAmount    *big.Int
	HighestBidder    common.Address
}, error) {
	return _OvercollateralizedAuctions.Contract.Auctions(&_OvercollateralizedAuctions.CallOpts, arg0)
}

// ComputeCommitment is a free data retrieval call binding the contract method 0xddd2ced5.
//
// Solidity: function computeCommitment(bytes32 blinding, address bidder, uint256 amount) pure returns(bytes32 commit)
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsCaller) ComputeCommitment(opts *bind.CallOpts, blinding [32]byte, bidder common.Address, amount *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _OvercollateralizedAuctions.contract.Call(opts, &out, "computeCommitment", blinding, bidder, amount)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ComputeCommitment is a free data retrieval call binding the contract method 0xddd2ced5.
//
// Solidity: function computeCommitment(bytes32 blinding, address bidder, uint256 amount) pure returns(bytes32 commit)
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsSession) ComputeCommitment(blinding [32]byte, bidder common.Address, amount *big.Int) ([32]byte, error) {
	return _OvercollateralizedAuctions.Contract.ComputeCommitment(&_OvercollateralizedAuctions.CallOpts, blinding, bidder, amount)
}

// ComputeCommitment is a free data retrieval call binding the contract method 0xddd2ced5.
//
// Solidity: function computeCommitment(bytes32 blinding, address bidder, uint256 amount) pure returns(bytes32 commit)
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsCallerSession) ComputeCommitment(blinding [32]byte, bidder common.Address, amount *big.Int) ([32]byte, error) {
	return _OvercollateralizedAuctions.Contract.ComputeCommitment(&_OvercollateralizedAuctions.CallOpts, blinding, bidder, amount)
}

// CommitBid is a paid mutator transaction binding the contract method 0x9468cb61.
//
// Solidity: function commitBid(uint256 auctionId, bytes32 commit) returns()
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsTransactor) CommitBid(opts *bind.TransactOpts, auctionId *big.Int, commit [32]byte) (*types.Transaction, error) {
	return _OvercollateralizedAuctions.contract.Transact(opts, "commitBid", auctionId, commit)
}

// CommitBid is a paid mutator transaction binding the contract method 0x9468cb61.
//
// Solidity: function commitBid(uint256 auctionId, bytes32 commit) returns()
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsSession) CommitBid(auctionId *big.Int, commit [32]byte) (*types.Transaction, error) {
	return _OvercollateralizedAuctions.Contract.CommitBid(&_OvercollateralizedAuctions.TransactOpts, auctionId, commit)
}

// CommitBid is a paid mutator transaction binding the contract method 0x9468cb61.
//
// Solidity: function commitBid(uint256 auctionId, bytes32 commit) returns()
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsTransactorSession) CommitBid(auctionId *big.Int, commit [32]byte) (*types.Transaction, error) {
	return _OvercollateralizedAuctions.Contract.CommitBid(&_OvercollateralizedAuctions.TransactOpts, auctionId, commit)
}

// RevealBid is a paid mutator transaction binding the contract method 0x88b79626.
//
// Solidity: function revealBid(uint256 auctionId, bytes32 blinding, uint256 amount) returns()
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsTransactor) RevealBid(opts *bind.TransactOpts, auctionId *big.Int, blinding [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _OvercollateralizedAuctions.contract.Transact(opts, "revealBid", auctionId, blinding, amount)
}

// RevealBid is a paid mutator transaction binding the contract method 0x88b79626.
//
// Solidity: function revealBid(uint256 auctionId, bytes32 blinding, uint256 amount) returns()
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsSession) RevealBid(auctionId *big.Int, blinding [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _OvercollateralizedAuctions.Contract.RevealBid(&_OvercollateralizedAuctions.TransactOpts, auctionId, blinding, amount)
}

// RevealBid is a paid mutator transaction binding the contract method 0x88b79626.
//
// Solidity: function revealBid(uint256 auctionId, bytes32 blinding, uint256 amount) returns()
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsTransactorSession) RevealBid(auctionId *big.Int, blinding [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _OvercollateralizedAuctions.Contract.RevealBid(&_OvercollateralizedAuctions.TransactOpts, auctionId, blinding, amount)
}

// Settle is a paid mutator transaction binding the contract method 0x8df82800.
//
// Solidity: function settle(uint256 auctionId) returns()
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsTransactor) Settle(opts *bind.TransactOpts, auctionId *big.Int) (*types.Transaction, error) {
	return _OvercollateralizedAuctions.contract.Transact(opts, "settle", auctionId)
}

// Settle is a paid mutator transaction binding the contract method 0x8df82800.
//
// Solidity: function settle(uint256 auctionId) returns()
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsSession) Settle(auctionId *big.Int) (*types.Transaction, error) {
	return _OvercollateralizedAuctions.Contract.Settle(&_OvercollateralizedAuctions.TransactOpts, auctionId)
}

// Settle is a paid mutator transaction binding the contract method 0x8df82800.
//
// Solidity: function settle(uint256 auctionId) returns()
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsTransactorSession) Settle(auctionId *big.Int) (*types.Transaction, error) {
	return _OvercollateralizedAuctions.Contract.Settle(&_OvercollateralizedAuctions.TransactOpts, auctionId)
}

// StartAuction is a paid mutator transaction binding the contract method 0x23df3b99.
//
// Solidity: function startAuction(address collection, uint256 tokenId, address bidToken, address proceedsReceiver) returns(uint256 auctionId)
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsTransactor) StartAuction(opts *bind.TransactOpts, collection common.Address, tokenId *big.Int, bidToken common.Address, proceedsReceiver common.Address) (*types.Transaction, error) {
	return _OvercollateralizedAuctions.contract.Transact(opts, "startAuction", collection, tokenId, bidToken, proceedsReceiver)
}

// StartAuction is a paid mutator transaction binding the contract method 0x23df3b99.
//
// Solidity: function startAuction(address collection, uint256 tokenId, address bidToken, address proceedsReceiver) returns(uint256 auctionId)
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsSession) StartAuction(collection common.Address, tokenId *big.Int, bidToken common.Address, proceedsReceiver common.Address) (*types.Transaction, error) {
	return _OvercollateralizedAuctions.Contract.StartAuction(&_OvercollateralizedAuctions.TransactOpts, collection, tokenId, bidToken, proceedsReceiver)
}

// StartAuction is a paid mutator transaction binding the contract method 0x23df3b99.
//
// Solidity: function startAuction(address collection, uint256 tokenId, address bidToken, address proceedsReceiver) returns(uint256 auctionId)
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsTransactorSession) StartAuction(collection common.Address, tokenId *big.Int, bidToken common.Address, proceedsReceiver common.Address) (*types.Transaction, error) {
	return _OvercollateralizedAuctions.Contract.StartAuction(&_OvercollateralizedAuctions.TransactOpts, collection, tokenId, bidToken, proceedsReceiver)
}

// OvercollateralizedAuctionsAuctionStartedIterator is returned from FilterAuctionStarted and is used to iterate over the raw logs and unpacked data for AuctionStarted events raised by the OvercollateralizedAuctions contract.
type OvercollateralizedAuctionsAuctionStartedIterator struct {
	Event *OvercollateralizedAuctionsAuctionStarted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OvercollateralizedAuctionsAuctionStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OvercollateralizedAuctionsAuctionStarted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OvercollateralizedAuctionsAuctionStarted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OvercollateralizedAuctionsAuctionStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OvercollateralizedAuctionsAuctionStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OvercollateralizedAuctionsAuctionStarted represents a AuctionStarted event raised by the OvercollateralizedAuctions contract.
type OvercollateralizedAuctionsAuctionStarted struct {
	AuctionId *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAuctionStarted is a free log retrieval operation binding the contract event 0x1bb96dff6ab5005aff98cdc0cf176bb7d8e0423cb48e02217d35b042cec81e9f.
//
// Solidity: event AuctionStarted(uint256 auctionId)
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsFilterer) FilterAuctionStarted(opts *bind.FilterOpts) (*OvercollateralizedAuctionsAuctionStartedIterator, error) {

	logs, sub, err := _OvercollateralizedAuctions.contract.FilterLogs(opts, "AuctionStarted")
	if err != nil {
		return nil, err
	}
	return &OvercollateralizedAuctionsAuctionStartedIterator{contract: _OvercollateralizedAuctions.contract, event: "AuctionStarted", logs: logs, sub: sub}, nil
}

// WatchAuctionStarted is a free log subscription operation binding the contract event 0x1bb96dff6ab5005aff98cdc0cf176bb7d8e0423cb48e02217d35b042cec81e9f.
//
// Solidity: event AuctionStarted(uint256 auctionId)
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsFilterer) WatchAuctionStarted(opts *bind.WatchOpts, sink chan<- *OvercollateralizedAuctionsAuctionStarted) (event.Subscription, error) {

	logs, sub, err := _OvercollateralizedAuctions.contract.WatchLogs(opts, "AuctionStarted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OvercollateralizedAuctionsAuctionStarted)
				if err := _OvercollateralizedAuctions.contract.UnpackLog(event, "AuctionStarted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAuctionStarted is a log parse operation binding the contract event 0x1bb96dff6ab5005aff98cdc0cf176bb7d8e0423cb48e02217d35b042cec81e9f.
//
// Solidity: event AuctionStarted(uint256 auctionId)
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsFilterer) ParseAuctionStarted(log types.Log) (*OvercollateralizedAuctionsAuctionStarted, error) {
	event := new(OvercollateralizedAuctionsAuctionStarted)
	if err := _OvercollateralizedAuctions.contract.UnpackLog(event, "AuctionStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OvercollateralizedAuctionsCommitIterator is returned from FilterCommit and is used to iterate over the raw logs and unpacked data for Commit events raised by the OvercollateralizedAuctions contract.
type OvercollateralizedAuctionsCommitIterator struct {
	Event *OvercollateralizedAuctionsCommit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OvercollateralizedAuctionsCommitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OvercollateralizedAuctionsCommit)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OvercollateralizedAuctionsCommit)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OvercollateralizedAuctionsCommitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OvercollateralizedAuctionsCommitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OvercollateralizedAuctionsCommit represents a Commit event raised by the OvercollateralizedAuctions contract.
type OvercollateralizedAuctionsCommit struct {
	AuctionId *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterCommit is a free log retrieval operation binding the contract event 0x5bdd2fc99022530157777690475b670d3872f32262eb1d47d9ba8000dad58f87.
//
// Solidity: event Commit(uint256 auctionId)
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsFilterer) FilterCommit(opts *bind.FilterOpts) (*OvercollateralizedAuctionsCommitIterator, error) {

	logs, sub, err := _OvercollateralizedAuctions.contract.FilterLogs(opts, "Commit")
	if err != nil {
		return nil, err
	}
	return &OvercollateralizedAuctionsCommitIterator{contract: _OvercollateralizedAuctions.contract, event: "Commit", logs: logs, sub: sub}, nil
}

// WatchCommit is a free log subscription operation binding the contract event 0x5bdd2fc99022530157777690475b670d3872f32262eb1d47d9ba8000dad58f87.
//
// Solidity: event Commit(uint256 auctionId)
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsFilterer) WatchCommit(opts *bind.WatchOpts, sink chan<- *OvercollateralizedAuctionsCommit) (event.Subscription, error) {

	logs, sub, err := _OvercollateralizedAuctions.contract.WatchLogs(opts, "Commit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OvercollateralizedAuctionsCommit)
				if err := _OvercollateralizedAuctions.contract.UnpackLog(event, "Commit", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCommit is a log parse operation binding the contract event 0x5bdd2fc99022530157777690475b670d3872f32262eb1d47d9ba8000dad58f87.
//
// Solidity: event Commit(uint256 auctionId)
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsFilterer) ParseCommit(log types.Log) (*OvercollateralizedAuctionsCommit, error) {
	event := new(OvercollateralizedAuctionsCommit)
	if err := _OvercollateralizedAuctions.contract.UnpackLog(event, "Commit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OvercollateralizedAuctionsRevealIterator is returned from FilterReveal and is used to iterate over the raw logs and unpacked data for Reveal events raised by the OvercollateralizedAuctions contract.
type OvercollateralizedAuctionsRevealIterator struct {
	Event *OvercollateralizedAuctionsReveal // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OvercollateralizedAuctionsRevealIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OvercollateralizedAuctionsReveal)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OvercollateralizedAuctionsReveal)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OvercollateralizedAuctionsRevealIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OvercollateralizedAuctionsRevealIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OvercollateralizedAuctionsReveal represents a Reveal event raised by the OvercollateralizedAuctions contract.
type OvercollateralizedAuctionsReveal struct {
	AuctionId *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterReveal is a free log retrieval operation binding the contract event 0x1747b48b6ade85d7dc97c0f523e0e780795930a468c01b18a51546791fdd3ac0.
//
// Solidity: event Reveal(uint256 auctionId)
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsFilterer) FilterReveal(opts *bind.FilterOpts) (*OvercollateralizedAuctionsRevealIterator, error) {

	logs, sub, err := _OvercollateralizedAuctions.contract.FilterLogs(opts, "Reveal")
	if err != nil {
		return nil, err
	}
	return &OvercollateralizedAuctionsRevealIterator{contract: _OvercollateralizedAuctions.contract, event: "Reveal", logs: logs, sub: sub}, nil
}

// WatchReveal is a free log subscription operation binding the contract event 0x1747b48b6ade85d7dc97c0f523e0e780795930a468c01b18a51546791fdd3ac0.
//
// Solidity: event Reveal(uint256 auctionId)
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsFilterer) WatchReveal(opts *bind.WatchOpts, sink chan<- *OvercollateralizedAuctionsReveal) (event.Subscription, error) {

	logs, sub, err := _OvercollateralizedAuctions.contract.WatchLogs(opts, "Reveal")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OvercollateralizedAuctionsReveal)
				if err := _OvercollateralizedAuctions.contract.UnpackLog(event, "Reveal", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseReveal is a log parse operation binding the contract event 0x1747b48b6ade85d7dc97c0f523e0e780795930a468c01b18a51546791fdd3ac0.
//
// Solidity: event Reveal(uint256 auctionId)
func (_OvercollateralizedAuctions *OvercollateralizedAuctionsFilterer) ParseReveal(log types.Log) (*OvercollateralizedAuctionsReveal, error) {
	event := new(OvercollateralizedAuctionsReveal)
	if err := _OvercollateralizedAuctions.contract.UnpackLog(event, "Reveal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
