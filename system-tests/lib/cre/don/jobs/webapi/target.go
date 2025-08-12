package webapi

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

const flag = cre.WebAPITargetCapability

// Web API target capability uses configurable rate limiter settings via TOML
const webAPITargetConfigTemplate = `"""
[rateLimiter]
GlobalRPS = {{.GlobalRPS}}
GlobalBurst = {{.GlobalBurst}}
PerSenderRPS = {{.PerSenderRPS}}
PerSenderBurst = {{.PerSenderBurst}}
"""`

var TargetJobSpecFn = func(input *cre.JobSpecFactoryInput) (cre.DonsToJobSpecs, error) {
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
		if !flags.HasFlag(donWithMetadata.Flags, cre.WebAPITargetCapability) {
			continue
		}

		webAPITargetConfig, ok := capabilitiesConfig[flag]
		if !ok {
			return nil, errors.Errorf("%s config not found in capabilities config", flag)
		}

		workflowNodeSet, err := libnode.FindManyWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: libnode.NodeTypeKey, Value: cre.WorkerNode}, libnode.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		// Merge global defaults with DON-specific overrides
		var donOverrides map[string]map[string]any
		if donIdx < len(nodeSetInput) && nodeSetInput[donIdx] != nil {
			donOverrides = nodeSetInput[donIdx].CapabilityOverrides
		}

		mergedConfig := cre.ResolveCapabilityConfigForDON(string(flag), webAPITargetConfig.Config, donOverrides)

		// Apply runtime values only for keys not specified by user
		templateData := don.ApplyRuntimeValues(mergedConfig, map[string]any{})

		// If no custom config provided, use default config
		var configStr string
		if len(templateData) == 0 {
			configStr = jobs.EmptyStdCapConfig
		} else {
			// Parse and execute template with custom config
			tmpl, err := template.New("webAPITargetConfig").Parse(webAPITargetConfigTemplate)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse %s config template", flag)
			}

			var configBuffer bytes.Buffer
			if err := tmpl.Execute(&configBuffer, templateData); err != nil {
				return nil, errors.Wrapf(err, "failed to execute %s config template", flag)
			}
			configStr = configBuffer.String()

			if err := don.ValidateTemplateSubstitution(configStr, flag); err != nil {
				return nil, errors.Wrapf(err, "%s template validation failed", flag)
			}
		}

		for _, workerNode := range workflowNodeSet {
			nodeID, nodeIDErr := libnode.FindLabelValue(workerNode, libnode.NodeIDKey)
			if nodeIDErr != nil {
				return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
			}

			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.WorkerStandardCapability(nodeID, cre.WebAPITargetCapability, "__builtin_web-api-target", configStr, ""))
		}
	}

	return donToJobSpecs, nil
}
