package verification

import (
	"testing"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"
	dsutil "github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/v0_5"
	"github.com/stretchr/testify/require"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	dsTypes "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
)

// DeployVerifierProxyAndVerifier deploys a VerifierProxy, deploys a Verifier,
// initializes the VerifierProxy with the Verifier, and returns the updated
// environment and the addresses of VerifierProxy and Verifier.
func DeployVerifierProxyAndVerifier(
	t *testing.T,
	e deployment.Environment,
) (env deployment.Environment, verifierProxyAddr common.Address, verifierAddr common.Address) {
	t.Helper()

	chainSelector := testutil.TestChain.Selector

	// 1) Deploy VerifierProxy
	deployProxyCfg := DeployVerifierProxyConfig{
		Version: deployment.Version0_5_0,
		ChainsToDeploy: map[uint64]DeployVerifierProxy{
			chainSelector: {
				AccessControllerAddress: common.Address{},
			},
		},
	}
	env, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			DeployVerifierProxyChangeset,
			deployProxyCfg,
		),
	)
	require.NoError(t, err, "deploying verifier proxy should not fail")

	// Get the VerifierProxy address
	record, err := env.DataStore.Addresses().Get(ds.NewAddressRefKey(chainSelector, ds.ContractType(dsTypes.VerifierProxy), &deployment.Version0_5_0, ""))
	require.NoError(t, err)
	verifierProxyAddr = common.HexToAddress(record.Address)

	// 2) Deploy Verifier
	deployVerifierCfg := DeployVerifierConfig{
		ChainsToDeploy: map[uint64]DeployVerifier{
			chainSelector: {
				VerifierProxyAddress: verifierProxyAddr,
			},
		},
	}
	env, err = commonChangesets.Apply(t, env, nil,
		commonChangesets.Configure(
			DeployVerifierChangeset,
			deployVerifierCfg,
		),
	)
	require.NoError(t, err, "deploying verifier should not fail")

	// Get the Verifier address
	record, err = env.DataStore.Addresses().Get(ds.NewAddressRefKey(chainSelector, ds.ContractType(dsTypes.Verifier), &deployment.Version0_5_0, ""))
	require.NoError(t, err)
	verifierAddr = common.HexToAddress(record.Address)

	// 3) Initialize the VerifierProxy
	initCfg := VerifierProxyInitializeVerifierConfig{
		ConfigPerChain: map[uint64][]InitializeVerifierConfig{
			chainSelector: {
				{
					VerifierAddress:      verifierAddr,
					VerifierProxyAddress: verifierProxyAddr,
				},
			},
		},
	}
	env, err = commonChangesets.Apply(t, env, nil,
		commonChangesets.Configure(
			InitializeVerifierChangeset,
			initCfg,
		),
	)
	require.NoError(t, err, "initializing verifier proxy should not fail")

	return env, verifierProxyAddr, verifierAddr
}

func VerifyState(
	t *testing.T,
	inDs ds.MutableDataStore[ds.DefaultMetadata, ds.DefaultMetadata],
	chainSelector uint64,
	contractAddress common.Address,
	expectedConfig SetConfig,
	shouldConfigBeActive bool,
) {

	envDatastore, err := ds.FromDefault[metadata.SerializedContractMetadata, ds.DefaultMetadata](inDs.Seal())
	require.NoError(t, err)

	// Retrieve contract metadata from datastore
	cm, err := envDatastore.ContractMetadata().Get(
		ds.NewContractMetadataKey(chainSelector, contractAddress.String()),
	)
	require.NoError(t, err, "Failed to get contract metadata")

	contractMetadata, err := metadata.DeserializeMetadata[v0_5.VerifierView](cm.Metadata)
	require.NoError(t, err, "Failed to convert contract metadata")

	configDigestString := dsutil.HexEncodeBytes32(expectedConfig.ConfigDigest)

	// Retrieve the config state
	configState, err := contractMetadata.View.GetVerifierState(configDigestString)
	require.NoError(t, err, "Failed to get config state")

	// Verify basic configuration properties
	require.Equal(t, expectedConfig.F, configState.F, "F value mismatch")
	require.Equal(t, configDigestString, configState.ConfigDigest, "ConfigDigest mismatch")
	require.Equal(t, shouldConfigBeActive, configState.IsActive, "IsActive mismatch")

	stringSigners := make([]string, len(expectedConfig.Signers))
	for i, signer := range expectedConfig.Signers {
		stringSigners[i] = signer.String()
	}

	require.Equal(t, stringSigners, configState.Signers, "Signers mismatch")

	t.Log("All state verifications passed")
}
