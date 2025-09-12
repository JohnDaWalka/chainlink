package lane_v2

import (
	"testing"

	"github.com/test-go/testify/require"
)

func TestMultipleLanes(t *testing.T) {
	cfgs := multipleFamiliesLanes()
	require.NotEmpty(t, cfgs)
}

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
