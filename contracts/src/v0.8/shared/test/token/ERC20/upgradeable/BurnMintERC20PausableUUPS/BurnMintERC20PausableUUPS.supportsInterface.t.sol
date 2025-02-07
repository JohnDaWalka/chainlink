// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IAccessControl, IBurnMintERC20Upgradeable, IERC165, IERC20, IGetCCIPAdmin} from "../../../../../token/ERC20/upgradeable/BurnMintERC20PausableUUPS.sol";
import {BurnMintERC20PausableUUPSSetup} from "./BurnMintERC20PausableUUPSSetup.t.sol";

contract BurnMintERC20PausableUUPS_supportsInterface is BurnMintERC20PausableUUPSSetup {
  function test_SupportsInterface() public view {
    assertTrue(s_burnMintERC20PausableUUPS.supportsInterface(type(IERC20).interfaceId));
    assertTrue(s_burnMintERC20PausableUUPS.supportsInterface(type(IBurnMintERC20Upgradeable).interfaceId));
    assertTrue(s_burnMintERC20PausableUUPS.supportsInterface(type(IERC165).interfaceId));
    assertTrue(s_burnMintERC20PausableUUPS.supportsInterface(type(IAccessControl).interfaceId));
    assertTrue(s_burnMintERC20PausableUUPS.supportsInterface(type(IGetCCIPAdmin).interfaceId));
  }
}
