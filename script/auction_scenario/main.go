package main

import (
	"context"
	"encoding/json"
	"io"
	"flag"
	"fmt"
	"os"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/miguelmota/go-ethereum-hdwallet"
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

	account, err := wallet.Derive(accounts.DefaultBaseDerivationPath, true)
	if err != nil {
		return err
	}

	fmt.Println(account)

	nonce, err := client.PendingNonceAt(context.Background(), account.Address)
	if err != nil {
		return err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	calldata := []byte{}
	to := addresses["weth"]
	val, _ := new(big.Int).SetString("10000000000000000000", 10)

	gasLimit := uint64(1000000) // FIXME: hardcoded


	tx := types.NewTx(&types.AccessListTx{
		ChainID:    chainID,
		Nonce:      nonce,
		GasPrice:   gasPrice,
		Gas:        gasLimit,
		Value:      val,
		To:	    &to,
		Data:       calldata,
	})

	// sign
	signedTx, err := wallet.SignTx(account, tx, chainID)
	if err != nil {
		return err
	}

	if err = client.SendTransaction(context.Background(), signedTx); err != nil {
		return err
	}


	return nil
}
