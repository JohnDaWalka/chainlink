// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IMessageTransmitter} from "./IMessageTransmitter.sol";
import {ITokenMessenger} from "./ITokenMessenger.sol";

import {Ownable2StepMsgSender} from "../../../shared/access/Ownable2StepMsgSender.sol";

contract CCTPMessageTransmitterProxy is Ownable2StepMsgSender {
  error OnlyCallableByTokenPool();

  IMessageTransmitter public immutable i_transmitter;
  address private s_tokenPool;

  // There's a one time cyclic dependency between the TokenPool and the MessageTransmitter.
  // We will deploy MessageTransmitter first and then deploy upgraded TokenPool.
  // We will set the address of the TokenPool to the MessageTransmitter using updateTokenPool.
  constructor(
    ITokenMessenger tokenMessenger
  ) {
    i_transmitter = IMessageTransmitter(tokenMessenger.localMessageTransmitter());
  }

  function receiveMessage(
    bytes calldata message,
    bytes calldata attestation
  ) external onlyTokenPool returns (bool success) {
    return i_transmitter.receiveMessage(message, attestation);
  }

  function getTokenPool() external view returns (address) {
    return s_tokenPool;
  }

  function updateTokenPool(
    address _tokenPool
  ) external onlyOwner {
    s_tokenPool = _tokenPool;
  }

  /// @dev only calls from the set router are accepted.
  modifier onlyTokenPool() {
    if (msg.sender != s_tokenPool) revert OnlyCallableByTokenPool();
    _;
  }
}
