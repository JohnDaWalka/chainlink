// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IAccessControl, IERC20} from "../../../../../token/ERC20/upgradeable/BurnMintERC20UUPS.sol";
import {BurnMintERC20UUPSSetup} from "./BurnMintERC20UUPSSetup.t.sol";

contract BurnMintERC20UUPS_burn is BurnMintERC20UUPSSetup {
  address s_mockPool;
  uint256 s_amount;

  function setUp() public virtual override {
    BurnMintERC20UUPSSetup.setUp();

    s_mockPool = makeAddr("s_mockPool");

    s_burnMintERC20UUPS.grantMintAndBurnRoles(s_mockPool);
    s_amount = 1e18;
  }

  function test_Burn() public {
    changePrank(s_mockPool);
    s_burnMintERC20UUPS.mint(STRANGER, s_amount);

    uint256 amountToBurn = s_amount / 2;

    changePrank(STRANGER);
    s_burnMintERC20UUPS.transfer(s_mockPool, amountToBurn);

    changePrank(s_mockPool);

    uint256 balanceBefore = s_burnMintERC20UUPS.balanceOf(s_mockPool);
    uint256 totalSupplyBefore = s_burnMintERC20UUPS.totalSupply();

    vm.expectEmit();
    emit IERC20.Transfer(s_mockPool, address(0), amountToBurn);

    s_burnMintERC20UUPS.burn(amountToBurn);

    assertEq(s_burnMintERC20UUPS.balanceOf(s_mockPool), balanceBefore - amountToBurn);
    assertEq(s_burnMintERC20UUPS.totalSupply(), totalSupplyBefore - amountToBurn);
  }

  function test_Burn_RevertWhen_CallerDoesNotHaveBurnerRole() public {
    changePrank(s_mockPool);
    s_burnMintERC20UUPS.mint(STRANGER, s_amount);

    changePrank(STRANGER);

    vm.expectRevert(
      abi.encodeWithSelector(
        IAccessControl.AccessControlUnauthorizedAccount.selector, STRANGER, s_burnMintERC20UUPS.BURNER_ROLE()
      )
    );

    s_burnMintERC20UUPS.burn(s_amount);
  }
}
