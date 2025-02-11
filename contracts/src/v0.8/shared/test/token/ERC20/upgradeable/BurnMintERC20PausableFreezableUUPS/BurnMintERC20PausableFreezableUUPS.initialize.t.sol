// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Initializable} from "../../../../../../vendor/openzeppelin-solidity-upgradeable/v5.0.2/contracts/proxy/utils/Initializable.sol";
import {BurnMintERC20PausableFreezableUUPS} from "../../../../../token/ERC20/upgradeable/BurnMintERC20PausableFreezableUUPS.sol";
import {BurnMintERC20PausableFreezableUUPSSetup} from "./BurnMintERC20PausableFreezableUUPSSetup.t.sol";

contract BurnMintERC20PausableFreezableUUPS_initialize is BurnMintERC20PausableFreezableUUPSSetup {
  function test_Initialize() public view {
    assertEq(s_burnMintERC20PausableFreezableUUPS.name(), s_name);
    assertEq(s_burnMintERC20PausableFreezableUUPS.symbol(), s_symbol);
    assertEq(s_burnMintERC20PausableFreezableUUPS.decimals(), s_decimals);
    assertEq(s_burnMintERC20PausableFreezableUUPS.maxSupply(), s_maxSupply);
    assertEq(s_burnMintERC20PausableFreezableUUPS.totalSupply(), s_preMint);

    assertTrue(
      s_burnMintERC20PausableFreezableUUPS.hasRole(
        s_burnMintERC20PausableFreezableUUPS.DEFAULT_ADMIN_ROLE(),
        s_defaultAdmin
      )
    );
    assertTrue(
      s_burnMintERC20PausableFreezableUUPS.hasRole(
        s_burnMintERC20PausableFreezableUUPS.UPGRADER_ROLE(),
        s_defaultUpgrader
      )
    );
    assertTrue(
      s_burnMintERC20PausableFreezableUUPS.hasRole(s_burnMintERC20PausableFreezableUUPS.PAUSER_ROLE(), s_defaultPauser)
    );
    assertTrue(
      s_burnMintERC20PausableFreezableUUPS.hasRole(
        s_burnMintERC20PausableFreezableUUPS.FREEZER_ROLE(),
        s_defaultFreezer
      )
    );
  }

  function test_Initialize_RevertWhen_AlreadyInitialized() public {
    vm.expectRevert(abi.encodeWithSelector(Initializable.InvalidInitialization.selector));

    s_burnMintERC20PausableFreezableUUPS.initialize(
      s_name,
      s_symbol,
      s_decimals,
      s_maxSupply,
      s_preMint,
      s_defaultAdmin,
      s_defaultUpgrader,
      s_defaultPauser,
      s_defaultFreezer
    );
  }

  /// @dev Adding _disableInitializers() function to implementation's constructor ensures that no one can call initialize directly on the implementation.
  /// @dev The initialize should be only callable through Proxy.
  /// @dev This test tests that case.
  function test_Initialize_RevertWhen_CallIsNotThroughProxy() public {
    BurnMintERC20PausableFreezableUUPS newBurnMintERC20PausableFreezableUUPS = new BurnMintERC20PausableFreezableUUPS();

    vm.expectRevert(abi.encodeWithSelector(Initializable.InvalidInitialization.selector));

    newBurnMintERC20PausableFreezableUUPS.initialize(
      s_name,
      s_symbol,
      s_decimals,
      s_maxSupply,
      s_preMint,
      s_defaultAdmin,
      s_defaultUpgrader,
      s_defaultPauser,
      s_defaultFreezer
    );
  }
}
