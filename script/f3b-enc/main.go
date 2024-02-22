package main

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math"
	"math/big"
	"path"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/f3b"
)

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

func sendEtherF3bEnc(client *ethclient.Client, nonce uint64, ks *keystore.KeyStore, from, to accounts.Account, val *big.Int, gasLimit uint64) {
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

	dkgcli := f3b.NewDkgCli()

	plaintext := []byte("dddddddd")
	label := from.Address.Bytes()
	label = binary.BigEndian.AppendUint64(label, nonce)

	encrypted_data, err := dkgcli.Encrypt(label, plaintext)
	if err != nil {
		log.Fatal(err)
	}

	enc := &types.EncryptedTx{
		ChainID:    chainID,
		Nonce:      nonce,
		GasFeeCap:  gasPrice,
		GasTipCap:  big.NewInt(10),
		Gas:        gasLimit,
		To:         &to.Address,
		Value:      val,
		Data:       encrypted_data,
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

// hd, err := client.HeaderByNumber(context.Background(), blockNumber)
// fmt.Printf("# header: %v, %v, parent: %v, root: %v\n", hd.Number, hd.TxHash, hd.ParentHash, hd.Root)
// if err != nil {
// 	log.Fatal(err)
// }

// txHash := common.HexToHash("0xbd917bdb05d6e174c6fec4e2d54ae4d419b077113c60c78cb51d7eebdca0a425")
// tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
// if err != nil {
// 	log.Fatal(err)
// }

// fmt.Println(tx.Hash().Hex()) // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
// fmt.Println(isPending)       // false

func main() {
	num := flag.Int("num", 1, "number of transactions")
	gethDir := flag.String("ethdir", ".ethereum/", "geth client directory")
	flag.Parse()

	var client *ethclient.Client
	var err error
	if client, err = ethclient.Dial(path.Join(*gethDir, "geth.ipc")); err != nil {
		log.Fatal(err)
	}

	// unlock the pre fund user account
	ks := keystore.NewKeyStore(path.Join(*gethDir, "keystore"), keystore.StandardScryptN, keystore.StandardScryptP)
	if err = ks.Unlock(ks.Accounts()[1], ""); err != nil {
		log.Fatal(err)
	}

	// print accounts
	fmt.Println("Existing accoutns:")
	for i, item := range ks.Accounts() {
		fmt.Printf("Accounts[%v]: %v\n", i, item.Address.Hex())
	}
	auth1 := ks.Accounts()[0].Address
	user := ks.Accounts()[1].Address
	auth2 := ks.Accounts()[2].Address

	// create receiver account
	rcv := "be76d61cd0a253d6ce6d363966abbde3b7b5a6a7"
	rcvHex, err := hex.DecodeString(rcv)
	if err != nil {
		panic(err)
	}
	rcvAddr := common.BytesToAddress(rcvHex)

	user_acc := accounts.Account{Address: user}
	rcv_acc := accounts.Account{Address: rcvAddr}

	_, eth_user := get_balance(client, user)
	_, eth_rcv := get_balance(client, rcvAddr)
	_, eth_auth1 := get_balance(client, auth1)
	_, eth_auth2 := get_balance(client, auth2)

	gasLimit := uint64(210000)
	val := big.NewInt(1e18)

	var nonce uint64

	// get nonce
	if nonce, err = client.PendingNonceAt(context.Background(), user_acc.Address); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[Gas Info] Gas Limit: %v\n", gasLimit)

	for i := 0; i < *num; i++ {
		sendEtherF3bEnc(client, nonce+uint64(i), ks, user_acc, rcv_acc, val, gasLimit)
	}

	fmt.Printf("[Balance before tx exec] User: %s, Receiver: %s, Auth1: %s, Auth2: %s\n",
		eth_user.String(), eth_rcv.String(), eth_auth1.String(), eth_auth2.String())

}
