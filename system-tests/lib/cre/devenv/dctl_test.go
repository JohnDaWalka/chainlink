package devenv

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
}

func TestDctlClient_WithOptions(t *testing.T) {
	var stdout, stderr bytes.Buffer

	client := NewDctlClient(
		WithBinaryPath("/custom/path/dctl"),
		WithStdout(&stdout),
		WithStderr(&stderr),
	)

	assert.Equal(t, "/custom/path/dctl", client.binaryPath)
	assert.Equal(t, &stdout, client.stdout)
	assert.Equal(t, &stderr, client.stderr)
}

func TestDctlClient_Apply(t *testing.T) {
	var stdout, stderr bytes.Buffer

	// Using echo as a mock dctl binary for testing
	client := NewDctlClient(
		WithBinaryPath("echo"),
		WithStdout(&stdout),
		WithStderr(&stderr),
	)

	err := client.Apply("test-config.yaml")
	assert.NoError(t, err)

	// Echo will output the arguments we passed to it
	assert.Contains(t, stdout.String(), "deploy apply -f test-config.yaml")
}

func TestDctlClient_Delete(t *testing.T) {
	var stdout, stderr bytes.Buffer

	client := NewDctlClient(
		WithBinaryPath("echo"),
		WithStdout(&stdout),
		WithStderr(&stderr),
	)

	err := client.Delete("test-config.yaml")
	assert.NoError(t, err)

	assert.Contains(t, stdout.String(), "deploy delete -f test-config.yaml")
}

func TestDctlClient_Status(t *testing.T) {
	var stdout, stderr bytes.Buffer

	client := NewDctlClient(
		WithBinaryPath("echo"),
		WithStdout(&stdout),
		WithStderr(&stderr),
	)

	err := client.Status()
	assert.NoError(t, err)

	assert.Contains(t, stdout.String(), "deploy status")
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
