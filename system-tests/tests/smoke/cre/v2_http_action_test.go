package cre

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	httpactionconfig "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/httpaction/config"
	thelpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// HTTP Action test cases for successful CRUD operations
type httpActionSuccessTest struct {
	name     string
	testCase string
	method   string
	body     string
	url      string
}

var httpActionSuccessTests = []httpActionSuccessTest{
	{
		name:     "POST operation",
		testCase: "crud-success",
		method:   "POST",
		body:     `{"name": "Test Resource", "type": "test"}`,
		url:      "http://<host>:<port>/api/resources",
	},
	{
		name:     "GET operation",
		testCase: "crud-success",
		method:   "GET",
		body:     ``,
		url:      "http://<host>:<port>/api/resources/test-resource-123",
	},
	{
		name:     "PUT operation",
		testCase: "crud-success",
		method:   "PUT",
		body:     `{"name": "Updated Test Resource", "type": "test"}`,
		url:      "http://<host>:<port>/api/resources/test-resource-123",
	},
	{
		name:     "DELETE operation",
		testCase: "crud-success",
		method:   "DELETE",
		body:     ``,
		url:      "http://<host>:<port>/api/resources/test-resource-123",
	},
}

// ExecuteHTTPActionCRUDSuccessTest executes HTTP Action CRUD operations success test
func ExecuteHTTPActionCRUDSuccessTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	// Use the pre-configured fake server port from the environment
	// This port is already whitelisted in the gateway configuration
	fakeServerPort := testEnv.Config.FakeHTTP.Port

	// Start fake HTTP server with CRUD endpoints
	fakeServer, err := thelpers.StartCRUDTestServer(t, fakeServerPort, false)
	require.NoError(t, err, "failed to start fake HTTP server")

	defer func() {
		if fakeServer != nil {
			testLogger.Info().Msgf("Cleaning up fake server on port %d", fakeServerPort)
		}
	}()

	for _, testCase := range httpActionSuccessTests {
		// Set dynamic URL with configured port and host.docker.internal
		testCase.url = strings.ReplaceAll(testCase.url, "<port>", strconv.Itoa(fakeServerPort))
		testCase.url = strings.ReplaceAll(testCase.url, "<host>", "host.docker.internal") // Use

		testName := "[v2] HTTP Action " + testCase.name
		t.Run(testName, func(t *testing.T) {
			HTTPActionSuccessTest(t, testEnv, testCase)
		})
	}
}

// HTTPActionSuccessTest executes a single HTTP Action success test case
func HTTPActionSuccessTest(t *testing.T, testEnv *ttypes.TestEnvironment, httpActionTest httpActionSuccessTest) {
	testLogger := framework.L
	const workflowFileLocation = "./httpaction/main.go"

	testLogger.Info().Msg("Creating HTTP Action success test workflow configuration...")

	workflowConfig := httpactionconfig.Config{
		URL:      httpActionTest.url,
		TestCase: httpActionTest.testCase,
		Method:   httpActionTest.method,
		Body:     httpActionTest.body,
	}

	testID := uuid.New().String()[0:8]
	workflowName := "http-action-success-workflow-" + httpActionTest.testCase + "-" + testID
	thelpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	// Start Beholder listener to capture workflow execution messages
	listenerCtx, messageChan, kafkaErrChan := thelpers.StartBeholder(t, testLogger, testEnv)

	// Wait for workflow execution to complete and verify success
	testLogger.Info().Msg("Waiting for HTTP Action CRUD operations to complete...")
	timeout := 60 * time.Second

	// Expect exact success message for this test case
	expectedMessage := "HTTP Action CRUD success test completed: " + httpActionTest.testCase
	err := thelpers.AssertBeholderMessage(listenerCtx, t, expectedMessage, testLogger, messageChan, kafkaErrChan, timeout)
	require.NoError(t, err, "HTTP Action CRUD success test failed")

	testLogger.Info().Msg("HTTP Action CRUD success test completed")
}
