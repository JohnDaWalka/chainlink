// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

/*
 * @title Common
 * @author Michael Fletcher
 * @notice Common functions and structs
 */
library CommonV5 {
  struct Config {
    // The ID of the Config
    bytes32 configDigest;
    // Fault tolerance of the DON
    uint8 f;
    // Whether the config is active
    bool isActive;
    // Map of signer addresses to configDigest
    mapping(address => bool) oracles;
  }
}
