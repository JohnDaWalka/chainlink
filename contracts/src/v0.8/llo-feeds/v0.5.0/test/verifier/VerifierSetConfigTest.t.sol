// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {Common} from "../../../libraries/Common.sol";
import {Verifier} from "../../Verifier.sol";
import {BaseTest} from "./BaseVerifierTest.t.sol";

contract VerifierSetConfigTest is BaseTest {
  function setUp() public virtual override {
    BaseTest.setUp();
  }

  function test_revertsIfCalledByNonOwner() public {
    vm.expectRevert("Only callable by owner");
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    changePrank(USER);
    s_verifier.setConfig(
      DEFAULT_CONFIG_DIGEST,
      _getSignerAddresses(signers),
      FAULT_TOLERANCE,
      new Common.AddressAndWeight[](0)
    );
  }

  function test_revertsIfSetWithTooManySigners() public {
    address[] memory signers = new address[](MAX_ORACLES + 1);
    vm.expectRevert(abi.encodeWithSelector(Verifier.ExcessSigners.selector, signers.length, MAX_ORACLES));
    s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, signers, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
  }

  function test_revertsIfFaultToleranceIsZero() public {
    vm.expectRevert(abi.encodeWithSelector(Verifier.FaultToleranceMustBePositive.selector));
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, _getSignerAddresses(signers), 0, new Common.AddressAndWeight[](0));
  }

  function test_revertsIfNotEnoughSigners() public {
    address[] memory signers = new address[](2);
    signers[0] = address(1000);
    signers[1] = address(1001);

    vm.expectRevert(
      abi.encodeWithSelector(Verifier.InsufficientSigners.selector, signers.length, FAULT_TOLERANCE * 3 + 1)
    );
    s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, signers, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
  }

  function test_revertsIfDuplicateSigners() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    signerAddrs[0] = signerAddrs[1];
    vm.expectRevert(abi.encodeWithSelector(Verifier.NonUniqueSignatures.selector));
    s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
  }

  function test_revertsIfSignerContainsZeroAddress() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    signerAddrs[0] = address(0);
    vm.expectRevert(abi.encodeWithSelector(Verifier.ZeroAddress.selector));
    s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
  }

  function test_setConfigActiveUnknownConfigId() public {
    vm.expectRevert(abi.encodeWithSelector(Verifier.ConfigDoesNotExist.selector));
    s_verifier.setConfigActive(bytes32(uint256(3)), true);
  }
}
