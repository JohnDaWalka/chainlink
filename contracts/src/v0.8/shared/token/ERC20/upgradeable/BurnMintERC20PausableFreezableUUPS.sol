// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {BurnMintERC20PausableUUPS} from "./BurnMintERC20PausableUUPS.sol";

contract BurnMintERC20PausableFreezableUUPS is BurnMintERC20PausableUUPS {
  event AccountFrozen(address indexed account);
  event AccountUnfrozen(address indexed account);

  error BurnMintERC20PausableFreezableUUPS__Paused();
  error BurnMintERC20PausableFreezableUUPS__InvalidRecipient(address recipient);
  error BurnMintERC20PausableFreezableUUPS__AccountFrozen(address account);

  bytes32 public constant FREEZER_ROLE = keccak256("FREEZER_ROLE");

  // ================================================================
  // │                          Storage                             │
  // ================================================================

  /// @dev Mapping to keep track of the freezed status of an address
  mapping(address => bool) internal s_isFrozen;

  /**
   * @dev This empty reserved space is put in place to allow future versions to add new
   * variables without shifting down storage in the inheritance chain.
   * See https://docs.openzeppelin.com/contracts/4.x/upgradeable#storage_gaps
   */
  uint256[49] private __gap;

  // ================================================================
  // │                            UUPS                              │
  // ================================================================

  /// @custom:oz-upgrades-unsafe-allow constructor
  constructor() {
    _disableInitializers();
  }

  function initialize(
    string memory name,
    string memory symbol,
    uint8 decimals,
    uint256 maxSupply,
    uint256 preMint,
    address defaultAdmin,
    address defaultUpgrader,
    address defaultPauser,
    address defaultFreezer
  ) public initializer {
    super.initialize(name, symbol, decimals, maxSupply, preMint, defaultAdmin, defaultUpgrader, defaultPauser);

    _grantRole(FREEZER_ROLE, defaultFreezer);
  }

  // ================================================================
  // │                         Freezing                             │
  // ================================================================

  /// @notice Freezes an account, disallowing transfers, minting and burning from/to it.
  /// @dev Requires the caller to have the FREEZER_ROLE.
  /// @dev Can be called even if the contract is paused.
  function freeze(
    address account
  ) public onlyRole(FREEZER_ROLE) {
    if (account == address(this)) revert BurnMintERC20PausableFreezableUUPS__InvalidRecipient(account);

    s_isFrozen[account] = true;

    emit AccountFrozen(account);
  }

  /// @notice Unfreezes an account
  /// @dev Requires the caller to have the FREEZER_ROLE.
  /// @dev Can be called even if the contract is paused.
  function unfreeze(
    address account
  ) public onlyRole(FREEZER_ROLE) {
    s_isFrozen[account] = false;

    emit AccountUnfrozen(account);
  }

  function isFrozen(
    address account
  ) public view returns (bool) {
    return s_isFrozen[account];
  }

  // ================================================================
  // │                            ERC20                             │
  // ================================================================

  /// @dev Uses OZ BurnMintERC20PausableUUPS _beforeTokenTransfer hook to disallow transfers, minting and burning from/to frozen addresses.
  function _beforeTokenTransfer(address from, address to, uint256 amount) internal virtual override {
    super._beforeTokenTransfer(from, to, amount);

    if (s_isFrozen[from]) revert BurnMintERC20PausableFreezableUUPS__AccountFrozen(from);
    if (s_isFrozen[to]) revert BurnMintERC20PausableFreezableUUPS__AccountFrozen(to);
  }

  /// @dev Uses OZ BurnMintERC20PausableUUPS _approve to disallow approving from and to frozen addresses.
  function _approve(address owner, address spender, uint256 amount) internal virtual override {
    if (s_isFrozen[owner]) revert BurnMintERC20PausableFreezableUUPS__AccountFrozen(owner);
    if (s_isFrozen[spender]) revert BurnMintERC20PausableFreezableUUPS__AccountFrozen(spender);

    super._approve(owner, spender, amount);
  }
}
