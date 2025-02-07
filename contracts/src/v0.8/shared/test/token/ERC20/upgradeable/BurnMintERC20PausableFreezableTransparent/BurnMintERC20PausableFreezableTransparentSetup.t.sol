// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Upgrades} from "../../../../../../vendor/openzeppelin-foundry-upgrades/v0.3.8/Upgrades.sol";
import {BurnMintERC20PausableFreezableTransparent} from
  "../../../../../token/ERC20/upgradeable/BurnMintERC20PausableFreezableTransparent.sol";
import {BaseTest} from "../../../../BaseTest.t.sol";

contract BurnMintERC20PausableFreezableTransparentSetup is BaseTest {
  BurnMintERC20PausableFreezableTransparent internal s_burnMintERC20PausableFreezableTransparent;
  address s_TransparentProxy;

  string s_name = "CCIP-BnM Upgradeable";
  string s_symbol = "CCIP-BnM";
  uint8 s_decimals = 18;
  uint256 s_maxSupply = 1e27;
  uint256 s_preMint = 0;
  address s_defaultAdmin = OWNER;
  address s_defaultPauser = OWNER;
  address s_defaultFreezer = OWNER;
  address s_initialOwnerAddressForProxyAdmin = OWNER;

  function setUp() public virtual override {
    BaseTest.setUp();

    s_TransparentProxy = Upgrades.deployTransparentProxy(
      "BurnMintERC20PausableFreezableTransparent.sol",
      s_initialOwnerAddressForProxyAdmin,
      abi.encodeCall(
        BurnMintERC20PausableFreezableTransparent.initialize,
        (s_name, s_symbol, s_decimals, s_maxSupply, s_preMint, s_defaultAdmin, s_defaultPauser, s_defaultFreezer)
      )
    );

    s_burnMintERC20PausableFreezableTransparent = BurnMintERC20PausableFreezableTransparent(s_TransparentProxy);
  }
}
