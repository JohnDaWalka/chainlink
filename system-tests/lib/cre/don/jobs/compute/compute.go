package compute

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	crenode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

const flag = cre.CustomComputeCapability

const customComputeConfigTemplate = `"""
NumWorkers = {{.NumWorkers}}
[rateLimiter]
globalRPS = {{.GlobalRPS}}
globalBurst = {{.GlobalBurst}}
perSenderRPS = {{.PerSenderRPS}}
perSenderBurst = {{.PerSenderBurst}}
"""`

var JobSpecFn = func(input *cre.JobSpecFactoryInput) (cre.DonsToJobSpecs, error) {
	return generateJobSpecs(
		input.DonTopology,
		input.InfraInput,
		input.CapabilityConfigs,
		input.CapabilitiesAwareNodeSets,
	)
}

func generateJobSpecs(donTopology *cre.DonTopology, _ *infra.Input, capabilitiesConfig cre.CapabilityConfigs, nodeSetInput []*cre.CapabilitiesAwareNodeSet) (cre.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(cre.DonsToJobSpecs)

	for donIdx, donWithMetadata := range donTopology.DonsWithMetadata {
		if !flags.HasFlag(donWithMetadata.Flags, cre.CustomComputeCapability) {
			continue
		}

		customComputeConfig, ok := capabilitiesConfig[flag]
		if !ok {
			return nil, errors.Errorf("%s config not found in capabilities config", flag)
		}

		workflowNodeSet, err := crenode.FindManyWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: crenode.NodeTypeKey, Value: cre.WorkerNode}, crenode.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		var donOverrides map[string]map[string]any
		if donIdx < len(nodeSetInput) && nodeSetInput[donIdx] != nil {
			donOverrides = nodeSetInput[donIdx].CapabilityOverrides
		}

		mergedConfig := cre.ResolveCapabilityConfigForDON(string(flag), customComputeConfig.Config, donOverrides)
		templateData := don.ApplyRuntimeValues(mergedConfig, map[string]any{})

		tmpl, err := template.New("customComputeConfig").Parse(customComputeConfigTemplate)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse %s config template", flag)
		}

		var configBuffer bytes.Buffer
		if err := tmpl.Execute(&configBuffer, templateData); err != nil {
			return nil, errors.Wrapf(err, "failed to execute %s config template", flag)
		}
		configStr := configBuffer.String()

		if err := don.ValidateTemplateSubstitution(configStr, flag); err != nil {
			return nil, errors.Wrapf(err, "%s template validation failed", flag)
		}

		for _, workerNode := range workflowNodeSet {
			nodeID, nodeIDErr := crenode.FindLabelValue(workerNode, crenode.NodeIDKey)
			if nodeIDErr != nil {
				return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
			}

			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.WorkerStandardCapability(nodeID, cre.CustomComputeCapability, "__builtin_custom-compute-action", configStr, ""))
		}
	}

	return donToJobSpecs, nil
}
