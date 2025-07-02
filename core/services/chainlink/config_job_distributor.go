package chainlink

import (
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

type jobDistributorConfig struct {
	c toml.JobDistributor
}

func (s jobDistributorConfig) DisplayName() string {
	if s.c.DisplayName == nil {
		return ""
	}
	return *s.c.DisplayName
}

func (s jobDistributorConfig) AllowedJobTransfers() []config.JobTransferRule {
	if s.c.AllowedJobTransfers == nil {
		return []config.JobTransferRule{}
	}

	rules := make([]config.JobTransferRule, len(s.c.AllowedJobTransfers))
	for i, rule := range s.c.AllowedJobTransfers {
		rules[i] = config.JobTransferRule{
			From: rule.From,
			To:   rule.To,
		}
	}
	return rules
}
