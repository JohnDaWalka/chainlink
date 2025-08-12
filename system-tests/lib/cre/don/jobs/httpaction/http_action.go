package httpaction

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	creregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	crenode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

const flag = cre.HTTPActionCapability

const httpActionConfigTemplate = `"""
{
	"proxyMode": "{{.ProxyMode}}",
	"incomingRateLimiter": {
		"globalBurst": {{.IncomingGlobalBurst}},
		"globalRPS": {{.IncomingGlobalRPS}},
		"perSenderBurst": {{.IncomingPerSenderBurst}},
		"perSenderRPS": {{.IncomingPerSenderRPS}}
	},
	"outgoingRateLimiter": {
		"globalBurst": {{.OutgoingGlobalBurst}},
		"globalRPS": {{.OutgoingGlobalRPS}},
		"perSenderBurst": {{.OutgoingPerSenderBurst}},
		"perSenderRPS": {{.OutgoingPerSenderRPS}}
	}
}
"""`

var JobSpecFn = func(input *cre.JobSpecFactoryInput) (cre.DonsToJobSpecs, error) {
	return generateJobSpecs(
		input.DonTopology,
		input.InfraInput,
		input.AdditionalCapabilityConfigs,
		input.CapabilitiesAwareNodeSets,
	)
}

func generateJobSpecs(donTopology *cre.DonTopology, infraInput *infra.Input, capabilitiesConfig cre.AdditionalCapabilitiesConfigs, nodeSetInput []*cre.CapabilitiesAwareNodeSet) (cre.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(cre.DonsToJobSpecs)

	for donIdx, donWithMetadata := range donTopology.DonsWithMetadata {
		if !flags.HasFlag(donWithMetadata.Flags, cre.HTTPActionCapability) {
			continue
		}

		httpActionConfig, ok := capabilitiesConfig[flag]
		if !ok {
			return nil, errors.Errorf("%s config not found in capabilities config", flag)
		}

		containerPath, pathErr := creregistry.DefaultContainerDirectory(infraInput.Type)
		if pathErr != nil {
			return nil, errors.Wrapf(pathErr, "failed to get default container directory for infra type %s", infraInput.Type)
		}

		httpActionBinaryPath := filepath.Join(containerPath, filepath.Base(httpActionConfig.BinaryPath))

		workflowNodeSet, err := crenode.FindManyWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: crenode.NodeTypeKey, Value: cre.WorkerNode}, crenode.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		var donOverrides map[string]map[string]any
		if donIdx < len(nodeSetInput) && nodeSetInput[donIdx] != nil {
			donOverrides = nodeSetInput[donIdx].CapabilityOverrides
		}

		mergedConfig := cre.ResolveCapabilityConfigForDON(flag, httpActionConfig.Config, donOverrides)
		templateData := don.ApplyRuntimeValues(mergedConfig, map[string]any{})

		tmpl, err := template.New("httpActionConfig").Parse(httpActionConfigTemplate)
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

			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.WorkerStandardCapability(nodeID, cre.HTTPActionCapability, httpActionBinaryPath, configStr, ""))
		}
	}

	return donToJobSpecs, nil
}
