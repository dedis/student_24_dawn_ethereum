package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"math/big"
	"io"
	"fmt"
	"os"
	"time"
	"sync"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/f3b"
	"github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/dedis/f3b-ethereum/bindings"
)


func main() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
	if err := Main(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type Addresses map[string]common.Address

func LoadAddresses() (Addresses, error) {
	addressesFile := os.Getenv("ADDRESSES_FILE")
	var obj map[string]string
	f, err := os.Open(addressesFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}
	addresses := make(Addresses, len(obj))
	for k, v := range obj {
		addresses[k] = common.HexToAddress(v)
	}

	return addresses, nil
}

// Ad-hoc functional options on TransactOpts
// https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html
func with(transactOpts *bind.TransactOpts, mods ...func(*bind.TransactOpts)) *bind.TransactOpts {
	ret := new(bind.TransactOpts)
	*ret = *transactOpts
	for _, f := range mods {
		f(ret)
	}
	return ret
}

func value(value *big.Int) func(*bind.TransactOpts) {
	return func(transactOpts *bind.TransactOpts) {
		transactOpts.Value = value
	}
}

func encrypt(targetBlock uint64) func(*bind.TransactOpts) {
	return func(transactOpts *bind.TransactOpts) {
		prevSigner := transactOpts.Signer
		transactOpts.Signer = func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
			tx, err := tx.Encrypt(addr, targetBlock)
			if err != nil {
				return nil, err
			}
			return prevSigner(addr, tx)
		}
	}
}

func gasLimit(limit uint64) func(*bind.TransactOpts) {
	return func(transactOpts *bind.TransactOpts) {
		transactOpts.GasLimit = limit
	}
}

type Scenario struct {
	Context context.Context
	Client  *ethclient.Client
	ChainID *big.Int
	Wallet  *hdwallet.Wallet
	Addresses Addresses

	WETH *bindings.WETH
	Auctions *bindings.Auctions
	OvercollateralizedAuctions *bindings.OvercollateralizedAuctions
	SimpleAuctions *bindings.SimpleAuctions
	Collection *bindings.Collection
	Params *f3b.FullParams

	BiddersReady sync.WaitGroup
	BiddersDone sync.WaitGroup
}

func (s *Scenario) bidderScript(account accounts.Account) {
	defer s.BiddersDone.Done()
	transactOpts, err := s.bidderScriptPrepare(account)
	// Auction starts
	if err == nil {
		err = s.bidderScriptBid(transactOpts)
	}
	if err != nil {
		log.Error("bidder failed", "err", err)
	}
}

func (s *Scenario) bidderScriptPrepare(account accounts.Account) (*bind.TransactOpts, error) {
	defer s.BiddersReady.Done()

	privkey, err := s.Wallet.PrivateKey(account)
	if err != nil {
		return nil, err
	}

	transactOpts, err := bind.NewKeyedTransactorWithChainID(privkey, s.ChainID)
	if err != nil {
		return nil, err
	}
	transactOpts.Context = s.Context

	maxBid, _ := new(big.Int).SetString("10000000000000000000", 10)
	_, err = s.WETH.Deposit(with(transactOpts, value(maxBid)))
	if err != nil {
		return nil, err
	}

	_, err = s.checkSuccess(s.WETH.Approve(transactOpts, s.Addresses["auctions"], maxBid))
	if err != nil {
		return nil, err
	}

	log.Info("bidder ready")

	return transactOpts, nil
}

func (s *Scenario) bidderScriptBid(transactOpts *bind.TransactOpts) error {
	auctionStarted, err := s.waitForAuction()
	if err != nil {
		return err
	}

	amount := common.Big3 // FIXME: hardcoded
	if s.OvercollateralizedAuctions != nil {
	err = s.waitForBlockNumber(auctionStarted.Opening)
	if err != nil {
		return err
	}

	var blinding [32]byte
	rand.Read(blinding[:])
	callOpts := &bind.CallOpts{Context: s.Context, From: transactOpts.From}
	commit, err := s.OvercollateralizedAuctions.ComputeCommitment(callOpts, blinding, transactOpts.From, amount)
	if err != nil {
		return err
	}

	_, err = s.checkSuccess(s.OvercollateralizedAuctions.CommitBid(transactOpts, auctionStarted.AuctionId, commit))
	if err != nil {
		return err
	}
	log.Info("bid committed")

	err = s.waitForBlockNumber(auctionStarted.CommitDeadline)
	if err != nil {
		return err
	}

	_, err = s.checkSuccess(s.OvercollateralizedAuctions.RevealBid(transactOpts, auctionStarted.AuctionId, blinding, amount))
	if err != nil {
		return err
	}
	log.Info("bid revealed")
} else {
	err = s.waitForBlockNumber(auctionStarted.Opening - s.Params.BlockDelay) // account for latency
	if err != nil {
		return err
	}
	// 21k base gas, bid itself should be up to ~117k, plus slack
	const limit = 21_000 + 120_000 + 10_000;
	targetBlock := auctionStarted.Opening
	_, err = s.checkSuccess(s.SimpleAuctions.Bid(with(transactOpts, encrypt(targetBlock), gasLimit(limit)), auctionStarted.AuctionId, amount))
	if err != nil {
		return err
	}
	log.Info("bid sent")
}


	return nil
}

func (s *Scenario) waitForAuction() (*bindings.AuctionsAuctionStarted, error) {
	// FIXME: should use Watch instead of looping
	for {
		it, err := s.Auctions.FilterAuctionStarted(&bind.FilterOpts{Start: 0, End: nil, Context: s.Context})
		if err != nil {
			return nil, err
		}

		for it.Next() {
			return it.Event, nil
		}
		if it.Error() != nil {
			return nil, err
		}
	}
}

func (s *Scenario) waitForBlockNumber(bn uint64) error {
	// short path
	cur, err := s.Client.BlockNumber(s.Context)
	if cur >= bn {
		return nil
	}
	ch := make(chan *types.Header)
	sub, err := s.Client.SubscribeNewHead(s.Context, ch)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	// in case of race condition
	cur, err = s.Client.BlockNumber(s.Context)
	if cur >= bn {
		return nil
	}

	for block := range ch {
		if block.Number.Cmp(big.NewInt(int64(bn))) >= 0 {
			break
		}
	}
	return nil
}

func (s *Scenario) operatorScript() error {
	ks := keystore.NewKeyStore("keystore/", keystore.StandardScryptN, keystore.StandardScryptP)
	transactOpts := new(bind.TransactOpts)
	transactOpts.From = common.HexToAddress("0x280F6B48E4d9aEe0Efdb04EeBe882023357f6434")
	transactOpts.Signer = func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
		from := accounts.Account{Address: addr}
		if err := ks.Unlock(from, ""); err != nil {
			return nil, err
		}
		return ks.SignTx(from, tx, s.ChainID)
	}
	transactOpts.Context = s.Context

	tokenId := common.Big1
	_, err := s.checkSuccess(s.Collection.Approve(transactOpts, s.Addresses["auctions"], tokenId))
	if err != nil {
		return err
	}

	s.BiddersReady.Wait()

	_, err = s.checkSuccess(s.Auctions.StartAuction(transactOpts, s.Addresses["collection"], tokenId, s.Addresses["weth"], common.Address{}))
	if err != nil {
		return err
	}

	log.Info("auction started")

	auctionStarted, err := s.waitForAuction()
	if err != nil {
		return err
	}
	fmt.Println(auctionStarted)

	return s.waitForBlockNumber(auctionStarted.RevealDeadline)
}

func (s *Scenario) checkSuccess(tx *types.Transaction, err error) (*types.Transaction, error) {
	if err != nil {
		return nil, err
	}

	var rcpt *types.Receipt
	for {
		rcpt, err = s.Client.TransactionReceipt(s.Context, tx.Hash())
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if rcpt.Status != types.ReceiptStatusSuccessful {
		return nil, fmt.Errorf("transaction failed")
	}
	return tx, nil
}


func Main() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := ethclient.Dial("ws://localhost:8546")
	if err != nil {
		return err
	}
	defer client.Close()

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return err
	}

	mnemonic := os.Getenv("MNEMONIC")
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return err
	}

	addresses, err := LoadAddresses()
	if err != nil {
		return err
	}

	weth, err := bindings.NewWETH(addresses["weth"], client)
	if err != nil {
		return err
	}

	auctions, err := bindings.NewAuctions(addresses["auctions"], client)
	if err != nil {
		return err
	}

	collection, err := bindings.NewCollection(addresses["collection"], client)
	if err != nil {
		return err
	}

	params, err := f3b.ReadParams()
	if err != nil {
		return err
	}

	s := Scenario{
		Context: ctx,
		Client:  client,
		ChainID: chainID,
		Wallet:  wallet,
		Addresses: addresses,
		WETH:    weth,
		Auctions: auctions,
		Collection: collection,
		Params: params,
	}

	if params.Protocol == "" {
		// no encryption, have to use overcollateralization
		s.OvercollateralizedAuctions, err = bindings.NewOvercollateralizedAuctions(addresses["auctions"], client)
	} else {
		s.SimpleAuctions, err = bindings.NewSimpleAuctions(addresses["auctions"], client)
	}
	if err != nil {
		return err
	}

	nBidders := params.NumBidders
	s.BiddersReady.Add(nBidders)
	s.BiddersDone.Add(nBidders)
	it := accounts.DefaultIterator(hdwallet.DefaultBaseDerivationPath)
	for i := 0; i < nBidders; i++ {
		account, err := s.Wallet.Derive(it(), true)
		if err != nil {
			return err
		}
		go s.bidderScript(account)
	}
	defer s.BiddersDone.Wait()

	return s.operatorScript()
}
