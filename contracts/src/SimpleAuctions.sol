// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {IERC721} from "forge-std/interfaces/IERC721.sol";
import {IERC20} from "forge-std/interfaces/IERC20.sol";

contract SimpleAuctions {
    struct Auction {
        IERC721 collection;
        uint256 tokenId;
        IERC20 bidToken;
        address proceedsReceiver;
        uint64 deadline;
        uint256 highestAmount;
        address highestBidder;
    }

    event AuctionStarted(
        uint256 auctionId,
        IERC721 collection,
        uint256 tokenId,
        IERC20 bidToken,
        address proceedsReceiver,
        uint64 deadline
    );

    Auction[] public auctions;

    uint64 constant delay = 60 seconds;

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
        auction.deadline = uint64(block.timestamp) + delay;

        collection.transferFrom(msg.sender, address(this), auction.tokenId);

        emit AuctionStarted(auctionId, collection, tokenId, bidToken, proceedsReceiver, auction.deadline);
    }

    function bid(uint256 auctionId, uint256 amount) external {
        Auction storage auction = auctions[auctionId];

        require(block.timestamp < auction.deadline, "late");

        if (amount > auction.highestAmount) {
            address prevHighestBidder = auction.highestBidder;
            uint256 prevHighestAmount = auction.highestAmount;
            auction.highestBidder = msg.sender;
            auction.highestAmount = amount;
            auction.bidToken.transferFrom(msg.sender, address(this), amount);
            auction.bidToken.transfer(prevHighestBidder, prevHighestAmount);
        }
    }

    function settle(uint256 auctionId) external {
        Auction storage auction = auctions[auctionId];

        require(block.timestamp >= auction.deadline, "early");
        require(address(auction.collection) != address(0));

        auction.bidToken.transfer(auction.proceedsReceiver, auction.highestAmount);
        auction.collection.transferFrom(address(this), auction.highestBidder, auction.tokenId);

        // prevent replays
        delete auction.collection;
    }
}
