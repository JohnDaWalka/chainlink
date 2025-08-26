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
)

func TestCron(t *testing.T) {
	// Connect to the cluster
	mockClient := mockcapability.NewMockCapabilityController(framework.L)

	err := mockClient.ConnectAll([]string{"192.168.48.19:7777"}, true, false)
	require.NoError(t, err, "connecting with mock client failed")

	// Use WASP to trigger registrations to the cron-trigger

	// We want to see n responses back in order to consider it sucessful
	// For example if we can sustain it for 5m then we consider it successful
	metadata := &pb2.Metadata{
		WorkflowID: "load-test",
	}
	payload, err := anypb.New(&cron.Config{Schedule: "*/30 * * * * *"})
	require.NoError(t, err, "creating payload failed")
	executionTime := time.Minute * 1
	vu := &VirtualUser{
		VUControl:      wasp.NewVUControl(),
		mockController: mockClient,
		triggerID:      "cron-trigger@1.0.0",
		executionTime:  executionTime,
		payload:        payload,
		metadata:       metadata,
		closeCh:        make(chan struct{}),
	}
	lokiURL := "http://localhost:3030/loki/api/v1/push"
	lokiToken := ""
	lokiTenant := ""
	generator, err := wasp.NewGenerator(&wasp.Config{
		CallTimeout: executionTime + time.Minute,
		T:           t,
		LoadType:    wasp.VU,
		VU:          vu,
		Schedule: wasp.Combine(
			wasp.Plain(3000, executionTime),
			wasp.Plain(1, executionTime),
			wasp.Plain(5000, executionTime),
		),
		LokiConfig: wasp.NewLokiConfig(&lokiURL, &lokiTenant, &lokiToken, &lokiToken),
	})
	require.NoError(t, err, "could not create generator")

	_, err = wasp.NewProfile().Add(generator, nil).Run(true)
	require.NoError(t, err)

	report, err := benchspy.NewStandardReport("cron-load-test", benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Direct),
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
	closeCh               chan struct{}
}

func (v *VirtualUser) Clone(l *wasp.Generator) wasp.VirtualUser {
	return &VirtualUser{
		VUControl:      wasp.NewVUControl(),
		mockController: v.mockController,
		triggerID:      v.triggerID,
		executionTime:  v.executionTime,
		payload:        v.payload,
		metadata:       v.metadata,
		closeCh:        make(chan struct{}),
	}
}

func (v *VirtualUser) Setup(l *wasp.Generator) error {
	v.triggerRegistrationID = uuid.New().String()
	v.closeCh = make(chan struct{})
	chList, err := v.mockController.RegisterTrigger(context.Background(), v.triggerID, v.metadata, nil, v.payload, "", v.triggerRegistrationID)
	if err != nil {
		return err
	}
	v.triggerCh = chList
	return nil
}

func (v *VirtualUser) Teardown(l *wasp.Generator) error {
	v.closeCh <- struct{}{}
	return v.mockController.UnregisterTrigger(context.Background(), v.triggerID, v.metadata, nil, v.payload, "", v.triggerRegistrationID)
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
				select {
				case msg, ok := <-ch:
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
				case <-v.closeCh:
					return
				}

			}
		}(i)
	}
	wg.Wait()
}
