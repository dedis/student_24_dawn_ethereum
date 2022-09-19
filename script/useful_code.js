// connect geth console
geth attach .ethereum/geth.ipc

// run a geth node with id=42
geth --nodiscover --networkid 42 --datadir .ethereum/ --unlock 0x6e531c7d9cbc97aaf3fa53839ed5049e458910be --mine

// unlock an account encrypted in geth keystore
web3.personal.unlockAccount("0x7c20c88bc915aa3aa59cb9721c1e84d212d8cac7", "", 1000);

// send eth to another
eth.getBalance("0x7c20c88bc915aa3aa59cb9721c1e84d212d8cac7")
// 9.04625697166532776746648320380374280103671755200316906558262375061821325312e+74

var amount_to_send_eth = web3.fromWei(eth.getBalance("0x7c20c88bc915aa3aa59cb9721c1e84d212d8cac7"), "ether")

var amount_to_send_wei = amount_to_send_eth *1000000000000000000

var transactionFee = web3.eth.gasPrice * 21001;

var total_amount_to_send_wei = transactionFee + amount_to_send_wei / 10

eth.sendTransaction({from:"0x7c20c88bc915aa3aa59cb9721c1e84d212d8cac7", to:"0x6e531c7d9cbc97aaf3fa53839ed5049e458910be", value: total_amount_to_send_wei});

// get a block with id
web3.eth.getBlock(162, function(e, r) { blockInfo = r; });

blockInfo

// last block with tx: 162

// install node on ubuntu
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

// start the chain
geth --nodiscover --networkid 42 --datadir .ethereum/ --unlock 0x280F6B48E4d9aEe0Efdb04EeBe882023357f6434 --mine

// interact with blockchain
geth attach //./pipe/geth.ipc