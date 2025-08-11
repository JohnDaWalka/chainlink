package environment

import (
	"errors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/s3provider"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

type Config struct {
	Blockchains            []*cre.WrappedBlockchainInput   `toml:"blockchains" validate:"required"`
	NodeSets               []*cre.CapabilitiesAwareNodeSet `toml:"nodesets" validate:"required"`
	JD                     *jd.Input                       `toml:"jd" validate:"required"`
	Infra                  *infra.Input                    `toml:"infra" validate:"required"`
	Fake                   *fake.Input                     `toml:"fake" validate:"required"`
	CapabilitiesConfig     CapabilitiesConfig              `toml:"capabilities_configs"`
	S3ProviderInput        *s3provider.Input               `toml:"s3provider"`
	AdditionalCapabilities map[string]cre.CapabilityConfig `toml:"additional_capabilities"` // capability flag -> capability config
}

type CapabilitiesConfig struct {
	EVM map[string]map[string]any `toml:"evm"`
}

func (c Config) Validate() error {
	if c.JD.CSAEncryptionKey == "" {
		return errors.New("jd.csa_encryption_key must be provided")
	}
	return nil
}
