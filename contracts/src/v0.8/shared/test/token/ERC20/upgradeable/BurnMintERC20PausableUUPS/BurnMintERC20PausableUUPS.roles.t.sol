// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BurnMintERC20PausableUUPS, IAccessControl} from "../../../../../token/ERC20/upgradeable/BurnMintERC20PausableUUPS.sol";
import {BurnMintERC20PausableUUPSSetup} from "./BurnMintERC20PausableUUPSSetup.t.sol";

contract BurnMintERC20PausableUUPS_roles is BurnMintERC20PausableUUPSSetup {
  function test_GrantMintAndBurnRoles() public {
    assertFalse(s_burnMintERC20PausableUUPS.hasRole(s_burnMintERC20PausableUUPS.MINTER_ROLE(), STRANGER));
    assertFalse(s_burnMintERC20PausableUUPS.hasRole(s_burnMintERC20PausableUUPS.BURNER_ROLE(), STRANGER));

    changePrank(s_defaultAdmin);

    vm.expectEmit();
    emit IAccessControl.RoleGranted(s_burnMintERC20PausableUUPS.MINTER_ROLE(), STRANGER, OWNER);
    vm.expectEmit();
    emit IAccessControl.RoleGranted(s_burnMintERC20PausableUUPS.BURNER_ROLE(), STRANGER, OWNER);

    s_burnMintERC20PausableUUPS.grantMintAndBurnRoles(STRANGER);

    assertTrue(s_burnMintERC20PausableUUPS.hasRole(s_burnMintERC20PausableUUPS.MINTER_ROLE(), STRANGER));
    assertTrue(s_burnMintERC20PausableUUPS.hasRole(s_burnMintERC20PausableUUPS.BURNER_ROLE(), STRANGER));
  }

  function test_GetCCIPAdmin() public view {
    assertEq(s_burnMintERC20PausableUUPS.getCCIPAdmin(), s_defaultAdmin);
  }

  function test_SetCCIPAdmin() public {
    changePrank(s_defaultAdmin);

    vm.expectEmit();
    emit BurnMintERC20PausableUUPS.CCIPAdminTransferred(s_defaultAdmin, STRANGER);

    s_burnMintERC20PausableUUPS.setCCIPAdmin(STRANGER);

    assertEq(s_burnMintERC20PausableUUPS.getCCIPAdmin(), STRANGER);
  }
}
