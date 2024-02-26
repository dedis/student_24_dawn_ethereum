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
geth -datadir "$GETH_DATADIR" --nodiscover --allow-insecure-unlock --unlock 0x280F6B48E4d9aEe0Efdb04EeBe882023357f6434 --password /dev/null --mine &

# send an encrypted transaction
go run ./script/f3b-enc -ethdir "$GETH_DATADIR"

bash
