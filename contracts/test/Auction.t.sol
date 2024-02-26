// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test, console} from "forge-std/Test.sol";
import {Auction} from "../src/Auction.sol";

contract AuctionTest is Test {
    Auction public auction;
    address payable constant bidder1 = payable(0x1111111111111111111111111111111111111111);
    address payable constant bidder2 = payable(0x2222222222222222222222222222222222222222);
    address payable constant proceedsReceiver = payable(0x3333333333333333333333333333333333333333);

    function setUp() public {
        auction = new Auction(proceedsReceiver);
    }

    function test_bid() public {
        vm.deal(bidder1, 1 wei);
        vm.prank(bidder1);
        auction.bid{value: 1 wei}();
        assertEq(auction.highestBidder(), bidder1);
        assertEq(auction.highestBidAmount(), 1 wei);
    }

    function test_claim() public {
        vm.deal(bidder1, 1 wei);
        vm.prank(bidder1);
        auction.bid{value: 1 wei}();
        vm.expectRevert();
        auction.claim();
        vm.warp(auction.deadline());
        auction.claim();
    }

    function test_failbid0() public {
        vm.prank(bidder1);
        vm.expectRevert();
        auction.bid{value: 0 wei}();
    }

    function test_failbid1() public {
        vm.deal(bidder1, 1 wei);
        vm.prank(bidder1);
        auction.bid{value: 1 wei}();
        vm.deal(bidder2, 1 wei);
        vm.prank(bidder2);
        vm.expectRevert();
        auction.bid{value: 1 wei}();
        assertEq(auction.highestBidder(), bidder1);
        assertEq(auction.highestBidAmount(), 1 wei);
    }

    function test_bid2() public {
        vm.deal(bidder1, 1 wei);
        vm.prank(bidder1);
        auction.bid{value: 1 wei}();
        vm.deal(bidder2, 2 wei);
        vm.prank(bidder2);
        auction.bid{value: 2 wei}();
        assertEq(auction.highestBidder(), bidder2);
        assertEq(auction.highestBidAmount(), 2 wei);
        assertEq(bidder1.balance, 1 wei);
    }
}
