# F3B-Geth
Delayed-execution Ethereum client.
Based on previous work by Shufan Wang.

# Architecture
`dela/` is a modified version of `dela`.
`go-ethereum/` is a modified version of `go-ethereum v1.10.23` which integrates with `dela`.

# Design (WIP)

- [x] When creating a transaction, the user agent IBE-encrypts the calldata with dela.
- [x] The IBE label is sender address concatenated with big-endian 64-bit nonce
- [x] The ciphertext is authenticated (HMAC-SHA256, encrypt-then-MAC)
- [ ] The chain only accepts encrypted transactions
- [ ] The `to` address is encrypted
- [ ] The transaction receipt contains a symmetric encryption key
- [ ] The execution layer can direct the SMC node to release an encryption label only after the transaction is finalized.
- [ ] TDH2, PVSS, beacon IBE options maybe

# Start the Dela nodes
Reference: [README](dela/dkg/pedersen_bn256/dkgcli/README.md)
```sh
go install go.dedis.ch/dela/dkg/pedersen_bn256/dkgcli

LLVL=info dkgcli --config /tmp/dela/node1 start --routing tree --listen tcp://127.0.0.1:2001 &
LLVL=info dkgcli --config /tmp/dela/node2 start --routing tree --listen tcp://127.0.0.1:2002 &
LLVL=info dkgcli --config /tmp/dela/node3 start --routing tree --listen tcp://127.0.0.1:2003 &

dkgcli --config /tmp/dela/node2 minogrpc join --address //127.0.0.1:2001 $(dkgcli --config /tmp/dela/node1 minogrpc token)
dkgcli --config /tmp/dela/node3 minogrpc join --address //127.0.0.1:2001 $(dkgcli --config /tmp/dela/node1 minogrpc token)
                                   
# Initialize DKG on each node. Do that in a 4th session.
dkgcli --config /tmp/dela/node1 dkg listen
dkgcli --config /tmp/dela/node2 dkg listen
dkgcli --config /tmp/dela/node3 dkg listen

# Do the setup in one of the node:
dkgcli --config /tmp/dela/node1 dkg setup \
    --authority $(cat /tmp/dela/node1/dkgauthority) \
    --authority $(cat /tmp/dela/node2/dkgauthority) \
    --authority $(cat /tmp/dela/node3/dkgauthority) \
    --threshold 2


# this is for other commands to be able to communicate
export F3B_DKG_PATH=/tmp/dela/node1
```
# Start the geth node

Please refer to the go-ethereum development [book](https://goethereumbook.org/) for more infomation about the commands.

```sh
go install github.com/ethereum/go-ethereum/cmd/geth

# commands for single geth simulation
geth --datadir .ethereum/ init clique.json

# start node
geth --nodiscover --networkid 42 --datadir .ethereum/ --unlock 0x280F6B48E4d9aEe0Efdb04EeBe882023357f6434 --mine
# password is the empty string
```

# Test Example

1. send two plaintext transaction to your node:
```sh
go run ./script/f3b-enc -num 2
```

2. send an encrypted transaction to your node:
```sh
go run ./script/f3b-enc -encrypted
```

# Config
There is a default parameter that is changeable in the system in `go-ethereum/core/types/transaction.go` as shown below. The EncryptedBlockDelay defines the block delay between the ordering block and the execution block.
```
const EncryptedBlockDelay uint64 = 2
```
