// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {AccessControllerInterface} from "../../../../shared/interfaces/AccessControllerInterface.sol";
import {Common} from "../../../libraries/Common.sol";
import {Verifier} from "../../Verifier.sol";
import {BaseTest} from "./BaseVerifierTest.t.sol";

contract VerifierVerifyTest is BaseTest {
  bytes32[3] internal s_reportContext;
  V3Report internal s_testReportThree;

  function setUp() public virtual override {
    BaseTest.setUp();

    s_testReportThree = V3Report({
      feedId: FEED_ID_V3,
      observationsTimestamp: OBSERVATIONS_TIMESTAMP,
      validFromTimestamp: uint32(block.timestamp),
      nativeFee: uint192(DEFAULT_REPORT_NATIVE_FEE),
      linkFee: uint192(DEFAULT_REPORT_LINK_FEE),
      expiresAt: uint32(block.timestamp),
      benchmarkPrice: MEDIAN,
      bid: BID,
      ask: ASK
    });
  }

  function test_verifyReport() public {
    // Simple use case just setting a config and verifying a report
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    s_reportContext[0] = DEFAULT_CONFIG_DIGEST;

    s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));

    bytes memory signedReport = _generateV3EncodedBlob(s_testReportThree, s_reportContext, signers);

    bytes memory verifierResponse = s_verifierProxy.verify(signedReport, abi.encode(native));
    assertReportsEqual(verifierResponse, s_testReportThree);
  }


  function test_verifyReportWithMultipleConfigs() public {
    // Simple use case just setting a config and verifying a report
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    Signer[] memory signers2 = _getSigners(MAX_ORACLES);
    Signer[] memory signers3 = _getSigners(MAX_ORACLES);

    address[] memory signerAddrs = _getSignerAddresses(signers);
    address[] memory signerAddrs2 = _getSignerAddresses(signers2);
    address[] memory signerAddrs3 = _getSignerAddresses(signers3);

    s_reportContext[0] = DEFAULT_CONFIG_DIGEST;

    bytes32 DUMMY_CONFIG_1 = keccak256("DUMMY_CONFIG_1");
    bytes32 DUMMY_CONFIG_2 = keccak256("DUMMY_CONFIG_2");

    s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
    s_verifier.setConfig(DUMMY_CONFIG_1, signerAddrs2, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
    s_verifier.setConfig(DUMMY_CONFIG_2, signerAddrs3, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));

    //verify
    bytes memory signedReport = _generateV3EncodedBlob(s_testReportThree, s_reportContext, signers);
    s_verifierProxy.verify(signedReport, abi.encode(native));

    s_reportContext[0] = DUMMY_CONFIG_1;
    bytes memory signedReport2 = _generateV3EncodedBlob(s_testReportThree, s_reportContext, signers2);
    s_verifierProxy.verify(signedReport2, abi.encode(native));

    s_reportContext[0] = DUMMY_CONFIG_2;
    bytes memory signedReport3 = _generateV3EncodedBlob(s_testReportThree, s_reportContext, signers3);
    s_verifierProxy.verify(signedReport3, abi.encode(native));


    vm.expectRevert(abi.encodeWithSelector(Verifier.BadVerification.selector));
    s_reportContext[0] = keccak256("UNKNOWN_CONFIG");
    bytes memory signedReport4 = _generateV3EncodedBlob(s_testReportThree, s_reportContext, signers);
    s_verifierProxy.verify(signedReport4, abi.encode(native));
  }

  function test_verifyTogglingActiveFlagsDonConfigs() public {
    // sets config
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    s_reportContext[0] = DEFAULT_CONFIG_DIGEST;
    bytes memory signedReport = _generateV3EncodedBlob(s_testReportThree, s_reportContext, signers);
    s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
    // verifies report
    bytes memory verifierResponse = s_verifierProxy.verify(signedReport, abi.encode(native));
    assertReportsEqual(verifierResponse, s_testReportThree);

    // test verifying via a config that is deactivated
    s_verifier.setConfigActive(DEFAULT_CONFIG_DIGEST, false);
    vm.expectRevert(abi.encodeWithSelector(Verifier.BadVerification.selector));
    verifierResponse = s_verifierProxy.verify(signedReport, abi.encode(native));

    // test verifying via a reactivated config
    s_verifier.setConfigActive(DEFAULT_CONFIG_DIGEST, true);
    verifierResponse = s_verifierProxy.verify(signedReport, abi.encode(native));
    assertReportsEqual(verifierResponse, s_testReportThree);
  }

  function test_failToVerifyReportIfNotEnoughSigners() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);

    uint8 MINIMAL_FAULT_TOLERANCE = 2;
    BaseTest.Signer[] memory signersSubset1 = new BaseTest.Signer[](7);
    signersSubset1[0] = signers[0];
    signersSubset1[1] = signers[1];
    signersSubset1[2] = signers[2];
    signersSubset1[3] = signers[3];
    signersSubset1[4] = signers[4];
    signersSubset1[5] = signers[5];
    signersSubset1[6] = signers[6];
    address[] memory signersAddrSubset1 = _getSignerAddresses(signersSubset1);
    s_verifier.setConfig(
      DEFAULT_CONFIG_DIGEST,
      signersAddrSubset1,
      MINIMAL_FAULT_TOLERANCE,
      new Common.AddressAndWeight[](0)
    );

    // only one signer, signers < MINIMAL_FAULT_TOLERANCE
    BaseTest.Signer[] memory signersSubset2 = new BaseTest.Signer[](1);
    signersSubset2[0] = signers[4];

    bytes memory signedReport = _generateV3EncodedBlob(s_testReportThree, s_reportContext, signersSubset2);
    vm.expectRevert(abi.encodeWithSelector(Verifier.BadVerification.selector));
    s_verifierProxy.verify(signedReport, abi.encode(native));
  }

  function test_failToVerifyReportIfNoSigners() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);

    uint8 MINIMAL_FAULT_TOLERANCE = 2;
    BaseTest.Signer[] memory signersSubset1 = new BaseTest.Signer[](7);
    signersSubset1[0] = signers[0];
    signersSubset1[1] = signers[1];
    signersSubset1[2] = signers[2];
    signersSubset1[3] = signers[3];
    signersSubset1[4] = signers[4];
    signersSubset1[5] = signers[5];
    signersSubset1[6] = signers[6];
    address[] memory signersAddrSubset1 = _getSignerAddresses(signersSubset1);
    s_verifier.setConfig(
      DEFAULT_CONFIG_DIGEST,
      signersAddrSubset1,
      MINIMAL_FAULT_TOLERANCE,
      new Common.AddressAndWeight[](0)
    );

    // No signers for this report
    BaseTest.Signer[] memory signersSubset2 = new BaseTest.Signer[](0);
    bytes memory signedReport = _generateV3EncodedBlob(s_testReportThree, s_reportContext, signersSubset2);

    vm.expectRevert(abi.encodeWithSelector(Verifier.NoSigners.selector));
    s_verifierProxy.verify(signedReport, abi.encode(native));
  }

  function test_failToVerifyReportIfDupSigners() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);

    uint8 MINIMAL_FAULT_TOLERANCE = 2;
    BaseTest.Signer[] memory signersSubset1 = new BaseTest.Signer[](7);
    signersSubset1[0] = signers[0];
    signersSubset1[1] = signers[1];
    signersSubset1[2] = signers[2];
    signersSubset1[3] = signers[3];
    signersSubset1[4] = signers[4];
    signersSubset1[5] = signers[5];
    signersSubset1[6] = signers[6];
    address[] memory signersAddrSubset1 = _getSignerAddresses(signersSubset1);
    s_verifier.setConfig(
      DEFAULT_CONFIG_DIGEST,
      signersAddrSubset1,
      MINIMAL_FAULT_TOLERANCE,
      new Common.AddressAndWeight[](0)
    );
    // One signer is repeated
    BaseTest.Signer[] memory signersSubset2 = new BaseTest.Signer[](4);
    signersSubset2[0] = signers[0];
    signersSubset2[1] = signers[1];
    // repeated signers
    signersSubset2[2] = signers[2];
    signersSubset2[3] = signers[2];

    bytes memory signedReport = _generateV3EncodedBlob(s_testReportThree, s_reportContext, signersSubset2);

    vm.expectRevert(abi.encodeWithSelector(Verifier.BadVerification.selector));
    s_verifierProxy.verify(signedReport, abi.encode(native));
  }

  function test_failToVerifyReportIfSignerNotInConfig() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);

    uint8 MINIMAL_FAULT_TOLERANCE = 2;
    BaseTest.Signer[] memory signersSubset1 = new BaseTest.Signer[](7);
    signersSubset1[0] = signers[0];
    signersSubset1[1] = signers[1];
    signersSubset1[2] = signers[2];
    signersSubset1[3] = signers[3];
    signersSubset1[4] = signers[4];
    signersSubset1[5] = signers[5];
    signersSubset1[6] = signers[6];
    address[] memory signersAddrSubset1 = _getSignerAddresses(signersSubset1);
    s_verifier.setConfig(
      DEFAULT_CONFIG_DIGEST,
      signersAddrSubset1,
      MINIMAL_FAULT_TOLERANCE,
      new Common.AddressAndWeight[](0)
    );

    // one report whose signer is not in the config
    BaseTest.Signer[] memory reportSigners = new BaseTest.Signer[](4);
    // these signers are part ofm the config
    reportSigners[0] = signers[4];
    reportSigners[1] = signers[5];
    reportSigners[2] = signers[6];
    // this single signer is not in the config
    reportSigners[3] = signers[7];

    bytes memory signedReport = _generateV3EncodedBlob(s_testReportThree, s_reportContext, reportSigners);

    vm.expectRevert(abi.encodeWithSelector(Verifier.BadVerification.selector));
    s_verifierProxy.verify(signedReport, abi.encode(native));
  }

   function test_revertsVerifyIfNoAccess() public {
     vm.mockCall(
       ACCESS_CONTROLLER_ADDRESS,
       abi.encodeWithSelector(AccessControllerInterface.hasAccess.selector, USER),
       abi.encode(false)
     );
     bytes memory signedReport = _generateV3EncodedBlob(
       s_testReportThree,
       s_reportContext,
       _getSigners(FAULT_TOLERANCE + 1)
     );

     vm.expectRevert(abi.encodeWithSelector(Verifier.AccessForbidden.selector));

     changePrank(USER);
     s_verifier.verify(signedReport, abi.encode(native), msg.sender);
   }

   function test_verifyFailsWhenReportIsOlderThanConfig() public {
     /*
           - SetConfig A at time T0
           - SetConfig B at time T1
           - tries verifing report issued at blocktimestmap < T0
          */
     Signer[] memory signers = _getSigners(MAX_ORACLES);
     address[] memory signerAddrs = _getSignerAddresses(signers);
     s_reportContext[0] = bytes32(abi.encode(uint32(5), uint8(1)));

     vm.warp(block.timestamp + 100);

     V3Report memory reportAtTMinus100 = V3Report({
       feedId: FEED_ID_V3,
       observationsTimestamp: OBSERVATIONS_TIMESTAMP,
       validFromTimestamp: uint32(block.timestamp - 100),
       nativeFee: uint192(DEFAULT_REPORT_NATIVE_FEE),
       linkFee: uint192(DEFAULT_REPORT_LINK_FEE),
       expiresAt: uint32(block.timestamp),
       benchmarkPrice: MEDIAN,
       bid: BID,
       ask: ASK
     });


     s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
     vm.warp(block.timestamp + 100);

     bytes32 DUMMY_CONFIG_1 = keccak256("DUMMY_CONFIG_1");
     s_verifier.setConfig(DUMMY_CONFIG_1, signerAddrs, FAULT_TOLERANCE - 1, new Common.AddressAndWeight[](0));

     bytes memory signedReport = _generateV3EncodedBlob(reportAtTMinus100, s_reportContext, signers);

     vm.expectRevert(abi.encodeWithSelector(Verifier.BadVerification.selector));
     s_verifierProxy.verify(signedReport, abi.encode(native));
   }

   function test_configUnsetAndSetStillVerifies() public {
     Signer[] memory signers = _getSigners(MAX_ORACLES);
     address[] memory signerAddrs = _getSignerAddresses(signers);
     s_reportContext[0] = DEFAULT_CONFIG_DIGEST;

     s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
     s_verifier.unsetConfig(DEFAULT_CONFIG_DIGEST, signerAddrs);

     bytes memory signedReport = _generateV3EncodedBlob(s_testReportThree, s_reportContext, signers);

     vm.expectRevert(abi.encodeWithSelector(Verifier.BadVerification.selector));
     s_verifierProxy.verify(signedReport, abi.encode(native));

     s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
     s_verifierProxy.verify(signedReport, abi.encode(native));
   }

  function test_configUnsetAndResetWithDifferentKeysDoesNotVerifyWithPreviousKeys() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    s_reportContext[0] = DEFAULT_CONFIG_DIGEST;

    s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
    s_verifier.unsetConfig(DEFAULT_CONFIG_DIGEST, signerAddrs);

    bytes memory signedReport = _generateV3EncodedBlob(s_testReportThree, s_reportContext, signers);

    vm.expectRevert(abi.encodeWithSelector(Verifier.BadVerification.selector));
    s_verifierProxy.verify(signedReport, abi.encode(native));

    Signer[] memory signers2 = _getSecondarySigners(MAX_ORACLES);
    address[] memory signerAddrs2 = _getSignerAddresses(signers2);
    s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, signerAddrs2, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));

    vm.expectRevert(abi.encodeWithSelector(Verifier.BadVerification.selector));
    s_verifierProxy.verify(signedReport, abi.encode(native));

    signedReport = _generateV3EncodedBlob(s_testReportThree, s_reportContext, signers2);
    s_verifierProxy.verify(signedReport, abi.encode(native));

  }
}
