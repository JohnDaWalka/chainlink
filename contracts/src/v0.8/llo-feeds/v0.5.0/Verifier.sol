// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {TypeAndVersionInterface} from "../../interfaces/TypeAndVersionInterface.sol";
import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";

import {IAccessController} from "../../shared/interfaces/IAccessController.sol";
import {IERC165} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {Common} from "../libraries/Common.sol";
import {CommonV5} from "./libraries/CommonV5.sol";

import {IVerifier} from "./interfaces/IVerifier.sol";

import {IVerifierFeeManager} from "./interfaces/IVerifierFeeManager.sol";
import {IVerifierProxy} from "./interfaces/IVerifierProxy.sol";
import {IVerifierProxyVerifier} from "./interfaces/IVerifierProxyVerifier.sol";

// OCR2 standard
uint256 constant MAX_NUM_ORACLES = 31;

/**
 * @title Verifier
 * @author Michael Fletcher
 * @author ad0ll
 * @notice This contract will be used to verify reports based on the oracle signatures. This is not the source verifier which required individual fee configurations, instead, this checks that a report has been signed by one of the configured configDigest.
 */
contract Verifier is IVerifier, IVerifierProxyVerifier, ConfirmedOwner, TypeAndVersionInterface {
  /// @notice Mapping of configDigest to Config, used to look up the verification configuration for the configDigest of the incoming report
  mapping(bytes32 => CommonV5.Config) private s_Configs;

  /// @notice Array of all configs, used in a convenience view function that returns all configs
  CommonV5.Config[] private s_allConfigs;

  /// @notice The address of the verifierProxy
  address public s_feeManager;

  /// @notice The address of the access controller
  address public s_accessController;

  /// @notice The address of the verifierProxy
  IVerifierProxy public immutable i_verifierProxy;

  /// @notice This error is thrown whenever trying to set a config
  /// with a fault tolerance of 0
  error FaultToleranceMustBePositive();

  /// @notice This error is thrown whenever a report is signed
  /// with more than the max number of signers
  /// @param numSigners The number of signers who have signed the report
  /// @param maxSigners The maximum number of signers that can sign a report
  error ExcessSigners(uint256 numSigners, uint256 maxSigners);

  /// @notice This error is thrown whenever a report is signed or expected to be signed with less than the minimum number of signers
  /// @param numSigners The number of signers who have signed the report
  /// @param minSigners The minimum number of signers that need to sign a report
  error InsufficientSigners(uint256 numSigners, uint256 minSigners);

  /// @notice This error is thrown whenever a report is submitted with no signatures
  error NoSigners();

  /// @notice This error is thrown whenever a Config already exists
  /// @param configDigest The ID of the Config that already exists
  error ConfigAlreadyExists(bytes32 configDigest);

  /// @notice This error is thrown whenever the R and S signer components
  /// have different lengths
  /// @param rsLength The number of r signature components
  /// @param ssLength The number of s signature components
  error MismatchedSignatures(uint256 rsLength, uint256 ssLength);

  /// @notice This error is thrown whenever setting a config with duplicate signatures
  error NonUniqueSignatures();

  /* @notice This error is thrown whenever a report fails to verify. This error be thrown for multiple reasons and it's purposely like
   * this to prevent information being leaked about the verification process which could be used to enable free verifications maliciously
   */
  error BadVerification();

  /// @notice This error is thrown whenever a zero address is passed
  error ZeroAddress();

  /// @notice This error is thrown when the fee manager at an address does
  /// not conform to the fee manager interface
  error FeeManagerInvalid();

  /// @notice This error is thrown whenever an address tries
  /// to execute a verification that it is not authorized to do so
  error AccessForbidden();

  /// @notice This error is thrown whenever a config does not exist
  error ConfigDoesNotExist();

  /// @notice This error is thrown when you try to call _setConfig with a configDigest of 0
  error BadConfigDigest();

  error FunctionDeprecated();

  /// @notice This event is emitted when a new report is verified.
  /// It is used to keep a historical record of verified reports.
  event ReportVerified(bytes32 indexed feedId, address requester);

  /// @notice This event is emitted whenever a configuration is activated or deactivated
  event ConfigActivated(bytes32 configDigest, bool isActive);

  /// @notice event is emitted whenever a new Config is set
  event ConfigSet(
    bytes32 indexed configDigest, address[] signers, uint8 f, Common.AddressAndWeight[] recipientAddressesAndWeights
  );

  /// @notice This event is emitted when a new fee manager is set
  /// @param oldFeeManager The old fee manager address
  /// @param newFeeManager The new fee manager address
  event FeeManagerSet(address oldFeeManager, address newFeeManager);

  /// @notice This event is emitted when a new access controller is set
  /// @param oldAccessController The old access controller address
  /// @param newAccessController The new access controller address
  event AccessControllerSet(address oldAccessController, address newAccessController);

  constructor(
    address verifierProxy
  ) ConfirmedOwner(msg.sender) {
    if (verifierProxy == address(0)) {
      revert ZeroAddress();
    }

    i_verifierProxy = IVerifierProxy(verifierProxy);
  }

  /// @inheritdoc IVerifierProxyVerifier
  function verify(
    bytes calldata signedReport,
    bytes calldata parameterPayload,
    address sender
  ) external payable override onlyProxy checkAccess(sender) returns (bytes memory) {
    (bytes memory verifierResponse, bytes32 configDigest) = _verify(signedReport, sender);

    address fm = s_feeManager;
    if (fm != address(0)) {
      //process the fee and catch the error
      try IVerifierFeeManager(fm).processFee{value: msg.value}(configDigest, signedReport, parameterPayload, sender) {
        //do nothing
      } catch {
        // we purposefully obfuscate the error here to prevent information leaking leading to free verifications
        revert BadVerification();
      }
    }

    return verifierResponse;
  }

  /// @inheritdoc IVerifierProxyVerifier
  function verifyBulk(
    bytes[] calldata signedReports,
    bytes calldata parameterPayload,
    address sender
  ) external payable override onlyProxy checkAccess(sender) returns (bytes[] memory) {
    bytes[] memory verifierResponses = new bytes[](signedReports.length);
    bytes32[] memory donConfigs = new bytes32[](signedReports.length);

    for (uint256 i; i < signedReports.length; ++i) {
      (bytes memory report, bytes32 config) = _verify(signedReports[i], sender);
      verifierResponses[i] = report;
      donConfigs[i] = config;
    }

    address fm = s_feeManager;
    if (fm != address(0)) {
      //process the fee and catch the error
      try IVerifierFeeManager(fm).processFeeBulk{value: msg.value}(donConfigs, signedReports, parameterPayload, sender)
      {
        //do nothing
      } catch {
        // we purposefully obfuscate the error here to prevent information leaking leading to free verifications
        revert BadVerification();
      }
    }

    return verifierResponses;
  }

  function _verify(bytes calldata signedReport, address sender) internal returns (bytes memory, bytes32) {
    (bytes32[3] memory reportContext, bytes memory reportData, bytes32[] memory rs, bytes32[] memory ss, bytes32 rawVs)
    = abi.decode(signedReport, (bytes32[3], bytes, bytes32[], bytes32[], bytes32));

    // Signature lengths must match
    if (rs.length != ss.length) revert MismatchedSignatures(rs.length, ss.length);

    //Must always be at least 1 signer
    if (rs.length == 0) revert NoSigners();

    // The payload is hashed and signed by the configDigest - we need to recover the addresses
    bytes32 signedPayload = keccak256(abi.encodePacked(keccak256(reportData), reportContext));
    address[] memory signers = new address[](rs.length);
    for (uint256 i; i < rs.length; ++i) {
      signers[i] = ecrecover(signedPayload, uint8(rawVs[i]) + 27, rs[i], ss[i]);
    }

    // Duplicate signatures are not allowed
    if (Common._hasDuplicateAddresses(signers)) {
      revert BadVerification();
    }

    // Find the config for this report
    //TODO should this be storage?
    CommonV5.Config storage config = s_Configs[reportContext[0]]; //reportContext[0] is the config digest

    // Check a config has been set
    if (config.configDigest == bytes32(0)) {
      revert BadVerification();
    }

    //check the config is active
    if (!config.isActive) {
      revert BadVerification();
    }

    //check we have enough signatures
    if (signers.length <= config.f) {
      revert BadVerification();
    }

    //check each signer is registered against the config
    for (uint256 i; i < signers.length; ++i) {
      if (!config.oracles[signers[i]]) {
        revert BadVerification();
      }
    }

    emit ReportVerified(bytes32(reportData), sender);

    return (reportData, config.configDigest);
  }

  /// @inheritdoc IVerifier
  function setConfig(
    bytes32 configDigest,
    address[] memory signers,
    uint8 f,
    Common.AddressAndWeight[] memory recipientAddressesAndWeights
  ) external override checkConfigValid(signers.length, f) onlyOwner {
    _setConfig(configDigest, signers, f, recipientAddressesAndWeights);
  }

  function _setConfig(
    bytes32 configDigest,
    address[] memory signers,
    uint8 f,
    Common.AddressAndWeight[] memory recipientAddressesAndWeights
  ) internal {
    if (configDigest == bytes32(0)) {
      revert BadConfigDigest();
    }
    //TODO Check if i_verifierProxy is v0.3 verifier. If it is, then we need to set the config on the verifierProxy
    // Below doesn't work because interface is different between v0.3 and v0.4
    // if(i_verifierProxy.typeAndVersion() == "VerifierProxy 0.3.0"){
    //   i_verifierProxy.setVerifier(bytes32(0), configDigest, recipientAddressesAndWeights); //First param is currentConfigDigest, since all configs are new, this is always 0
    // }

    // Duplicate addresses would break protocol rules
    if (Common._hasDuplicateAddresses(signers)) {
      revert NonUniqueSignatures();
    }

    // Check the config we're setting isn't already set as the current active config as this will increase search costs unnecessarily when verifying historic reports
    // TODO TEMPORARILY COMMENTED OUT SO WE CAN SMOKE TEST WITHOUT HAVING TO DEPLOY A NEW CONFIG
    // UNCOMMENT BEFORE AUDIT
    if (s_Configs[configDigest].configDigest != bytes32(0)) {
      revert ConfigAlreadyExists(configDigest);
    }

    // We may want to register these later or skip this step in the unlikely scenario they've previously been registered in the RewardsManager
    if (recipientAddressesAndWeights.length != 0) {
      if (s_feeManager == address(0)) {
        revert FeeManagerInvalid();
      }
      IVerifierFeeManager(s_feeManager).setFeeRecipients(configDigest, recipientAddressesAndWeights);
    }

    // Store the config fields individually instead of struct assignment
    s_Configs[configDigest].configDigest = configDigest;
    s_Configs[configDigest].f = f;
    s_Configs[configDigest].isActive = true;

    // Generate signers mapping for the config, used to efficiently lookup whether a signer is allowed to appear when verifying a report
    for (uint256 i; i < signers.length; ++i) {
      if (signers[i] == address(0)) revert ZeroAddress();
      s_Configs[configDigest].oracles[signers[i]] = true;
    }
    // Note: the oracles mapping is already populated above

    //TODO Nested mapping giving me trouble
    // push the config to the convenience array used in getAllConfigs
    // s_allConfigs.push(s_Configs[configDigest]);

    emit ConfigSet(configDigest, signers, f, recipientAddressesAndWeights);
  }

  /// @inheritdoc IVerifier
  function setFeeManager(
    address feeManager
  ) external override onlyOwner {
    if (!IERC165(feeManager).supportsInterface(type(IVerifierFeeManager).interfaceId)) revert FeeManagerInvalid();

    address oldFeeManager = s_feeManager;
    s_feeManager = feeManager;

    emit FeeManagerSet(oldFeeManager, feeManager);
  }

  /// @inheritdoc IVerifier
  function setAccessController(
    address accessController
  ) external override onlyOwner {
    address oldAccessController = s_accessController;
    s_accessController = accessController;
    emit AccessControllerSet(oldAccessController, accessController);
  }

  /// @inheritdoc IVerifier
  function setConfigActive(bytes32 configDigest, bool isActive) external onlyOwner {
    // Fetch config
    CommonV5.Config storage config = s_Configs[configDigest];

    // Check config is set before update
    if (config.configDigest == bytes32(0)) {
      revert ConfigDoesNotExist();
    }
    config.isActive = isActive;

    emit ConfigActivated(config.configDigest, isActive);
  }

  //TODO Nested mappings giving me trouble
  // /// @inheritdoc IVerifier
  // function getAllConfigs() external view returns (CommonV5.Config[] memory) {
  //   return s_allConfigs;
  // }

  modifier checkConfigValid(uint256 numSigners, uint256 f) {
    if (f == 0) revert FaultToleranceMustBePositive();
    if (numSigners > MAX_NUM_ORACLES) revert ExcessSigners(numSigners, MAX_NUM_ORACLES);
    if (numSigners <= 3 * f) revert InsufficientSigners(numSigners, 3 * f + 1);
    _;
  }

  modifier onlyProxy() {
    if (address(i_verifierProxy) != msg.sender) {
      revert AccessForbidden();
    }
    _;
  }

  modifier checkAccess(
    address sender
  ) {
    address ac = s_accessController;
    if (address(ac) != address(0) && !IAccessController(ac).hasAccess(sender, msg.data)) revert AccessForbidden();
    _;
  }

  /// @inheritdoc IERC165
  function supportsInterface(
    bytes4 interfaceId
  ) public pure override returns (bool) {
    return interfaceId == type(IVerifier).interfaceId || interfaceId == type(IVerifierProxyVerifier).interfaceId;
  }

  /// @inheritdoc TypeAndVersionInterface
  function typeAndVersion() external pure override returns (string memory) {
    return "Verifier 0.5.0";
  }
}
