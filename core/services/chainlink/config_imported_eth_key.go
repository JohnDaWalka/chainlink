package chainlink

import (
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

type importedEthKeyConfig struct {
	s toml.EthKey
}

func (t *importedEthKeyConfig) JSON() string {
	if t.s.JSON == nil {
		return ""
	}
	return string(*t.s.JSON)
}

func (t *importedEthKeyConfig) ChainDetails() chain_selectors.ChainDetails {
	if t.s.ChainDetails == nil {
		return chain_selectors.ChainDetails{}
	}
	return *t.s.ChainDetails
}

func (t *importedEthKeyConfig) Password() string {
	if t.s.Password == nil {
		return ""
	}
	return string(*t.s.Password)
}

type importedEthKeyConfigs struct {
	s toml.EthKeysWrapper
}

func (t *importedEthKeyConfigs) List() []config.EthKeyConfig {
	res := make([]config.EthKeyConfig, len(t.s.EthKeys))

	if len(t.s.EthKeys) == 0 {
		return res
	}

	for i, v := range t.s.EthKeys {
		res[i] = &importedEthKeyConfig{s: *v}
	}
	return res
}
