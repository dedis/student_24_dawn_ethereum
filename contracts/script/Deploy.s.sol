// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Script, console} from "forge-std/Script.sol";
import {WETH} from "solady/tokens/WETH.sol";
import {Auction} from "src/Auction.sol";

contract Deploy is Script {
    WETH private weth;
    Auction private auction;

    function run() public {
        vm.startBroadcast();
        weth = new WETH();
        auction = new Auction(payable(msg.sender));
        console.log("WETH deployed at:", address(weth));
        console.log("Auction deployed at:", address(auction));
    }
}
