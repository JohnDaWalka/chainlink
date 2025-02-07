// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BurnMintERC20PausableFreezableUUPS} from "../../../../../token/ERC20/upgradeable/BurnMintERC20PausableFreezableUUPS.sol";
import {IAccessControl} from "../../../../../token/ERC20/upgradeable/BurnMintERC20PausableUUPS.sol";
import {BurnMintERC20PausableFreezableUUPSSetup} from "./BurnMintERC20PausableFreezableUUPSSetup.t.sol";

contract BurnMintERC20PausableFreezableUUPS_freeze is BurnMintERC20PausableFreezableUUPSSetup {
  uint256 s_amount = 1e18;

  function setUp() public override {
    super.setUp();

    changePrank(s_defaultAdmin);
    s_burnMintERC20PausableFreezableUUPS.grantMintAndBurnRoles(s_defaultAdmin);
  }

  function test_Freeze() public {
    changePrank(s_defaultFreezer);

    vm.expectEmit();
    emit BurnMintERC20PausableFreezableUUPS.AccountFrozen(OWNER);
    s_burnMintERC20PausableFreezableUUPS.freeze(OWNER);

    assertTrue(s_burnMintERC20PausableFreezableUUPS.isFrozen(OWNER));
  }

  function test_Freeze_EvenWhenImplementationIsPaused() public {
    changePrank(s_defaultPauser);
    s_burnMintERC20PausableFreezableUUPS.pause();
    assertTrue(s_burnMintERC20PausableFreezableUUPS.paused());

    test_Freeze();
  }

  function test_Freeze_RevertWhen_CallerDoesNotHaveFreezerRole() public {
    changePrank(STRANGER);

    vm.expectRevert(
      abi.encodeWithSelector(
        IAccessControl.AccessControlUnauthorizedAccount.selector,
        STRANGER,
        s_burnMintERC20PausableFreezableUUPS.FREEZER_ROLE()
      )
    );

    s_burnMintERC20PausableFreezableUUPS.freeze(OWNER);
  }

  function test_Freeze_RevertWhen_RecipientIsImplementationItself() public {
    changePrank(s_defaultFreezer);

    vm.expectRevert(
      abi.encodeWithSelector(
        BurnMintERC20PausableFreezableUUPS.BurnMintERC20PausableFreezableUUPS__InvalidRecipient.selector,
        address(s_burnMintERC20PausableFreezableUUPS)
      )
    );

    s_burnMintERC20PausableFreezableUUPS.freeze(address(s_burnMintERC20PausableFreezableUUPS));
  }

  function test_Mint_RevertWhen_AccountIsFrozen() public {
    changePrank(s_defaultFreezer);
    s_burnMintERC20PausableFreezableUUPS.freeze(OWNER);

    changePrank(s_defaultAdmin);

    vm.expectRevert(
      abi.encodeWithSelector(
        BurnMintERC20PausableFreezableUUPS.BurnMintERC20PausableFreezableUUPS__AccountFrozen.selector,
        OWNER
      )
    );

    s_burnMintERC20PausableFreezableUUPS.mint(OWNER, s_amount);
  }

  function test_Approve_RevertWhen_AccountIsFrozen() public {
    changePrank(s_defaultAdmin);
    s_burnMintERC20PausableFreezableUUPS.mint(OWNER, s_amount);

    changePrank(s_defaultFreezer);
    s_burnMintERC20PausableFreezableUUPS.freeze(OWNER);

    changePrank(OWNER);

    vm.expectRevert(
      abi.encodeWithSelector(
        BurnMintERC20PausableFreezableUUPS.BurnMintERC20PausableFreezableUUPS__AccountFrozen.selector,
        OWNER
      )
    );

    s_burnMintERC20PausableFreezableUUPS.approve(STRANGER, s_amount);
  }
}
