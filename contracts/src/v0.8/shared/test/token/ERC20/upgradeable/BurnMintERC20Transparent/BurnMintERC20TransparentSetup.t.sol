// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {TransparentUpgradeableProxy} from
  "../../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import {BurnMintERC20Transparent} from "../../../../../token/ERC20/upgradeable/BurnMintERC20Transparent.sol";
import {BaseTest} from "../../../../BaseTest.t.sol";

contract BurnMintERC20TransparentSetup is BaseTest {
  BurnMintERC20Transparent internal s_burnMintERC20Transparent;
  address s_TransparentProxy;

  string s_name = "CCIP-BnM Upgradeable";
  string s_symbol = "CCIP-BnM";
  uint8 s_decimals = 18;
  uint256 s_maxSupply = 1e27;
  uint256 s_preMint = 0;
  address s_defaultAdmin = OWNER;
  address s_initialOwnerAddressForProxyAdmin = OWNER;

  function setUp() public virtual override {
    BaseTest.setUp();

    address implementation = address(new BurnMintERC20Transparent());

    s_TransparentProxy = address(
      new TransparentUpgradeableProxy(
        implementation,
        s_initialOwnerAddressForProxyAdmin,
        abi.encodeCall(
          BurnMintERC20Transparent.initialize, (s_name, s_symbol, s_decimals, s_maxSupply, s_preMint, s_defaultAdmin)
        )
      )
    );

    s_burnMintERC20Transparent = BurnMintERC20Transparent(s_TransparentProxy);
  }
}
