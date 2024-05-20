# F3B-Geth
Delayed-execution Ethereum client.
Based on previous work by Shufan Wang.

# Architecture
`smc/` contains the Secret Management Committee code and provides the `smccli` command.
`go-ethereum/` is a modified version of `go-ethereum v1.10.23` which integrates with the SMC.

# Design (WIP)

- [x] When creating a transaction, the user agent IBE-encrypts the calldata with dela.
- [x] The IBE label is sender address concatenated with big-endian 64-bit nonce
- [x] The ciphertext is authenticated (HMAC-SHA256, encrypt-then-MAC)
- [x] ~~The chain only accepts encrypted transactions~~ delayed execution for all transactions
- [x] The `to` address is encrypted
- [ ] ~~The transaction receipt contains a symmetric encryption key~~
- [ ] The execution layer can direct the SMC node to release an encryption label only after the transaction is finalized.
- [ ] ~~TDH2, PVSS, beacon IBE options maybe~~ sticking to pairings for compact decryption proofs
- [ ] Direct contract creation doesn't work ? (nil address thing)

# Running
Run `setup.sh` to build and install `smc` and the modified `go-ethereum` to `$GOPATH`.

Make sure [Foundry](https://getfoundry.sh/) is installed.
Run `git submodule update --init --recursive contracts` to make sure the foundry dependencies are ready.

Run `run.sh` to start all the services and give a shell.

# Config
There is a default parameter that is changeable in the system in `go-ethereum/core/types/transaction.go` as shown below. The EncryptedBlockDelay defines the block delay between the ordering block and the execution block.
```
const EncryptedBlockDelay uint64 = 2
```
