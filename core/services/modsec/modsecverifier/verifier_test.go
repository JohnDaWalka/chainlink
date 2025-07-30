package modsecverifier

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/verifier_events"
	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads/headstest"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/modsec/modsecstorage"
	"github.com/smartcontractkit/chainlink/v2/core/services/modsec/modsectypes"
	"github.com/test-go/testify/require"
	"go.uber.org/zap/zapcore"
)

const (
	sourceChainSelector = 1
	destChainSelector   = 2
)

type testUniverse struct {
	simBackend     *simulated.Backend
	simClient      client.Client
	auth           *bind.TransactOpts
	verifierEvents *verifier_events.VerifierEvents
	lp             logpoller.LogPoller
	verifier       *verifier
}

func (u *testUniverse) emitCCIPMessageSent(t *testing.T, destChainSelector uint64, sequenceNumber uint64, message verifier_events.InternalEVM2AnyCommitVerifierMessage) {
	_, err := u.verifierEvents.EmitCCIPMessageSent(
		u.auth,
		destChainSelector,
		sequenceNumber,
		message,
	)
	require.NoError(t, err)
	u.simBackend.Commit()
}

func setupSimulatedBackendAndAuth(t *testing.T) (*simulated.Backend, client.Client, *bind.TransactOpts, *verifier_events.VerifierEvents) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	blnc, ok := big.NewInt(0).SetString("999999999999999999999999999999999999", 10)
	require.True(t, ok)

	alloc := map[common.Address]ethtypes.Account{crypto.PubkeyToAddress(privateKey.PublicKey): {Balance: blnc}}
	simulatedBackend := simulated.NewBackend(alloc, simulated.WithBlockGasLimit(8000000))

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	require.NoError(t, err)
	auth.GasLimit = uint64(6000000)

	simClient := client.NewSimulatedBackendClient(t, simulatedBackend, big.NewInt(1337))

	// adjust time so the queries by the log poller make sense.
	head, err := simClient.HeadByNumber(t.Context(), nil)
	require.NoError(t, err)
	simulatedBackend.AdjustTime(time.Since(head.Timestamp) - time.Hour)
	simulatedBackend.Commit()

	addr, _, _, err := verifier_events.DeployVerifierEvents(auth, simClient)
	require.NoError(t, err)
	simulatedBackend.Commit()

	verifierEvents, err := verifier_events.NewVerifierEvents(addr, simClient)
	require.NoError(t, err)

	return simulatedBackend, simClient, auth, verifierEvents
}

func testSetup(t *testing.T) *testUniverse {
	simBackend, simClient, auth, verifierEvents := setupSimulatedBackendAndAuth(t)

	db := pgtest.NewSqlxDB(t)
	lpOpts := logpoller.Opts{
		PollPeriod:               time.Millisecond,
		FinalityDepth:            1,
		BackfillBatchSize:        10,
		RPCBatchSize:             10,
		KeepFinalizedBlocksDepth: 100000,
	}
	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(zapcore.ErrorLevel)

	headTracker := headstest.NewSimulatedHeadTracker(simClient, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
	lp := logpoller.NewLogPoller(
		logpoller.NewORM(big.NewInt(0).SetUint64(uint64(1337)), db, lggr),
		simClient,
		lggr,
		headTracker,
		lpOpts,
	)
	require.NoError(t, lp.Start(t.Context()))
	t.Cleanup(func() {
		require.NoError(t, lp.Close())
	})

	// mine some blocks
	for range 10 {
		simBackend.Commit()
	}

	key, err := ethkey.NewV2()
	require.NoError(t, err)

	v := &verifier{
		lggr:          logger.TestLogger(t),
		lp:            lp,
		eventSig:      verifier_events.VerifierEventsCCIPMessageSent{}.Topic().Hex(),
		onRampAddress: verifierEvents.Address().Hex(),
		storage:       modsecstorage.NewTestStorage(),
		parser:        modsectypes.NewEVMCCIPMessageParser(),
		signerKey:     key,
	}

	return &testUniverse{
		simBackend:     simBackend,
		simClient:      simClient,
		auth:           auth,
		verifierEvents: verifierEvents,
		lp:             lp,
		verifier:       v,
	}
}

func Test_registerFilterIfNeeded(t *testing.T) {
	universe := testSetup(t)

	filtersBefore := universe.lp.GetFilters()
	require.Len(t, filtersBefore, 0)

	require.NoError(t, universe.verifier.registerFilterIfNeeded(t.Context()))

	filtersAfter := universe.lp.GetFilters()
	require.Len(t, filtersAfter, 1)
	var fltr logpoller.Filter
	for _, filter := range filtersAfter {
		require.Equal(t, universe.verifier.eventSig, filter.EventSigs[0].Hex())
		require.Equal(t, universe.verifier.onRampAddress, filter.Addresses[0].Hex())
		fltr = filter
	}

	// re-register should not add a new filter
	require.NoError(t, universe.verifier.registerFilterIfNeeded(t.Context()))

	filtersAfter2 := universe.lp.GetFilters()
	require.Len(t, filtersAfter2, 1)
	for _, filter := range filtersAfter2 {
		require.Equal(t, fltr, filter)
	}
}

func Test_initializeLastProcessedBlock_allVerified(t *testing.T) {
	universe := testSetup(t)

	require.NoError(t, universe.verifier.registerFilterIfNeeded(t.Context()))

	for seqNr := range 10 {
		universe.emitCCIPMessageSent(t, destChainSelector, uint64(seqNr), verifier_events.InternalEVM2AnyCommitVerifierMessage{
			Header: verifier_events.InternalHeader{
				MessageId:           [32]byte{byte(seqNr)},
				SourceChainSelector: sourceChainSelector,
				DestChainSelector:   destChainSelector,
				SequenceNumber:      uint64(seqNr),
			},
			Sender:             common.HexToAddress("0x1234"),
			Data:               nil,
			Receiver:           common.LeftPadBytes(common.HexToAddress("0x5678").Bytes(), 32),
			DestChainExtraArgs: nil,
			VerifierExtraArgs:  nil,
			FeeToken:           common.HexToAddress("0x0"),
			FeeTokenAmount:     big.NewInt(13371337),
			FeeValueJuels:      big.NewInt(13371337),
			TokenAmounts:       nil,
			RequiredVerifiers:  nil,
		})
	}

	universe.simBackend.Commit()

	// Wait for log poller to index the logs.
	// TODO: should be programmatic / event driven instead.
	time.Sleep(5 * time.Second)

	// verify everything in storage
	for i := range 10 {
		messageID := [32]byte{byte(i)}
		universe.verifier.storage.Set(t.Context(), hexutil.Encode(messageID[:]), []byte("verified"))
	}

	latestBlock, err := universe.lp.LatestBlock(t.Context())
	require.NoError(t, err)

	lastProcessedBlock, err := universe.verifier.initializeLastProcessedBlock(t.Context())
	require.NoError(t, err)
	require.Equal(t, latestBlock.FinalizedBlockNumber, lastProcessedBlock)
}

func Test_initializeLastProcessedBlock_someVerified(t *testing.T) {
	universe := testSetup(t)

	require.NoError(t, universe.verifier.registerFilterIfNeeded(t.Context()))

	for seqNr := range 10 {
		universe.emitCCIPMessageSent(t, destChainSelector, uint64(seqNr), verifier_events.InternalEVM2AnyCommitVerifierMessage{
			Header: verifier_events.InternalHeader{
				MessageId:           [32]byte{byte(seqNr)},
				SourceChainSelector: sourceChainSelector,
				DestChainSelector:   destChainSelector,
				SequenceNumber:      uint64(seqNr),
			},
			Sender:             common.HexToAddress("0x1234"),
			Data:               nil,
			Receiver:           common.LeftPadBytes(common.HexToAddress("0x5678").Bytes(), 32),
			DestChainExtraArgs: nil,
			VerifierExtraArgs:  nil,
			FeeToken:           common.HexToAddress("0x0"),
			FeeTokenAmount:     big.NewInt(13371337),
			FeeValueJuels:      big.NewInt(13371337),
			TokenAmounts:       nil,
			RequiredVerifiers:  nil,
		})
	}

	universe.simBackend.Commit()

	// Wait for log poller to index the logs.
	// TODO: should be programmatic / event driven instead.
	time.Sleep(5 * time.Second)

	// put some stuff in storage
	for i := range 5 {
		messageID := [32]byte{byte(i)}
		universe.verifier.storage.Set(t.Context(), hexutil.Encode(messageID[:]), []byte("verified"))
	}

	// query using the eth client to determine the earliest log that was emitted.
	// since we didn't verify anything just yet, the earliest log should be the one
	// whose block number is returned.
	logs, err := universe.simClient.FilterLogs(t.Context(), ethereum.FilterQuery{
		FromBlock: big.NewInt(1),
		Addresses: []common.Address{universe.verifierEvents.Address()},
		Topics: [][]common.Hash{
			{verifier_events.VerifierEventsCCIPMessageSent{}.Topic()},
		},
	})
	require.NoError(t, err)
	// seq nrs 0 to 4 should already be verified, so the earliest log should be the one
	// right after 4, so 5.
	require.Len(t, logs, 10)
	var wantBlockNumber uint64
	for _, log := range logs {
		parsed, err := universe.verifierEvents.ParseCCIPMessageSent(log)
		require.NoError(t, err)
		if parsed.Message.Header.SequenceNumber == 5 {
			wantBlockNumber = log.BlockNumber
		}
	}

	lastProcessedBlock, err := universe.verifier.initializeLastProcessedBlock(t.Context())
	require.NoError(t, err)
	require.Equal(t, new(big.Int).SetUint64(wantBlockNumber), big.NewInt(lastProcessedBlock))
}

func Test_initializeLastProcessedBlock_noneVerified(t *testing.T) {
	universe := testSetup(t)

	require.NoError(t, universe.verifier.registerFilterIfNeeded(t.Context()))

	for seqNr := range 10 {
		universe.emitCCIPMessageSent(t, destChainSelector, uint64(seqNr), verifier_events.InternalEVM2AnyCommitVerifierMessage{
			Header: verifier_events.InternalHeader{
				MessageId:           [32]byte{byte(seqNr)},
				SourceChainSelector: sourceChainSelector,
				DestChainSelector:   destChainSelector,
				SequenceNumber:      uint64(seqNr),
			},
			Sender:             common.HexToAddress("0x1234"),
			Data:               nil,
			Receiver:           common.LeftPadBytes(common.HexToAddress("0x5678").Bytes(), 32),
			DestChainExtraArgs: nil,
			VerifierExtraArgs:  nil,
			FeeToken:           common.HexToAddress("0x0"),
			FeeTokenAmount:     big.NewInt(13371337),
			FeeValueJuels:      big.NewInt(13371337),
			TokenAmounts:       nil,
			RequiredVerifiers:  nil,
		})
	}

	universe.simBackend.Commit()

	// Wait for log poller to index the logs.
	// TODO: should be programmatic / event driven instead.
	time.Sleep(5 * time.Second)

	// query using the eth client to determine the earliest log that was emitted.
	// since we didn't verify anything just yet, the earliest log should be the one
	// whose block number is returned.
	logs, err := universe.simClient.FilterLogs(t.Context(), ethereum.FilterQuery{
		FromBlock: big.NewInt(1),
		Addresses: []common.Address{universe.verifierEvents.Address()},
		Topics: [][]common.Hash{
			{verifier_events.VerifierEventsCCIPMessageSent{}.Topic()},
		},
	})
	require.NoError(t, err)
	require.Len(t, logs, 10)
	head, err := universe.simClient.HeadByNumber(t.Context(), new(big.Int).SetUint64(logs[0].BlockNumber))
	require.NoError(t, err)
	t.Logf("earliest log block number: %d, and time stamp: %s", logs[0].BlockNumber, head.Timestamp)

	lastProcessedBlock, err := universe.verifier.initializeLastProcessedBlock(t.Context())
	require.NoError(t, err)
	require.Equal(t, new(big.Int).SetUint64(logs[0].BlockNumber), big.NewInt(lastProcessedBlock))
}
