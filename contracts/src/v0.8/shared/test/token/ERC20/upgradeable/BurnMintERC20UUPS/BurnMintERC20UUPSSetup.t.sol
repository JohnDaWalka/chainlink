// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {ERC1967Proxy} from "../../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {BurnMintERC20UUPS} from "../../../../../token/ERC20/upgradeable/BurnMintERC20UUPS.sol";
import {BaseTest} from "../../../../BaseTest.t.sol";

contract BurnMintERC20UUPSSetup is BaseTest {
  BurnMintERC20UUPS internal s_burnMintERC20UUPS;
  address s_uupsProxy;

  string s_name = "CCIP-BnM Upgradeable";
  string s_symbol = "CCIP-BnM";
  uint8 s_decimals = 18;
  uint256 s_maxSupply = 1e27;
  uint256 s_preMint = 0;
  address s_defaultAdmin = OWNER;
  address s_defaultUpgrader = OWNER;

  function setUp() public virtual override {
    BaseTest.setUp();

    address implementation = address(new BurnMintERC20UUPS());

    s_uupsProxy = address(
      new ERC1967Proxy(
        implementation,
        abi.encodeCall(
          BurnMintERC20UUPS.initialize,
          (s_name, s_symbol, s_decimals, s_maxSupply, s_preMint, s_defaultAdmin, s_defaultUpgrader)
        )
      )
    );

    s_burnMintERC20UUPS = BurnMintERC20UUPS(s_uupsProxy);
  }
}
