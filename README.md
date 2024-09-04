# Dawn
Delayed execution Ethereum, part of ["Optimizing Commit-and-reveal for Smart Contracts", Master's Thesis, Julie Bettens.](https://blog.bbjubjub.fr/thesis.pdf)
Based on previous work by Shufan Wang.

The goal of the project is to enable decentralized applications that require pre-execution privacy.
Execution is delayed until the order of transactions is finalized by the blockchain.
Users can encrypt the payload of their transactions in such a way that they can be decrypted before execution.
See the report for more details.

NOTE: this is academic code, it is trivially exploitable if put in production!

# Architecture
`smc/` contains the Secret Management Committee code and provides the `smccli` command.
This is based on previous work by Mahsa Bastankhah.

`go-ethereum/` is a modified version of `go-ethereum v1.10.23` which integrates with the SMC.
[`git-subrepo`](https://github.com/ingydotnet/git-subrepo) can be used to track differences with upstream.

# Running
Make sure [Foundry](https://getfoundry.sh/) is installed.
Run `git submodule update --init --recursive contracts` to make sure the foundry dependencies are ready.

Run `setup.sh` to build and install `smccli` and `geth` to `$GOPATH`.

Run `run.sh` to start all the services and give a shell.

Run `gen.sh` to update generated files.

Run `measure.sh` to perform a throughput benchmark.
Other benchmarks are in `go-ethereum/f3b` and can be run with `go test -bench`.

# Parameters
The `params.toml` file can be edited to change the experiment parameters.
In particular, it allows selecting the encryption method or reverting to hash-based commit-and-reveal.

# Copying
Original work by DEDIS is dual-licensed MIT and Apache 2.
Geth and Kyber retain their original licenses.

# Known bugs

`cast` seems to have issues with transaction signatures during the experiments.
This applies even when the Lausanne hard fork is off.
Try for example `cast tx 0xedf4e0d37e54c49710960396abff350f90c6b6aa83100ebba236e488e39dd586`
and note that the `from` value is the zero address.
The cause of this is currently unknown.

Context binding is not implemented in the TIBE case.
The simplest solution would be to add the transaction label to the HKDF step.

The Geth miner will normally prepare an empty block before starting block building.
We disabled this due to the shadow block rule.

EIP-1559 is forcefully disabled.

It may happen that Clique decrees a block is *lost* if decryption causes too much delay.

More of a quirk: the SMC uses Schnorr over BN254 as a way to authenticate nodes.
It is a bit odd, it would make more sense to do Ed25519.
