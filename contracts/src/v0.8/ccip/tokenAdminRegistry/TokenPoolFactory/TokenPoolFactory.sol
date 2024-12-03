// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IOwnable} from "../../../shared/interfaces/IOwnable.sol";
import {ITypeAndVersion} from "../../../shared/interfaces/ITypeAndVersion.sol";
import {ITokenAdminRegistry} from "../../interfaces/ITokenAdminRegistry.sol";

import {RateLimiter} from "../../libraries/RateLimiter.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {RegistryModuleOwnerCustom} from "../RegistryModuleOwnerCustom.sol";
import {FactoryBurnMintERC20} from "./FactoryBurnMintERC20.sol";

import {Create2} from "../../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/Create2.sol";

/// @notice A contract for deploying new tokens and token pools, and configuring them with the token admin registry
/// @dev At the end of the transaction, the ownership transfer process will begin, but the user must accept the
/// ownership transfer in a separate transaction.
/// @dev The address prediction mechanism is only capable of deploying and predicting addresses for EVM based chains.
/// adding compatibility for other chains will require additional offchain computation.
contract TokenPoolFactory is ITypeAndVersion {
  using Create2 for bytes32;

  error InvalidZeroAddress();

  /// @dev This struct will only ever exist in memory and as calldata, and therefore does not need to be efficiently packed for storage. The struct is used to pass information to the create2 address generation function.
  struct RemoteTokenPoolInfo {
    uint64 remoteChainSelector; // The CCIP specific selector for the remote chain
    bytes[] remotePoolAddresses; // The address of the remote pool to either deploy or use as is. If empty, address
    // will be predicted
    bytes remoteTokenAddress; // EVM address for remote token. If empty, the address will be predicted
    RateLimiter.Config rateLimiterConfig; // Token Pool rate limit. Values will be applied on incoming an outgoing messages
  }

  string public constant typeAndVersion = "TokenPoolFactory 1.5.1";

  ITokenAdminRegistry private immutable i_tokenAdminRegistry;
  RegistryModuleOwnerCustom private immutable i_registryModuleOwnerCustom;

  /// @notice Construct the TokenPoolFactory
  /// @param tokenAdminRegistry The address of the token admin registry
  /// @param tokenAdminModule The address of the token admin module which can register the token via ownership module
  constructor(ITokenAdminRegistry tokenAdminRegistry, RegistryModuleOwnerCustom tokenAdminModule) {
    if (address(tokenAdminRegistry) == address(0) || address(tokenAdminModule) == address(0)) {
      revert InvalidZeroAddress();
    }

    i_tokenAdminRegistry = ITokenAdminRegistry(tokenAdminRegistry);
    i_registryModuleOwnerCustom = RegistryModuleOwnerCustom(tokenAdminModule);
  }

  // ================================================================
  // │                   Top-Level Deployment                       │
  // ================================================================

  /// @notice Deploys a token and token pool with the given token information and configures it with remote token pools
  /// @dev The token and token pool are deployed in the same transaction, and the token pool is configured with the
  /// remote token pools. The token pool is then set in the token admin registry. Ownership of the everything is transferred
  /// to the msg.sender, but must be accepted in a separate transaction due to 2-step ownership transfer.
  /// @param remoteTokenPools An array of remote token pools info to be used in the pool's applyChainUpdates function
  /// or to be predicted if the pool has not been deployed yet on the remote chain
  /// @param tokenInitCode The creation code for the token, which includes the constructor parameters already appended
  /// @param tokenPoolInitCode The creation code for the token pool, without the constructor parameters appended
  /// @param salt The salt to be used in the create2 deployment of the token and token pool to ensure a unique address
  /// @return token The address of the token that was deployed
  /// @return pool The address of the token pool that was deployed
  function deployTokenAndTokenPool(
    RemoteTokenPoolInfo[] calldata remoteTokenPools,
    bytes memory tokenInitCode,
    bytes calldata tokenPoolInitCode,
    bytes32 salt
  ) external returns (address, address) {
    // Ensure a unique deployment between senders even if the same input parameter is used to prevent
    // DOS/front running attacks
    salt = keccak256(abi.encodePacked(salt, msg.sender));

    // Deploy the token. The constructor parameters are already provided in the tokenInitCode
    address token = Create2.deploy(0, salt, tokenInitCode);

    // Deploy the token pool
    address pool = _createTokenPool(remoteTokenPools, tokenPoolInitCode, salt);

    // Grant the mint and burn roles to the pool for the token
    FactoryBurnMintERC20(token).grantMintAndBurnRoles(pool);

    // Set the token pool for token in the token admin registry since this contract is the token and pool owner
    _setTokenPoolInTokenAdminRegistry(token, pool);

    // Begin the 2 step ownership transfer of the newly deployed token to the msg.sender
    IOwnable(token).transferOwnership(msg.sender);

    return (token, pool);
  }

  /// @notice Deploys a token pool with an existing ERC20 token
  /// @dev Since the token already exists, this contract is not the owner and therefore cannot configure the
  /// token pool in the token admin registry in the same transaction. The user must invoke the calls to the
  /// tokenAdminRegistry manually
  /// @dev since the token already exists, the owner must grant the mint and burn roles to the pool manually
  /// @param remoteTokenPools An array of remote token pools info to be used in the pool's applyChainUpdates function
  /// @param tokenPoolInitCode The creation code for the token pool
  /// @param salt The salt to be used in the create2 deployment of the token pool
  /// @return poolAddress The address of the token pool that was deployed
  function deployTokenPoolWithExistingToken(
    RemoteTokenPoolInfo[] calldata remoteTokenPools,
    bytes calldata tokenPoolInitCode,
    bytes32 salt
  ) external returns (address poolAddress) {
    // Ensure a unique deployment between senders even if the same input parameter is used to prevent
    // DOS/front running attacks
    salt = keccak256(abi.encodePacked(salt, msg.sender));

    // create the token pool and return the address
    return _createTokenPool(remoteTokenPools, tokenPoolInitCode, salt);
  }

  // ================================================================
  // │                Pool Deployment/Configuration                 │
  // ================================================================

  /// @notice Deploys a token pool with the given token information and remote token pools
  /// @param remoteTokenPools An array of remote token pools info to be used in the pool's applyChainUpdates function
  /// @param tokenPoolInitCode The creation code for the token pool
  /// @param salt The salt to be used in the create2 deployment of the token pool
  /// @return poolAddress The address of the token pool that was deployed
  function _createTokenPool(
    RemoteTokenPoolInfo[] calldata remoteTokenPools,
    bytes calldata tokenPoolInitCode,
    bytes32 salt
  ) private returns (address) {
    // Create an array of chain updates to apply to the token pool
    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](remoteTokenPools.length);

    RemoteTokenPoolInfo memory remoteTokenPool;
    for (uint256 i = 0; i < remoteTokenPools.length; ++i) {
      remoteTokenPool = remoteTokenPools[i];

      chainUpdates[i] = TokenPool.ChainUpdate({
        remoteChainSelector: remoteTokenPool.remoteChainSelector,
        remotePoolAddresses: remoteTokenPool.remotePoolAddresses,
        remoteTokenAddress: remoteTokenPool.remoteTokenAddress,
        outboundRateLimiterConfig: remoteTokenPool.rateLimiterConfig,
        inboundRateLimiterConfig: remoteTokenPool.rateLimiterConfig
      });
    }

    // Construct the deployment code from the initCode and the initArgs and then deploy
    address poolAddress = Create2.deploy(0, salt, tokenPoolInitCode);

    // Apply the chain updates to the token pool
    TokenPool(poolAddress).applyChainUpdates(new uint64[](0), chainUpdates);

    // Begin the 2 step ownership transfer of the token pool to the msg.sender.
    IOwnable(poolAddress).transferOwnership(address(msg.sender)); // 2 step ownership transfer

    return poolAddress;
  }

  /// @notice Sets the token pool address in the token admin registry for a newly deployed token pool.
  /// @dev this function should only be called when the token is deployed by this contract as well, otherwise
  /// the token pool will not be able to be set in the token admin registry, and this function will revert.
  /// @param token The address of the token to set the pool for
  /// @param pool The address of the pool to set in the token admin registry
  function _setTokenPoolInTokenAdminRegistry(address token, address pool) private {
    i_registryModuleOwnerCustom.registerAdminViaOwner(token);
    i_tokenAdminRegistry.acceptAdminRole(token);
    i_tokenAdminRegistry.setPool(token, pool);

    // Begin the 2 admin transfer process which must be accepted in a separate tx.
    i_tokenAdminRegistry.transferAdminRole(token, msg.sender);
  }
}
