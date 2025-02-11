// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IAccessControl, IERC20} from "../../../../../token/ERC20/upgradeable/BurnMintERC20UUPS.sol";
import {BurnMintERC20UUPSSetup} from "./BurnMintERC20UUPSSetup.t.sol";

contract BurnMintERC20UUPS_burnFrom_alias is BurnMintERC20UUPSSetup {
  address s_mockPool;
  uint256 s_amount;

  function setUp() public virtual override {
    BurnMintERC20UUPSSetup.setUp();

    s_mockPool = makeAddr("s_mockPool");

    s_burnMintERC20UUPS.grantMintAndBurnRoles(s_mockPool);
    s_amount = 1e18;
  }

  function test_BurnFrom_alias() public {
    changePrank(s_mockPool);
    s_burnMintERC20UUPS.mint(STRANGER, s_amount);

    uint256 balanceBefore = s_burnMintERC20UUPS.balanceOf(STRANGER);
    uint256 totalSupplyBefore = s_burnMintERC20UUPS.totalSupply();
    uint256 amountToBurn = s_amount / 2;

    changePrank(STRANGER);
    s_burnMintERC20UUPS.approve(s_mockPool, amountToBurn);

    changePrank(s_mockPool);

    vm.expectEmit();
    emit IERC20.Transfer(STRANGER, address(0), amountToBurn);

    // burn(account, amount) is alias for burnFrom(account, amount)
    s_burnMintERC20UUPS.burn(STRANGER, amountToBurn);

    assertEq(s_burnMintERC20UUPS.balanceOf(STRANGER), balanceBefore - amountToBurn);
    assertEq(s_burnMintERC20UUPS.totalSupply(), totalSupplyBefore - amountToBurn);
  }

  function test_BurnFrom_alias_RevertWhen_CallerDoesNotHaveBurnerRole() public {
    changePrank(s_mockPool);
    s_burnMintERC20UUPS.mint(STRANGER, s_amount);

    changePrank(STRANGER);

    vm.expectRevert(
      abi.encodeWithSelector(
        IAccessControl.AccessControlUnauthorizedAccount.selector, STRANGER, s_burnMintERC20UUPS.BURNER_ROLE()
      )
    );

    s_burnMintERC20UUPS.burn(STRANGER, s_amount);
  }
}
