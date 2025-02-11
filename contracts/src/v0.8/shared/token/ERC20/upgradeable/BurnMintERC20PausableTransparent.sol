// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {PausableUpgradeable} from
  "../../../../vendor/openzeppelin-solidity-upgradeable/v5.0.2/contracts/utils/PausableUpgradeable.sol";
import {BurnMintERC20Transparent} from "./BurnMintERC20Transparent.sol";

contract BurnMintERC20PausableTransparent is BurnMintERC20Transparent, PausableUpgradeable {
  error BurnMintERC20PausableTransparent__Paused();

  event Paused();
  event Unpaused();

  bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");

  // ================================================================
  // │                         Transparent                          │
  // ================================================================

  /// @custom:oz-upgrades-unsafe-allow constructor
  constructor() {
    _disableInitializers();
  }

  /// @dev the underscores in parameter names are used to suppress compiler warnings about shadowing ERC20 functions
  function initialize(
    string memory name,
    string memory symbol,
    uint8 decimals_,
    uint256 maxSupply_,
    uint256 preMint,
    address defaultAdmin,
    address defaultPauser
  ) public initializer {
    super.initialize(name, symbol, decimals_, maxSupply_, preMint, defaultAdmin);

    _grantRole(PAUSER_ROLE, defaultPauser);
  }

  // ================================================================
  // │                          Pausing                             │
  // ================================================================

  /// @notice Pauses the implementation.
  /// @dev Requires the caller to have the PAUSER_ROLE.
  function pause() public onlyRole(PAUSER_ROLE) {
    _pause();

    emit Paused();
  }

  /// @notice Unpauses the implementation.
  /// @dev Requires the caller to have the PAUSER_ROLE.
  function unpause() public onlyRole(PAUSER_ROLE) {
    _unpause();

    emit Unpaused();
  }

  // ================================================================
  // │                            ERC20                             │
  // ================================================================

  /// @dev Disallows sending, minting and burning if implementation is paused.
  function _update(address from, address to, uint256 value) internal virtual override {
    if (paused()) revert BurnMintERC20PausableTransparent__Paused();

    super._update(from, to, value);
  }

  /// @dev Disallows approving if implementation is paused.
  function _approve(address owner, address spender, uint256 value, bool emitEvent) internal virtual override {
    if (paused()) revert BurnMintERC20PausableTransparent__Paused();

    super._approve(owner, spender, value, emitEvent);
  }
}
