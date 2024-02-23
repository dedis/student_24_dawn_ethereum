package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"math/big"

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
		fmt.Println("tx.Data: ", string(tx.Data()[:]))
		fmt.Println("tx.To: ", tx.To().Hex())
		fmt.Println("tx Key: ", tx.Key())
		rc, _ := client.TransactionReceipt(context.Background(), tx.Hash())
		fmt.Printf("%+v\n", rc)
		fmt.Println("receipt Key: ", rc.Key)
	}
}

func prettyPrintReceipt(client *ethclient.Client, num *big.Int) {
	var block *types.Block
	var err error
	if block, err = client.BlockByNumber(context.Background(), num); err != nil {
		log.Fatal(err)
	}
	rc, err := client.GetAllReceipt(context.Background(), block.Hash())
	fmt.Printf("%+v\n", rc)
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
	id := flag.Int64("id", 0, "the block number")
	who := flag.Int64("who", 1, "the node number")
	flag.Parse()
	var client *ethclient.Client
	var err error
	if client, err = ethclient.Dial(fmt.Sprintf("pipe/geth%v.ipc", *who)); err != nil {
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

	wei_user, eth_user := get_balance(client, user)
	wei_auth, eth_auth := get_balance(client, authority)

	fmt.Printf("Before User: %s, %s\nAuthority: %s, %s\n", wei_user.String(), eth_user.String(), wei_auth.String(), eth_auth.String())

	// encrypted tx: 7
	// plain tx: 5

	num0 := big.NewInt(*id)
	fmt.Println("\n# Block number: ", id, " from ", who)
	prettyPrintBlock(client, num0)
	// prettyPrintReceipt(client, num0)
}
