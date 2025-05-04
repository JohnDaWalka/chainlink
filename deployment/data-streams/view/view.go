package view

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/v0_5"
)

type Address = string
type ChainView struct {
	Configurator  map[Address]v0_5.ConfiguratorView  `json:"configurator,omitempty"`
	RewardManager map[Address]v0_5.RewardManagerView `json:"rewardManager,omitempty"`
	Verifier      map[Address]v0_5.VerifierView      `json:"verifier,omitempty"`
	VerifierProxy map[Address]v0_5.VerifierProxyView `json:"verifierProxy,omitempty"`
}

func NewChain() ChainView {
	return ChainView{
		Configurator:  make(map[Address]v0_5.ConfiguratorView),
		RewardManager: make(map[Address]v0_5.RewardManagerView),
		Verifier:      make(map[Address]v0_5.VerifierView),
		VerifierProxy: make(map[Address]v0_5.VerifierProxyView),
	}
}

type DataStreamsView struct {
	Chains map[uint64]ChainView `json:"chains,omitempty"`
}

func (v DataStreamsView) MarshalJSON() ([]byte, error) {
	// Alias to avoid recursive calls.
	type Alias DataStreamsView
	return json.MarshalIndent(&struct{ Alias }{Alias: Alias(v)}, "", " ")
}
