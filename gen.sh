#!/bin/sh

get_abi() (
cd contracts > /dev/null
forge inspect $1 abi
)

get_abi WETH | abigen -abi - -pkg bindings -type WETH -out bindings/weth.go
get_abi Auctions | abigen -abi - -pkg bindings -type Auctions -out bindings/auctions.go
get_abi OvercollateralizedAuctions | abigen -abi - -pkg bindings -type OvercollateralizedAuctions -out bindings/overcollateralized_auctions.go
get_abi SimpleAuctions | abigen -abi - -pkg bindings -type SimpleAuctions -out bindings/simple_auctions.go
get_abi Collection | abigen -abi - -pkg bindings -type Collection -out bindings/collection.go
