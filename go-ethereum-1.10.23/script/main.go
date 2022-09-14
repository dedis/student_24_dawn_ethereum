package main

import (
	"context"
	"fmt"
	"io/ioutil"
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

func importKs(file string) *keystore.KeyStore {
	// file := "./tmp/UTC--2018-07-04T09-58-30.122808598Z--20f8d42fb0f667f2e53930fed426f225752453b3"
	ks := keystore.NewKeyStore("../.ethereum/keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	jsonBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	password := ""
	account, err := ks.Import(jsonBytes, password, password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(account.Address.Hex())

	return ks

	// if err := os.Remove(file); err != nil {
	// 	log.Fatal(err)
	// }
}

// func dummy_enc_tx(nonce uint64, gasPrice *big.Int, gasLimit uint64) (tx *types.Transaction) {
// 	return tx
// }

func main() {
	client, err := ethclient.Dial("//./pipe/geth.ipc")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("we have a connection")

	// unlock the pre fund user account
	ks := keystore.NewKeyStore("../.ethereum/keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	// ks := importKs(user_file)
	err = ks.Unlock(ks.Accounts()[1], "")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("keys of accounts: %v\n", ks.Accounts())

	authority := ks.Accounts()[0].Address
	user := ks.Accounts()[1].Address
	// auth_acc := accounts.Account{Address: authority}
	user_acc := ks.Accounts()[1]

	wei_user, eth_user := get_balance(client, user)
	wei_auth, eth_auth := get_balance(client, authority)

	fmt.Printf("User: %s, %s\nAuthority: %s, %s\n", wei_user.String(), eth_user.String(), wei_auth.String(), eth_auth.String())

	// send ether from user to authority
	nonce_user, err := client.PendingNonceAt(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	// tx params
	val := big.NewInt(1e18)
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// legacy tx: send ether to auth
	// tx := types.NewTransaction(nonce_user, authority, val, gasLimit, gasPrice, nil)

	// dummy encrypted tx
	enc := &types.EncryptedTx{
		ChainID:   big.NewInt(42),
		Nonce:     nonce_user,
		GasFeeCap: gasPrice,
		Gas:       gasLimit,
		To:        &authority,
		Value:     val,
		Data:      nil,
	}
	tx := types.NewTx(enc)

	signedTx, err := ks.SignTx(user_acc, tx, chainID)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex())

	wei_user, eth_user = get_balance(client, user)
	wei_auth, eth_auth = get_balance(client, authority)

	fmt.Printf("User: %s, %s\nAuthority: %s, %s\n", wei_user.String(), eth_user.String(), wei_auth.String(), eth_auth.String())

}
