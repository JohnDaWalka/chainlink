//go:build wasip1

package main

import (
	"fmt"

	evm "github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	cron "github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"

	//"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

type None struct{}

func main() {
	wasm.NewRunner(func(configBytes []byte) (None, error) {
		return None{}, nil
	}).Run(RunSimpleEvmActionWorkflow)
}

func RunSimpleEvmActionWorkflow(wcx *sdk.Environment[None]) (sdk.Workflow[None], error) {
	wcx.Logger.Info("EVM_CAP RunSimpleEvmActionWorkflow EVM")
	fmt.Println("EVM_CAP RunSimpleEvmActionWorkflow2 EVM")
	workflows := sdk.Workflow[None]{
		sdk.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}),
			onTrigger,
		),
	}
	return workflows, nil
}

func onTrigger(wcx *sdk.Environment[None], runtime sdk.Runtime, trigger *cron.Payload) (string, error) {
	wcx.Logger.Info("Cron trigger called us EVM!")
	fmt.Println("EVM_CAP 001 onTrigger()")

	message := fmt.Sprintf(
		"RESPONSE: Cron trigger called us EVM! %s",
		trigger.ScheduledExecutionTime,
	)

	client := &evm.Client{
		ChainSelector: 3379446385462418246, // geth chain selector
	}
	hash := common.HexToHash("0x5a38eb835b8340fceccf86e048ee07d9ca469770eb675dad6487fde90d7e7ff5")
	fmt.Printf("EVM_CAP 002.1 preparing to call GetTransactionByHash with hash %s\n", hash.Hex())
	promiseGetTxByHash := client.GetTransactionByHash(runtime, &evm.GetTransactionByHashRequest{
		Hash: hash.Bytes(),
	})
	fmt.Println("EVM_CAP 002.2 method GetTransactionByHash has been called")
	tx, err := promiseGetTxByHash.Await()
	fmt.Println("EVM_CAP 002.3 after Await()")
	if err != nil {
		fmt.Println("EVM_CAP 003 ERROR 001 failed to get GetTransactionByHash")
		fmt.Println("EVM_CAP 003 ERROR 002 failed to get GetTransactionByHash:", err)
		return "", fmt.Errorf("EVM_CAP ERROR failed to get GetTransactionByHash: %w", err)
	}
	fmt.Println("EVM_CAP 003 after GetTransactionByHash")
	fmt.Printf("EVM_CAP 003 after GetTransactionByHash: %v", tx)
	message = fmt.Sprintf(
		"RESPONSE: Transaction Hash: %s, nonce: %d",
		tx.GetTransaction().Hash,
		tx.GetTransaction().Nonce,
	)
	fmt.Printf("EVM_CAP 004 message created: %q\n", message)

	wcx.Logger.Info(message)
	fmt.Println(message)
	return "such a lovely disaster: " + message, nil
}
