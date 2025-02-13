// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BurnMintERC20PausableFreezableTransparent} from "../../../../../token/ERC20/upgradeable/BurnMintERC20PausableFreezableTransparent.sol";
import {BurnMintERC20Transparent} from "../../../../../token/ERC20/upgradeable/BurnMintERC20Transparent.sol";
import {BaseTest} from "../../../../BaseTest.t.sol";

import {TransparentUpgradeableProxy} from "../../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

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

    address implementation = address(new BurnMintERC20PausableFreezableTransparent());

    s_TransparentProxy = address(
      new TransparentUpgradeableProxy(
        implementation,
        s_initialOwnerAddressForProxyAdmin,
        abi.encodeCall(
          BurnMintERC20Transparent.initialize,
          (s_name, s_symbol, s_decimals, s_maxSupply, s_preMint, s_defaultAdmin)
        )
      )
    );

    s_burnMintERC20PausableFreezableTransparent = BurnMintERC20PausableFreezableTransparent(s_TransparentProxy);

    s_burnMintERC20PausableFreezableTransparent.grantRole(
      s_burnMintERC20PausableFreezableTransparent.PAUSER_ROLE(),
      s_defaultPauser
    );
    s_burnMintERC20PausableFreezableTransparent.grantRole(
      s_burnMintERC20PausableFreezableTransparent.FREEZER_ROLE(),
      s_defaultFreezer
    );
  }
}
