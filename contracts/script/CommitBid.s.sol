// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {OvercollateralizedAuctions} from "src/OvercollateralizedAuctions.sol";

import {Script} from "forge-std/Script.sol";

contract CommitBid is Script {
    function run(OvercollateralizedAuctions auctions, uint256 auctionId, bytes32 blinding, uint256 amount) public {
        bytes32 commit = auctions.computeCommitment(blinding, msg.sender, amount);
        vm.broadcast();
        auctions.commitBid(auctionId, commit);
    }
}
