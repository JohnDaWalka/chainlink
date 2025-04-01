package types

import "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

type NodeType = string

const (
	BootstrapNode NodeType = "bootstrap"
	WorkerNode    NodeType = "worker"
	GatewayNode   NodeType = "gateway"
)

type Label struct {
	Key   string
	Value string
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
