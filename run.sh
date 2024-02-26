#!/bin/sh

set -e

tempdir=$(mktemp -dt f3b.XXXXXX)
cleanup() {
        pkill -P $$ # kill all child processes
        rm -rf $tempdir
}
trap cleanup EXIT


LLVL=info dkgcli --config $tempdir/dela/node1 start --routing tree --listen tcp://127.0.0.1:2001 &
LLVL=info dkgcli --config $tempdir/dela/node2 start --routing tree --listen tcp://127.0.0.1:2002 &
LLVL=info dkgcli --config $tempdir/dela/node3 start --routing tree --listen tcp://127.0.0.1:2003 &
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

geth -datadir "$GETH_DATADIR" init clique.json

cp -R .ethereum/keystore/* -t "$GETH_DATADIR/keystore/"
geth -datadir "$GETH_DATADIR" --nodiscover --http --rpc.allow-unprotected-txs --allow-insecure-unlock --mine &
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

# send an encrypted transaction
go run ./script/f3b-enc -ethdir "$GETH_DATADIR"

bash
