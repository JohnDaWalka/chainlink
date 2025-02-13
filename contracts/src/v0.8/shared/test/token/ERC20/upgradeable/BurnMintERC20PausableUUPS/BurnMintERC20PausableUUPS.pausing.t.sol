// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BurnMintERC20PausableUUPS} from "../../../../../token/ERC20/upgradeable/BurnMintERC20PausableUUPS.sol";
import {IERC20} from "../../../../../token/ERC20/upgradeable/BurnMintERC20UUPS.sol";
import {BurnMintERC20PausableUUPSSetup} from "./BurnMintERC20PausableUUPSSetup.t.sol";

import {PausableUpgradeable} from "../../../../../../vendor/openzeppelin-solidity-upgradeable/v5.0.2/contracts/utils/PausableUpgradeable.sol";

contract BurnMintERC20PausableUUPS_pausing is BurnMintERC20PausableUUPSSetup {
  address s_mockPool;
  uint256 s_amount;

  function setUp() public virtual override {
    BurnMintERC20PausableUUPSSetup.setUp();

    s_mockPool = makeAddr("s_mockPool");

    s_burnMintERC20PausableUUPS.grantMintAndBurnRoles(s_mockPool);
    s_amount = 1e18;
  }

  function test_Pause() public {
    changePrank(s_defaultPauser);

    vm.expectEmit();
    emit PausableUpgradeable.Paused(s_defaultPauser);

    s_burnMintERC20PausableUUPS.pause();

    assertTrue(s_burnMintERC20PausableUUPS.paused());
  }

  function test_Unpause() public {
    changePrank(s_defaultPauser);
    s_burnMintERC20PausableUUPS.pause();

    vm.expectEmit();
    emit PausableUpgradeable.Unpaused(s_defaultPauser);
    s_burnMintERC20PausableUUPS.unpause();

    assertFalse(s_burnMintERC20PausableUUPS.paused());

    changePrank(s_mockPool);
    vm.expectEmit();
    emit IERC20.Transfer(address(0), STRANGER, s_amount);
    s_burnMintERC20PausableUUPS.mint(STRANGER, s_amount);
  }

  function test_Mint_RevertWhen_ImplementationIsPaused() public {
    changePrank(s_defaultPauser);
    s_burnMintERC20PausableUUPS.pause();

    changePrank(s_mockPool);

    vm.expectRevert(abi.encodeWithSelector(PausableUpgradeable.EnforcedPause.selector));
    s_burnMintERC20PausableUUPS.mint(STRANGER, s_amount);
  }

  function test_Transfer_RevertWhen_ImplementationIsPaused() public {
    changePrank(s_mockPool);
    s_burnMintERC20PausableUUPS.mint(STRANGER, s_amount);

    changePrank(s_defaultPauser);
    s_burnMintERC20PausableUUPS.pause();

    changePrank(STRANGER);
    vm.expectRevert(abi.encodeWithSelector(PausableUpgradeable.EnforcedPause.selector));
    s_burnMintERC20PausableUUPS.transfer(OWNER, s_amount);
  }

  function test_Burn_RevertWhen_ImplementationIsPaused() public {
    changePrank(s_defaultPauser);
    s_burnMintERC20PausableUUPS.pause();

    changePrank(s_mockPool);

    vm.expectRevert(abi.encodeWithSelector(PausableUpgradeable.EnforcedPause.selector));

    s_burnMintERC20PausableUUPS.burn(0);
  }

  function test_BurnFrom_RevertWhen_ImplementationIsPaused() public {
    changePrank(s_defaultPauser);
    s_burnMintERC20PausableUUPS.pause();

    changePrank(s_mockPool);

    vm.expectRevert(abi.encodeWithSelector(PausableUpgradeable.EnforcedPause.selector));
    s_burnMintERC20PausableUUPS.burnFrom(STRANGER, 0);
  }

  function test_Approve_RevertWhen_ImplementationIsPaused() public {
    changePrank(s_defaultPauser);
    s_burnMintERC20PausableUUPS.pause();

    changePrank(STRANGER);

    vm.expectRevert(abi.encodeWithSelector(PausableUpgradeable.EnforcedPause.selector));

    s_burnMintERC20PausableUUPS.approve(s_mockPool, s_amount);
  }
}
