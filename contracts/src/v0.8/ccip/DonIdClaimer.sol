// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

interface ICapabilityRegistry {
    /// @notice Gets the next available DON ID from the CapabilitiesRegistry
    /// @return uint32 The next available DON ID
    function getNextDONId() external view returns (uint32);
}

/// @notice DonIDClaimer contract is used to claim and manage DON IDs. It interacts with
/// the CapabilitiesRegistry to fetch the next available DON ID and allows
/// for synchronization of the DON ID with an optional offset to skip certain
/// DON IDs as needed. The contract provides functionality for claiming, 
/// retrieving, and syncing DON IDs, ensuring that multiple workflows can
/// manage DON IDs without conflict or accidental reuse.
/// @dev The contract maintains its own internal counter for DON IDs and ensures
/// the next available ID is claimed and tracked by the contract. The sync function
/// allows for alignment with the CapabilitiesRegistry.
contract DonIDClaimer is Ownable2StepMsgSender {
  error ZeroAddressNotAllowed();

  string public constant override typeAndVersion = "DonIdClaimer 1.0.0-dev";
  /// @notice The next available DON ID that is claimed and incremented
  uint32 private s_nextDONId;  

  /// @notice The address of the CapabilitiesRegistry contract used to fetch the next DON ID
  address private immutable i_capabilitiesRegistry;

  /// @notice Initializes the contract with the CapabilitiesRegistry address
  /// @param capabilitiesRegistry The address of the CapabilitiesRegistry contract
  constructor(address capabilitiesRegistry) {
    if (capabilitiesRegistry == address(0)) {
      revert ZeroAddressNotAllowed();  
    }
    i_capabilitiesRegistry = capabilitiesRegistry;
    
    // Sync the initial s_nextDONId from the CapabilitiesRegistry contract
    s_nextDONId = ICapabilityRegistry(i_capabilitiesRegistry).getNextDONId();
  }

  /// @notice Claims the next available DON ID and increments the internal counter
  /// @dev The function increments s_nextDONId after returning the current value
  /// @return uint32 The DON ID that was claimed
  function claimNextDONId() external onlyOwner returns (uint32) {
    return s_nextDONId++;
  }

  /// @notice Returns the next available donID
  /// @return uint32 The next available donID to be claimed
  function getNextDONId() external view returns (uint32) {
      return s_nextDONId; 
  }

  /// @notice Synchronizes the next donID with the CapabilitiesRegistry and applies an offset
  /// @param offset The offset to adjust the donID (useful when certain DON IDs are dropped)
  /// @dev This can be used to synchronize with the CapabilitiesRegistry after some actions have occurred
  function syncNextDONIdWithOffset(uint32 offset) external onlyOwner {
      address capabilitiesRegistry = i_capabilitiesRegistry;
      s_nextDONId = ICapabilityRegistry(capabilitiesRegistry).getNextDONId() + offset;
  } 
}
