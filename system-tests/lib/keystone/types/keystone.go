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

type DonJobs = map[JobDescription][]*jobv1.ProposeJobRequest
type DonsToJobSpecs = map[uint32]DonJobs

type NodeIndexToConfigOverrides = map[int]string
type DonsToConfigOverrides = map[uint32]NodeIndexToConfigOverrides

type KeystoneContractsInput struct {
	ChainSelector uint64
	CldEnv        *deployment.Environment
}

type KeystoneContractOutput struct {
	CapabilitiesRegistryAddress common.Address
	ForwarderAddress            common.Address
	OCR3CapabilityAddress       common.Address
	WorkflowRegistryAddress     common.Address
}

type WorkflowRegistryInput struct {
	ChainSelector  uint64
	CldEnv         *deployment.Environment
	AllowedDonIDs  []uint32
	WorkflowOwners []common.Address
}

type DeployFeedConsumerInput struct {
	ChainSelector uint64
	CldEnv        *deployment.Environment
}

type DeployFeedConsumerOutput struct {
	Address common.Address
}

type ConfigureFeedConsumerInput struct {
	SethClient            *seth.Client
	FeedConsumerAddress   common.Address
	AllowedSenders        []common.Address
	AllowedWorkflowOwners []common.Address
	AllowedWorkflowNames  []string
}

type WrappedNodeOutput struct {
	*ns.Output
	NodeSetName  string
	Capabilities []string
}

type ConfigureDonInput struct {
	CldEnv               *deployment.Environment
	BlockchainOutput     *blockchain.Output
	DonTopology          *DonTopology
	JdOutput             *jd.Output
	DonToJobSpecs        DonsToJobSpecs
	DonToConfigOverrides DonsToConfigOverrides
}

type ConfigureDonOutput struct {
	JdOutput *deployment.OffchainClient
}

type DebugInput struct {
	DonTopology      *DonTopology
	BlockchainOutput *blockchain.Output
}

type ConfigureKeystoneInput struct {
	ChainSelector uint64
	DonTopology   *DonTopology
	CldEnv        *deployment.Environment
}

type GatewayConnectorOutput struct {
	Host string // do not set, it will be set dynamically
	Path string
	Port int
}

// DonWithMetadata is a struct that holds the DON references and various metadata
type DonWithMetadata struct {
	DON        *devenv.DON
	NodeInput  *CapabilitiesAwareNodeSet
	NodeOutput *WrappedNodeOutput
	ID         uint32
	Flags      []string
}

type DonTopology struct {
	WorkflowDONID uint32
	MetaDons      []*DonWithMetadata
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
