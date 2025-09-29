package devenv

import (
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// DctlClient is a thin wrapper around the dctl binary tool
type DctlClient struct {
	binaryPath string
	stdout     io.Writer
	stderr     io.Writer
	verbose    bool
}

// DctlOption configures the DctlClient
type DctlOption func(*DctlClient)

// NewDctlClient creates a new dctl client with optional configuration
func NewDctlClient(opts ...DctlOption) *DctlClient {
	client := &DctlClient{
		binaryPath: "dctl", // default assumes dctl is in PATH
		stdout:     os.Stdout,
		stderr:     os.Stderr,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// WithBinaryPath sets a custom path to the dctl binary
func WithBinaryPath(path string) DctlOption {
	return func(c *DctlClient) {
		c.binaryPath = path
	}
}

// WithStdout sets a custom stdout writer
func WithStdout(w io.Writer) DctlOption {
	return func(c *DctlClient) {
		c.stdout = w
	}
}

// WithStderr sets a custom stderr writer
func WithStderr(w io.Writer) DctlOption {
	return func(c *DctlClient) {
		c.stderr = w
	}
}

func WithVerboseOutput(verbose bool) DctlOption {
	return func(c *DctlClient) {
		c.verbose = verbose
	}
}

// Apply executes "dctl deploy apply -f <configFile>"
func (c *DctlClient) Apply(configFile string, templateName string, namespace string) error {
	args := []string{
		"deploy", "apply", "-f", configFile, "-t", templateName, "-n", namespace,
	}
	if c.verbose {
		args = append(args, "--verbose")
	}

	cmd := exec.Command(c.binaryPath, args...)
	cmd.Stdout = c.stdout
	cmd.Stderr = c.stderr

	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "failed to execute dctl deploy apply -f %s", configFile)
	}

	return nil
}

// Delete executes "dctl deploy delete -f <configFile>"
func (c *DctlClient) Delete(configFile string) error {
	cmd := exec.Command(c.binaryPath, "deploy", "delete", "-f", configFile)
	cmd.Stdout = c.stdout
	cmd.Stderr = c.stderr

	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "failed to execute dctl deploy delete -f %s", configFile)
	}

	return nil
}

// Status executes "dctl deploy status"
func (c *DctlClient) Status() error {
	cmd := exec.Command(c.binaryPath, "deploy", "status")
	cmd.Stdout = c.stdout
	cmd.Stderr = c.stderr

	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "failed to execute dctl deploy status")
	}

	return nil
}

// Execute runs a custom dctl command with the provided arguments
func (c *DctlClient) Execute(args ...string) error {
	cmd := exec.Command(c.binaryPath, args...)
	cmd.Stdout = c.stdout
	cmd.Stderr = c.stderr

	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "failed to execute dctl %v", args)
	}

	return nil
}
