// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTest} from "./BaseVerifierTest.t.sol";
import {VerifierProxy} from "../../VerifierProxy.sol";

contract VerifierProxyInitializeVerifierTest is BaseTest {
  function test_setVerifierCalledByNoOwner() public {
    address STRANGER = address(999);
    changePrank(STRANGER);
    vm.expectRevert(bytes("Only callable by owner"));
    s_verifierProxy.setVerifier(address(s_verifier));
  }

  function test_setVerifierWhichDoesntHonourInterface() public {
    vm.expectRevert(abi.encodeWithSelector(VerifierProxy.VerifierInvalid.selector, address(rewardManager)));
    s_verifierProxy.setVerifier(address(rewardManager));
  }

  function test_setVerifierOk() public {
    s_verifierProxy.setVerifier(address(s_verifier));
    assertEq(s_verifierProxy.s_feeManager(), s_verifier.s_feeManager());
    assertEq(s_verifierProxy.s_accessController(), s_verifier.s_accessController());
  }

  function test_correctlySetsTheOwner() public {
    VerifierProxy proxy = new VerifierProxy();
    assertEq(proxy.owner(), ADMIN);
  }

  function test_correctlySetsVersion() public view {
    string memory version = s_verifierProxy.typeAndVersion();
    assertEq(version, "VerifierProxy 0.5.0");
  }
}
