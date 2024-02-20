# DELA F3B-IBE

Software repository delivered as part of the "Optimizing Frontrunning Protection" research project at @dedis.

The main point of interest is [dkg/pedersen_bn256](./dkg/pedersen_bn256).

## Demo

A simple demo simulating a frontrunning protected decentralized exchange is available.
Simply run `docker compose run demo` from the top-level of the repository.
Alternatively, `./demo.sh` can be run directly inside [tmux],
with the dependencies being Vim 9.0 (for `xxd`) and Go.

[tmux]: https://tmux.github.io

## Code Coverage

```
$ go test  -covermode=count -coverprofile=profile.cov ./dkg/pedersen_bn256/
ok  	go.dedis.ch/dela/dkg/pedersen_bn256	93.347s	coverage: 96.7% of statements
```
