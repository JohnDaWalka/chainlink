package chainlink

import (
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

var _ config.EAMetricsReporter = (*eaMetricsReporterConfig)(nil)

type eaMetricsReporterConfig struct {
	c toml.EAMetricsReporter
}

func (e *eaMetricsReporterConfig) Enabled() bool {
	return *e.c.Enabled
}

func (e *eaMetricsReporterConfig) MetricsPath() string {
	return *e.c.MetricsPath
}

func (e *eaMetricsReporterConfig) PollingInterval() time.Duration {
	return e.c.PollingInterval.Duration()
}
