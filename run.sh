#!/bin/sh

set -e

tempdir=$(mktemp -dt f3b.XXXXXX)
cleanup() {
        pkill -P $$ # kill all child processes
        rm -rf $tempdir
}
trap cleanup EXIT

deployer=0xF5f341CD21350259A8666B3A5fE47132efF57838
address1=0x280F6B48E4d9aEe0Efdb04EeBe882023357f6434
address2=0xa9ca84343c8dB08d596400d35A7034027A5F4b31
export ETH_KEYSTORE="$(pwd)/keystore"
touch $tempdir/password
export ETH_PASSWORD="$tempdir/password"

(cd contracts
forge compile
)

tmux neww env LLVL=info dkgcli --config $tempdir/dela/node1 start --routing tree --listen tcp://127.0.0.1:2001
tmux neww env LLVL=info dkgcli --config $tempdir/dela/node2 start --routing tree --listen tcp://127.0.0.1:2002
tmux neww env LLVL=info dkgcli --config $tempdir/dela/node3 start --routing tree --listen tcp://127.0.0.1:2003
sleep 1

dkgcli --config $tempdir/dela/node2 minogrpc join --address //127.0.0.1:2001 $(dkgcli --config $tempdir/dela/node1 minogrpc token)
dkgcli --config $tempdir/dela/node3 minogrpc join --address //127.0.0.1:2001 $(dkgcli --config $tempdir/dela/node1 minogrpc token)
                                   
dkgcli --config $tempdir/dela/node1 dkg listen
dkgcli --config $tempdir/dela/node2 dkg listen
dkgcli --config $tempdir/dela/node3 dkg listen

dkgcli --config $tempdir/dela/node1 dkg setup \
    --authority $(cat $tempdir/dela/node1/dkgauthority) \
    --authority $(cat $tempdir/dela/node2/dkgauthority) \
    --authority $(cat $tempdir/dela/node3/dkgauthority) \
    --threshold 2

export F3B_DKG_PATH=$tempdir/dela/node1
export GETH_DATADIR=$tempdir/ethereum
export ETH_RPC_URL=http://localhost:8545

go install github.com/ethereum/go-ethereum/cmd/geth

geth -datadir "$GETH_DATADIR" init clique.json

cp -R keystore/* -t "$GETH_DATADIR/keystore/"
tmux neww env F3B_DKG_PATH="$F3B_DKG_PATH" geth -datadir "$GETH_DATADIR" --nodiscover --http -dev --mine -verbosity 4

# wait for geth to start
while ! cast block-number 2> /dev/null; do
sleep 1
done

(cd contracts
forge script --broadcast --legacy -f $ETH_RPC_URL --keystore "$ETH_KEYSTORE/$deployer" --sender $deployer script/Deploy.s.sol
)

weth_contract=0x8bD0539849AaA50C85f418825a1962208F5C30fA
auction_contract=0x3712327B0E9fAE301cFED65eD6BDEf03629fCCFa

cast send --async --legacy --keystore $ETH_KEYSTORE/$deployer --from $deployer $auction_contract 'start()'

# send an encrypted bid
go run ./script/send_enc -sender $address1 -value 1 $auction_contract $(cast sig 'bid()')
go run ./script/send_enc -sender $address2 -value 2 $auction_contract $(cast sig 'bid()')

sleep 40

(cd contracts
forge script --broadcast --legacy -f $ETH_RPC_URL --keystore "$ETH_KEYSTORE/$deployer" --sender $deployer script/CloseAuction.s.sol
)

bash
