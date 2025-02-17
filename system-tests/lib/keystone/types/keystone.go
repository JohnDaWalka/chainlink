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

type KeystoneConfiguration interface {
	JdInput() (*jd.Input, error)
	NodeSetInput() ([]*CapabilitiesAwareNodeSet, error)
	BlockchainInput() (*blockchain.Input, error)
}

type KeystoneEnvironmentConsumerFn = func(keystoneEnv *KeystoneEnvironment) error

type JobAndConfigProducingFn = func(keystoneEnv *KeystoneEnvironment) (DonsToConfigOverrides, DonsToJobSpecs, error)

type KeystoneContractAddresses struct {
	CapabilitiesRegistryAddress common.Address
	ForwarderAddress            common.Address
	OCR3CapabilityAddress       common.Address
	WorkflowRegistryAddress     common.Address
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

func (k *KeystoneEnvironment) MustCLDEnvironment() *deployment.Environment {
	if k.Environment == nil {
		panic("CLD environment must be set")
	}
	return k.Environment
}

func (k *KeystoneEnvironment) MustBlockchain() *blockchain.Output {
	if k.Blockchain == nil {
		panic("blockchain must be set")
	}
	return k.Blockchain
}

func (k *KeystoneEnvironment) MustWrappedNodeOutput() []*WrappedNodeOutput {
	if k.WrappedNodeOutput == nil {
		panic("wrapped node output must be set")
	}
	return k.WrappedNodeOutput
}

func (k *KeystoneEnvironment) MustJD() *jd.Output {
	if k.JD == nil {
		panic("job distributor must be set")
	}
	return k.JD
}

func (k *KeystoneEnvironment) MustSethClient() *seth.Client {
	if k.SethClient == nil {
		panic("seth client must be set")
	}
	return k.SethClient
}

func (k *KeystoneEnvironment) MustDONTopology() []*DONTopology {
	if k.DONTopology == nil {
		panic("DON topology must not be empty")
	}
	return k.DONTopology
}

func (k *KeystoneEnvironment) MustKeystoneContractAddresses() *KeystoneContractAddresses {
	if k.KeystoneContractAddresses == nil {
		panic("keystone contract addresses must be set")
	}
	return k.KeystoneContractAddresses
}

func (k *KeystoneEnvironment) MustGatewayConnectorData() *GatewayConnectorData {
	if k.GatewayConnectorData == nil {
		panic("gateway connector data must be set")
	}
	return k.GatewayConnectorData
}

func (k *KeystoneEnvironment) MustNodeInput() []*CapabilitiesAwareNodeSet {
	if k.NodeInput == nil {
		panic("node input must be set")
	}
	return k.NodeInput
}

func (k *KeystoneEnvironment) MustDons() []*devenv.DON {
	if k.Dons == nil {
		panic("dons must be set")
	}
	return k.Dons
}

func (k *KeystoneEnvironment) MustWorkflowDONID() uint32 {
	if k.WorkflowDONID == 0 {
		panic("workflow DON ID must be set")
	}
	return k.WorkflowDONID
}

func (k *KeystoneEnvironment) MustChainSelector() uint64 {
	if k.ChainSelector == 0 {
		panic("chain selector must be set")
	}
	return k.ChainSelector
}
