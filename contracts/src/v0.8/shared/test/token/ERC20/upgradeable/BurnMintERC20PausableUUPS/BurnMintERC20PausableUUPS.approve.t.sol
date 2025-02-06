// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BurnMintERC20PausableUUPS, IERC20} from "../../../../../token/ERC20/upgradeable/BurnMintERC20PausableUUPS.sol";
import {BurnMintERC20PausableUUPSSetup} from "./BurnMintERC20PausableUUPSSetup.t.sol";

contract BurnMintERC20PausableUUPS_approve is BurnMintERC20PausableUUPSSetup {
  address s_mockPool;
  uint256 s_amount;

  function setUp() public virtual override {
    BurnMintERC20PausableUUPSSetup.setUp();

    s_mockPool = makeAddr("s_mockPool");

    s_burnMintERC20PausableUUPS.grantMintAndBurnRoles(s_mockPool);
    s_amount = 1e18;

    changePrank(s_mockPool);
    s_burnMintERC20PausableUUPS.mint(STRANGER, s_amount);
  }

  function test_Approve() public {
    changePrank(STRANGER);

    vm.expectEmit();
    emit IERC20.Approval(STRANGER, s_mockPool, s_amount);

    s_burnMintERC20PausableUUPS.approve(s_mockPool, s_amount);

    assertEq(s_burnMintERC20PausableUUPS.allowance(STRANGER, s_mockPool), s_amount);
  }

  function test_Approve_RevertWhen_ImplementationIsPaused() public {
    changePrank(s_defaultPauser);
    s_burnMintERC20PausableUUPS.pause();

    changePrank(STRANGER);

    vm.expectRevert(abi.encodeWithSelector(BurnMintERC20PausableUUPS.BurnMintERC20PausableUUPS__Paused.selector));

    s_burnMintERC20PausableUUPS.approve(s_mockPool, s_amount);
  }

  function test_Approve_RevertWhen_RecipientIsImplementationItself() public {
    changePrank(STRANGER);

    vm.expectRevert(
      abi.encodeWithSelector(
        BurnMintERC20PausableUUPS.BurnMintERC20PausableUUPS__InvalidRecipient.selector,
        address(s_burnMintERC20PausableUUPS)
      )
    );

    s_burnMintERC20PausableUUPS.approve(address(s_burnMintERC20PausableUUPS), s_amount);
  }
}
