package eastatusreporter

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"encoding/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	bridgeMocks "github.com/smartcontractkit/chainlink/v2/core/bridges/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/eastatusreporter/events"
	"github.com/smartcontractkit/chainlink/v2/core/services/eastatusreporter/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"google.golang.org/protobuf/proto"
)

// Test constants and fixtures
const (
	testStatusPath      = "/status"
	testPollingInterval = 5 * time.Minute
	testBridgeName1     = "bridge1"
	testBridgeName2     = "bridge2"
	testBridgeURL1      = "http://bridge1.example.com"
	testBridgeURL2      = "http://bridge2.example.com"
)

// loadFixture loads a JSON fixture file
func loadFixture(t *testing.T, filename string) string {
	t.Helper()

	fixturePath := filepath.Join("fixtures", filename)
	data, err := os.ReadFile(fixturePath)
	require.NoError(t, err, "Failed to read fixture file: %s", fixturePath)

	return string(data)
}

// loadFixtureAsEAStatusResponse loads and unmarshals fixture data
func loadFixtureAsEAStatusResponse(t *testing.T, filename string) EAStatusResponse {
	fixtureData := loadFixture(t, filename)

	var status EAStatusResponse
	err := json.Unmarshal([]byte(fixtureData), &status)
	require.NoError(t, err, "Failed to unmarshal test fixture")

	return status
}

// parseWebURL creates WebURL from string
func parseWebURL(s string) models.WebURL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return models.WebURL(*u)
}

// Test fixtures
var (
	testBridge1 = bridges.BridgeType{
		Name: bridges.MustParseBridgeName(testBridgeName1),
		URL:  parseWebURL(testBridgeURL1),
	}
	testBridge2 = bridges.BridgeType{
		Name: bridges.MustParseBridgeName(testBridgeName2),
		URL:  parseWebURL(testBridgeURL2),
	}
	testBridges = []bridges.BridgeType{testBridge1, testBridge2}
)

// setupTestService creates a test service with mocks
func setupTestService(t *testing.T, enabled bool, pollingInterval time.Duration, httpClient *http.Client) (*Service, *bridgeMocks.ORM, *mocks.MessageEmitter) {
	t.Helper()

	eaConfig := mocks.NewTestEAStatusReporterConfig(enabled, testStatusPath, pollingInterval)

	bridgeORM := bridgeMocks.NewORM(t)
	emitter := mocks.NewMessageEmitter()
	lggr := logger.TestLogger(t)

	// Reduce log noise
	lggr.SetLogLevel(zapcore.ErrorLevel)

	service := NewEaStatusReporter(eaConfig, bridgeORM, httpClient, emitter, lggr)

	return service, bridgeORM, emitter
}

func TestNewEaStatusReporter(t *testing.T) {
	httpClient := &http.Client{}
	service, _, _ := setupTestService(t, true, testPollingInterval, httpClient)

	assert.NotNil(t, service)
	assert.Equal(t, ServiceName, service.Name())
}

func TestService_Start_Disabled(t *testing.T) {
	httpClient := &http.Client{}
	service, _, _ := setupTestService(t, false, testPollingInterval, httpClient)

	ctx := context.Background()
	err := service.Start(ctx)
	assert.NoError(t, err)

	err = service.Close()
	assert.NoError(t, err)
}

func TestService_Start_Enabled(t *testing.T) {
	httpClient := &http.Client{}
	service, _, _ := setupTestService(t, true, 100*time.Millisecond, httpClient)

	ctx := context.Background()
	err := service.Start(ctx)
	assert.NoError(t, err)

	err = service.Close()
	assert.NoError(t, err)
}

func TestService_HealthReport(t *testing.T) {
	httpClient := &http.Client{}
	service, _, _ := setupTestService(t, true, testPollingInterval, httpClient)

	health := service.HealthReport()
	assert.Contains(t, health, service.Name())
}

func TestService_pollAllBridges_NoBridges(t *testing.T) {
	httpClient := &http.Client{}
	service, bridgeORM, _ := setupTestService(t, true, testPollingInterval, httpClient)

	bridgeORM.On("BridgeTypes", mock.Anything, 0, 1000).Return([]bridges.BridgeType{}, 0, nil)

	ctx := context.Background()

	// Should handle no bridges gracefully
	assert.NotPanics(t, func() {
		service.pollAllBridges(ctx)
	})

	bridgeORM.AssertExpectations(t)
}

func TestService_pollAllBridges_WithBridges(t *testing.T) {
	httpClient := mocks.NewMockHTTPClient(loadFixture(t, "ea_status_response.json"), http.StatusOK)
	service, bridgeORM, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	bridgeORM.On("BridgeTypes", mock.Anything, 0, 1000).Return(testBridges, len(testBridges), nil)

	emitter.On("With", mock.Anything).Return(emitter)
	emitter.On("Emit", mock.Anything, mock.Anything).Return(nil)

	ctx := context.Background()
	service.pollAllBridges(ctx)

	bridgeORM.AssertExpectations(t)
	emitter.AssertExpectations(t)
}

func TestService_pollAllBridges_FetchError(t *testing.T) {
	httpClient := &http.Client{}
	service, bridgeORM, _ := setupTestService(t, true, testPollingInterval, httpClient)

	bridgeORM.On("BridgeTypes", mock.Anything, 0, 1000).Return([]bridges.BridgeType{}, 0, assert.AnError)

	ctx := context.Background()

	// Should handle bridge ORM error gracefully
	assert.NotPanics(t, func() {
		service.pollAllBridges(ctx)
	})

	bridgeORM.AssertExpectations(t)
}

func TestService_pollBridge_Success(t *testing.T) {
	httpClient := mocks.NewMockHTTPClient(loadFixture(t, "ea_status_response.json"), http.StatusOK)
	service, _, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	emitter.On("With", mock.Anything).Return(emitter)
	emitter.On("Emit", mock.Anything, mock.Anything).Return(nil)

	ctx := context.Background()
	service.pollBridge(ctx, "test-bridge", "http://example.com")

	emitter.AssertExpectations(t)
}

func TestService_pollBridge_HTTPError(t *testing.T) {
	httpClient := &http.Client{}
	service, _, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	emitter.On("With", mock.Anything).Return(emitter).Maybe()
	emitter.On("Emit", mock.Anything, mock.Anything).Return(nil).Maybe()

	ctx := context.Background()

	// Should handle HTTP error gracefully
	assert.NotPanics(t, func() {
		service.pollBridge(ctx, "test-bridge", "http://127.0.0.1:8080")
	})

	emitter.AssertNotCalled(t, "Emit", mock.Anything, mock.Anything)
}

func TestService_pollBridge_InvalidJSON(t *testing.T) {
	httpClient := mocks.NewMockHTTPClient("invalid json", http.StatusOK)
	service, _, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	emitter.On("With", mock.Anything).Return(emitter).Maybe()
	emitter.On("Emit", mock.Anything, mock.Anything).Return(nil).Maybe()

	ctx := context.Background()

	assert.NotPanics(t, func() {
		service.pollBridge(ctx, "test-bridge", "http://example.com")
	})
	emitter.AssertNotCalled(t, "Emit", mock.Anything, mock.Anything)
}

func TestService_pollBridge_InvalidURL(t *testing.T) {
	httpClient := &http.Client{}
	service, _, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	emitter.On("With", mock.Anything).Return(emitter).Maybe()
	emitter.On("Emit", mock.Anything, mock.Anything).Return(nil).Maybe()

	ctx := context.Background()

	assert.NotPanics(t, func() {
		service.pollBridge(ctx, "test-bridge", "://invalid-url")
	})
	emitter.AssertNotCalled(t, "Emit", mock.Anything, mock.Anything)
}

func TestService_pollBridge_Non200Status(t *testing.T) {
	httpClient := mocks.NewMockHTTPClient("Not Found", http.StatusNotFound)
	service, _, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	emitter.On("With", mock.Anything).Return(emitter).Maybe()
	emitter.On("Emit", mock.Anything, mock.Anything).Return(nil).Maybe()

	ctx := context.Background()

	assert.NotPanics(t, func() {
		service.pollBridge(ctx, "test-bridge", "http://example.com")
	})
	emitter.AssertNotCalled(t, "Emit", mock.Anything, mock.Anything)
}

func TestService_emitEAStatus_Success(t *testing.T) {
	httpClient := &http.Client{}
	service, _, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	emitter.On("With", mock.Anything).Return(emitter)
	emitter.On("Emit", mock.Anything, mock.Anything).Return(nil)

	ctx := context.Background()
	service.emitEAStatus(ctx, "test-bridge", loadFixtureAsEAStatusResponse(t, "ea_status_response.json"))

	emitter.AssertExpectations(t)
}

func TestService_emitEAStatus_EmitError(t *testing.T) {
	httpClient := &http.Client{}
	service, _, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	emitter.On("With", mock.Anything).Return(emitter)
	emitter.On("Emit", mock.Anything, mock.Anything).Return(assert.AnError)

	ctx := context.Background()
	assert.NotPanics(t, func() {
		service.emitEAStatus(ctx, "test-bridge", loadFixtureAsEAStatusResponse(t, "ea_status_response.json"))
	})

	emitter.AssertExpectations(t)
}

func TestService_pollAllBridges_RefreshError(t *testing.T) {
	httpClient := &http.Client{}
	service, bridgeORM, _ := setupTestService(t, true, testPollingInterval, httpClient)

	// Setup bridge ORM mock to return error
	bridgeORM.On("BridgeTypes", mock.Anything, 0, 1000).Return([]bridges.BridgeType{}, 0, assert.AnError)

	ctx := context.Background()

	// Should handle bridge refresh error gracefully (no panic)
	assert.NotPanics(t, func() {
		service.pollAllBridges(ctx)
	})

	bridgeORM.AssertExpectations(t)
}

func TestService_pollAllBridges_MultipleBridges(t *testing.T) {
	httpClient := mocks.NewMockHTTPClient(loadFixture(t, "ea_status_response.json"), http.StatusOK)
	service, bridgeORM, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Setup bridge ORM mock to return our test bridges
	bridgeORM.On("BridgeTypes", mock.Anything, 0, 1000).Return(testBridges, len(testBridges), nil)

	// Track emitted bridge names from protobuf events
	var emittedBridgeNamesMutex sync.Mutex
	emittedBridgeNames := []string{}

	// Setup emitter mock to capture protobuf events and extract bridge names
	emitter.On("With", mock.AnythingOfType("[]string")).Return(emitter)
	emitter.On("Emit", mock.Anything, mock.AnythingOfType("string")).Return(nil).Run(func(args mock.Arguments) {
		// Unmarshal protobuf to extract bridge name
		protobufBytes := []byte(args.Get(1).(string))
		var event events.EAStatusEvent
		if err := proto.Unmarshal(protobufBytes, &event); err == nil {
			emittedBridgeNamesMutex.Lock()
			emittedBridgeNames = append(emittedBridgeNames, event.BridgeName)
			emittedBridgeNamesMutex.Unlock()
		}
	})

	ctx := context.Background()
	service.pollAllBridges(ctx)

	bridgeORM.AssertExpectations(t)

	// Verify we emitted events for both bridges
	expectedBridgeNames := []string{testBridgeName1, testBridgeName2}
	assert.ElementsMatch(t, expectedBridgeNames, emittedBridgeNames, "Should emit telemetry for each bridge")

	emitter.AssertExpectations(t)
}

func TestService_emitEAStatus_CaptureOutput(t *testing.T) {
	emitter := mocks.NewMessageEmitter()
	var capturedProtobufBytes []byte

	// Capture protobuf metadata labels
	emitter.On("With", mock.AnythingOfType("[]string")).Return(emitter)
	emitter.On("Emit", mock.Anything, mock.AnythingOfType("string")).Return(nil).Run(func(args mock.Arguments) {
		capturedProtobufBytes = []byte(args.Get(1).(string))
	})

	config := mocks.NewTestEAStatusReporterConfig(true, "/status", 5*time.Minute)
	service := &Service{
		config:  config,
		emitter: emitter,
		lggr:    logger.TestLogger(t),
	}

	// Load fixture and emit
	ctx := context.Background()
	status := loadFixtureAsEAStatusResponse(t, "ea_status_response.json")
	service.emitEAStatus(ctx, "test-bridge", status)

	// Unmarshal and verify protobuf matches fixture values
	require.NotEmpty(t, capturedProtobufBytes)
	var event events.EAStatusEvent
	err := proto.Unmarshal(capturedProtobufBytes, &event)
	require.NoError(t, err)

	// Verify key fields match fixture
	assert.Equal(t, "test-bridge", event.BridgeName)
	assert.Equal(t, status.Adapter.Name, event.AdapterName)
	assert.Equal(t, status.Adapter.Version, event.AdapterVersion)
	assert.Equal(t, status.Adapter.UptimeSeconds, event.AdapterUptimeSeconds)

	//Verify Endpoints
	for i, endpoint := range status.Endpoints {
		assert.Equal(t, endpoint.Name, event.Endpoints[i].Name)
		assert.Equal(t, endpoint.Aliases, event.Endpoints[i].Aliases)
		assert.Equal(t, endpoint.Transports, event.Endpoints[i].Transports)
	}

	//Verify Default Endpoint
	assert.Equal(t, status.DefaultEndpoint, event.DefaultEndpoint)

	// Verify configuration
	for i, configuration := range status.Configuration {
		assert.Equal(t, configuration.Name, event.Configuration[i].Name)
		assert.Equal(t, fmt.Sprintf("%v", configuration.Value), event.Configuration[i].Value) // Values are converted to strings
		assert.Equal(t, configuration.Type, event.Configuration[i].Type)
		assert.Equal(t, configuration.Description, event.Configuration[i].Description)
		assert.Equal(t, configuration.Required, event.Configuration[i].Required)
		assert.Equal(t, fmt.Sprintf("%v", configuration.Default), event.Configuration[i].DefaultValue) // Defaults converted to strings
		assert.Equal(t, configuration.CustomSetting, event.Configuration[i].CustomSetting)
		assert.Equal(t, fmt.Sprintf("%v", configuration.EnvDefaultOverride), event.Configuration[i].EnvDefaultOverride) // Overrides converted to strings
	}

	//Verify Runtime
	assert.Equal(t, status.Runtime.NodeVersion, event.Runtime.NodeVersion)
	assert.Equal(t, status.Runtime.Platform, event.Runtime.Platform)
	assert.Equal(t, status.Runtime.Architecture, event.Runtime.Architecture)
	assert.Equal(t, status.Runtime.Hostname, event.Runtime.Hostname)

	// Verify Metrics
	assert.Equal(t, status.Metrics.Enabled, event.Metrics.Enabled)

	emitter.AssertExpectations(t)
}

func TestService_Start_AlreadyStarted(t *testing.T) {
	httpClient := &http.Client{}
	service, _, _ := setupTestService(t, true, testPollingInterval, httpClient)

	ctx := context.Background()

	err := service.Start(ctx)
	assert.NoError(t, err)
	err = service.Start(ctx)
	// services.StateMachine prevents double start, should return error
	assert.Error(t, err)

	// Clean up
	err = service.Close()
	assert.NoError(t, err)
}

func TestService_Close_AlreadyClosed(t *testing.T) {
	httpClient := &http.Client{}
	service, _, _ := setupTestService(t, true, testPollingInterval, httpClient)

	ctx := context.Background()

	err := service.Start(ctx)
	assert.NoError(t, err)

	err = service.Close()
	assert.NoError(t, err)
	err = service.Close()

	// services.StateMachine prevents double close, should return error
	assert.Error(t, err)
}

func TestService_PollAllBridges_3000Bridges(t *testing.T) {
	httpClient := mocks.NewMockHTTPClient(loadFixture(t, "ea_status_response.json"), http.StatusOK)
	service, mockORM, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	numBridges := 3000
	var allBridges []bridges.BridgeType
	for i := 0; i < numBridges; i++ {
		u, _ := url.Parse(fmt.Sprintf("http://bridge%d.example.com", i))
		bridge := bridges.BridgeType{
			Name: bridges.MustParseBridgeName(fmt.Sprintf("bridge%d", i)),
			URL:  models.WebURL(*u),
		}
		allBridges = append(allBridges, bridge)
	}

	// Page 1: bridges 0-999 (1000 bridges)
	page1 := allBridges[0:1000]
	mockORM.On("BridgeTypes", mock.Anything, 0, bridgePollPageSize).Return(page1, 3000, nil).Once()

	// Page 2: bridges 1000-1999 (1000 bridges)
	page2 := allBridges[1000:2000]
	mockORM.On("BridgeTypes", mock.Anything, 1000, bridgePollPageSize).Return(page2, 3000, nil).Once()

	// Page 3: bridges 2000-2999 (1000 bridges)
	page3 := allBridges[2000:3000]
	mockORM.On("BridgeTypes", mock.Anything, 2000, bridgePollPageSize).Return(page3, 3000, nil).Once()

	// Page 4: empty (end of results)
	mockORM.On("BridgeTypes", mock.Anything, 3000, bridgePollPageSize).Return([]bridges.BridgeType{}, 3000, nil).Once()

	// Expect 3000 telemetry emissions
	emitter.On("With", mock.Anything).Return(emitter)
	emitter.On("Emit", mock.Anything, mock.AnythingOfType("string")).Return(nil).Times(numBridges)

	ctx := context.Background()

	service.pollAllBridges(ctx)
	mockORM.AssertExpectations(t)
	emitter.AssertExpectations(t)
}

func TestService_PollAllBridges_ContextTimeout(t *testing.T) {
	httpClient := &http.Client{}
	service, mockORM, _ := setupTestService(t, true, testPollingInterval, httpClient)

	numBridges := 5
	var allBridges []bridges.BridgeType
	for i := 0; i < numBridges; i++ {
		u, _ := url.Parse(fmt.Sprintf("http://bridge%d.example.com", i))
		bridge := bridges.BridgeType{
			Name: bridges.MustParseBridgeName(fmt.Sprintf("bridge%d", i)),
			URL:  models.WebURL(*u),
		}
		allBridges = append(allBridges, bridge)
	}

	mockORM.On("BridgeTypes", mock.Anything, 0, bridgePollPageSize).Return(allBridges, numBridges, nil).Once()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Fail(t, "HTTP handler should not complete due to context cancellation")
	}))
	defer server.Close()

	serverURL, _ := url.Parse(server.URL)
	for i := range allBridges {
		allBridges[i].URL = models.WebURL(*serverURL)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	service.pollAllBridges(ctx)
	mockORM.AssertExpectations(t)
}

func TestService_emitEAStatus_EmptyFields(t *testing.T) {
	emitter := mocks.NewMessageEmitter()
	var capturedProtobufBytes []byte

	// Capture protobuf metadata labels
	emitter.On("With", mock.AnythingOfType("[]string")).Return(emitter)
	emitter.On("Emit", mock.Anything, mock.AnythingOfType("string")).Return(nil).Run(func(args mock.Arguments) {
		capturedProtobufBytes = []byte(args.Get(1).(string))
	})

	config := mocks.NewTestEAStatusReporterConfig(true, "/status", 5*time.Minute)
	service := &Service{
		config:  config,
		emitter: emitter,
		lggr:    logger.TestLogger(t),
	}

	// Load empty fixture and emit
	ctx := context.Background()
	status := loadFixtureAsEAStatusResponse(t, "ea_status_empty.json")
	service.emitEAStatus(ctx, "empty-bridge", status)

	// Unmarshal and verify protobuf handles empty values correctly
	require.NotEmpty(t, capturedProtobufBytes)
	var event events.EAStatusEvent
	err := proto.Unmarshal(capturedProtobufBytes, &event)
	require.NoError(t, err)

	// Verify empty/minimal values are handled correctly
	assert.Equal(t, "empty-bridge", event.BridgeName)
	assert.Equal(t, "", event.AdapterName)
	assert.Equal(t, "", event.AdapterVersion)
	assert.Equal(t, int64(0), event.AdapterUptimeSeconds)
	assert.Equal(t, "", event.DefaultEndpoint)

	// Verify empty runtime info
	require.NotNil(t, event.Runtime)
	assert.Equal(t, "", event.Runtime.NodeVersion)
	assert.Equal(t, "", event.Runtime.Platform)
	assert.Equal(t, "", event.Runtime.Architecture)
	assert.Equal(t, "", event.Runtime.Hostname)

	// Verify metrics with false enabled
	require.NotNil(t, event.Metrics)
	assert.False(t, event.Metrics.Enabled)

	// Verify empty arrays
	assert.Empty(t, event.Endpoints)
	assert.Empty(t, event.Configuration)

	emitter.AssertExpectations(t)
}
