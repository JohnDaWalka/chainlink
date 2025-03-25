// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {ICapabilityConfiguration} from "../keystone/interfaces/ICapabilityConfiguration.sol";

contract DonIDClaimer {
  error ZeroAddressNotAllowed();

  uint256 public donID;  
  address internal immutable i_capabilitiesRegistry;

  constructor(
    address capabilitiesRegistry
  ) {
    if (capabilitiesRegistry == address(0)) {
      revert ZeroAddressNotAllowed();
    }
    i_capabilitiesRegistry = capabilitiesRegistry;
  }

  function claimNextDonID() external returns (uint256) {
    donID = donID + 1;
    return donID; 
  }

  function setDONID(uint256 donId) external pure returns (uint256) {
      return donId; 
  }

  function getDonID() external view returns (uint256) {
      return donID; 
  }
  
  function syncDonIdWithCapReg() external {
      address capabilitiesRegistry = i_capabilitiesRegistry;
      uint32 id = ICapabilityConfiguration(capabilitiesRegistry).getNextDONId();
      donID = id; 
  } 
}
 