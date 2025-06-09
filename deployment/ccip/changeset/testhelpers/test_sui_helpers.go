package testhelpers

import (
	"testing"

	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/sui"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

func DeployChainContractsToSuiCS(t *testing.T, e DeployedEnv, chainSelector uint64) commonchangeset.ConfiguredChangeSet {
	// Set mock link token address on Address book (to skip deploying)
	const mockLinkContract = "0xa"
	err := e.Env.ExistingAddresses.Save(chainSelector, mockLinkContract, cldf.NewTypeAndVersion(commontypes.LinkToken, deployment.Version1_6_0))
	require.NoError(t, err)

	//  Deploy contracts
	ccipConfig := sui.DeploySuiChainConfig{
		ContractParamsPerChain: map[uint64]sui.ChainContractParams{
			chainSelector: sui.GetMockChainContractParams(t, chainSelector),
		},
	}

	return commonchangeset.Configure(sui.DeploySuiChain{}, ccipConfig)
}
