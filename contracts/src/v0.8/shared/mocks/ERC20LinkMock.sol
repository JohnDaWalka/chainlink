// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {LinkTokenInterface} from "../interfaces/LinkTokenInterface.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract MockLinkToken is LinkTokenInterface, ERC20 {
    constructor() ERC20("Mock LINK", "mLINK") {}

    function allowance(address owner, address spender) public view override(ERC20, LinkTokenInterface) returns (uint256) {
        return super.allowance(owner, spender);
    }

    function approve(address spender, uint256 value) public override(ERC20, LinkTokenInterface) returns (bool) {
        return super.approve(spender, value);
    }

    function balanceOf(address owner) public view override(ERC20, LinkTokenInterface) returns (uint256) {
        return super.balanceOf(owner);
    }

    function decimals() public view override(ERC20, LinkTokenInterface) returns (uint8) {
        return super.decimals();
    }

    function decreaseApproval(address spender, uint256 subtractedValue) external override returns (bool) {
        return approve(spender, allowance(msg.sender, spender) - subtractedValue);
    }

    function increaseApproval(address spender, uint256 addedValue) external override {
        approve(spender, allowance(msg.sender, spender) + addedValue);
    }

    function name() public view override(ERC20, LinkTokenInterface) returns (string memory) {
        return super.name();
    }

    function symbol() public view override(ERC20, LinkTokenInterface) returns (string memory) {
        return super.symbol();
    }

    function totalSupply() public view override(ERC20, LinkTokenInterface) returns (uint256) {
        return super.totalSupply();
    }

    function transfer(address to, uint256 value) public override(ERC20, LinkTokenInterface) returns (bool) {
        return super.transfer(to, value);
    }

    function transferAndCall(address to, uint256 value, bytes calldata) external override returns (bool) {
        bool success = transfer(to, value);
        if (success) {
            // Add mock implementation for transferAndCall if needed
        }
        return success;
    }

    function transferFrom(address from, address to, uint256 value) public override(ERC20, LinkTokenInterface) returns (bool) {
        return super.transferFrom(from, to, value);
    }

    // Additional function to mint tokens for testing
    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}