package mocks

import (
	"time"
)

// TestEAStatusReporterConfig implements config.EAStatusReporter for testing
type TestEAStatusReporterConfig struct {
	enabled         bool
	statusPath      string
	pollingInterval time.Duration
}

func NewTestEAStatusReporterConfig(enabled bool, statusPath string, pollingInterval time.Duration) *TestEAStatusReporterConfig {
	return &TestEAStatusReporterConfig{
		enabled:         enabled,
		statusPath:      statusPath,
		pollingInterval: pollingInterval,
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
