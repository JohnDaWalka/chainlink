package changeset_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/token_pool"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeployAndConfigureBnMHelperE2E(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)

	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     1,
		Nodes:      4,
	})

	selectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	chainSelector := selectors[0]

	// Deploy prerequisites first
	prereqCfg := changeset.DeployPrerequisiteConfig{
		Configs: []changeset.DeployPrerequisiteConfigPerChain{
			{
				ChainSelector: chainSelector,
				Opts: []changeset.PrerequisiteOpt{
					changeset.WithTokenPoolFactoryEnabled(),
				},
			},
		},
	}
	output, err := changeset.DeployPrerequisitesChangeset(e, prereqCfg)
	require.NoError(t, err)
	err = e.ExistingAddresses.Merge(output.AddressBook) //nolint
	require.NoError(t, err)

	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)
	chainState, ok := state.EVMChainState(chainSelector)
	require.True(t, ok)
	require.NotNil(t, chainState.RegistryModules1_6)
	require.Len(t, chainState.RegistryModules1_6, 1)

	var registryModuleAddress common.Address
	for _, registry := range chainState.RegistryModules1_6 {
		registryModuleAddress = registry.Address()
		break
	}

	// Token Pool Cfg
	rateLimiterConfig := token_pool.RateLimiterConfig{
		IsEnabled: false,
		Capacity:  big.NewInt(0),
		Rate:      big.NewInt(0),
	}

	tokenPoolConfig := v1_5_1.ConfigureTokenPoolContractsConfig{
		TokenSymbol: shared.CCIPBnMSymbol,
		PoolUpdates: map[uint64]v1_5_1.TokenPoolConfig{
			chainSelector: {
				ChainUpdates: v1_5_1.RateLimiterPerChain{
					chainSelector: v1_5_1.RateLimiterConfig{
						Inbound:  rateLimiterConfig,
						Outbound: rateLimiterConfig,
					},
				},
				Type:                    shared.BurnMintTokenPool,
				Version:                 deployment.Version1_5_1,
				SkipOwnershipValidation: true, // Skip ownership validation for testing
			},
		},
	}

	// BnM Helper token deployment and config input
	cfg := changeset.DeployAndConfigureBnMSelfServeConfig{
		Selector:                  chainSelector,
		TokenName:                 "CCIP-BnM",
		TokenSymbol:               string(shared.CCIPBnMSymbol),
		RegistryModuleOwnerCustom: registryModuleAddress.Hex(),
		TokenPoolConfig:           tokenPoolConfig,
	}

	err = changeset.DeployAndConfigureBnMSelfServe.VerifyPreconditions(e, cfg)
	require.NoError(t, err)

	changesetOutput, err := changeset.DeployAndConfigureBnMSelfServe.Apply(e, cfg)
	require.NoError(t, err)
	require.NotNil(t, changesetOutput)
	require.NotNil(t, changesetOutput.AddressBook) //nolint

	chainAddresses, err := changesetOutput.AddressBook.AddressesForChain(chainSelector) //nolint
	require.NoError(t, err)
	require.NotEmpty(t, chainAddresses, "Should have deployed contracts on the chain")

	// Reload state to verify the deployment and configuration
	updatedState, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)
	updatedChainState, ok := updatedState.EVMChainState(chainSelector)
	require.True(t, ok)

	// Verify BnM token was deployed and is present in the env
	bnmToken, exists := updatedChainState.BurnMintTokens677[shared.CCIPBnMSymbol]
	require.True(t, exists, "BnM token should exist in state")
	require.Contains(t, chainAddresses, bnmToken.Address().Hex(), "BnM token address should match deployed address")

	// Verify BnM token pool was deployed
	chain := e.BlockChains.EVMChains()[cfg.Selector]
	bnmTokenPoolAddr, exists := v1_5_1.GetTokenPoolAddressFromSymbolTypeAndVersion(updatedChainState, chain, shared.CCIPBnMSymbol, shared.BurnMintTokenPool, deployment.Version1_5_1)
	require.True(t, exists, "BnM token pool should exist in state")
	require.Contains(t, chainAddresses, bnmTokenPoolAddr.Hex(), "BnM token pool address should match deployed address")

	// Verify that the BnM token pool was deployed and configured correctly in state
	require.NotNil(t, updatedChainState.BurnMintTokenPools)
	bnmTokenPools, exists := updatedChainState.BurnMintTokenPools[shared.CCIPBnMSymbol]
	require.True(t, exists, "BnM token pools should exist in state")

	tokenPool, exists := bnmTokenPools[shared.CurrentTokenPoolVersion]
	require.True(t, exists, "BnM token pool with current version should exist in state")
	require.NotNil(t, tokenPool)

	// Verify that token admin registry is configured properly
	require.NotNil(t, updatedChainState.TokenAdminRegistry)
	tokenAdminReg := updatedChainState.TokenAdminRegistry
	require.True(t, exists, "TokenAdminRegistry should exist for the chain")
	isAdmin, err := tokenAdminReg.IsAdministrator(nil, bnmToken.Address(), chain.DeployerKey.From)
	require.NoError(t, err)
	require.True(t, isAdmin, "deployer should be an Owner of the token in TokenAdminRegistry")

	poolAddr, err := tokenAdminReg.GetPool(nil, bnmToken.Address())
	require.NoError(t, err)
	require.Equal(t, poolAddr, bnmTokenPoolAddr)
}
