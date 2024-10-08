package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/f3b"
	"github.com/ethereum/go-ethereum/params"
	"github.com/miguelmota/go-ethereum-hdwallet"
)

func main() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
	if err := Main(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Main() error {
	deployer := common.HexToAddress("280f6b48e4d9aee0efdb04eebe882023357f6434")

	mnemonic := os.Getenv("MNEMONIC")
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return err
	}

	alloc := map[common.Address]core.GenesisAccount{
		deployer: {Balance: big.NewInt(params.Ether)},
	}

	p, err := f3b.ReadParams()
	if err != nil {
		return err
	}
	nBidders := p.NumBidders
	amount := new(big.Int).Mul(big.NewInt(11), big.NewInt(params.Ether))
	it := accounts.DefaultIterator(hdwallet.DefaultBaseDerivationPath)
	for i := 0; i < nBidders; i++ {
		account, err := wallet.Derive(it(), true)
		if err != nil {
			return err
		}
		alloc[account.Address] = core.GenesisAccount{Balance: amount}
	}

	genesis := core.Genesis{
		Config: &params.ChainConfig{
			ChainID:             big.NewInt(1337),
			HomesteadBlock:      common.Big0,
			EIP150Block:         common.Big0,
			EIP155Block:         common.Big0,
			EIP158Block:         common.Big0,
			ByzantiumBlock:      common.Big0,
			ConstantinopleBlock: common.Big0,
			PetersburgBlock:     common.Big0,
			IstanbulBlock:       common.Big0,
			BerlinBlock:         common.Big0,
			LondonBlock:         common.Big0,
			Clique: &params.CliqueConfig{
				Period: p.BlockTime,
				Epoch:  30000,
			},
		},
		Nonce:      0,
		Timestamp:  0x63206b0f,
		ExtraData:  common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000280f6b48e4d9aee0efdb04eebe882023357f64340000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		GasLimit:   p.BlockTime * p.GasPerSecond,
		Difficulty: big.NewInt(1),
		Alloc:      alloc,
	}

	if p.Protocol != "" {
		genesis.Config.LausanneBlock = common.Big0
	}

	obj, err := json.Marshal(genesis)
	if err != nil {
		return err
	}
	os.Stdout.Write(obj)
	return nil
}
