#!/bin/sh

set -e

tmux set -g remain-on-exit failed ||
tmux set -g remain-on-exit on

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

export F3B_PROTOCOL=tpke

export MNEMONIC="candy maple cake sugar pudding cream honey rich smooth crumble sweet treat"
go run ./script/write_genesis > $tempdir/clique.json

geth -datadir "$producer_datadir" -verbosity 1 init $tempdir/clique.json

case "$F3B_PROTOCOL" in
	tpke | tibe )
tmux neww -d env LLVL=info smccli --config $tempdir/dela/node1 start --routing tree --listen tcp://127.0.0.1:2001
tmux neww -d env LLVL=info smccli --config $tempdir/dela/node2 start --routing tree --listen tcp://127.0.0.1:2002
tmux neww -d env LLVL=info smccli --config $tempdir/dela/node3 start --routing tree --listen tcp://127.0.0.1:2003
sleep 1

smccli --config $tempdir/dela/node2 minogrpc join --address //127.0.0.1:2001 $(smccli --config $tempdir/dela/node1 minogrpc token)
smccli --config $tempdir/dela/node3 minogrpc join --address //127.0.0.1:2001 $(smccli --config $tempdir/dela/node1 minogrpc token)
                                   
smccli --config $tempdir/dela/node1 dkg listen
smccli --config $tempdir/dela/node2 dkg listen
smccli --config $tempdir/dela/node3 dkg listen

smccli --config $tempdir/dela/node1 dkg setup \
    --authority $(cat $tempdir/dela/node1/dkgauthority) \
    --authority $(cat $tempdir/dela/node2/dkgauthority) \
    --authority $(cat $tempdir/dela/node3/dkgauthority) \
    --threshold 2
export F3B_DKG_PATH=$tempdir/dela/node1
    ;;
esac


cp keystore/$coinbase $producer_datadir/keystore
tmux neww -d env F3B_DKG_PATH="$F3B_DKG_PATH" geth -datadir "$producer_datadir" -nodiscover -mine -password /dev/null -unlock $coinbase -nodekeyhex $producer_nodekey -nat none

observer_datadir=$tempdir/observer
geth -datadir "$observer_datadir" -verbosity 1 init $tempdir/clique.json
tmux neww -d geth -datadir "$observer_datadir" -http -port 0 -authrpc.port 0 -bootnodes $producer_addr
export ETH_RPC_URL=http://localhost:8545

# wait for geth to start
while ! cast block-number 2> /dev/null; do
	sleep 1
done

export ADDRESSES_FILE=$tempdir/addresses
(cd contracts
	visibly 'forge script --keystore "$ETH_KEYSTORE/$deployer" --sender $deployer -f $ETH_RPC_URL --broadcast script/Setup.s.sol'
)
auctions_address=$(jq -r .auctions <$ADDRESSES_FILE)
weth_address=$(jq -r .weth <$ADDRESSES_FILE)
collection_address=$(jq -r .collection <$ADDRESSES_FILE)

visibly 'go run ./script/auction_scenario'

auction_id=0 # FIXME: hardcoded
cast call --trace $auctions_address 'settle(uint256)' $auction_id
