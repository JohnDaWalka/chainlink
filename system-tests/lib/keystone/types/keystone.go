package types

import (
	"github.com/ethereum/go-ethereum/common"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
)

type KeystoneContractAddresses struct {
	CapabilitiesRegistryAddress common.Address
	ForwarderAddress            common.Address
	OCR3CapabilityAddress       common.Address
	WorkflowRegistryAddress     common.Address
	FeedsConsumerAddress        common.Address
}

type DonJobs = map[JobDescription][]*jobv1.ProposeJobRequest
type DonsToJobSpecs = map[uint32]DonJobs

type NodeIndexToConfigOverrides = map[int]string
type DonsToConfigOverrides = map[uint32]NodeIndexToConfigOverrides

type KeystoneEnvironment struct {
	*deployment.Environment
	// TODO multiple blockchains support, think if we can just use deployment.Environment
	Blockchain                *blockchain.Output
	SethClient                *seth.Client
	ChainSelector             uint64
	DeployerPrivateKey        string
	KeystoneContractAddresses *KeystoneContractAddresses

	JD *jd.Output

	GatewayConnectorData *GatewayConnectorData

	NodeInput         []*CapabilitiesAwareNodeSet
	WrappedNodeOutput []*WrappedNodeOutput
	DONTopology       []*DONTopology
	Dons              []*devenv.DON
	WorkflowDONID     uint32
	JobSpecs          map[JobDescription]*jobv1.ProposeJobRequest
}

type NodeType = string

const (
	BootstrapNode NodeType = "bootstrap"
	WorkerNode    NodeType = "worker"
)

type JobDescription struct {
	Flag     CapabilityFlag
	NodeType string
}

type ConfigDescription struct {
	Flag     CapabilityFlag
	NodeType string
}

type WrappedNodeOutput struct {
	*ns.Output
	NodeSetName  string
	Capabilities []string
}

// DONTopology is a struct that holds the DON references and various metadata
type DONTopology struct {
	DON        *devenv.DON
	NodeInput  *CapabilitiesAwareNodeSet
	NodeOutput *WrappedNodeOutput
	ID         uint32
	Flags      []string
}

type CapabilitiesAwareNodeSet struct {
	*ns.Input
	Capabilities []string `toml:"capabilities"`
	DONType      string   `toml:"don_type"`
}

type PeeringData struct {
	GlobalBootstraperPeerID string
	GlobalBootstraperHost   string
	Port                    int
}

type OCRPeeringData struct {
	OCRBootstraperPeerID string
	OCRBootstraperHost   string
	Port                 int
}

type GatewayConnectorData struct {
	Host string // do not set, it will be set dynamically
	Path string
	Port int
}
