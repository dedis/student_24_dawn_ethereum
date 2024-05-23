// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {IERC721} from "forge-std/interfaces/IERC721.sol";
import {IERC20} from "forge-std/interfaces/IERC20.sol";

// common interface for Simple- and OvercollateralizedAuctions
interface Auctions {
    event AuctionStarted( // block after which commits are accepted
        // last block where commits can bet included
        // last block where reveals can be included
        uint256 auctionId,
        IERC721 collection,
        uint256 tokenId,
        IERC20 bidToken,
        address proceedsReceiver,
        uint64 opening,
        uint64 commitDeadline,
        uint64 revealDeadline,
        // In case of SimpleAuction, revealDeadline = commitDeadline
        uint256 maxBid
    );
    // In case of SimpleAuction, maxBid = type(uint256).max

    function startAuction(IERC721 collection, uint256 tokenId, IERC20 bidToken, address proceedsReceiver)
        external
        returns (uint256 auctionId);

    function settle(uint256 auctionId) external;
}
