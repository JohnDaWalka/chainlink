// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {OnRampSetup} from "./OnRampSetup.t.sol";
import {OnRamp} from "../../../onRamp/OnRamp.sol";
import {MessageTransformerOnRamp} from "../../../onRamp/MessageTransformerOnRamp.sol";
import {MessageTransformerHelper} from "../../helpers/MessageTransformerHelper.sol";
import {Client} from "../../../libraries/Client.sol";
import {Internal} from "../../../libraries/Internal.sol";
import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC20.sol";
import {IMessageTransformer} from "../../../interfaces/IMessageTransformer.sol";
import {Router} from "../../../Router.sol";
import {NonceManager} from "../../../NonceManager.sol";
import {AuthorizedCallers} from "../../../../shared/access/AuthorizedCallers.sol";

contract MessageTransformerOnRamp_forwardFromRouter is OnRampSetup {

    MessageTransformerOnRamp internal s_messageTransformerOnRamp;
    MessageTransformerHelper internal s_messageTransformer;

    function setUp() public virtual override {
        super.setUp();
        s_messageTransformer = new MessageTransformerHelper();
        s_messageTransformerOnRamp = new MessageTransformerOnRamp(
            s_onRamp.getStaticConfig(),
            s_onRamp.getDynamicConfig(),
            _generateDestChainConfigArgs(s_sourceRouter),
            address(s_messageTransformer)
        );
        s_metadataHash = keccak256(abi.encode(Internal.EVM_2_ANY_MESSAGE_HASH, SOURCE_CHAIN_SELECTOR, DEST_CHAIN_SELECTOR, address(s_messageTransformerOnRamp)));
        address[] memory authorizedCallers = new address[](1);
        authorizedCallers[0] = address(s_messageTransformerOnRamp);

        NonceManager(s_outboundNonceManager).applyAuthorizedCallerUpdates(
            AuthorizedCallers.AuthorizedCallerArgs({addedCallers: authorizedCallers, removedCallers: new address[](0)})
        );

        Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](1);
        onRampUpdates[0] = Router.OnRamp({destChainSelector: DEST_CHAIN_SELECTOR, onRamp: address(s_messageTransformerOnRamp)});

        Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](2);
        offRampUpdates[0] = Router.OffRamp({sourceChainSelector: SOURCE_CHAIN_SELECTOR, offRamp: makeAddr("offRamp0")});
        offRampUpdates[1] = Router.OffRamp({sourceChainSelector: SOURCE_CHAIN_SELECTOR, offRamp: makeAddr("offRamp1")});
        s_sourceRouter.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), offRampUpdates);
        vm.startPrank(address(s_sourceRouter));
    }

    function test_forwardFromRouter_WithMessageTransformer_Success() public {
        Client.EVM2AnyMessage memory message = _generateEmptyMessage();
        message.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT * 2}));
        uint256 feeAmount = 1234567890;
        message.tokenAmounts = new Client.EVMTokenAmount[](1);
        message.tokenAmounts[0].amount = 1e18;
        message.tokenAmounts[0].token = s_sourceTokens[0];
        IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_messageTransformerOnRamp), feeAmount);
        vm.expectEmit();
        emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, 1, _messageToEvent(message, 1, 1, feeAmount, OWNER));
        s_messageTransformerOnRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
    }

    function test_MessageTransformerError_Revert() public {
        Client.EVM2AnyMessage memory message = _generateEmptyMessage();
        message.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT * 2}));
        uint256 feeAmount = 1234567890;
        message.tokenAmounts = new Client.EVMTokenAmount[](1);
        message.tokenAmounts[0].amount = 1e18;
        message.tokenAmounts[0].token = s_sourceTokens[0];
        IERC20(s_sourceFeeToken).transferFrom(OWNER, address(s_messageTransformerOnRamp), feeAmount);

        // Fail with any error (UnknownChain in this case) to check if OnRamp wraps the error with MessageTransformError during the revert
        s_messageTransformer.setShouldRevert(true);
        vm.expectRevert(
            abi.encodeWithSelector(
                IMessageTransformer.MessageTransformError.selector,
                abi.encodeWithSelector(MessageTransformerHelper.UnknownChain.selector)
            )
        );
        s_messageTransformerOnRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
    }
}