package cre

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"strings"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	writetarget "github.com/smartcontractkit/chainlink-solana/pkg/solana/write_target"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	df "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	df_sol "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/solana"
	ks_sol "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
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
	Selector           uint64
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
	wFName        = "testwf"
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
	require.False(t, s.ForwarderProgramID.IsZero(), "failed to receive forwarder program id from blockchains output")
	s.Selector = solChain.SolChain.ChainSelector

	framework.L.Info().Msg("Deploy and configure data-feeds cache programs...")
	deployAndConfigureCache(t, &s, *fullCldEnvOutput.Environment, solChain)
	framework.L.Info().Msg("Successfully deployed and configured")

	framework.L.Info().Msg("Generate and propose secure mint job...")
	jobSpec := createSecureMintWorkflowJobSpec(t, &s, solChain)
	proposeSecureMintJob(t, fullCldEnvOutput.Environment.Offchain, jobSpec)
	framework.L.Info().Msgf("Secure mint job is succesfully posted. Job spec:\n %v", jobSpec)
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

const secureMintWorkflowTemplate = `
name: "{{.WorkflowName}}"
owner: "{{.WorkflowOwner}}"
triggers:
  - id: "securemint-trigger@1.0.0" #currently mocked 
    config:
      maxFrequencyMs: 5000
actions:
  - id: "{{.DeriveID}}"
    ref: "solana_data_feeds_cache_accounts"
    inputs: 
      trigger_output: $(trigger.outputs) # don't really need it, but without inputs can't pass wf validation
    config:
      Receiver: "{{.DFCacheAddr}}"
      State: "{{.CacheStateID}}"
      FeedIDs: ["{{.FeedID}}"]
consensus:
  - id: "offchain_reporting@1.0.0"
    ref: "secure-mint-consensus"
    inputs:
      observations:
        - "$(trigger.outputs)"
      solana:
        account_context:
          - "$(solana_data_feeds_cache_accounts.outputs)"
    config:
      report_id: "0003"  
      key_id: "solana"
      aggregation_method: "secure_mint" 
      aggregation_config:
        targetChainSelector: "{{.ChainSelector}}" # CHAIN_ID_FOR_WRITE_TARGET: NEW Param, to match write target
        solana:
          account_context: "$(inputs.solana.account_context)"
      encoder: "Borsh"
      encoder_config:
        abi: "(bytes16 DataID, uint32 Timestamp, uint224 Answer)[] Reports"

targets:
  - id: "{{.SolanaWriteTargetID}}"
    inputs:
      signed_report: $(secure-mint-consensus.outputs)
    config:
      address: "{{.DFCacheAddr}}"
      params: ["$(report)"]
      abi: "receive(report bytes)"
      deltaStage: 1s
      schedule: oneAtATime
`

func createSecureMintWorkflowJobSpec(t *testing.T, s *setup, solChain *cre.WrappedBlockchainOutput) string {
	tmpl, err := template.New("secureMintWorkflow").Parse(secureMintWorkflowTemplate)
	require.NoError(t, err)
	chainID, err := solChain.SolClient.GetGenesisHash(context.Background())
	require.NoError(t, err, "failed to receive genesis hash")
	deriveCapabilityID := writetarget.GenerateDeriveRemainingName(chainID.String())
	writeCapabilityID := writetarget.GenerateWriteTargetName(chainID.String())
	data := map[string]any{
		"WorkflowName":        s.WFName,
		"WorkflowOwner":       common.BytesToAddress(s.WFOwner[:]).Hex(),
		"ChainSelector":       s.Selector,
		"DFCacheAddr":         s.CacheProgramID.String(),
		"CacheStateID":        s.CacheState.String(),
		"SolanaWriteTargetID": writeCapabilityID,
		"DeriveID":            deriveCapabilityID,
		"FeedID":              s.FeedID,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	require.NoError(t, err)

	spec := buf.String()
	workflowJobSpec := testspecs.GenerateWorkflowJobSpec(t, spec)
	return fmt.Sprintf(
		`
		externaljobid   		 	=  "123e4567-e89b-12d3-a456-426655440002"
		%s
	`, workflowJobSpec.Toml())
}

func proposeSecureMintJob(t *testing.T, offchain offchain.Client, jobSpec string) {
	nodes, err := offchain.ListNodes(t.Context(), &node.ListNodesRequest{})
	require.NoError(t, err, "failed to get list nodes")
	var specs cre.DonJobs
	for _, n := range nodes.GetNodes() {
		if strings.Contains(n.Name, "bootstrap") {
			continue
		}
		specs = append(specs, &job.ProposeJobRequest{
			Spec:   jobSpec,
			NodeId: n.Id,
		})

	}
	err = jobs.Create(t.Context(), offchain, specs)
	if err != nil && strings.Contains(err.Error(), "is already approved") {
		return
	}
	require.NoError(t, err, "failed to propose jobs")
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
