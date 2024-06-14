#!/usr/bin/env bash

set -e

. script/prepare_smc.sh

cp keystore/$coinbase $producer_datadir/keystore
tmux splitw -hd geth -datadir "$producer_datadir" -nodiscover -mine -password /dev/null -unlock $coinbase -nodekeyhex $producer_nodekey -nat none -http -ws -allow-insecure-unlock

export ETH_RPC_URL=http://localhost:8545

# wait for geth to start
while ! cast block-number 2> /dev/null; do
	sleep 1
done

export ADDRESSES_FILE=$tempdir/addresses
(cd contracts
	F3B_PROTOCOL=$protocol F3B_BLOCKDELAY=$blockdelay visibly 'forge script --keystore "$ETH_KEYSTORE/$deployer" --sender $deployer -f $ETH_RPC_URL --broadcast script/Setup.s.sol'
)
auctions_address=$(jq -r .auctions <$ADDRESSES_FILE)
weth_address=$(jq -r .weth <$ADDRESSES_FILE)
collection_address=$(jq -r .collection <$ADDRESSES_FILE)

visibly 'go run ./script/auction_scenario'
