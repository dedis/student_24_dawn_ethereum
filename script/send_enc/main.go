package main

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/f3b"
)

func logMeasurement(elapsed time.Duration) error {
	f, err := os.OpenFile("enctiming.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("%s\n", elapsed))
	if err != nil {
		return err
	}
	return nil
}

func sendEtherF3bEnc(client *ethclient.Client, ks *keystore.KeyStore, from accounts.Account, to common.Address, val *big.Int, gasLimit uint64, calldata []byte) (error) {
	nonce, err := client.PendingNonceAt(context.Background(), from.Address); 
	if err != nil {
		return err
	}

	// get gas price
	gasPrice, err := client.SuggestGasPrice(context.Background()); 
	if err != nil {
		return err
	}
	// get chainID
	chainID, err := client.ChainID(context.Background()); 
	if err != nil {
		return err
	}

	start := time.Now()
	dkgcli := f3b.NewDkgCli()

	label := binary.BigEndian.AppendUint64(from.Address.Bytes(), nonce)
	plaintext := append(to.Bytes(), calldata...)

	ciphertext, err := dkgcli.Encrypt(label, plaintext)
	if err != nil {
		return err
	}
	elapsed := time.Since(start)
	logMeasurement(elapsed)

	enc := &types.EncryptedTx{
		ChainID:    chainID,
		Nonce:      nonce,
		GasFeeCap:  gasPrice,
		GasTipCap:  big.NewInt(10),
		Gas:        gasLimit,
		Value:      val,
		Payload:    ciphertext,
	}
	tx := types.NewTx(enc)

	// sign
	signedTx, err := ks.SignTx(from, tx, chainID);
	if err != nil {
		return err
	}

	if err = client.SendTransaction(context.Background(), signedTx); err != nil {
		return err
	}

	fmt.Println(signedTx.Hash().Hex())
	return nil
}

func usage() {
	fmt.Println("Usage: send_enc [options] <to> [calldata]")
	fmt.Println("issue an encrypted transaction to the given address with optional calldata")
	flag.PrintDefaults()
}

func main() {
	if err := main2(); err != nil {
		log.Fatal(err)
	}
}

func main2() error {
	var to common.Address
	var calldata []byte

	sender := flag.String("sender", "", "sender address")
	var value *big.Int
	flag.Func("value", "call value in wei", func(s string) error {
		var ok bool
		if 
		value, ok = new(big.Int).SetString(s, 10);
		!ok {
			return fmt.Errorf("invalid value: %s", s)
		}
		return nil
	})
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
		return fmt.Errorf("missing receiver address")
	}
	if arg := flag.Arg(0); common.IsHexAddress(arg) {
		to = common.HexToAddress(arg)
	}
	// optional, defaults to ""
	calldata = common.FromHex(flag.Arg(1))

	// unlock the pre fund user account
	ks := keystore.NewKeyStore(os.Getenv("ETH_KEYSTORE"), keystore.StandardScryptN, keystore.StandardScryptP)
	var from accounts.Account
	if sender == nil {
		from = ks.Accounts()[0]
	} else if !common.IsHexAddress(*sender) {
		return fmt.Errorf("invalid sender address: %s", *sender)
	} else {
	addr := common.HexToAddress(*sender)
	if !ks.HasAddress(addr) {
		return fmt.Errorf("no key for *sender address: %s", *sender)
	}
	from = accounts.Account{Address: addr}
}
	if err := ks.Unlock(from, ""); err != nil {
		return err
	}

	client, err := ethclient.Dial(os.Getenv("ETH_RPC_URL"))
	if err != nil {
		return err
	}
	defer client.Close()

	gasLimit := uint64(1000000) // FIXME: hardcoded

	return sendEtherF3bEnc(client, ks, from, to, value, gasLimit, calldata)
}
