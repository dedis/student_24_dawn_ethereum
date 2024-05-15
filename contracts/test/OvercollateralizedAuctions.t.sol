// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test, console} from "forge-std/Test.sol";
import {MockERC721} from "forge-std/mocks/MockERC721.sol";
import {MockERC20} from "forge-std/mocks/MockERC20.sol";
import {OvercollateralizedAuctions} from "../src/OvercollateralizedAuctions.sol";
import {IERC721} from "forge-std/interfaces/IERC721.sol";
import {IERC20} from "forge-std/interfaces/IERC20.sol";

contract Collection is MockERC721 {
    constructor() {
        initialize("Daredevil Iguana Squad", "DDIS");
        for (uint256 i = 0; i < 10; i++) {
            _mint(msg.sender, i);
        }
    }
}

contract OvercollateralizedAuctionsTest is Test {
    OvercollateralizedAuctions auctions;
    IERC721 collection;
    IERC20 bidToken;

    address payable constant bidder1 = payable(0x1111111111111111111111111111111111111111);
    address payable constant bidder2 = payable(0x2222222222222222222222222222222222222222);
    address payable constant proceedsReceiver = payable(0x3333333333333333333333333333333333333333);

    function setUp() public {
        auctions = new OvercollateralizedAuctions();
        collection = IERC721(address(new Collection()));
        bidToken = IERC20(address(new MockERC20()));
    }

    function startAuction(uint256 tokenId) internal returns (uint256 auctionId) {
        collection.approve(address(auctions), tokenId);
        return auctions.startAuction(collection, tokenId, bidToken, proceedsReceiver);
    }

    function prepareCommit(address bidder, uint256 amount) internal returns (bytes32 blinding, bytes32 commit) {
        blinding = hex"1234";
        commit = keccak256(abi.encodePacked(blinding, bidder, amount));
    }

    function doCommit(uint256 auctionId, address bidder, bytes32 commit) internal {
        deal(address(bidToken), bidder, 10 ether);
        vm.startPrank(bidder);
        bidToken.approve(address(auctions), 10 ether);
        auctions.commitBid(auctionId, commit);
        vm.stopPrank();
    }

    function testCommitBid() public {
        uint256 tokenId = 1;
        uint256 auctionId = startAuction(tokenId);
        uint256 amount = 1 ether;
        (bytes32 blinding, bytes32 commit) = prepareCommit(bidder1, amount);
        doCommit(auctionId, bidder1, commit);
        assertEq(bidToken.balanceOf(address(auctions)), 10 ether);
        assertEq(bidToken.balanceOf(bidder1), 0);
    }
}
