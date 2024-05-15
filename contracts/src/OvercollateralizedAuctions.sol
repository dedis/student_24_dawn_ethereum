// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {IERC721} from "forge-std/interfaces/IERC721.sol";
import {IERC20} from "forge-std/interfaces/IERC20.sol";

contract OvercollateralizedAuctions {
    struct Auction {
        IERC721 collection;
        uint256 tokenId;
        IERC20 bidToken;
        address proceedsReceiver;
        uint256 commitDeadline;
        uint256 revealDeadline;
        uint256 maxBid;
        uint256 highestAmount;
        address highestBidder;
        mapping(address => bytes32) commits;
    }

    Auction[] public auctions;

    uint256 constant delay = 60 seconds;

    function startAuction(IERC721 collection, uint256 tokenId, IERC20 bidToken, address proceedsReceiver)
        external
        returns (uint256 auctionId)
    {
        auctionId = auctions.length;
        auctions.push();
        Auction storage auction = auctions[auctionId];

        auction.collection = collection;
        auction.tokenId = tokenId;
        auction.bidToken = bidToken;
        auction.proceedsReceiver = proceedsReceiver;
        auction.commitDeadline = block.timestamp + delay;
        auction.revealDeadline = auction.commitDeadline + delay;
        auction.maxBid = 10 ether; // FIXME: hardcoded

        collection.transferFrom(msg.sender, address(this), auction.tokenId);
    }

    function computeCommitment(bytes32 blinding, address bidder, uint256 amount) public pure returns (bytes32 commit) {
        commit = keccak256(abi.encode(blinding, bidder, amount));
    }

    function commitBid(uint256 auctionId, bytes32 commit) external {
        Auction storage auction = auctions[auctionId];

        require(block.timestamp < auction.commitDeadline, "late");

        require(auction.bidToken.transferFrom(msg.sender, address(this), auction.maxBid));

        // NOTE: bidders can self-grief by overwriting their commit
        auction.commits[msg.sender] = commit;
    }

    function revealBid(uint256 auctionId, bytes32 blinding, uint256 amount) external {
        Auction storage auction = auctions[auctionId];

        require(block.timestamp >= auction.commitDeadline, "early");
        require(block.timestamp < auction.revealDeadline, "late");

        bytes32 commit = computeCommitment(blinding, msg.sender, amount);
        require(auction.commits[msg.sender] == commit, "commit");
        auction.commits[msg.sender] = "";

        if (amount > auction.highestAmount) {
            address prevHighestBidder = auction.highestBidder;
            uint256 prevHighestAmount = auction.highestAmount;
            auction.highestBidder = msg.sender;
            auction.highestAmount = amount;
            auction.bidToken.transfer(msg.sender, auction.maxBid - amount);
            auction.bidToken.transfer(prevHighestBidder, prevHighestAmount);
        } else {
            auction.bidToken.transfer(msg.sender, auction.maxBid);
        }
    }

    function settle(uint256 auctionId) external {
        Auction storage auction = auctions[auctionId];

        require(block.timestamp >= auction.revealDeadline, "early");
        require(address(auction.collection) != address(0));

        auction.bidToken.transfer(auction.proceedsReceiver, auction.highestAmount);
        auction.collection.transferFrom(address(this), auction.highestBidder, auction.tokenId);

        // prevent replays
        delete auction.collection;
    }
}
