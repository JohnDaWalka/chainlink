package mocks

import (
	"time"
)

// TestEAMetricsReporterConfig implements config.EAMetricsReporter for testing
type TestEAMetricsReporterConfig struct {
	enabled         bool
	metricsPath     string
	pollingInterval time.Duration
}

func NewTestEAMetricsReporterConfig(enabled bool, metricsPath string, pollingInterval time.Duration) *TestEAMetricsReporterConfig {
	return &TestEAMetricsReporterConfig{
		enabled:         enabled,
		metricsPath:     metricsPath,
		pollingInterval: pollingInterval,
	}
}

func (e *TestEAMetricsReporterConfig) Enabled() bool {
	return e.enabled
}

func (e *TestEAMetricsReporterConfig) MetricsPath() string {
	return e.metricsPath
}

func (e *TestEAMetricsReporterConfig) PollingInterval() time.Duration {
	return e.pollingInterval
}
