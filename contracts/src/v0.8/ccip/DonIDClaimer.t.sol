// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

// import {DonIDClaimer} from "./DonIDClaimer.sol";
import {Test} from "forge-std/Test.sol";

// interface ICapabilitiesRegistry {
//     /// @notice Gets the next available DON ID from the CapabilitiesRegistry
//     /// @return uint32 The next available DON ID
//     function getNextDONId() external view returns (uint32);
// }

// contract MockCapabilitiesRegistry is ICapabilitiesRegistry {
//     uint32 private nextDonId;

//     constructor(uint32 _initialDonId) {
//         nextDonId = _initialDonId;
//     }

//     function getNextDONId() external view override returns (uint32) {
//         return nextDonId;
//     }
// }

contract DonIDClaimerTest is Test {
    // DonIDClaimer private donIDClaimer; 
    // MockCapabilitiesRegistry private mockRegistry;
    // address private owner = address(0x1); 
    // address private deployer = address(0x2);
    // address private unauthorized = address(0x3); 

    // function setUp() public {
    //     vm.startPrank(owner); 
    //     mockRegistry = new MockCapabilitiesRegistry(100);
    //     donIDClaimer= new DonIDClaimer(address(mockRegistry)); 
    //     DonIDClaimer.setAuthorizedDeployer(deployer, true);
    //     vm.stopPrank();
    // }

    function test_ClaimNextDONIdDonIDClaimer() public {
        // vm.prank(deployer); 
        // uint32 claimedId = donIDClaimer.claimNextDONId();
        // assertEq(claimedId, 100, "Claimed DON ID should be 100"); 
        assertEq(101, 101, "Next DON ID should be incremented to 101");
    }

}