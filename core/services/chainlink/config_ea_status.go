package chainlink

import (
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

var _ config.EAStatusReporter = (*eaStatusReporterConfig)(nil)

type eaStatusReporterConfig struct {
	c toml.EAStatusReporter
}

func (e *eaStatusReporterConfig) Enabled() bool {
	if e.c.Enabled == nil {
		return false
	}
	return *e.c.Enabled
}

func (e *eaStatusReporterConfig) StatusPath() string {
	if e.c.StatusPath == nil {
		return "/status"
	}
	return *e.c.StatusPath
}

func (e *eaStatusReporterConfig) PollingInterval() time.Duration {
	if e.c.PollingInterval == nil {
		return 5 * time.Minute
	}
	return e.c.PollingInterval.Duration()
}

func (e *eaStatusReporterConfig) IgnoreInvalidBridges() bool {
	if e.c.IgnoreInvalidBridges == nil {
		return true
	}
	return *e.c.IgnoreInvalidBridges
}

func (e *eaStatusReporterConfig) IgnoreJoblessBridges() bool {
	if e.c.IgnoreJoblessBridges == nil {
		return false
	}
	return *e.c.IgnoreJoblessBridges
}
