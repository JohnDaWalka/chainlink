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
	return *e.c.Enabled
}

func (e *eaStatusReporterConfig) StatusPath() string {
	return *e.c.StatusPath
}

func (e *eaStatusReporterConfig) PollingInterval() time.Duration {
	return e.c.PollingInterval.Duration()
}
