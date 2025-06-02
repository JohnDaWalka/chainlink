package environment

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	pkgerrors "github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	libcaps "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	libdon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	creconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config"
	cresecrets "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/secrets"
	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

type BuildTopologyOpDeps struct {
	BlockchainsOutput      []*BlockchainOutput
	AddressBook            deployment.AddressBook
	ConfigFactoryFunctions []cretypes.ConfigFactoryFn
	CustomBinariesPaths    map[cretypes.CapabilityFlag]string
}

type BuildTopologyOpInput struct {
	// Should InfraInput be a dep instead? Do we want it to be serialized? (in/out are serialized)
	InfraInput                libtypes.InfraInput
	CapabilitiesAwareNodeSets []*cretypes.CapabilitiesAwareNodeSet
	ChainSelector             uint64
}

type BuildTopologyOpOutput struct {
	Topology                  *cretypes.Topology
	CapabilitiesAwareNodeSets []*cretypes.CapabilitiesAwareNodeSet
}

var BuildTopologyOp = operations.NewOperation[BuildTopologyOpInput, BuildTopologyOpOutput, BuildTopologyOpDeps](
	"build-topology-op",
	semver.MustParse("1.0.0"),
	"Build Topology",
	func(b operations.Bundle, deps BuildTopologyOpDeps, input BuildTopologyOpInput) (BuildTopologyOpOutput, error) {
		topologyErr := libdon.ValidateTopology(input.CapabilitiesAwareNodeSets, input.InfraInput)
		if topologyErr != nil {
			return BuildTopologyOpOutput{}, pkgerrors.Wrap(topologyErr, "failed to validate topology")
		}

		topology, topoErr := libdon.BuildTopology(input.CapabilitiesAwareNodeSets, input.InfraInput, input.ChainSelector)
		if topoErr != nil {
			return BuildTopologyOpOutput{}, pkgerrors.Wrap(topoErr, "failed to build topology")
		}

		// Generate EVM and P2P keys or read them from the config
		// That way we can pass them final configs and do away with restarting the nodes
		var keys *cretypes.GenerateKeysOutput

		keysOutput, keysOutputErr := cresecrets.KeysOutputFromConfig(input.CapabilitiesAwareNodeSets)
		if keysOutputErr != nil {
			return BuildTopologyOpOutput{}, pkgerrors.Wrap(keysOutputErr, "failed to generate keys output")
		}

		// get chainIDs, they'll be used for identifying ETH keys and Forwarder addresses
		// and also for creating the CLD environment
		chainIDs := make([]int, 0)
		bcOuts := make(map[uint64]*blockchain.Output)
		sethClients := make(map[uint64]*seth.Client)
		for _, bcOut := range deps.BlockchainsOutput {
			chainIDs = append(chainIDs, libc.MustSafeInt(bcOut.ChainID))
			bcOuts[bcOut.ChainSelector] = bcOut.BlockchainOutput
			sethClients[bcOut.ChainSelector] = bcOut.SethClient
		}

		generateKeysInput := &cretypes.GenerateKeysInput{
			GenerateEVMKeysForChainIDs: chainIDs,
			GenerateP2PKeys:            true,
			Topology:                   topology,
			Password:                   "", // since the test runs on private ephemeral blockchain we don't use real keys and do not care a lot about the password
			Out:                        keysOutput,
		}
		keys, keysErr := cresecrets.GenereteKeys(generateKeysInput)
		if keysErr != nil {
			return BuildTopologyOpOutput{}, pkgerrors.Wrap(keysErr, "failed to generate keys")
		}

		topology, addKeysErr := cresecrets.AddKeysToTopology(topology, keys)
		if addKeysErr != nil {
			return BuildTopologyOpOutput{}, pkgerrors.Wrap(addKeysErr, "failed to add keys to topology")
		}

		peeringData, peeringErr := libdon.FindPeeringData(topology)
		if peeringErr != nil {
			return BuildTopologyOpOutput{}, pkgerrors.Wrap(peeringErr, "failed to find peering data")
		}

		for i, donMetadata := range topology.DonsMetadata {
			configsFound := 0
			secretsFound := 0
			for _, nodeSpec := range input.CapabilitiesAwareNodeSets[i].NodeSpecs {
				if nodeSpec.Node.TestConfigOverrides != "" {
					configsFound++
				}
				if nodeSpec.Node.TestSecretsOverrides != "" {
					secretsFound++
				}
			}
			if configsFound != 0 && configsFound != len(input.CapabilitiesAwareNodeSets[i].NodeSpecs) {
				return BuildTopologyOpOutput{}, fmt.Errorf("%d out of %d node specs have config overrides. Either provide overrides for all nodes or none at all", configsFound, len(input.CapabilitiesAwareNodeSets[i].NodeSpecs))
			}

			if secretsFound != 0 && secretsFound != len(input.CapabilitiesAwareNodeSets[i].NodeSpecs) {
				return BuildTopologyOpOutput{}, fmt.Errorf("%d out of %d node specs have secrets overrides. Either provide overrides for all nodes or none at all", secretsFound, len(input.CapabilitiesAwareNodeSets[i].NodeSpecs))
			}

			// Allow providing only secrets, because we can decode them and use them to generate configs
			// We can't allow providing only configs, because we can't replace secret-related values in the configs
			// If both are provided, we assume that the user knows what they are doing and we don't need to validate anything
			// And that configs match the secrets
			if configsFound > 0 && secretsFound == 0 {
				return BuildTopologyOpOutput{}, fmt.Errorf("nodese config overrides are provided for DON %d, but not secrets. You need to either provide both, only secrets or nothing at all", donMetadata.ID)
			}

			// generate configs only if they are not provided
			if configsFound == 0 {
				config, configErr := creconfig.Generate(
					cretypes.GenerateConfigsInput{
						DonMetadata:            donMetadata,
						BlockchainOutput:       bcOuts,
						Flags:                  donMetadata.Flags,
						PeeringData:            peeringData,
						AddressBook:            deps.AddressBook, //nolint:staticcheck // won't migrate now
						HomeChainSelector:      topology.HomeChainSelector,
						GatewayConnectorOutput: topology.GatewayConnectorOutput,
					},
					deps.ConfigFactoryFunctions,
				)
				if configErr != nil {
					return BuildTopologyOpOutput{}, pkgerrors.Wrap(configErr, "failed to generate config")
				}

				for j := range donMetadata.NodesMetadata {
					input.CapabilitiesAwareNodeSets[i].NodeSpecs[j].Node.TestConfigOverrides = config[j]
				}
			}

			// generate secrets only if they are not provided
			if secretsFound == 0 {
				secretsInput := &cretypes.GenerateSecretsInput{
					DonMetadata: donMetadata,
				}

				if evmKeys, ok := keys.EVMKeys[donMetadata.ID]; ok {
					secretsInput.EVMKeys = evmKeys
				}

				if p2pKeys, ok := keys.P2PKeys[donMetadata.ID]; ok {
					secretsInput.P2PKeys = p2pKeys
				}

				// EVM and P2P keys will be provided to nodes as secrets
				secrets, secretsErr := cresecrets.GenerateSecrets(
					secretsInput,
				)
				if secretsErr != nil {
					return BuildTopologyOpOutput{}, pkgerrors.Wrap(secretsErr, "failed to generate secrets")
				}

				for j := range donMetadata.NodesMetadata {
					input.CapabilitiesAwareNodeSets[i].NodeSpecs[j].Node.TestSecretsOverrides = secrets[j]
				}
			}

			var appendErr error
			input.CapabilitiesAwareNodeSets[i], appendErr = libcaps.AppendBinariesPathsNodeSpec(input.CapabilitiesAwareNodeSets[i], donMetadata, deps.CustomBinariesPaths)
			if appendErr != nil {
				return BuildTopologyOpOutput{}, pkgerrors.Wrapf(appendErr, "failed to append binaries paths to node spec for DON %d", donMetadata.ID)
			}
		}

		// TODO: has to return an updated `CapabilitiesAwareNodeSets`

		return BuildTopologyOpOutput{Topology: topology}, nil
	},
)
