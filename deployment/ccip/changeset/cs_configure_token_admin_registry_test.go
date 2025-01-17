package changeset_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestValidateRegistryConfig(t *testing.T) {
	t.Parallel()

	e, selectorA, _, tokens, timelockContracts := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)
	invalidPoolAddress := utils.RandomAddress()
	administrator := utils.RandomAddress()

	e = testhelpers.DeployTestTokenPools(t, e, map[uint64]changeset.DeployTokenPoolInput{
		selectorA: {
			Type:               changeset.BurnMintTokenPool,
			TokenAddress:       tokens[selectorA].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
	}, true)

	// Deploy another token pool with force enabled
	e = testhelpers.DeployTestTokenPools(t, e, map[uint64]changeset.DeployTokenPoolInput{
		selectorA: {
			Type:               changeset.BurnMintTokenPool,
			TokenAddress:       tokens[selectorA].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
			ForceDeployment:    true,
		},
	}, true)

	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	poolAddress1 := state.Chains[selectorA].BurnMintTokenPools[testhelpers.TestTokenSymbol][0].Address()
	poolAddress2 := state.Chains[selectorA].BurnMintTokenPools[testhelpers.TestTokenSymbol][1].Address()
	tokenAdminRegistry := state.Chains[selectorA].TokenAdminRegistry

	// We want to transfer the admin role of the pool to the deployer to get a validation failure in the last test
	e, err = commonchangeset.ApplyChangesets(t, e, timelockContracts, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.ConfigureTokenAdminRegistry),
			Config: changeset.ConfigureTokenAdminRegistryConfig{
				TokenSymbol: testhelpers.TestTokenSymbol,
				MCMS: &changeset.MCMSConfig{
					MinDelay: 0 * time.Second,
				},
				RegistryUpdates: map[uint64]changeset.RegistryConfig{
					selectorA: {
						PoolAddress:   poolAddress1,
						Administrator: e.Chains[selectorA].DeployerKey.From,
					},
				},
			},
		},
	})
	require.NoError(t, err)
	tx, err := tokenAdminRegistry.AcceptAdminRole(e.Chains[selectorA].DeployerKey, tokens[selectorA].Address)
	require.NoError(t, err)
	_, err = e.Chains[selectorA].Confirm(tx)
	require.NoError(t, err)

	tests := []struct {
		Msg            string
		UseMcms        bool
		TokenSymbol    changeset.TokenSymbol
		RegistryConfig changeset.RegistryConfig
		ErrStr         string
	}{
		{
			Msg:            "Pool address is missing",
			RegistryConfig: changeset.RegistryConfig{},
			ErrStr:         "pool address must be defined",
		},
		{
			Msg:         "Pool doesn't exist",
			TokenSymbol: testhelpers.TestTokenSymbol,
			RegistryConfig: changeset.RegistryConfig{
				PoolAddress:   invalidPoolAddress,
				Administrator: administrator,
			},
			ErrStr: fmt.Sprintf("failed to find token pool on %d with symbol %s and address %s", selectorA, testhelpers.TestTokenSymbol, invalidPoolAddress),
		},
		{
			Msg:         "Token admin registry is not owned by required address",
			TokenSymbol: testhelpers.TestTokenSymbol,
			RegistryConfig: changeset.RegistryConfig{
				PoolAddress:   poolAddress2,
				Administrator: administrator,
			},
			ErrStr: fmt.Sprintf("token admin registry failed ownership validation on %d", selectorA),
		},
		{
			Msg:         "Owner can't become admin of token",
			TokenSymbol: testhelpers.TestTokenSymbol,
			UseMcms:     true,
			RegistryConfig: changeset.RegistryConfig{
				PoolAddress:   poolAddress2,
				Administrator: administrator,
			},
			ErrStr: fmt.Sprintf("address %s is unable to be the admin of %s on %d", state.Chains[selectorA].Timelock.Address(), testhelpers.TestTokenSymbol, selectorA),
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			err := test.RegistryConfig.Validate(e.GetContext(), e.Chains[selectorA], state.Chains[selectorA], test.UseMcms, test.TokenSymbol)
			require.Error(t, err)
			require.ErrorContains(t, err, test.ErrStr)
		})
	}
}

func TestValidateConfigureTokenAdminRegistryConfig(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 2,
	})

	tests := []struct {
		TokenSymbol changeset.TokenSymbol
		Input       changeset.ConfigureTokenAdminRegistryConfig
		ErrStr      string
		Msg         string
	}{
		{
			Msg:    "Token symbol is missing",
			Input:  changeset.ConfigureTokenAdminRegistryConfig{},
			ErrStr: "token symbol must be defined",
		},
		{
			Msg: "Chain selector is invalid",
			Input: changeset.ConfigureTokenAdminRegistryConfig{
				TokenSymbol: testhelpers.TestTokenSymbol,
				RegistryUpdates: map[uint64]changeset.RegistryConfig{
					0: changeset.RegistryConfig{},
				},
			},
			ErrStr: "failed to validate chain selector 0",
		},
		{
			Msg: "Chain selector doesn't exist in environment",
			Input: changeset.ConfigureTokenAdminRegistryConfig{
				TokenSymbol: testhelpers.TestTokenSymbol,
				RegistryUpdates: map[uint64]changeset.RegistryConfig{
					5009297550715157269: changeset.RegistryConfig{},
				},
			},
			ErrStr: "chain with selector 5009297550715157269 does not exist in environment",
		},
		{
			Msg: "Token admin registry is missing",
			Input: changeset.ConfigureTokenAdminRegistryConfig{
				TokenSymbol: testhelpers.TestTokenSymbol,
				RegistryUpdates: map[uint64]changeset.RegistryConfig{
					e.AllChainSelectors()[0]: changeset.RegistryConfig{},
				},
			},
			ErrStr: fmt.Sprintf("missing tokenAdminRegistry on %d", e.AllChainSelectors()[0]),
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			err := test.Input.Validate(e)
			require.Contains(t, err.Error(), test.ErrStr)
		})
	}
}

func TestConfigureTokenAdminRegistry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Administrator common.Address
		MCMS          *changeset.MCMSConfig
		Msg           string
	}{
		{
			Msg: "Configure with MCMS",
			MCMS: &changeset.MCMSConfig{
				MinDelay: 0 * time.Second,
			},
		},
		{
			Msg: "Configure without MCMS",
		},
		{
			Msg: "Configure with MCMS & transfer",
			MCMS: &changeset.MCMSConfig{
				MinDelay: 0 * time.Second,
			},
			Administrator: utils.RandomAddress(),
		},
		{
			Msg:           "Configure without MCMS & transfer",
			Administrator: utils.RandomAddress(),
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			e, selectorA, _, tokens, timelockContracts := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), test.MCMS != nil)

			e = testhelpers.DeployTestTokenPools(t, e, map[uint64]changeset.DeployTokenPoolInput{
				selectorA: {
					Type:               changeset.BurnMintTokenPool,
					TokenAddress:       tokens[selectorA].Address,
					LocalTokenDecimals: testhelpers.LocalTokenDecimals,
				},
			}, test.MCMS != nil)

			state, err := changeset.LoadOnchainState(e)
			require.NoError(t, err)

			tokenAddress := tokens[selectorA].Address
			timelockAddress := state.Chains[selectorA].Timelock.Address()
			poolAddress := state.Chains[selectorA].BurnMintTokenPools[testhelpers.TestTokenSymbol][0].Address()
			tokenAdminRegistry := state.Chains[selectorA].TokenAdminRegistry
			deployerKey := e.Chains[selectorA].DeployerKey.From

			e, err = commonchangeset.ApplyChangesets(t, e, timelockContracts, []commonchangeset.ChangesetApplication{
				{
					Changeset: commonchangeset.WrapChangeSet(changeset.ConfigureTokenAdminRegistry),
					Config: changeset.ConfigureTokenAdminRegistryConfig{
						TokenSymbol: testhelpers.TestTokenSymbol,
						MCMS:        test.MCMS,
						RegistryUpdates: map[uint64]changeset.RegistryConfig{
							selectorA: {
								PoolAddress:   poolAddress,
								Administrator: test.Administrator,
							},
						},
					},
				},
			})
			require.NoError(t, err)

			tokenConfig, err := tokenAdminRegistry.GetTokenConfig(nil, tokenAddress)
			require.NoError(t, err)

			if test.Administrator != utils.ZeroAddress {
				require.Equal(t, test.Administrator, tokenConfig.PendingAdministrator)
			}
			if test.MCMS == nil {
				require.Equal(t, deployerKey, tokenConfig.Administrator)
			} else {
				require.Equal(t, timelockAddress, tokenConfig.Administrator)
			}
			require.Equal(t, poolAddress, tokenConfig.TokenPool)
		})
	}
}
