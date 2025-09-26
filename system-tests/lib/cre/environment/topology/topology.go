package topology

import (
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
)

func NewDonTopology(registryChainSelector uint64, topology *cre.Topology, dons []*don.DON) *DonTopology {
	return &DonTopology{
		WorkflowDonID:          topology.WorkflowDONID,
		HomeChainSelector:      registryChainSelector,
		Dons:                   dons,
		GatewayConnectorOutput: topology.GatewayConnectorOutput,
	}
}

// TODO refactor it to only contain []DON, once we have our own DON struct
// and maybe the GatewayConnectorOutput
type DonTopology struct {
	WorkflowDonID          uint64                      `toml:"workflow_don_id" json:"workflow_don_id"`
	HomeChainSelector      uint64                      `toml:"home_chain_selector" json:"home_chain_selector"`
	Dons                   []*don.DON                  `toml:"dons" json:"dons"`
	GatewayConnectorOutput *cre.GatewayConnectorOutput `toml:"gateway_connector_output" json:"gateway_connector_output"`
}

// BootstrapNode returns the metadata for the node that should be used as the bootstrap node for P2P peering
// Currently only one bootstrap is supported.
func (t *DonTopology) BootstrapNode() (*cre.NodeMetadata, error) {
	// for _, don := range t.Dons {
	// 	if don.ContainsBootstrapNode() {
	// 		return don.BootstrapNode()
	// 	}
	// }
	// return nil, errors.New("no don contains a bootstrap node")

	return nil, nil // fix me
}

func (t *DonTopology) ToDonMetadata() []*cre.DonMetadata {
	metadata := []*cre.DonMetadata{}
	// TODO fix me
	// metadata = append(metadata, t.Dons.DonMetadata...)

	return metadata
}

type Environment struct {
	CldfEnvironment *cldf.Environment
	DonTopology     *DonTopology
}
