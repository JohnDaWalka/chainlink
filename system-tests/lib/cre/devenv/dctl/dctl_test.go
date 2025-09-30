package dctl

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDctlClient_NewClient(t *testing.T) {
	client := NewDctlClient()

	assert.Equal(t, "dctl", client.binaryPath)
	assert.NotNil(t, client.stdout)
	assert.NotNil(t, client.stderr)
	assert.Equal(t, "kind-griddle-dev", client.kubeContext)
	assert.Equal(t, "info", client.logLevel)
	assert.False(t, client.verbose)
	assert.False(t, client.quiet)
}

func TestDctlClient_WithOptions(t *testing.T) {
	var stdout, stderr bytes.Buffer

	client := NewDctlClient(
		WithBinaryPath("/custom/path/dctl"),
		WithStdout(&stdout),
		WithStderr(&stderr),
		WithKubeContext("custom-context"),
		WithLogLevel("debug"),
		WithVerboseOutput(true),
		WithQuietOutput(true),
	)

	assert.Equal(t, "/custom/path/dctl", client.binaryPath)
	assert.Equal(t, &stdout, client.stdout)
	assert.Equal(t, &stderr, client.stderr)
	assert.Equal(t, "custom-context", client.kubeContext)
	assert.Equal(t, "debug", client.logLevel)
	assert.True(t, client.verbose)
	assert.True(t, client.quiet)
}

func TestDctlClient_NetworkConnect(t *testing.T) {
	var stdout, stderr bytes.Buffer

	client := NewDctlClient(
		WithBinaryPath("echo"),
		WithStdout(&stdout),
		WithStderr(&stderr),
	)

	err := client.NetworkConnect()
	assert.NoError(t, err)

	assert.Contains(t, stdout.String(), "network connect")
}

func TestDctlClient_NetworkDisconnect(t *testing.T) {
	var stdout, stderr bytes.Buffer

	client := NewDctlClient(
		WithBinaryPath("echo"),
		WithStdout(&stdout),
		WithStderr(&stderr),
	)

	err := client.NetworkDisconnect()
	assert.NoError(t, err)

	assert.Contains(t, stdout.String(), "network disconnect")
}

func TestDctlClient_NetworkStatus(t *testing.T) {
	var stdout, stderr bytes.Buffer

	client := NewDctlClient(
		WithBinaryPath("echo"),
		WithStdout(&stdout),
		WithStderr(&stderr),
	)

	err := client.NetworkStatus()
	assert.NoError(t, err)

	assert.Contains(t, stdout.String(), "network status")
}

func TestDctlClient_Init(t *testing.T) {
	var stdout, stderr bytes.Buffer

	client := NewDctlClient(
		WithBinaryPath("echo"),
		WithStdout(&stdout),
		WithStderr(&stderr),
	)

	err := client.Init("custom-context")
	assert.NoError(t, err)

	output := stdout.String()
	assert.Contains(t, output, "init")
	assert.Contains(t, output, "--name custom-context")
}

func TestDctlClient_InitWithEmptyContext(t *testing.T) {
	var stdout, stderr bytes.Buffer

	client := NewDctlClient(
		WithBinaryPath("echo"),
		WithStdout(&stdout),
		WithStderr(&stderr),
	)

	err := client.Init("")
	assert.NoError(t, err)

	output := stdout.String()
	assert.Contains(t, output, "init")
	assert.NotContains(t, output, "--name")
}

func TestDctlClient_Apply(t *testing.T) {
	var stdout, stderr bytes.Buffer

	client := NewDctlClient(
		WithBinaryPath("echo"),
		WithStdout(&stdout),
		WithStderr(&stderr),
	)

	err := client.DeployApply("test-config.yaml", "template-name", "test-namespace", nil)
	assert.NoError(t, err)

	output := stdout.String()
	assert.Contains(t, output, "deploy apply -f test-config.yaml")
	assert.Contains(t, output, "-t template-name")
	assert.Contains(t, output, "-n test-namespace")
}

func TestDctlClient_ApplyWithOptions(t *testing.T) {
	var stdout, stderr bytes.Buffer

	client := NewDctlClient(
		WithBinaryPath("echo"),
		WithStdout(&stdout),
		WithStderr(&stderr),
		WithKubeContext("custom-context"),
	)

	opts := &DeployOptions{
		Bootstrap: true,
		SkipCache: true,
		SkipDeps:  true,
		SkipInit:  true,
	}

	err := client.DeployApply("test-config.yaml", "template-name", "test-namespace", opts)
	assert.NoError(t, err)

	output := stdout.String()
	assert.Contains(t, output, "deploy apply")
	assert.Contains(t, output, "--bootstrap")
	assert.Contains(t, output, "--skip-cache")
	assert.Contains(t, output, "--skip-deps")
	assert.Contains(t, output, "--skip-init")
	assert.Contains(t, output, "--kube-context custom-context")
}

func TestDctlClient_Delete(t *testing.T) {
	var stdout, stderr bytes.Buffer

	client := NewDctlClient(
		WithBinaryPath("echo"),
		WithStdout(&stdout),
		WithStderr(&stderr),
	)

	err := client.Delete("test-config.yaml", nil)
	assert.NoError(t, err)

	assert.Contains(t, stdout.String(), "deploy delete -f test-config.yaml")
}

func TestDctlClient_DeployDescribe(t *testing.T) {
	var stdout, stderr bytes.Buffer

	client := NewDctlClient(
		WithBinaryPath("echo"),
		WithStdout(&stdout),
		WithStderr(&stderr),
	)

	err := client.DeployDescribe("component-name", nil)
	assert.NoError(t, err)

	assert.Contains(t, stdout.String(), "deploy describe component-name")
}

func TestDctlClient_DeployList(t *testing.T) {
	var stdout, stderr bytes.Buffer

	client := NewDctlClient(
		WithBinaryPath("echo"),
		WithStdout(&stdout),
		WithStderr(&stderr),
	)

	err := client.DeployList(nil)
	assert.NoError(t, err)

	assert.Contains(t, stdout.String(), "deploy list")
}

func TestDctlClient_DeployStatus(t *testing.T) {
	var stdout, stderr bytes.Buffer

	client := NewDctlClient(
		WithBinaryPath("echo"),
		WithStdout(&stdout),
		WithStderr(&stderr),
	)

	err := client.DeployStatus(nil)
	assert.NoError(t, err)

	assert.Contains(t, stdout.String(), "deploy status")
}

func TestDctlClient_DeployWatch(t *testing.T) {
	var stdout, stderr bytes.Buffer

	client := NewDctlClient(
		WithBinaryPath("echo"),
		WithStdout(&stdout),
		WithStderr(&stderr),
	)

	err := client.DeployWatch("test-config.yaml", nil)
	assert.NoError(t, err)

	assert.Contains(t, stdout.String(), "deploy watch -f test-config.yaml")
}

func TestDctlClient_GlobalFlags(t *testing.T) {
	var stdout, stderr bytes.Buffer

	client := NewDctlClient(
		WithBinaryPath("echo"),
		WithStdout(&stdout),
		WithStderr(&stderr),
		WithVerboseOutput(true),
		WithLogLevel("debug"),
	)

	err := client.NetworkStatus()
	assert.NoError(t, err)

	output := stdout.String()
	assert.Contains(t, output, "--verbose")
	assert.Contains(t, output, "--log-level debug")
}

func TestDctlClient_QuietFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer

	client := NewDctlClient(
		WithBinaryPath("echo"),
		WithStdout(&stdout),
		WithStderr(&stderr),
		WithQuietOutput(true),
	)

	err := client.NetworkStatus()
	assert.NoError(t, err)

	assert.Contains(t, stdout.String(), "--quiet")
}

func TestDctlClient_Execute(t *testing.T) {
	var stdout, stderr bytes.Buffer

	client := NewDctlClient(
		WithBinaryPath("echo"),
		WithStdout(&stdout),
		WithStderr(&stderr),
	)

	err := client.Execute("custom", "command", "--flag")
	assert.NoError(t, err)

	assert.Contains(t, stdout.String(), "custom command --flag")
}

func TestDctlClient_DefaultKubeContext(t *testing.T) {
	var stdout, stderr bytes.Buffer

	client := NewDctlClient(
		WithBinaryPath("echo"),
		WithStdout(&stdout),
		WithStderr(&stderr),
	)

	err := client.DeployApply("test.yaml", "template", "namespace", nil)
	assert.NoError(t, err)

	// Default context should not appear in output since it matches the default
	output := stdout.String()
	assert.NotContains(t, output, "--kube-context")
}
