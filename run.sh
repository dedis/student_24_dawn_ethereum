#!/bin/sh

set -e

tempdir=$(mktemp -dt f3b.XXXXXX)
cleanup() {
        pkill -P $$ # kill all child processes
        rm -rf $tempdir
}
trap cleanup EXIT

coinbase=0x280F6B48E4d9aEe0Efdb04EeBe882023357f6434
deployer=$coinbase
address1=0xF5f341CD21350259A8666B3A5fE47132efF57838
address2=0xa9ca84343c8dB08d596400d35A7034027A5F4b31
export ETH_KEYSTORE="$(pwd)/keystore"
touch $tempdir/password
export ETH_PASSWORD="$tempdir/password"

verbosely() {
	echo
	echo $'\e[1m$' "$*"$'\e[0m'
	eval "$@"
}

export GETH_DATADIR=$tempdir/ethereum
export ETH_RPC_URL=http://localhost:8545

verbosely 'geth -datadir "$GETH_DATADIR" init clique.json'

cp keystore/$coinbase $GETH_DATADIR/keystore
tmux neww -d geth -datadir "$GETH_DATADIR" -dev -http

# wait for geth to start
while ! cast block-number 2> /dev/null; do
	sleep 1
done

(cd contracts
	verbosely 'forge create --keystore "$ETH_KEYSTORE/$deployer" --from $deployer Auction'
)

auction_contract=0xef434c1405f66997CBf4a04FDDed518C28a6a013

verbosely 'cast send --async --keystore $ETH_KEYSTORE/$deployer --from $deployer $auction_contract "start()"'

# send an encrypted bid
verbosely 'go run ./script/send_enc -sender $address1 -value 1 $auction_contract $(cast sig "bid()")'
#verbosely 'go run ./script/send_enc -sender $address2 -value 2 $auction_contract $(cast sig "bid()")'

sleep 40

(cd contracts
	verbosely 'forge script --broadcast -f $ETH_RPC_URL --keystore "$ETH_KEYSTORE/$deployer" --sender $deployer script/CloseAuction.s.sol'
)

bash
