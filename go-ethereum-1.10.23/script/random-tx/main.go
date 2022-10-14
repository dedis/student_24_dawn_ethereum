package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const authority_file = "../.ethereum/keystore/UTC--2022-09-13T11-34-29.303731400Z--280f6b48e4d9aee0efdb04eebe882023357f6434"
const user_file = "../.ethereum/keystore/UTC--2022-09-13T11-35-11.765870700Z--f5f341cd21350259a8666b3a5fe47132eff57838"

func get_balance(client *ethclient.Client, addr common.Address) (*big.Int, *big.Float) {
	weiValue, err := client.BalanceAt(context.Background(), addr, nil)
	if err != nil {
		log.Fatal(err)
	}
	fbalance := new(big.Float)
	fbalance.SetString(weiValue.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	return weiValue, ethValue
}

func sendEtherPlaintext(client *ethclient.Client, nonce uint64, ks *keystore.KeyStore, from, to accounts.Account, val *big.Int, gasLimit uint64) {
	var err error
	var gasPrice, chainID *big.Int
	var signedTx *types.Transaction

	// get gas price
	if gasPrice, err = client.SuggestGasPrice(context.Background()); err != nil {
		log.Fatal(err)
	}
	// get chainID
	if chainID, err = client.ChainID(context.Background()); err != nil {
		log.Fatal(err)
	}

	tx := types.NewTransaction(nonce, to.Address, val, gasLimit, gasPrice, nil)

	// sign
	if signedTx, err = ks.SignTx(from, tx, chainID); err != nil {
		log.Fatal(err)
	}

	if err = client.SendTransaction(context.Background(), signedTx); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Plaintext Transaction send: %v\n", signedTx.Hash().Hex())
}

func sendEtherEncrypted(client *ethclient.Client, nonce uint64, ks *keystore.KeyStore, from, to accounts.Account, val *big.Int, gasLimit uint64) {
	var err error
	var gasPrice, chainID *big.Int
	var signedTx *types.Transaction

	// get gas price
	if gasPrice, err = client.SuggestGasPrice(context.Background()); err != nil {
		log.Fatal(err)
	}
	// get chainID
	if chainID, err = client.ChainID(context.Background()); err != nil {
		log.Fatal(err)
	}

	// dummy encrypted tx
	addr := common.HexToAddress("0x0000000000000000000000000000000000000001")
	accesses := types.AccessList{types.AccessTuple{
		Address: addr,
		StorageKeys: []common.Hash{
			{0},
		},
	}}
	enc := &types.EncryptedTx{
		ChainID:    chainID,
		Nonce:      nonce,
		GasFeeCap:  gasPrice,
		GasTipCap:  big.NewInt(10),
		Gas:        gasLimit,
		To:         &to.Address,
		Value:      val,
		Data:       []byte("encrytped: hello world"),
		AccessList: accesses,
	}
	tx := types.NewTx(enc)

	// sign
	if signedTx, err = ks.SignTx(from, tx, chainID); err != nil {
		log.Fatal(err)
	}

	if err = client.SendTransaction(context.Background(), signedTx); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Encrypted Transaction send: %v\n", signedTx.Hash().Hex())
}

func prettyPrintBlock(client *ethclient.Client, num *big.Int) {
	var block *types.Block
	var err error
	if block, err = client.BlockByNumber(context.Background(), num); err != nil {
		log.Fatal(err)
	}
	fmt.Println("block.Number: ", block.Number().Uint64())
	fmt.Println("block.Time: ", block.Time())
	fmt.Println("block.Difficulty: ", block.Difficulty().Uint64())
	fmt.Println("block.Hash: ", block.Hash().Hex())
	fmt.Println("block.Transactions num: ", len(block.Transactions()))

	for _, tx := range block.Transactions() {
		fmt.Println("tx.Hash: ", tx.Hash().Hex())
		fmt.Println("tx.Value: ", tx.Value().String())
		fmt.Println("tx.Gas: ", tx.Gas())
		fmt.Println("tx.GasPrice: ", tx.GasPrice().Uint64())
		fmt.Println("tx.Nonce: ", tx.Nonce())
		fmt.Println("tx.Data: ", tx.Data())
		fmt.Println("tx.To: ", tx.To().Hex())
	}
}

func main() {
	var client *ethclient.Client
	var err error
	if client, err = ethclient.Dial("//./pipe/geth.ipc"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection established")

	// unlock the pre fund user account
	ks := keystore.NewKeyStore("../.ethereum/keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	if err = ks.Unlock(ks.Accounts()[1], ""); err != nil {
		log.Fatal(err)
	}

	// print accounts
	for i, item := range ks.Accounts() {
		fmt.Printf("Accounts[%v]: %v\n", i, item.Address.Hex())
	}

	authority := ks.Accounts()[0].Address
	user := ks.Accounts()[1].Address

	user_acc := accounts.Account{Address: user}
	auth_acc := accounts.Account{Address: authority}

	wei_user, eth_user := get_balance(client, user)
	wei_auth, eth_auth := get_balance(client, authority)

	fmt.Printf("Before User: %s, %s\nAuthority: %s, %s\n", wei_user.String(), eth_user.String(), wei_auth.String(), eth_auth.String())

	gasLimit := uint64(210000)
	val := big.NewInt(1e18)
	mix := flag.Bool("mix", false, "if send an encrypted transaction")
	amount := flag.Uint64("batch", 10, "batch size of tx")
	flag.Parse()

	var nonce uint64

	var (
		total_amounts uint64
		i             uint64
	)

	for round := 0; round < 10; round++ {
		// total_amounts = rand.Uint64() % 1000
		total_amounts = *amount
		i = 0

		// get nonce
		if nonce, err = client.PendingNonceAt(context.Background(), user_acc.Address); err != nil {
			log.Fatal(err)
		}

		for i = 0; i < total_amounts; i++ {
			var r uint32
			if *mix {
				r = rand.Uint32() % 2
			} else {
				r = 0
			}

			if r == 1 {
				sendEtherPlaintext(client, nonce+i, ks, user_acc, auth_acc, val, gasLimit)
			} else {
				sendEtherEncrypted(client, nonce+i, ks, user_acc, auth_acc, val, gasLimit)
			}
		}

		time.Sleep(40 * time.Second)
	}

	wei_user, eth_user = get_balance(client, user)
	wei_auth, eth_auth = get_balance(client, authority)

	fmt.Printf("After User: %s, %s\nAuthority: %s, %s\n", wei_user.String(), eth_user.String(), wei_auth.String(), eth_auth.String())

}
