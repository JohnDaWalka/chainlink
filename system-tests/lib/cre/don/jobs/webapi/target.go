package webapi

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
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

var WebAPITargetJobSpecFactoryFn = func(input *cre.JobSpecFactoryInput) (cre.DonsToJobSpecs, error) {
	return generateJobSpecs(
		input.DonTopology,
		input.InfraInput,
		input.AdditionalCapabilities,
	)
}

func generateJobSpecs(donTopology *cre.DonTopology, _ *infra.Input, capabilitiesConfig cre.AdditionalCapabilitiesConfigs) (cre.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(cre.DonsToJobSpecs)

	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		if !flags.HasFlag(donWithMetadata.Flags, cre.WebAPITargetCapability) {
			continue
		}

		webAPITargetConfig, ok := capabilitiesConfig[flag]
		if !ok {
			return nil, errors.New("web-api-target config not found in capabilities config")
		}
		workflowNodeSet, err := libnode.FindManyWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: libnode.NodeTypeKey, Value: cre.WorkerNode}, libnode.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		// Apply runtime fallbacks only for keys not specified by user
		templateData := jobs.ApplyRuntimeValues(webAPITargetConfig.Config, map[string]any{})

		// If no custom config provided, use default config
		var configStr string
		if len(templateData) == 0 {
			configStr = jobs.EmptyStdCapConfig
		} else {
			// Parse and execute template with custom config
			tmpl, err := template.New("webAPITargetConfig").Parse(webAPITargetConfigTemplate)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse web-api-target config template")
			}

			var configBuffer bytes.Buffer
			if err := tmpl.Execute(&configBuffer, templateData); err != nil {
				return nil, errors.Wrap(err, "failed to execute web-api-target config template")
			}
			configStr = configBuffer.String()
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
