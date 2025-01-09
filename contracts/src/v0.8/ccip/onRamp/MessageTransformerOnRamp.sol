// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {OnRamp} from "./OnRamp.sol";
import {IMessageTransformer} from "../interfaces/IMessageTransformer.sol";
import {Internal} from "../libraries/Internal.sol";

contract MessageTransformerOnRamp is OnRamp {

  address internal s_messageTransformer;

  error ZeroAddressNotAllowed();

  constructor(
    StaticConfig memory staticConfig,
    DynamicConfig memory dynamicConfig,
    DestChainConfigArgs[] memory destChainConfigs,
    address messageTransformerAddr
  ) OnRamp(staticConfig, dynamicConfig, destChainConfigs) {
    if (address(messageTransformerAddr) == address(0)) {
      revert ZeroAddressNotAllowed();
    }
    s_messageTransformer = messageTransformerAddr;
  }

  function getMessageTransformerAddress() external view returns (address) {
    return s_messageTransformer;
  } 

  function _postProcessMessage(
    Internal.EVM2AnyRampMessage memory message
  ) internal override returns (Internal.EVM2AnyRampMessage memory transformedMessage) {
    try IMessageTransformer(s_messageTransformer).transformOutboundMessage(
      message
    ) returns (Internal.EVM2AnyRampMessage memory m) {
      transformedMessage = m;
    } catch (bytes memory err) {
      revert IMessageTransformer.MessageTransformError(err);
    }
    return transformedMessage;
  }
}