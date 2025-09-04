package cre

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	mockcapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock/pb"
)

const (
	// cronSchedule defines the cron expression for triggering every 30 seconds
	cronSchedule = "*/30 * * * * *"
	// triggerID is the identifier for the cron trigger capability
	triggerID = "cron-trigger@1.0.0"
	// workflowID is the unique identifier for the test workflow
	workflowID = "cron-load-test"
	// defaultLokiURL is the default endpoint for Loki logging integration
	defaultLokiURL       = "http://localhost:3030/loki/api/v1/push"
	defaultPrometheusURL = "http://localhost:9099"
	// defaultMockEndpoint is the default address for the mock capability controller
	defaultMockEndpointPort = 7777
)

func TestOCR(t *testing.T) {
	// Connect to the cluster
	ip, err := getDockerContainerIP("workflow-node1")
	require.NoError(t, err, "could not get container IP")
	mockClient := mockcapability.NewMockCapabilityController(framework.L)
	err = mockClient.ConnectAll([]string{fmt.Sprintf("%s:%d", ip, defaultMockEndpointPort)}, true, false)
	require.NoError(t, err, "connecting with mock client failed")

	// offchain_reporting@1.0.0
	/*
	   config:
	     report_id: '0001'
	     key_id: 'evm'
	     aggregation_method: data_feeds
	     aggregation_config:
	       allowedPartialStaleness: '0.5'
	       feeds:
	         '0x000351de403f638036014add21a5abd5f464bf21d11aa356dfc6dbe4e2384e4e':  # BTC/USD
	           deviation: '0.01'
	           heartbeat: 600
	           remappedID: '0x666666666666'
	         '0x0003f2f4cae1891f647db8d73c87a7a03888bd176afdb7206853da9abfc92874': # ETH/USD
	           deviation: '0.01'
	           heartbeat: 600
	           remappedID: '0x777777777777'
	         '0x00034db6355441c80b613f666757c63777dae7743885a9c594ca25d9f9b896ca': # LINK/USD
	           deviation: '0.01'
	           heartbeat: 600
	     encoder: EVM
	     encoder_config:
	       abi: (bytes32 FeedID, uint224 Price, uint32 Timestamp)[] Reports
	*/
	config, err := values.NewMap(
		map[string]any{
			"aggregation_method": "data_feeds", // data_feeds OR reduce
			"aggregation_config": map[string]any{},
			"encoder_config":     map[string]any{
				// "abi": "(bytes32 FeedID, uint224 Price, uint32 Timestamp)[] Reports",
			},
			"encoder":   "EVM",
			"report_id": "ffff",
			"key_id":    "evm",
		},
	)
	require.NoError(t, err, "failed to create config")

	configBytes, err := mockcapability.MapToBytes(config)
	require.NoError(t, err, "failed to convert config to bytes")

	// Register to workflow
	err = mockClient.RegisterToWorkflow(context.Background(), &pb.RegisterToWorkflowRequest{
		ID:                   "offchain_reporting@1.0.0",
		CapabilityType:       3,
		RegistrationMetadata: nil,
		Config:               configBytes,
	})
	require.NoError(t, err, "registering to workflow failed")
	// Execute
	resps, err := mockClient.Execute(context.Background(), &pb.ExecutableRequest{
		ID:              "offchain_reporting@1.0.0",
		CapabilityType:  3,
		RequestMetadata: &pb.Metadata{},
		Config:          nil,
		Inputs:          nil,
		Payload:         nil,
		ConfigPayload:   nil,
		Method:          "",
		CapabilityId:    "",
	})
	require.NoError(t, err, "executing failed")
	spew.Dump(resps)
}

func getDockerContainerIP(containerName string) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	inspect, err := cli.ContainerInspect(context.Background(), containerName)
	if err != nil {
		return "", err
	}

	n, ok := inspect.NetworkSettings.Networks["ctf"]
	if !ok {
		return "", errors.New("ctf network not found")
	}

	return n.IPAddress, nil
}
