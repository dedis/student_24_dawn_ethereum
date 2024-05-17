package main

import (
	"context"
	"encoding/json"
	"math/big"
	"io"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
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

func LoadAddresses() (map[string]common.Address, error) {
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
	addresses := make(map[string]common.Address, len(obj))
	for k, v := range obj {
		addresses[k] = common.HexToAddress(v)
	}

	return addresses, nil
}

func Main() error {
	client, err := ethclient.Dial(os.Getenv("ETH_RPC_URL"))
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := context.Background()

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return err
	}

	mnemonic := os.Getenv("MNEMONIC")
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return err
	}

	account, err := wallet.Derive(accounts.DefaultBaseDerivationPath, true)
	if err != nil {
		return err
	}

	privkey, err := wallet.PrivateKey(account)
	if err != nil {
		return err
	}

	addresses, err := LoadAddresses()
	if err != nil {
		return err
	}

	contract, err := bindings.NewWETH(addresses["weth"], client)
	if err != nil {
		return err
	}

	transactOpts, err := bind.NewKeyedTransactorWithChainID(privkey, chainID)
	if err != nil {
		return err
	}
	weth := bindings.WETHSession{Contract: contract, TransactOpts: *transactOpts}
	maxBid, _ := new(big.Int).SetString("10000000000000000000", 10)
	weth.TransactOpts.Value = maxBid
	tx, err := weth.Deposit()
	if err != nil {
		return err
	}
	weth.TransactOpts.Value = nil

	tx, err = weth.Approve(addresses["auctions"], maxBid)
	if err != nil {
		return err
	}

	var rcpt *types.Receipt
	for i := 0; i < 2000; i++ {
		rcpt, err = client.TransactionReceipt(ctx, tx.Hash())
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println(rcpt)

	return nil
}
