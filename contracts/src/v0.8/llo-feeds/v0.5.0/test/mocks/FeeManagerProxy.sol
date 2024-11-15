// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {IVerifierFeeManager} from "../../interfaces/IVerifierFeeManager.sol";

contract FeeManagerProxy {
  IVerifierFeeManager internal s_feeManager;

  function processFee(bytes32 poolId, bytes calldata payload, bytes calldata parameterPayload) public payable {
    s_feeManager.processFee{value: msg.value}(poolId, payload, parameterPayload, msg.sender);
  }

  function processFeeBulk(
    bytes32[] memory poolIds,
    bytes[] calldata payloads,
    bytes calldata parameterPayload
  ) public payable {
    s_feeManager.processFeeBulk{value: msg.value}(poolIds, payloads, parameterPayload, msg.sender);
  }

  function setFeeManager(address feeManager) public {
    s_feeManager = IVerifierFeeManager(feeManager);
  }
}
