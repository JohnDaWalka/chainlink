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
	"github.com/smartcontractkit/chainlink/v2/core/services/eastatusreporter/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
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

	// Track emitted bridge names
	var emittedBridgeNamesMutex sync.Mutex
	emittedBridgeNames := []string{}

	var emittedMessagesMutex sync.Mutex
	emittedMessages := []string{}

	// Setup emitter mock to capture bridge names and messages
	emitter.On("With", mock.AnythingOfType("[]string")).Return(emitter).Run(func(args mock.Arguments) {
		kvs := args.Get(0).([]string)
		for i := 0; i < len(kvs); i += 2 {
			if i+1 < len(kvs) && kvs[i] == "bridge_name" {
				emittedBridgeNamesMutex.Lock()
				emittedBridgeNames = append(emittedBridgeNames, kvs[i+1])
				emittedBridgeNamesMutex.Unlock()
			}
		}
	})

	emitter.On("Emit", mock.Anything, mock.AnythingOfType("string")).Return(nil).Run(func(args mock.Arguments) {
		emittedMessagesMutex.Lock()
		emittedMessages = append(emittedMessages, args.Get(1).(string))
		emittedMessagesMutex.Unlock()
	})

	ctx := context.Background()
	service.pollAllBridges(ctx)

	bridgeORM.AssertExpectations(t)

	expectedBridgeNames := []string{testBridgeName1, testBridgeName2}
	assert.ElementsMatch(t, expectedBridgeNames, emittedBridgeNames, "Should emit telemetry for each bridge")
	assert.Len(t, emittedMessages, 2)
	for _, msg := range emittedMessages {
		assert.Contains(t, msg, "EA Status - Bridge:")
		assert.Contains(t, msg, "Adapter: test-adapter")
		assert.Contains(t, msg, "Version: 1.0.0")
	}

	emitter.AssertExpectations(t)
}

func TestService_emitEAStatus_CaptureOutput(t *testing.T) {
	emitter := mocks.NewMessageEmitter()
	var capturedLabelsMutex sync.Mutex
	capturedLabels := make(map[string]string)

	emitter.On("With", mock.AnythingOfType("[]string")).Return(emitter).Run(func(args mock.Arguments) {
		kvs := args.Get(0).([]string)
		// Process key-value pairs
		capturedLabelsMutex.Lock()
		for i := 0; i < len(kvs); i += 2 {
			if i+1 < len(kvs) {
				capturedLabels[kvs[i]] = kvs[i+1]
			}
		}
		capturedLabelsMutex.Unlock()
	})

	capturedMessage := ""
	emitter.On("Emit", mock.Anything, mock.AnythingOfType("string")).Return(nil).Run(func(args mock.Arguments) {
		capturedMessage = args.Get(1).(string)
	})
	config := mocks.NewTestEAStatusReporterConfig(true, "/status", 5*time.Minute)
	service := &Service{
		config:  config,
		emitter: emitter,
		lggr:    logger.TestLogger(t),
	}

	// Call emitEAStatus with test fixture
	ctx := context.Background()
	status := loadFixtureAsEAStatusResponse(t, "ea_status_response.json")
	service.emitEAStatus(ctx, "test-bridge", status)

	expectedMessage := "EA Status - Bridge: test-bridge, Adapter: test-adapter, Version: 1.0.0"
	assert.Equal(t, expectedMessage, capturedMessage)
	assert.Equal(t, "test-bridge", capturedLabels["bridge_name"])
	assert.Equal(t, "test-adapter", capturedLabels["adapter_name"])
	assert.Equal(t, "1.0.0", capturedLabels["adapter_version"])
	assert.Equal(t, "3600", capturedLabels["adapter_uptime_seconds"])
	assert.Equal(t, "linux", capturedLabels["runtime_platform"])
	assert.Equal(t, "x64", capturedLabels["runtime_architecture"])
	assert.Equal(t, "18.17.0", capturedLabels["runtime_node_version"])
	assert.Equal(t, "ea-adapter-01", capturedLabels["runtime_hostname"])
	assert.Equal(t, "true", capturedLabels["metrics_enabled"])

	assert.Contains(t, capturedLabels, "endpoints")
	assert.Contains(t, capturedLabels, "configuration")

	// Parse and verify JSON endpoints
	var endpoints []struct {
		Name       string   `json:"name"`
		Aliases    []string `json:"aliases"`
		Transports []string `json:"transports"`
	}
	err := json.Unmarshal([]byte(capturedLabels["endpoints"]), &endpoints)
	assert.NoError(t, err)
	assert.Len(t, endpoints, 2)
	assert.Equal(t, "price", endpoints[0].Name)
	assert.Equal(t, "volume", endpoints[1].Name)

	// Parse and verify JSON configurations
	var configuration []struct {
		Name               string      `json:"name"`
		Value              interface{} `json:"value"`
		Type               string      `json:"type"`
		Description        string      `json:"description"`
		Required           bool        `json:"required"`
		Default            interface{} `json:"default"`
		CustomSetting      bool        `json:"customSetting"`
		EnvDefaultOverride interface{} `json:"envDefaultOverride"`
	}
	err = json.Unmarshal([]byte(capturedLabels["configuration"]), &configuration)
	assert.NoError(t, err)
	assert.Len(t, configuration, 3)
	assert.Equal(t, "API_KEY", configuration[0].Name)
	assert.Equal(t, "TIMEOUT", configuration[1].Name)
	assert.Equal(t, "CACHE_ENABLED", configuration[2].Name)

	emitter.AssertExpectations(t)
}

func TestService_emitEAStatus_MissingFields(t *testing.T) {
	httpClient := &http.Client{}
	service, _, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Load the fixture and remove one field to test robustness
	fixtureData := loadFixture(t, "ea_status_response.json")

	var responseData map[string]interface{}
	err := json.Unmarshal([]byte(fixtureData), &responseData)
	require.NoError(t, err)

	if adapter, ok := responseData["adapter"].(map[string]interface{}); ok {
		delete(adapter, "version")
	}

	// Marshal back to JSON and unmarshal into our struct
	modifiedJSON, err := json.Marshal(responseData)
	require.NoError(t, err)

	var status EAStatusResponse
	err = json.Unmarshal(modifiedJSON, &status)
	require.NoError(t, err)

	// Setup emitter mock - collect all fields from all With() calls
	var allFieldsMutex sync.Mutex
	var allFields []string

	emitter.On("With", mock.AnythingOfType("[]string")).Return(emitter).Run(func(args mock.Arguments) {
		kvs := args.Get(0).([]string)
		allFieldsMutex.Lock()
		allFields = append(allFields, kvs...)
		allFieldsMutex.Unlock()
	})

	emitter.On("Emit", mock.Anything, mock.AnythingOfType("string")).Return(nil)

	ctx := context.Background()

	// Should emit successfully even with missing field
	service.emitEAStatus(ctx, "test-bridge", status)

	// Verify adapter_version field is present and empty
	assert.Contains(t, allFields, "adapter_version", "Missing adapter_version field should be present in emission")
	// Find the version value and verify it's empty
	for i := 0; i < len(allFields)-1; i += 2 {
		if allFields[i] == "adapter_version" {
			assert.Equal(t, "", allFields[i+1], "Missing version field should be empty string")
			break
		}
	}

	emitter.AssertExpectations(t)
}

func TestService_emitEAStatus_ExtraFields(t *testing.T) {
	httpClient := &http.Client{}
	service, _, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Load the fixture and add one extra field to test forward compatibility
	fixtureData := loadFixture(t, "ea_status_response.json")

	var responseData map[string]interface{}
	err := json.Unmarshal([]byte(fixtureData), &responseData)
	require.NoError(t, err)

	if adapter, ok := responseData["adapter"].(map[string]interface{}); ok {
		adapter["buildNumber"] = "12345"
	}

	// Marshal back to JSON and unmarshal into our struct - extra field should be ignored
	modifiedJSON, err := json.Marshal(responseData)
	require.NoError(t, err)

	var status EAStatusResponse
	err = json.Unmarshal(modifiedJSON, &status)
	require.NoError(t, err)

	// Setup emitter mock - should not include the extra field
	emitter.On("With", mock.AnythingOfType("[]string")).Return(emitter).Run(func(args mock.Arguments) {
		kvs := args.Get(0).([]string)
		// Extra field (buildNumber) should not be in any emission
		assert.NotContains(t, kvs, "buildNumber")
	})
	emitter.On("Emit", mock.Anything, mock.AnythingOfType("string")).Return(nil)

	ctx := context.Background()

	// Should emit successfully, ignoring extra field
	service.emitEAStatus(ctx, "test-bridge", status)

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
