package evm

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/Masterminds/semver/v3"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

const flag = cre.EVMCapability
const evmConfigTemplate = `'{"chainId":{{.ChainID}},"network":"{{.NetworkFamily}}","logTriggerPollInterval":{{.LogTriggerPollInterval}}, "creForwarderAddress":"{{.CreForwarderAddress}}","receiverGasMinimum":{{.ReceiverGasMinimum}},"nodeAddress":"{{.NodeAddress}}"}'`

// buildEVMRuntimeFallbacks creates runtime-generated fallback values for any keys not specified in TOML
func buildEVMRuntimeFallbacks(chainID int, networkFamily, creForwarderAddress, nodeAddress string) map[string]any {
	return map[string]any{
		"ChainID":             chainID,
		"NetworkFamily":       networkFamily,
		"CreForwarderAddress": creForwarderAddress,
		"NodeAddress":         nodeAddress,
	}
}

var EVMJobSpecFactoryFn = func(input *cre.JobSpecFactoryInput) (cre.DonsToJobSpecs, error) {
	return generateJobSpecs(input.DonTopology, input.CldEnvironment.DataStore, *input.InfraInput, input.AdditionalCapabilities, input.CapabilitiesAwareNodeSets)
}

var jobName = func(chainID string) string {
	return "evm-capability-" + chainID
}

func generateJobSpecs(
	donTopology *cre.DonTopology,
	ds datastore.DataStore,
	infraInput infra.Input,
	capabilitiesConfig cre.AdditionalCapabilitiesConfigs,
	nodeSetInput []*cre.CapabilitiesAwareNodeSet,
) (cre.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(cre.DonsToJobSpecs)

	logger := framework.L

	for donIdx, donWithMetadata := range donTopology.DonsWithMetadata {
		// EVM capability is enabled strictly per-chain via ChainCapabilities
		if donIdx >= len(nodeSetInput) || nodeSetInput[donIdx] == nil || nodeSetInput[donIdx].ChainCapabilities == nil {
			continue
		}
		if cc, ok := nodeSetInput[donIdx].ChainCapabilities[string(flag)]; !ok || cc == nil || len(cc.EnabledChains) == 0 {
			continue
		}

		evmConfig, ok := capabilitiesConfig[flag]
		if !ok {
			return nil, errors.New("evm config not found in capabilities config")
		}

		containerPath, pathErr := crecapabilities.DefaultContainerDirectory(infraInput.Type)
		if pathErr != nil {
			return nil, errors.Wrapf(pathErr, "failed to get default container directory for infra type %s", infraInput.Type)
		}

		evmBinaryPath := filepath.Join(containerPath, filepath.Base(evmConfig.BinaryPath))

		evmOCR3Key := datastore.NewAddressRefKey(
			donTopology.HomeChainSelector,
			datastore.ContractType(keystone_changeset.OCR3Capability.String()),
			semver.MustParse("1.0.0"),
			"capability_evm",
		)
		evmOCR3CapabilityAddress, err := ds.Addresses().Get(evmOCR3Key)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get EVM capability address")
		}

		internalHostsBS := getBoostrapWorkflowNames(donWithMetadata, nodeSetInput, donIdx, infraInput)
		if len(internalHostsBS) == 0 {
			return nil, fmt.Errorf("no bootstrap node found for DON %s (there should be at least 1)", donWithMetadata.Name)
		}

		workflowNodeSet, err := node.FindManyWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: node.NodeTypeKey, Value: cre.WorkerNode}, node.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		// look for boostrap node and then for required values in its labels
		bootstrapNode, bootErr := node.FindOneWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: node.NodeTypeKey, Value: cre.BootstrapNode}, node.EqualLabels)
		if bootErr != nil {
			return nil, errors.Wrap(bootErr, "failed to find bootstrap node")
		}

		bootstrapNodeID, nodeIDErr := node.FindLabelValue(bootstrapNode, node.NodeIDKey)
		if nodeIDErr != nil {
			return nil, errors.Wrap(nodeIDErr, "failed to get bootstrap node id from labels")
		}

		// New: iterate enabled chains from nodeset chain capabilities resolver
		nodeSet := nodeSetInput[donIdx]

		// Defaults for capability configs come from AdditionalCapabilities[flag].Config
		// These are global defaults per capability.
		defaults := map[string]map[string]any{}
		if capCfg, ok := capabilitiesConfig[flag]; ok {
			defaults[string(flag)] = capCfg.Config
		}

		enabledChains := []uint64{}
		if nodeSet.ChainCapabilities != nil {
			if cc, ok := nodeSet.ChainCapabilities[string(flag)]; ok {
				enabledChains = append(enabledChains, cc.EnabledChains...)
			}
		}

		for _, chainIDUint64 := range enabledChains {
			chainID := int(chainIDUint64)
			chainIDStr := strconv.Itoa(chainID)
			chain, ok := chainsel.ChainByEvmChainID(chainIDUint64)
			if !ok {
				return nil, fmt.Errorf("failed to get chain selector for chain ID %d", chainIDUint64)
			}

			// create job specs for the bootstrap node
			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.BootstrapOCR3(bootstrapNodeID, "evm-capability", evmOCR3CapabilityAddress.Address, chainIDUint64))
			logger.Debug().Msgf("Deployed EVM OCR3 contract on chain %d at %s", chainIDUint64, evmOCR3CapabilityAddress.Address)

			creForwarderKey := datastore.NewAddressRefKey(
				chain.Selector,
				datastore.ContractType(keystone_changeset.KeystoneForwarder.String()),
				semver.MustParse("1.0.0"),
				"",
			)
			creForwarderAddress, err := ds.Addresses().Get(creForwarderKey)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get CRE Forwarder address")
			}

			logger.Debug().Msgf("Found CRE Forwarder contract on chain %d at %s", chainID, creForwarderAddress.Address)

			// Build user configuration from defaults + chain overrides
			enabled, mergedConfig, rErr := cre.ResolveCapabilityForChain(string(flag), nodeSet.ChainCapabilities, defaults, chainIDUint64)
			if rErr != nil {
				return nil, errors.Wrap(rErr, "failed to resolve capability config for chain")
			}
			if !enabled {
				// should not happen because we derived enabledChains from the same source, but guard anyway
				continue
			}

			// To preserve current behavior, treat mergedConfig as the "user" config post-merge
			// and pass globalConfig = mergedConfig, chain-specific section unused here
			userConfig, err := jobs.BuildGlobalConfigFromTOML(map[string]any{"config": mergedConfig})
			if err != nil {
				return nil, errors.Wrap(err, "failed to build config from TOML")
			}

			for _, workerNode := range workflowNodeSet {
				nodeID, nodeIDErr := node.FindLabelValue(workerNode, node.NodeIDKey)
				if nodeIDErr != nil {
					return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
				}

				transmitterAddress, tErr := node.FindLabelValue(workerNode, node.AddressKeyFromSelector(chain.Selector))
				if tErr != nil {
					return nil, errors.Wrap(tErr, "failed to get transmitter address from bootstrap node labels")
				}

				keyBundle, kErr := node.FindLabelValue(workerNode, node.NodeOCR2KeyBundleIDKey)
				if kErr != nil {
					return nil, errors.Wrap(kErr, "failed to get key bundle id from worker node labels")
				}

				keyNodeAddress := node.AddressKeyFromSelector(chain.Selector)
				nodeAddress, nodeAddressErr := node.FindLabelValue(workerNode, keyNodeAddress)
				if nodeAddressErr != nil {
					return nil, errors.Wrap(nodeAddressErr, "failed to get node address from labels")
				}
				logger.Debug().Msgf("Deployed node on chain %d/%d at %s", chainID, chain.Selector, nodeAddress)

				bootstrapNodeP2pKeyID, pErr := node.FindLabelValue(bootstrapNode, node.NodeP2PIDKey)
				if pErr != nil {
					return nil, errors.Wrap(pErr, "failed to get p2p key id from bootstrap node labels")
				}
				// remove the prefix if it exists, to match the expected format
				bootstrapNodeP2pKeyID = strings.TrimPrefix(bootstrapNodeP2pKeyID, "p2p_")
				bootstrapPeers := make([]string, len(internalHostsBS))
				for i, workflowName := range internalHostsBS {
					bootstrapPeers[i] = fmt.Sprintf("%s@%s:5001", bootstrapNodeP2pKeyID, workflowName)
				}

				oracleFactoryConfigInstance := job.OracleFactoryConfig{
					Enabled:            true,
					ChainID:            chainIDStr,
					BootstrapPeers:     bootstrapPeers,
					OCRContractAddress: evmOCR3CapabilityAddress.Address,
					OCRKeyBundleID:     keyBundle,
					TransmitterID:      transmitterAddress,
					OnchainSigning: job.OnchainSigningStrategy{
						StrategyName: "single-chain",
						Config:       map[string]string{"evm": keyBundle},
					},
				}

				type OracleFactoryConfigWrapper struct {
					OracleFactory job.OracleFactoryConfig `toml:"oracle_factory"`
				}
				wrapper := OracleFactoryConfigWrapper{OracleFactory: oracleFactoryConfigInstance}

				var oracleBuffer bytes.Buffer
				if errEncoder := toml.NewEncoder(&oracleBuffer).Encode(wrapper); errEncoder != nil {
					return nil, errors.Wrap(errEncoder, "failed to encode oracle factory config to TOML")
				}
				oracleStr := strings.ReplaceAll(oracleBuffer.String(), "\n", "\n\t")

				// Build runtime fallbacks for any missing values
				runtimeFallbacks := buildEVMRuntimeFallbacks(chainID, "evm", creForwarderAddress.Address, nodeAddress)

				// Apply runtime fallbacks only for keys not specified by user
				templateData := jobs.ApplyRuntimeFallbacks(userConfig, runtimeFallbacks)

				// Parse and execute template
				tmpl, err := template.New("evmConfig").Parse(evmConfigTemplate)
				if err != nil {
					return nil, errors.Wrap(err, "failed to parse EVM config template")
				}

				var configBuffer bytes.Buffer
				if err := tmpl.Execute(&configBuffer, templateData); err != nil {
					return nil, errors.Wrap(err, "failed to execute EVM config template")
				}
				configStr := configBuffer.String()

				logger.Debug().Msgf("Creating EVM Capability job spec for chainID: %d, selector: %d, DON:%q, node:%q", chainID, chain.Selector, donWithMetadata.Name, nodeID)

				jobSpec := jobs.WorkerStandardCapability(nodeID, jobName(chainIDStr), evmBinaryPath,
					configStr,
					oracleStr,
				)

				if _, ok := donToJobSpecs[donWithMetadata.ID]; !ok {
					donToJobSpecs[donWithMetadata.ID] = make(cre.DonJobs, 0)
				}

				donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobSpec)
			}
		}
	}

	return donToJobSpecs, nil
}

func getBoostrapWorkflowNames(donWithMetadata *cre.DonWithMetadata, nodeSetInput []*cre.CapabilitiesAwareNodeSet, donIdx int, infraInput infra.Input) []string {
	internalHostsBS := make([]string, 0)
	for nodeIdx := range donWithMetadata.NodesMetadata {
		if nodeSetInput[donIdx].BootstrapNodeIndex != -1 && nodeIdx == nodeSetInput[donIdx].BootstrapNodeIndex {
			internalHostBS := don.InternalHost(nodeIdx, cre.BootstrapNode, donWithMetadata.Name, infraInput)
			internalHostsBS = append(internalHostsBS, internalHostBS)
		}
	}
	return internalHostsBS
}
