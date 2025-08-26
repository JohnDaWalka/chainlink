package cre

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"

	texttmpl "text/template"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	writetarget "github.com/smartcontractkit/chainlink-solana/pkg/solana/write_target"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	df "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	df_sol "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	ks_sol "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	mock_capability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock/pb"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
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

	// configure contracts
	framework.L.Info().Msg("Deploy and configure data-feeds cache programs...")
	deployAndConfigureCache(t, &s, *fullCldEnvOutput.Environment, solChain)
	framework.L.Info().Msg("Successfully deployed and configured")

	// configure workflow
	framework.L.Info().Msg("Generate and propose secure mint job...")
	jobSpec := createSecureMintWorkflowJobSpec(t, &s, solChain)
	proposeSecureMintJob(t, fullCldEnvOutput.Environment.Offchain, jobSpec)
	framework.L.Info().Msgf("Secure mint job is succesfully posted. Job spec:\n %v", jobSpec)

	// trigger workflow
	trigger := createFakeTrigger(t, &s, fullCldEnvOutput.DonTopology)
	trigger.Call(t)

	// wait for price update
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

const reportSchema = `{
      "kind": "struct",
      "fields": [
        { "name": "payload", "type": { "vec": { "defined": "DecimalReport" } } }
      ]
    }`
const definedTypes = `
     [
      {
        "name":"DecimalReport",
         "type":{
          "kind":"struct",
          "fields":[
            { "name":"timestamp", "type":"u32" },
            { "name":"answer",    "type":"u128" },
            { "name": "dataId",   "type": {"array": ["u8",16]}}
          ]
        }
      }
    ]`

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
      encoder: "borsh"
      encoder_config:
        report_schema: |
          {
            "kind": "struct",
            "fields": [
              { "name": "payload", "type": { "vec": { "defined": "DecimalReport" } } }
            ]
          }
        defined_types: |
          [
            {
              "name": "DecimalReport",
              "type": {
                "kind": "struct",
                "fields": [
                  { "name": "timestamp", "type": "u32" },
                  { "name": "answer",    "type": "u128" },
                  { "name": "dataId",    "type": { "array": ["u8", 16] } }
                ]
              }
            }
          ]

targets:
  - id: "{{.SolanaWriteTargetID}}"
    inputs:
      signed_report: $(secure-mint-consensus.outputs)
      remaining_accounts: $(solana_data_feeds_cache_accounts.outputs)
    config:
      address: "{{.DFCacheAddr}}"
      params: ["$(report)"]
      deltaStage: 1s
      schedule: oneAtATime
`

func createSecureMintWorkflowJobSpec(t *testing.T, s *setup, solChain *cre.WrappedBlockchainOutput) string {
	tmpl, err := texttmpl.New("secureMintWorkflow").Parse(secureMintWorkflowTemplate)
	require.NoError(t, err)

	chainID, err := solChain.SolClient.GetGenesisHash(context.Background())
	require.NoError(t, err, "failed to receive genesis hash")

	deriveCapabilityID := writetarget.GenerateDeriveRemainingName(chainID.String())
	writeCapabilityID := writetarget.GenerateWriteTargetName(chainID.String())
	owner := hex.EncodeToString(s.WFOwner[:])
	data := map[string]any{
		"WorkflowName":        s.WFName,
		"WorkflowOwner":       "0x" + owner,
		"ChainSelector":       s.Selector,
		"DFCacheAddr":         s.CacheProgramID.String(),
		"CacheStateID":        s.CacheState.String(),
		"SolanaWriteTargetID": writeCapabilityID,
		"DeriveID":            deriveCapabilityID,
		"FeedID":              s.FeedID,
		"ReportSchema":        reportSchema,
		"DefinedTypes":        definedTypes,
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

type fakeTrigger struct {
	triggerCap *mock_capability.Controller
	setup      *setup
	triggerID  string
	keys       []ocr2key.KeyBundle
}

func (f *fakeTrigger) Call(t *testing.T) {
	outputs, err := f.createReport()
	require.NoError(t, err, "failed to create fake report")

	outputsBytes, err := mock_capability.MapToBytes(outputs)
	require.NoError(t, err, "failed to convert map to bytes")

	message := pb.SendTriggerEventRequest{
		TriggerID: f.triggerID,
		ID:        "fake_trigger",
		Outputs:   outputsBytes,
	}

	err = f.triggerCap.SendTrigger(t.Context(), &message)
	require.NoError(t, err, "failed to send trigger to workflow")
}

func (f *fakeTrigger) createReport() (*values.Map, error) {
	type secureMintReport struct {
		ConfigDigest ocr2types.ConfigDigest
		SeqNr        uint64
		Block        uint64
		Mintable     *big.Int
	}
	configDigest, _ := hex.DecodeString("000e8613ec1ad47912a636904823e77e92bf4f580fe6fe7f2e6eff395118623c")
	report := &secureMintReport{
		ConfigDigest: ocr2types.ConfigDigest(configDigest),
		SeqNr:        0,
		Block:        10,
		Mintable:     big.NewInt(15),
	}

	reportBytes, err := json.Marshal(report)
	if err != nil {
		return nil, err
	}

	ocr3Report := &ocr3types.ReportWithInfo[uint64]{
		Report: ocr2types.Report(reportBytes),
		Info:   f.setup.Selector,
	}

	jsonReport, err := json.Marshal(ocr3Report)
	if err != nil {
		return nil, err
	}

	var sigs []capabilities.OCRAttributedOnchainSignature
	for i, key := range f.keys {
		sig, err2 := key.Sign3(ocr2types.ConfigDigest(configDigest), 0, reportBytes)
		if err2 != nil {
			return nil, err2
		}
		sigs = append(sigs, capabilities.OCRAttributedOnchainSignature{
			Signer:    uint32(i), //nolint:gosec // G115 don't care in test code
			Signature: sig,
		})
	}

	event := &capabilities.OCRTriggerEvent{
		ConfigDigest: configDigest,
		SeqNr:        0,
		Report:       jsonReport,
		Sigs:         sigs,
	}

	outputs, err := event.ToMap()
	if err != nil {
		return nil, err
	}

	return outputs, nil
}

func createFakeTrigger(t *testing.T, s *setup, dons *cre.DonTopology) *fakeTrigger {
	client := createMockClient(t)
	keys := exportOcr2Keys(t, dons)
	require.NotEqual(t, len(keys), 0)
	framework.L.Info().Msg("Successfully exported ocr2 keys")

	return &fakeTrigger{
		triggerCap: client,
		keys:       keys,
		setup:      s,
		triggerID:  "securemint-trigger@1.0.0",
	}
}

func exportOcr2Keys(t *testing.T, dons *cre.DonTopology) []ocr2key.KeyBundle {
	kb := make([]ocr2key.KeyBundle, 0)
	for _, don := range dons.DonsWithMetadata {
		if flags.HasFlag(don.Flags, cre.MockCapability) {
			for _, n := range don.DON.Nodes {
				key, err2 := n.ExportOCR2Keys(n.Ocr2KeyBundleID)
				if err2 == nil {
					b, err2 := json.Marshal(key)
					require.NoError(t, err2, "could not marshal OCR2 key")
					kk, err3 := ocr2key.FromEncryptedJSON(b, nodeclient.ChainlinkKeyPassword)
					require.NoError(t, err3, "could not decrypt OCR2 key json")
					kb = append(kb, kk)
				} else {
					framework.L.Error().Msgf("Could not export OCR2 key: %s", err2)
				}
			}
		}
	}

	return kb
}

func createMockClient(t *testing.T) *mock_capability.Controller {
	in, err := framework.Load[envconfig.Config](nil)
	require.NoError(t, err, "couldn't load environment state")
	mockClientsAddress := make([]string, 0)
	for _, nodeSet := range in.NodeSets {
		for i, n := range nodeSet.NodeSpecs {
			if i == 0 {
				continue
			}
			if len(n.Node.CustomPorts) == 0 {
				panic("no custom port specified, mock capability running in kind must have a custom port in order to connect")
			}
			ports := strings.Split(n.Node.CustomPorts[0], ":")
			mockClientsAddress = append(mockClientsAddress, "127.0.0.1:"+ports[0])
		}
	}

	mocksClient := mock_capability.NewMockCapabilityController(framework.L)
	require.NoError(t, mocksClient.ConnectAll(mockClientsAddress, true, true), " failed to connect mock client")

	return mocksClient
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
