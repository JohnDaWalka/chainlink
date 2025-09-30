package dctl

import (
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// DctlClient is a thin wrapper around the dctl binary tool
type DctlClient struct {
	binaryPath  string
	stdout      io.Writer
	stderr      io.Writer
	kubeContext string
	logLevel    string
	verbose     bool
	quiet       bool
}

// DctlOption configures the DctlClient
type DctlOption func(*DctlClient)

// NewDctlClient creates a new dctl client with optional configuration
func NewDctlClient(opts ...DctlOption) *DctlClient {
	client := &DctlClient{
		binaryPath:  "dctl",
		stdout:      os.Stdout,
		stderr:      os.Stderr,
		kubeContext: "kind-griddle-dev",
		logLevel:    "info",
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

// WithVerboseOutput enables verbose logging
func WithVerboseOutput(verbose bool) DctlOption {
	return func(c *DctlClient) {
		c.verbose = verbose
	}
}

// WithQuietOutput suppresses all output except errors
func WithQuietOutput(quiet bool) DctlOption {
	return func(c *DctlClient) {
		c.quiet = quiet
	}
}

// WithLogLevel sets the logging level
func WithLogLevel(level string) DctlOption {
	return func(c *DctlClient) {
		c.logLevel = level
	}
}

// WithKubeContext sets the kubernetes context
func WithKubeContext(context string) DctlOption {
	return func(c *DctlClient) {
		c.kubeContext = context
	}
}

// addGlobalFlags adds global flags to command args
func (c *DctlClient) addGlobalFlags(args []string) []string {
	if c.verbose {
		args = append(args, "--verbose")
	}
	if c.quiet {
		args = append(args, "--quiet")
	}
	if c.logLevel != "info" {
		args = append(args, "--log-level", c.logLevel)
	}
	return args
}

// addDeployFlags adds deploy-specific flags to command args
func (c *DctlClient) addDeployFlags(args []string, opts *DeployOptions) []string {
	if opts != nil {
		if opts.Bootstrap {
			args = append(args, "--bootstrap")
		}
		if opts.SkipCache {
			args = append(args, "--skip-cache")
		}
		if opts.SkipDeps {
			args = append(args, "--skip-deps")
		}
		if opts.SkipInit {
			args = append(args, "--skip-init")
		}
	}
	if c.kubeContext != "kind-griddle-dev" {
		args = append(args, "--kube-context", c.kubeContext)
	}
	return args
}

// DeployOptions holds deploy command options
type DeployOptions struct {
	Bootstrap bool
	SkipCache bool
	SkipDeps  bool
	SkipInit  bool
}

// Network Commands

// NetworkConnect connects to the cluster using Telepresence
func (c *DctlClient) NetworkConnect() error {
	args := []string{"network", "connect"}
	args = c.addGlobalFlags(args)
	return c.execute(args...)
}

// NetworkDisconnect disconnects from the cluster
func (c *DctlClient) NetworkDisconnect() error {
	args := []string{"network", "disconnect"}
	args = c.addGlobalFlags(args)
	return c.execute(args...)
}

// NetworkStatus shows Telepresence connection status
func (c *DctlClient) NetworkStatus() error {
	args := []string{"network", "status"}
	args = c.addGlobalFlags(args)
	return c.execute(args...)
}

// Init Commands

// Init executes "dctl init" command to initialize the environment
func (c *DctlClient) Init(contextName string) error {
	args := []string{"init"}
	if contextName != "" {
		args = append(args, "--name", contextName)
	}
	args = c.addGlobalFlags(args)
	return c.execute(args...)
}

// Deploy Commands

// DeployApply executes "dctl deploy apply -f <configFile>"
func (c *DctlClient) DeployApply(configFile string, templateName string, namespace string, opts *DeployOptions) error {
	args := []string{"deploy", "apply", "-f", configFile, "-t", templateName, "-n", namespace}
	args = c.addDeployFlags(args, opts)
	args = c.addGlobalFlags(args)
	return c.execute(args...)
}

// Delete executes "dctl deploy delete -f <configFile>"
func (c *DctlClient) Delete(configFile string, opts *DeployOptions) error {
	args := []string{"deploy", "delete", "-f", configFile}
	args = c.addDeployFlags(args, opts)
	args = c.addGlobalFlags(args)
	return c.execute(args...)
}

// DeployDescribe describes a deployed component
func (c *DctlClient) DeployDescribe(componentName string, opts *DeployOptions) error {
	args := []string{"deploy", "describe", componentName}
	args = c.addDeployFlags(args, opts)
	args = c.addGlobalFlags(args)
	return c.execute(args...)
}

// DeployList lists deployed components
func (c *DctlClient) DeployList(opts *DeployOptions) error {
	args := []string{"deploy", "list"}
	args = c.addDeployFlags(args, opts)
	args = c.addGlobalFlags(args)
	return c.execute(args...)
}

// DeployStatus checks the status of deployed components
func (c *DctlClient) DeployStatus(opts *DeployOptions) error {
	args := []string{"deploy", "status"}
	args = c.addDeployFlags(args, opts)
	args = c.addGlobalFlags(args)
	return c.execute(args...)
}

// DeployWatch watches config files and automatically applies changes
func (c *DctlClient) DeployWatch(configFile string, opts *DeployOptions) error {
	args := []string{"deploy", "watch", "-f", configFile}
	args = c.addDeployFlags(args, opts)
	args = c.addGlobalFlags(args)
	return c.execute(args...)
}

// execute runs a dctl command with the provided arguments
func (c *DctlClient) execute(args ...string) error {
	cmd := exec.Command(c.binaryPath, args...)
	cmd.Stdout = c.stdout
	cmd.Stderr = c.stderr

	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "failed to execute dctl %v", args)
	}
	return nil
}

// Execute runs a custom dctl command with the provided arguments
func (c *DctlClient) Execute(args ...string) error {
	return c.execute(args...)
}
