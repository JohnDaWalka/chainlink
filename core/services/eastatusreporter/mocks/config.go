package mocks

import (
	"time"
)

// TestEAStatusReporterConfig implements config.EAStatusReporter for testing
type TestEAStatusReporterConfig struct {
	enabled              bool
	statusPath           string
	pollingInterval      time.Duration
	ignoreInvalidBridges bool
	ignoreJoblessBridges bool
}

func NewTestEAStatusReporterConfig(enabled bool, statusPath string, pollingInterval time.Duration) *TestEAStatusReporterConfig {
	return &TestEAStatusReporterConfig{
		enabled:              enabled,
		statusPath:           statusPath,
		pollingInterval:      pollingInterval,
		ignoreInvalidBridges: true,  
		ignoreJoblessBridges: false,
	}
}

func NewTestEAStatusReporterConfigWithSkip(enabled bool, statusPath string, pollingInterval time.Duration, ignoreInvalidBridges bool, ignoreJoblessBridges bool) *TestEAStatusReporterConfig {
	return &TestEAStatusReporterConfig{
		enabled:              enabled,
		statusPath:           statusPath,
		pollingInterval:      pollingInterval,
		ignoreInvalidBridges: ignoreInvalidBridges,
		ignoreJoblessBridges: ignoreJoblessBridges,
	}
}

func (e *TestEAStatusReporterConfig) Enabled() bool {
	return e.enabled
}

func (e *TestEAStatusReporterConfig) StatusPath() string {
	return e.statusPath
}

func (e *TestEAStatusReporterConfig) PollingInterval() time.Duration {
	return e.pollingInterval
}

func (e *TestEAStatusReporterConfig) IgnoreInvalidBridges() bool {
	return e.ignoreInvalidBridges
}

func (e *TestEAStatusReporterConfig) IgnoreJoblessBridges() bool {
	return e.ignoreJoblessBridges
}
