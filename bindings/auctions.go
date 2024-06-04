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

// AuctionsMetaData contains all meta data concerning the Auctions contract.
var AuctionsMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"auctions\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"collection\",\"type\":\"address\",\"internalType\":\"contractIERC721\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"bidToken\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"proceedsReceiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"opening\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commitDeadline\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revealDeadline\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxBid\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"highestAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"highestBidder\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"settle\",\"inputs\":[{\"name\":\"auctionId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"startAuction\",\"inputs\":[{\"name\":\"collection\",\"type\":\"address\",\"internalType\":\"contractIERC721\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"bidToken\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"proceedsReceiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"auctionId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AuctionStarted\",\"inputs\":[{\"name\":\"auctionId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
}

// AuctionsABI is the input ABI used to generate the binding from.
// Deprecated: Use AuctionsMetaData.ABI instead.
var AuctionsABI = AuctionsMetaData.ABI

// Auctions is an auto generated Go binding around an Ethereum contract.
type Auctions struct {
	AuctionsCaller     // Read-only binding to the contract
	AuctionsTransactor // Write-only binding to the contract
	AuctionsFilterer   // Log filterer for contract events
}

// AuctionsCaller is an auto generated read-only Go binding around an Ethereum contract.
type AuctionsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AuctionsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AuctionsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AuctionsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AuctionsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AuctionsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AuctionsSession struct {
	Contract     *Auctions         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AuctionsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AuctionsCallerSession struct {
	Contract *AuctionsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// AuctionsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AuctionsTransactorSession struct {
	Contract     *AuctionsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// AuctionsRaw is an auto generated low-level Go binding around an Ethereum contract.
type AuctionsRaw struct {
	Contract *Auctions // Generic contract binding to access the raw methods on
}

// AuctionsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AuctionsCallerRaw struct {
	Contract *AuctionsCaller // Generic read-only contract binding to access the raw methods on
}

// AuctionsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AuctionsTransactorRaw struct {
	Contract *AuctionsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAuctions creates a new instance of Auctions, bound to a specific deployed contract.
func NewAuctions(address common.Address, backend bind.ContractBackend) (*Auctions, error) {
	contract, err := bindAuctions(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Auctions{AuctionsCaller: AuctionsCaller{contract: contract}, AuctionsTransactor: AuctionsTransactor{contract: contract}, AuctionsFilterer: AuctionsFilterer{contract: contract}}, nil
}

// NewAuctionsCaller creates a new read-only instance of Auctions, bound to a specific deployed contract.
func NewAuctionsCaller(address common.Address, caller bind.ContractCaller) (*AuctionsCaller, error) {
	contract, err := bindAuctions(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AuctionsCaller{contract: contract}, nil
}

// NewAuctionsTransactor creates a new write-only instance of Auctions, bound to a specific deployed contract.
func NewAuctionsTransactor(address common.Address, transactor bind.ContractTransactor) (*AuctionsTransactor, error) {
	contract, err := bindAuctions(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AuctionsTransactor{contract: contract}, nil
}

// NewAuctionsFilterer creates a new log filterer instance of Auctions, bound to a specific deployed contract.
func NewAuctionsFilterer(address common.Address, filterer bind.ContractFilterer) (*AuctionsFilterer, error) {
	contract, err := bindAuctions(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AuctionsFilterer{contract: contract}, nil
}

// bindAuctions binds a generic wrapper to an already deployed contract.
func bindAuctions(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AuctionsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Auctions *AuctionsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Auctions.Contract.AuctionsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Auctions *AuctionsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Auctions.Contract.AuctionsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Auctions *AuctionsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Auctions.Contract.AuctionsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Auctions *AuctionsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Auctions.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Auctions *AuctionsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Auctions.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Auctions *AuctionsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Auctions.Contract.contract.Transact(opts, method, params...)
}

// Auctions is a free data retrieval call binding the contract method 0x571a26a0.
//
// Solidity: function auctions(uint256 ) view returns(address collection, uint256 tokenId, address bidToken, address proceedsReceiver, uint64 opening, uint64 commitDeadline, uint64 revealDeadline, uint256 maxBid, uint256 highestAmount, address highestBidder)
func (_Auctions *AuctionsCaller) Auctions(opts *bind.CallOpts, arg0 *big.Int) (struct {
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
	err := _Auctions.contract.Call(opts, &out, "auctions", arg0)

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
func (_Auctions *AuctionsSession) Auctions(arg0 *big.Int) (struct {
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
	return _Auctions.Contract.Auctions(&_Auctions.CallOpts, arg0)
}

// Auctions is a free data retrieval call binding the contract method 0x571a26a0.
//
// Solidity: function auctions(uint256 ) view returns(address collection, uint256 tokenId, address bidToken, address proceedsReceiver, uint64 opening, uint64 commitDeadline, uint64 revealDeadline, uint256 maxBid, uint256 highestAmount, address highestBidder)
func (_Auctions *AuctionsCallerSession) Auctions(arg0 *big.Int) (struct {
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
	return _Auctions.Contract.Auctions(&_Auctions.CallOpts, arg0)
}

// Settle is a paid mutator transaction binding the contract method 0x8df82800.
//
// Solidity: function settle(uint256 auctionId) returns()
func (_Auctions *AuctionsTransactor) Settle(opts *bind.TransactOpts, auctionId *big.Int) (*types.Transaction, error) {
	return _Auctions.contract.Transact(opts, "settle", auctionId)
}

// Settle is a paid mutator transaction binding the contract method 0x8df82800.
//
// Solidity: function settle(uint256 auctionId) returns()
func (_Auctions *AuctionsSession) Settle(auctionId *big.Int) (*types.Transaction, error) {
	return _Auctions.Contract.Settle(&_Auctions.TransactOpts, auctionId)
}

// Settle is a paid mutator transaction binding the contract method 0x8df82800.
//
// Solidity: function settle(uint256 auctionId) returns()
func (_Auctions *AuctionsTransactorSession) Settle(auctionId *big.Int) (*types.Transaction, error) {
	return _Auctions.Contract.Settle(&_Auctions.TransactOpts, auctionId)
}

// StartAuction is a paid mutator transaction binding the contract method 0x23df3b99.
//
// Solidity: function startAuction(address collection, uint256 tokenId, address bidToken, address proceedsReceiver) returns(uint256 auctionId)
func (_Auctions *AuctionsTransactor) StartAuction(opts *bind.TransactOpts, collection common.Address, tokenId *big.Int, bidToken common.Address, proceedsReceiver common.Address) (*types.Transaction, error) {
	return _Auctions.contract.Transact(opts, "startAuction", collection, tokenId, bidToken, proceedsReceiver)
}

// StartAuction is a paid mutator transaction binding the contract method 0x23df3b99.
//
// Solidity: function startAuction(address collection, uint256 tokenId, address bidToken, address proceedsReceiver) returns(uint256 auctionId)
func (_Auctions *AuctionsSession) StartAuction(collection common.Address, tokenId *big.Int, bidToken common.Address, proceedsReceiver common.Address) (*types.Transaction, error) {
	return _Auctions.Contract.StartAuction(&_Auctions.TransactOpts, collection, tokenId, bidToken, proceedsReceiver)
}

// StartAuction is a paid mutator transaction binding the contract method 0x23df3b99.
//
// Solidity: function startAuction(address collection, uint256 tokenId, address bidToken, address proceedsReceiver) returns(uint256 auctionId)
func (_Auctions *AuctionsTransactorSession) StartAuction(collection common.Address, tokenId *big.Int, bidToken common.Address, proceedsReceiver common.Address) (*types.Transaction, error) {
	return _Auctions.Contract.StartAuction(&_Auctions.TransactOpts, collection, tokenId, bidToken, proceedsReceiver)
}

// AuctionsAuctionStartedIterator is returned from FilterAuctionStarted and is used to iterate over the raw logs and unpacked data for AuctionStarted events raised by the Auctions contract.
type AuctionsAuctionStartedIterator struct {
	Event *AuctionsAuctionStarted // Event containing the contract specifics and raw log

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
func (it *AuctionsAuctionStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AuctionsAuctionStarted)
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
		it.Event = new(AuctionsAuctionStarted)
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
func (it *AuctionsAuctionStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AuctionsAuctionStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AuctionsAuctionStarted represents a AuctionStarted event raised by the Auctions contract.
type AuctionsAuctionStarted struct {
	AuctionId *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAuctionStarted is a free log retrieval operation binding the contract event 0x1bb96dff6ab5005aff98cdc0cf176bb7d8e0423cb48e02217d35b042cec81e9f.
//
// Solidity: event AuctionStarted(uint256 auctionId)
func (_Auctions *AuctionsFilterer) FilterAuctionStarted(opts *bind.FilterOpts) (*AuctionsAuctionStartedIterator, error) {

	logs, sub, err := _Auctions.contract.FilterLogs(opts, "AuctionStarted")
	if err != nil {
		return nil, err
	}
	return &AuctionsAuctionStartedIterator{contract: _Auctions.contract, event: "AuctionStarted", logs: logs, sub: sub}, nil
}

// WatchAuctionStarted is a free log subscription operation binding the contract event 0x1bb96dff6ab5005aff98cdc0cf176bb7d8e0423cb48e02217d35b042cec81e9f.
//
// Solidity: event AuctionStarted(uint256 auctionId)
func (_Auctions *AuctionsFilterer) WatchAuctionStarted(opts *bind.WatchOpts, sink chan<- *AuctionsAuctionStarted) (event.Subscription, error) {

	logs, sub, err := _Auctions.contract.WatchLogs(opts, "AuctionStarted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AuctionsAuctionStarted)
				if err := _Auctions.contract.UnpackLog(event, "AuctionStarted", log); err != nil {
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
func (_Auctions *AuctionsFilterer) ParseAuctionStarted(log types.Log) (*AuctionsAuctionStarted, error) {
	event := new(AuctionsAuctionStarted)
	if err := _Auctions.contract.UnpackLog(event, "AuctionStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
