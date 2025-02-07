// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {
  BurnMintERC20PausableUUPS,
  Initializable
} from "../../../../../token/ERC20/upgradeable/BurnMintERC20PausableUUPS.sol";
import {BurnMintERC20PausableUUPSSetup} from "./BurnMintERC20PausableUUPSSetup.t.sol";

contract BurnMintERC20PausableUUPS_initialize is BurnMintERC20PausableUUPSSetup {
  function test_Initialize() public view {
    assertEq(s_burnMintERC20PausableUUPS.name(), s_name);
    assertEq(s_burnMintERC20PausableUUPS.symbol(), s_symbol);
    assertEq(s_burnMintERC20PausableUUPS.decimals(), s_decimals);
    assertEq(s_burnMintERC20PausableUUPS.maxSupply(), s_maxSupply);
    assertEq(s_burnMintERC20PausableUUPS.totalSupply(), s_preMint);

    assertTrue(s_burnMintERC20PausableUUPS.hasRole(s_burnMintERC20PausableUUPS.DEFAULT_ADMIN_ROLE(), s_defaultAdmin));
    assertTrue(s_burnMintERC20PausableUUPS.hasRole(s_burnMintERC20PausableUUPS.UPGRADER_ROLE(), s_defaultUpgrader));
    assertTrue(s_burnMintERC20PausableUUPS.hasRole(s_burnMintERC20PausableUUPS.PAUSER_ROLE(), s_defaultPauser));
  }

  function test_Initialize_RevertWhen_AlreadyInitialized() public {
    vm.expectRevert(abi.encodeWithSelector(Initializable.InvalidInitialization.selector));

    s_burnMintERC20PausableUUPS.initialize(
      s_name, s_symbol, s_decimals, s_maxSupply, s_preMint, s_defaultAdmin, s_defaultUpgrader, s_defaultPauser
    );
  }

  /// @dev Adding _disableInitializers() function to implementation's constructor ensures that no one can call initialize directly on the implementation.
  /// @dev The initialize should be only callable through Proxy.
  /// @dev This test tests that case.
  function test_Initialize_RevertWhen_CallIsNotThroughProxy() public {
    BurnMintERC20PausableUUPS newBurnMintERC20PausableUUPS = new BurnMintERC20PausableUUPS();

    vm.expectRevert(abi.encodeWithSelector(Initializable.InvalidInitialization.selector));

    newBurnMintERC20PausableUUPS.initialize(
      s_name, s_symbol, s_decimals, s_maxSupply, s_preMint, s_defaultAdmin, s_defaultUpgrader, s_defaultPauser
    );
  }
}
