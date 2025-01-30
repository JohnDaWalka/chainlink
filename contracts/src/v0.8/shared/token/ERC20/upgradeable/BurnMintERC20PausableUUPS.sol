// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {IGetCCIPAdmin} from "../../../../ccip/interfaces/IGetCCIPAdmin.sol";
import {IBurnMintERC20Upgradeable} from "../../../../shared/token/ERC20/upgradeable/IBurnMintERC20Upgradeable.sol";

import {Initializable} from
  "../../../../vendor/openzeppelin-solidity-upgradeable/v4.8.3/contracts/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from
  "../../../../vendor/openzeppelin-solidity-upgradeable/v4.8.3/contracts/proxy/utils/UUPSUpgradeable.sol";

import {AccessControlUpgradeable} from
  "../../../../vendor/openzeppelin-solidity-upgradeable/v4.8.3/contracts/access/AccessControlUpgradeable.sol";
import {IAccessControlUpgradeable} from
  "../../../../vendor/openzeppelin-solidity-upgradeable/v4.8.3/contracts/access/IAccessControlUpgradeable.sol";
import {PausableUpgradeable} from
  "../../../../vendor/openzeppelin-solidity-upgradeable/v4.8.3/contracts/security/PausableUpgradeable.sol";
import {ERC20Upgradeable} from
  "../../../../vendor/openzeppelin-solidity-upgradeable/v4.8.3/contracts/token/ERC20/ERC20Upgradeable.sol";
import {IERC20Upgradeable} from
  "../../../../vendor/openzeppelin-solidity-upgradeable/v4.8.3/contracts/token/ERC20/IERC20Upgradeable.sol";
import {ERC20BurnableUpgradeable} from
  "../../../../vendor/openzeppelin-solidity-upgradeable/v4.8.3/contracts/token/ERC20/extensions/ERC20BurnableUpgradeable.sol";
import {IERC165Upgradeable} from
  "../../../../vendor/openzeppelin-solidity-upgradeable/v4.8.3/contracts/utils/introspection/IERC165Upgradeable.sol";

contract BurnMintERC20PausableUUPS is
  Initializable,
  UUPSUpgradeable,
  IBurnMintERC20Upgradeable,
  IGetCCIPAdmin,
  IERC165Upgradeable,
  ERC20BurnableUpgradeable,
  AccessControlUpgradeable,
  PausableUpgradeable
{
  error BurnMintERC20PausableUUPS__MaxSupplyExceeded(uint256 supplyAfterMint);
  error BurnMintERC20PausableUUPS__InvalidRecipient(address recipient);
  error BurnMintERC20PausableUUPS__Paused();

  event CCIPAdminTransferred(address indexed previousAdmin, address indexed newAdmin);

  bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
  bytes32 public constant BURNER_ROLE = keccak256("BURNER_ROLE");
  bytes32 public constant UPGRADER_ROLE = keccak256("UPGRADER_ROLE");
  bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");

  // ================================================================
  // │                         Storage                              │
  // ================================================================

  /// @dev the CCIPAdmin can be used to register with the CCIP token admin registry, but has no other special powers,
  /// and can only be transferred by the owner.
  address internal s_ccipAdmin;

  /// @dev The number of decimals for the token
  uint8 internal s_decimals;

  /// @dev The maximum supply of the token, 0 if unlimited
  uint256 internal s_maxSupply;

  /**
   * @dev This empty reserved space is put in place to allow future versions to add new
   * variables without shifting down storage in the inheritance chain.
   * See https://docs.openzeppelin.com/contracts/4.x/upgradeable#storage_gaps
   */
  uint256[47] private __gap;

  // ================================================================
  // │                            UUPS                              │
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
    address defaultUpgrader,
    address defaultPauser
  ) public initializer {
    __ERC20_init(name, symbol);
    __ERC20Burnable_init();
    __AccessControl_init();
    __Pausable_init();
    __UUPSUpgradeable_init();

    s_decimals = decimals_;
    s_maxSupply = maxSupply_;

    s_ccipAdmin = defaultAdmin;

    if (preMint != 0) {
      _mint(defaultAdmin, preMint);
    }

    _grantRole(DEFAULT_ADMIN_ROLE, defaultAdmin);
    _grantRole(UPGRADER_ROLE, defaultUpgrader);
    _grantRole(PAUSER_ROLE, defaultPauser);
  }

  function _authorizeUpgrade(
    address newImplementation
  ) internal override onlyRole(UPGRADER_ROLE) {}

  // ================================================================
  // │                           ERC165                             │
  // ================================================================

  /// @inheritdoc IERC165Upgradeable
  function supportsInterface(
    bytes4 interfaceId
  ) public pure virtual override(AccessControlUpgradeable, IERC165Upgradeable) returns (bool) {
    return interfaceId == type(IERC20Upgradeable).interfaceId
      || interfaceId == type(IBurnMintERC20Upgradeable).interfaceId || interfaceId == type(IERC165Upgradeable).interfaceId
      || interfaceId == type(IAccessControlUpgradeable).interfaceId || interfaceId == type(IGetCCIPAdmin).interfaceId;
  }

  // ================================================================
  // │                            ERC20                             │
  // ================================================================

  /// @dev Returns the number of decimals used in its user representation.
  function decimals() public view virtual override returns (uint8) {
    return s_decimals;
  }

  /// @dev Returns the max supply of the token, 0 if unlimited.
  function maxSupply() public view virtual returns (uint256) {
    return s_maxSupply;
  }

  /// @dev Uses OZ ERC20Upgradeable _beforeTokenTransfer hook to disallow transfers, minting and burning if implementation is paused.
  /// @dev Disallows sending to address(this)
  function _beforeTokenTransfer(address from, address to, uint256 amount) internal virtual override {
    super._beforeTokenTransfer(from, to, amount);

    if (paused()) revert BurnMintERC20PausableUUPS__Paused();
    if (to == address(this)) revert BurnMintERC20PausableUUPS__InvalidRecipient(to);
  }

  /// @dev Uses OZ ERC20Upgradeable _approve to disallow approving for address(0).
  /// @dev Disallows approving if implementation is paused.
  /// @dev Disallows approving for address(this)
  function _approve(address owner, address spender, uint256 amount) internal virtual override {
    if (paused()) revert BurnMintERC20PausableUUPS__Paused();
    if (spender == address(this)) revert BurnMintERC20PausableUUPS__InvalidRecipient(spender);

    super._approve(owner, spender, amount);
  }

  // ================================================================
  // │                      Burning & minting                       │
  // ================================================================

  /// @inheritdoc ERC20BurnableUpgradeable
  /// @dev Uses OZ ERC20Upgradeable _burn to disallow burning from address(0).
  /// @dev Decreases the total supply.
  function burn(
    uint256 amount
  ) public override(IBurnMintERC20Upgradeable, ERC20BurnableUpgradeable) onlyRole(BURNER_ROLE) {
    super.burn(amount);
  }

  /// @inheritdoc IBurnMintERC20Upgradeable
  /// @dev Alias for BurnFrom for compatibility with the older naming convention.
  /// @dev Uses burnFrom for all validation & logic.
  function burn(address account, uint256 amount) public virtual override {
    burnFrom(account, amount);
  }

  /// @inheritdoc ERC20BurnableUpgradeable
  /// @dev Uses OZ ERC20Upgradeable _burn to disallow burning from address(0).
  /// @dev Decreases the total supply.
  function burnFrom(
    address account,
    uint256 amount
  ) public override(IBurnMintERC20Upgradeable, ERC20BurnableUpgradeable) onlyRole(BURNER_ROLE) {
    super.burnFrom(account, amount);
  }

  /// @inheritdoc IBurnMintERC20Upgradeable
  /// @dev Uses OZ ERC20Upgradeable _mint to disallow minting to address(0).
  /// @dev Disallows minting to address(this) via _beforeTokenTransfer hook.
  /// @dev Increases the total supply.
  function mint(address account, uint256 amount) external override onlyRole(MINTER_ROLE) {
    uint256 maxSupply_ = s_maxSupply;
    uint256 totalSupply_ = totalSupply();

    if (maxSupply_ != 0 && totalSupply_ + amount > maxSupply_) {
      revert BurnMintERC20PausableUUPS__MaxSupplyExceeded(totalSupply_ + amount);
    }

    _mint(account, amount);
  }

  // ================================================================
  // │                          Pausing                             │
  // ================================================================

  /// @notice Pauses the implementation.
  /// @dev Requires the caller to have the PAUSER_ROLE.
  function pause() public onlyRole(PAUSER_ROLE) {
    _pause();
  }

  /// @notice Unpauses the implementation.
  /// @dev Requires the caller to have the PAUSER_ROLE.
  function unpause() public onlyRole(PAUSER_ROLE) {
    _unpause();
  }

  // ================================================================
  // │                            Roles                             │
  // ================================================================

  /// @notice grants both mint and burn roles to `burnAndMinter`.
  /// @dev calls public functions so this function does not require
  /// access controls. This is handled in the inner functions.
  function grantMintAndBurnRoles(
    address burnAndMinter
  ) external {
    grantRole(MINTER_ROLE, burnAndMinter);
    grantRole(BURNER_ROLE, burnAndMinter);
  }

  /// @notice Returns the current CCIPAdmin
  function getCCIPAdmin() external view returns (address) {
    return s_ccipAdmin;
  }

  /// @notice Transfers the CCIPAdmin role to a new address
  /// @dev only the owner can call this function, NOT the current ccipAdmin, and 1-step ownership transfer is used.
  /// @param newAdmin The address to transfer the CCIPAdmin role to. Setting to address(0) is a valid way to revoke
  /// the role
  function setCCIPAdmin(
    address newAdmin
  ) external onlyRole(DEFAULT_ADMIN_ROLE) {
    address currentAdmin = s_ccipAdmin;

    s_ccipAdmin = newAdmin;

    emit CCIPAdminTransferred(currentAdmin, newAdmin);
  }
}
