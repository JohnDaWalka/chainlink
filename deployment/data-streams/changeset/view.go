package changeset

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink/deployment"
	dsstate "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/state"
	dsView "github.com/smartcontractkit/chainlink/deployment/data-streams/view"
)

var _ deployment.ViewState = ViewDataStreams

func ViewDataStreams(e deployment.Environment) (json.Marshaler, error) {
	state, err := dsstate.LoadOnchainState(e)
	if err != nil {
		return nil, err
	}
	chainView, err := state.View(e.GetContext(), e.AllChainSelectors())
	if err != nil {
		return nil, err
	}
	return dsView.DataStreamsView{
		Chains: chainView,
	}, nil
}
