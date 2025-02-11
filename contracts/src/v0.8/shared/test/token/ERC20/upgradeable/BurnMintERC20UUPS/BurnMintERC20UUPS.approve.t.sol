// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BurnMintERC20UUPS, IERC20} from "../../../../../token/ERC20/upgradeable/BurnMintERC20UUPS.sol";
import {BurnMintERC20UUPSSetup} from "./BurnMintERC20UUPSSetup.t.sol";

contract BurnMintERC20UUPS_approve is BurnMintERC20UUPSSetup {
  address s_mockPool;
  uint256 s_amount;

  function setUp() public virtual override {
    BurnMintERC20UUPSSetup.setUp();

    s_mockPool = makeAddr("s_mockPool");

    s_burnMintERC20UUPS.grantMintAndBurnRoles(s_mockPool);
    s_amount = 1e18;

    changePrank(s_mockPool);
    s_burnMintERC20UUPS.mint(STRANGER, s_amount);
  }

  function test_Approve() public {
    changePrank(STRANGER);

    vm.expectEmit();
    emit IERC20.Approval(STRANGER, s_mockPool, s_amount);

    s_burnMintERC20UUPS.approve(s_mockPool, s_amount);

    assertEq(s_burnMintERC20UUPS.allowance(STRANGER, s_mockPool), s_amount);
  }

  function test_Approve_RevertWhen_RecipientIsImplementationItself() public {
    changePrank(STRANGER);

    vm.expectRevert(
      abi.encodeWithSelector(
        BurnMintERC20UUPS.BurnMintERC20UUPS__InvalidRecipient.selector, address(s_burnMintERC20UUPS)
      )
    );

    s_burnMintERC20UUPS.approve(address(s_burnMintERC20UUPS), s_amount);
  }
}
