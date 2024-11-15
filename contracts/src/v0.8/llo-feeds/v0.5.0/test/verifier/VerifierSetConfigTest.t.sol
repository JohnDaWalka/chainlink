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
      bytes32(uint256(1)), _getSignerAddresses(signers), FAULT_TOLERANCE, new Common.AddressAndWeight[](0)
    );
  }

  function test_revertsIfSetWithTooManySigners() public {
    address[] memory signers = new address[](MAX_ORACLES + 1);
    vm.expectRevert(abi.encodeWithSelector(Verifier.ExcessSigners.selector, signers.length, MAX_ORACLES));
    s_verifier.setConfig(bytes32(uint256(1)), signers, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
  }

  function test_revertsIfFaultToleranceIsZero() public {
    vm.expectRevert(abi.encodeWithSelector(Verifier.FaultToleranceMustBePositive.selector));
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    s_verifier.setConfig(bytes32(uint256(1)), _getSignerAddresses(signers), 0, new Common.AddressAndWeight[](0));
  }

  function test_revertsIfNotEnoughSigners() public {
    address[] memory signers = new address[](2);
    signers[0] = address(1000);
    signers[1] = address(1001);

    vm.expectRevert(
      abi.encodeWithSelector(Verifier.InsufficientSigners.selector, signers.length, FAULT_TOLERANCE * 3 + 1)
    );
    s_verifier.setConfig(bytes32(uint256(1)), signers, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
  }

  function test_revertsIfDuplicateSigners() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    signerAddrs[0] = signerAddrs[1];
    vm.expectRevert(abi.encodeWithSelector(Verifier.NonUniqueSignatures.selector));
    s_verifier.setConfig(bytes32(uint256(1)), signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
  }

  function test_revertsIfSignerContainsZeroAddress() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    signerAddrs[0] = address(0);
    vm.expectRevert(abi.encodeWithSelector(Verifier.ZeroAddress.selector));
    s_verifier.setConfig(bytes32(uint256(1)), signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
  }

  //TODO Is this function still relevant?
  // function test_donConfigIdIsSameForSignersInDifferentOrder() public {
  //   Signer[] memory signers = _getSigners(MAX_ORACLES);
  //   address[] memory signerAddrs = _getSignerAddresses(signers);

  //   bytes24 expectedDonConfigId = _donConfigIdFromConfigData(signerAddrs, FAULT_TOLERANCE);

  //   s_verifier.setConfig(bytes32(uint256(1)), signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
  //   vm.warp(block.timestamp + 1);

  //   address temp = signerAddrs[0];
  //   signerAddrs[0] = signerAddrs[1];
  //   signerAddrs[1] = temp;

  //   vm.expectRevert(abi.encodeWithSelector(Verifier.ConfigAlreadyExists.selector, expectedDonConfigId));

  //   s_verifier.setConfig(bytes32(uint256(1)), signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
  // }
  //TODO Function is not likely relevant anymore
  // function test_NoConfigAlreadyExists() public {
  //   Signer[] memory signers = _getSigners(MAX_ORACLES);
  //   address[] memory signerAddrs = _getSignerAddresses(signers);

  //   s_verifier.setConfig(bytes32(uint256(1)), signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));

  //   vm.warp(block.timestamp + 1);

  //   // testing adding same set of Signers but different FAULT_TOLERENCE does not result in ConfigAlreadyExists revert
  //   s_verifier.setConfig(bytes32(uint256(2)), signerAddrs, FAULT_TOLERANCE - 1, new Common.AddressAndWeight[](0));

  //   vm.warp(block.timestamp + 1);

  //   // testing adding a different set of Signers with same FAULT_TOLERENCE does not result in ConfigAlreadyExists revert
  //   address[] memory signerAddrsMinusOne = new address[](signerAddrs.length - 1);
  //   for (uint256 i = 0; i < signerAddrs.length - 1; i++) {
  //     signerAddrsMinusOne[i] = signerAddrs[i];
  //   }
  //   s_verifier.setConfig(
  //     bytes32(uint256(1)), signerAddrsMinusOne, FAULT_TOLERANCE - 1, new Common.AddressAndWeight[](0)
  //   );
  // }

  //TODO Is this function still relevant?
  // function test_addressesAndWeightsDoNotProduceSideEffectsInDonConfigIds() public {
  //   Signer[] memory signers = _getSigners(MAX_ORACLES);
  //   address[] memory signerAddrs = _getSignerAddresses(signers);

  //   s_verifier.setConfig(bytes32(uint256(1)), signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
  //   vm.warp(block.timestamp + 1);

  //   bytes24 expectedDonConfigId = _donConfigIdFromConfigData(signerAddrs, FAULT_TOLERANCE);

  //   vm.expectRevert(abi.encodeWithSelector(Verifier.ConfigAlreadyExists.selector, expectedDonConfigId));

  //   // Same call to setConfig with different addressAndWeights do not entail a new DonConfigID
  //   // Resulting in a ConfigAlreadyExists error
  //   Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](1);
  //   weights[0] = Common.AddressAndWeight(signers[0].signerAddress, 1);
  //   s_verifier.setConfig(bytes32(uint256(1)), signerAddrs, FAULT_TOLERANCE, weights);
  // }

  function test_setConfigActiveUnknownConfigId() public {
    vm.expectRevert(abi.encodeWithSelector(Verifier.ConfigDoesNotExist.selector));
    s_verifier.setConfigActive(bytes32(uint256(3)), true);
  }
}
