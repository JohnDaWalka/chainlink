// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BurnMintERC20PausableFreezableUUPS} from
  "../../../../../token/ERC20/upgradeable/BurnMintERC20PausableFreezableUUPS.sol";
import {IAccessControl} from "../../../../../token/ERC20/upgradeable/BurnMintERC20UUPS.sol";
import {BurnMintERC20PausableFreezableUUPSSetup} from "./BurnMintERC20PausableFreezableUUPSSetup.t.sol";

contract BurnMintERC20PausableFreezableUUPS_unfreeze is BurnMintERC20PausableFreezableUUPSSetup {
  uint256 s_amount = 1e18;

  function setUp() public override {
    super.setUp();

    changePrank(s_defaultFreezer);
    s_burnMintERC20PausableFreezableUUPS.freeze(OWNER);
  }

  function test_Unfreeze() public {
    changePrank(s_defaultFreezer);

    vm.expectEmit();
    emit BurnMintERC20PausableFreezableUUPS.AccountUnfrozen(OWNER);
    s_burnMintERC20PausableFreezableUUPS.unfreeze(OWNER);

    assertFalse(s_burnMintERC20PausableFreezableUUPS.isFrozen(OWNER));

    changePrank(s_defaultAdmin);
    s_burnMintERC20PausableFreezableUUPS.grantMintAndBurnRoles(s_defaultAdmin);
    s_burnMintERC20PausableFreezableUUPS.mint(OWNER, s_amount);
    assertEq(s_burnMintERC20PausableFreezableUUPS.balanceOf(OWNER), s_amount);
  }

  function test_Unfreeze_EvenWhenImplementationIsPaused() public {
    changePrank(s_defaultPauser);
    s_burnMintERC20PausableFreezableUUPS.pause();
    assertTrue(s_burnMintERC20PausableFreezableUUPS.paused());

    changePrank(s_defaultFreezer);
    s_burnMintERC20PausableFreezableUUPS.unfreeze(OWNER);
    assertFalse(s_burnMintERC20PausableFreezableUUPS.isFrozen(OWNER));
  }

  function test_Unfreeze_RevertWhen_CallerDoesNotHaveFreezerRole() public {
    changePrank(STRANGER);

    vm.expectRevert(
      abi.encodeWithSelector(
        IAccessControl.AccessControlUnauthorizedAccount.selector,
        STRANGER,
        s_burnMintERC20PausableFreezableUUPS.FREEZER_ROLE()
      )
    );

    s_burnMintERC20PausableFreezableUUPS.unfreeze(OWNER);
  }

  function test_Unfreeze_RevertWhen_AccountIsNotFrozen() public {
    changePrank(s_defaultFreezer);

    assertFalse(s_burnMintERC20PausableFreezableUUPS.isFrozen(STRANGER));

    vm.expectRevert(
      abi.encodeWithSelector(
        BurnMintERC20PausableFreezableUUPS.BurnMintERC20PausableFreezableUUPS__AccountNotFrozen.selector, STRANGER
      )
    );

    s_burnMintERC20PausableFreezableUUPS.unfreeze(STRANGER);
  }
}
