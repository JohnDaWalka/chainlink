// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Ownable2Step} from "../../shared/access/Ownable2Step.sol";
import {DonIDClaimer} from "../DonIDClaimer.sol";
import {Test} from "forge-std/Test.sol";

contract MockCapabilitiesRegistry {
  uint32 private s_nextDonId;

  constructor(
    uint32 _initialDonId
  ) {
    s_nextDonId = _initialDonId;
  }

  function getNextDONId() external view returns (uint32) {
    return s_nextDonId;
  }
}

contract DonIDClaimerTest is Test {
  DonIDClaimer private s_donIDClaimer;
  MockCapabilitiesRegistry private s_mockRegistry;
  address private s_owner = address(0x1);
  address private s_deployer = address(0x2);
  address private s_unauthorized = address(0x3);

  function setUp() public {
    vm.startPrank(s_owner);
    s_mockRegistry = new MockCapabilitiesRegistry(100);
    s_donIDClaimer = new DonIDClaimer(address(s_mockRegistry));
    s_donIDClaimer.setAuthorizedDeployer(s_deployer, true);
    vm.stopPrank();
  }

  function test_Constructor() public {
    // Check the revert if constructor is called with a zero address
    vm.expectRevert(abi.encodeWithSelector(DonIDClaimer.ZeroAddressNotAllowed.selector));
    new DonIDClaimer(address(0));

    // Now test the normal constructor behavior with a valid address
    DonIDClaimer validDonIDClaimer = new DonIDClaimer(address(s_mockRegistry));
    assertEq(validDonIDClaimer.getNextDONId(), 100, "Initial DON ID should be set correctly from the registry");
  }

  function test_ClaimNextDONId() public {
    vm.expectEmit(true, true, true, true);
    emit DonIDClaimer.DonIDClaimed(s_deployer, 100);

    vm.prank(s_deployer);
    uint32 claimedId = s_donIDClaimer.claimNextDONId();
    assertEq(claimedId, 100, "Claimed DON ID should be 100");
    assertEq(s_donIDClaimer.getNextDONId(), 101, "Next DON ID should be incremented to 101");
  }

  function test_SyncNextDONIdWithOffset() public {
    vm.expectEmit(true, true, true, true);
    emit DonIDClaimer.DonIDSynced(110);

    vm.prank(s_deployer);
    s_donIDClaimer.syncNextDONIdWithOffset(10);
    assertEq(s_donIDClaimer.getNextDONId(), 110, "Next DON ID should be 110 after offset");
  }

  function test_SetAuthorizedDeployer() public {
    vm.expectEmit(true, true, true, true);
    emit DonIDClaimer.AuthorizedDeployerSet(s_unauthorized, true);

    vm.prank(s_owner);
    s_donIDClaimer.setAuthorizedDeployer(s_unauthorized, true);
    assertTrue(s_donIDClaimer.isAuthorizedDeployer(s_unauthorized), "Address should be authorized");
  }

  // Reverts
  function test_RevertWhen_UnauthorizedSenderClaimReverts() public {
    vm.expectRevert(abi.encodeWithSelector(DonIDClaimer.AccessForbidden.selector, s_unauthorized));
    vm.prank(s_unauthorized);
    s_donIDClaimer.claimNextDONId();
  }

  function test_RevertWhen_UnauthorizedSetAuthorizedDeployer() public {
    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    vm.prank(s_unauthorized);
    s_donIDClaimer.setAuthorizedDeployer(s_unauthorized, true);
  }
}
