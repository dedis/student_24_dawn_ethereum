// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Script, console} from "forge-std/Script.sol";
import {Auction} from "src/Auction.sol";

contract CloseAuction is Script {
    Auction private auction = Auction(0x3712327B0E9fAE301cFED65eD6BDEf03629fCCFa);

    function run() public {
        vm.startBroadcast();
        console.log("Highest bidder:", auction.highestBidder());
        console.log("bid:", auction.highestBidAmount());
	auction.claim();
    }
}
