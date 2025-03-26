// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {ITypeAndVersion} from "../shared/interfaces/ITypeAndVersion.sol";

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
contract DonIDClaimer is ITypeAndVersion {
  error ZeroAddressNotAllowed();
  error Unauthorized(address caller);

  string public constant override typeAndVersion = "DonIdClaimer 1.0.0-dev";
  /// @notice The next available DON ID that is claimed and incremented
  uint32 private s_nextDONId;  

  /// @notice The address of the CapabilitiesRegistry contract used to fetch the next DON ID
  address private immutable i_capabilitiesRegistry;

  /// @notice Mapping to track authorized deployed keys 
  mapping(address => bool) private authorizedDeployer; 

  /// @notice Initializes the contract with the CapabilitiesRegistry address
  /// @param capabilitiesRegistry The address of the CapabilitiesRegistry contract
  constructor(address capabilitiesRegistry) {
    if (capabilitiesRegistry == address(0)) revert ZeroAddressNotAllowed();  
    i_capabilitiesRegistry = capabilitiesRegistry;

    // Initializing the deployer authorization (owner can be the initial deployer)
    authorizedDeployer[msg.sender] = true;

    // Sync the initial s_nextDONId from the CapabilitiesRegistry contract
    s_nextDONId = ICapabilityRegistry(i_capabilitiesRegistry).getNextDONId();
  }

  /// @notice Modifier to check if the caller is an authorized deployer
  modifier onlyAuthorizedDeployer() {
    if (!authorizedDeployer[msg.sender]) {
      revert Unauthorized(msg.sender); 
    }
    _;
  }

  /// @notice Claims the next available DON ID and increments the internal counter
  /// @dev The function increments s_nextDONId after returning the current value
  /// @return uint32 The DON ID that was claimed
  function claimNextDONId() external onlyAuthorizedDeployer returns (uint32) {
    return s_nextDONId++;
  }

  /// @notice Synchronizes the next donID with the CapabilitiesRegistry and applies an offset
  /// @param offset The offset to adjust the donID (useful when certain DON IDs are dropped)
  /// @dev This can be used to synchronize with the CapabilitiesRegistry after some actions have occurred
  function syncNextDONIdWithOffset(uint32 offset) external onlyAuthorizedDeployer {
      address capabilitiesRegistry = i_capabilitiesRegistry;
      s_nextDONId = ICapabilityRegistry(capabilitiesRegistry).getNextDONId() + offset;
  } 

  /// @notice Sets authorization status for a deployer address
  /// @param senderAddress The address to be added or removed as an authorized deployer
  /// @param allowed Boolean indicating whether the address is authorized (true) or revoked (false)
  /// @dev Can only be called by an existing authorized deployer
  function setAuthorizedDeployer(address senderAddress, bool allowed) external onlyAuthorizedDeployer {
    if (senderAddress == address(0)) revert ZeroAddressNotAllowed();  
    authorizedDeployer[senderAddress] = allowed;
  }

  /// @notice Returns the next available donID
  /// @return uint32 The next available donID to be claimed
  function getNextDONId() external view returns (uint32) {
      return s_nextDONId; 
  }

  /// @notice Checks if an address is an authorized deployer
  /// @param senderAddress The address to check for authorization
  /// @return bool True if the address is an authorized deployer, false otherwise
  function isAuthorizedDeployer(address senderAddress) external view returns (bool) {
    return authorizedDeployer[senderAddress];
  }
}
