package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math"
	"math/big"
	"os/exec"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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

func sendEtherF3bEnc(client *ethclient.Client, ks *keystore.KeyStore, from, to accounts.Account, val *big.Int, gasLimit uint64) {
	var nonce uint64
	var err error
	var gasPrice, chainID *big.Int
	var signedTx *types.Transaction

	// get nonce
	if nonce, err = client.PendingNonceAt(context.Background(), from.Address); err != nil {
		log.Fatal(err)
	}
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

	// dkgcli --config ./tmp/node1 dkg setup
	// --authority RjEyNy4wLjAuMToyMDAx:r7I444R4a+sh+a0n6PV/FBRDRho702USB9MCxNWkTOk=
	// --authority RjEyNy4wLjAuMToyMDAy:iVZVJ4vvL3We94Y75eG23SsXOOBgTupi1eKzeg66BbE=
	// --authority RjEyNy4wLjAuMToyMDAz:9mZrUFDI2c6LML1t8qGCIw3hDCrYormquNxBbgcEaNg=

	pt := "D:/EPFL/master_thesis/dela/dkg/pedersen/dkgcli/tmp/node1/"

	node := filepath.Dir(pt)

	// auth := []string{"", "", ""}

	// args := []string{"dkgcli", "--config", node, "dkg", "setup", "--authority", auth[0], "--authority", auth[1], "--authority", auth[2]}

	// cmd := strings.Join(args[:], " ")

	plaintext := "dddddddd"

	args_enc := []string{"dkgcli", "--config", node, "dkg", "encrypt", "--message", plaintext}

	encrypted_data, err := exec.Command(args_enc[0], args_enc[1:]...).Output()

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("## Encrypted data: ", string(encrypted_data))

	args_dec := []string{"dkgcli", "--config", node, "dkg", "decrypt", "--encrypted", string(encrypted_data)}

	decrypted_data, err := exec.Command(args_enc[0], args_dec[1:]...).Output()

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("## Decrypted data: ", string(decrypted_data))

	enc := &types.EncryptedTx{
		ChainID:    chainID,
		Nonce:      nonce,
		GasFeeCap:  gasPrice,
		GasTipCap:  big.NewInt(10),
		Gas:        gasLimit,
		To:         &to.Address,
		Value:      val,
		Data:       encrypted_data,
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

func sendEtherF3bVerifiedEnc(client *ethclient.Client, nonce uint64, ks *keystore.KeyStore, from, to accounts.Account, val *big.Int, gasLimit uint64) {
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

	// dkgcli --config ./tmp/node1 dkg setup
	// --authority RjEyNy4wLjAuMToyMDAx:r7I444R4a+sh+a0n6PV/FBRDRho702USB9MCxNWkTOk=
	// --authority RjEyNy4wLjAuMToyMDAy:iVZVJ4vvL3We94Y75eG23SsXOOBgTupi1eKzeg66BbE=
	// --authority RjEyNy4wLjAuMToyMDAz:9mZrUFDI2c6LML1t8qGCIw3hDCrYormquNxBbgcEaNg=

	pt := "D:/EPFL/master_thesis/dela/dkg/pedersen/dkgcli/tmp/node1/"

	node := filepath.Dir(pt)

	// auth := []string{"", "", ""}

	// args := []string{"dkgcli", "--config", node, "dkg", "setup", "--authority", auth[0], "--authority", auth[1], "--authority", auth[2]}

	// cmd := strings.Join(args[:], " ")

	msg := "Merry Christmas!"
	fmt.Println("## Plaintext message: ", msg)
	symKey := make([]byte, types.SymKeyLen)
	_, err = rand.Read(symKey)
	if err != nil {
		panic(fmt.Sprintf("failed on load random key: %v", err))
	}

	symKeyStr := hex.EncodeToString(symKey)
	// msgStr := hex.EncodeToString([]byte(msg))

	encryptedMsg := crypto.EncryptAES(symKey, []byte(msg))
	// compute hash of cleartext key
	khash := sha256.Sum256([]byte(symKey))
	khash[0] = 0

	gBar := "1d0194fdc2fa2ffcc041d3ff12045b73c86e4ff95ff662a5eee82abdf44a53c7"

	args_enc := []string{"dkgcli", "--config", node, "dkg", "verifiableEncrypt", "--GBar", gBar, "--message", symKeyStr}

	encrypted_data, err := exec.Command(args_enc[0], args_enc[1:]...).Output()

	encrypted_data = encrypted_data[:len(encrypted_data)-2]

	encKeyWithHash := append(khash[:], encrypted_data...)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("## Random generated msg.Key: ", symKeyStr)
	fmt.Println("## Hash: ", hex.EncodeToString(khash[:]))
	fmt.Println("## Encrypted Key: ", string(encrypted_data))
	fmt.Println("## hash | encrypted key: ", hex.EncodeToString(encKeyWithHash))

	// TODO: uncomment following to display decryption result
	// args_dec := []string{"dkgcli", "--config", node, "dkg", "verifiableDecrypt", "--GBar", gBar, "--ciphertexts", string(encrypted_data)}

	// fmt.Printf("args to decrypt: %v\n", args_dec)

	// decrypted_data, err := exec.Command(args_enc[0], args_dec[1:]...).Output()

	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println("## Decrypted data: ", string(decrypted_data))

	weiValue := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))
	weiFloat := new(big.Float).SetInt(weiValue)
	ethValue := new(big.Float).Quo(weiFloat, big.NewFloat(math.Pow10(18)))
	fmt.Println("[GasLimit in Ether]: ", ethValue.String())

	enc := &types.EncryptedTx{
		ChainID:    chainID,
		Nonce:      nonce,
		GasFeeCap:  gasPrice,
		GasTipCap:  big.NewInt(10),
		Gas:        gasLimit,
		To:         &to.Address,
		Value:      val,
		Key:        encKeyWithHash,
		Data:       encryptedMsg,
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
	encrypted := flag.Bool("encrypted", false, "if send an encrypted transaction")
	num := flag.Int("num", 1, "number of transactions")
	cid := flag.String("id", "", "id of geth client")
	flag.Parse()

	var client *ethclient.Client
	var err error
	if client, err = ethclient.Dial(fmt.Sprintf("pipe/geth%s.ipc", *cid)); err != nil {
		log.Fatal(err)
	}

	// unlock the pre fund user account
	ks := keystore.NewKeyStore("../.ethereum/keystore", keystore.StandardScryptN, keystore.StandardScryptP)
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
		if !*encrypted {
			sendEtherPlaintext(client, nonce+uint64(i), ks, user_acc, rcv_acc, val, gasLimit)
		} else {
			sendEtherF3bVerifiedEnc(client, nonce+uint64(i), ks, user_acc, rcv_acc, val, gasLimit)
		}
	}

	fmt.Printf("[Balance before tx exec] User: %s, Receiver: %s, Auth1: %s, Auth2: %s\n",
		eth_user.String(), eth_rcv.String(), eth_auth1.String(), eth_auth2.String())

}
