package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/naoina/toml"

	"github.com/ethereum/go-ethereum/f3b"
	"github.com/ethereum/go-ethereum/log"
)

func main() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
	if err := Main(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Main() error {
	file, err := os.Open("params.toml")
	if err != nil {
		return err
	}
	defer file.Close()

	var p f3b.FullParams

	decoder := toml.NewDecoder(file)
	err = decoder.Decode(&p.Params)
	if err != nil {
		return err
	}
	p.SmcPath = os.Getenv("F3B_SMC_PATH")

	file, err = os.Create("params.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(&p)
	if err != nil {
		return err
	}
	return nil
}
