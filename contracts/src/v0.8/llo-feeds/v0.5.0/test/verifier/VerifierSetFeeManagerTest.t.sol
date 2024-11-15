// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTest} from "./BaseVerifierTest.t.sol";
import {Verifier} from "../../Verifier.sol";

contract VerifierSetAccessControllerTest is BaseTest {
  event FeeManagerSet(address oldFeeManager, address newFeeManager);

  function test_revertsIfCalledByNonOwner() public {
    vm.expectRevert("Only callable by owner");
    changePrank(USER);
    s_verifier.setFeeManager(address(feeManager));
  }

  function test_successfullySetsNewFeeManager() public {
    vm.expectEmit(true, false, false, false);
    emit FeeManagerSet(address(0), ACCESS_CONTROLLER_ADDRESS);
    s_verifier.setFeeManager(address(feeManager));
    address ac = s_verifier.s_feeManager();
    assertEq(ac, address(feeManager));
  }

  function test_setFeeManagerWhichDoesntHonourInterface() public {
    vm.expectRevert(abi.encodeWithSelector(Verifier.FeeManagerInvalid.selector));
    s_verifier.setFeeManager(address(rewardManager));
  }
}
