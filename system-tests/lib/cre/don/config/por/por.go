package por

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	keystoneflags "github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

func GenerateConfigs(input cretypes.GeneratePoRConfigsInput) (cretypes.NodeIndexToOverride, error) {
	// if err := input.Validate(); err != nil {
	// 	return nil, errors.Wrap(err, "input validation failed")
	// }
	configOverrides := make(cretypes.NodeIndexToOverride)

	chainIDInt, err := strconv.Atoi(input.BlockchainOutput.ChainID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert chain ID to int")
	}
	chainIDUint64 := libc.MustSafeUint64(int64(chainIDInt))

	// find bootstrap node
	bootstrapNode, err := node.FindOneWithLabel(input.DonMetadata.NodesMetadata, &ptypes.Label{Key: devenv.NodeLabelKeyType, Value: ptr.Ptr(string(devenv.NodeLabelValueBootstrap))})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find bootstrap node")
	}

	donBootstrapNodePeerID, err := node.ToP2PID(bootstrapNode, node.KeyExtractingTransformFn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get bootstrap node peer ID")
	}

	var donBootstrapNodeHost string
	for _, label := range bootstrapNode.Labels {
		if label.Key == node.HostLabelKey {
			donBootstrapNodeHost = *label.Value
			break
		}
	}

	if donBootstrapNodeHost == "" {
		return nil, errors.New("failed to get bootstrap node host from labels")
	}

	var nodeIndex int
	for _, label := range bootstrapNode.Labels {
		if label.Key == node.IndexKey {
			nodeIndex, err = strconv.Atoi(*label.Value)
			if err != nil {
				return nil, errors.Wrap(err, "failed to convert node index to int")
			}
			break
		}
	}

	// generat configuration for the bootstrap node
	configOverrides[nodeIndex] = config.BootstrapEVM(donBootstrapNodePeerID, chainIDUint64, input.CapabilitiesRegistryAddress, input.BlockchainOutput.Nodes[0].DockerInternalHTTPUrl, input.BlockchainOutput.Nodes[0].DockerInternalWSUrl)

	if keystoneflags.HasFlag(input.Flags, cretypes.WorkflowDON) {
		configOverrides[nodeIndex] += config.BoostrapDon2DonPeering(input.PeeringData)

		if input.GatewayConnectorOutput == nil {
			return nil, errors.New("GatewayConnectorOutput is required for Workflow DON")
		}
		input.GatewayConnectorOutput.Host = donBootstrapNodeHost
	}

	// find worker nodes
	workflowNodeSet, err := node.FindManyWithLabel(input.DonMetadata.NodesMetadata, &ptypes.Label{Key: devenv.NodeLabelKeyType, Value: ptr.Ptr(string(devenv.NodeLabelValuePlugin))})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find worker nodes")
	}

	for i := range workflowNodeSet {
		var nodeIndex int
		for _, label := range workflowNodeSet[i].Labels {
			if label.Key == node.IndexKey {
				nodeIndex, err = strconv.Atoi(*label.Value)
				if err != nil {
					return nil, errors.Wrap(err, "failed to convert node index to int")
				}
			}
		}

		configOverrides[nodeIndex] = config.WorkerEVM(donBootstrapNodePeerID, donBootstrapNodeHost, input.PeeringData, chainIDUint64, input.CapabilitiesRegistryAddress, input.BlockchainOutput.Nodes[0].DockerInternalHTTPUrl, input.BlockchainOutput.Nodes[0].DockerInternalWSUrl)
		var nodeEthAddr common.Address
		for _, label := range workflowNodeSet[i].Labels {
			if label.Key == node.EthAddressKey {
				if label.Value == nil {
					return nil, errors.New("eth address label value is nil")
				}
				if *label.Value == "" {
					return nil, errors.New("eth address label value is empty")
				}
				nodeEthAddr = common.HexToAddress(*label.Value)
				break
			}
		}

		if keystoneflags.HasFlag(input.Flags, cretypes.WriteEVMCapability) {
			configOverrides[nodeIndex] += config.WorkerWriteEMV(
				nodeEthAddr,
				input.ForwarderAddress,
			)
		}

		// if it's workflow DON configure workflow registry
		if keystoneflags.HasFlag(input.Flags, cretypes.WorkflowDON) {
			configOverrides[nodeIndex] += config.WorkerWorkflowRegistry(
				input.WorkflowRegistryAddress, chainIDUint64)
		}

		// workflow DON nodes always needs gateway connector, otherwise they won't be able to fetch the workflow
		// it's also required by custom compute, which can only run on workflow DON nodes
		if keystoneflags.HasFlag(input.Flags, cretypes.WorkflowDON) || keystoneflags.HasFlag(input.Flags, cretypes.CustomComputeCapability) {
			configOverrides[nodeIndex] += config.WorkerGateway(
				nodeEthAddr,
				chainIDUint64,
				input.DonID,
				*input.GatewayConnectorOutput,
			)
		}
	}

	return configOverrides, nil
}

func GenerateSecrets(input *cretypes.GenerateSecretsInput) (cretypes.NodeIndexToOverride, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}

	overrides := make(cretypes.NodeIndexToOverride)

	for i := range input.DonMetadata.NodesMetadata {
		nodeSecret := types.NodeSecret{}
		if input.EVMKeys != nil {
			nodeSecret.EthKey = types.NodeEthKey{
				JSON:     string((*input.EVMKeys).EncryptedJSONs[i]),
				Password: input.EVMKeys.Password,
				Selector: types.NodeEthKeySelector{
					ChainSelector: input.EVMKeys.ChainSelector,
				},
			}
		}

		if input.P2PKeys != nil {
			nodeSecret.P2PKey = types.NodeP2PKey{
				JSON:     string((*input.P2PKeys).EncryptedJSONs[i]),
				Password: input.P2PKeys.Password,
			}
		}

		nodeSecretString, err := toml.Marshal(nodeSecret)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal node secrets")
		}

		overrides[i] = string(nodeSecretString)
	}

	return overrides, nil
}
