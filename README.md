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
- [x] The `to` address is encrypted
- [ ] The transaction receipt contains a symmetric encryption key
- [ ] The execution layer can direct the SMC node to release an encryption label only after the transaction is finalized.
- [ ] TDH2, PVSS, beacon IBE options maybe
- [ ] Direct contract creation doesn't work ? (nil address thing)

# Running
Run `setup.sh` to build and install to `$GOPATH` the modified `dela and `go-ethereum`.

Run `git submodule update --init --recursive` to make sure the foundry dependencies are ready.

Run `run.sh` to start all the services and give a shell.

# Config
There is a default parameter that is changeable in the system in `go-ethereum/core/types/transaction.go` as shown below. The EncryptedBlockDelay defines the block delay between the ordering block and the execution block.
```
const EncryptedBlockDelay uint64 = 2
```
