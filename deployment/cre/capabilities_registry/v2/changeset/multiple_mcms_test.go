package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/cre"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
)

// TestMultipleMCMSDeploymentsConflict demonstrates the gap where GetMCMSContracts
// cannot distinguish between multiple MCMS deployments on the same chain
func TestMultipleMCMSDeploymentsConflict(t *testing.T) {
	lggr := logger.Test(t)
	env, chainSelector := cre.BuildMinimalEnvironment(t, lggr)

	t.Log("=== Setting up Team A's MCMS infrastructure ===")

	// Deploy Team A's MCMS infrastructure
	teamATimelockCfgs := map[uint64]commontypes.MCMSWithTimelockConfigV2{
		chainSelector: proposalutils.SingleGroupTimelockConfigV2(t),
	}

	teamAEnv, err := commonchangeset.Apply(t, env,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
			teamATimelockCfgs,
		),
	)
	require.NoError(t, err, "failed to deploy Team A's MCMS infrastructure")
	t.Log("Team A's MCMS infrastructure deployed successfully")

	// Get Team A's MCMS contracts
	teamAMCMSContracts, err := strategies.GetMCMSContracts(teamAEnv, chainSelector)
	require.NoError(t, err, "should be able to get Team A's MCMS contracts")
	require.NotNil(t, teamAMCMSContracts, "Team A's MCMS contracts should not be nil")

	teamATimelockAddr := teamAMCMSContracts.Timelock.Address()
	teamAProposerAddr := teamAMCMSContracts.ProposerMcm.Address()
	t.Logf("Team A - Timelock: %s, Proposer: %s", teamATimelockAddr.Hex(), teamAProposerAddr.Hex())

	t.Log("=== Setting up Team B's MCMS infrastructure ===")

	// Deploy Team B's MCMS infrastructure ON THE SAME CHAIN
	// This simulates two independent teams deploying their own governance
	teamBTimelockCfgs := map[uint64]commontypes.MCMSWithTimelockConfigV2{
		chainSelector: proposalutils.SingleGroupTimelockConfigV2(t),
	}

	teamBEnv, err := commonchangeset.Apply(t, teamAEnv, // Build on top of Team A's environment
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
			teamBTimelockCfgs,
		),
	)
	require.NoError(t, err, "failed to deploy Team B's MCMS infrastructure")
	t.Log("Team B's MCMS infrastructure deployed successfully")

	// Both teams deploy their registries in the SAME environment (teamBEnv)
	// This simulates the real-world scenario where both teams' contracts exist on same chain
	teamARegistry, err := DeployCapabilitiesRegistry{}.Apply(teamBEnv, DeployCapabilitiesRegistryInput{
		ChainSelector: chainSelector,
		Qualifier:     "team-a-registry",
	})
	require.NoError(t, err, "failed to deploy Team A's capabilities registry")

	teamBRegistry, err := DeployCapabilitiesRegistry{}.Apply(teamBEnv, DeployCapabilitiesRegistryInput{
		ChainSelector: chainSelector,
		Qualifier:     "team-b-registry",
	})
	require.NoError(t, err, "failed to deploy Team B's capabilities registry")

	teamARegistryAddr := teamARegistry.Reports[0].Output.(contracts.DeployCapabilitiesRegistryOutput).Address
	teamBRegistryAddr := teamBRegistry.Reports[0].Output.(contracts.DeployCapabilitiesRegistryOutput).Address

	t.Logf("Team A Registry: %s", teamARegistryAddr)
	t.Logf("Team B Registry: %s", teamBRegistryAddr)

	// Team B tries to configure THEIR registry
	// But GetMCMSContracts() might return Team A's governance!
	teamBConfigInput := ConfigureCapabilitiesRegistryInput{
		ChainSelector:               chainSelector,
		CapabilitiesRegistryAddress: teamBRegistryAddr, // Team B's registry
		UseMCMS:                     true,
		MCMSConfig:                  &strategies.MCMSConfig{MinDuration: "30s"},
		Description:                 "Team B trying to configure THEIR OWN registry",
		Nops: []CapabilitiesRegistryNodeOperator{
			{
				Admin: common.HexToAddress("0x2222222222222222222222222222222222222222"),
				Name:  "Team B NOP",
			},
		},
	}

	// Get MCMS contracts that will be used for Team B's configuration
	teamBMCMSContracts, err := strategies.GetMCMSContracts(teamBEnv, chainSelector)
	require.NoError(t, err, "failed to get MCMS contracts for Team B")

	usedTimelockAddr := teamBMCMSContracts.Timelock.Address()
	usedProposerAddr := teamBMCMSContracts.ProposerMcm.Address()

	t.Logf("GetMCMSContracts returned: Timelock=%s, Proposer=%s", usedTimelockAddr.Hex(), usedProposerAddr.Hex())
	require.NotEqual(t, usedTimelockAddr, teamATimelockAddr, "GetMCMSContracts should return Team B's Timelock, not Team A's")

	_, err = ConfigureCapabilitiesRegistry{}.Apply(teamBEnv, teamBConfigInput)
	require.NoError(t, err, "Team B should be able to configure their registry successfully")
}
