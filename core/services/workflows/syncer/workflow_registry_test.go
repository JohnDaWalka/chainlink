package syncer

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	errChanClosed = "channel closed"
	errNoEvent    = "no event"
)

type mockService struct{}

func (m *mockService) Start(context.Context) error { return nil }

func (m *mockService) Close() error { return nil }

func (m *mockService) HealthReport() map[string]error { return map[string]error{"svc": nil} }

func (m *mockService) Ready() error { return nil }

func (m *mockService) Name() string { return "svc" }

func readEventCh(ch chan Event) (Event, error) {
	select {
	case event, ok := <-ch:
		if ok {
			return event, nil
		} else {
			return nil, errors.New(errChanClosed)
		}
	default:
		return nil, errors.New(errNoEvent)
	}
}

func Test_workflowMetadataToEvents(t *testing.T) {
	t.Run("WorkflowRegisteredEvent", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		donID := uint32(1)
		workflowDonNotifier := capabilities.NewDonNotifier()
		// No engines are in the workflow registry
		er := NewEngineRegistry()
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (ContractReader, error) {
				return nil, nil
			},
			"",
			WorkflowEventPollerConfig{
				QueryCount: 20,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		// Capture messages on event channel
		eventCh := make(chan Event, 100)
		wr.eventCh = eventCh

		wfID := [32]byte{1}
		owner := []byte{}
		status := uint8(0)
		wfName := "wf name 1"
		binaryURL := "b1"
		configURL := "c1"
		secretsURL := "s1"
		metadata := []GetWorkflowMetadata{
			{
				WorkflowID:   wfID,
				Owner:        owner,
				DonID:        donID,
				Status:       status,
				WorkflowName: wfName,
				BinaryURL:    binaryURL,
				ConfigURL:    configURL,
				SecretsURL:   secretsURL,
			},
		}

		wr.workflowMetadataToEvents(ctx, metadata, donID)

		// The first event is WorkflowRegisteredEvent
		event, err := readEventCh(eventCh)
		require.NoError(t, err)
		require.Equal(t, WorkflowRegisteredEvent, event.GetEventType())
		expectedRegisteredEvent := WorkflowRegistryWorkflowRegisteredV1{
			WorkflowID:    wfID,
			WorkflowOwner: owner,
			DonID:         donID,
			Status:        status,
			WorkflowName:  wfName,
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
			SecretsURL:    secretsURL,
		}
		require.Equal(t, expectedRegisteredEvent, event.GetData())

		// No other events on the channel
		_, err = readEventCh(eventCh)
		require.ErrorContains(t, err, errNoEvent)
	})

	t.Run("WorkflowUpdatedEvent", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		donID := uint32(1)
		workflowDonNotifier := capabilities.NewDonNotifier()
		// Engine already in the workflow registry
		er := NewEngineRegistry()
		wfID := [32]byte{1}
		owner := []byte{}
		wfName := "wf name 1"
		er.Add(EngineRegistryKey{Owner: owner, Name: wfName}, &mockService{}, wfID)
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (ContractReader, error) {
				return nil, nil
			},
			"",
			WorkflowEventPollerConfig{
				QueryCount: 20,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		// Capture messages on event channel
		eventCh := make(chan Event, 100)
		wr.eventCh = eventCh

		// The workflow metadata gets updated
		wfID2 := [32]byte{2}
		status := uint8(0)
		binaryURL2 := "b2"
		configURL := "c1"
		secretsURL := "s1"
		metadata := []GetWorkflowMetadata{
			{
				WorkflowID:   wfID2,
				Owner:        owner,
				DonID:        donID,
				Status:       status,
				WorkflowName: wfName,
				BinaryURL:    binaryURL2,
				ConfigURL:    configURL,
				SecretsURL:   secretsURL,
			},
		}

		wr.workflowMetadataToEvents(ctx, metadata, donID)

		// The first event is WorkflowUpdatedEvent
		event, err := readEventCh(eventCh)
		require.NoError(t, err)
		require.Equal(t, WorkflowUpdatedEvent, event.GetEventType())
		expectedUpdatedEvent := WorkflowRegistryWorkflowUpdatedV1{
			OldWorkflowID: wfID,
			NewWorkflowID: wfID2,
			WorkflowOwner: owner,
			DonID:         donID,
			WorkflowName:  wfName,
			BinaryURL:     binaryURL2,
			ConfigURL:     configURL,
			SecretsURL:    secretsURL,
		}
		require.Equal(t, expectedUpdatedEvent, event.GetData())

		// No other events on the channel
		event, err = readEventCh(eventCh)
		require.ErrorContains(t, err, errNoEvent)
	})

	t.Run("WorkflowDeletedEvent", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		donID := uint32(1)
		workflowDonNotifier := capabilities.NewDonNotifier()
		// Engine already in the workflow registry
		er := NewEngineRegistry()
		wfID := [32]byte{1}
		owner := []byte{}
		wfName := "wf name 1"
		er.Add(EngineRegistryKey{Owner: owner, Name: wfName}, &mockService{}, wfID)
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (ContractReader, error) {
				return nil, nil
			},
			"",
			WorkflowEventPollerConfig{
				QueryCount: 20,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		// Capture messages on event channel
		eventCh := make(chan Event, 100)
		wr.eventCh = eventCh

		// The workflow metadata is empty
		metadata := []GetWorkflowMetadata{}

		wr.workflowMetadataToEvents(ctx, metadata, donID)

		// The first event is WorkflowDeletedEvent
		event, err := readEventCh(eventCh)
		require.NoError(t, err)
		require.Equal(t, WorkflowDeletedEvent, event.GetEventType())
		expectedUpdatedEvent := WorkflowRegistryWorkflowDeletedV1{
			WorkflowID:    wfID,
			WorkflowOwner: owner,
			DonID:         donID,
			WorkflowName:  wfName,
		}
		require.Equal(t, expectedUpdatedEvent, event.GetData())

		// No other events on the channel
		event, err = readEventCh(eventCh)
		require.ErrorContains(t, err, errNoEvent)
	})

	t.Run("No change", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		donID := uint32(1)
		workflowDonNotifier := capabilities.NewDonNotifier()
		// No engines are in the workflow registry
		er := NewEngineRegistry()
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (ContractReader, error) {
				return nil, nil
			},
			"",
			WorkflowEventPollerConfig{
				QueryCount: 20,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		// Capture messages on event channel
		eventCh := make(chan Event, 100)
		wr.eventCh = eventCh

		wfID := [32]byte{1}
		owner := []byte{}
		status := uint8(0)
		wfName := "wf name 1"
		binaryURL := "b1"
		configURL := "c1"
		secretsURL := "s1"
		metadata := []GetWorkflowMetadata{
			{
				WorkflowID:   wfID,
				Owner:        owner,
				DonID:        donID,
				Status:       status,
				WorkflowName: wfName,
				BinaryURL:    binaryURL,
				ConfigURL:    configURL,
				SecretsURL:   secretsURL,
			},
		}

		wr.workflowMetadataToEvents(ctx, metadata, donID)

		// The first event is WorkflowRegisteredEvent
		event, err := readEventCh(eventCh)
		require.NoError(t, err)
		require.Equal(t, WorkflowRegisteredEvent, event.GetEventType())
		expectedRegisteredEvent := WorkflowRegistryWorkflowRegisteredV1{
			WorkflowID:    wfID,
			WorkflowOwner: owner,
			DonID:         donID,
			Status:        status,
			WorkflowName:  wfName,
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
			SecretsURL:    secretsURL,
		}
		require.Equal(t, expectedRegisteredEvent, event.GetData())

		// No other events on the channel
		_, err = readEventCh(eventCh)
		require.ErrorContains(t, err, errNoEvent)

		// Add the workflow to the engine registry as the handler would
		er.Add(EngineRegistryKey{Owner: owner, Name: wfName}, &mockService{}, wfID)

		// Repeated ticks do not make any new events
		wr.workflowMetadataToEvents(ctx, metadata, donID)
		_, err = readEventCh(eventCh)
		require.ErrorContains(t, err, errNoEvent)
	})
}
