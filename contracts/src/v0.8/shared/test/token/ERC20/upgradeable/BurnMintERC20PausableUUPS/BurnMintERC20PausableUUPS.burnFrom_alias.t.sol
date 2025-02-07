// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BurnMintERC20PausableUUPS, IAccessControl, IERC20} from "../../../../../token/ERC20/upgradeable/BurnMintERC20PausableUUPS.sol";
import {BurnMintERC20PausableUUPSSetup} from "./BurnMintERC20PausableUUPSSetup.t.sol";

contract BurnMintERC20PausableUUPS_burnFrom_alias is BurnMintERC20PausableUUPSSetup {
  address s_mockPool;
  uint256 s_amount;

  function setUp() public virtual override {
    BurnMintERC20PausableUUPSSetup.setUp();

    s_mockPool = makeAddr("s_mockPool");

    s_burnMintERC20PausableUUPS.grantMintAndBurnRoles(s_mockPool);
    s_amount = 1e18;
  }

  function test_BurnFrom_alias() public {
    changePrank(s_mockPool);
    s_burnMintERC20PausableUUPS.mint(STRANGER, s_amount);

    uint256 balanceBefore = s_burnMintERC20PausableUUPS.balanceOf(STRANGER);
    uint256 totalSupplyBefore = s_burnMintERC20PausableUUPS.totalSupply();
    uint256 amountToBurn = s_amount / 2;

    changePrank(STRANGER);
    s_burnMintERC20PausableUUPS.approve(s_mockPool, amountToBurn);

    changePrank(s_mockPool);

    vm.expectEmit();
    emit IERC20.Transfer(STRANGER, address(0), amountToBurn);

    // burn(account, amount) is alias for burnFrom(account, amount)
    s_burnMintERC20PausableUUPS.burn(STRANGER, amountToBurn);

    assertEq(s_burnMintERC20PausableUUPS.balanceOf(STRANGER), balanceBefore - amountToBurn);
    assertEq(s_burnMintERC20PausableUUPS.totalSupply(), totalSupplyBefore - amountToBurn);
  }

  function test_BurnFrom_alias_RevertWhen_CallerDoesNotHaveBurnerRole() public {
    changePrank(s_mockPool);
    s_burnMintERC20PausableUUPS.mint(STRANGER, s_amount);

    changePrank(STRANGER);

    vm.expectRevert(
      abi.encodeWithSelector(
        IAccessControl.AccessControlUnauthorizedAccount.selector,
        STRANGER,
        s_burnMintERC20PausableUUPS.BURNER_ROLE()
      )
    );

    s_burnMintERC20PausableUUPS.burn(STRANGER, s_amount);
  }

  function test_BurnFrom_alias_RevertWhen_ImplementationIsPaused() public {
    changePrank(s_defaultPauser);
    s_burnMintERC20PausableUUPS.pause();

    changePrank(s_mockPool);

    vm.expectRevert(abi.encodeWithSelector(BurnMintERC20PausableUUPS.BurnMintERC20PausableUUPS__Paused.selector));

    s_burnMintERC20PausableUUPS.burn(STRANGER, 0);
  }
}
