// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {OvercollateralizedAuctions} from "src/OvercollateralizedAuctions.sol";
import {WETH} from "solady/tokens/WETH.sol";
import {Collection} from "src/Collection.sol";

import {Script} from "forge-std/Script.sol";

contract Setup is Script {
    function run() public {
        vm.startBroadcast();
        WETH weth = new WETH();
        OvercollateralizedAuctions auctions = new OvercollateralizedAuctions();
        Collection collection = new Collection();
        vm.stopBroadcast();

        string memory obj = "arbitrary";
        vm.serializeAddress(obj, "auctions", address(auctions));
        vm.serializeAddress(obj, "weth", address(weth));
        string memory addresses = vm.serializeAddress(obj, "collection", address(collection));
        vm.writeJson(addresses, vm.envString("ADDRESSES_FILE"));
    }
}
