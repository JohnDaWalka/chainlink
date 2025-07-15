package config

import "time"

type EAMetricsReporter interface {
	Enabled() bool
	MetricsPath() string
	PollingInterval() time.Duration
}
