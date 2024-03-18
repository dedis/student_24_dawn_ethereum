// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Script, console} from "forge-std/Script.sol";
import {Auction} from "src/Auction.sol";

contract CloseAuction is Script {
    Auction private auction = Auction(0xef434c1405f66997CBf4a04FDDed518C28a6a013);
    

    function run() public {
        vm.startBroadcast();
        console.log("Highest bidder:", auction.highestBidder());
        console.log("bid:", auction.highestBidAmount());
	auction.claim();
    }
}
