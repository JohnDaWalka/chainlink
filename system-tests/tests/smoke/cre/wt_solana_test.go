package cre

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	ocr3types "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/report"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	sol "github.com/smartcontractkit/chainlink-solana/pkg/solana"
	writetarget "github.com/smartcontractkit/chainlink-solana/pkg/solana/write_target"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	df_cs "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_solana "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	consensuscap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/consensus"
	mockcap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/mock"
	gatewayconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config/gateway"
	solwriterconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config/writesolana"
	creconsensus "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/consensus"
	cregateway "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/gateway"
	cremock "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/mock"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	mock_capability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock"
	mockcapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock/pb"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/stretchr/testify/require"
)

var SinglePoRDonCapabilitiesFlagsSolana = []string{cre.OCR3Capability, cre.WriteSolanaCapability, cre.MockCapability}

func Test_WT_solana_with_mocked_capabilities(t *testing.T) {
	configErr := setConfigurationIfMissing("environment-one-don-multichain-solana-ci.toml", "workflow")
	require.NoError(t, configErr, "failed to set CTF config")
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[environment.Config](t)
	require.NoError(t, err, "couldn't load test config")
	//validateEnvVars(t)
	require.Len(t, in.NodeSets, 1, "expected 1 node set in the test config")
	// Assign all capabilities to the single node set
	var solChains []string
	for _, chain := range in.Blockchains {
		if chain.Type == "solana" {
			solChains = append(solChains, chain.ChainID)
		}
	}
	mustSetCapabilitiesFn := func(input []*ns.Input) []*cre.CapabilitiesAwareNodeSet {
		return []*cre.CapabilitiesAwareNodeSet{
			{
				Input:              input[0],
				SupportedSolChains: solChains,
				Capabilities:       SinglePoRDonCapabilitiesFlagsSolana,
				DONTypes:           []string{cre.WorkflowDON, cre.CapabilitiesDON, cre.GatewayDON},
				BootstrapNodeIndex: 0, // not required, but set to make the configuration explicit
				GatewayNodeIndex:   0, // not required, but set to make the configuration explicit
			},
		}
	}

	feedIDs := make([]string, 0, len(in.WorkflowConfigs))
	for _, wc := range in.WorkflowConfigs {
		feedIDs = append(feedIDs, wc.FeedID)
	}

	capabilityFactoryFns := []cre.DONCapabilityWithConfigFactoryFn{
		consensuscap.OCR3CapabilityFactoryFn,
		mockcap.CapabilityFactoryFn,
	}

	setupOut := setupWTTestEnvironment(
		t,
		testLogger,
		in,
		mustSetCapabilitiesFn,
		capabilityFactoryFns,
	)

	kb := make([]ocr2key.KeyBundle, 0)
	for _, don := range setupOut.DonTopology.DonsWithMetadata {
		if flags.HasFlag(don.Flags, cre.OCR3Capability) {
			for i, n := range don.DON.Nodes {
				if i == 0 {
					continue
				}
				key, err2 := n.ExportOCR2Keys(n.Ocr2KeyBundleID)
				if err2 != nil {
					testLogger.Error().Err(err).Msgf("failed to export ocr2 keys for node %d", i)
					continue
				}

				b, err2 := json.Marshal(key)
				require.NoError(t, err2, "failed to marshal ocr2 key")
				kk, err3 := ocr2key.FromEncryptedJSON(b, nodeclient.ChainlinkKeyPassword)
				require.NoError(t, err3, "could not decrypt OCR2 key json")
				kb = append(kb, kk)
				fmt.Println("setup key ", common.BytesToAddress(kk.PublicKey()))
			}
		}

	}

	fmt.Println("ocr2 keys", kb)

	mocksClient := mock_capability.NewMockCapabilityController(testLogger)
	mockClientsAddress := make([]string, 0)
	if in.Infra.Type == "docker" {
		for _, nodeSet := range in.NodeSets {
			if nodeSet.Name == "workflow" {
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
		}
	}

	require.NoError(t, mocksClient.ConnectAll(mockClientsAddress, true, true), "could not connect to mock capabilities")
	fmt.Println("cap name", setupOut.WriteCap)
	fmt.Println("forwarder address", setupOut.ForwarderAddress, "forwarder state", setupOut.ForwarderState)
	fmt.Println("cache address", setupOut.CacheAddress, "cache state", setupOut.CacheState)

	writer := newWriter(mocksClient, kb, setupOut)

	err = writer.Call()
	require.NoError(t, err, "failed to call write solana")
}

type writer struct {
	mocksClient *mock_capability.Controller

	wfID    string //
	wfOwner string
	wfName  string

	seqNr    uint64
	feedID   string
	reportID uint8

	cacheState   string
	cacheAddress string

	writeCapability           string
	deriveRemainingCapability string

	keys []ocr2key.KeyBundle
}

func newWriter(mclient *mock_capability.Controller, keys []ocr2key.KeyBundle, setup *setupWTOutput) *writer {
	return &writer{
		mocksClient: mclient,
		keys:        keys,
		wfOwner:     setup.WFOwner,
		wfName:      setup.WFName,
		wfID:        "5dbe5f217ff07d6b1dddb43519fe7bf13ccb10b540578fafdbea86b508abbd71",

		seqNr:    1,
		feedID:   setup.FeedID,
		reportID: 1,

		cacheState:   setup.CacheState,
		cacheAddress: setup.CacheAddress,

		writeCapability:           setup.WriteCap,
		deriveRemainingCapability: setup.DeriveRemaining,
	}
}

func (w *writer) Call() error {
	remainings, err := w.getRemainings()
	if err != nil {
		return fmt.Errorf("failed to get remaining accounts: %w", err)
	}

	metadata, err := w.createWorkflowMetadata()
	if err != nil {
		return fmt.Errorf("failed to create workflow metadata: %w", err)
	}

	repCtx := report.GenerateReportContext(w.seqNr, [32]byte{1})

	encodedReport, err := w.createEncodedReport(metadata, remainings)
	if err != nil {
		return fmt.Errorf("failed to create encoded report: %w", err)
	}

	sigs, err := w.generateSignatures(encodedReport, repCtx)
	if err != nil {
		return fmt.Errorf("failed to generate signatures: %w", err)
	}

	err = w.executeRequest(remainings, metadata, repCtx, encodedReport, sigs)
	if err != nil {
		return fmt.Errorf("failed to execute write capability: %w", err)
	}

	return nil
}

func (w *writer) executeRequest(remainings solana.AccountMetaSlice, metadata *pb.Metadata, repCtx []byte, encReport []byte, sigs [][]byte) error {
	config, input, err := w.createRequestInputs(remainings, repCtx, encReport, sigs)
	if err != nil {
		return err
	}

	req := &pb.ExecutableRequest{
		ID:              w.writeCapability,
		CapabilityType:  4,
		RequestMetadata: metadata,
		Config:          config,
		Inputs:          input,
	}

	return w.mocksClient.Execute(context.TODO(), req)
}

func (w *writer) createRequestInputs(remainings solana.AccountMetaSlice, repCtx []byte, encReport []byte, sigs [][]byte) ([]byte, []byte, error) {
	inputs, err := values.NewMap(map[string]any{
		"signed_report": map[string]any{
			"report":     encReport,
			"signatures": sigs,
			"context":    repCtx,
			"id":         [2]byte{0, w.reportID},
		},
		"remaining_accounts": remainings,
	})
	if err != nil {
		return nil, nil, err
	}
	retInputs, err := mock_capability.MapToBytes(inputs)
	if err != nil {
		return nil, nil, err
	}

	config, err := values.NewMap(map[string]any{
		"Address": w.cacheAddress,
	})
	if err != nil {
		return nil, nil, err
	}

	retConfig, err := mock_capability.MapToBytes(config)
	if err != nil {
		return nil, nil, err
	}

	return retConfig, retInputs, nil
}

func (w *writer) generateSignatures(report []byte, reportCtx []byte) ([][]byte, error) {
	var ret [][]byte
	sigData := append(report, reportCtx...)
	for _, k := range w.keys {
		hashed := sha256.Sum256(sigData)
		sig, err := k.SignBlob(hashed[:])
		if err != nil {
			return nil, err
		}
		ret = append(ret, sig)
	}

	return ret, nil
}

type decimalReport struct {
	Timestamp uint32
	Answer    *big.Int
	DataID    [16]byte
}

func (w *writer) getRemainings() (solana.AccountMetaSlice, error) {
	val, err := values.WrapMap(struct {
		State    string // State pubkey of df cache
		Receiver string // df cache programID
		FeedIDs  []string
	}{
		w.cacheState,
		w.cacheAddress,
		[]string{w.feedID},
	})
	if err != nil {
		return nil, err
	}
	cfg, err := mockcapability.MapToBytes(val)
	if err != nil {
		return nil, err
	}
	fmt.Println("remaining ID", w.deriveRemainingCapability)
	// get remaining accounts
	ret, err := w.mocksClient.Nodes[1].API.Execute(context.TODO(), &pb.ExecutableRequest{
		ID:             w.deriveRemainingCapability,
		CapabilityType: 4,
		Config:         cfg,
		Inputs:         []byte{},
		RequestMetadata: &pb.Metadata{
			WorkflowOwner: w.wfOwner,
			WorkflowName:  w.wfName,
		},
	})
	if err != nil {
		return nil, err
	}
	val, err = mockcapability.BytesToMap(ret.Value)
	if err != nil {
		return nil, err
	}

	var remainings solana.AccountMetaSlice

	err = val.Underlying["remaining_accounts"].UnwrapTo(&remainings)
	if err != nil {
		return nil, err
	}

	return remainings, nil
}

func (w *writer) createEncodedReport(m *pb.Metadata, remainings solana.AccountMetaSlice) ([]byte, error) {
	encoder, err := w.createEncoder()
	if err != nil {
		return nil, fmt.Errorf("failed to create encoder: %w", err)
	}
	ts := uint32(time.Now().Unix()) //nolint:gosec // disable G115

	reports := w.generateReports(ts)

	var buff []byte
	for _, acc := range remainings {
		buff = append(buff, acc.PublicKey.Bytes()...)
	}

	accsHash := sha256.Sum256(buff)
	fakeReport := w.createFakeReport(reports, m, ts, accsHash)
	wrappedFakeReport, err := values.NewMap(fakeReport)
	if err != nil {
		return nil, fmt.Errorf("failed to wrap fake report: %w", err)
	}

	return encoder.Encode(context.TODO(), *wrappedFakeReport)
}

func (w *writer) createFakeReport(reports []any, m *pb.Metadata, ts uint32, accHash [32]byte) map[string]any {
	meta := ocr3types.Metadata{
		Version:          1,
		ExecutionID:      m.WorkflowExecutionID,
		Timestamp:        ts,
		DONID:            1,
		DONConfigVersion: 1,
		WorkflowID:       m.WorkflowID,
		WorkflowName:     m.WorkflowName,
		WorkflowOwner:    m.WorkflowOwner,
		ReportID:         fmt.Sprintf("%04x", w.reportID),
	}

	return map[string]any{
		"account_ctx_hash":          accHash,
		"payload":                   reports,
		ocr3types.MetadataFieldName: meta,
	}

}

func (w *writer) generateReports(ts uint32) []any {
	var reports []any

	dataID, _ := new(big.Int).SetString(w.feedID, 0)
	var data [16]byte
	copy(data[:], dataID.Bytes())

	reports = append(reports, map[string]any{
		"timestamp": ts,
		"answer":    big.NewInt(12),
		"dataId":    data,
	})

	return reports
}

func (w *writer) createEncoder() (ocr3types.Encoder, error) {
	encoderConfig := map[string]any{
		"report_schema": `{
      "kind": "struct",
      "fields": [
        { "name": "payload", "type": { "vec": { "defined": "DecimalReport" } } }
      ]
    }`,
		"defined_types": `
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
    ]`,
	}
	wrappedEncoderConfig, err := values.NewMap(encoderConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create encoder config: %w", err)
	}

	encoder, err := sol.NewEncoder(wrappedEncoderConfig)
	if err != nil {
		return nil, err
	}

	return encoder, nil
}

func (w *writer) createWorkflowMetadata() (*pb.Metadata, error) {
	executionID, err := workflows.EncodeExecutionID("someID", uuid.NewString())
	if err != nil {
		return nil, err
	}

	return &pb.Metadata{
		WorkflowID:    w.wfID,
		WorkflowOwner: w.wfOwner,

		WorkflowName:        w.wfName, // already has correct format
		WorkflowExecutionID: executionID,

		WorkflowDonID:            1,
		WorkflowDonConfigVersion: 1,
		DecodedWorkflowName:      w.wfName,
	}, nil
}

func convertToHashedWorkflowName(input string) string {
	// Create SHA256 hash
	hash := sha256.New()
	hash.Write([]byte(input))

	// Get the hex string of the hash
	hashHex := hex.EncodeToString(hash.Sum(nil))

	return hex.EncodeToString([]byte(hashHex[:10]))
}

type setupWTOutput struct {
	WriteCap        string
	DeriveRemaining string
	SolChainID      string

	ForwarderAddress string
	ForwarderState   string

	CacheAddress string
	CacheState   string

	WFName      string
	WFOwner     string
	FeedID      string
	DonTopology *cre.DonTopology
}

func setupWTTestEnvironment(
	t *testing.T,
	testLogger zerolog.Logger,
	in *environment.Config,
	mustSetCapabilitiesFn func(input []*ns.Input) []*cre.CapabilitiesAwareNodeSet,
	capabilityFactoryFns []func([]string) []keystone_changeset.DONCapabilityWithConfig,
) *setupWTOutput {
	extraAllowedGatewayPorts := []int{}

	customBinariesPaths := map[string]string{}
	//containerPath, pathErr := capabilities.DefaultContainerDirectory(in.Infra.Type)
	//require.NoError(t, pathErr, "failed to get default container directory")
	//var mockBinaryPathInTheContainer string
	if in.ExtraCapabilities.MockCapapilityBinaryPath != "" {
		// where cron binary is located in the container
		//	mockBinaryPathInTheContainer = filepath.Join(containerPath, filepath.Base(in.DependenciesConfig.MockCapapilityBinaryPath))
		// where cron binary is located on the host
		customBinariesPaths[cre.MockCapability] = in.ExtraCapabilities.MockCapapilityBinaryPath
	}

	t.Log("customBinariesPaths", customBinariesPaths)
	firstBlockchain := in.Blockchains[0]

	chainIDInt, err := strconv.Atoi(firstBlockchain.ChainID)
	require.NoError(t, err, "failed to convert chain ID to int")
	chainIDUint64 := libc.MustSafeUint64(int64(chainIDInt))
	cfg := solwriterconfig.GetGenerateConfig()

	universalSetupInput := creenv.SetupInput{
		CapabilitiesAwareNodeSets:            mustSetCapabilitiesFn(in.NodeSets),
		CapabilitiesContractFactoryFunctions: capabilityFactoryFns,
		BlockchainsInput:                     in.Blockchains,
		JdInput:                              *in.JD,
		InfraInput:                           *in.Infra,
		CustomBinariesPaths:                  customBinariesPaths,
		JobSpecFactoryFunctions: []cre.JobSpecFactoryFn{
			creconsensus.ConsensusJobSpecFactoryFn(chainIDUint64),
			cregateway.GatewayJobSpecFactoryFn(extraAllowedGatewayPorts, []string{}, []string{"0.0.0.0/0"}),
			cremock.MockJobSpecFactoryFn("mock"),
		},
		ConfigFactoryFunctions: []cre.ConfigFactoryFn{
			gatewayconfig.GenerateConfigFn,
			cfg,
		},
	}

	universalSetupOutput, setupErr := creenv.SetupTestEnvironment(t.Context(), testLogger, cldlogger.NewSingleFileLogger(t), universalSetupInput)
	require.NoError(t, setupErr, "failed to setup test environment")

	out := &setupWTOutput{}
	out.DonTopology = universalSetupOutput.DonTopology
	wfDescription := [][32]uint8{{2, 3, 4}}
	wfOwner := [20]byte{222, 173, 190}
	name := convertToHashedWorkflowName("wf_name")
	out.WFName = name
	bname, _ := hex.DecodeString(name)
	var nname [10]uint8
	copy(nname[:], bname)
	wfName := [][10]uint8{nname}

	out.WFOwner = hex.EncodeToString(wfOwner[:])
	out.FeedID = in.WorkflowConfigs[0].FeedID
	for _, bo := range universalSetupOutput.BlockchainOutput {
		if bo.ReadOnly {
			continue
		}

		if bo.SolChain != nil {
			chainID, err := bo.SolClient.GetGenesisHash(context.Background())
			require.NoError(t, err, "failed to get genesis hash")
			out.WriteCap = writetarget.GenerateWriteTargetName(chainID.String())
			out.DeriveRemaining = writetarget.GenerateDeriveRemainingName(chainID.String())
			forwarder, err := universalSetupOutput.CldEnvironment.DataStore.Addresses().Get(datastore.NewAddressRefKey(
				bo.SolChain.ChainSelector,
				ks_solana.ForwarderContract,
				semver.MustParse("1.0.0"),
				"test-forwarder",
			))
			require.NoError(t, err, "forwarder not found")
			forwarderState, err := universalSetupOutput.CldEnvironment.DataStore.Addresses().Get(datastore.NewAddressRefKey(
				bo.SolChain.ChainSelector,
				ks_solana.ForwarderState,
				semver.MustParse("1.0.0"),
				"test-forwarder",
			))
			out.ForwarderAddress = forwarder.Address
			out.ForwarderState = forwarderState.Address

			//df cache
			dfQualifier := "df-cache-qualifier"
			dfDeployOut, err := commonchangeset.RunChangeset(df_cs.DeployCache{}, *universalSetupOutput.CldEnvironment, &df_cs.DeployCacheRequest{
				ChainSel:   bo.SolChain.ChainSelector,
				Qualifier:  dfQualifier,
				Version:    "1.0.0",
				FeedAdmins: []solana.PublicKey{bo.SolChain.PrivateKey.PublicKey()},
			})
			require.NoError(t, err, "failed to deploy df cache")
			cacheID, err := dfDeployOut.DataStore.Addresses().Get(
				datastore.NewAddressRefKey(bo.SolChain.ChainSelector, df_cs.CacheContract, semver.MustParse("1.0.0"), dfQualifier))
			require.NoError(t, err, "df cache address not found")
			out.CacheAddress = cacheID.Address

			cacheState, err := dfDeployOut.DataStore.Addresses().Get(
				datastore.NewAddressRefKey(bo.SolChain.ChainSelector, df_cs.CacheState, semver.MustParse("1.0.0"), dfQualifier))
			require.NoError(t, err, "df cache state not found")

			out.CacheAddress = cacheID.Address
			out.CacheState = cacheState.Address
			ds := datastore.NewMemoryDataStore()
			ds.Merge(dfDeployOut.DataStore.Seal())
			ds.Merge(universalSetupOutput.CldEnvironment.DataStore)
			universalSetupOutput.CldEnvironment.DataStore = ds.Seal()
			feedIDin, ok := new(big.Int).SetString(in.WorkflowConfigs[0].FeedID, 0)
			require.True(t, ok, "invalid feedID")
			require.LessOrEqual(t, feedIDin.BitLen(), 128, "invalid feedID len")
			var feedID [16]uint8
			copy(feedID[:], feedIDin.Bytes())
			require.NoError(t, err, "failed to decode FeedID")
			_, err = commonchangeset.RunChangeset(df_cs.InitCacheDecimalReport{}, *universalSetupOutput.CldEnvironment,
				&df_cs.InitCacheDecimalReportRequest{
					ChainSel:  bo.SolChain.ChainSelector,
					Qualifier: dfQualifier,
					Version:   "1.0.0",
					FeedAdmin: bo.SolChain.PrivateKey.PublicKey(),
					DataIDs:   [][16]uint8{feedID},
				},
			)
			require.NoError(t, err, "failed to init decimal report")

			_, err = commonchangeset.RunChangeset(df_cs.ConfigureCacheDecimalReport{}, *universalSetupOutput.CldEnvironment,
				&df_cs.ConfigureCacheDecimalReportRequest{
					ChainSel:  bo.SolChain.ChainSelector,
					Qualifier: dfQualifier,
					Version:   "1.0.0",
					SenderList: []df_cs.Sender{
						{
							ProgramID: solana.MustPublicKeyFromBase58(forwarder.Address),
							StateID:   solana.MustPublicKeyFromBase58(forwarderState.Address),
						},
					},
					FeedAdmin:            bo.SolChain.PrivateKey.PublicKey(),
					DataIDs:              [][16]uint8{feedID},
					AllowedWorkflowOwner: [][20]uint8{wfOwner},
					AllowedWorkflowName:  wfName,
					Descriptions:         wfDescription,
				})
			require.NoError(t, err, "failed to configure decimal report")
		}
	}

	return out
}
