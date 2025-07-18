package config

import "time"

const MINIMUM_POLLING_INTERVAL = time.Minute

type EAStatusReporter interface {
	Enabled() bool
	StatusPath() string
	PollingInterval() time.Duration
	IgnoreInvalidBridges() bool
	IgnoreJoblessBridges() bool
}
