package environment

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	creconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config"
	cresecrets "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/secrets"
	creflags "github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

func PrepareNodeTOMLConfigurations(
	registryChainSelector uint64,
	nodeSets []*cre.CapabilitiesAwareNodeSet,
	infraInput infra.Provider,
	blockchainOutputs []*cre.WrappedBlockchainOutput,
	addressBook deployment.AddressBook,
	datastore datastore.DataStore,
	capabilities []cre.InstallableCapability,
	capabilityConfigs cre.CapabilityConfigs,
	copyCapabilityBinaries bool,
) (*cre.Topology, []*cre.CapabilitiesAwareNodeSet, error) {
	topology, tErr := cre.NewTopology(nodeSets, infraInput)
	if tErr != nil {
		return nil, nil, errors.Wrap(tErr, "failed to create topology")
	}

	bt, err := topology.BootstrapNode()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to find bootstrap node")
	}

	capabilitiesPeeringData, ocrPeeringData, peeringErr := cre.PeeringCfgs(bt)
	if peeringErr != nil {
		return nil, nil, errors.Wrap(peeringErr, "failed to find peering data")
	}

	localNodeSets := topology.CapabilitiesAwareNodeSets()
	chainPerSelector := make(map[uint64]*cre.WrappedBlockchainOutput)
	for _, bcOut := range blockchainOutputs {
		if bcOut.SolChain != nil {
			sel := bcOut.SolChain.ChainSelector
			chainPerSelector[sel] = bcOut
			chainPerSelector[sel].ChainSelector = sel
			chainPerSelector[sel].SolChain = bcOut.SolChain
			chainPerSelector[sel].SolChain.ArtifactsDir = bcOut.SolChain.ArtifactsDir
			continue
		}
		chainPerSelector[bcOut.ChainSelector] = bcOut
	}

	for i, donMetadata := range topology.DonsMetadata.List() {
		configsFound := 0
		secretsFound := 0
		nodeSet := localNodeSets[i]
		for _, nodeSpec := range nodeSet.NodeSpecs {
			if nodeSpec.Node.TestConfigOverrides != "" {
				configsFound++
			}
			if nodeSpec.Node.TestSecretsOverrides != "" {
				secretsFound++
			}
		}
		if configsFound != 0 && configsFound != len(localNodeSets[i].NodeSpecs) {
			return nil, nil, fmt.Errorf("%d out of %d node specs have config overrides. Either provide overrides for all nodes or none at all", configsFound, len(localNodeSets[i].NodeSpecs))
		}

		if secretsFound != 0 && secretsFound != len(localNodeSets[i].NodeSpecs) {
			return nil, nil, fmt.Errorf("%d out of %d node specs have secrets overrides. Either provide overrides for all nodes or none at all", secretsFound, len(localNodeSets[i].NodeSpecs))
		}

		// Allow providing only secrets, because we can decode them and use them to generate configs
		// We can't allow providing only configs, because we can't replace secret-related values in the configs
		// If both are provided, we assume that the user knows what they are doing and we don't need to validate anything
		// And that configs match the secrets
		if configsFound > 0 && secretsFound == 0 {
			return nil, nil, fmt.Errorf("nodese config overrides are provided for DON %d, but not secrets. You need to either provide both, only secrets or nothing at all", donMetadata.ID)
		}

		configFactoryFunctions := make([]cre.NodeConfigTransformerFn, 0)
		for _, capability := range capabilities {
			configFactoryFunctions = append(configFactoryFunctions, capability.NodeConfigTransformerFn())
		}

		// generate configs only if they are not provided
		if configsFound == 0 {
			config, configErr := creconfig.Generate(
				cre.GenerateConfigsInput{
					AddressBook:             addressBook,
					Datastore:               datastore,
					DonMetadata:             donMetadata,
					BlockchainOutput:        chainPerSelector,
					Flags:                   donMetadata.Flags,
					CapabilitiesPeeringData: capabilitiesPeeringData,
					OCRPeeringData:          ocrPeeringData,
					HomeChainSelector:       registryChainSelector,
					GatewayConnectorOutput:  topology.GatewayConnectorOutput,
					NodeSet:                 localNodeSets[i],
					CapabilityConfigs:       capabilityConfigs,
				},
				configFactoryFunctions,
			)
			if configErr != nil {
				return nil, nil, errors.Wrap(configErr, "failed to generate config")
			}

			for j := range donMetadata.NodesMetadata {
				localNodeSets[i].NodeSpecs[j].Node.TestConfigOverrides = config[j]
			}
		}

		// generate secrets only if they are not provided
		if secretsFound == 0 {
			for j := range donMetadata.NodesMetadata {
				wnode := donMetadata.NodesMetadata[j]
				nodeSecret := cresecrets.NewNodeSecret(wnode.Keys)
				toml, err := nodeSecret.Toml()
				if err != nil {
					return nil, nil, errors.Wrap(err, "failed to marshal node secrets")
				}
				localNodeSets[i].NodeSpecs[j].Node.TestSecretsOverrides = toml
			}
		}

		if !copyCapabilityBinaries {
			continue
		}

		customBinariesPaths := make(map[cre.CapabilityFlag]string)
		for flag, config := range capabilityConfigs {
			if creflags.HasFlagForAnyChain(donMetadata.Flags, flag) && config.BinaryPath != "" {
				customBinariesPaths[flag] = config.BinaryPath
			}
		}

		executableErr := crecapabilities.MakeBinariesExecutable(customBinariesPaths)
		if executableErr != nil {
			return nil, nil, errors.Wrap(executableErr, "failed to make binaries executable")
		}

		var err error
		ns, err := crecapabilities.AppendBinariesPathsNodeSpec(nodeSet, donMetadata, customBinariesPaths)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to append binaries paths to node spec for DON %d", donMetadata.ID)
		}
		localNodeSets[i] = ns
	}

	// Add env vars, which were provided programmatically, to the node specs
	// or fail, if node specs already had some env vars set in the TOML config
	for donIdx, donMetadata := range topology.DonsMetadata.List() {
		hasEnvVarsInTomlConfig := false
		for nodeIdx, nodeSpec := range localNodeSets[donIdx].NodeSpecs {
			if len(nodeSpec.Node.EnvVars) > 0 {
				hasEnvVarsInTomlConfig = true
				break
			}

			localNodeSets[donIdx].NodeSpecs[nodeIdx].Node.EnvVars = localNodeSets[donIdx].EnvVars
		}

		if hasEnvVarsInTomlConfig && len(localNodeSets[donIdx].EnvVars) > 0 {
			return nil, nil, fmt.Errorf("extra env vars for Chainlink Nodes are provided in the TOML config for the %s DON, but you tried to provide them programatically. Please set them only in one place", donMetadata.Name)
		}
	}

	// Deploy the DONs
	// Hack for CI that allows us to dynamically set the chainlink image and version
	// CTFv2 currently doesn't support dynamic image and version setting
	if os.Getenv("CI") == "true" {
		// Due to how we pass custom env vars to reusable workflow we need to use placeholders, so first we need to resolve what's the name of the target environment variable
		// that stores chainlink version and then we can use it to resolve the image name
		for i := range localNodeSets {
			image := fmt.Sprintf("%s:%s", os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV), ctfconfig.MustReadEnvVar_String(ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV))
			for j := range localNodeSets[i].NodeSpecs {
				localNodeSets[i].NodeSpecs[j].Node.Image = image
				// unset docker context and file path, so that we can use the image from the registry
				localNodeSets[i].NodeSpecs[j].Node.DockerContext = ""
				localNodeSets[i].NodeSpecs[j].Node.DockerFilePath = ""
			}
		}
	}

	return topology, localNodeSets, nil
}
