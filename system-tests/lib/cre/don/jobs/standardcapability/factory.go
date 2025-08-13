package factory

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	creregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	crenode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

type JobSpecFactory interface {
	GenerateJobSpecs(input *cre.JobSpecInput) (cre.DonsToJobSpecs, error)
}

// Type aliases for cleaner function signatures
type RuntimeValuesExtractor func(chainID uint64, nodeMetadata *cre.NodeMetadata) map[string]any
type CommandBuilderFn func(input *cre.JobSpecInput, capabilityConfig cre.CapabilityConfig) (string, error)

// NoOpExtractor is a no-operation runtime values extractor for DON-level capabilities
// that don't need runtime values extraction from node metadata
var NoOpExtractor RuntimeValuesExtractor = func(chainID uint64, nodeMetadata *cre.NodeMetadata) map[string]any {
	return map[string]any{} // Return empty map - DON-level capabilities typically don't need runtime values
}

// BinaryPathBuilder constructs the container path for capability binaries by combining
// the default container directory with the base name of the capability's binary path
var BinaryPathBuilder CommandBuilderFn = func(input *cre.JobSpecInput, capabilityConfig cre.CapabilityConfig) (string, error) {
	containerPath, pathErr := creregistry.DefaultContainerDirectory(input.InfraInput.Type)
	if pathErr != nil {
		return "", errors.Wrapf(pathErr, "failed to get default container directory for infra type %s", input.InfraInput.Type)
	}

	return filepath.Join(containerPath, filepath.Base(capabilityConfig.BinaryPath)), nil
}

func NewDonLevelCapabilityJobSpecFactory(
	capabilityFlag cre.CapabilityFlag,
	configTemplate string,
	runtimeValuesExtractor RuntimeValuesExtractor,
	commandBuilder CommandBuilderFn,
) JobSpecFactory {
	return &CapabilityJobSpecFactory{
		capabilityFlag:         capabilityFlag,
		configTemplate:         configTemplate,
		runtimeValuesExtractor: runtimeValuesExtractor,
		commandBuilder:         commandBuilder,
		enabledFn:              enabledForDonFn,
		enabledChainsFn:        enabledChainsForDonFn,
		configResolverFn:       perDonConfigResolverFn,
		jobNameFn: func(chainID uint64, flag cre.CapabilityFlag) string {
			return flag
		},
	}
}

func NewChainSpecificCapabilityJobSpecFactory(
	capabilityFlag cre.CapabilityFlag,
	configTemplate string,
	runtimeValuesExtractor RuntimeValuesExtractor,
	commandBuilder CommandBuilderFn,
) JobSpecFactory {
	return &CapabilityJobSpecFactory{
		capabilityFlag:         capabilityFlag,
		configTemplate:         configTemplate,
		runtimeValuesExtractor: runtimeValuesExtractor,
		commandBuilder:         commandBuilder,
		enabledFn:              enabledForChainsFn,
		enabledChainsFn:        enabledChainIdsFn,
		configResolverFn:       perChainConfigResolverFn,
		jobNameFn: func(chainID uint64, flag cre.CapabilityFlag) string {
			return fmt.Sprintf("%s-%d", flag, chainID)
		},
	}
}

type CapabilityJobSpecFactory struct {
	capabilityFlag         cre.CapabilityFlag
	configTemplate         string
	runtimeValuesExtractor RuntimeValuesExtractor
	commandBuilder         CommandBuilderFn
	jobNameFn              func(chainID uint64, flag cre.CapabilityFlag) string
	enabledFn              func(donWithMetadata *cre.DonWithMetadata, nodeSet *cre.CapabilitiesAwareNodeSet, flag cre.CapabilityFlag) bool
	enabledChainsFn        func(donTopology *cre.DonTopology, nodeSetInput *cre.CapabilitiesAwareNodeSet, flag cre.CapabilityFlag) []uint64
	configResolverFn       func(nodeSetInput *cre.CapabilitiesAwareNodeSet, capabilityConfig cre.CapabilityConfig, chainID uint64, flag cre.CapabilityFlag) (bool, map[string]any, error)
}

var enabledForChainsFn = func(donWithMetadata *cre.DonWithMetadata, nodeSet *cre.CapabilitiesAwareNodeSet, flag cre.CapabilityFlag) bool {
	// Check if this capability is enabled for any chains on this DON
	if donWithMetadata == nil || nodeSet == nil || nodeSet.ChainCapabilities == nil {
		return false
	}

	chainCapConfig, ok := nodeSet.ChainCapabilities[flag]
	if !ok || chainCapConfig == nil || len(chainCapConfig.EnabledChains) == 0 {
		return false
	}

	return true
}

var enabledForDonFn = func(donWithMetadata *cre.DonWithMetadata, nodeSet *cre.CapabilitiesAwareNodeSet, flag cre.CapabilityFlag) bool {
	// Check if this DON has the capability enabled
	return flags.HasFlag(donWithMetadata.Flags, flag)
}

var enabledChainIdsFn = func(donTopology *cre.DonTopology, nodeSetInput *cre.CapabilitiesAwareNodeSet, flag cre.CapabilityFlag) []uint64 {
	chainCapConfig, ok := nodeSetInput.ChainCapabilities[flag]
	if !ok || chainCapConfig == nil {
		return []uint64{}
	}

	return chainCapConfig.EnabledChains
}

var enabledChainsForDonFn = func(donTopology *cre.DonTopology, nodeSetInput *cre.CapabilitiesAwareNodeSet, flag cre.CapabilityFlag) []uint64 {
	return []uint64{donTopology.HomeChainSelector}
}

var perChainConfigResolverFn = func(nodeSetInput *cre.CapabilitiesAwareNodeSet, capabilityConfig cre.CapabilityConfig, chainID uint64, flag cre.CapabilityFlag) (bool, map[string]any, error) {
	// Resolve capability config for this chain
	enabled, mergedConfig, rErr := cre.ResolveCapabilityForChain(
		flag,
		nodeSetInput.ChainCapabilities,
		capabilityConfig.Config,
		chainID,
	)
	if rErr != nil {
		return false, nil, errors.Wrap(rErr, "failed to resolve capability config for chain")
	}
	if !enabled {
		return false, nil, errors.New("capability not enabled for chain")
	}

	return true, mergedConfig, nil
}

var perDonConfigResolverFn = func(nodeSetInput *cre.CapabilitiesAwareNodeSet, capabilityConfig cre.CapabilityConfig, _ uint64, flag cre.CapabilityFlag) (bool, map[string]any, error) {
	// Merge global defaults with DON-specific overrides
	if nodeSetInput == nil {
		return false, nil, errors.New("node set input is nil")
	}

	return true, cre.ResolveCapabilityConfigForDON(flag, capabilityConfig.Config, nodeSetInput.CapabilityOverrides), nil
}

func (f *CapabilityJobSpecFactory) GenerateJobSpecs(input *cre.JobSpecInput) (cre.DonsToJobSpecs, error) {
	if input.DonTopology == nil {
		return nil, errors.New("topology is nil")
	}

	donToJobSpecs := make(cre.DonsToJobSpecs)

	for donIdx, donWithMetadata := range input.DonTopology.DonsWithMetadata {
		// // Check if this capability is enabled for any chains on this DON
		// if donIdx >= len(input.CapabilitiesAwareNodeSets) ||
		// 	input.CapabilitiesAwareNodeSets[donIdx] == nil ||
		// 	input.CapabilitiesAwareNodeSets[donIdx].ChainCapabilities == nil {
		// 	continue
		// }

		// chainCapConfig, ok := input.CapabilitiesAwareNodeSets[donIdx].ChainCapabilities[f.capabilityFlag]
		// if !ok || chainCapConfig == nil || len(chainCapConfig.EnabledChains) == 0 {
		// 	continue
		// }

		if donIdx >= len(input.CapabilitiesAwareNodeSets) || input.CapabilitiesAwareNodeSets[donIdx] == nil {
			continue
		}

		if f.enabledFn != nil && !f.enabledFn(donWithMetadata, input.CapabilitiesAwareNodeSets[donIdx], f.capabilityFlag) {
			continue
		}

		// Get capability config
		capabilityConfig, ok := input.CapabilityConfigs[f.capabilityFlag]
		if !ok {
			return nil, errors.Errorf("%s config not found in capabilities config", f.capabilityFlag)
		}

		// Get capability command
		command, cmdErr := f.commandBuilder(input, capabilityConfig)
		if cmdErr != nil {
			return nil, errors.Wrap(cmdErr, "failed to get capability command")
		}

		// Find worker nodes
		workflowNodeSet, setErr := crenode.FindManyWithLabel(
			donWithMetadata.NodesMetadata,
			&cre.Label{Key: crenode.NodeTypeKey, Value: cre.WorkerNode},
			crenode.EqualLabels,
		)

		if setErr != nil {
			return nil, errors.Wrap(setErr, "failed to find worker nodes")
		}

		// Generate job specs for each enabled chain
		for _, chainIDUint64 := range f.enabledChainsFn(input.DonTopology, input.CapabilitiesAwareNodeSets[donIdx], f.capabilityFlag) {

			// // Resolve capability config for this chain
			// enabled, mergedConfig, rErr := cre.ResolveCapabilityForChain(
			// 	f.capabilityFlag,
			// 	input.CapabilitiesAwareNodeSets[donIdx].ChainCapabilities,
			// 	capabilityConfig.Config,
			// 	chainIDUint64,
			// )
			// if rErr != nil {
			// 	return nil, errors.Wrap(rErr, "failed to resolve capability config for chain")
			// }
			// if !enabled {
			// 	continue
			// }

			enabled, mergedConfig, rErr := f.configResolverFn(input.CapabilitiesAwareNodeSets[donIdx], capabilityConfig, chainIDUint64, f.capabilityFlag)
			if rErr != nil {
				return nil, errors.Wrap(rErr, "failed to resolve capability config for chain")
			}
			if !enabled {
				continue
			}

			// Create job specs for each worker node
			for _, workerNode := range workflowNodeSet {
				nodeID, nodeIDErr := crenode.FindLabelValue(workerNode, crenode.NodeIDKey)
				if nodeIDErr != nil {
					return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
				}

				// Apply runtime values to merged config using the runtime value builder
				templateData := don.ApplyRuntimeValues(mergedConfig, f.runtimeValuesExtractor(chainIDUint64, workerNode))

				// Parse and execute template
				tmpl, tmplErr := template.New(string(f.capabilityFlag) + "-config").Parse(f.configTemplate)
				if tmplErr != nil {
					return nil, errors.Wrapf(tmplErr, "failed to parse %s config template", f.capabilityFlag)
				}

				var configBuffer bytes.Buffer
				if err := tmpl.Execute(&configBuffer, templateData); err != nil {
					return nil, errors.Wrapf(err, "failed to execute %s config template", f.capabilityFlag)
				}
				configStr := configBuffer.String()

				if err := don.ValidateTemplateSubstitution(configStr, f.capabilityFlag); err != nil {
					return nil, errors.Wrapf(err, "%s template validation failed", f.capabilityFlag)
				}

				jobSpec := jobs.WorkerStandardCapability(nodeID, f.jobNameFn(chainIDUint64, f.capabilityFlag), command, configStr, "")
				donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobSpec)
			}
		}
	}

	return donToJobSpecs, nil
}
