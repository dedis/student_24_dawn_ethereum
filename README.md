# Geth-F3B
Geth with F3B front-running protection achieved in consensus layer

based on go-ethereum v1.10.23

# Usefule command

```
// install node on ubuntu server
curl -sL https://deb.nodesource.com/setup_16.x | sudo -E bash -
sudo apt-get install gcc g++ make
sudo apt-get install nodejs
npm i web3 log-timestamp eth
git clone https://github.com/ethereum/go-ethereum.git -b release/1.10

// build
go install -v ./...
go install -v ./cmd/geth

// install golang
sudo apt update
sudo apt upgrade
sudo apt install golang-go

// set up the chain 
https://arctouch.com/blog/how-to-set-up-ethereum-blockchain

geth --datadir .ethereum/ account new

puppeth

// init
geth --datadir .ethereum/ init clique.json

geth --datadir .ethereum1/ init ../multiclique.json
geth --datadir .ethereum2/ init ../multiclique.json

// start the chain
geth --nodiscover --networkid 42 --datadir .ethereum/ --unlock 0x280F6B48E4d9aEe0Efdb04EeBe882023357f6434 --mine

// multiclique
geth --nodiscover --networkid 43 --datadir .ethereum1/ --unlock 0x280F6B48E4d9aEe0Efdb04EeBe882023357f6434 --mine --ipcpath \\.\pipe\geth1.ipc --authrpc.port 8551 --port 30303

geth --nodiscover --networkid 43 --datadir .ethereum2/ --unlock 0xa9ca84343c8dB08d596400d35A7034027A5F4b31 --mine --ipcpath \\.\pipe\geth2.ipc --authrpc.port 8552 --port 30304 --miner.etherbase 0xa9ca84343c8dB08d596400d35A7034027A5F4b31 --syncmode full

// addPeers
geth attach \\.\pipe\geth1.ipc

admin.nodeInfo.enode

geth attach \\.\pipe\geth2.ipc

admin.addPeer

admin.peers


TODO: check addpeer
// interact with blockchain
geth attach //./pipe/geth.ipc

// generate vanity address / contract address
npm install -g vanity-eth

vanity -i f413 --contract

vanity -n 10 -i DEADbeef -c
```

Execution time:
    avg     err
1

10

100