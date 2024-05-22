#!/bin/sh

jq --version > /dev/null || { echo "Please install jq"; exit 1; }

go install github.com/ethereum/go-ethereum/cmd/geth \
           github.com/ethereum/go-ethereum/cmd/abigen \
           go.dedis.ch/f3b/smc/smccli
