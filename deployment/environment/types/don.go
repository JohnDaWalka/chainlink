package types

import (
	"errors"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

type ChainIDToBlockchainOutputs map[string]*blockchain.Output

type NodeType = string

const (
	BootstrapNode NodeType = "bootstrap"
	WorkerNode    NodeType = "worker"
	GatewayNode   NodeType = "gateway"
)

const (
	NodeTypeKey            = "type"
	HostLabelKey           = "host"
	IndexKey               = "node_index"
	EthAddressKey          = "eth_address"
	ExtraRolesKey          = "extra_roles"
	NodeIDKey              = "node_id"
	NodeOCR2KeyBundleIDKey = "ocr2_key_bundle_id"
	NodeP2PIDKey           = "p2p_id"
)

type Label struct {
	Key   string
	Value string
}

func LabelFromProto(p *ptypes.Label) (*Label, error) {
	if p.Value == nil {
		return nil, errors.New("value not set")
	}
	return &Label{
		Key:   p.Key,
		Value: *p.Value,
	}, nil
}

type NodeMetadata struct {
	Labels []*Label
}

type DonMetadata struct {
	NodesMetadata []*NodeMetadata
	Flags         []string
	ID            uint32
	Name          string
}

type GatewayConnectorDons struct {
	MembersEthAddresses []string
	ID                  uint32
}

type GatewayConnectorOutput struct {
	Dons []GatewayConnectorDons // do not set, it will be set dynamically
	Host string                 // do not set, it will be set dynamically
	Path string
	Port int
}

type Topology struct {
	WorkflowDONID          uint32
	DonsMetadata           []*DonMetadata
	GatewayConnectorOutput *GatewayConnectorOutput
}

type WrappedNodeOutput struct {
	*simple_node_set.Output
	NodeSetName  string
	Capabilities []string
}
