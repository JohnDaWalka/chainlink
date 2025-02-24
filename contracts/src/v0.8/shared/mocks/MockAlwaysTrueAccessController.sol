// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {AccessControllerInterface} from "../interfaces/AccessControllerInterface.sol";

contract MockAlwaysTrueAccessController is AccessControllerInterface {
    /// @notice Always returns true for access checks
    /// @dev Ignores the input parameters and always returns true
    /// @return True, indicating access is always granted
    function hasAccess(address, bytes calldata) external pure override returns (bool) {
        return true;
    }
}