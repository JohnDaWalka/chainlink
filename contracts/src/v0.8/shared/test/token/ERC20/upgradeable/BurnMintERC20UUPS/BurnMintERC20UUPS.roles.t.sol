// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BurnMintERC20UUPS, IAccessControl} from "../../../../../token/ERC20/upgradeable/BurnMintERC20UUPS.sol";
import {BurnMintERC20UUPSSetup} from "./BurnMintERC20UUPSSetup.t.sol";

contract BurnMintERC20UUPS_roles is BurnMintERC20UUPSSetup {
  function test_GrantMintAndBurnRoles() public {
    assertFalse(s_burnMintERC20UUPS.hasRole(s_burnMintERC20UUPS.MINTER_ROLE(), STRANGER));
    assertFalse(s_burnMintERC20UUPS.hasRole(s_burnMintERC20UUPS.BURNER_ROLE(), STRANGER));

    changePrank(s_defaultAdmin);

    vm.expectEmit();
    emit IAccessControl.RoleGranted(s_burnMintERC20UUPS.MINTER_ROLE(), STRANGER, s_defaultAdmin);
    vm.expectEmit();
    emit IAccessControl.RoleGranted(s_burnMintERC20UUPS.BURNER_ROLE(), STRANGER, s_defaultAdmin);

    s_burnMintERC20UUPS.grantMintAndBurnRoles(STRANGER);

    assertTrue(s_burnMintERC20UUPS.hasRole(s_burnMintERC20UUPS.MINTER_ROLE(), STRANGER));
    assertTrue(s_burnMintERC20UUPS.hasRole(s_burnMintERC20UUPS.BURNER_ROLE(), STRANGER));
  }

  function test_GetCCIPAdmin() public view {
    assertEq(s_burnMintERC20UUPS.getCCIPAdmin(), s_defaultAdmin);
  }

  function test_SetCCIPAdmin() public {
    changePrank(s_defaultAdmin);

    vm.expectEmit();
    emit BurnMintERC20UUPS.CCIPAdminTransferred(s_defaultAdmin, STRANGER);

    s_burnMintERC20UUPS.setCCIPAdmin(STRANGER);

    assertEq(s_burnMintERC20UUPS.getCCIPAdmin(), STRANGER);
  }
}
