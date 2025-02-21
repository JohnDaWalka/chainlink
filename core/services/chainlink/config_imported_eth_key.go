package chainlink

import (
	chain_selectors "github.com/smartcontractkit/chain-selectors"
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
	if t.s.Selector == nil {
		return chain_selectors.ChainDetails{}
	}
	return *t.s.Selector
}

func (t *importedEthKeyConfig) Password() string {
	if t.s.Password == nil {
		return ""
	}
	return string(*t.s.Password)
}
