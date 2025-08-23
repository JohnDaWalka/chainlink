package cre

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	df "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	df_sol "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/solana"
	ks_sol "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/stretchr/testify/require"
)

var (
	solWriterDonConfig = "../../../../core/scripts/cre/environment/configs/workflow-solana-don-cache.toml"
)

func Test_CRE_WorkflowDon_WriteSolana(t *testing.T) {
	confErr := setConfigurationIfMissing(solWriterDonConfig, "workflow-solana")
	require.NoError(t, confErr, "failed to set configuration")

	configurationFiles := os.Getenv("CTF_CONFIGS")
	require.NotEmpty(t, configurationFiles, "CTF_CONFIGS env var is not set")

	topology := os.Getenv("CRE_TOPOLOGY")
	require.NotEmpty(t, topology, "CRE_TOPOLOGY env var is not set")

	createErr := createEnvironmentIfNotExists(configurationFiles, "../../../../core/scripts/cre/environment", topology)
	require.NoError(t, createErr, "failed to create environment")

	/*
		LOAD ENVIRONMENT STATE
	*/
	in, err := framework.Load[envconfig.Config](nil)
	require.NoError(t, err, "couldn't load environment state")

	var envArtifact environment.EnvArtifact
	artFile, err := os.ReadFile(os.Getenv("ENV_ARTIFACT_PATH"))
	require.NoError(t, err, "failed to read artifact file")
	err = json.Unmarshal(artFile, &envArtifact)
	require.NoError(t, err, "failed to unmarshal artifact file")

	executeSecureMintTest(t, in, envArtifact)
}

type setup struct {
	ForwarderProgramID solana.PublicKey
	ForwarderState     solana.PublicKey
	CacheProgramID     solana.PublicKey
	CacheState         solana.PublicKey

	FeedID       string
	Descriptions [][32]byte
	WFOwner      [20]byte
	WFName       string
}

var (
	feedID        = "0x018e16c39e00032000000"
	wFName        = "test workflow"
	wFDescription = "securemint test"
	wFOwner       = [20]byte{1, 2, 3}
)

func executeSecureMintTest(t *testing.T, in *envconfig.Config, envArtifact environment.EnvArtifact) {
	cldLogger := cldlogger.NewSingleFileLogger(t)

	fullCldEnvOutput, wrappedBlockchainOutputs, loadErr := environment.BuildFromSavedState(t.Context(), cldLogger, in, envArtifact)
	require.NoError(t, loadErr, "failed to load environment")
	ds := fullCldEnvOutput.Environment.DataStore

	// prevalidate environment
	forwarders := fullCldEnvOutput.Environment.DataStore.Addresses().Filter(
		datastore.AddressRefByQualifier(ks_sol.DefaultForwarderQualifier),
		datastore.AddressRefByType(ks_sol.ForwarderContract))
	require.Len(t, forwarders, 1)

	forwarderStates := fullCldEnvOutput.Environment.DataStore.Addresses().Filter(
		datastore.AddressRefByQualifier(ks_sol.DefaultForwarderQualifier),
		datastore.AddressRefByType(ks_sol.ForwarderState))
	require.Len(t, forwarderStates, 1)

	var s setup
	var solChain *cre.WrappedBlockchainOutput
	for _, w := range wrappedBlockchainOutputs {
		if w.BlockchainOutput.Type != blockchain.FamilySolana {
			continue
		}

		s.ForwarderProgramID = mustGetContract(t, ds, w.SolChain.ChainSelector, ks_sol.ForwarderContract)
		s.ForwarderState = mustGetContract(t, ds, w.SolChain.ChainSelector, ks_sol.ForwarderState)
		// we assume we always have just 1 solana chain
		solChain = w
		break
	}

	deployAndConfigureCache(t, &s, *fullCldEnvOutput.Environment, solChain)
}

func deployAndConfigureCache(t *testing.T, s *setup, env cldf.Environment, solChain *cre.WrappedBlockchainOutput) {
	var d [32]byte
	copy(d[:], []byte(wFDescription))
	s.Descriptions = append(s.Descriptions, d)
	s.WFName = wFName
	s.WFOwner = wFOwner
	s.FeedID = feedID
	// deploy df cache
	deployCS := commonchangeset.Configure(df_sol.DeployCache{}, &df_sol.DeployCacheRequest{
		ChainSel:           solChain.SolChain.ChainSelector,
		Qualifier:          ks_sol.DefaultForwarderQualifier,
		Version:            "1.0.0",
		FeedAdmins:         []solana.PublicKey{solChain.SolChain.PrivateKey.PublicKey()},
		ForwarderProgramID: s.ForwarderProgramID,
	})

	// init decimal report
	initCS := commonchangeset.Configure(df_sol.InitCacheDecimalReport{},
		&df_sol.InitCacheDecimalReportRequest{
			ChainSel:  solChain.SolChain.ChainSelector,
			Qualifier: ks_sol.DefaultForwarderQualifier,
			Version:   "1.0.0",
			FeedAdmin: solChain.SolChain.PrivateKey.PublicKey(),
			DataIDs:   []string{s.FeedID},
		})

	// configure decimal report
	configureCS := commonchangeset.Configure(df_sol.ConfigureCacheDecimalReport{},
		&df_sol.ConfigureCacheDecimalReportRequest{
			ChainSel:  solChain.SolChain.ChainSelector,
			Qualifier: ks_sol.DefaultForwarderQualifier,
			Version:   "1.0.0",
			SenderList: []df_sol.Sender{
				{
					ProgramID: s.ForwarderProgramID,
					StateID:   s.ForwarderState,
				},
			},
			FeedAdmin:            solChain.SolChain.PrivateKey.PublicKey(),
			DataIDs:              []string{s.FeedID},
			AllowedWorkflowOwner: [][20]byte{s.WFOwner},
			AllowedWorkflowName:  [][10]byte{df.HashedWorkflowName(s.WFName)},
			Descriptions:         s.Descriptions,
		})

	env, _, cacheErr := commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{deployCS, initCS, configureCS})
	require.NoError(t, cacheErr)
	s.CacheProgramID = mustGetContract(t, env.DataStore, solChain.SolChain.ChainSelector, df_sol.CacheContract)
	s.CacheState = mustGetContract(t, env.DataStore, solChain.SolChain.ChainSelector, df_sol.CacheState)
}

func mustGetContract(t *testing.T, ds datastore.DataStore, sel uint64, ctype datastore.ContractType) solana.PublicKey {
	key := datastore.NewAddressRefKey(
		sel,
		ctype,
		semver.MustParse("1.0.0"),
		ks_sol.DefaultForwarderQualifier,
	)
	contract, err := ds.Addresses().Get(key)

	require.NoError(t, err)

	return solana.MustPublicKeyFromBase58(contract.Address)
}
