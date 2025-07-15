package config

import "time"

type EAStatusReporter interface {
	Enabled() bool
	StatusPath() string
	PollingInterval() time.Duration
}
