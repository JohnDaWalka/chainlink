package logevent

import (
	"bytes"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	libjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

const flag = cre.LogTriggerCapability
const logEventTriggerConfigTemplate = `'{"chainId":"{{.ChainID}}","network":"{{.NetworkFamily}}","lookbackBlocks":{{.LookbackBlocks}},"pollPeriod":{{.PollPeriod}}}'`

// Log event trigger is now fully configurable via TOML - no runtime fallbacks needed

var LogEventTriggerJobSpecFactoryFn = func(input *cre.JobSpecFactoryInput) (cre.DonsToJobSpecs, error) {
	return generateJobSpecs(
		input.DonTopology,
		*input.InfraInput,
		input.AdditionalCapabilities,
	)
}

var LogEventTriggerJobName = func(chainID string) string {
	return "log-event-trigger-" + chainID
}

func generateJobSpecs(donTopology *cre.DonTopology, infraInput infra.Input, capabilitiesConfig cre.AdditionalCapabilitiesConfigs) (cre.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(cre.DonsToJobSpecs)

	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		if !flags.HasFlag(donWithMetadata.Flags, cre.LogTriggerCapability) {
			continue
		}

		logEventConfig, ok := capabilitiesConfig[flag]
		if !ok {
			return nil, errors.New("log event trigger config not found in capabilities config")
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

		// Build user configuration from TOML (optional for cron)
		globalConfig, err := jobs.BuildGlobalConfigFromTOML(logEventConfig.Config)
		if err != nil {
			return nil, errors.Wrap(err, "failed to build config from TOML")
		}

		for _, chainIDStr := range logEventConfig.Chains {
			chainID, err := strconv.Atoi(chainIDStr)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to convert chain ID %s to int", chainIDStr)
			}

			// Extract chain ID from template data to pass to BuildFromTOML
			// We need to get it from the first chain in the chains list or from config
			templateData, err := jobs.BuildConfigFromTOML(globalConfig, logEventConfig.Config, chainID)
			if err != nil {
				return nil, errors.Wrap(err, "failed to build config from TOML")
			}

			// Parse and execute template
			tmpl, err := template.New("logEventTriggerConfig").Parse(logEventTriggerConfigTemplate)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse log event trigger config template")
			}

			var configBuffer bytes.Buffer
			if err := tmpl.Execute(&configBuffer, templateData); err != nil {
				return nil, errors.Wrap(err, "failed to execute log event trigger config template")
			}
			configStr := configBuffer.String()

			for _, workerNode := range workflowNodeSet {
				nodeID, nodeIDErr := libnode.FindLabelValue(workerNode, libnode.NodeIDKey)
				if nodeIDErr != nil {
					return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
				}

				jobSpec := libjobs.WorkerStandardCapability(nodeID, LogEventTriggerJobName(chainIDStr), logEventTriggerBinaryPath, configStr, "")

				if _, ok := donToJobSpecs[donWithMetadata.ID]; !ok {
					donToJobSpecs[donWithMetadata.ID] = make(cre.DonJobs, 0)
				}

				donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobSpec)
			}
		}
	}

	return donToJobSpecs, nil
}
