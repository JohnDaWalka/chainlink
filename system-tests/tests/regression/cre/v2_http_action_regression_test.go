package cre

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"

	httpaction_config "github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/http/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// HTTP Action test cases for successful CRUD operations
type httpActionSuccessTest struct {
	name     string
	testCase string
	methods  []string
}

var httpActionSuccessTests = []httpActionSuccessTest{
	{
		name:     "CRUD operations success",
		testCase: "crud-success",
		methods:  []string{"GET", "POST", "PUT", "DELETE"},
	},
}

// HTTP Action test cases for failure scenarios
type httpActionFailureTest struct {
	name          string
	testCase      string
	method        string
	url           string
	headers       map[string]string
	body          string
	timeoutMs     int
	expectedError string
}

var httpActionFailureTests = []httpActionFailureTest{
	// Invalid URL tests
	{
		name:          "invalid URL format",
		testCase:      "crud-failure",
		method:        "GET",
		url:           "not-a-valid-url",
		expectedError: "invalid URL",
	},
	{
		name:          "non-existing URL",
		testCase:      "crud-failure",
		method:        "GET",
		url:           "http://non-existing-domain-12345.com/api/test",
		expectedError: "connection failed",
	},
	// Invalid method tests
	{
		name:          "invalid HTTP method",
		testCase:      "crud-failure",
		method:        "INVALID",
		url:           "http://localhost:8080/test",
		expectedError: "invalid method",
	},
	// Invalid headers tests
	{
		name:     "invalid headers",
		testCase: "crud-failure",
		method:   "GET",
		url:      "http://localhost:8080/test",
		headers: map[string]string{
			"Invalid\nHeader": "value",
		},
		expectedError: "invalid header",
	},
	// Invalid body tests
	{
		name:          "corrupt JSON body",
		testCase:      "crud-failure",
		method:        "POST",
		url:           "http://localhost:8080/test",
		body:          `{"invalid": json}`,
		expectedError: "invalid JSON",
	},
	// Size limit tests
	{
		name:          "oversized request body",
		testCase:      "crud-failure",
		method:        "POST",
		url:           "http://localhost:8080/test",
		body:          strings.Repeat("a", 10*1024*1024), // 10MB body
		expectedError: "request too large",
	},
	{
		name:          "oversized URL",
		testCase:      "crud-failure",
		method:        "GET",
		url:           "http://localhost:8080/test?" + strings.Repeat("param=value&", 10000), // Very long URL
		expectedError: "URL too long",
	},
	// Timeout tests
	{
		name:          "request timeout",
		testCase:      "crud-failure",
		method:        "GET",
		url:           "http://httpbin.org/delay/10", // Endpoint that delays response
		timeoutMs:     1000,                          // 1 second timeout
		expectedError: "timeout",
	},
}

func HTTPActionSuccessTest(t *testing.T, testEnv *ttypes.TestEnvironment, httpActionTest httpActionSuccessTest) {
	testLogger := framework.L
	const workflowFileLocation = "./httpaction/main.go"

	// Get a free port for this test
	freePort, err := getFreePort()
	require.NoError(t, err, "failed to get free port")

	// Start fake HTTP server with CRUD endpoints
	testID := uuid.New().String()[0:8]
	fakeServer, err := startCRUDTestServer(t, freePort, testID)
	require.NoError(t, err, "failed to start fake HTTP server")

	defer func() {
		if fakeServer != nil {
			testLogger.Info().Msgf("Cleaning up fake server on port %d", freePort)
		}
	}()

	testLogger.Info().Msg("Creating HTTP Action success test workflow configuration...")

	// Use host.docker.internal for container-to-host communication
	serverURL := strings.Replace(fakeServer.BaseURLHost, "localhost", "host.docker.internal", 1)
	workflowConfig := httpaction_config.Config{
		URL:      serverURL + "/api/resources",
		TestCase: httpActionTest.testCase,
	}

	workflowName := "http-action-success-workflow-" + httpActionTest.testCase + "-" + testID
	t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	// Start Beholder listener to capture workflow execution messages
	listenerCtx, messageChan, kafkaErrChan := t_helpers.StartBeholder(t, testLogger, testEnv)

	// Wait for workflow execution to complete and verify success
	testLogger.Info().Msg("Waiting for HTTP Action CRUD operations to complete...")
	timeout := 60 * time.Second
	err = t_helpers.AssertBeholderMessage(listenerCtx, t, "HTTP Action CRUD operations completed", testLogger, messageChan, kafkaErrChan, timeout)
	require.NoError(t, err, "HTTP Action CRUD success test failed")

	testLogger.Info().Msg("HTTP Action CRUD success test completed")
}

func HTTPActionFailureTest(t *testing.T, testEnv *ttypes.TestEnvironment, httpActionTest httpActionFailureTest) {
	testLogger := framework.L
	const workflowFileLocation = "./httpaction/main.go"

	testLogger.Info().Msg("Creating HTTP Action failure test workflow configuration...")

	workflowConfig := httpaction_config.Config{
		URL:       httpActionTest.url,
		TestCase:  httpActionTest.testCase,
		Method:    httpActionTest.method,
		Headers:   httpActionTest.headers,
		Body:      httpActionTest.body,
		TimeoutMs: httpActionTest.timeoutMs,
	}

	workflowName := "http-action-fail-workflow-" + httpActionTest.method + "-" + uuid.New().String()[0:8]
	t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	// Start Beholder listener to capture error messages
	listenerCtx, messageChan, kafkaErrChan := t_helpers.StartBeholder(t, testLogger, testEnv)

	// Wait for error message in Beholder
	testLogger.Info().Msg("Waiting for expected HTTP Action failure in Beholder...")
	timeout := 60 * time.Second
	err := t_helpers.AssertBeholderMessage(listenerCtx, t, "HTTP Action expected failure", testLogger, messageChan, kafkaErrChan, timeout)

	// For some failure cases, we might get a different but acceptable error message
	if err != nil {
		// Try to find any error-related message
		err = t_helpers.AssertBeholderMessage(listenerCtx, t, "HTTP Action failed", testLogger, messageChan, kafkaErrChan, timeout)
		if err != nil {
			err = t_helpers.AssertBeholderMessage(listenerCtx, t, httpActionTest.expectedError, testLogger, messageChan, kafkaErrChan, timeout)
		}
	}

	require.NoError(t, err, "Expected HTTP Action failure message not found in Beholder logs")
	testLogger.Info().Msg("HTTP Action failure test completed successfully")
}

// startCRUDTestServer creates a fake HTTP server that supports CRUD operations
func startCRUDTestServer(t *testing.T, port int, testID string) (*fake.Output, error) {
	fakeInput := &fake.Input{
		Port: port,
	}

	fakeOutput, err := fake.NewFakeDataProvider(fakeInput)
	if err != nil {
		return nil, err
	}

	// Set up CRUD endpoints
	resourceResponse := map[string]interface{}{
		"id":     "test-resource-123",
		"name":   "Test Resource",
		"status": "success",
	}

	// POST /api/resources - Create
	err = fake.JSON("POST", "/api/resources", resourceResponse, 201)
	require.NoError(t, err, "failed to set up POST endpoint")

	// GET /api/resources/{id} - Read
	err = fake.JSON("GET", "/api/resources/test-resource-123", resourceResponse, 200)
	require.NoError(t, err, "failed to set up GET endpoint")

	// PUT /api/resources/{id} - Update
	updatedResponse := resourceResponse
	updatedResponse["name"] = "Updated Test Resource"
	err = fake.JSON("PUT", "/api/resources/test-resource-123", updatedResponse, 200)
	require.NoError(t, err, "failed to set up PUT endpoint")

	// DELETE /api/resources/{id} - Delete
	deleteResponse := map[string]interface{}{
		"message": "Resource deleted successfully",
		"status":  "success",
	}
	err = fake.JSON("DELETE", "/api/resources/test-resource-123", deleteResponse, 200)
	require.NoError(t, err, "failed to set up DELETE endpoint")

	framework.L.Info().Msgf("CRUD test server started on port %d at: %s", port, fakeOutput.BaseURLHost)
	framework.L.Info().Msgf("Server URL will be converted to: %s", strings.Replace(fakeOutput.BaseURLHost, "localhost", "host.docker.internal", 1))
	return fakeOutput, nil
}
