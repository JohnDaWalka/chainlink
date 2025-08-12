package logevent

import (
	"bytes"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"

	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

const flag = cre.LogTriggerCapability
const logEventTriggerConfigTemplate = `'{"chainId":"{{.ChainID}}","network":"{{.NetworkFamily}}","lookbackBlocks":{{.LookbackBlocks}},"pollPeriod":{{.PollPeriod}}}'`

// Log event trigger capability with per-chain configuration support

var JobSpecFn = func(input *cre.JobSpecInput) (cre.DonsToJobSpecs, error) {
	return generateJobSpecs(
		input.DonTopology,
		*input.InfraInput,
		input.CapabilityConfigs,
		input.CapabilitiesAwareNodeSets,
	)
}

var LogEventTriggerJobName = func(chainID string) string {
	return "log-event-trigger-" + chainID
}

func generateJobSpecs(donTopology *cre.DonTopology, infraInput infra.Input, capabilitiesConfig cre.CapabilityConfigs, nodeSetInput []*cre.CapabilitiesAwareNodeSet) (cre.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(cre.DonsToJobSpecs)

	for donIdx, donWithMetadata := range donTopology.DonsWithMetadata {
		// Log event trigger capability is enabled strictly per-chain via ChainCapabilities
		if donIdx >= len(nodeSetInput) || nodeSetInput[donIdx] == nil || nodeSetInput[donIdx].ChainCapabilities == nil {
			continue
		}
		if cc, ok := nodeSetInput[donIdx].ChainCapabilities[flag]; !ok || cc == nil || len(cc.EnabledChains) == 0 {
			continue
		}

		logEventConfig, ok := capabilitiesConfig[flag]
		if !ok {
			return nil, errors.Errorf("%s config not found in capabilities config", flag)
		}

		containerPath, pathErr := crecapabilities.DefaultContainerDirectory(infraInput.Type)
		if pathErr != nil {
			return nil, errors.Wrapf(pathErr, "failed to get default container directory for infra type %s", infraInput.Type)
		}

		logEventTriggerBinaryPath := filepath.Join(containerPath, filepath.Base(logEventConfig.BinaryPath))

		workflowNodeSet, err := libnode.FindManyWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: libnode.NodeTypeKey, Value: cre.WorkerNode}, libnode.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		for _, chainIDUint64 := range nodeSetInput[donIdx].ChainCapabilities[flag].EnabledChains {
			chainID := int(chainIDUint64)
			chainIDStr := strconv.Itoa(chainID)

			// Build user configuration from defaults + chain overrides
			enabled, mergedConfig, rErr := cre.ResolveCapabilityForChain(flag, nodeSetInput[donIdx].ChainCapabilities, logEventConfig.Config, chainIDUint64)
			if rErr != nil {
				return nil, errors.Wrap(rErr, "failed to resolve capability config for chain")
			}
			if !enabled {
				// should not happen because we derived enabledChains from the same source, but guard anyway
				continue
			}

			// Build runtime values for any missing values
			runtimeFallbacks := map[string]any{
				"ChainID":       strconv.Itoa(chainID), // string for logevent template
				"NetworkFamily": "evm",
			}

			// Apply runtime values only for keys not specified by user
			templateData := don.ApplyRuntimeValues(mergedConfig, runtimeFallbacks)

			// Parse and execute template
			tmpl, err := template.New("logEventTriggerConfig").Parse(logEventTriggerConfigTemplate)
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
				nodeID, nodeIDErr := libnode.FindLabelValue(workerNode, libnode.NodeIDKey)
				if nodeIDErr != nil {
					return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
				}

				jobSpec := jobs.WorkerStandardCapability(nodeID, LogEventTriggerJobName(chainIDStr), logEventTriggerBinaryPath, configStr, "")

				if _, ok := donToJobSpecs[donWithMetadata.ID]; !ok {
					donToJobSpecs[donWithMetadata.ID] = make(cre.DonJobs, 0)
				}

				donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobSpec)
			}
		}
	}

	return donToJobSpecs, nil
}
