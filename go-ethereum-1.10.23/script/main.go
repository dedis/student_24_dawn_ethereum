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

	for i, item := range ks.Accounts() {
		fmt.Printf("Accounts[%v]: %v\n", i, item.Address.Hex())
	}

	authority := ks.Accounts()[0].Address
	user := ks.Accounts()[1].Address
	// auth_acc := accounts.Account{Address: authority}

	wei_user, eth_user := get_balance(client, user)
	wei_auth, eth_auth := get_balance(client, authority)

	fmt.Printf("Before User: %s, %s\nAuthority: %s, %s\n", wei_user.String(), eth_user.String(), wei_auth.String(), eth_auth.String())

	// send ether from user to authority
	// user_acc := ks.Accounts()[1]
	// nonce_user, err := client.PendingNonceAt(context.Background(), user)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // tx params
	// val := big.NewInt(1e18)
	// gasLimit := uint64(210000)
	// gasPrice, err := client.SuggestGasPrice(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// chainID, err := client.NetworkID(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // legacy tx: send ether to auth
	// // tx := types.NewTransaction(nonce_user, authority, val, gasLimit, gasPrice, nil)

	// fmt.Println("point0", chainID)

	// // dummy encrypted tx
	// addr := common.HexToAddress("0x0000000000000000000000000000000000000001")
	// accesses := types.AccessList{types.AccessTuple{
	// 	Address: addr,
	// 	StorageKeys: []common.Hash{
	// 		{0},
	// 	},
	// }}
	// enc := &types.EncryptedTx{
	// 	ChainID:    big.NewInt(42),
	// 	Nonce:      nonce_user,
	// 	GasFeeCap:  gasPrice,
	// 	GasTipCap:  big.NewInt(10),
	// 	Gas:        gasLimit,
	// 	To:         &authority,
	// 	Value:      val,
	// 	Data:       []byte{},
	// 	AccessList: accesses,
	// }

	// tx := types.NewTx(enc)
	// signedTx, err := ks.SignTx(user_acc, tx, chainID)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("point1")

	// err = client.SendTransaction(context.Background(), signedTx)

	// fmt.Println("point2")

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex())

	// // get block and tx
	blockNumber := big.NewInt(3)
	fmt.Println("point0")
	hd, err := client.HeaderByNumber(context.Background(), blockNumber)
	// block, err := client.BlockByNumber(context.Background(), blockNumber)
	fmt.Printf("# header: %v, %v, parent: %v, root: %v\n", hd.Number, hd.TxHash, hd.ParentHash, hd.Root)
	// fmt.Println("point1")
	// fmt.Printf("fatal err: %v\n", err)
	if err != nil {
		log.Fatal(err)
	}

	txHash := common.HexToHash("0xbd917bdb05d6e174c6fec4e2d54ae4d419b077113c60c78cb51d7eebdca0a425")
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tx.Hash().Hex()) // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
	fmt.Println(isPending)       // false
	// fmt.Println("point2")
	// fmt.Println(block.Number().Uint64())     // 5671744
	// fmt.Println(block.Time())                // 1527211625
	// fmt.Println(block.Difficulty().Uint64()) // 3217000136609065
	// fmt.Println(block.Hash().Hex())          // 0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9
	// fmt.Println(len(block.Transactions()))   // 144
	// fmt.Println("point3")
	// for _, tx := range block.Transactions() {
	// 	fmt.Println(tx.Hash().Hex())        // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
	// 	fmt.Println(tx.Value().String())    // 10000000000000000
	// 	fmt.Println(tx.Gas())               // 105000
	// 	fmt.Println(tx.GasPrice().Uint64()) // 102000000000
	// 	fmt.Println(tx.Nonce())             // 110644
	// 	fmt.Println(tx.Data())              // []
	// 	fmt.Println(tx.To().Hex())          // 0x55fE59D8Ad77035154dDd0AD0388D09Dd4047A8e
	// }
	// fmt.Println("point4")

	wei_user, eth_user = get_balance(client, user)
	wei_auth, eth_auth = get_balance(client, authority)

	fmt.Printf("After User: %s, %s\nAuthority: %s, %s\n", wei_user.String(), eth_user.String(), wei_auth.String(), eth_auth.String())

}
