#!/bin/sh

go install github.com/ethereum/go-ethereum/cmd/geth \
           github.com/ethereum/go-ethereum/cmd/abigen \
           go.dedis.ch/f3b/smc/smccli

(cd contracts > /dev/null; forge inspect WETH abi) | abigen -abi - -pkg bindings -type WETH -out bindings/weth.go
