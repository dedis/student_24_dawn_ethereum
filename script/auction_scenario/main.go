package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"math/big"
	"io"
	"flag"
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
	"github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/dedis/f3b-ethereum/bindings"
)


func usage() {
	fmt.Println("Usage: send_enc [options] <to> [calldata]")
	fmt.Println("issue an encrypted transaction to the given address with optional calldata")
	flag.PrintDefaults()
}

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

type Scenario struct {
	Context context.Context
	Client  *ethclient.Client
	ChainID *big.Int
	Wallet  *hdwallet.Wallet
	Addresses Addresses

	WETH *bindings.WETH
	Auctions *bindings.OvercollateralizedAuctions
	Collection *bindings.Collection

	BiddersReady sync.WaitGroup
	AuctionId *big.Int
	AuctionStarted chan struct{}
}

func (s *Scenario) bidderScript(account accounts.Account) {
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
	transactOpts.Value = maxBid
	_, err = s.WETH.Deposit(transactOpts)
	if err != nil {
		return nil, err
	}
	transactOpts.Value = nil

	_, err = s.checkSuccess(s.WETH.Approve(transactOpts, s.Addresses["auctions"], maxBid))
	if err != nil {
		return nil, err
	}

	log.Info("bidder ready")

	return transactOpts, nil
}

func (s *Scenario) bidderScriptBid(transactOpts *bind.TransactOpts) error {
	<-s.AuctionStarted
	var blinding [32]byte
	rand.Read(blinding[:])
	amount := common.Big3 // FIXME: hardcoded
	callOpts := &bind.CallOpts{Context: s.Context, From: transactOpts.From}
	commit, err := s.Auctions.ComputeCommitment(callOpts, blinding, transactOpts.From, amount)
	if err != nil {
		return err
	}

	_, err = s.checkSuccess(s.Auctions.CommitBid(transactOpts, s.AuctionId, commit))
	log.Info("bid committed")

	time.Sleep(60 * time.Second)
	_, err = s.checkSuccess(s.Auctions.RevealBid(transactOpts, s.AuctionId, blinding, amount))
	log.Info("bid revealed")

	return nil
}

func (s *Scenario) operatorScript() error {
	s.AuctionStarted = make(chan struct{})
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

	s.AuctionId = common.Big0 // FIXME: hardcoded
	close(s.AuctionStarted)

	log.Info("auction started")

	time.Sleep(120 * time.Second)

	return nil
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

	client, err := ethclient.Dial(os.Getenv("ETH_RPC_URL"))
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

	auctions, err := bindings.NewOvercollateralizedAuctions(addresses["auctions"], client)
	if err != nil {
		return err
	}

	collection, err := bindings.NewCollection(addresses["collection"], client)
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
	}

	nBidders := 10
	s.BiddersReady.Add(nBidders)
	it := accounts.DefaultIterator(hdwallet.DefaultBaseDerivationPath)
	for i := 0; i < nBidders; i++ {
		account, err := s.Wallet.Derive(it(), true)
		if err != nil {
			return err
		}
		go s.bidderScript(account)
	}

	return s.operatorScript()
}
