package crib

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment/environment/types"
)

const (
	AddressBookFileName       = "address-book.json"
	NodesDetailsFileName      = "nodes-details.json"
	ChainsConfigsFileName     = "chains-details.json"
	RMNNodeIdentitiesFileName = "rmn-node-identities.json"
	JDOutputFileName          = "jd-output.json"
	BlockChainsOutputFileName = "blockchains-output.json"
	NodeSetOutputFileName     = "nodeset-output.json"
)

type CRIBEnv struct {
	lggr                logger.Logger
	cribEnvStateDirPath string
}

func NewCRIBEnvFromStateDir(lggr logger.Logger, envStateDir string) CRIBEnv {
	return CRIBEnv{
		lggr:                lggr,
		cribEnvStateDirPath: envStateDir,
	}
}

func (c CRIBEnv) GetConfig(key string) (CCIPInfraAndOnChainDeployOutput, error) {
	reader := NewOutputReader(c.cribEnvStateDirPath)
	nodesDetails, err := reader.ReadNodesDetails()
	if err != nil {
		c.lggr.Warn("No nodes details found, not necessary for testing.. continuing...", err)
	}
	chainConfigs, err := reader.ReadChainConfigs()
	if err != nil {
		return CCIPInfraAndOnChainDeployOutput{}, errors.Wrapf(err, "failed to read chain configs from %s", c.cribEnvStateDirPath)
	}
	for i, chain := range chainConfigs {
		err := chain.SetDeployerKey(&key)
		if err != nil {
			return CCIPInfraAndOnChainDeployOutput{}, err
		}
		chainConfigs[i] = chain
	}

	addressBook, err := reader.ReadAddressBook()
	if err != nil {
		return CCIPInfraAndOnChainDeployOutput{}, errors.Wrapf(err, "failed to read address book")
	}

	jdOutput, err := reader.ReadJDOutput()
	if err != nil {
		return CCIPInfraAndOnChainDeployOutput{}, errors.Wrap(err, "error reading jd output")
	}

	blockchainOutputs, err := reader.ReadBlockchainOutputs()
	if err != nil {
		return CCIPInfraAndOnChainDeployOutput{}, errors.Wrap(err, "error reading blockchain outputs")
	}

	nodeSetOutput, err := reader.ReadNodeSetOutput()
	if err != nil {
		return CCIPInfraAndOnChainDeployOutput{}, errors.Wrap(err, "error reading node set output")
	}

	return CCIPInfraAndOnChainDeployOutput{
		NodeIDs:           nodesDetails.NodeIDs,
		Chains:            chainConfigs,
		AddressBook:       addressBook,
		JDOutput:          jdOutput,
		BlockchainOutputs: types.ChainIDToBlockchainOutputsFromArray(blockchainOutputs),
		NodesetOutput:     nodeSetOutput,
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
