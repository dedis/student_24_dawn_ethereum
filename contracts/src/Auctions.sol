// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {IERC721} from "forge-std/interfaces/IERC721.sol";
import {IERC20} from "forge-std/interfaces/IERC20.sol";

// common interface for Simple- and OvercollateralizedAuctions
interface Auctions {
    event AuctionStarted(
        uint256 auctionId,
        IERC721 collection,
        uint256 tokenId,
        IERC20 bidToken,
        address proceedsReceiver,
        uint64 commitDeadline,
        // In case of SimpleAuction, revealDeadline = commitDeadline
        uint64 revealDeadline,
        // In case of SimpleAuction, maxBid = type(uint256).max
        uint256 maxBid
    );

    function startAuction(IERC721 collection, uint256 tokenId, IERC20 bidToken, address proceedsReceiver)
        external
        returns (uint256 auctionId);

    function settle(uint256 auctionId) external;
}
