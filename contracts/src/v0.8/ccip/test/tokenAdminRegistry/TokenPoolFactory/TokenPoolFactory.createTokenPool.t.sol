// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IOwner} from "../../../interfaces/IOwner.sol";

import {Ownable2Step} from "../../../../shared/access/Ownable2Step.sol";

import {Router} from "../../../Router.sol";
import {RateLimiter} from "../../../libraries/RateLimiter.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";
import {TokenPoolFactory} from "../../../tokenAdminRegistry/TokenPoolFactory/TokenPoolFactory.sol";
import {TokenPoolFactorySetup} from "./TokenPoolFactorySetup.t.sol";

import {BurnMintERC20} from "../../../../shared/token/ERC20/BurnMintERC20.sol";
import {Create2} from "../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/Create2.sol";

contract TokenPoolFactory_createTokenPool is TokenPoolFactorySetup {
  using Create2 for bytes32;

  uint8 private constant LOCAL_TOKEN_DECIMALS = 18;
  uint8 private constant REMOTE_TOKEN_DECIMALS = 6;

  address internal s_burnMintOffRamp = makeAddr("burn_mint_offRamp");

  function setUp() public override {
    TokenPoolFactorySetup.setUp();

    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](1);
    offRampUpdates[0] = Router.OffRamp({sourceChainSelector: DEST_CHAIN_SELECTOR, offRamp: s_burnMintOffRamp});
    s_sourceRouter.applyRampUpdates(new Router.OnRamp[](0), new Router.OffRamp[](0), offRampUpdates);
  }

  function test_createTokenPool_WithNoExistingToken_Success() public {
    vm.startPrank(OWNER);

    bytes32 dynamicSalt = keccak256(abi.encodePacked(FAKE_SALT, OWNER));

    address predictedTokenAddress =
      Create2.computeAddress(dynamicSalt, keccak256(s_tokenInitCode), address(s_tokenPoolFactory));

    // Create the constructor params for the predicted pool
    bytes memory poolCreationParams =
      abi.encode(predictedTokenAddress, LOCAL_TOKEN_DECIMALS, new address[](0), s_rmnProxy, s_sourceRouter);

    // Predict the address of the pool before we make the tx by using the init code and the params
    bytes memory predictedPoolInitCode = abi.encodePacked(s_poolInitCode, poolCreationParams);

    (address tokenAddress, address poolAddress) = s_tokenPoolFactory.deployTokenAndTokenPool(
      new TokenPoolFactory.RemoteTokenPoolInfo[](0), s_tokenInitCode, predictedPoolInitCode, FAKE_SALT
    );

    assertNotEq(address(0), tokenAddress, "Token Address should not be 0");
    assertNotEq(address(0), poolAddress, "Pool Address should not be 0");

    assertEq(predictedTokenAddress, tokenAddress, "Token Address should have been predicted");

    s_tokenAdminRegistry.acceptAdminRole(tokenAddress);
    Ownable2Step(tokenAddress).acceptOwnership();
    Ownable2Step(poolAddress).acceptOwnership();

    assertEq(poolAddress, s_tokenAdminRegistry.getPool(tokenAddress), "Token Pool should be set");
    assertEq(IOwner(tokenAddress).owner(), OWNER, "Token should be owned by the owner");
    assertEq(IOwner(poolAddress).owner(), OWNER, "Token should be owned by the owner");
  }

  function test_createTokenPool_WithExistingToken_Success() public {
    vm.startPrank(OWNER);

    address TOKEN_ADDRESS = address(new BurnMintERC20("FAKE TOKEN", "FAKE", 18, type(uint256).max, 0));

    bytes memory RANDOM_TOKEN_ADDRESS = abi.encode(makeAddr("RANDOM_TOKEN_ADDRESS"));
    bytes memory RANDOM_POOL_ADDRESS = abi.encode(makeAddr("RANDOM_POOL"));

    // Create an array of remote pools with some fake addresses
    TokenPoolFactory.RemoteTokenPoolInfo[] memory remoteTokenPools = new TokenPoolFactory.RemoteTokenPoolInfo[](1);

    bytes[] memory remotePools = new bytes[](1);
    remotePools[0] = RANDOM_POOL_ADDRESS;

    remoteTokenPools[0] = TokenPoolFactory.RemoteTokenPoolInfo(
      DEST_CHAIN_SELECTOR, // remoteChainSelector
      remotePools, // remotePoolAddress
      RANDOM_TOKEN_ADDRESS, // remoteTokenAddress
      RateLimiter.Config(false, 0, 0) // rateLimiterConfig
    );

    // Create the constructor params for the predicted pool
    bytes memory poolCreationParams =
      abi.encode(TOKEN_ADDRESS, LOCAL_TOKEN_DECIMALS, new address[](0), s_rmnProxy, s_sourceRouter);

    // Predict the address of the pool before we make the tx by using the init code and the params
    bytes memory poolInitcode = abi.encodePacked(s_poolInitCode, poolCreationParams);

    address poolAddress = s_tokenPoolFactory.deployTokenPoolWithExistingToken(remoteTokenPools, poolInitcode, FAKE_SALT);

    assertNotEq(address(0), poolAddress, "Pool Address should not be 0");

    Ownable2Step(poolAddress).acceptOwnership();

    assertEq(address(TokenPool(poolAddress).getToken()), TOKEN_ADDRESS, "local tToken address should have been set");

    assertEq(
      TokenPool(poolAddress).getRemotePools(DEST_CHAIN_SELECTOR)[0],
      RANDOM_POOL_ADDRESS,
      "Remote Pool Address should have been set"
    );

    assertEq(IOwner(poolAddress).owner(), OWNER, "Pool should be owned by the owner");
  }
}
