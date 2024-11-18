// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {Common} from "../../../libraries/Common.sol";
import {VerifierProxy} from "../../VerifierProxy.sol";
import {MultipleVerifierWithMultipleFeeManagers} from "./BaseVerifierTest.t.sol";
import {RewardManager} from "../../RewardManager.sol";

contract MultiVerifierBillingTests is MultipleVerifierWithMultipleFeeManagers {
  uint8 MINIMAL_FAULT_TOLERANCE = 2;
  address internal constant DEFAULT_RECIPIENT_1 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_1"))));
  address internal constant DEFAULT_RECIPIENT_2 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_2"))));
  address internal constant DEFAULT_RECIPIENT_3 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_3"))));
  address internal constant DEFAULT_RECIPIENT_4 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_4"))));
  address internal constant DEFAULT_RECIPIENT_5 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_5"))));
  address internal constant DEFAULT_RECIPIENT_6 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_6"))));
  address internal constant DEFAULT_RECIPIENT_7 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_7"))));

  bytes32[3] internal s_reportContext;
  V3Report internal s_testReport;

  function setUp() public virtual override {
    MultipleVerifierWithMultipleFeeManagers.setUp();
    s_reportContext[0] = DEFAULT_CONFIG_DIGEST;
    s_testReport = generateReportAtTimestamp(block.timestamp);
  }

  function _verify(
    VerifierProxy proxy,
    bytes memory payload,
    address feeAddress,
    uint256 wrappedNativeValue,
    address sender
  ) internal {
    address originalAddr = msg.sender;
    changePrank(sender);

    proxy.verify{value: wrappedNativeValue}(payload, abi.encode(feeAddress));

    changePrank(originalAddr);
  }

  function generateReportAtTimestamp(uint256 timestamp) public pure returns (V3Report memory) {
    return
      V3Report({
        feedId: FEED_ID_V3,
        observationsTimestamp: OBSERVATIONS_TIMESTAMP,
        validFromTimestamp: uint32(timestamp),
        nativeFee: uint192(DEFAULT_REPORT_NATIVE_FEE),
        linkFee: uint192(DEFAULT_REPORT_LINK_FEE),
        // ask michael about this expires at, is it usually set at what blocks
        expiresAt: uint32(timestamp) + 500,
        benchmarkPrice: MEDIAN,
        bid: BID,
        ask: ASK
      });
  }

  function payRecipients(bytes32 poolId, address[] memory recipients, address sender) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(sender);

    //pay the recipients
    rewardManager.payRecipients(poolId, recipients);

    //change back to the original address
    changePrank(originalAddr);
  }

  function test_multipleFeeManagersAndVerifiers() public {
    /*
       In this test we got:
        - three verifiers (verifier, verifier2, verifier3).
        - two fee managers (feeManager, feeManager2)
        - one reward manager
        
       we glue:
       - feeManager is used by verifier1 and verifier2
       - feeManager is used by verifier3
       - Rewardmanager is used by feeManager and feeManager2
      
      In this test we do verificatons via verifier1, verifier2 and verifier3 and check that rewards are set accordingly
   
    */
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](1);
    weights[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, ONE_PERCENT * 100);

    Common.AddressAndWeight[] memory weights2 = new Common.AddressAndWeight[](1);
    weights2[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_2, ONE_PERCENT * 100);

    Common.AddressAndWeight[] memory weights3 = new Common.AddressAndWeight[](1);
    weights3[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_3, ONE_PERCENT * 100);

    bytes32 DUMMY_CONFIG_DIGEST_1 = keccak256("DUMMY_CONFIG_DIGEST_1");
    bytes32 DUMMY_CONFIG_DIGEST_2 = keccak256("DUMMY_CONFIG_DIGEST_2");

    s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, signerAddrs, FAULT_TOLERANCE, weights);
    s_verifier2.setConfig(DUMMY_CONFIG_DIGEST_1, signerAddrs, MINIMAL_FAULT_TOLERANCE, weights2);
    s_verifier3.setConfig(DUMMY_CONFIG_DIGEST_2, signerAddrs, MINIMAL_FAULT_TOLERANCE + 1, weights3);

    bytes memory signedReport = _generateV3EncodedBlob(s_testReport, s_reportContext, signers);
    s_reportContext[0] = DUMMY_CONFIG_DIGEST_1;
    bytes memory signedReport2 = _generateV3EncodedBlob(s_testReport, s_reportContext, signers);
    s_reportContext[0] = DUMMY_CONFIG_DIGEST_2;
    bytes memory signedReport3 = _generateV3EncodedBlob(s_testReport, s_reportContext, signers);

    _approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);
    _verify(s_verifierProxy, signedReport, address(link), 0, USER);
    assertEq(link.balanceOf(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE);

    // internal state checks
    assertEq(feeManager.s_linkDeficit(DEFAULT_CONFIG_DIGEST), 0);
    assertEq(rewardManager.s_totalRewardRecipientFees(DEFAULT_CONFIG_DIGEST), DEFAULT_REPORT_LINK_FEE);
    assertEq(link.balanceOf(address(rewardManager)), DEFAULT_REPORT_LINK_FEE);

    // check the recipients are paid according to weights
    // These rewards happened through verifier1 and feeManager1
    address[] memory recipients = new address[](1);
    recipients[0] = DEFAULT_RECIPIENT_1;
    payRecipients(DEFAULT_CONFIG_DIGEST, recipients, ADMIN);
    assertEq(link.balanceOf(recipients[0]), DEFAULT_REPORT_LINK_FEE);
    assertEq(link.balanceOf(address(rewardManager)), 0);

    // these rewards happened through verifier2 and feeManager1
    address[] memory recipients2 = new address[](1);
    recipients2[0] = DEFAULT_RECIPIENT_2;
    _approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);
    _verify(s_verifierProxy2, signedReport2, address(link), 0, USER);
    payRecipients(DUMMY_CONFIG_DIGEST_1, recipients2, ADMIN);
    assertEq(link.balanceOf(recipients2[0]), DEFAULT_REPORT_LINK_FEE);
    assertEq(link.balanceOf(address(rewardManager)), 0);

    // these rewards happened through verifier3 and feeManager2
    address[] memory recipients3 = new address[](1);
    recipients3[0] = DEFAULT_RECIPIENT_3;
    _approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);
    _verify(s_verifierProxy3, signedReport3, address(link), 0, USER);
    payRecipients(DUMMY_CONFIG_DIGEST_2, recipients3, ADMIN);
    assertEq(link.balanceOf(recipients3[0]), DEFAULT_REPORT_LINK_FEE);
    assertEq(link.balanceOf(address(rewardManager)), 0);
  }

  function test_multipleFeeManagersAndVerifiersWithSameAddress() public {
    /*
       In this test we got:
        - three verifiers (verifier, verifier2, verifier3).
        - two fee managers (feeManager, feeManager2)
        - one reward manager

       we glue:
       - feeManager is used by verifier1 and verifier2
       - feeManager is used by verifier3
       - Rewardmanager is used by feeManager and feeManager2

      In this test we do verificatons via verifier1, verifier2 and verifier3 and check that rewards are set accordingly

    */
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](1);
    weights[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, ONE_PERCENT * 100);

    Common.AddressAndWeight[] memory weights2 = new Common.AddressAndWeight[](1);
    weights2[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_2, ONE_PERCENT * 100);

    Common.AddressAndWeight[] memory weights3 = new Common.AddressAndWeight[](1);
    weights3[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_3, ONE_PERCENT * 100);

    s_verifier.setConfig(DEFAULT_CONFIG_DIGEST, signerAddrs, FAULT_TOLERANCE, weights);

    // should fail with InvalidPoolId
    vm.expectRevert(abi.encodeWithSelector(RewardManager.InvalidPoolId.selector));

    s_verifier2.setConfig(DEFAULT_CONFIG_DIGEST, signerAddrs, FAULT_TOLERANCE, weights2);
  }
}
