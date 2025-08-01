package solreader

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	sel "github.com/smartcontractkit/chain-selectors"

	typepkgmock "github.com/smartcontractkit/chainlink-ccip/mocks/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pkg/contractreader"
	"github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/chainreader"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/logpoller"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/logpoller/types"

	soltest "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/testutils"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_event_emitter"

	solconf "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/configs/solana"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
)

// How to run this test:
//   cd ./core/capabilities/ccip/ccip_integration_tests/usdcreader/solreader
//   go test -v -run=Test_USDCReader_MessageHashes ./...

const (
	// obtained from Anchor.toml in chainlink-ccip
	SolanaTestEventEmitterProgramID = "EGfB7iiotGoDVpQvByFD8AD11BhTpc9WMCyUL5q64smj"

	// points to a Git release asset with the test_event_emitter program in chainlink-ccip
	SolanaCCIPProgramsGitShaVersion = "47fa4bc350ce"

	// arbitrary timeout - can be anything
	DefaultLogPollerRequestTimeout = time.Second * 5
)

func Test_USDCReader_MessageHashes(t *testing.T) {
	ctx := testutils.Context(t)

	lgr, err := logger.New()
	require.NoError(t, err)

	solBlockchain := ccipocr3.ChainSelector(sel.SOLANA_MAINNET.Selector)
	solDomainCCTP := reader.CCTPDestDomains[uint64(solBlockchain)]
	mockAddrCodec := typepkgmock.NewMockAddressCodec(t)

	cfg, err := solconf.DestContractReaderConfig()
	require.NoError(t, err)

	vars := setup(ctx, t, solBlockchain, cfg, false)
	t.Log("finished setting up environment - proceeding with test")

	mockAddrCodec.
		On("AddressBytesToString", mock.Anything, mock.Anything).
		Return(
			func(addr ccipocr3.UnknownAddress, _ ccipocr3.ChainSelector) string {
				return solana.PublicKeyFromBytes(addr).String()
			},
		).
		Maybe()

	mockAddrCodec.
		On("AddressStringToBytes", mock.Anything, mock.Anything).
		Return(
			func(addr string, _ ccipocr3.ChainSelector) (ccipocr3.UnknownAddress, error) {
				if pubKey, err := solana.PublicKeyFromBase58(addr); err != nil {
					return ccipocr3.UnknownAddress{}, err
				} else {
					return pubKey.Bytes(), nil
				}
			},
		).
		Maybe()

	t.Log("creating USDC message reader for Solana")
	usdcReader, err := reader.NewUSDCMessageReader(
		ctx,
		lgr,
		map[ccipocr3.ChainSelector]pluginconfig.USDCCCTPTokenConfig{
			solBlockchain: {
				SourceMessageTransmitterAddr: vars.contract.String(),
			},
		},
		map[ccipocr3.ChainSelector]contractreader.Extended{
			solBlockchain: vars.reader,
		},
		mockAddrCodec,
	)
	require.NoError(t, err)

	t.Log("emitting avalanche events")
	avalancheBlockchain := ccipocr3.ChainSelector(sel.AVALANCHE_MAINNET.Selector)
	avalancheChainSlctr := uint64(avalancheBlockchain)
	emitMessageSent(t, vars, solDomainCCTP, avalancheChainSlctr, 11)
	emitMessageSent(t, vars, solDomainCCTP, avalancheChainSlctr, 21)
	emitMessageSent(t, vars, solDomainCCTP, avalancheChainSlctr, 31)
	emitMessageSent(t, vars, solDomainCCTP, avalancheChainSlctr, 41)
	emitMessageSent(t, vars, solDomainCCTP, avalancheChainSlctr, 51)

	t.Log("emitting polygon events")
	polygonBlockchain := ccipocr3.ChainSelector(sel.POLYGON_MAINNET.Selector)
	polygonChainSlctr := uint64(polygonBlockchain)
	emitMessageSent(t, vars, solDomainCCTP, polygonChainSlctr, 31)
	emitMessageSent(t, vars, solDomainCCTP, polygonChainSlctr, 41)

	// Replicating comment from EVM side in case it is relevant:
	//   Need to replay as sometimes the logs are not picked up by the log poller (?)
	//   Maybe another situation where chain reader doesn't register filters as expected.
	t.Log("requesting block replay")
	vars.lp.Replay(1)

	// Wait for replay to finish
	t.Log("waiting for block replay to finish")
	for true {
		time.Sleep(time.Second)
		switch vars.lp.ReplayStatus() {
		case types.ReplayStatusNoRequest:
			t.Fatal("replay request was not found")
			continue
		case types.ReplayStatusComplete:
			t.Log("replay complete")
			break
		}
	}

	// Define test cases
	tt := []struct {
		name           string
		tokens         map[reader.MessageTokenID]ccipocr3.RampTokenAmount
		sourceChain    ccipocr3.ChainSelector
		destChain      ccipocr3.ChainSelector
		expectedMsgIDs []reader.MessageTokenID
	}{
		{
			name:           "empty messages should return empty response",
			tokens:         map[reader.MessageTokenID]ccipocr3.RampTokenAmount{},
			sourceChain:    solBlockchain,
			destChain:      avalancheBlockchain,
			expectedMsgIDs: []reader.MessageTokenID{},
		},
		{
			name: "single token message",
			tokens: map[reader.MessageTokenID]ccipocr3.RampTokenAmount{
				reader.NewMessageTokenID(1, 1): {
					ExtraData: reader.NewSourceTokenDataPayload(11, solDomainCCTP).ToBytes(),
				},
			},
			sourceChain:    solBlockchain,
			destChain:      avalancheBlockchain,
			expectedMsgIDs: []reader.MessageTokenID{reader.NewMessageTokenID(1, 1)},
		},
		{
			name: "single token message but different chain",
			tokens: map[reader.MessageTokenID]ccipocr3.RampTokenAmount{
				reader.NewMessageTokenID(1, 2): {
					ExtraData: reader.NewSourceTokenDataPayload(31, solDomainCCTP).ToBytes(),
				},
			},
			sourceChain:    solBlockchain,
			destChain:      polygonBlockchain,
			expectedMsgIDs: []reader.MessageTokenID{reader.NewMessageTokenID(1, 2)},
		},
		{
			name: "message without matching nonce",
			tokens: map[reader.MessageTokenID]ccipocr3.RampTokenAmount{
				reader.NewMessageTokenID(1, 1): {
					ExtraData: reader.NewSourceTokenDataPayload(1234, solDomainCCTP).ToBytes(),
				},
			},
			sourceChain:    solBlockchain,
			destChain:      avalancheBlockchain,
			expectedMsgIDs: []reader.MessageTokenID{},
		},
		{
			name: "message without matching source domain",
			tokens: map[reader.MessageTokenID]ccipocr3.RampTokenAmount{
				reader.NewMessageTokenID(1, 1): {
					ExtraData: reader.NewSourceTokenDataPayload(11, 12910).ToBytes(),
				},
			},
			sourceChain:    solBlockchain,
			destChain:      avalancheBlockchain,
			expectedMsgIDs: []reader.MessageTokenID{},
		},
		{
			name: "message with multiple tokens",
			tokens: map[reader.MessageTokenID]ccipocr3.RampTokenAmount{
				reader.NewMessageTokenID(1, 1): {
					ExtraData: reader.NewSourceTokenDataPayload(11, solDomainCCTP).ToBytes(),
				},
				reader.NewMessageTokenID(1, 2): {
					ExtraData: reader.NewSourceTokenDataPayload(21, solDomainCCTP).ToBytes(),
				},
			},
			sourceChain: solBlockchain,
			destChain:   avalancheBlockchain,
			expectedMsgIDs: []reader.MessageTokenID{
				reader.NewMessageTokenID(1, 1),
				reader.NewMessageTokenID(1, 2),
			},
		},
		{
			name: "message with multiple tokens, one without matching nonce",
			tokens: map[reader.MessageTokenID]ccipocr3.RampTokenAmount{
				reader.NewMessageTokenID(1, 1): {
					ExtraData: reader.NewSourceTokenDataPayload(11, solDomainCCTP).ToBytes(),
				},
				reader.NewMessageTokenID(1, 2): {
					ExtraData: reader.NewSourceTokenDataPayload(12, solDomainCCTP).ToBytes(),
				},
				reader.NewMessageTokenID(1, 3): {
					ExtraData: reader.NewSourceTokenDataPayload(31, solDomainCCTP).ToBytes(),
				},
			},
			sourceChain: solBlockchain,
			destChain:   avalancheBlockchain,
			expectedMsgIDs: []reader.MessageTokenID{
				reader.NewMessageTokenID(1, 1),
				reader.NewMessageTokenID(1, 3),
			},
		},
		{
			name: "not finalized events are not returned",
			tokens: map[reader.MessageTokenID]ccipocr3.RampTokenAmount{
				reader.NewMessageTokenID(1, 5): {
					ExtraData: reader.NewSourceTokenDataPayload(51, solDomainCCTP).ToBytes(),
				},
			},
			sourceChain:    solBlockchain,
			destChain:      avalancheBlockchain,
			expectedMsgIDs: []reader.MessageTokenID{},
		},
	}

	// Run test cases
	t.Log("running test cases")
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Skip()

			hashes, queryErr := usdcReader.MessagesByTokenID(ctx, tc.sourceChain, tc.destChain, tc.tokens)
			require.NoError(t, queryErr)

			require.Equal(t, len(tc.expectedMsgIDs), len(hashes))
			for _, msgID := range tc.expectedMsgIDs {
				_, ok := hashes[msgID]
				require.True(t, ok)
			}
		})
	}
}

func emitMessageSent(t *testing.T, testEnv *testvars, source uint32, dest uint64, nonce uint64) {
	randSender, err := solana.NewRandomPrivateKey()
	randPrvKey, err := solana.NewRandomPrivateKey()

	args := test_event_emitter.CcipCctpMessageSentEventArgs{
		OriginalSender:      randSender.PublicKey(),
		EventAddress:        randPrvKey.PublicKey(),
		RemoteChainSelector: dest,
		MessageSentBytes:    nil,
		SourceDomain:        source,
		MsgTotalNonce:       nonce,
		CctpNonce:           nonce,
	}

	ix, err := test_event_emitter.NewEmitCcipCctpMsgSentInstruction(args, solana.SysVarClockPubkey).ValidateAndBuild()
	require.NoError(t, err)

	require.NotNil(t,
		soltest.SendAndConfirm(
			t.Context(),
			t,
			testEnv.rpc,
			[]solana.Instruction{ix},
			testEnv.auth,
			rpc.CommitmentConfirmed,
		),
	)
}

func setup(ctx context.Context, t *testing.T, readerChain ccipocr3.ChainSelector, cfg config.ContractReader, useHeavyDB bool) *testvars {
	// Create a logger for off chain services
	lggr, err := logger.New()
	require.NoError(t, err)

	// Parameterize database selection
	var db *sqlx.DB
	if useHeavyDB {
		_, db = heavyweight.FullTestDBV2(t, nil) // Use heavyweight database for benchmarks
	} else {
		db = pgtest.NewSqlxDB(t) // Use simple in-memory DB for tests
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Logf("failed to close db: %v", err)
		}
	})

	// Set program ID
	programID, err := solana.PublicKeyFromBase58(SolanaTestEventEmitterProgramID)
	require.NoError(t, err)
	test_event_emitter.SetProgramID(programID)

	// Create a deployer account
	prvKey, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)

	// Download compiled program file (test_event_emitter.so)
	programsDir := setupSolanaCCIPProgramArtifactsDir(ctx, t, SolanaCCIPProgramsGitShaVersion)
	programFile := filepath.Join(programsDir, "test_event_emitter.so")
	t.Logf("Deploying program at: %s", programFile)

	// Start a local solana node and deploy the test_event_emitter program to it
	tomlConfig := config.NewDefault()
	nodeRpcUrl := NewSolanaValidator().
		AddProgram(programID, programFile, prvKey.PublicKey()).
		WithTestDefaults(t).
		Run(t).
		RpcUrlString()

	// Setup logpoller client
	client, rpc, err := client.NewTestClient(nodeRpcUrl, tomlConfig, DefaultLogPollerRequestTimeout, lggr)
	require.NoError(t, err)

	// Fund the deployer account
	soltest.FundAccounts(ctx, []solana.PrivateKey{prvKey}, rpc, t)

	// Create ORM client and logpoller instance
	orm := logpoller.NewORM(readerChain.String(), db, lggr)
	lp := logpoller.New(
		logger.Sugared(lggr),
		orm,
		client,
		tomlConfig,
	)

	// Start log poller
	require.NoError(t, lp.Start(ctx))
	t.Cleanup(func() {
		if err := lp.Close(); err != nil {
			t.Logf("failed to close log poller: %v", err)
		}
	})

	// Create a contract reader instance
	wrp := &chainreader.RPCClientWrapper{AccountReader: rpc}
	cr, err := chainreader.NewContractReaderService(
		lggr,
		wrp,
		cfg,
		lp,
	)
	require.NoError(t, err)

	// Start chain reader
	require.NoError(t, cr.Start(ctx))
	t.Cleanup(func() {
		if err := cr.Close(); err != nil {
			t.Logf("failed to close chain reader: %v", err)
		}
	})

	// Convert to the extended contract reader interface
	reader := contractreader.NewExtendedContractReader((contractreader.ContractReaderFacade)(cr))
	return &testvars{
		contract: programID,
		client:   client,
		reader:   reader,
		auth:     prvKey,
		rpc:      rpc,
		orm:      orm,
		db:       db,
		lp:       lp,
	}
}

type testvars struct {
	contract solana.PublicKey
	client   *client.Client
	reader   contractreader.Extended
	auth     solana.PrivateKey
	rpc      *rpc.Client
	orm      logpoller.ORM
	db       *sqlx.DB
	lp       *logpoller.Service
}
