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

const triggerFlag = cre.WebAPITriggerCapability

// Web API trigger capability uses empty config by default, but can be overridden via TOML
const webAPITriggerConfigTemplate = `""` // Empty config by default

var WebAPITriggerJobSpecFactoryFn = func(input *cre.JobSpecFactoryInput) (cre.DonsToJobSpecs, error) {
	return generateTriggerJobSpecs(
		input.DonTopology,
		input.InfraInput,
		input.AdditionalCapabilities,
	)
}

func generateTriggerJobSpecs(donTopology *cre.DonTopology, _ *infra.Input, capabilitiesConfig cre.AdditionalCapabilitiesConfigs) (cre.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(cre.DonsToJobSpecs)

	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		if !flags.HasFlag(donWithMetadata.Flags, cre.WebAPITriggerCapability) {
			continue
		}

		webAPITriggerConfig, ok := capabilitiesConfig[triggerFlag]
		if !ok {
			return nil, errors.New("web-api-trigger config not found in capabilities config")
		}

		workflowNodeSet, err := libnode.FindManyWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: libnode.NodeTypeKey, Value: cre.WorkerNode}, libnode.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		// Apply runtime values only for keys not specified by user
		templateData := jobs.ApplyRuntimeValues(webAPITriggerConfig.Config, map[string]any{})

		// If no custom config provided, use empty config (jobs.EmptyStdCapConfig)
		var configStr string
		if len(templateData) == 0 {
			configStr = jobs.EmptyStdCapConfig
		} else {
			// Parse and execute template with custom config
			tmpl, err := template.New("webAPITriggerConfig").Parse(webAPITriggerConfigTemplate)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse web-api-trigger config template")
			}

			var configBuffer bytes.Buffer
			if err := tmpl.Execute(&configBuffer, templateData); err != nil {
				return nil, errors.Wrap(err, "failed to execute web-api-trigger config template")
			}
			configStr = configBuffer.String()
		}

		for _, workerNode := range workflowNodeSet {
			nodeID, nodeIDErr := libnode.FindLabelValue(workerNode, libnode.NodeIDKey)
			if nodeIDErr != nil {
				return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
			}

			if _, ok := donToJobSpecs[donWithMetadata.ID]; !ok {
				donToJobSpecs[donWithMetadata.ID] = make(cre.DonJobs, 0)
			}
			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.WorkerStandardCapability(nodeID, cre.WebAPITriggerCapability, "__builtin_web-api-trigger", configStr, ""))
		}
	}

	return donToJobSpecs, nil
}
