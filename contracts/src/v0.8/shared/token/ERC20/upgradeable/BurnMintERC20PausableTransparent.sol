// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {IGetCCIPAdmin} from "../../../../ccip/interfaces/IGetCCIPAdmin.sol";
import {IBurnMintERC20Upgradeable} from "../../../../shared/token/ERC20/upgradeable/IBurnMintERC20Upgradeable.sol";

import {Initializable} from
  "../../../../vendor/openzeppelin-solidity-upgradeable/v5.0.2/contracts/proxy/utils/Initializable.sol";

import {AccessControlUpgradeable} from
  "../../../../vendor/openzeppelin-solidity-upgradeable/v5.0.2/contracts/access/AccessControlUpgradeable.sol";
import {ERC20BurnableUpgradeable} from
  "../../../../vendor/openzeppelin-solidity-upgradeable/v5.0.2/contracts/token/ERC20/extensions/ERC20BurnableUpgradeable.sol";
import {PausableUpgradeable} from
  "../../../../vendor/openzeppelin-solidity-upgradeable/v5.0.2/contracts/utils/PausableUpgradeable.sol";
import {IAccessControl} from "../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/access/IAccessControl.sol";
import {IERC20} from "../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/interfaces/IERC20.sol";
import {IERC165} from "../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/introspection/IERC165.sol";

contract BurnMintERC20PausableTransparent is
  Initializable,
  IBurnMintERC20Upgradeable,
  IGetCCIPAdmin,
  IERC165,
  ERC20BurnableUpgradeable,
  AccessControlUpgradeable,
  PausableUpgradeable
{
  error BurnMintERC20PausableTransparent__MaxSupplyExceeded(uint256 supplyAfterMint);
  error BurnMintERC20PausableTransparent__InvalidRecipient(address recipient);
  error BurnMintERC20PausableTransparent__Paused();

  event CCIPAdminTransferred(address indexed previousAdmin, address indexed newAdmin);
  event Paused();
  event Unpaused();

  bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
  bytes32 public constant BURNER_ROLE = keccak256("BURNER_ROLE");
  bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");

  // ================================================================
  // │                         Storage                              │
  // ================================================================

  /// @custom:storage-location erc7201:chainlink.storage.BurnMintERC20PausableTransparent
  struct BurnMintERC20PausableTransparentStorage {
    /// @dev the CCIPAdmin can be used to register with the CCIP token admin registry, but has no other special powers, and can only be transferred by the owner.
    address s_ccipAdmin;
    /// @dev The number of decimals for the token
    uint8 s_decimals;
    /// @dev The maximum supply of the token, 0 if unlimited
    uint256 s_maxSupply;
  }

  // keccak256(abi.encode(uint256(keccak256("chainlink.storage.BurnMintERC20PausableTransparent")) - 1)) & ~bytes32(uint256(0xff));
  bytes32 private constant BURN_MINT_ERC20_PAUSABLE_TRANSPARENT_STORAGE_LOCATION =
    0x59b4c9c2cce0d798b79dd0fc5ebe1928d8919a5c6a224a033d19ec4c09b58b00;

  // solhint-disable-next-line chainlink-solidity/explicit-returns
  function _getBurnMintERC20PausableTransparentStorage()
    private
    pure
    returns (BurnMintERC20PausableTransparentStorage storage $)
  {
    assembly {
      $.slot := BURN_MINT_ERC20_PAUSABLE_TRANSPARENT_STORAGE_LOCATION
    }
  }

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
    __ERC20_init(name, symbol);
    __ERC20Burnable_init();
    __AccessControl_init();
    __Pausable_init();

    BurnMintERC20PausableTransparentStorage storage $ = _getBurnMintERC20PausableTransparentStorage();

    $.s_decimals = decimals_;
    $.s_maxSupply = maxSupply_;

    $.s_ccipAdmin = defaultAdmin;

    if (preMint != 0) {
      _mint(defaultAdmin, preMint);
    }

    _grantRole(DEFAULT_ADMIN_ROLE, defaultAdmin);
    _grantRole(PAUSER_ROLE, defaultPauser);
  }

  // ================================================================
  // │                           ERC165                             │
  // ================================================================

  /// @inheritdoc IERC165
  function supportsInterface(
    bytes4 interfaceId
  ) public pure virtual override(AccessControlUpgradeable, IERC165) returns (bool) {
    return interfaceId == type(IERC20).interfaceId || interfaceId == type(IBurnMintERC20Upgradeable).interfaceId
      || interfaceId == type(IERC165).interfaceId || interfaceId == type(IAccessControl).interfaceId
      || interfaceId == type(IGetCCIPAdmin).interfaceId;
  }

  // ================================================================
  // │                            ERC20                             │
  // ================================================================

  /// @dev Returns the number of decimals used in its user representation.
  function decimals() public view virtual override returns (uint8) {
    BurnMintERC20PausableTransparentStorage storage $ = _getBurnMintERC20PausableTransparentStorage();
    return $.s_decimals;
  }

  /// @dev Returns the max supply of the token, 0 if unlimited.
  function maxSupply() public view virtual returns (uint256) {
    BurnMintERC20PausableTransparentStorage storage $ = _getBurnMintERC20PausableTransparentStorage();
    return $.s_maxSupply;
  }

  /// @dev Uses OZ ERC20Upgradeable _update hook to disallow transfers, minting and burning if implementation is paused.
  /// @dev Disallows sending, minting and burning if implementation is paused.
  /// @dev Disallows sending to address(this)
  function _update(address from, address to, uint256 value) internal virtual override {
    if (paused()) revert BurnMintERC20PausableTransparent__Paused();
    if (to == address(this)) revert BurnMintERC20PausableTransparent__InvalidRecipient(to);

    super._update(from, to, value);
  }

  /// @dev Uses OZ ERC20Upgradeable _approve to disallow approving for address(0).
  /// @dev Disallows approving if implementation is paused.
  /// @dev Disallows approving for address(this)
  function _approve(address owner, address spender, uint256 value, bool emitEvent) internal virtual override {
    if (paused()) revert BurnMintERC20PausableTransparent__Paused();
    if (spender == address(this)) revert BurnMintERC20PausableTransparent__InvalidRecipient(spender);

    super._approve(owner, spender, value, emitEvent);
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
    BurnMintERC20PausableTransparentStorage storage $ = _getBurnMintERC20PausableTransparentStorage();
    uint256 _maxSupply = $.s_maxSupply;
    uint256 _totalSupply = totalSupply();

    if (_maxSupply != 0 && _totalSupply + amount > _maxSupply) {
      revert BurnMintERC20PausableTransparent__MaxSupplyExceeded(_totalSupply + amount);
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

    emit Paused();
  }

  /// @notice Unpauses the implementation.
  /// @dev Requires the caller to have the PAUSER_ROLE.
  function unpause() public onlyRole(PAUSER_ROLE) {
    _unpause();

    emit Unpaused();
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
    BurnMintERC20PausableTransparentStorage storage $ = _getBurnMintERC20PausableTransparentStorage();
    return $.s_ccipAdmin;
  }

  /// @notice Transfers the CCIPAdmin role to a new address
  /// @dev only the owner can call this function, NOT the current ccipAdmin, and 1-step ownership transfer is used.
  /// @param newAdmin The address to transfer the CCIPAdmin role to. Setting to address(0) is a valid way to revoke
  /// the role
  function setCCIPAdmin(
    address newAdmin
  ) external onlyRole(DEFAULT_ADMIN_ROLE) {
    BurnMintERC20PausableTransparentStorage storage $ = _getBurnMintERC20PausableTransparentStorage();
    address currentAdmin = $.s_ccipAdmin;

    $.s_ccipAdmin = newAdmin;

    emit CCIPAdminTransferred(currentAdmin, newAdmin);
  }
}
