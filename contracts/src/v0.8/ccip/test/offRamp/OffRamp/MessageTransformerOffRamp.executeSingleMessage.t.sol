// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {OffRamp} from "../../../offRamp/OffRamp.sol";
import {MessageTransformerOffRamp} from "../../../offRamp/MessageTransformerOffRamp.sol";
import {OffRampSetup} from "./OffRampSetup.t.sol";
import {Internal} from "../../../libraries/Internal.sol";
import {IMessageTransformer} from "../../../interfaces/IMessageTransformer.sol";
import {MessageTransformerHelper} from "../../helpers/MessageTransformerHelper.sol";
import {MultiOCR3Base} from "../../../ocr/MultiOCR3Base.sol";
import {Router} from "../../../Router.sol";

contract MessageTransformerOffRamp_executeSingleMessage is OffRampSetup {

  MessageTransformerOffRamp internal s_messageTransformerOffRamp;


  function setUp() public virtual override {
    super.setUp();
    s_messageTransformerOffRamp = new MessageTransformerOffRamp(
      s_offRamp.getStaticConfig(),
      s_offRamp.getDynamicConfig(),
      new OffRamp.SourceChainConfigArgs[](0),
      address(s_inboundMessageTransformer)
    );
    
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true
    });
    s_messageTransformerOffRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](0);
    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](2 * sourceChainConfigs.length);

    for (uint256 i = 0; i < sourceChainConfigs.length; ++i) {
      uint64 sourceChainSelector = sourceChainConfigs[i].sourceChainSelector;

      offRampUpdates[2 * i] = Router.OffRamp({sourceChainSelector: sourceChainSelector, offRamp: address(s_messageTransformerOffRamp)});
      offRampUpdates[2 * i + 1] = Router.OffRamp({
        sourceChainSelector: sourceChainSelector,
        offRamp: s_inboundNonceManager.getPreviousRamps(sourceChainSelector).prevOffRamp
      });
    }

    s_destRouter.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), offRampUpdates);
  }

  function test_executeSingleMessage_WithMessageTransformer() public {
    vm.stopPrank();
    vm.startPrank(address(s_messageTransformerOffRamp));
    Internal.Any2EVMRampMessage memory message = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    s_messageTransformerOffRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
  }

  function test_executeSingleMessage_WithMessageTransformer_RevertWhen_UnknownChain() public {
    vm.stopPrank();
    vm.startPrank(address(s_messageTransformerOffRamp));
    Internal.Any2EVMRampMessage memory message = _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    // Fail with any error (UnknownChain in this case) to check if OffRamp wraps the error with MessageTransformError during the revert
    s_inboundMessageTransformer.setShouldRevert(true);
    vm.expectRevert(
      abi.encodeWithSelector(
        IMessageTransformer.MessageTransformError.selector,
        abi.encodeWithSelector(MessageTransformerHelper.UnknownChain.selector)
      )
    );
    s_messageTransformerOffRamp.executeSingleMessage(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
  }
}