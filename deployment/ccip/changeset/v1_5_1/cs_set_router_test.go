package v1_5_1_test

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/stretchr/testify/require"
)

func TestSetRouterOnTokenPool_TestRouterWithMCMS(t *testing.T) {
	t.Parallel()

	e, selectorA, _, tokens, timelockContracts := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	// Deploy the test router and save to env
	address, _, _, err := router.DeployRouter(e.Chains[selectorA].DeployerKey, e.Chains[selectorA].Client, common.Address{}, common.Address{})
	require.NoError(t, err)
	e.E=.Save(selectorA, address.Hex(), deployment.NewTypeAndVersion(changeset.TestRouter, deployment.Version1_2_0))

	e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
		selectorA: {
			Type:               changeset.BurnMintTokenPool,
			TokenAddress:       tokens[selectorA].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
	}, true, true)

	mcmsConfig := &proposalutils.TimelockConfig{
		MinDelay: 0 * time.Second,
	}

	tests := []struct {
		Config changeset.TokenAdminRegistryChangesetConfig
		ErrStr string
		Msg    string
	}{
		{
			Msg: "Chain selector is invalid",
			Config: changeset.TokenAdminRegistryChangesetConfig{
				Pools: map[uint64]map[changeset.TokenSymbol]changeset.TokenPoolInfo{
					0: {},
				},
			},
			ErrStr: "failed to validate chain selector 0",
		},
		{
			Msg: "Chain selector doesn't exist in environment",
			Config: changeset.TokenAdminRegistryChangesetConfig{
				Pools: map[uint64]map[changeset.TokenSymbol]changeset.TokenPoolInfo{
					5009297550715157269: {},
				},
			},
			ErrStr: "does not exist in environment",
		},
		{
			Msg: "Invalid pool type",
			Config: changeset.TokenAdminRegistryChangesetConfig{
				MCMS: mcmsConfig,
				Pools: map[uint64]map[changeset.TokenSymbol]changeset.TokenPoolInfo{
					selectorA: {
						testhelpers.TestTokenSymbol: {
							Type:    "InvalidType",
							Version: deployment.Version1_5_1,
						},
					},
				},
			},
			ErrStr: "InvalidType is not a known token pool type",
		},
		{
			Msg: "Invalid pool version",
			Config: changeset.TokenAdminRegistryChangesetConfig{
				MCMS: mcmsConfig,
				Pools: map[uint64]map[changeset.TokenSymbol]changeset.TokenPoolInfo{
					selectorA: {
						testhelpers.TestTokenSymbol: {
							Type:    changeset.BurnMintTokenPool,
							Version: deployment.Version1_0_0,
						},
					},
				},
			},
			ErrStr: "1.0.0 is not a known token pool version",
		},
		{
			Msg: "External admin undefined",
			Config: changeset.TokenAdminRegistryChangesetConfig{
				MCMS: mcmsConfig,
				Pools: map[uint64]map[changeset.TokenSymbol]changeset.TokenPoolInfo{
					selectorA: {
						testhelpers.TestTokenSymbol: {
							Type:    changeset.BurnMintTokenPool,
							Version: deployment.Version1_5_1,
						},
					},
				},
			},
			ErrStr: "external admin must be defined",
		},
		{
			Msg: "Not admin",
			Config: changeset.TokenAdminRegistryChangesetConfig{
				MCMS: mcmsConfig,
				Pools: map[uint64]map[changeset.TokenSymbol]changeset.TokenPoolInfo{
					selectorA: {
						testhelpers.TestTokenSymbol: {
							Type:          changeset.BurnMintTokenPool,
							Version:       deployment.Version1_5_1,
							ExternalAdmin: utils.RandomAddress(),
						},
					},
				},
			},
			ErrStr: "is not the administrator",
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			_, err := commonchangeset.Apply(t, e, timelockContracts,
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_5_1.TransferAdminRoleChangeset),
					test.Config,
				),
			)
			require.Error(t, err)
			require.ErrorContains(t, err, test.ErrStr)
		})
	}
}
