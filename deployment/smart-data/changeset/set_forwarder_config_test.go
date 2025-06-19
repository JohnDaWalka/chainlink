package changeset_test

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commonTypes "github.com/smartcontractkit/chainlink/deployment/common/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	kschangeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/smart-data/changeset"
	"github.com/smartcontractkit/chainlink/deployment/smart-data/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/smart-data/changeset/types"
)

func TestSetForwarderConfig(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Chains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	chainSelector := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]

	newEnv, err := commonChangesets.Apply(t, env, commonChangesets.Configure(
		cldf.CreateLegacyChangeSet(kschangeset.DeployForwarderV2),
		&kschangeset.DeployRequestV2{
			ChainSel:  chainSelector,
			Qualifier: "my-test-forwarder",
		},
	), commonChangesets.Configure(
		cldf.CreateLegacyChangeSet(commonChangesets.DeployMCMSWithTimelockV2),
		map[uint64]commonTypes.MCMSWithTimelockConfigV2{
			chainSelector: proposalutils.SingleGroupTimelockConfigV2(t),
		},
	))
	require.NoError(t, err)

	forwarderAddress, err := cldf.SearchAddressBook(newEnv.ExistingAddresses, chainSelector, globals.KeystoneForwarder)
	require.NoError(t, err)

	signers := make([]common.Address, 10)
	for i := 0; i < 10; i++ {
		signers[i] = common.HexToAddress(fmt.Sprintf("0x%040x", i+1))
	}

	// without MCMS
	resp, err := commonChangesets.Apply(t, newEnv,
		commonChangesets.Configure(
			changeset.SetForwarderConfigChangeset,
			types.SetForwarderConfig{
				ForwarderAddress: common.HexToAddress(forwarderAddress),
				ChainSelector:    chainSelector,
				DonID:            1234,
				ConfigVersion:    1,
				F:                2,
				Signers:          signers,
			},
		),
	)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// with MCMS
	newEnv, err = commonChangesets.Apply(t, newEnv,
		commonChangesets.Configure(
			cldf.CreateLegacyChangeSet(commonChangesets.TransferToMCMSWithTimelockV2),
			commonChangesets.TransferToMCMSWithTimelockConfig{
				ContractsByChain: map[uint64][]common.Address{
					chainSelector: {common.HexToAddress(forwarderAddress)},
				},
				MCMSConfig: proposalutils.TimelockConfig{MinDelay: 0},
			},
		),
	)
	require.NoError(t, err)

	resp, err = commonChangesets.Apply(t, newEnv,
		commonChangesets.Configure(
			changeset.SetForwarderConfigChangeset,
			types.SetForwarderConfig{
				ForwarderAddress: common.HexToAddress(forwarderAddress),
				ChainSelector:    chainSelector,
				DonID:            5678,
				ConfigVersion:    2,
				F:                3,
				Signers:          signers,
				McmsConfig: &types.MCMSConfig{
					MinDelay: 0,
				},
			},
		),
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
}
