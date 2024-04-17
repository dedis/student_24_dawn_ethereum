#!/bin/sh

set -e

tmux set -g remain-on-exit failed

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

visibly() {
	echo
	echo $'\e[1m$' "$*"$'\e[0m'
	eval "$@"
}

producer_datadir=$tempdir/producer
producer_nodekey="e74976d3e1d9069b85d6659038105fe601696a0ddcb63f0407b11328e341a47c"
producer_addr="enode://3d1bb945ae2e250f5fe23f6da3f150b1af4d425bd280bdbfc3e7626ae4625cac2cfb3a59469b67528765a50237c0f434bc3cebcb63118b21949e4139de6b9fb1@127.0.0.1:30303"

geth -datadir "$producer_datadir" -verbosity 1 init clique.json

cp keystore/$coinbase $producer_datadir/keystore
tmux neww -d geth -datadir "$producer_datadir" -nodiscover -mine -password /dev/null -unlock $coinbase -nodekeyhex $producer_nodekey -nat none

observer_datadir=$tempdir/observer
geth -datadir "$observer_datadir" -verbosity 1 init clique.json
tmux neww -d geth -datadir "$observer_datadir" -http -port 0 -authrpc.port 0 -bootnodes $producer_addr -verbosity 4
export ETH_RPC_URL=http://localhost:8545

# wait for geth to start
while ! cast block-number 2> /dev/null; do
	sleep 1
done

(cd contracts
	visibly 'forge create --keystore "$ETH_KEYSTORE/$deployer" --from $deployer Auction'
)

auction_contract=0xef434c1405f66997CBf4a04FDDed518C28a6a013

visibly 'cast send --async --keystore $ETH_KEYSTORE/$deployer --from $deployer $auction_contract "start()"'

# send an encrypted bid
visibly 'go run ./script/send_enc -sender $address1 -value 1 $auction_contract $(cast sig "bid()")'
visibly 'go run ./script/send_enc -sender $address2 -value 2 $auction_contract $(cast sig "bid()")'

sleep 40

(cd contracts
	visibly 'forge script --broadcast -f $ETH_RPC_URL --keystore "$ETH_KEYSTORE/$deployer" --sender $deployer script/CloseAuction.s.sol'
)

bash
