// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Upgrades} from "../../../../../../vendor/openzeppelin-foundry-upgrades/v0.3.8/Upgrades.sol";
import {
  BurnMintERC20PausableUUPS,
  IAccessControl,
  IERC20
} from "../../../../../token/ERC20/upgradeable/BurnMintERC20PausableUUPS.sol";
import {BurnMintERC20PausableUUPSSetup} from "./BurnMintERC20PausableUUPSSetup.t.sol";

contract BurnMintERC20PausableUUPS_mint is BurnMintERC20PausableUUPSSetup {
  address s_mockPool;
  uint256 s_amount;

  function setUp() public virtual override {
    BurnMintERC20PausableUUPSSetup.setUp();

    s_mockPool = makeAddr("s_mockPool");

    s_burnMintERC20PausableUUPS.grantMintAndBurnRoles(s_mockPool);
    s_amount = 1e18;
  }

  function test_Mint() public {
    changePrank(s_mockPool);

    uint256 balanceBefore = s_burnMintERC20PausableUUPS.balanceOf(STRANGER);

    vm.expectEmit();
    emit IERC20.Transfer(address(0), STRANGER, s_amount);

    s_burnMintERC20PausableUUPS.mint(STRANGER, s_amount);

    assertEq(s_burnMintERC20PausableUUPS.balanceOf(STRANGER), balanceBefore + s_amount);
    assertEq(s_burnMintERC20PausableUUPS.totalSupply(), s_preMint + s_amount);
  }

  function test_Mint_RevertWhen_AmountExceedsMaxSupply() public {
    changePrank(s_mockPool);

    uint256 amountToMint = s_burnMintERC20PausableUUPS.maxSupply() + s_amount;

    vm.expectRevert(
      abi.encodeWithSelector(
        BurnMintERC20PausableUUPS.BurnMintERC20PausableUUPS__MaxSupplyExceeded.selector, amountToMint
      )
    );

    s_burnMintERC20PausableUUPS.mint(STRANGER, amountToMint);
  }

  function test_Mint_RevertWhen_CallerDoesNotHaveMinterRole() public {
    changePrank(STRANGER);

    vm.expectRevert(
      abi.encodeWithSelector(
        IAccessControl.AccessControlUnauthorizedAccount.selector, STRANGER, s_burnMintERC20PausableUUPS.MINTER_ROLE()
      )
    );

    s_burnMintERC20PausableUUPS.mint(STRANGER, s_amount);
  }

  function test_Mint_RevertWhen_ImplementationIsPaused() public {
    changePrank(s_defaultPauser);
    s_burnMintERC20PausableUUPS.pause();

    changePrank(s_mockPool);

    vm.expectRevert(abi.encodeWithSelector(BurnMintERC20PausableUUPS.BurnMintERC20PausableUUPS__Paused.selector));

    s_burnMintERC20PausableUUPS.mint(STRANGER, s_amount);
  }

  function test_Mint_RevertWhen_RecipientIsImplementationItself() public {
    changePrank(s_mockPool);

    vm.expectRevert(
      abi.encodeWithSelector(
        BurnMintERC20PausableUUPS.BurnMintERC20PausableUUPS__InvalidRecipient.selector,
        address(s_burnMintERC20PausableUUPS)
      )
    );

    s_burnMintERC20PausableUUPS.mint(address(s_burnMintERC20PausableUUPS), s_amount);
  }
}
