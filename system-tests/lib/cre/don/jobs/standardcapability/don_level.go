package factory

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	crenode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

// DonLevelFactory implements DON-level job spec generation for capabilities
// that don't require chain-specific configuration
type DonLevelFactory struct {
	capabilityFlag         cre.CapabilityFlag
	configTemplate         string
	runtimeValuesExtractor RuntimeValuesExtractor
	commandBuilder         CommandBuilder
}

// NewDonLevelFactory creates a new DON-level factory
func NewDonLevelFactory(
	capabilityFlag cre.CapabilityFlag,
	configTemplate string,
	runtimeValuesExtractor RuntimeValuesExtractor,
	commandBuilder CommandBuilder,
) *DonLevelFactory {
	if runtimeValuesExtractor == nil {
		runtimeValuesExtractor = NoOpExtractor
	}

	if commandBuilder == nil {
		commandBuilder = BinaryPathBuilder
	}

	return &DonLevelFactory{
		capabilityFlag:         capabilityFlag,
		configTemplate:         configTemplate,
		runtimeValuesExtractor: runtimeValuesExtractor,
		commandBuilder:         commandBuilder,
	}
}

func (f *DonLevelFactory) CapabilityFlag() cre.CapabilityFlag {
	return f.capabilityFlag
}

func (f *DonLevelFactory) ConfigTemplate() string {
	return f.configTemplate
}

func (f *DonLevelFactory) RuntimeValuesExtractor() RuntimeValuesExtractor {
	return f.runtimeValuesExtractor
}

func (f *DonLevelFactory) CommandBuilder() CommandBuilder {
	return f.commandBuilder
}

func (f *DonLevelFactory) GenerateJobSpecs(input *cre.JobSpecInput) (cre.DonsToJobSpecs, error) {
	if input.DonTopology == nil {
		return nil, errors.New("topology is nil")
	}

	donToJobSpecs := make(cre.DonsToJobSpecs)

	for donIdx, donWithMetadata := range input.DonTopology.DonsWithMetadata {
		// Check if this DON has the capability enabled
		if !flags.HasFlag(donWithMetadata.Flags, f.capabilityFlag) {
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
		workflowNodeSet, err := crenode.FindManyWithLabel(
			donWithMetadata.NodesMetadata,
			&cre.Label{Key: crenode.NodeTypeKey, Value: cre.WorkerNode},
			crenode.EqualLabels,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		// Merge global defaults with DON-specific overrides
		var donOverrides map[string]map[string]any
		if donIdx < len(input.CapabilitiesAwareNodeSets) && input.CapabilitiesAwareNodeSets[donIdx] != nil {
			donOverrides = input.CapabilitiesAwareNodeSets[donIdx].CapabilityOverrides
		}

		mergedConfig := cre.ResolveCapabilityConfigForDON(f.capabilityFlag, capabilityConfig.Config, donOverrides)

		// Create job specs for each worker node
		for _, workerNode := range workflowNodeSet {
			nodeID, nodeIDErr := crenode.FindLabelValue(workerNode, crenode.NodeIDKey)
			if nodeIDErr != nil {
				return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
			}

			// Apply runtime values to merged config using the runtime value builder
			templateData := don.ApplyRuntimeValues(mergedConfig, f.runtimeValuesExtractor(0, workerNode))

			// Generate config string for this specific node
			var nodeConfigStr string
			if f.configTemplate == "" {
				// Empty template means use empty config
				nodeConfigStr = jobs.EmptyStdCapConfig
			} else {
				// Parse and execute template
				tmpl, err := template.New(string(f.capabilityFlag) + "Config").Parse(f.configTemplate)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to parse %s config template", f.capabilityFlag)
				}

				var configBuffer bytes.Buffer
				if err := tmpl.Execute(&configBuffer, templateData); err != nil {
					return nil, errors.Wrapf(err, "failed to execute %s config template", f.capabilityFlag)
				}
				nodeConfigStr = configBuffer.String()

				if err := don.ValidateTemplateSubstitution(nodeConfigStr, f.capabilityFlag); err != nil {
					return nil, errors.Wrapf(err, "%s template validation failed", f.capabilityFlag)
				}
			}

			jobSpec := jobs.WorkerStandardCapability(nodeID, f.capabilityFlag, command, nodeConfigStr, "")
			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobSpec)
		}
	}

	return donToJobSpecs, nil
}
