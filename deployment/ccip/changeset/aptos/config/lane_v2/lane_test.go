package lane_v2

import (
	"errors"
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/test-go/testify/require"
)

func TestMultipleLanes(t *testing.T) {
	cfgs := multipleFamiliesLanes()
	require.NotEmpty(t, cfgs)
	// Call changesets based on selector family
	// This can be used in a multi-family changeset (TBD) or inside a durable pipeline
	for selector, cfg := range cfgs {
		family, _ := chain_selectors.GetSelectorFamily(selector)
		switch family {
		case chain_selectors.FamilyAptos:
			mockAptosUpdateCS(cfg.(*UpdateAptosLanes))
		default:
			panic(errors.New("unsupported chain family"))
		}

	}
}

func mockAptosUpdateCS(_ *UpdateAptosLanes) {}

func multipleFamiliesLanes() map[uint64]UpdateLanesCfg {
	lanes := []LaneConfig{
		{
			Source: ChainDefinition{
				Selector: 743186221051783445, // Aptos
			},
			Dest: ChainDefinition{
				Selector: 3478487238524512106, // EVM 1
			},
		},
		{
			Source: ChainDefinition{
				Selector: 743186221051783445, // Aptos
			},
			Dest: ChainDefinition{
				Selector: 5224473277236331295, // EVM 2
			},
		},
		{
			Source: ChainDefinition{
				Selector: 5224473277236331295, // EVM 2
			},
			Dest: ChainDefinition{
				Selector: 743186221051783445, // Aptos
			},
		},
	}
	updateCfgs := map[uint64]UpdateLanesCfg{}

	for _, lane := range lanes {
		source := lane.Source.Selector
		if _, exists := updateCfgs[source]; !exists {
			updateCfgs[source], _ = NewUpdateLanesCfg(source)
		}
		dest := lane.Dest.Selector
		if _, exists := updateCfgs[dest]; !exists {
			updateCfgs[dest], _ = NewUpdateLanesCfg(dest)
		}
		err := ToUpdateLanesConfig(lane, updateCfgs[source], updateCfgs[dest])
		if err != nil {
			panic(err)
		}
	}

	return updateCfgs
}
