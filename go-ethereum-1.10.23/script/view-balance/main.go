package main

import (
	"context"
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

func main() {
	var client *ethclient.Client
	var err error
	if client, err = ethclient.Dial("//./pipe/geth1.ipc"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection established")

	// unlock the pre fund user account
	ks := keystore.NewKeyStore("../.ethereum/keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	if err = ks.Unlock(ks.Accounts()[1], ""); err != nil {
		log.Fatal(err)
	}

	// print accounts
	// for i, item := range ks.Accounts() {
	// 	fmt.Printf("Accounts[%v]: %v\n", i, item.Address.Hex())
	// }

	authority := ks.Accounts()[0].Address
	user := ks.Accounts()[1].Address

	_, eth_user := get_balance(client, user)
	_, eth_auth := get_balance(client, authority)

	fmt.Printf("[Balance] User: %s\n", eth_user.String())

	fmt.Printf("[Balance] Authority: %s\n", eth_auth.String())
}
