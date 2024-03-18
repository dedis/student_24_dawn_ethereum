#!/bin/sh

set -e

tempdir=$(mktemp -dt f3b.XXXXXX)
cleanup() {
        pkill -P $$ # kill all child processes
        rm -rf $tempdir
}
trap cleanup EXIT

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
tmux neww env F3B_DKG_PATH="$F3B_DKG_PATH" geth -datadir "$GETH_DATADIR" --nodiscover --http --rpc.allow-unprotected-txs --allow-insecure-unlock --mine -verbosity 4
sleep 1 # give geth time to start

geth attach  -datadir "$GETH_DATADIR" <<END
// unlock all accounts (empty password)
for (i in eth.accounts) {
  personal.unlockAccount(eth.accounts[i], "", 0)
}
END

(cd contracts
forge script --broadcast --legacy --unlocked -f http://localhost:8545 --sender 0xF5f341CD21350259A8666B3A5fE47132efF57838 script/Deploy.s.sol
)

cast send --async --legacy --unlocked -f 0xF5f341CD21350259A8666B3A5fE47132efF57838 0x3712327B0E9fAE301cFED65eD6BDEf03629fCCFa 'start()'

# send an encrypted bid
go run ./script/send_enc -ethdir "$GETH_DATADIR" --sender 0xF5f341CD21350259A8666B3A5fE47132efF57838 -value 1 0x3712327B0E9fAE301cFED65eD6BDEf03629fCCFa $(cast sig 'bid()')
go run ./script/send_enc -ethdir "$GETH_DATADIR" -value 2 0x3712327B0E9fAE301cFED65eD6BDEf03629fCCFa $(cast sig 'bid()')

sleep 40

(cd contracts
forge script --broadcast --legacy --unlocked -f http://localhost:8545 --sender 0xF5f341CD21350259A8666B3A5fE47132efF57838 script/CloseAuction.s.sol
)

bash
