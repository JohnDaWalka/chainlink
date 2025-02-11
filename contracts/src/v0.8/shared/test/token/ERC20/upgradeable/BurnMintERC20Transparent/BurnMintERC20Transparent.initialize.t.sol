// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BurnMintERC20Transparent} from "../../../../../token/ERC20/upgradeable/BurnMintERC20Transparent.sol";
import {Initializable} from "../../../../../token/ERC20/upgradeable/BurnMintERC20Transparent.sol";
import {BurnMintERC20TransparentSetup} from "./BurnMintERC20TransparentSetup.t.sol";

contract BurnMintERC20Transparent_initialize is BurnMintERC20TransparentSetup {
  function test_Initialize() public view {
    assertEq(s_burnMintERC20Transparent.name(), s_name);
    assertEq(s_burnMintERC20Transparent.symbol(), s_symbol);
    assertEq(s_burnMintERC20Transparent.decimals(), s_decimals);
    assertEq(s_burnMintERC20Transparent.maxSupply(), s_maxSupply);
    assertEq(s_burnMintERC20Transparent.totalSupply(), s_preMint);

    assertTrue(s_burnMintERC20Transparent.hasRole(s_burnMintERC20Transparent.DEFAULT_ADMIN_ROLE(), s_defaultAdmin));
  }

  function test_Initialize_RevertWhen_AlreadyInitialized() public {
    vm.expectRevert(abi.encodeWithSelector(Initializable.InvalidInitialization.selector));

    s_burnMintERC20Transparent.initialize(s_name, s_symbol, s_decimals, s_maxSupply, s_preMint, s_defaultAdmin);
  }

  /// @dev Adding _disableInitializers() function to implementation's constructor ensures that no one can call initialize directly on the implementation.
  /// @dev The initialize should be only callable through Proxy.
  /// @dev This test tests that case.
  function test_Initialize_RevertWhen_CallIsNotThroughProxy() public {
    BurnMintERC20Transparent newBurnMintERC20Transparent = new BurnMintERC20Transparent();

    vm.expectRevert(abi.encodeWithSelector(Initializable.InvalidInitialization.selector));

    newBurnMintERC20Transparent.initialize(s_name, s_symbol, s_decimals, s_maxSupply, s_preMint, s_defaultAdmin);
  }
}
