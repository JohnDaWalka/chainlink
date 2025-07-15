package eametricsreporter

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"encoding/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	bridgeMocks "github.com/smartcontractkit/chainlink/v2/core/bridges/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/eametricsreporter/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

// Test constants and fixtures
const (
	testMetricsPath     = "/metrics"
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

// loadTestEAStatusResponse loads the test fixture from JSON using existing loadFixture helper
func loadFixtureAsEAStatusResponse(t *testing.T, filename string) EAStatusResponse {
	fixtureData := loadFixture(t, filename)

	var status EAStatusResponse
	err := json.Unmarshal([]byte(fixtureData), &status)
	require.NoError(t, err, "Failed to unmarshal test fixture")

	return status
}

// Helper function to create WebURL for testing
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
		Name: bridges.BridgeName(testBridgeName1),
		URL:  parseWebURL(testBridgeURL1),
	}
	testBridge2 = bridges.BridgeType{
		Name: bridges.BridgeName(testBridgeName2),
		URL:  parseWebURL(testBridgeURL2),
	}
	testBridges = []bridges.BridgeType{testBridge1, testBridge2}
)

// setupTestService creates a test service with the given configuration overrides
func setupTestService(t *testing.T, enabled bool, pollingInterval time.Duration, httpClient *http.Client) (*Service, *bridgeMocks.ORM, *mocks.MessageEmitter) {
	t.Helper()

	eaConfig := mocks.NewTestEAMetricsReporterConfig(enabled, testMetricsPath, pollingInterval)

	bridgeORM := bridgeMocks.NewORM(t)
	emitter := mocks.NewMessageEmitter()
	lggr := logger.TestLogger(t)

	service := NewEaMetricsReporter(eaConfig, bridgeORM, httpClient, emitter, lggr)

	return service, bridgeORM, emitter
}

func TestNewEaMetricsReporter(t *testing.T) {
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

func TestService_refreshBridgeURLs_NoBridges(t *testing.T) {
	httpClient := &http.Client{}
	service, bridgeORM, _ := setupTestService(t, true, testPollingInterval, httpClient)

	bridgeORM.On("BridgeTypes", mock.Anything, 0, 1000).Return([]bridges.BridgeType{}, 0, nil)

	ctx := context.Background()
	err := service.refreshBridgeURLs(ctx)
	assert.NoError(t, err)

	service.mu.RLock()
	assert.Empty(t, service.bridgeURLs)
	service.mu.RUnlock()
}

func TestService_refreshBridgeURLs_WithBridges(t *testing.T) {
	httpClient := &http.Client{}
	service, bridgeORM, _ := setupTestService(t, true, testPollingInterval, httpClient)

	bridgeORM.On("BridgeTypes", mock.Anything, 0, 1000).Return(testBridges, len(testBridges), nil)

	ctx := context.Background()
	err := service.refreshBridgeURLs(ctx)
	assert.NoError(t, err)

	service.mu.RLock()
	assert.Len(t, service.bridgeURLs, 2)
	assert.Equal(t, testBridgeURL1, service.bridgeURLs[testBridgeName1])
	assert.Equal(t, testBridgeURL2, service.bridgeURLs[testBridgeName2])
	service.mu.RUnlock()
}

func TestService_refreshBridgeURLs_Error(t *testing.T) {
	httpClient := &http.Client{}
	service, bridgeORM, _ := setupTestService(t, true, testPollingInterval, httpClient)

	bridgeORM.On("BridgeTypes", mock.Anything, 0, 1000).Return([]bridges.BridgeType{}, 0, assert.AnError)

	ctx := context.Background()

	// Should handle bridge ORM error gracefully (no panic)
	var err error
	assert.NotPanics(t, func() {
		err = service.refreshBridgeURLs(ctx)
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch bridges")
}

func TestService_pollBridge_Success(t *testing.T) {
	httpClient := mocks.NewMockHTTPClient(loadFixture(t, "ea_status_response.json"), http.StatusOK)
	service, _, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Setup emitter mock
	emitter.On("With", mock.Anything).Return(emitter)
	emitter.On("Emit", mock.Anything, mock.Anything).Return(nil)

	ctx := context.Background()
	service.pollBridge(ctx, "test-bridge", "http://example.com")

	emitter.AssertExpectations(t)
}

func TestService_pollBridge_HTTPError(t *testing.T) {
	httpClient := &http.Client{} // Real client will fail to connect
	service, _, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Expect no emitter calls on HTTP error
	emitter.On("With", mock.Anything).Return(emitter).Maybe()
	emitter.On("Emit", mock.Anything, mock.Anything).Return(nil).Maybe()

	ctx := context.Background()

	// Should handle HTTP error gracefully (no panic, no emission)
	assert.NotPanics(t, func() {
		service.pollBridge(ctx, "test-bridge", "http://127.0.0.1:8080")
	})

	// The emitter should not be called since the HTTP request fails
	emitter.AssertNotCalled(t, "Emit", mock.Anything, mock.Anything)
}

func TestService_pollBridge_InvalidJSON(t *testing.T) {
	httpClient := mocks.NewMockHTTPClient("invalid json", http.StatusOK)
	service, _, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Expect no emitter calls on JSON parse error
	emitter.On("With", mock.Anything).Return(emitter).Maybe()
	emitter.On("Emit", mock.Anything, mock.Anything).Return(nil).Maybe()

	ctx := context.Background()

	// Should handle invalid JSON gracefully (no panic, no emission)
	assert.NotPanics(t, func() {
		service.pollBridge(ctx, "test-bridge", "http://example.com")
	})

	// The emitter should not be called since JSON parsing fails
	emitter.AssertNotCalled(t, "Emit", mock.Anything, mock.Anything)
}

func TestService_pollBridge_InvalidURL(t *testing.T) {
	httpClient := &http.Client{}
	service, _, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Expect no emitter calls on invalid URL
	emitter.On("With", mock.Anything).Return(emitter).Maybe()
	emitter.On("Emit", mock.Anything, mock.Anything).Return(nil).Maybe()

	ctx := context.Background()

	// Should handle invalid URL gracefully (no panic, no emission)
	assert.NotPanics(t, func() {
		service.pollBridge(ctx, "test-bridge", "://invalid-url")
	})

	// The emitter should not be called since URL parsing fails
	emitter.AssertNotCalled(t, "Emit", mock.Anything, mock.Anything)
}

func TestService_pollBridge_Timeout(t *testing.T) {
	httpClient := &http.Client{
		Timeout: 1 * time.Millisecond, // Very short timeout
	}
	service, _, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Expect no emitter calls on timeout
	emitter.On("With", mock.Anything).Return(emitter).Maybe()
	emitter.On("Emit", mock.Anything, mock.Anything).Return(nil).Maybe()

	ctx := context.Background()

	// Should handle timeout gracefully (no panic, no emission)
	assert.NotPanics(t, func() {
		service.pollBridge(ctx, "test-bridge", "http://127.0.0.1:8080")
	})

	// Wait longer than the timeout to ensure the request actually times out
	// This prevents flakiness
	time.Sleep(10 * time.Millisecond)

	// The emitter should not be called since HTTP request times out
	emitter.AssertNotCalled(t, "Emit", mock.Anything, mock.Anything)
}

func TestService_pollBridge_Non200Status(t *testing.T) {
	httpClient := mocks.NewMockHTTPClient("Not Found", http.StatusNotFound)
	service, _, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Expect no emitter calls on non-200 status
	emitter.On("With", mock.Anything).Return(emitter).Maybe()
	emitter.On("Emit", mock.Anything, mock.Anything).Return(nil).Maybe()

	ctx := context.Background()

	// Should handle non-200 status gracefully (no panic, no emission)
	assert.NotPanics(t, func() {
		service.pollBridge(ctx, "test-bridge", "http://example.com")
	})

	// The emitter should not be called since status is 404 Not Found
	emitter.AssertNotCalled(t, "Emit", mock.Anything, mock.Anything)
}

func TestService_emitEAMetrics_Success(t *testing.T) {
	httpClient := &http.Client{}
	service, _, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Setup emitter mock
	emitter.On("With", mock.Anything).Return(emitter)
	emitter.On("Emit", mock.Anything, mock.Anything).Return(nil)

	ctx := context.Background()
	service.emitEAMetrics(ctx, "test-bridge", loadFixtureAsEAStatusResponse(t, "ea_status_response.json"))

	emitter.AssertExpectations(t)
}

func TestService_emitEAMetrics_EmitError(t *testing.T) {
	httpClient := &http.Client{}
	service, _, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Setup emitter mock to return error
	emitter.On("With", mock.Anything).Return(emitter)
	emitter.On("Emit", mock.Anything, mock.Anything).Return(assert.AnError)

	ctx := context.Background()

	// This should not panic or return an error, even when emitter fails
	assert.NotPanics(t, func() {
		service.emitEAMetrics(ctx, "test-bridge", loadFixtureAsEAStatusResponse(t, "ea_status_response.json"))
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

	// Track emitted bridge names to verify each bridge was processed
	emittedBridgeNames := []string{}
	emittedMessages := []string{}

	// Setup emitter mock to capture bridge names and messages
	emitter.On("With", mock.AnythingOfType("[]string")).Return(emitter).Run(func(args mock.Arguments) {
		kvs := args.Get(0).([]string)
		for i := 0; i < len(kvs); i += 2 {
			if i+1 < len(kvs) && kvs[i] == "bridge_name" {
				emittedBridgeNames = append(emittedBridgeNames, kvs[i+1])
			}
		}
	})

	emitter.On("Emit", mock.Anything, mock.AnythingOfType("string")).Return(nil).Run(func(args mock.Arguments) {
		emittedMessages = append(emittedMessages, args.Get(1).(string))
	})

	ctx := context.Background()
	service.pollAllBridges(ctx)

	// Verify bridge ORM was called
	bridgeORM.AssertExpectations(t)

	// Verify telemetry was emitted for each bridge
	expectedBridgeNames := []string{testBridgeName1, testBridgeName2}
	assert.ElementsMatch(t, expectedBridgeNames, emittedBridgeNames, "Should emit telemetry for each bridge")

	// Verify correct number of emissions (once per bridge)
	assert.Len(t, emittedMessages, 2, "Should emit telemetry exactly twice")
	for _, msg := range emittedMessages {
		assert.Contains(t, msg, "EA Metrics - Bridge:", "Should emit correct message format")
		assert.Contains(t, msg, "Adapter: test-adapter", "Should include adapter name")
		assert.Contains(t, msg, "Version: 1.0.0", "Should include version")
	}

	emitter.AssertExpectations(t)
}

func TestService_emitEAMetrics_CaptureOutput(t *testing.T) {
	// Create a mock emitter that captures the data
	emitter := mocks.NewMessageEmitter()

	// Mock the With method to capture labels
	capturedLabels := make(map[string]string)
	emitter.On("With", mock.AnythingOfType("[]string")).Return(emitter).Run(func(args mock.Arguments) {
		kvs := args.Get(0).([]string)
		// Process key-value pairs
		for i := 0; i < len(kvs); i += 2 {
			if i+1 < len(kvs) {
				capturedLabels[kvs[i]] = kvs[i+1]
			}
		}
	})

	// Mock the Emit method to capture the message
	capturedMessage := ""
	emitter.On("Emit", mock.Anything, mock.AnythingOfType("string")).Return(nil).Run(func(args mock.Arguments) {
		capturedMessage = args.Get(1).(string)
	})

	// Create service with mock emitter
	config := mocks.NewTestEAMetricsReporterConfig(true, "/metrics", 5*time.Minute)
	service := &Service{
		config:  config,
		emitter: emitter,
		lggr:    logger.TestLogger(t),
	}

	// Call emitEAMetrics with test fixture
	ctx := context.Background()
	status := loadFixtureAsEAStatusResponse(t, "ea_status_response.json")
	service.emitEAMetrics(ctx, "test-bridge", status)

	// Verify the expected message format (based on fixture data)
	expectedMessage := "EA Metrics - Bridge: test-bridge, Adapter: test-adapter, Version: 1.0.0"
	assert.Equal(t, expectedMessage, capturedMessage)

	// Verify key labels are present (based on fixture data)
	assert.Equal(t, "test-bridge", capturedLabels["bridge_name"])
	assert.Equal(t, "test-adapter", capturedLabels["adapter_name"])
	assert.Equal(t, "1.0.0", capturedLabels["adapter_version"])
	assert.Equal(t, "3600", capturedLabels["adapter_uptime_seconds"])
	assert.Equal(t, "linux", capturedLabels["runtime_platform"])
	assert.Equal(t, "x64", capturedLabels["runtime_architecture"])
	assert.Equal(t, "18.17.0", capturedLabels["runtime_node_version"])
	assert.Equal(t, "ea-adapter-01", capturedLabels["runtime_hostname"])
	assert.Equal(t, "true", capturedLabels["metrics_enabled"])

	// Verify JSON data is present
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

	// Verify all mocks were called as expected
	emitter.AssertExpectations(t)
}

func TestService_emitEAMetrics_MissingFields(t *testing.T) {
	httpClient := &http.Client{}
	service, _, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Load the fixture and remove one field to test robustness
	fixtureData := loadFixture(t, "ea_status_response.json")

	var responseData map[string]interface{}
	err := json.Unmarshal([]byte(fixtureData), &responseData)
	require.NoError(t, err)

	// Remove the version field from adapter section
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
	var allFields []string
	emitter.On("With", mock.AnythingOfType("[]string")).Return(emitter).Run(func(args mock.Arguments) {
		kvs := args.Get(0).([]string)
		allFields = append(allFields, kvs...)
	})
	emitter.On("Emit", mock.Anything, mock.AnythingOfType("string")).Return(nil)

	ctx := context.Background()

	// Should emit successfully even with missing field
	service.emitEAMetrics(ctx, "test-bridge", status)

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

func TestService_emitEAMetrics_ExtraFields(t *testing.T) {
	httpClient := &http.Client{}
	service, _, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Load the fixture and add one extra field to test forward compatibility
	fixtureData := loadFixture(t, "ea_status_response.json")

	var responseData map[string]interface{}
	err := json.Unmarshal([]byte(fixtureData), &responseData)
	require.NoError(t, err)

	// Add one extra field to the adapter section
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
	service.emitEAMetrics(ctx, "test-bridge", status)

	emitter.AssertExpectations(t)
}

func TestService_Start_AlreadyStarted(t *testing.T) {
	httpClient := &http.Client{}
	service, _, _ := setupTestService(t, true, testPollingInterval, httpClient)

	ctx := context.Background()

	// Start service first time
	err := service.Start(ctx)
	assert.NoError(t, err)

	// Attempt to start again - should return error or be no-op
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

	// Start and close service
	err := service.Start(ctx)
	assert.NoError(t, err)

	err = service.Close()
	assert.NoError(t, err)

	// Attempt to close again - should return error or be no-op
	err = service.Close()

	// services.StateMachine prevents double close, should return error
	assert.Error(t, err)
}
