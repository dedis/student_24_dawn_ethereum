// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {ERC721} from "solady/tokens/ERC721.sol";
import {SafeTransferLib} from "solady/utils/SafeTransferLib.sol";

using SafeTransferLib for address payable;

contract Auction is ERC721 {
    uint256 public highestBidAmount;
    address public highestBidder;
    uint256 public deadline;
    uint256 public constant duration = 4 /* blocks */;
    address payable public immutable proceedsReceiver;

    function name() public pure override returns (string memory) {
        return "Daredevil Iguana Squad";
    }

    function symbol() public pure override returns (string memory) {
        return "DDIS";
    }

    function tokenURI(uint256) public pure override returns (string memory) {
        return "https://dedis.ch";
    }

    constructor() {
        proceedsReceiver = payable(msg.sender);
        deadline = 0;
    }

    function start() external {
	require(deadline <= block.number, "Auction ongoing");
	deadline = block.number + duration;
    }

    function bid() external payable {
        require(block.number < deadline, "Auction has ended");
        require(msg.value > highestBidAmount, "Bid too low");

        address prevHighestBidder = highestBidder;
        uint256 prevHighestBidAmount = highestBidAmount;
        highestBidAmount = msg.value;
        highestBidder = msg.sender;

        if (prevHighestBidAmount > 0) {
            payable(prevHighestBidder).forceSafeTransferETH(prevHighestBidAmount);
        }
    }

    function claim() external {
        require(block.number >= deadline, "Auction ongoing");
        require(highestBidder != address(0), "No bids received");
        _mint(highestBidder, 0);
        proceedsReceiver.safeTransferETH(highestBidAmount);

        // reset auction
        highestBidAmount = 0;
        highestBidder = address(0);
    }
}
