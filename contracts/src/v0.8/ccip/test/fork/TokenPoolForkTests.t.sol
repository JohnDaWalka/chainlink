// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IRouter} from "../../interfaces/IRouter.sol";
import {IMessageTransmitter} from "../../pools/USDC/IMessageTransmitter.sol";
import {ITokenMessenger} from "../../pools/USDC/ITokenMessenger.sol";

import {Pool} from "../../libraries/Pool.sol";
import {RateLimiter} from "../../libraries/RateLimiter.sol";
import {LockReleaseTokenPool} from "../../pools/LockReleaseTokenPool.sol";
import {SiloedLockReleaseTokenPool} from "../../pools/SiloedLockReleaseTokenPool.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {HybridLockReleaseUSDCTokenPool} from "../../pools/USDC/HybridLockReleaseUSDCTokenPool.sol";
import {LOCK_RELEASE_FLAG} from "../../pools/USDC/HybridLockReleaseUSDCTokenPool.sol";
import {TokenAdminRegistry} from "../../tokenAdminRegistry/TokenAdminRegistry.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC20.sol";
import {SafeERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";
import {Test} from "forge-std/Test.sol";

contract RoninForkTests is Test {
  using SafeERC20 for IERC20;

  TokenAdminRegistry public constant TOKEN_ADMIN_REGISTRY =
    TokenAdminRegistry(0xb22764f98dD05c789929716D677382Df22C05Cb6);

  IRouter public constant ROUTER = IRouter(0x80226fc0Ee2b096224EeAc085Bb9a8cba1146f7D);
  address public constant RMNPRoxy = 0x411dE17f12D1A34ecC7F45f49844626267c75e81;

  HybridLockReleaseUSDCTokenPool public constant USDCTokenPool =
    HybridLockReleaseUSDCTokenPool(0xc2e3A3C18ccb634622B57fF119a1C8C7f12e8C0c);

  address public constant USDC = 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48;
  address public constant WETH = 0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2;
  address public constant RONIN_OFF_RAMP = 0x9a3Ed7007809CfD666999e439076B4Ce4120528D;
  address public constant LEGACY_WETH_POOL = 0x69c24c970B65e22Ac26864aF10b2295B7d78f93A;

  address public constant ARB_OFF_RAMP = 0xdf615eF8D4C64d0ED8Fd7824BBEd2f6a10245aC9;
  bytes public constant ARB_USDC_POOL = hex"000000000000000000000000f46beff26e1c4552fb4ffb00314bdf175fbe97e4";

  uint64 public constant RONIN_CHAIN_SELECTOR = 6916147374840168594;
  uint64 public constant ARB_CHAIN_SELECTOR = 4949039107694359620;
  uint256 public constant FORK_BLOCK = 21680971;

  address public constant MCMS_MULTISIG = 0x44835bBBA9D40DEDa9b64858095EcFB2693c9449;
  address public constant RONIN_USDC_LIQUIDITY_PROVIDER = address(0x1234);

  SiloedLockReleaseTokenPool public s_siloedTokenPool;

  address public constant SILOED_REBALANCER = address(0xdeadbeef);

  function setUp() public {
    vm.createSelectFork(vm.envString("ETHEREUM_RPC_URL"), FORK_BLOCK);

    vm.mockCall(
      address(USDCTokenPool.i_messageTransmitter()),
      abi.encodeWithSelector(IMessageTransmitter.receiveMessage.selector),
      abi.encode(false)
    );

    s_siloedTokenPool = new SiloedLockReleaseTokenPool(IERC20(WETH), 18, new address[](0), RMNPRoxy, address(ROUTER));
    vm.makePersistent(address(s_siloedTokenPool));
  }

  function test_LockReleaseTokenPool_Migrations() public {
    vm.createSelectFork(vm.envString("ARBITRUM_RPC_URL"), 301236399);

    address ArbRouter = 0x141fa059441E0ca23ce184B6A78bafD2A517DdE8;
    address ArbTokenAdminRegistry = 0x39AE1032cF4B334a1Ed41cdD0833bdD7c7E7751E;
    address ARBWeth = 0x82aF49447D8a07e3bd95BD0d56f35241523fBab1;
    address ArbRMN = 0xC311a21e6fEf769344EB1515588B9d535662a145;

    address currentWethPool = TokenAdminRegistry(ArbTokenAdminRegistry).getPool(ARBWeth);

    LockReleaseTokenPool arbWethPool = LockReleaseTokenPool(currentWethPool);
    LockReleaseTokenPool newArbWethPool = new LockReleaseTokenPool(IERC20(ARBWeth), 18, new address[](0), ArbRMN, true, ArbRouter);

    vm.makePersistent(address(arbWethPool));
    vm.makePersistent(address(newArbWethPool));

    address currentOwner = arbWethPool.owner();

    vm.startPrank(currentOwner);
    
    arbWethPool.setRebalancer(address(newArbWethPool));

    uint256 liquidityBalance = IERC20(ARBWeth).balanceOf(address(arbWethPool));

    vm.stopPrank();
    newArbWethPool.transferLiquidity(address(arbWethPool), liquidityBalance);

    assertEq(IERC20(ARBWeth).balanceOf(address(newArbWethPool)), liquidityBalance);
    assertEq(IERC20(ARBWeth).balanceOf(address(arbWethPool)), 0);

    TokenAdminRegistry.TokenConfig memory config = TokenAdminRegistry(ArbTokenAdminRegistry).getTokenConfig(ARBWeth);

    vm.startPrank(config.administrator);
    TokenAdminRegistry(ArbTokenAdminRegistry).setPool(ARBWeth, address(newArbWethPool));

    assertEq(TokenAdminRegistry(ArbTokenAdminRegistry).getPool(ARBWeth), address(newArbWethPool));

  }

  function test_SiloedLockReleaseTokenPool() public {
    address currentWethPool = TOKEN_ADMIN_REGISTRY.getPool(WETH);

    // Get the address of the rebalancer that can withdraw from the pool
    address rebalancer = LockReleaseTokenPool(LEGACY_WETH_POOL).getRebalancer();

    // Set the rebalancer on the new pool to be equal to the current finance multisig
    s_siloedTokenPool.setRebalancer(rebalancer);

    // Add ronin to the list of allowed chains
    RateLimiter.Config memory rateLimiterConfig = RateLimiter.Config({isEnabled: false, capacity: 0, rate: 0});

    bytes[] memory poolAddresses = new bytes[](1);
    poolAddresses[0] = abi.encode("FAKE_POOL");

    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: RONIN_CHAIN_SELECTOR,
      remotePoolAddresses: poolAddresses,
      remoteTokenAddress: "FAKE_TOKEN",
      outboundRateLimiterConfig: rateLimiterConfig,
      inboundRateLimiterConfig: rateLimiterConfig
    });

    s_siloedTokenPool.applyChainUpdates(new uint64[](0), chainUpdates);

    vm.startPrank(MCMS_MULTISIG);

    // Set the new siloed lock release token pool as the official pool for WETH on mainnet.
    vm.expectEmit();
    emit TokenAdminRegistry.PoolSet(WETH, currentWethPool, address(s_siloedTokenPool));
    TOKEN_ADMIN_REGISTRY.setPool(WETH, address(s_siloedTokenPool));

    assertEq(address(s_siloedTokenPool), TOKEN_ADMIN_REGISTRY.getPool(WETH));

    // Get the balance of the legacy pool now which will be migrated
    uint256 liquidityBalance = IERC20(WETH).balanceOf(address(LEGACY_WETH_POOL));

    // Do the liquidity transfer from the existing pool to the new one with unsiloed liquidity
    vm.startPrank(rebalancer);

    LockReleaseTokenPool(LEGACY_WETH_POOL).withdrawLiquidity(liquidityBalance);
    IERC20(WETH).safeApprove(address(s_siloedTokenPool), type(uint256).max);

    s_siloedTokenPool.provideLiquidity(liquidityBalance);

    // Check that the liquidity was migrated correctly.
    assertEq(IERC20(WETH).balanceOf(address(s_siloedTokenPool)), liquidityBalance);
    assertEq(s_siloedTokenPool.getUnsiloedLiquidity(), liquidityBalance);

    // Set the Remote chain to be siloed
    SiloedLockReleaseTokenPool.SiloConfigUpdate[] memory updates = new SiloedLockReleaseTokenPool.SiloConfigUpdate[](1);
    updates[0] = SiloedLockReleaseTokenPool.SiloConfigUpdate({
      remoteChainSelector: RONIN_CHAIN_SELECTOR,
      rebalancer: SILOED_REBALANCER
    });

    vm.stopPrank();
    s_siloedTokenPool.updateSiloDesignations(new uint64[](0), updates);
    assertTrue(s_siloedTokenPool.isSiloed(RONIN_CHAIN_SELECTOR));

    // As the Ronin rebalance role, approve the siloed token pool to spend the WETH
    vm.startPrank(SILOED_REBALANCER);
    IERC20(WETH).safeApprove(address(s_siloedTokenPool), type(uint256).max);

    // give 1e24 worth of WETH to the silo rebalancer for testing
    uint256 dealAmount = 1e24;
    deal(WETH, SILOED_REBALANCER, dealAmount);

    // Provide the siloed liquidity for Ronin
    s_siloedTokenPool.provideSiloedLiquidity(RONIN_CHAIN_SELECTOR, dealAmount);

    // Check that it was properly provided
    assertEq(IERC20(WETH).balanceOf(address(s_siloedTokenPool)), dealAmount + liquidityBalance);
    assertEq(s_siloedTokenPool.getAvailableTokens(RONIN_CHAIN_SELECTOR), dealAmount);
    assertEq(s_siloedTokenPool.getUnsiloedLiquidity(), liquidityBalance);

    // Attempt a release of the siloed liquidity
    vm.startPrank(RONIN_OFF_RAMP);

    vm.expectEmit();
    emit TokenPool.Released(RONIN_OFF_RAMP, address(0xdeadbeef), dealAmount);

    s_siloedTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: abi.encode("FAKE_SENDER"),
        remoteChainSelector: RONIN_CHAIN_SELECTOR,
        receiver: address(0xdeadbeef),
        amount: dealAmount,
        localToken: WETH,
        sourcePoolAddress: abi.encode("FAKE_POOL"),
        sourcePoolData: "",
        offchainTokenData: ""
      })
    );

    // Check that the tokens were actually released
    assertEq(IERC20(WETH).balanceOf(address(s_siloedTokenPool)), liquidityBalance);
    assertEq(s_siloedTokenPool.getAvailableTokens(RONIN_CHAIN_SELECTOR), 0);
    assertEq(s_siloedTokenPool.getUnsiloedLiquidity(), liquidityBalance);
  }

  function test_HybridLockReleaseUSDCTokenPool() public {
    vm.startPrank(MCMS_MULTISIG);
    uint256 sendAmount = 1e6;

    // Check that Ronin is disabled by default for the pool
    assertFalse(USDCTokenPool.shouldUseLockRelease(RONIN_CHAIN_SELECTOR));
    assertFalse(USDCTokenPool.isSupportedChain(RONIN_CHAIN_SELECTOR));
    RateLimiter.Config memory rateLimiterConfig = RateLimiter.Config({isEnabled: false, capacity: 0, rate: 0});

    bytes[] memory poolAddresses = new bytes[](1);
    poolAddresses[0] = abi.encode("FAKE_POOL");

    TokenPool.ChainUpdate[] memory updates = new TokenPool.ChainUpdate[](1);
    updates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: RONIN_CHAIN_SELECTOR,
      remotePoolAddresses: poolAddresses,
      remoteTokenAddress: "FAKE_TOKEN",
      outboundRateLimiterConfig: rateLimiterConfig,
      inboundRateLimiterConfig: rateLimiterConfig
    });

    // Enable the pool first and then enable the selector mechanism to use L/R
    USDCTokenPool.applyChainUpdates(new uint64[](0), updates);

    // Check that Ronin has been enabled.
    assertTrue(USDCTokenPool.isSupportedChain(RONIN_CHAIN_SELECTOR));

    uint64[] memory chainSelectors = new uint64[](1);
    chainSelectors[0] = RONIN_CHAIN_SELECTOR;

    // Update the mechanism to indicate use L/R for Ronin
    USDCTokenPool.updateChainSelectorMechanisms(new uint64[](0), chainSelectors);

    assertTrue(USDCTokenPool.shouldUseLockRelease(RONIN_CHAIN_SELECTOR));

    // Get the on-ramp address and impersonate to send message to the token pool
    address roninOnRamp = ROUTER.getOnRamp(RONIN_CHAIN_SELECTOR);
    deal(USDC, address(USDCTokenPool), sendAmount);

    vm.startPrank(roninOnRamp);

    vm.expectEmit();
    emit TokenPool.Locked(roninOnRamp, sendAmount);

    // Lock the tokens using the L/R Mechanism
    USDCTokenPool.lockOrBurn(
      Pool.LockOrBurnInV1({
        receiver: abi.encode(address(0xdeadbeef)),
        remoteChainSelector: RONIN_CHAIN_SELECTOR,
        originalSender: address(0xdeadbeef),
        amount: sendAmount,
        localToken: USDC
      })
    );

    // Assert that the correct path was followed
    assertEq(IERC20(USDC).balanceOf(address(USDCTokenPool)), sendAmount);
    assertEq(USDCTokenPool.getLockedTokensForChain(RONIN_CHAIN_SELECTOR), sendAmount);

    // Deal overwrites balance, not adding more so we deal twice the amount
    // to prevent overwriting the tokens already locked in the previous L/R call.
    deal(USDC, address(USDCTokenPool), sendAmount * 2);

    // Change the On-Ramp to ARB
    address arbOnRamp = ROUTER.getOnRamp(ARB_CHAIN_SELECTOR);
    vm.startPrank(arbOnRamp);

    vm.expectEmit(false, true, true, false);
    emit ITokenMessenger.DepositForBurn(
      0, USDC, sendAmount, address(USDCTokenPool), bytes32(0), 0, bytes32(0), bytes32(0)
    );

    // Burn the tokens using the B/M CCTP Mechanism
    USDCTokenPool.lockOrBurn(
      Pool.LockOrBurnInV1({
        receiver: abi.encode(address(0xdeadbeef)),
        remoteChainSelector: ARB_CHAIN_SELECTOR,
        originalSender: address(0xdeadbeef),
        amount: sendAmount,
        localToken: USDC
      })
    );

    // Check that the tokens were actually burned
    assertEq(
      IERC20(USDC).balanceOf(address(USDCTokenPool)),
      sendAmount,
      "Token Balance should be sendAmount for locked tokens on Ronin"
    );

    // Check that no internal accounting updates occured for CCTP
    assertEq(USDCTokenPool.getLockedTokensForChain(ARB_CHAIN_SELECTOR), 0);

    // Mock the call to the CCTP transmitter so that it's true
    vm.startPrank(ARB_OFF_RAMP);

    // The call should will revert because the parsing of the attestation data will fail since I have not provided
    // any since this is a live-fork. However, it is only really necessary to check that the B/M path was taken since the source pool data is different from the BurnMintWithLockReleaseFlag pool.
    vm.expectRevert();

    USDCTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: abi.encode("FAKE_SENDER"),
        remoteChainSelector: ARB_CHAIN_SELECTOR,
        receiver: address(0xdeadbeef),
        amount: sendAmount,
        localToken: USDC,
        sourcePoolAddress: ARB_USDC_POOL,
        sourcePoolData: "FAKE_DATA",
        offchainTokenData: "OFFCHAIN_TOKEN_DATA"
      })
    );

    vm.startPrank(RONIN_OFF_RAMP);

    assertTrue(USDCTokenPool.shouldUseLockRelease(RONIN_CHAIN_SELECTOR));
    assertEq(USDCTokenPool.getLockedTokensForChain(RONIN_CHAIN_SELECTOR), sendAmount);

    vm.expectEmit();
    emit TokenPool.Released(RONIN_OFF_RAMP, address(0xdeadbeef), sendAmount);

    // Use the LOCK_RELEASE_FLAG to trigger the L/R mechanism on the pool from Ronin
    USDCTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: abi.encode("FAKE_SENDER"),
        remoteChainSelector: RONIN_CHAIN_SELECTOR,
        receiver: address(0xdeadbeef),
        amount: sendAmount,
        localToken: USDC,
        sourcePoolAddress: abi.encode("FAKE_POOL"),
        sourcePoolData: abi.encode(LOCK_RELEASE_FLAG),
        offchainTokenData: ""
      })
    );

    // Check that the tokens were actually released
    assertEq(IERC20(USDC).balanceOf(address(USDCTokenPool)), 0);

    // Test the providing of USDC Liquidity
    vm.startPrank(MCMS_MULTISIG);
    
    // Set the Ronin chain as a liquidity provider
    uint256 liquidityAmount = 1e12;
    USDCTokenPool.setLiquidityProvider(RONIN_CHAIN_SELECTOR, RONIN_USDC_LIQUIDITY_PROVIDER);
    deal(USDC, RONIN_USDC_LIQUIDITY_PROVIDER, liquidityAmount);
    
    // Provide the liquidity
    vm.startPrank(RONIN_USDC_LIQUIDITY_PROVIDER);
    IERC20(USDC).safeApprove(address(USDCTokenPool), type(uint256).max);
    USDCTokenPool.provideLiquidity(RONIN_CHAIN_SELECTOR, liquidityAmount);

    assertEq(USDCTokenPool.getLockedTokensForChain(RONIN_CHAIN_SELECTOR), liquidityAmount);

    vm.startPrank(RONIN_OFF_RAMP);

    vm.expectEmit();
    emit TokenPool.Released(RONIN_OFF_RAMP, address(0xdeadbeef), liquidityAmount);

    // Attempt an incoming message from Ronin to release the liquidity
    USDCTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: abi.encode("FAKE_SENDER"),
        remoteChainSelector: RONIN_CHAIN_SELECTOR,
        receiver: address(0xdeadbeef),
        amount: liquidityAmount,
        localToken: USDC,
        sourcePoolAddress: abi.encode("FAKE_POOL"),
        sourcePoolData: abi.encode(LOCK_RELEASE_FLAG),
        offchainTokenData: ""
      })
    );

    // Assert that the liquidity was released
    assertEq(IERC20(USDC).balanceOf(address(USDCTokenPool)), 0);
    assertEq(USDCTokenPool.getLockedTokensForChain(RONIN_CHAIN_SELECTOR), 0);

  }
}
