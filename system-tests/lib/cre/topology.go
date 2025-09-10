package cre

import (
	"fmt"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"

	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

const (
	OCRPeeringPort          = 5001
	CapabilitiesPeeringPort = 6690
	GatewayIncomingPort     = 5002
	GatewayOutgoingPort     = 5003
)

var (
	NodeTypeKey            = "type"
	HostLabelKey           = "host"
	IndexKey               = "node_index"
	ExtraRolesKey          = "extra_roles"
	NodeIDKey              = "node_id"
	NodeOCR2KeyBundleIDKey = "ocr2_key_bundle_id"
	NodeP2PIDKey           = "p2p_id"
	DONIDKey               = "don_id"
	EnvironmentKey         = "environment"
	ProductKey             = "product"
	DONNameKey             = "don_name"
)

type Topology struct {
	WorkflowDONID           uint64                  `toml:"workflow_don_id" json:"workflow_don_id"`
	HomeChainSelector       uint64                  `toml:"home_chain_selector" json:"home_chain_selector"`
	DonsMetadata            []*DonMetadata          `toml:"dons_metadata" json:"dons_metadata"`
	CapabilitiesPeeringData CapabilitiesPeeringData `toml:"capabilities_peering_data" json:"capabilities_peering_data"`
	OCRPeeringData          OCRPeeringData          `toml:"ocr_peering_data" json:"ocr_peering_data"`
	GatewayConnectorOutput  *GatewayConnectorOutput `toml:"gateway_connector_output" json:"gateway_connector_output"`
}

func NewTopology(nodeSetInput []*CapabilitiesAwareNodeSet, infraInput infra.Input, homeChainSelector uint64) (*Topology, error) {
	// TODO this setup is awkward, consider an withInfra opt to constructor
	dm := make([]*DonMetadata, len(nodeSetInput))
	for i := range nodeSetInput {
		// TODO take more care about the ID assignment, it should match what the capabilities registry will assign
		// currently we optimistically set the id to the that which the capabilities registry will assign it
		d := NewDonMetadata(nodeSetInput[i], libc.MustSafeUint64FromInt(i+1))
		d.labelNodes(infraInput)
		dm[i] = d
	}

	donsMetadata, err := NewDonsMetadata(dm, infraInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create DONs metadata: %w", err)
	}

	topology := &Topology{}
	if donsMetadata.GatewayRequired() {
		topology.GatewayConnectorOutput = NewGatewayConnectorOutput()
	}

	for _, d := range donsMetadata.List() {
		if d.IsGateway() {
			gc, err := d.GatewayConfig(infraInput)
			if err != nil {
				return nil, fmt.Errorf("failed to get gateway config for DON %s: %w", d.Name, err)
			}
			topology.GatewayConnectorOutput.Configurations = append(topology.GatewayConnectorOutput.Configurations, gc)
		}
	}
	wfDon, err := donsMetadata.GetWorkflowDON()
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow DON: %w", err)
	}
	bt, err := wfDon.GetBootstrapNode()
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow DON bootstrap node: %w", err)
	}
	capPeeringCfg, ocrPeeringCfg, err := peeringCfgs(bt)
	if err != nil {
		return nil, fmt.Errorf("failed to get peering data: %w", err)
	}
	topology.DonsMetadata = dm
	topology.WorkflowDONID = wfDon.ID
	topology.HomeChainSelector = homeChainSelector
	topology.CapabilitiesPeeringData = capPeeringCfg
	topology.OCRPeeringData = ocrPeeringCfg
	return topology, nil
}

func peeringCfgs(bt *NodeMetadata) (CapabilitiesPeeringData, OCRPeeringData, error) {

	return CapabilitiesPeeringData{
			GlobalBootstraperPeerID: bt.P2P,
			GlobalBootstraperHost:   bt.Host,
			Port:                    CapabilitiesPeeringPort,
		}, OCRPeeringData{
			OCRBootstraperPeerID: bt.P2P,
			OCRBootstraperHost:   bt.Host,
			Port:                 OCRPeeringPort,
		}, nil
}
