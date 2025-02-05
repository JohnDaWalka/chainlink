// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IMessageTransmitter} from "./IMessageTransmitter.sol";
import {ITokenMessenger} from "./ITokenMessenger.sol";

import {Ownable2StepMsgSender} from "../../../shared/access/Ownable2StepMsgSender.sol";

/// @title CCTP Message Transmitter Proxy
/// @notice A proxy contract for handling messages transmitted via the Cross Chain Transfer Protocol (CCTP).
/// @dev This contract is responsible for receiving messages from the `IMessageTransmitter` and ensuring only the Token Pool can invoke it.
contract CCTPMessageTransmitterProxy is Ownable2StepMsgSender {
  /// @notice Error thrown when a function is called by an unauthorized entity.
  error OnlyCallableByTokenPool();

  /// @notice Immutable reference to the `IMessageTransmitter` contract.
  IMessageTransmitter public immutable i_transmitter;

  /// @notice Address of the Token Pool that is allowed to call `receiveMessage`.
  address private s_tokenPool;

  /// @notice One-time cyclic dependency between TokenPool and MessageTransmitter.
  /// @dev The deployment sequence is:
  /// 1. Deploy MessageTransmitter first.
  /// 2. Deploy the upgraded TokenPool.
  /// 3. Set the TokenPool address in MessageTransmitter using `updateTokenPool`.
  constructor(
    ITokenMessenger tokenMessenger
  ) {
    i_transmitter = IMessageTransmitter(tokenMessenger.localMessageTransmitter());
  }

  /// @notice Receives a message from the `IMessageTransmitter` contract and validates it.
  /// @dev Can only be called by the Token Pool to process incoming messages.
  /// @param message The payload of the message being received.
  /// @param attestation The cryptographic proof validating the message.
  /// @return success A boolean indicating if the message was successfully processed.
  function receiveMessage(
    bytes calldata message,
    bytes calldata attestation
  ) external onlyTokenPool returns (bool success) {
    return i_transmitter.receiveMessage(message, attestation);
  }

  /// @notice Retrieves the address of the current Token Pool.
  /// @return The address of the Token Pool contract.
  function getTokenPool() external view returns (address) {
    return s_tokenPool;
  }

  /// @notice Updates the Token Pool address.
  /// @dev Can only be called by the contract owner.
  /// @param _tokenPool The new address of the Token Pool.
  function updateTokenPool(
    address _tokenPool
  ) external onlyOwner {
    s_tokenPool = _tokenPool;
  }

  /// @notice Ensures that only the authorized Token Pool can call certain functions.
  modifier onlyTokenPool() {
    if (msg.sender != s_tokenPool) revert OnlyCallableByTokenPool();
    _;
  }
}
