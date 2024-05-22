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

// SimpleAuctionsMetaData contains all meta data concerning the SimpleAuctions contract.
var SimpleAuctionsMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"auctions\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"collection\",\"type\":\"address\",\"internalType\":\"contractIERC721\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"bidToken\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"proceedsReceiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"deadline\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"highestAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"highestBidder\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"bid\",\"inputs\":[{\"name\":\"auctionId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"settle\",\"inputs\":[{\"name\":\"auctionId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"startAuction\",\"inputs\":[{\"name\":\"collection\",\"type\":\"address\",\"internalType\":\"contractIERC721\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"bidToken\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"proceedsReceiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"auctionId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AuctionStarted\",\"inputs\":[{\"name\":\"auctionId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"collection\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIERC721\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"bidToken\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIERC20\"},{\"name\":\"proceedsReceiver\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"commitDeadline\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"revealDeadline\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"maxBid\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
}

// SimpleAuctionsABI is the input ABI used to generate the binding from.
// Deprecated: Use SimpleAuctionsMetaData.ABI instead.
var SimpleAuctionsABI = SimpleAuctionsMetaData.ABI

// SimpleAuctions is an auto generated Go binding around an Ethereum contract.
type SimpleAuctions struct {
	SimpleAuctionsCaller     // Read-only binding to the contract
	SimpleAuctionsTransactor // Write-only binding to the contract
	SimpleAuctionsFilterer   // Log filterer for contract events
}

// SimpleAuctionsCaller is an auto generated read-only Go binding around an Ethereum contract.
type SimpleAuctionsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleAuctionsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SimpleAuctionsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleAuctionsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SimpleAuctionsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleAuctionsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SimpleAuctionsSession struct {
	Contract     *SimpleAuctions   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SimpleAuctionsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SimpleAuctionsCallerSession struct {
	Contract *SimpleAuctionsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// SimpleAuctionsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SimpleAuctionsTransactorSession struct {
	Contract     *SimpleAuctionsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// SimpleAuctionsRaw is an auto generated low-level Go binding around an Ethereum contract.
type SimpleAuctionsRaw struct {
	Contract *SimpleAuctions // Generic contract binding to access the raw methods on
}

// SimpleAuctionsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SimpleAuctionsCallerRaw struct {
	Contract *SimpleAuctionsCaller // Generic read-only contract binding to access the raw methods on
}

// SimpleAuctionsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SimpleAuctionsTransactorRaw struct {
	Contract *SimpleAuctionsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSimpleAuctions creates a new instance of SimpleAuctions, bound to a specific deployed contract.
func NewSimpleAuctions(address common.Address, backend bind.ContractBackend) (*SimpleAuctions, error) {
	contract, err := bindSimpleAuctions(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SimpleAuctions{SimpleAuctionsCaller: SimpleAuctionsCaller{contract: contract}, SimpleAuctionsTransactor: SimpleAuctionsTransactor{contract: contract}, SimpleAuctionsFilterer: SimpleAuctionsFilterer{contract: contract}}, nil
}

// NewSimpleAuctionsCaller creates a new read-only instance of SimpleAuctions, bound to a specific deployed contract.
func NewSimpleAuctionsCaller(address common.Address, caller bind.ContractCaller) (*SimpleAuctionsCaller, error) {
	contract, err := bindSimpleAuctions(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleAuctionsCaller{contract: contract}, nil
}

// NewSimpleAuctionsTransactor creates a new write-only instance of SimpleAuctions, bound to a specific deployed contract.
func NewSimpleAuctionsTransactor(address common.Address, transactor bind.ContractTransactor) (*SimpleAuctionsTransactor, error) {
	contract, err := bindSimpleAuctions(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleAuctionsTransactor{contract: contract}, nil
}

// NewSimpleAuctionsFilterer creates a new log filterer instance of SimpleAuctions, bound to a specific deployed contract.
func NewSimpleAuctionsFilterer(address common.Address, filterer bind.ContractFilterer) (*SimpleAuctionsFilterer, error) {
	contract, err := bindSimpleAuctions(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimpleAuctionsFilterer{contract: contract}, nil
}

// bindSimpleAuctions binds a generic wrapper to an already deployed contract.
func bindSimpleAuctions(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SimpleAuctionsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimpleAuctions *SimpleAuctionsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleAuctions.Contract.SimpleAuctionsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimpleAuctions *SimpleAuctionsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleAuctions.Contract.SimpleAuctionsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimpleAuctions *SimpleAuctionsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleAuctions.Contract.SimpleAuctionsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimpleAuctions *SimpleAuctionsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleAuctions.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimpleAuctions *SimpleAuctionsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleAuctions.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimpleAuctions *SimpleAuctionsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleAuctions.Contract.contract.Transact(opts, method, params...)
}

// Auctions is a free data retrieval call binding the contract method 0x571a26a0.
//
// Solidity: function auctions(uint256 ) view returns(address collection, uint256 tokenId, address bidToken, address proceedsReceiver, uint64 deadline, uint256 highestAmount, address highestBidder)
func (_SimpleAuctions *SimpleAuctionsCaller) Auctions(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Collection       common.Address
	TokenId          *big.Int
	BidToken         common.Address
	ProceedsReceiver common.Address
	Deadline         uint64
	HighestAmount    *big.Int
	HighestBidder    common.Address
}, error) {
	var out []interface{}
	err := _SimpleAuctions.contract.Call(opts, &out, "auctions", arg0)

	outstruct := new(struct {
		Collection       common.Address
		TokenId          *big.Int
		BidToken         common.Address
		ProceedsReceiver common.Address
		Deadline         uint64
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
	outstruct.Deadline = *abi.ConvertType(out[4], new(uint64)).(*uint64)
	outstruct.HighestAmount = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.HighestBidder = *abi.ConvertType(out[6], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// Auctions is a free data retrieval call binding the contract method 0x571a26a0.
//
// Solidity: function auctions(uint256 ) view returns(address collection, uint256 tokenId, address bidToken, address proceedsReceiver, uint64 deadline, uint256 highestAmount, address highestBidder)
func (_SimpleAuctions *SimpleAuctionsSession) Auctions(arg0 *big.Int) (struct {
	Collection       common.Address
	TokenId          *big.Int
	BidToken         common.Address
	ProceedsReceiver common.Address
	Deadline         uint64
	HighestAmount    *big.Int
	HighestBidder    common.Address
}, error) {
	return _SimpleAuctions.Contract.Auctions(&_SimpleAuctions.CallOpts, arg0)
}

// Auctions is a free data retrieval call binding the contract method 0x571a26a0.
//
// Solidity: function auctions(uint256 ) view returns(address collection, uint256 tokenId, address bidToken, address proceedsReceiver, uint64 deadline, uint256 highestAmount, address highestBidder)
func (_SimpleAuctions *SimpleAuctionsCallerSession) Auctions(arg0 *big.Int) (struct {
	Collection       common.Address
	TokenId          *big.Int
	BidToken         common.Address
	ProceedsReceiver common.Address
	Deadline         uint64
	HighestAmount    *big.Int
	HighestBidder    common.Address
}, error) {
	return _SimpleAuctions.Contract.Auctions(&_SimpleAuctions.CallOpts, arg0)
}

// Bid is a paid mutator transaction binding the contract method 0x598647f8.
//
// Solidity: function bid(uint256 auctionId, uint256 amount) returns()
func (_SimpleAuctions *SimpleAuctionsTransactor) Bid(opts *bind.TransactOpts, auctionId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _SimpleAuctions.contract.Transact(opts, "bid", auctionId, amount)
}

// Bid is a paid mutator transaction binding the contract method 0x598647f8.
//
// Solidity: function bid(uint256 auctionId, uint256 amount) returns()
func (_SimpleAuctions *SimpleAuctionsSession) Bid(auctionId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _SimpleAuctions.Contract.Bid(&_SimpleAuctions.TransactOpts, auctionId, amount)
}

// Bid is a paid mutator transaction binding the contract method 0x598647f8.
//
// Solidity: function bid(uint256 auctionId, uint256 amount) returns()
func (_SimpleAuctions *SimpleAuctionsTransactorSession) Bid(auctionId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _SimpleAuctions.Contract.Bid(&_SimpleAuctions.TransactOpts, auctionId, amount)
}

// Settle is a paid mutator transaction binding the contract method 0x8df82800.
//
// Solidity: function settle(uint256 auctionId) returns()
func (_SimpleAuctions *SimpleAuctionsTransactor) Settle(opts *bind.TransactOpts, auctionId *big.Int) (*types.Transaction, error) {
	return _SimpleAuctions.contract.Transact(opts, "settle", auctionId)
}

// Settle is a paid mutator transaction binding the contract method 0x8df82800.
//
// Solidity: function settle(uint256 auctionId) returns()
func (_SimpleAuctions *SimpleAuctionsSession) Settle(auctionId *big.Int) (*types.Transaction, error) {
	return _SimpleAuctions.Contract.Settle(&_SimpleAuctions.TransactOpts, auctionId)
}

// Settle is a paid mutator transaction binding the contract method 0x8df82800.
//
// Solidity: function settle(uint256 auctionId) returns()
func (_SimpleAuctions *SimpleAuctionsTransactorSession) Settle(auctionId *big.Int) (*types.Transaction, error) {
	return _SimpleAuctions.Contract.Settle(&_SimpleAuctions.TransactOpts, auctionId)
}

// StartAuction is a paid mutator transaction binding the contract method 0x23df3b99.
//
// Solidity: function startAuction(address collection, uint256 tokenId, address bidToken, address proceedsReceiver) returns(uint256 auctionId)
func (_SimpleAuctions *SimpleAuctionsTransactor) StartAuction(opts *bind.TransactOpts, collection common.Address, tokenId *big.Int, bidToken common.Address, proceedsReceiver common.Address) (*types.Transaction, error) {
	return _SimpleAuctions.contract.Transact(opts, "startAuction", collection, tokenId, bidToken, proceedsReceiver)
}

// StartAuction is a paid mutator transaction binding the contract method 0x23df3b99.
//
// Solidity: function startAuction(address collection, uint256 tokenId, address bidToken, address proceedsReceiver) returns(uint256 auctionId)
func (_SimpleAuctions *SimpleAuctionsSession) StartAuction(collection common.Address, tokenId *big.Int, bidToken common.Address, proceedsReceiver common.Address) (*types.Transaction, error) {
	return _SimpleAuctions.Contract.StartAuction(&_SimpleAuctions.TransactOpts, collection, tokenId, bidToken, proceedsReceiver)
}

// StartAuction is a paid mutator transaction binding the contract method 0x23df3b99.
//
// Solidity: function startAuction(address collection, uint256 tokenId, address bidToken, address proceedsReceiver) returns(uint256 auctionId)
func (_SimpleAuctions *SimpleAuctionsTransactorSession) StartAuction(collection common.Address, tokenId *big.Int, bidToken common.Address, proceedsReceiver common.Address) (*types.Transaction, error) {
	return _SimpleAuctions.Contract.StartAuction(&_SimpleAuctions.TransactOpts, collection, tokenId, bidToken, proceedsReceiver)
}

// SimpleAuctionsAuctionStartedIterator is returned from FilterAuctionStarted and is used to iterate over the raw logs and unpacked data for AuctionStarted events raised by the SimpleAuctions contract.
type SimpleAuctionsAuctionStartedIterator struct {
	Event *SimpleAuctionsAuctionStarted // Event containing the contract specifics and raw log

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
func (it *SimpleAuctionsAuctionStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleAuctionsAuctionStarted)
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
		it.Event = new(SimpleAuctionsAuctionStarted)
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
func (it *SimpleAuctionsAuctionStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleAuctionsAuctionStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleAuctionsAuctionStarted represents a AuctionStarted event raised by the SimpleAuctions contract.
type SimpleAuctionsAuctionStarted struct {
	AuctionId        *big.Int
	Collection       common.Address
	TokenId          *big.Int
	BidToken         common.Address
	ProceedsReceiver common.Address
	CommitDeadline   uint64
	RevealDeadline   uint64
	MaxBid           *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterAuctionStarted is a free log retrieval operation binding the contract event 0x3bab4ad29d487858b174244ae27b86da6c757e2779461dc416762cfc65ad8a28.
//
// Solidity: event AuctionStarted(uint256 auctionId, address collection, uint256 tokenId, address bidToken, address proceedsReceiver, uint64 commitDeadline, uint64 revealDeadline, uint256 maxBid)
func (_SimpleAuctions *SimpleAuctionsFilterer) FilterAuctionStarted(opts *bind.FilterOpts) (*SimpleAuctionsAuctionStartedIterator, error) {

	logs, sub, err := _SimpleAuctions.contract.FilterLogs(opts, "AuctionStarted")
	if err != nil {
		return nil, err
	}
	return &SimpleAuctionsAuctionStartedIterator{contract: _SimpleAuctions.contract, event: "AuctionStarted", logs: logs, sub: sub}, nil
}

// WatchAuctionStarted is a free log subscription operation binding the contract event 0x3bab4ad29d487858b174244ae27b86da6c757e2779461dc416762cfc65ad8a28.
//
// Solidity: event AuctionStarted(uint256 auctionId, address collection, uint256 tokenId, address bidToken, address proceedsReceiver, uint64 commitDeadline, uint64 revealDeadline, uint256 maxBid)
func (_SimpleAuctions *SimpleAuctionsFilterer) WatchAuctionStarted(opts *bind.WatchOpts, sink chan<- *SimpleAuctionsAuctionStarted) (event.Subscription, error) {

	logs, sub, err := _SimpleAuctions.contract.WatchLogs(opts, "AuctionStarted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleAuctionsAuctionStarted)
				if err := _SimpleAuctions.contract.UnpackLog(event, "AuctionStarted", log); err != nil {
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

// ParseAuctionStarted is a log parse operation binding the contract event 0x3bab4ad29d487858b174244ae27b86da6c757e2779461dc416762cfc65ad8a28.
//
// Solidity: event AuctionStarted(uint256 auctionId, address collection, uint256 tokenId, address bidToken, address proceedsReceiver, uint64 commitDeadline, uint64 revealDeadline, uint256 maxBid)
func (_SimpleAuctions *SimpleAuctionsFilterer) ParseAuctionStarted(log types.Log) (*SimpleAuctionsAuctionStarted, error) {
	event := new(SimpleAuctionsAuctionStarted)
	if err := _SimpleAuctions.contract.UnpackLog(event, "AuctionStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
