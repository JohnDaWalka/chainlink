// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {ERC1967Proxy} from
  "../../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {BurnMintERC20PausableUUPS} from "../../../../../token/ERC20/upgradeable/BurnMintERC20PausableUUPS.sol";
import {BaseTest} from "../../../../BaseTest.t.sol";

contract BurnMintERC20PausableUUPSSetup is BaseTest {
  BurnMintERC20PausableUUPS internal s_burnMintERC20PausableUUPS;
  address s_uupsProxy;

  string s_name = "CCIP-BnM Upgradeable";
  string s_symbol = "CCIP-BnM";
  uint8 s_decimals = 18;
  uint256 s_maxSupply = 1e27;
  uint256 s_preMint = 0;
  address s_defaultAdmin = OWNER;
  address s_defaultUpgrader = OWNER;
  address s_defaultPauser = OWNER;

  function setUp() public virtual override {
    BaseTest.setUp();

    address implementation = address(new BurnMintERC20PausableUUPS());

    s_uupsProxy = address(
      new ERC1967Proxy(
        implementation,
        abi.encodeCall(
          BurnMintERC20PausableUUPS.initialize,
          (s_name, s_symbol, s_decimals, s_maxSupply, s_preMint, s_defaultAdmin, s_defaultUpgrader, s_defaultPauser)
        )
      )
    );

    s_burnMintERC20PausableUUPS = BurnMintERC20PausableUUPS(s_uupsProxy);
  }
}
