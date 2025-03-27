// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {DonIDClaimer} from "./DonIDClaimer.sol";
import {Test} from "forge-std/Test.sol";

interface ICapabilitiesRegistry {
    /// @notice Gets the next available DON ID from the CapabilitiesRegistry
    /// @return uint32 The next available DON ID
    function getNextDONId() external view returns (uint32);
}

contract MockCapabilitiesRegistry is ICapabilitiesRegistry {
    uint32 private nextDonId;

    constructor(uint32 _initialDonId) {
        nextDonId = _initialDonId;
    }

    function getNextDONId() external view override returns (uint32) {
        return nextDonId;
    }
}

contract DonIDClaimerTest is Test {
    DonIDClaimer private donIDClaimer; 
    MockCapabilitiesRegistry private mockRegistry;
    address private owner = address(0x1); 
    address private deployer = address(0x2);
    address private unauthorized = address(0x3); 

    function setUp() public {
        vm.startPrank(owner); 
        mockRegistry = new MockCapabilitiesRegistry(100);
        donIDClaimer= new DonIDClaimer(address(mockRegistry)); 
        donIDClaimer.setAuthorizedDeployer(deployer, true);
        vm.stopPrank();
    }

    function test_Constructor() public {
        // Check the revert if constructor is called with a zero address
        vm.expectRevert(abi.encodeWithSelector(DonIDClaimer.ZeroAddressNotAllowed.selector));
        new DonIDClaimer(address(0));  

        // Now test the normal constructor behavior with a valid address
        DonIDClaimer validDonIDClaimer = new DonIDClaimer(address(mockRegistry));
        assertEq(validDonIDClaimer.getNextDONId(), 100, "Initial DON ID should be set correctly from the registry");
    }


    function test_ClaimNextDONId() public {
        vm.expectEmit(true, true, true, true);
        emit DonIDClaimer.DonIDClaimed(deployer, 100); 

        vm.prank(deployer); 
        uint32 claimedId = donIDClaimer.claimNextDONId();
        assertEq(claimedId, 100, "Claimed DON ID should be 100"); 
        assertEq(donIDClaimer.getNextDONId(), 101, "Next DON ID should be incremented to 101");
    }
 
    function test_SyncNextDONIdWithOffset() public {
        vm.expectEmit(true, true, true, true);
        emit DonIDClaimer.DonIDSynced(110);   

        vm.prank(deployer); 
        donIDClaimer.syncNextDONIdWithOffset(10);
        assertEq(donIDClaimer.getNextDONId(), 110, "Next DON ID should be 110 after offset");
    }

    function test_SetAuthorizedDeployer() public {
        vm.expectEmit(true, true, true, true);
        emit DonIDClaimer.AuthorizedDeployerSet(unauthorized, true); 

        vm.prank(owner); 
        donIDClaimer.setAuthorizedDeployer(unauthorized, true);
        assertTrue(donIDClaimer.isAuthorizedDeployer(unauthorized), "Address should be authorized");
    }

    // Reverts 
    function test_RevertWhen_UnauthorizedSenderClaimReverts() public {
        vm.expectRevert(abi.encodeWithSelector(DonIDClaimer.AccessForbidden.selector, unauthorized));
        vm.prank(unauthorized); 
        donIDClaimer.claimNextDONId();
    } 

    function test_RevertWhen_UnauthorizedSetAuthorizedDeployer() public {
        vm.expectRevert("Ownable: caller is not the owner"); 
        vm.prank(unauthorized); 
        donIDClaimer.setAuthorizedDeployer(unauthorized, true);
    }
}