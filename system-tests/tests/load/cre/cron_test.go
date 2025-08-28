package cre

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp/benchspy"
	mockcapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock"
	pb2 "github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock/pb"

	"github.com/docker/docker/client"
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

// TestCron performs a comprehensive load test of the cron trigger capability.
// It validates that the cron trigger can handle sustained load and provides
// timing for scheduled executions across multiple nodes.
//
// The test performs the following steps:
// 1. Sets up a mock capability controller
// 2. Configures a cron trigger with 30-second intervals
// 3. Runs a load test with multiple virtual users
// 4. Validates trigger responses and timing accuracy
// 5. Generates performance reports
func TestCron(t *testing.T) {
	// Connect to the cluster
	ip, err := getDockerContainerIP("workflow-node1")
	require.NoError(t, err, "could not get container IP")
	mockClient := mockcapability.NewMockCapabilityController(framework.L)
	err = mockClient.ConnectAll([]string{fmt.Sprintf("%s:%d", ip, defaultMockEndpointPort)}, true, false)
	require.NoError(t, err, "connecting with mock client failed")

	// Use WASP to trigger registrations to the cron-trigger

	// We want to see n responses back in order to consider it sucessful,
	// For example, if we can sustain it for 5m then we consider it successful
	payload, err := anypb.New(&cron.Config{Schedule: cronSchedule})
	require.NoError(t, err, "creating payload failed")
	executionTime := time.Minute * 1
	vu := &VirtualUser{
		VUControl:      wasp.NewVUControl(),
		mockController: mockClient,
		triggerID:      triggerID,
		executionTime:  executionTime,
		payload:        payload,
		metadata: &pb2.Metadata{
			WorkflowID: workflowID,
		},
	}
	lokiURL := defaultLokiURL
	emptyString := ""
	lokiConfig := wasp.NewLokiConfig(&lokiURL, &emptyString, &emptyString, &emptyString)
	generator, err := wasp.NewGenerator(&wasp.Config{
		GenName:     "cron-load-test",
		Labels:      map[string]string{"branch": "profile-check", "commit": "profile-check"},
		CallTimeout: executionTime + time.Minute,
		T:           t,
		LoadType:    wasp.VU,
		VU:          vu,
		Schedule: wasp.Combine(
			wasp.Plain(100, executionTime),
			wasp.Plain(200, executionTime),
			wasp.Plain(300, executionTime),
			wasp.Plain(400, executionTime),
			wasp.Plain(500, executionTime),
			wasp.Plain(1000, executionTime),
			wasp.Plain(1500, executionTime),
			wasp.Plain(2000, executionTime),
		),
		LokiConfig: lokiConfig,
	})
	require.NoError(t, err, "could not create generator")

	_, err = wasp.NewProfile().Add(generator, nil).Run(true)
	require.NoError(t, err)

	prometheusExecutor, err := benchspy.NewPrometheusQueryExecutor(
		map[string]string{
			"cpu":         `rate(container_cpu_usage_seconds_total{name="workflow-node1"}[5m]) * 100`,
			"mem":         `container_memory_working_set_bytes{name="workflow-node1"} /1024/1024`,
			"mem_rss":     `container_memory_rss{name="workflow-node1"} /1024/1024`,
			"network_tx":  `rate(container_network_transmit_bytes_total{name="workflow-node1"}[5m])`,
			"network_rx":  `rate(container_network_receive_bytes_total{name="workflow-node1"}[5m])`,
			"disk_reads":  `rate(container_fs_reads_bytes_total{name="workflow-node1"}[5m])`,
			"disk_writes": `rate(container_fs_writes_bytes_total{name="workflow-node1"}[5m])`,
		},
		&benchspy.PrometheusConfig{
			Url:               defaultPrometheusURL,
			NameRegexPatterns: []string{},
		},
	)
	require.NoError(t, err)

	report, err := benchspy.NewStandardReport("profile-check",
		benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Loki),
		benchspy.WithQueryExecutors(prometheusExecutor),
		benchspy.WithGenerators(generator))
	require.NoError(t, err, "creating report failed")
	store, err := report.Store()
	require.NoError(t, err, "storing report failed")
	fmt.Printf("Report stored at %s\n", store)
}

type VirtualUser struct {
	*wasp.VUControl
	mockController        *mockcapability.Controller
	triggerID             string
	executionTime         time.Duration
	triggerRegistrationID string
	triggerCh             []chan *capabilities.TriggerResponse
	metadata              *pb2.Metadata
	payload               *anypb.Any
}

func (v *VirtualUser) Clone(l *wasp.Generator) wasp.VirtualUser {
	return &VirtualUser{
		VUControl:      wasp.NewVUControl(),
		mockController: v.mockController,
		triggerID:      v.triggerID,
		executionTime:  v.executionTime,
		payload:        v.payload,
		metadata:       v.metadata,
	}
}

func (v *VirtualUser) Setup(l *wasp.Generator) error {
	v.triggerRegistrationID = uuid.New().String()
	chList, err := v.mockController.RegisterTrigger(context.Background(), v.triggerID, v.metadata, nil, v.payload, "", v.triggerRegistrationID)
	if err != nil {
		return err
	}
	v.triggerCh = chList
	return nil
}

func (v *VirtualUser) Teardown(l *wasp.Generator) error {
	// return v.mockController.UnregisterTrigger(context.Background(), v.triggerID, v.metadata, nil, v.payload, "", v.triggerRegistrationID)
	return nil
}

func (v *VirtualUser) Call(l *wasp.Generator) {
	// Calculate the number of thick we expect to get back to consider the call successful
	expectedCalls := int(v.executionTime.Seconds() / 30)
	confirmedCalls := make([]int, len(v.mockController.Nodes))
	lastTicks := make([]time.Time, len(v.mockController.Nodes))

	for i := range lastTicks {
		lastTicks[i] = time.Now()
	}

	wg := sync.WaitGroup{}
	wg.Add(len(v.triggerCh))

	for i, ch := range v.triggerCh {
		go func(i int) {
			defer wg.Done()
			for {
				msg, ok := <-ch
				if !ok {
					l.Responses.Err(&resty.Response{Request: &resty.Request{}}, "virtual-user-call-generation", errors.New("channel closed"))
					return
				}

				lastTickDiff := time.Since(lastTicks[i])
				lastTicks[i] = time.Now()
				if msg.Err != nil {
					l.Responses.Err(&resty.Response{Request: &resty.Request{}}, "virtual-user-call-generation", msg.Err)
					return
				}
				confirmedCalls[i]++
				l.ResponsesChan <- &wasp.Response{Data: v, Duration: lastTickDiff}
				if confirmedCalls[i] == expectedCalls {
					return
				}
			}
		}(i)
	}
	wg.Wait()
	err := v.Teardown(l)
	if err != nil {
		l.Responses.Err(&resty.Response{Request: &resty.Request{}}, "virtual-user-call-generation", err)
	}
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
