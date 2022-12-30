package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const authority_file = "../.ethereum/keystore/UTC--2022-09-13T11-34-29.303731400Z--280f6b48e4d9aee0efdb04eebe882023357f6434"
const user_file = "../.ethereum/keystore/UTC--2022-09-13T11-35-11.765870700Z--f5f341cd21350259a8666b3a5fe47132eff57838"

func get_balance(client *ethclient.Client, addr common.Address, bn *big.Int) (*big.Int, *big.Float) {
	weiValue, err := client.BalanceAt(context.Background(), addr, bn)
	if err != nil {
		log.Fatal(err)
	}
	fbalance := new(big.Float)
	fbalance.SetString(weiValue.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	return weiValue, ethValue
}

func main() {
	var client *ethclient.Client
	var err error
	cid := flag.String("id", "", "id of geth client")
	bn := flag.Int64("bn", 0, "block number") // 0 as default to query latest block
	flag.Parse()
	if client, err = ethclient.Dial(fmt.Sprintf("//./pipe/geth%s.ipc", *cid)); err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Connection established")

	// unlock the pre fund user account
	ks := keystore.NewKeyStore("../.ethereum/keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	if err = ks.Unlock(ks.Accounts()[1], ""); err != nil {
		log.Fatal(err)
	}

	// create receiver account
	rcv := "be76d61cd0a253d6ce6d363966abbde3b7b5a6a7"
	rcvHex, err := hex.DecodeString(rcv)
	if err != nil {
		panic(err)
	}
	rcvAddr := common.BytesToAddress(rcvHex)

	auth1 := ks.Accounts()[0].Address
	user := ks.Accounts()[1].Address
	auth2 := ks.Accounts()[2].Address

	var num *big.Int
	if *bn != 0 {
		num = big.NewInt(*bn)
	}

	_, eth_user := get_balance(client, user, num)
	_, eth_rcv := get_balance(client, rcvAddr, num)
	_, eth_auth1 := get_balance(client, auth1, num)
	_, eth_auth2 := get_balance(client, auth2, num)

	fmt.Printf("[Balance at block %v] User: %s, Receiver: %s, Auth1: %s, Auth2: %s\n",
		*bn, eth_user.String(), eth_rcv.String(), eth_auth1.String(), eth_auth2.String())
}
