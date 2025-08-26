//go:build wasip1

package main

import (
	"log/slog"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

type None struct{}

func main() {
	wasm.NewRunner(func(configBytes []byte) (None, error) {
		return None{}, nil
	}).Run(RunSimpleCronWorkflow)
}

func RunSimpleCronWorkflow(_ None, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[None], error) {
	workflows := cre.Workflow[None]{
		cre.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}),
			onTrigger,
		),
	}
	return workflows, nil
}

func onTrigger(cfg None, runtime cre.Runtime, _ *cron.Payload) (string, error) {
	runtime.Logger().Info("Now triggered fetch of of chain value")

	mathPromise := cre.RunInNodeMode(cfg, runtime, fetchData, cre.ConsensusIdenticalAggregation[float64]())
	offchainValue, err := mathPromise.Await()
	if err != nil {
		runtime.Logger().Info("Got an error oh no2", "error", err)
		return "", err
	}
	runtime.Logger().Info("Successfully fetched offchain value and reached consensus", "result", offchainValue)

	return "such a lovely disaster", nil
}

func fetchData(cfg None, nodeRuntime cre.NodeRuntime) (float64, error) {
	// pretend we're fetching some node-mode data
	return 420.69, nil
}
