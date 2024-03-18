// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Script, console} from "forge-std/Script.sol";
import {Auction} from "src/Auction.sol";

contract CloseAuction is Script {
    Auction private auction = Auction(0xF31b6eF875a924508bAB7A5922F6e34Ae2F65801);
    

    function run() public {
        vm.startBroadcast();
        console.log("Highest bidder:", auction.highestBidder());
        console.log("bid:", auction.highestBidAmount());
	auction.claim();
    }
}
