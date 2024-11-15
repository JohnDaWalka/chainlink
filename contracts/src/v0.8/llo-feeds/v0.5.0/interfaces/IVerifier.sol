// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {Common} from "../../libraries/Common.sol";
import {CommonV5} from "../libraries/CommonV5.sol";

interface IVerifier is IERC165 {
  /**
   * @notice sets off-chain reporting protocol configuration incl. participating oracles
   * @param signers addresses with which oracles sign the reports
   * @param f number of faulty oracles the system can tolerate
   * @param recipientAddressesAndWeights the addresses and weights of all the recipients to receive rewards
   */
  function setConfig(
    bytes32 configDigest,
    address[] memory signers,
    uint8 f,
    Common.AddressAndWeight[] memory recipientAddressesAndWeights
  ) external;

  /**
   * @notice Sets the fee manager address
   * @param feeManager The address of the fee manager
   */
  function setFeeManager(
    address feeManager
  ) external;

  /**
   * @notice Sets the access controller address
   * @param accessController The address of the access controller
   */
  function setAccessController(
    address accessController
  ) external;

  /**
   * @notice Updates the config active status
   * @param configDigest The ID of the config to update
   * @param isActive The new config active status
   */
  function setConfigActive(bytes32 configDigest, bool isActive) external;

  //TODO Nested config giving me trouble
  // /**
  //  * @notice Returns all DON configurations
  //  * @return array of DON configurations
  //  */
  // function getAllConfigs() external view returns (CommonV5.Config[] memory);
}
