// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BurnMintERC20UUPS, IAccessControl, IERC20} from "../../../../../token/ERC20/upgradeable/BurnMintERC20UUPS.sol";
import {BurnMintERC20UUPSSetup} from "./BurnMintERC20UUPSSetup.t.sol";

contract BurnMintERC20UUPS_mint is BurnMintERC20UUPSSetup {
  address s_mockPool;
  uint256 s_amount;

  function setUp() public virtual override {
    BurnMintERC20UUPSSetup.setUp();

    s_mockPool = makeAddr("s_mockPool");

    s_burnMintERC20UUPS.grantMintAndBurnRoles(s_mockPool);
    s_amount = 1e18;
  }

  function test_Mint() public {
    changePrank(s_mockPool);

    uint256 balanceBefore = s_burnMintERC20UUPS.balanceOf(STRANGER);

    vm.expectEmit();
    emit IERC20.Transfer(address(0), STRANGER, s_amount);

    s_burnMintERC20UUPS.mint(STRANGER, s_amount);

    assertEq(s_burnMintERC20UUPS.balanceOf(STRANGER), balanceBefore + s_amount);
    assertEq(s_burnMintERC20UUPS.totalSupply(), s_preMint + s_amount);
  }

  function test_Mint_RevertWhen_AmountExceedsMaxSupply() public {
    changePrank(s_mockPool);

    uint256 amountToMint = s_burnMintERC20UUPS.maxSupply() + s_amount;

    vm.expectRevert(
      abi.encodeWithSelector(BurnMintERC20UUPS.BurnMintERC20UUPS__MaxSupplyExceeded.selector, amountToMint)
    );

    s_burnMintERC20UUPS.mint(STRANGER, amountToMint);
  }

  function test_Mint_RevertWhen_CallerDoesNotHaveMinterRole() public {
    changePrank(STRANGER);

    vm.expectRevert(
      abi.encodeWithSelector(
        IAccessControl.AccessControlUnauthorizedAccount.selector,
        STRANGER,
        s_burnMintERC20UUPS.MINTER_ROLE()
      )
    );

    s_burnMintERC20UUPS.mint(STRANGER, s_amount);
  }

  function test_Mint_RevertWhen_RecipientIsImplementationItself() public {
    changePrank(s_mockPool);

    vm.expectRevert(
      abi.encodeWithSelector(
        BurnMintERC20UUPS.BurnMintERC20UUPS__InvalidRecipient.selector,
        address(s_burnMintERC20UUPS)
      )
    );

    s_burnMintERC20UUPS.mint(address(s_burnMintERC20UUPS), s_amount);
  }
}
