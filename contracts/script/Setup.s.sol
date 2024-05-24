// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {OvercollateralizedAuctions} from "src/OvercollateralizedAuctions.sol";
import {SimpleAuctions} from "src/SimpleAuctions.sol";
import {Auctions} from "src/Auctions.sol";
import {WETH} from "solady/tokens/WETH.sol";
import {Collection} from "src/Collection.sol";

import {Script} from "forge-std/Script.sol";

function streq(string memory a, string memory b) pure returns (bool) {
    return keccak256(abi.encodePacked(a)) == keccak256(abi.encodePacked(b));
}

contract Setup is Script {
    function run() public {
        vm.startBroadcast();
        WETH weth = new WETH();
        Auctions auctions;
	uint64 blockDelay = uint64(vm.envUint("F3B_BLOCKDELAY"));
        if (streq(vm.envString("F3B_PROTOCOL"), "")) {
            auctions = new OvercollateralizedAuctions(blockDelay);
        } else {
            auctions = new SimpleAuctions(blockDelay);
        }
        Collection collection = new Collection();
        vm.stopBroadcast();

        string memory obj = "arbitrary";
        vm.serializeAddress(obj, "auctions", address(auctions));
        vm.serializeAddress(obj, "weth", address(weth));
        string memory addresses = vm.serializeAddress(obj, "collection", address(collection));
        vm.writeJson(addresses, vm.envString("ADDRESSES_FILE"));
    }
}
