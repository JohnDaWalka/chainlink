//go:build wasip1

package main

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"
	"runtime/debug"
	"strings"

	"TestBR/contracts/evm/src/generated/logger_tester"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/google/go-cmp/cmp"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm/bindings"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	sdk "github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"

	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
)

func main() {
	wasm.NewRunner(sdk.ParseJSON[Config]).Run(RunTesterWorkflow)
}

func RunTesterWorkflow(cfg *Config, _ *slog.Logger, _ sdk.SecretsProvider) (sdk.Workflow[*Config], error) {
	eventSig, err := hex.DecodeString(strings.TrimPrefix(cfg.LogTriggerEventSig, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to decode event signature: %w", err)
	}

	ltCfg := &evm.FilterLogTriggerRequest{
		Addresses: [][]byte{common.HexToAddress(cfg.LoggerTesterAddress).Bytes()},
		Topics: []*evm.TopicValues{
			{
				Values: [][]byte{eventSig},
			},
		},
		Confidence: evm.ConfidenceLevel_CONFIDENCE_LEVEL_LATEST,
	}

	return sdk.Workflow[*Config]{
		sdk.Handler(
			evm.LogTrigger(cfg.ChainSelector, ltCfg),
			onReadTrigger,
		),
	}, nil
}

func onReadTrigger(cfg *Config, runtime sdk.Runtime, outputs *evm.Log) (_ any, _ error) {
	runtime.Logger().Info("onReadTrigger called", "payload", logToString(outputs))
	defer func() {
		if r := recover(); r != nil {
			runtime.Logger().Error("recovered from panic", "recovered", r, "stack", string(debug.Stack()))
		}
	}()

	t := &T{Logger: runtime.Logger()}
	client := evm.Client{ChainSelector: cfg.ChainSelector}

	loggerTesterAddress := common.HexToAddress(cfg.LoggerTesterAddress)
	loggerTester, err := logger_tester.NewLoggerTester(&client, loggerTesterAddress, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create contract instance: %w", err)
	}

	// Verify BalanceAt
	requireBalance(t, runtime, cfg.Balance, client)
	runtime.Logger().Info("Successfully got balance")

	// Verify HeaderByNumber
	latestHeadNumber := requireLatestBlockNumber(t, runtime, client)
	runtime.Logger().Info("Successfully got latestHeadNumber")

	// Verify FilterLogs
	requireEvent(t, runtime, latestHeadNumber, loggerTester)
	runtime.Logger().Info("Successfully got event")

	// Verify ContractCall
	requireContractCall(t, cfg, runtime, loggerTester, client)
	runtime.Logger().Info("Successfully called contract")

	// Verify GetTransactionReceipt
	requireReceipt(t, runtime, cfg.Tx, client)
	runtime.Logger().Info("Successfully got receipt")

	// Verify GetTransactionByHash
	var expectedTx types.Transaction
	require.NoError(t, expectedTx.UnmarshalBinary(common.FromHex(cfg.Tx.ExpectedBinaryTx)))
	requireTx(t, runtime, &expectedTx, client)
	runtime.Logger().Info("Successfully got transaction")

	// Verify EstimateGas
	requireEstimatedGas(t, runtime, cfg, expectedTx.Data(), client)
	runtime.Logger().Info("Successfully estimated gas")

	// Verify error on non-existing transaction
	requireError(t, runtime, client)
	runtime.Logger().Info("Successfully got error for non-existing transaction")

	// Verify WriteReport
	txHash := sendTx(t, runtime, loggerTester)
	runtime.Logger().Info("Successfully sent transaction", "hash", common.Hash(txHash).String())
	return
}

func requireBalance(t *T, runtime sdk.Runtime, balance Balance, client evm.Client) {
	balanceReply, err := client.BalanceAt(runtime, &evm.BalanceAtRequest{
		Account:     common.HexToAddress(balance.Address).Bytes(),
		BlockNumber: pb.NewBigIntFromInt(big.NewInt(rpc.LatestBlockNumber.Int64())),
	}).Await()
	require.NoError(t, err, "failed to get balance")
	require.NotNil(t, balanceReply, "BalanceAtReply should not be nil")
	require.NotNil(t, balanceReply.Balance, "Balance should not be nil")
	require.Equal(t, balance.ExpectedBalance.String(), pb.NewIntFromBigInt(balanceReply.Balance).String(), fmt.Sprintf("Balance should match expected value"))
}

func requireError(t *T, runtime sdk.Runtime, client evm.Client) {
	txReply, err := client.GetTransactionByHash(runtime, &evm.GetTransactionByHashRequest{Hash: make([]byte, common.HashLength)}).Await()
	require.NotNil(t, err, "expected error when getting non existing transaction by hash")
	require.Nil(t, txReply, "txReply expected to be nil")
	runtime.Logger().Info("Successfully got error for non-existing transaction", "error", err)
	require.ErrorContains(t, err, "not found", "expected error to be of type 'not found', got %s", err.Error())
}

func requireEstimatedGas(t *T, runtime sdk.Runtime, cfg *Config, txData []byte, client evm.Client) {
	estimatedGasReply, err := client.EstimateGas(runtime, &evm.EstimateGasRequest{
		Msg: &evm.CallMsg{
			From: common.HexToAddress(cfg.Tx.From).Bytes(),
			To:   common.HexToAddress(cfg.LoggerTesterAddress).Bytes(),
			Data: txData,
		},
	}).Await()
	require.NoError(t, err, "failed to estimate gas")
	require.NotNil(t, estimatedGasReply, "EstimateGasReply should not be nil")
	require.Greater(t, estimatedGasReply.Gas, uint64(0), "Estimated gas should greater than 0")
}

func requireTx(t *T, runtime sdk.Runtime, expectedTx *types.Transaction, client evm.Client) {
	txReply, err := client.GetTransactionByHash(runtime, &evm.GetTransactionByHashRequest{Hash: expectedTx.Hash().Bytes()}).Await()
	require.NoError(t, err, "failed to get transaction by hash")
	require.NotNil(t, txReply, "GetTransactionByHashReply should not be nil")
	require.NotNil(t, txReply.Transaction, "Transaction should not be nil")
	sdkExpectedTx := &evm.Transaction{
		Nonce:    expectedTx.Nonce(),
		Gas:      expectedTx.Gas(),
		To:       expectedTx.To().Bytes(),
		Data:     expectedTx.Data(),
		Hash:     expectedTx.Hash().Bytes(),
		Value:    pb.NewBigIntFromInt(expectedTx.Value()),
		GasPrice: pb.NewBigIntFromInt(expectedTx.GasPrice()),
	}
	require.Empty(t, cmp.Diff(txReply.Transaction, sdkExpectedTx, protocmp.Transform()))
}

func gethToSDKReceipt(r *types.Receipt) *evm.Receipt {
	return &evm.Receipt{
		Status:            r.Status,
		Logs:              make([]*evm.Log, len(r.Logs)), // workflow compares only number of logs, not their content
		TxHash:            r.TxHash.Bytes(),
		ContractAddress:   r.ContractAddress.Bytes(),
		GasUsed:           r.GasUsed,
		BlockHash:         r.BlockHash.Bytes(),
		BlockNumber:       pb.NewBigIntFromInt(r.BlockNumber),
		TxIndex:           uint64(r.TransactionIndex),
		EffectiveGasPrice: pb.NewBigIntFromInt(r.EffectiveGasPrice),
	}
}

func requireReceipt(t *T, runtime sdk.Runtime, tx Tx, client evm.Client) {
	receiptReply, err := client.GetTransactionReceipt(runtime, &evm.GetTransactionReceiptRequest{Hash: common.HexToHash(tx.TxHash).Bytes()}).Await()
	require.NoError(t, err, "failed to get transaction receipt")
	require.NotNil(t, receiptReply, "TransactionReceiptReply should not be nil")
	require.NotNil(t, receiptReply.Receipt, "TransactionReceipt should not be nil")
	require.Equal(t, len(tx.ExpectedReceipt.Logs), len(receiptReply.Receipt.Logs), "Logs length should match expected value")
	//tx.ExpectedReceipt.Logs = nil
	//receiptReply.Receipt.Logs = nil
	expectedReceipt := gethToSDKReceipt(tx.ExpectedReceipt)
	require.Empty(t, cmp.Diff(receiptReply.Receipt, expectedReceipt, protocmp.Transform()))
}

func requireContractCall(t *T, cfg *Config, runtime sdk.Runtime, loggerTester *logger_tester.LoggerTester, client evm.Client) {
	//  [4]byte{128, 95, 33, 50} is the interface ID of IReceiver
	data, err := loggerTester.Codec.EncodeSupportsInterfaceMethodCall(logger_tester.SupportsInterfaceInput{InterfaceId: [4]byte{128, 95, 33, 50}})
	require.NoError(t, err, "failed to encode SupportsInterface method call")

	callContractReply, err := client.CallContract(runtime, &evm.CallContractRequest{
		Call: &evm.CallMsg{
			To:   common.HexToAddress(cfg.LoggerTesterAddress).Bytes(),
			Data: data,
		},
	}).Await()
	require.NoError(t, err, "failed to call contract")
	require.NotNil(t, callContractReply, "CallContractReply should not be nil")

	isSupported, err := loggerTester.Codec.DecodeSupportsInterfaceMethodOutput(callContractReply.Data)
	require.NoError(t, err, "failed to unpack into result")
	require.Equal(t, true, isSupported, "SupportsInterface returns: ")
}

func requireLatestBlockNumber(t *T, runtime sdk.Runtime, client evm.Client) int64 {
	headerToFetch := []rpc.BlockNumber{rpc.FinalizedBlockNumber, rpc.SafeBlockNumber, rpc.LatestBlockNumber}
	var prevHeaderNumber *big.Int
	for _, headToFetch := range headerToFetch {
		runtime.Logger().Info("Fetching header", "headToFetch", headToFetch)
		headerReply, err := client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{BlockNumber: pb.NewBigIntFromInt(big.NewInt(headToFetch.Int64()))}).Await()
		require.NoError(t, err)
		require.NotNil(t, headerReply, "HeaderByNumberReply should not be nil %s", headToFetch)
		require.NotNil(t, headerReply.Header, "Header should not be nil %s", headToFetch)
		headerNumber := pb.NewIntFromBigInt(headerReply.Header.BlockNumber)
		runtime.Logger().Info("Header fetched", "blockNumber", headerNumber.String())
		if prevHeaderNumber != nil {
			require.True(t, headerNumber.Cmp(prevHeaderNumber) >= 0,
				"Expected prev head to have higher or equal block number. Current header: %s, Previous header: %s. HeadToFetch",
				headerNumber, prevHeaderNumber, headerToFetch)
		}
		prevHeaderNumber = headerNumber
	}
	return prevHeaderNumber.Int64()
}

func sendTx(t *T, runtime sdk.Runtime, loggerTester *logger_tester.LoggerTester) []byte {
	report, err := runtime.GenerateReport(&sdkpb.ReportRequest{
		EncodedPayload: []byte("empty"),
		EncoderName:    "evm",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	}).Await()
	require.NoError(t, err, "failed to generate report")
	reportReply, err := loggerTester.WriteReport(runtime, report, &evm.GasConfig{GasLimit: 500_000}).Await()
	require.NoError(t, err, "failed to write report")
	require.NotNil(t, reportReply)
	return reportReply.TxHash
}

func requireEvent(t *T, runtime sdk.Runtime, latestHeadNumber int64, loggerTester *logger_tester.LoggerTester) {
	const blocksStep = 100
	foundEvent := false
	for ; latestHeadNumber > 0; latestHeadNumber -= blocksStep {
		logsReply, err := loggerTester.FilterLogsLogEmitted(runtime, &bindings.FilterOptions{
			BlockHash: nil,
			FromBlock: big.NewInt(max(latestHeadNumber-blocksStep, 1)),
			ToBlock:   big.NewInt(latestHeadNumber),
		}).Await()
		require.NoError(t, err, "failed to filter logs")
		require.NotNil(t, logsReply, "FilterLogsReply should not be nil")
		if len(logsReply.Logs) > 0 {
			foundEvent = true
			break
		}
	}
	require.True(t, foundEvent, "Failed to find at least one event")
}

type T struct {
	*slog.Logger
}

func (t *T) Errorf(format string, args ...interface{}) {
	// if the log was produced by require/assert we need to split it, as engine does not allow logs longer than 1k bytes
	if len(args) > 0 {
		if msg, ok := args[0].(string); ok && strings.Contains(msg, "Error:") && strings.Contains(msg, "Error Trace:") {
			for _, line := range strings.Split(msg, "Error:") {
				t.Logger.Error(line)
			}
			return
		}
	}
	t.Logger.Error(fmt.Sprintf(format, args...))
	panic(fmt.Sprintf(format, args...)) // panic to stop execution
}

func (t *T) FailNow() {
	panic("Test failed. Panic to stop execution")
}

func logToString(l *evm.Log) string {
	formatBytes := func(b []byte) string {
		if len(b) == 0 {
			return "nil"
		}
		return "0x" + hex.EncodeToString(b)
	}

	formatTopics := func(topics [][]byte) string {
		if len(topics) == 0 {
			return "[]"
		}
		out := make([]string, len(topics))
		for i, t := range topics {
			out[i] = formatBytes(t)
		}
		return "[" + strings.Join(out, ", ") + "]"
	}

	blockNum := "nil"
	if l.BlockNumber != nil {
		blockNum = l.BlockNumber.String()
	}

	return fmt.Sprintf(
		"Log{Address:%s, Topics:%s, TxHash:%s, BlockHash:%s, Data:%s, EventSig:%s, BlockNumber:%s, TxIndex:%d, Index:%d, Removed:%t}",
		formatBytes(l.Address),
		formatTopics(l.Topics),
		formatBytes(l.TxHash),
		formatBytes(l.BlockHash),
		formatBytes(l.Data),
		formatBytes(l.EventSig),
		blockNum,
		l.TxIndex,
		l.Index,
		l.Removed,
	)
}

// Tx fields should all be from the same transaction.
type Tx struct {
	TxHash           string         `json:"txHash"`
	ExpectedReceipt  *types.Receipt `json:"expectedReceipt"`
	ExpectedBinaryTx string         `json:"expectedBinaryTx"`
	From             string         `json:"from"`
	To               string         `json:"to"`
}

type Balance struct {
	Address         string   `json:"address"`
	ExpectedBalance *big.Int `json:"expectedBalance"`
}
type Config struct {
	ChainSelector       uint64  `json:"chainSelector"`
	LogTriggerEventSig  string  `json:"logTriggerEventSig"`
	LoggerTesterAddress string  `json:"loggerTesterAddress"`
	Balance             Balance `json:"balance"`
	Tx                  Tx      `json:"tx"`
}
