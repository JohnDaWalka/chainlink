package factory

import (
	"bytes"
	"fmt"
	"strconv"
	"text/template"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	crenode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
)

// ChainSpecificFactory implements chain-specific job spec generation for capabilities
// that require per-chain configuration and deployment
type ChainSpecificFactory struct {
	capabilityFlag         cre.CapabilityFlag
	configTemplate         string
	runtimeValuesExtractor RuntimeValuesExtractor
	commandBuilder         CommandBuilder
}

// NewChainSpecificFactory creates a new chain-specific factory
func NewChainSpecificFactory(
	capabilityFlag cre.CapabilityFlag,
	configTemplate string,
	runtimeValuesExtractor RuntimeValuesExtractor,
	commandBuilder CommandBuilder,
) *ChainSpecificFactory {
	if runtimeValuesExtractor == nil {
		runtimeValuesExtractor = NoOpExtractor
	}

	return &ChainSpecificFactory{
		capabilityFlag:         capabilityFlag,
		configTemplate:         configTemplate,
		runtimeValuesExtractor: runtimeValuesExtractor,
		commandBuilder:         commandBuilder,
	}
}

func (f *ChainSpecificFactory) CapabilityFlag() cre.CapabilityFlag {
	return f.capabilityFlag
}

func (f *ChainSpecificFactory) ConfigTemplate() string {
	return f.configTemplate
}

func (f *ChainSpecificFactory) RuntimeValuesExtractor() RuntimeValuesExtractor {
	return f.runtimeValuesExtractor
}

func (f *ChainSpecificFactory) CommandBuilder() CommandBuilder {
	return f.commandBuilder
}

func (f *ChainSpecificFactory) GenerateJobSpecs(input *cre.JobSpecInput) (cre.DonsToJobSpecs, error) {
	if input.DonTopology == nil {
		return nil, errors.New("topology is nil")
	}

	donToJobSpecs := make(cre.DonsToJobSpecs)

	for donIdx, donWithMetadata := range input.DonTopology.DonsWithMetadata {
		// Check if this capability is enabled for any chains on this DON
		if donIdx >= len(input.CapabilitiesAwareNodeSets) ||
			input.CapabilitiesAwareNodeSets[donIdx] == nil ||
			input.CapabilitiesAwareNodeSets[donIdx].ChainCapabilities == nil {
			continue
		}

		chainCapConfig, ok := input.CapabilitiesAwareNodeSets[donIdx].ChainCapabilities[f.capabilityFlag]
		if !ok || chainCapConfig == nil || len(chainCapConfig.EnabledChains) == 0 {
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
		for _, chainIDUint64 := range chainCapConfig.EnabledChains {
			chainID := int(chainIDUint64)
			chainIDStr := strconv.Itoa(chainID)

			// Resolve capability config for this chain
			enabled, mergedConfig, rErr := cre.ResolveCapabilityForChain(
				f.capabilityFlag,
				input.CapabilitiesAwareNodeSets[donIdx].ChainCapabilities,
				capabilityConfig.Config,
				chainIDUint64,
			)
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

				// Create job name with chain suffix
				jobName := fmt.Sprintf("%s-%s", f.capabilityFlag, chainIDStr)

				jobSpec := jobs.WorkerStandardCapability(nodeID, jobName, command, configStr, "")
				donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobSpec)
			}
		}
	}

	return donToJobSpecs, nil
}
