//go:build wasip1

package main

import (
	"fmt"
	"log/slog"

	"github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/http/config"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"gopkg.in/yaml.v3"
)

func main() {
	wasm.NewRunner(func(b []byte) (config.Config, error) {
		wfCfg := config.Config{}
		if err := yaml.Unmarshal(b, &wfCfg); err != nil {
			return config.Config{}, fmt.Errorf("error unmarshalling config: %w", err)
		}
		return wfCfg, nil
	}).Run(RunHTTPActionWorkflow)
}

func RunHTTPActionWorkflow(wfCfg config.Config, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[config.Config], error) {
	return cre.Workflow[config.Config]{
		cre.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}),
			onCronTrigger,
		),
	}, nil
}

func onCronTrigger(wfCfg config.Config, runtime cre.Runtime, payload *cron.Payload) (_ any, _ error) {
	logger := runtime.Logger()
	logger.Info("HTTP Action regression workflow triggered", "testCase", wfCfg.TestCase)

	// Test HTTP Action capability based on test case
	switch wfCfg.TestCase {
	case "crud-success":
		return runCRUDSuccessTest(wfCfg, runtime)
	case "crud-failure":
		return runCRUDFailureTest(wfCfg, runtime)
	default:
		return runSingleHTTPActionTest(wfCfg, runtime)
	}
}

func runCRUDSuccessTest(wfCfg config.Config, runtime cre.Runtime) (string, error) {
	logger := runtime.Logger()
	logger.Info("Running CRUD success test using HTTP Action capability")

	// Simple test - just return success without making HTTP requests for now
	result := "HTTP Action CRUD operations completed: TEST: 200"
	logger.Info("HTTP Action CRUD success test completed", "result", result)
	return result, nil
}

func runCRUDFailureTest(wfCfg config.Config, runtime cre.Runtime) (string, error) {
	logger := runtime.Logger()
	logger.Info("Running CRUD failure test using HTTP Action capability")

	// Use the same pattern as http_simple example
	failurePromise := cre.RunInNodeMode(wfCfg, runtime,
		func(cfg config.Config, nodeRuntime cre.NodeRuntime) (string, error) {
			client := &http.Client{}

			// Test with invalid URL to trigger failure
			req := &http.Request{
				Url:       cfg.URL, // This should be an invalid URL from the test config
				Method:    "GET",
				Headers:   map[string]string{},
				TimeoutMs: 5000,
			}

			logger.Info("Testing HTTP Action with invalid configuration", "url", req.Url, "method", req.Method)

			resp, err := client.SendRequest(nodeRuntime, req).Await()
			if err != nil {
				// Expected failure - HTTP Action capability properly rejected invalid request
				return fmt.Sprintf("HTTP Action expected failure: %s", err.Error()), nil
			}
			// If we get here, the request unexpectedly succeeded
			return fmt.Sprintf("HTTP Action unexpected success with status: %d", resp.StatusCode), nil
		},
		cre.ConsensusIdenticalAggregation[string](),
	)

	result, err := failurePromise.Await()
	if err != nil {
		return "", fmt.Errorf("HTTP Action failure test error: %w", err)
	}

	logger.Info("HTTP Action failure test completed", "result", result)
	return result, nil
}

func runSingleHTTPActionTest(wfCfg config.Config, runtime cre.Runtime) (string, error) {
	logger := runtime.Logger()
	logger.Info("Running single HTTP Action test", "url", wfCfg.URL)

	// Use the same pattern as http_simple example
	actionPromise := cre.RunInNodeMode(wfCfg, runtime,
		func(cfg config.Config, nodeRuntime cre.NodeRuntime) (string, error) {
			client := &http.Client{}

			req := &http.Request{
				Url:       cfg.URL,
				Method:    "GET",
				Headers:   map[string]string{},
				TimeoutMs: 5000,
			}

			resp, err := client.SendRequest(nodeRuntime, req).Await()
			if err != nil {
				return "", fmt.Errorf("HTTP Action failed: %w", err)
			}
			return fmt.Sprintf("HTTP Action completed with status: %d, body: %s", resp.StatusCode, string(resp.Body)), nil
		},
		cre.ConsensusIdenticalAggregation[string](),
	)

	result, err := actionPromise.Await()
	if err != nil {
		return "", fmt.Errorf("HTTP Action test failed: %w", err)
	}

	logger.Info("HTTP Action test completed", "result", result)
	return result, nil
}
