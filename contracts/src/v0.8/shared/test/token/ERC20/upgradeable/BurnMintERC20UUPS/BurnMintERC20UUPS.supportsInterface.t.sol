// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {
  IAccessControl,
  IBurnMintERC20Upgradeable,
  IERC165,
  IERC20,
  IGetCCIPAdmin
} from "../../../../../token/ERC20/upgradeable/BurnMintERC20UUPS.sol";
import {BurnMintERC20UUPSSetup} from "./BurnMintERC20UUPSSetup.t.sol";

contract BurnMintERC20UUPS_supportsInterface is BurnMintERC20UUPSSetup {
  function test_SupportsInterface() public view {
    assertTrue(s_burnMintERC20UUPS.supportsInterface(type(IERC20).interfaceId));
    assertTrue(s_burnMintERC20UUPS.supportsInterface(type(IBurnMintERC20Upgradeable).interfaceId));
    assertTrue(s_burnMintERC20UUPS.supportsInterface(type(IERC165).interfaceId));
    assertTrue(s_burnMintERC20UUPS.supportsInterface(type(IAccessControl).interfaceId));
    assertTrue(s_burnMintERC20UUPS.supportsInterface(type(IGetCCIPAdmin).interfaceId));
  }
}
