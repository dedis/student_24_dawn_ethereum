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

# Parameters
The `params.toml` file can be edited to change the experiment parameters.
In particular, it allows selecting the encryption method or reverting to C&R.

# Known bugs

`cast` seems to have issues with transaction signatures during the experiments.
This applies even when the Lausanne hard fork is off.
Try for example `cast tx 0xedf4e0d37e54c49710960396abff350f90c6b6aa83100ebba236e488e39dd586`
and note that the `from` value is the zero address.
The cause of this is currently unknown.

Context binding is not implemented in the TIBE case.
The simplest solution would be to add the transaction label to a HKDF step.

The Geth miner will normally prepare an empty block before starting block building.
We disabled this due to the shadow blocks rule.

It may happen that Clique decrees a block is *lost* if decryption causes too much delay.

More of a quirk: the SMC uses Schnorr over BN254 as a way to authenticate nodes.
It is a bit odd, it would make more sense to do Ed25519.
