// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IMessageTransformer} from "../interfaces/IMessageTransformer.sol";
import {Internal} from "../libraries/Internal.sol";
import {OffRamp} from "./OffRamp.sol";

contract MessageTransformerOffRamp is OffRamp {
  address internal s_messageTransformer;

  constructor(
    StaticConfig memory staticConfig,
    DynamicConfig memory dynamicConfig,
    SourceChainConfigArgs[] memory sourceChainConfigs,
    address messageTransformerAddr
  ) OffRamp(staticConfig, dynamicConfig, sourceChainConfigs) {
    if (address(messageTransformerAddr) == address(0)) {
      revert ZeroAddressNotAllowed();
    }
    s_messageTransformer = messageTransformerAddr;
  }

  function getMessageTransformerAddress() external view returns (address) {
    return s_messageTransformer;
  }

  function _beforeExecuteSingleMessage(
    Internal.Any2EVMRampMessage memory message
  ) internal override returns (Internal.Any2EVMRampMessage memory transformedMessage) {
    try IMessageTransformer(s_messageTransformer).transformInboundMessage(message) returns (
      Internal.Any2EVMRampMessage memory m
    ) {
      transformedMessage = m;
    } catch (bytes memory err) {
      revert IMessageTransformer.MessageTransformError(err);
    }
    return transformedMessage;
  }
}
