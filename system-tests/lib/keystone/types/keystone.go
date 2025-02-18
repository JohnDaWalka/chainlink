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
	ChainSelector uint64                  `toml:"-"`
	CldEnv        *deployment.Environment `toml:"-"`
	Out           *KeystoneContractOutput `toml:"out"`
}

type KeystoneContractOutput struct {
	UseCache                    bool           `toml:"use_cache"`
	CapabilitiesRegistryAddress common.Address `toml:"capabilities_registry_address"`
	ForwarderAddress            common.Address `toml:"forwarder_address"`
	OCR3CapabilityAddress       common.Address `toml:"ocr3_capability_address"`
	WorkflowRegistryAddress     common.Address `toml:"workflow_registry_address"`
}

type WorkflowRegistryInput struct {
	ChainSelector  uint64                  `toml:"-"`
	CldEnv         *deployment.Environment `toml:"-"`
	AllowedDonIDs  []uint32                `toml:"-"`
	WorkflowOwners []common.Address        `toml:"-"`
	Out            *WorkflowRegistryOutput `toml:"out"`
}

type WorkflowRegistryOutput struct {
	UseCache       bool             `toml:"use_cache"`
	ChainSelector  uint64           `toml:"chain_selector"`
	AllowedDonIDs  []uint32         `toml:"allowed_don_ids"`
	WorkflowOwners []common.Address `toml:"workflow_owners"`
}

type DeployFeedConsumerInput struct {
	ChainSelector uint64                    `toml:"-"`
	CldEnv        *deployment.Environment   `toml:"-"`
	Out           *DeployFeedConsumerOutput `toml:"out"`
}

type DeployFeedConsumerOutput struct {
	UseCache            bool           `toml:"use_cache"`
	FeedConsumerAddress common.Address `toml:"feed_consumer_address"`
}

type ConfigureFeedConsumerInput struct {
	SethClient            *seth.Client                 `toml:"-"`
	FeedConsumerAddress   common.Address               `toml:"-"`
	AllowedSenders        []common.Address             `toml:"-"`
	AllowedWorkflowOwners []common.Address             `toml:"-"`
	AllowedWorkflowNames  []string                     `toml:"-"`
	Out                   *ConfigureFeedConsumerOutput `toml:"out"`
}

type ConfigureFeedConsumerOutput struct {
	UseCache              bool             `toml:"use_cache"`
	FeedConsumerAddress   common.Address   `toml:"feed_consumer_address"`
	AllowedSenders        []common.Address `toml:"allowed_senders"`
	AllowedWorkflowOwners []common.Address `toml:"allowed_workflow_owners"`
	AllowedWorkflowNames  []string         `toml:"allowed_workflow_names"`
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
