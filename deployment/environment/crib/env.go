package crib

const (
	AddressBookFileName       = "address-book.json"
	NodesDetailsFileName      = "nodes-details.json"
	ChainsConfigsFileName     = "chains-details.json"
	RMNNodeIdentitiesFileName = "rmn-node-identities.json"
)

type CRIBEnv struct {
	cribEnvStateDirPath string
}

func NewDevspaceEnvFromStateDir(envStateDir string) CRIBEnv {
	return CRIBEnv{
		cribEnvStateDirPath: envStateDir,
	}
}

func (c CRIBEnv) GetConfig(key string) (DeployOutput, error) {
	// todo: read new fields
	// JDOutput:          nil,
	// BlockchainOutputs: nil,
	// NodesetOutput:     nil,

	reader := NewOutputReader(c.cribEnvStateDirPath)
	nodesDetails := reader.ReadNodesDetails()
	chainConfigs := reader.ReadChainConfigs()
	for i, chain := range chainConfigs {
		err := chain.SetDeployerKey(&key)
		if err != nil {
			return DeployOutput{}, err
		}
		chainConfigs[i] = chain
	}

	return DeployOutput{
		NodeIDs:           nodesDetails.NodeIDs,
		Chains:            chainConfigs,
		AddressBook:       reader.ReadAddressBook(),
		JDOutput:          nil,
		BlockchainOutputs: nil,
		NodesetOutput:     nil,
	}, nil
}

type RPC struct {
	External *string
	Internal *string
}

type ChainConfig struct {
	ChainID   uint64 // chain id as per EIP-155, mainly applicable for EVM chains
	ChainName string // name of the chain populated from chainselector repo
	ChainType string // should denote the chain family. Acceptable values are EVM, COSMOS, SOLANA, STARKNET, APTOS etc
	WSRPCs    []RPC  // websocket rpcs to connect to the chain
	HTTPRPCs  []RPC  // http rpcs to connect to the chain
}

type BootstrapNode struct {
	P2PID        string
	InternalHost string
	Port         string
}

type NodesDetails struct {
	NodeIDs       []string
	BootstrapNode BootstrapNode
}
