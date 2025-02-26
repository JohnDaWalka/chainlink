// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

/// @dev this file exposes structs that are defined in the Aptos CCIP contracts.
/// because Aptos does not have an abi.encode or abi.decode equivalent on-chain,
/// it has to be done manually. in the Aptos CCIP contracts, this is done using
/// Solidity ABI format, allowing us to use this file to automatically generate
/// bindings for these structs in order to encode and decode with type safety
/// in offchain code.
abstract contract AptosUtils {
  struct RampMessageHeader {
    bytes32 messageId;
    uint64 sourceChainSelector;
    uint64 destChainSelector;
    uint64 sequenceNumber;
    uint64 nonce;
  }

  struct Any2AptosTokenTransfer {
    bytes sourcePoolAddress;
    bytes32 destTokenAddress;
    uint32 destGasAmount;
    bytes extraData;
    uint256 amount;
  }

  struct Any2AptosRampMessage {
    RampMessageHeader header;
    bytes sender;
    bytes data;
    bytes32 receiver;
    uint256 gasLimit;
    Any2AptosTokenTransfer[] tokenAmounts;
  }

  struct ExecutionReport {
    uint64 sourceChainSelector;
    Any2AptosRampMessage message;
    bytes[] offchainTokenData;
    bytes32[] proofs;
  }

  struct PriceUpdates {
    TokenPriceUpdate[] tokenPriceUpdates;
    GasPriceUpdate[] gasPriceUpdates;
  }

  struct TokenPriceUpdate {
    bytes32 sourceToken;
    uint256 usdPerToken;
  }

  struct GasPriceUpdate {
    uint64 destChainSelector;
    uint256 usdPerUnitGas;
  }

  // solhint-disable-next-line gas-struct-packing
  struct MerkleRoot {
    uint64 sourceChainSelector;
    bytes onRampAddress;
    uint64 minSequenceNumber;
    uint64 maxSequenceNumber;
    bytes32 merkleRoot;
  }

  struct RMNSignature {
    bytes32 r;
    bytes32 s;
  }

  struct CommitReport {
    PriceUpdates priceUpdates;
    MerkleRoot[] blessedMerkleRoots;
    MerkleRoot[] unblessedMerkleRoots;
    RMNSignature[] rmnSignatures;
    bytes32 offrampAddress;
  }

  struct EVMExtraArgsV1 {
    uint256 gasLimit;
  }

  struct EVMExtraArgsV2 {
    uint256 gasLimit;
    bool allowOutOfOrderExecution;
  }

  struct SVMExtraArgsV1 {
    uint32 computeUnits;
    uint64 accountIsWritableBitmap;
    bool allowOutOfOrderExecution;
    bytes32 tokenReceiver;
    bytes32[] accounts;
  }

  /// @dev used to encode and decode commit reports.
  function exposeCommitReport(
    CommitReport memory commitReport
  ) external view virtual returns (bytes memory);

  /// @dev used to encode and decode execution reports.
  function exposeExecutionReport(
    ExecutionReport[] memory executionReport
  ) external view virtual returns (bytes memory);

  /// @dev used to encode/decode EVM extra args v1
  function exposeEVMExtraArgsV1(
    EVMExtraArgsV1 memory evmExtraArgsV1
  ) external view virtual returns (bytes memory);

  /// @dev used to encode/decode EVM extra args v2
  function exposeEVMExtraArgsV2(
    EVMExtraArgsV2 memory evmExtraArgsV2
  ) external view virtual returns (bytes memory);

  /// @dev used to encode/decode SVM extra args v1
  function exposeSVMExtraArgsV1(
    SVMExtraArgsV1 memory svmExtraArgsV1
  ) external view virtual returns (bytes memory);
}
