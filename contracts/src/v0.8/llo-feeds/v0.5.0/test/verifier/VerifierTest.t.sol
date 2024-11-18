// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTest} from "./BaseVerifierTest.t.sol";
import {Verifier} from "../../Verifier.sol";
import {IVerifier} from "../../interfaces/IVerifier.sol";
import {IVerifierProxyVerifier} from "../../interfaces/IVerifierProxyVerifier.sol";
import {Common} from "../../../libraries/Common.sol";

contract VerifierConstructorTest is BaseTest {
  bytes32[3] internal s_reportContext;

  function test_revertsIfInitializedWithEmptyVerifierProxy() public {
    vm.expectRevert(abi.encodeWithSelector(Verifier.ZeroAddress.selector));
    new Verifier(address(0));
  }

  function test_typeAndVersion() public {
    Verifier v = new Verifier(address(s_verifierProxy));
    assertEq(v.owner(), ADMIN);
    string memory typeAndVersion = s_verifier.typeAndVersion();
    assertEq(typeAndVersion, "Verifier 0.5.0");
  }

  function test_falseIfIsNotCorrectInterface() public view {
    bool isInterface = s_verifier.supportsInterface(bytes4("abcd"));
    assertEq(isInterface, false);
  }

  function test_trueIfIsCorrectInterface() public view {
    bool isInterface = s_verifier.supportsInterface(type(IVerifier).interfaceId) &&
      s_verifier.supportsInterface(type(IVerifierProxyVerifier).interfaceId);
    assertEq(isInterface, true);
  }
}
