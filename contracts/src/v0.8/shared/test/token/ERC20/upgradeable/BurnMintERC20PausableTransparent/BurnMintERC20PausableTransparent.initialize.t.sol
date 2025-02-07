// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Upgrades} from "../../../../../../vendor/openzeppelin-foundry-upgrades/v0.3.8/Upgrades.sol";
import {
  BurnMintERC20PausableTransparent,
  Initializable
} from "../../../../../token/ERC20/upgradeable/BurnMintERC20PausableTransparent.sol";
import {BurnMintERC20PausableTransparentSetup} from "./BurnMintERC20PausableTransparentSetup.t.sol";

contract BurnMintERC20PausableTransparent_initialize is BurnMintERC20PausableTransparentSetup {
  function test_Initialize() public view {
    assertFalse(Upgrades.getAdminAddress(s_TransparentProxy) == address(0));

    assertEq(s_burnMintERC20PausableTransparent.name(), s_name);
    assertEq(s_burnMintERC20PausableTransparent.symbol(), s_symbol);
    assertEq(s_burnMintERC20PausableTransparent.decimals(), s_decimals);
    assertEq(s_burnMintERC20PausableTransparent.maxSupply(), s_maxSupply);
    assertEq(s_burnMintERC20PausableTransparent.totalSupply(), s_preMint);

    assertTrue(
      s_burnMintERC20PausableTransparent.hasRole(
        s_burnMintERC20PausableTransparent.DEFAULT_ADMIN_ROLE(), s_defaultAdmin
      )
    );
    assertTrue(
      s_burnMintERC20PausableTransparent.hasRole(s_burnMintERC20PausableTransparent.PAUSER_ROLE(), s_defaultPauser)
    );
  }

  function test_Initialize_RevertWhen_AlreadyInitialized() public {
    vm.expectRevert(abi.encodeWithSelector(Initializable.InvalidInitialization.selector));

    s_burnMintERC20PausableTransparent.initialize(
      s_name, s_symbol, s_decimals, s_maxSupply, s_preMint, s_defaultAdmin, s_defaultPauser
    );
  }

  /// @dev Adding _disableInitializers() function to implementation's constructor ensures that no one can call initialize directly on the implementation.
  /// @dev The initialize should be only callable through Proxy.
  /// @dev This test tests that case.
  function test_Initialize_RevertWhen_CallIsNotThroughProxy() public {
    BurnMintERC20PausableTransparent newBurnMintERC20PausableTransparent = new BurnMintERC20PausableTransparent();

    vm.expectRevert(abi.encodeWithSelector(Initializable.InvalidInitialization.selector));

    newBurnMintERC20PausableTransparent.initialize(
      s_name, s_symbol, s_decimals, s_maxSupply, s_preMint, s_defaultAdmin, s_defaultPauser
    );
  }
}
