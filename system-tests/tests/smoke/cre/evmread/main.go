//go:build wasip1

package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"runtime/debug"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/google/go-cmp/cmp"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	sdk "github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evmread/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evmread/types"
)

func main() {
	wasm.NewRunner(func(b []byte) (types.Config, error) {
		cfg := types.Config{}
		if err := json.Unmarshal(b, &cfg); err != nil {
			return types.Config{}, fmt.Errorf("error unmarshalling config: %w", err)
		}
		return cfg, nil
	}).Run(RunReadWorkflow)
}

func RunReadWorkflow(cfg types.Config, logger *slog.Logger, secretsProvider sdk.SecretsProvider) (sdk.Workflow[types.Config], error) {
	return sdk.Workflow[types.Config]{
		sdk.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}),
			onReadTrigger,
		),
	}, nil
}

func onReadTrigger(cfg types.Config, runtime sdk.Runtime, payload *cron.Payload) (_ any, _ error) {
	runtime.Logger().Info("onReadTrigger called", "payload", payload)
	defer func() {
		if r := recover(); r != nil {
			runtime.Logger().Error("recovered from panic", "recovered", r, "stack", string(debug.Stack()))
		}
	}()
	t := &T{Logger: runtime.Logger()}
	client := evm.Client{ChainSelector: cfg.ChainSelector}
	requireBalance(t, runtime, cfg, client)
	runtime.Logger().Info("Successfully got balance")

	latestHeadNumber := requireLatestBlockNumber(t, runtime, client)
	runtime.Logger().Info("Successfully got latestHeadNumber")

	requireEvent(t, cfg, runtime, latestHeadNumber, client)
	runtime.Logger().Info("Successfully got event")
	requireContractCall(t, cfg, runtime, client)
	runtime.Logger().Info("Successfully called contract")
	requireReceipt(t, runtime, cfg, client)
	runtime.Logger().Info("Successfully got receipt")
	requireTx(t, runtime, cfg, client)
	runtime.Logger().Info("Successfully got transaction")
	requireEstimatedGas(t, runtime, cfg, client)
	runtime.Logger().Info("Successfully estimated gas")
	requireError(t, runtime, cfg, client)
	runtime.Logger().Info("Successfully got error for non-existing transaction")
	txHash := sendTx(t, runtime, cfg, client, "EVM read workflow executed successfully")
	runtime.Logger().Info("Successfully sent transaction", "hash", common.Hash(txHash).String())
	return
}

func requireBalance(t *T, runtime sdk.Runtime, cfg types.Config, client evm.Client) {
	balanceReply, err := client.BalanceAt(runtime, &evm.BalanceAtRequest{
		Account:     cfg.AccountAddress,
		BlockNumber: nil,
	}).Await()
	require.NoError(t, err, "failed to get balance")
	require.NotNil(t, balanceReply, "BalanceAtReply should not be nil")
	require.NotNil(t, balanceReply.Balance, "Balance should not be nil")
	require.Equal(t, cfg.ExpectedBalance.String(), pb.NewIntFromBigInt(balanceReply.Balance).String(), "Balance should match expected value")
}

func requireError(t *T, runtime sdk.Runtime, cfg types.Config, client evm.Client) {
	txReply, err := client.GetTransactionByHash(runtime, &evm.GetTransactionByHashRequest{Hash: make([]byte, len(cfg.TxHash))}).Await()
	require.NotNil(t, err, "expected error when getting non existing transaction by hash")
	require.Nil(t, txReply, "txReply expected to be nil")
	require.ErrorContains(t, err, "not found", "expected error to be of type 'not found', got %s", err.Error())
	runtime.Logger().Info("Successfully got error for non-existing transaction", "error", err)
}

func requireEstimatedGas(t *T, runtime sdk.Runtime, cfg types.Config, client evm.Client) {
	estimatedGasReply, err := client.EstimateGas(runtime, &evm.EstimateGasRequest{
		Msg: &evm.CallMsg{
			To:   cfg.ContractAddress,
			Data: cfg.ExpectedTx.Data,
		},
	}).Await()
	require.NoError(t, err, "failed to estimate gas")
	require.NotNil(t, estimatedGasReply, "EstimateGasReply should not be nil")
	require.Greater(t, estimatedGasReply.Gas, uint64(0), "Estimated gas should greater than 0")
}

func requireTx(t *T, runtime sdk.Runtime, cfg types.Config, client evm.Client) {
	txReply, err := client.GetTransactionByHash(runtime, &evm.GetTransactionByHashRequest{Hash: cfg.TxHash}).Await()
	require.NoError(t, err, "failed to get transaction by hash")
	require.NotNil(t, txReply, "GetTransactionByHashReply should not be nil")
	require.NotNil(t, txReply.Transaction, "Transaction should not be nil")
	require.Empty(t, cmp.Diff(txReply.Transaction, cfg.ExpectedTx, protocmp.Transform()))
}

func requireReceipt(t *T, runtime sdk.Runtime, cfg types.Config, client evm.Client) {
	receiptReply, err := client.GetTransactionReceipt(runtime, &evm.GetTransactionReceiptRequest{Hash: cfg.TxHash}).Await()
	require.NoError(t, err, "failed to get transaction receipt")
	require.NotNil(t, receiptReply, "TransactionReceiptReply should not be nil")
	require.NotNil(t, receiptReply.Receipt, "TransactionReceipt should not be nil")
	require.Equal(t, len(cfg.ExpectedReceipt.Logs), len(receiptReply.Receipt.Logs), "Logs length should match expected value")
	cfg.ExpectedReceipt.Logs = nil
	receiptReply.Receipt.Logs = nil
	require.Empty(t, cmp.Diff(receiptReply.Receipt, cfg.ExpectedReceipt, protocmp.Transform()))
}

func requireContractCall(t *T, cfg types.Config, runtime sdk.Runtime, client evm.Client) {
	parsed, err := abi.JSON(strings.NewReader(contracts.MessageEmitterMetaData.ABI))
	require.NoError(t, err, "failed to parse api")
	const callArg = "Hey CRE"
	const methodName = "getMessage"
	packed, err := parsed.Pack(methodName, callArg)
	require.NoError(t, err, "failed to pack getMessage")
	callContractReply, err := client.CallContract(runtime, &evm.CallContractRequest{
		Call: &evm.CallMsg{
			To:   cfg.ContractAddress,
			Data: packed,
		},
	}).Await()
	require.NoError(t, err, "failed to call contract")
	require.NotNil(t, callContractReply, "CallContractReply should not be nil")
	var result string
	err = parsed.UnpackIntoInterface(&result, methodName, callContractReply.Data)
	require.NoError(t, err, "failed to unpack into result")
	require.Equal(t, "getMessage returns: "+callArg, string(result))
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

func sendTx(t *T, runtime sdk.Runtime, cfg types.Config, client evm.Client, msg string) []byte {
	// NOTE: This is not a right way to send a transaction. Msg must be properly encoded to trigger a proper receiver contract call.
	// In this case we just need to see transaction on chain, so it's sufficient.
	report, err := runtime.GenerateReport(&sdkpb.ReportRequest{
		EncodedPayload: []byte(msg),
		EncoderName:    "evm",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	}).Await()
	require.NoError(t, err, "failed to generate report")
	reportReply, err := client.WriteReport(runtime, &evm.WriteCreReportRequest{
		Receiver:  cfg.ContractAddress,
		Report:    report,
		GasConfig: &evm.GasConfig{GasLimit: 1e6},
	}).Await()
	require.NoError(t, err, "failed to write report")
	require.NotNil(t, reportReply)
	return reportReply.TxHash
}

func requireEvent(t *T, cfg types.Config, runtime sdk.Runtime, latestHeadNumber int64, client evm.Client) {
	const blocksStep = 100
	foundEvent := false
	for ; latestHeadNumber > 0; latestHeadNumber -= blocksStep {
		eventsReply, err := client.FilterLogs(runtime, &evm.FilterLogsRequest{FilterQuery: &evm.FilterQuery{
			FromBlock: pb.NewBigIntFromInt(big.NewInt(max(latestHeadNumber-blocksStep, 1))),
			ToBlock:   pb.NewBigIntFromInt(big.NewInt(latestHeadNumber)),
			Addresses: [][]byte{cfg.ContractAddress},
		}}).Await()
		require.NoError(t, err, "failed to filter logs")
		require.NotNil(t, eventsReply, "FilterLogsReply should not be nil")
		if len(eventsReply.Logs) > 0 {
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
