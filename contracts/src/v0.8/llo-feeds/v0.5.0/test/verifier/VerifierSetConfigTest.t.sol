// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {Common} from "../../../libraries/Common.sol";
import {Verifier} from "../../Verifier.sol";
import {BaseTest} from "./BaseVerifierTest.t.sol";

contract VerifierSetConfigTest is BaseTest {
  event ConfigUnset(bytes32 configDigest, address[] signers);
  event ConfigSet(bytes32 indexed configDigest, address[] signers, uint8 f, Common.AddressAndWeight[] recipientAddressesAndWeights);

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

  function test_settingDuplicateConfigFails() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](0);


    s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, signerAddrs, FAULT_TOLERANCE, weights);
    vm.expectRevert(abi.encodeWithSelector(Verifier.ConfigAlreadyExists.selector, DEFAULT_CONFIG_DIGEST));
    s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, signerAddrs, FAULT_TOLERANCE, weights);
  }

  function test_settingZeroConfigDigestFails() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](0);

    vm.expectRevert(abi.encodeWithSelector(Verifier.ZeroAddress.selector));
    s_verifier.setConfig(bytes32(0), signerAddrs, FAULT_TOLERANCE, weights);
  }

  function test_feeManagerCanBeRemoved() public {
    s_verifier.setFeeManager(address(0));
    assertEq(s_verifier.s_feeManager(), address(0));
  }

  function test_setConfigSucceedsAfterUnsetConfig() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](0);

    s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, signerAddrs, FAULT_TOLERANCE, weights);
    s_verifier.unsetConfig(DEFAULT_CONFIG_DIGEST, signerAddrs);
    s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, signerAddrs, FAULT_TOLERANCE, weights);
  }

  function test_unsetConfigFailsWithInvalidSigners() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](0);

    s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, signerAddrs, FAULT_TOLERANCE, weights);
    signerAddrs[0] = address(0);
    vm.expectRevert(abi.encodeWithSelector(Verifier.MismatchedSigners.selector));
    s_verifier.unsetConfig(DEFAULT_CONFIG_DIGEST, signerAddrs);
  }

  function test_unsetConfigFailsWithDuplicateSigners() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](0);

    s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, signerAddrs, FAULT_TOLERANCE, weights);
    signerAddrs[0] = signerAddrs[1];
    vm.expectRevert(abi.encodeWithSelector(Verifier.MismatchedSigners.selector));
    s_verifier.unsetConfig(DEFAULT_CONFIG_DIGEST, signerAddrs);
  }

  function test_unsetConfigFailsWithSubsetOfSigners() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](0);

    s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, signerAddrs, FAULT_TOLERANCE, weights);
    address[] memory subSignerAddres = new address[](2);
    subSignerAddres[0] = signerAddrs[0];
    subSignerAddres[1] = signerAddrs[1];
    vm.expectRevert(abi.encodeWithSelector(Verifier.MismatchedSigners.selector));
    s_verifier.unsetConfig(DEFAULT_CONFIG_DIGEST, subSignerAddres);
  }

  function test_unsetSignersEmitsEvent() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](0);

    s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, signerAddrs, FAULT_TOLERANCE, weights);

    vm.expectEmit();

    emit ConfigUnset(DEFAULT_CONFIG_DIGEST, signerAddrs);

    s_verifier.unsetConfig(DEFAULT_CONFIG_DIGEST, signerAddrs);
  }

  function test_setConfigEmitsEvents() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](0);

    vm.expectEmit();

    emit ConfigSet(DEFAULT_CONFIG_DIGEST, signerAddrs, FAULT_TOLERANCE, weights);

    s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, signerAddrs, FAULT_TOLERANCE, weights);
  }

  function test_unsetUnsetConfig() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](0);

    s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, signerAddrs, FAULT_TOLERANCE, weights);
    s_verifier.unsetConfig(DEFAULT_CONFIG_DIGEST, signerAddrs);
    vm.expectRevert(abi.encodeWithSelector(Verifier.ConfigDoesNotExist.selector));
    s_verifier.unsetConfig(DEFAULT_CONFIG_DIGEST, signerAddrs);
  }

  function test_onlyAdminCanUnsetConfig() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](0);

    s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, signerAddrs, FAULT_TOLERANCE, weights);
    changePrank(USER);
    vm.expectRevert("Only callable by owner");
    s_verifier.unsetConfig(DEFAULT_CONFIG_DIGEST, signerAddrs);
  }
}
