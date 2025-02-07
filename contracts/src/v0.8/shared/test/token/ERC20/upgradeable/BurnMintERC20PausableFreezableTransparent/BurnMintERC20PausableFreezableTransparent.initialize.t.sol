// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Initializable} from "../../../../../../vendor/openzeppelin-solidity-upgradeable/v5.0.2/contracts/proxy/utils/Initializable.sol";
import {BurnMintERC20PausableFreezableTransparent} from "../../../../../token/ERC20/upgradeable/BurnMintERC20PausableFreezableTransparent.sol";
import {BurnMintERC20PausableFreezableTransparentSetup} from "./BurnMintERC20PausableFreezableTransparentSetup.t.sol";

contract BurnMintERC20PausableFreezableTransparent_initialize is BurnMintERC20PausableFreezableTransparentSetup {
  function test_Initialize() public view {
    assertEq(s_burnMintERC20PausableFreezableTransparent.name(), s_name);
    assertEq(s_burnMintERC20PausableFreezableTransparent.symbol(), s_symbol);
    assertEq(s_burnMintERC20PausableFreezableTransparent.decimals(), s_decimals);
    assertEq(s_burnMintERC20PausableFreezableTransparent.maxSupply(), s_maxSupply);
    assertEq(s_burnMintERC20PausableFreezableTransparent.totalSupply(), s_preMint);

    assertTrue(
      s_burnMintERC20PausableFreezableTransparent.hasRole(
        s_burnMintERC20PausableFreezableTransparent.DEFAULT_ADMIN_ROLE(),
        s_defaultAdmin
      )
    );
    assertTrue(
      s_burnMintERC20PausableFreezableTransparent.hasRole(
        s_burnMintERC20PausableFreezableTransparent.PAUSER_ROLE(),
        s_defaultPauser
      )
    );
    assertTrue(
      s_burnMintERC20PausableFreezableTransparent.hasRole(
        s_burnMintERC20PausableFreezableTransparent.FREEZER_ROLE(),
        s_defaultFreezer
      )
    );
  }

  function test_Initialize_RevertWhen_AlreadyInitialized() public {
    vm.expectRevert(abi.encodeWithSelector(Initializable.InvalidInitialization.selector));

    s_burnMintERC20PausableFreezableTransparent.initialize(
      s_name,
      s_symbol,
      s_decimals,
      s_maxSupply,
      s_preMint,
      s_defaultAdmin,
      s_defaultPauser,
      s_defaultFreezer
    );
  }

  /// @dev Adding _disableInitializers() function to implementation's constructor ensures that no one can call initialize directly on the implementation.
  /// @dev The initialize should be only callable through Proxy.
  /// @dev This test tests that case.
  function test_Initialize_RevertWhen_CallIsNotThroughProxy() public {
    BurnMintERC20PausableFreezableTransparent newBurnMintERC20PausableFreezableTransparent = new BurnMintERC20PausableFreezableTransparent();

    vm.expectRevert(abi.encodeWithSelector(Initializable.InvalidInitialization.selector));

    newBurnMintERC20PausableFreezableTransparent.initialize(
      s_name,
      s_symbol,
      s_decimals,
      s_maxSupply,
      s_preMint,
      s_defaultAdmin,
      s_defaultPauser,
      s_defaultFreezer
    );
  }
}
