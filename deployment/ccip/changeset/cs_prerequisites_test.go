package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeployPrerequisites(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name   string
		config memory.MemoryEnvironmentConfig
	}{{
		name: "TestDeployPrerequisitesEVM",
		config: memory.MemoryEnvironmentConfig{
			Bootstraps: 1,
			Chains:     2,
			Nodes:      4,
		},
	}, { // this is failing now...
		name: "TestDeployPrerequisitesZKVM",
		config: memory.MemoryEnvironmentConfig{
			Bootstraps: 1,
			Chains:     0,
			Nodes:      4,
			ZKChains:   2,
		},
	}}

	lggr := logger.TestLogger(t)
	for _, testCase := range testCases {
		t.Run(testCase.name, testDeployPrerequisitesForChain(memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, testCase.config)))
	}
}

func testDeployPrerequisitesForChain(e deployment.Environment) func(t *testing.T) {
	return func(t *testing.T) {
		newChain := e.AllChainSelectors()[0]
		cfg := changeset.DeployPrerequisiteConfig{
			Configs: []changeset.DeployPrerequisiteConfigPerChain{
				{
					ChainSelector: newChain,
				},
			},
		}
		output, err := changeset.DeployPrerequisitesChangeset(e, cfg)
		require.NoError(t, err)
		err = e.ExistingAddresses.Merge(output.AddressBook)
		require.NoError(t, err)
		state, err := changeset.LoadOnchainState(e)
		require.NoError(t, err)
		require.NotNil(t, state.Chains[newChain].Weth9)
		require.NotNil(t, state.Chains[newChain].TokenAdminRegistry)
		require.NotNil(t, state.Chains[newChain].RegistryModule)
		require.NotNil(t, state.Chains[newChain].Router)
	}
}
