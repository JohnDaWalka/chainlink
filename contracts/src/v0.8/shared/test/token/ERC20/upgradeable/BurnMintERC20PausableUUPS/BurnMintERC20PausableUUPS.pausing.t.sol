// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BurnMintERC20PausableUUPS, IERC20} from "../../../../../token/ERC20/upgradeable/BurnMintERC20PausableUUPS.sol";
import {BurnMintERC20PausableUUPSSetup} from "./BurnMintERC20PausableUUPSSetup.t.sol";

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
    emit BurnMintERC20PausableUUPS.Paused();
    s_burnMintERC20PausableUUPS.pause();

    assertTrue(s_burnMintERC20PausableUUPS.paused());

    changePrank(s_mockPool);
    vm.expectRevert(abi.encodeWithSelector(BurnMintERC20PausableUUPS.BurnMintERC20PausableUUPS__Paused.selector));
    s_burnMintERC20PausableUUPS.mint(STRANGER, s_amount);
  }

  function test_Unpause() public {
    changePrank(s_defaultPauser);
    s_burnMintERC20PausableUUPS.pause();

    vm.expectEmit();
    emit BurnMintERC20PausableUUPS.Unpaused();
    s_burnMintERC20PausableUUPS.unpause();

    assertFalse(s_burnMintERC20PausableUUPS.paused());

    changePrank(s_mockPool);
    vm.expectEmit();
    emit IERC20.Transfer(address(0), STRANGER, s_amount);
    s_burnMintERC20PausableUUPS.mint(STRANGER, s_amount);
  }
}
