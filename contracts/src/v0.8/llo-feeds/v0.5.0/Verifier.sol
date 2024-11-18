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
import {IVerifierProxyV03} from "./interfaces/IVerifierProxyV03.sol";
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
  mapping(bytes32 => CommonV5.Config) private s_configs;

  /// @notice Array of all configs, used in a convenience view function that returns all configs
  bytes32[] private s_allConfigs;

  /// @notice The address of the verifierProxy
  address public s_feeManager;

  /// @notice The address of the access controller
  address public s_accessController;

  /// @notice The address of the V0.3 verifierProxy
  address public immutable i_verifierProxy;

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

  /// @notice This error is thrown when the proxy is neither a valid v0.3 or v0.4 proxy
  error VerifierProxyInvalid();

  /// @notice This error is thrown whenever an address tries
  /// to execute a verification that it is not authorized to do so
  error AccessForbidden();

  /// @notice This error is thrown whenever a config does not exist
  error ConfigDoesNotExist();

  /// @notice This error is thrown when you try to call _setConfig with a configDigest of 0
  error BadConfigDigest();

  /// @notice This error is thrown when trying to access a config that is out of bounds externally
  error InvalidIndex();

  /// @notice This event is emitted when a new report is verified.
  /// It is used to keep a historical record of verified reports.
  event ReportVerified(bytes32 indexed feedId, address requester);

  /// @notice This event is emitted whenever a configuration is activated or deactivated
  event ConfigActivated(bytes32 configDigest, bool isActive);

  /// @notice event is emitted whenever a new Config is set
  event ConfigSet(
    bytes32 indexed configDigest,
    address[] signers,
    uint8 f,
    Common.AddressAndWeight[] recipientAddressesAndWeights
  );

  /// @notice This event is emitted when a new fee manager is set
  /// @param oldFeeManager The old fee manager address
  /// @param newFeeManager The new fee manager address
  event FeeManagerSet(address oldFeeManager, address newFeeManager);

  /// @notice This event is emitted when a new access controller is set
  /// @param oldAccessController The old access controller address
  /// @param newAccessController The new access controller address
  event AccessControllerSet(address oldAccessController, address newAccessController);

  bytes32 constant V03_PROXY_TYPE_AND_VERSION = keccak256(bytes("VerifierProxy 2.0.0"));

  constructor(address verifierProxy) ConfirmedOwner(msg.sender) {
    if (verifierProxy == address(0)) {
      revert ZeroAddress();
    }

    // Proxy should support TypeAndVersion as we need to identify which proxy is calling
    if(!IERC165(verifierProxy).supportsInterface(type(TypeAndVersionInterface).interfaceId))
      revert VerifierProxyInvalid();

    // If it's the v0.3 Proxy check it implements the V03 Interface
    if (keccak256(bytes(TypeAndVersionInterface(verifierProxy).typeAndVersion())) == V03_PROXY_TYPE_AND_VERSION) {
      if(!IERC165(verifierProxy).supportsInterface(type(IVerifierProxyV03).interfaceId))
        revert VerifierProxyInvalid();
    }

    i_verifierProxy = verifierProxy;
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
      try
        IVerifierFeeManager(fm).processFeeBulk{value: msg.value}(donConfigs, signedReports, parameterPayload, sender)
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
    (
      bytes32[3] memory reportContext,
      bytes memory reportData,
      bytes32[] memory rs,
      bytes32[] memory ss,
      bytes32 rawVs
    ) = abi.decode(signedReport, (bytes32[3], bytes, bytes32[], bytes32[], bytes32));

    // Signature lengths must match
    if (rs.length != ss.length) revert MismatchedSignatures(rs.length, ss.length);

    //Must always be at least 1 signer
    if (rs.length == 0) revert NoSigners();

    // The payload is hashed and signed by the oracles - we need to recover the addresses
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
    CommonV5.Config storage config = s_configs[reportContext[0]]; //reportContext[0] is the config digest

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

    // If it's a v0.3 interface, register the verifier
     if(keccak256(bytes(TypeAndVersionInterface(i_verifierProxy).typeAndVersion())) == V03_PROXY_TYPE_AND_VERSION){
       IVerifierProxyV03(i_verifierProxy).setVerifier(bytes32(0), configDigest, recipientAddressesAndWeights); //First param is currentConfigDigest, since all configs are new, this is always 0
     }

    // Duplicate addresses would break protocol rules
    if (Common._hasDuplicateAddresses(signers)) {
      revert NonUniqueSignatures();
    }

    // Check the config does not already exist
    if (s_configs[configDigest].configDigest != bytes32(0)) {
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
    s_configs[configDigest].configDigest = configDigest;
    s_configs[configDigest].f = f;
    s_configs[configDigest].isActive = true;

    // Generate signers mapping for the config, used to efficiently lookup whether a signer is allowed to appear when verifying a report
    for (uint256 i; i < signers.length; ++i) {
      if (signers[i] == address(0)) revert ZeroAddress();
      s_configs[configDigest].oracles[signers[i]] = true;
    }

    // Keep track of all digests to enable easier querying off-chain
    s_allConfigs.push(configDigest);

    emit ConfigSet(configDigest, signers, f, recipientAddressesAndWeights);
  }

  /// @inheritdoc IVerifier
  function setFeeManager(address feeManager) external override onlyOwner {
    if (!IERC165(feeManager).supportsInterface(type(IVerifierFeeManager).interfaceId)) revert FeeManagerInvalid();

    address oldFeeManager = s_feeManager;
    s_feeManager = feeManager;

    emit FeeManagerSet(oldFeeManager, feeManager);
  }

  /// @inheritdoc IVerifier
  function setAccessController(address accessController) external override onlyOwner {
    address oldAccessController = s_accessController;
    s_accessController = accessController;
    emit AccessControllerSet(oldAccessController, accessController);
  }

  /// @inheritdoc IVerifier
  function setConfigActive(bytes32 configDigest, bool isActive) external onlyOwner {
    // Fetch config
    CommonV5.Config storage config = s_configs[configDigest];

    // Check config is set before update
    if (config.configDigest == bytes32(0)) {
      revert ConfigDoesNotExist();
    }
    config.isActive = isActive;

    emit ConfigActivated(config.configDigest, isActive);
  }

  modifier checkConfigValid(uint256 numSigners, uint256 f) {
    if (f == 0) revert FaultToleranceMustBePositive();
    if (numSigners > MAX_NUM_ORACLES) revert ExcessSigners(numSigners, MAX_NUM_ORACLES);
    if (numSigners <= 3 * f) revert InsufficientSigners(numSigners, 3 * f + 1);
    _;
  }

  modifier onlyProxy() {
    if (i_verifierProxy != msg.sender) {
      revert AccessForbidden();
    }
    _;
  }

  modifier checkAccess(address sender) {
    address ac = s_accessController;
    if (address(ac) != address(0) && !IAccessController(ac).hasAccess(sender, msg.data)) revert AccessForbidden();
    _;
  }

  /// @inheritdoc IERC165
  function supportsInterface(bytes4 interfaceId) public pure override returns (bool) {
    return interfaceId == type(IVerifier).interfaceId || interfaceId == type(IVerifierProxyVerifier).interfaceId;
  }

  /// @inheritdoc TypeAndVersionInterface
  function typeAndVersion() external pure override returns (string memory) {
    return "Verifier 0.5.0";
  }

  //  /// Utility function to get all configs off-chain
  // TODO should we expose?
  function getAllConfigs(
    uint256 startIndex,
    uint256 endIndex
  ) external view returns (bytes32[] memory) {
    // Only EOA can read configs
    if (msg.sender != tx.origin) revert AccessForbidden();

    // Check bounds
    if (startIndex > endIndex) revert InvalidIndex();
    if (endIndex >= s_allConfigs.length) revert InvalidIndex();

    // Calculate size of return array
    uint256 size = endIndex - startIndex + 1;
    bytes32[] memory configs = new bytes32[](size);

    // Copy requested range
    for (uint256 i; i < size; ++i) {
      configs[i] = s_allConfigs[startIndex + i];
    }

    return configs;
  }
}
