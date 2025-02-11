// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BurnMintERC20UUPS, Initializable} from "../../../../../token/ERC20/upgradeable/BurnMintERC20UUPS.sol";
import {BurnMintERC20UUPSSetup} from "./BurnMintERC20UUPSSetup.t.sol";

contract BurnMintERC20UUPS_initialize is BurnMintERC20UUPSSetup {
  function test_Initialize() public view {
    assertEq(s_burnMintERC20UUPS.name(), s_name);
    assertEq(s_burnMintERC20UUPS.symbol(), s_symbol);
    assertEq(s_burnMintERC20UUPS.decimals(), s_decimals);
    assertEq(s_burnMintERC20UUPS.maxSupply(), s_maxSupply);
    assertEq(s_burnMintERC20UUPS.totalSupply(), s_preMint);

    assertTrue(s_burnMintERC20UUPS.hasRole(s_burnMintERC20UUPS.DEFAULT_ADMIN_ROLE(), s_defaultAdmin));
    assertTrue(s_burnMintERC20UUPS.hasRole(s_burnMintERC20UUPS.UPGRADER_ROLE(), s_defaultUpgrader));
  }

  function test_Initialize_RevertWhen_AlreadyInitialized() public {
    vm.expectRevert(abi.encodeWithSelector(Initializable.InvalidInitialization.selector));

    s_burnMintERC20UUPS.initialize(
      s_name,
      s_symbol,
      s_decimals,
      s_maxSupply,
      s_preMint,
      s_defaultAdmin,
      s_defaultUpgrader
    );
  }

  /// @dev Adding _disableInitializers() function to implementation's constructor ensures that no one can call initialize directly on the implementation.
  /// @dev The initialize should be only callable through Proxy.
  /// @dev This test tests that case.
  function test_Initialize_RevertWhen_CallIsNotThroughProxy() public {
    BurnMintERC20UUPS newBurnMintERC20UUPS = new BurnMintERC20UUPS();

    vm.expectRevert(abi.encodeWithSelector(Initializable.InvalidInitialization.selector));

    newBurnMintERC20UUPS.initialize(
      s_name,
      s_symbol,
      s_decimals,
      s_maxSupply,
      s_preMint,
      s_defaultAdmin,
      s_defaultUpgrader
    );
  }
}
