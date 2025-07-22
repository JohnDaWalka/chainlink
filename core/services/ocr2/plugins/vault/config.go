package vault

import (
	"errors"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
)

type DKGConfig struct {
	N     int
	T     int
	Index int
}

type Config struct {
	RequestExpiryDuration commonconfig.Duration `json:"requestExpiryDuration"`
	DKG                   *DKGConfig
}

func (c *Config) Validate() error {
	if c.RequestExpiryDuration.Duration() <= 0 {
		return errors.New("request expiry duration cannot be 0")
	}
	return nil
}
