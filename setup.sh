#!/bin/sh

status=0
if ! forge --version > /dev/null; then
	echo "Please install foundry (https://getfoundry.sh)"
	status=1
fi
if ! tmux -V > /dev/null; then
	echo "Please install tmux"
	status=1
fi
if ! jq --version > /dev/null; then
	echo "Please install jq"
	status=1
fi
if ! go version > /dev/null; then
	echo "Please install go"
	status=1
else
	go install github.com/ethereum/go-ethereum/cmd/geth \
		github.com/ethereum/go-ethereum/cmd/abigen \
		go.dedis.ch/f3b/smc/smccli \
		|| status=1
fi

exit $status
