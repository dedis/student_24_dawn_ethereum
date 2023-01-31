# F3B-Geth
Execution layer based front-running protection on Ethereum (Master thesis) at DEDIS EPFL by Shufan Wang.

Geth with F3B front-running protection achieved in execution layer based on go-ethereum v1.10.23.

# Config before use
There are several default parameters that are changeable in the system in core/types/transactions.go as shown below. The EncryptedBlockDelay defines the block delay between the ordering block and the execution block. The GBar is a publicly known constant we directly take from Dela dkg. The NodePath specifies the path to the SMC node directory which would be used to launch the verifiable encryption. Make sure you adapt NodePath according to your machine.
```
const EncryptedBlockDelay uint64 = 2

const GBar string = "1d0194fdc2fa2ffcc041d3ff12045b73c86e4ff95ff662a5eee82abdf44a53c7"

const NodePath string = "D:/EPFL/master\_thesis/dela/dkg/pedersen/dkgcli/tmp/node1/"
```

# Start the geth node

Please refer to the go-ethereum development [book](https://goethereumbook.org/) for more infomation about the commands.

```
// build geth
go install -v ./cmd/geth


// Commands for single geth simulation:
geth --datadir .ethereum/ init clique.json

// start node
geth --nodiscover --networkid 42 --datadir .ethereum/ --unlock 0x280F6B48E4d9aEe0Efdb04EeBe882023357f6434 --mine


// Commands for multiple geth simulation:
geth --datadir .ethereum1/ init ../multiclique.json
geth --datadir .ethereum2/ init ../multiclique.json

// start node1
geth --nodiscover --networkid 43 --datadir .ethereum1/ --unlock 0x280F6B48E4d9aEe0Efdb04EeBe882023357f6434 --mine --ipcpath \\.\pipe\geth1.ipc --authrpc.port 8551 --port 30303

// start node2
geth --nodiscover --networkid 43 --datadir .ethereum2/ --unlock 0xa9ca84343c8dB08d596400d35A7034027A5F4b31 --mine --ipcpath \\.\pipe\geth2.ipc --authrpc.port 8552 --port 30304 --miner.etherbase 0xa9ca84343c8dB08d596400d35A7034027A5F4b31 --syncmode full

// addPeers
geth attach \\.\pipe\geth1.ipc

admin.addPeer

```

# Start the Dela nodes
Please follow the [instructions](https://github.com/Mahsa-Bastankhah/dela/tree/5593c8d782ae14910343212447956d8b46ea958b/dkg/pedersen/dkgcli) to run several F3B committee members.

# Test Example

1. send two plaintext transaction to node1:
```go run script/f3b-enc/main.go -num 2 -id 1```

1. send one encrypted transaction to node2: 
```go run script/f3b-enc/main.go -num 1 -id 2 -encrypted```

1. query balance of existing accounts at block 3 from node1:
```go run script/view-balance/main.go -id 1 -bn 3```
