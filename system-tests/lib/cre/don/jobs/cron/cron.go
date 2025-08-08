package cron

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	crenode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

const flag = cre.CronCapability

// Cron capability uses empty config by default, but can be overridden via TOML
const cronConfigTemplate = `""` // Empty config by default

var CronJobSpecFactoryFn = func(input *cre.JobSpecFactoryInput) (cre.DonsToJobSpecs, error) {
	return generateJobSpecs(
		input.DonTopology,
		input.InfraInput,
		input.AdditionalCapabilities,
	)
}

func generateJobSpecs(donTopology *cre.DonTopology, infraInput *infra.Input, capabilitiesConfig cre.AdditionalCapabilitiesConfigs) (cre.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(cre.DonsToJobSpecs)

	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		if !flags.HasFlag(donWithMetadata.Flags, cre.CronCapability) {
			continue
		}

		cronConfig, ok := capabilitiesConfig[flag]
		if !ok {
			return nil, errors.New("cron config not found in capabilities config")
		}

		containerPath, pathErr := crecapabilities.DefaultContainerDirectory(infraInput.Type)
		if pathErr != nil {
			return nil, errors.Wrapf(pathErr, "failed to get default container directory for infra type %s", infraInput.Type)
		}

		cronBinaryPath := filepath.Join(containerPath, filepath.Base(cronConfig.BinaryPath))

		workflowNodeSet, err := crenode.FindManyWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: crenode.NodeTypeKey, Value: cre.WorkerNode}, crenode.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		// Build user configuration from TOML (optional for cron)
		userConfig, err := jobs.BuildGlobalConfigFromTOML(cronConfig.Config)
		if err != nil {
			return nil, errors.Wrap(err, "failed to build config from TOML")
		}

		// Apply runtime fallbacks only for keys not specified by user
		templateData := jobs.ApplyRuntimeFallbacks(userConfig, map[string]any{})

		// If no custom config provided, use empty config (jobs.EmptyStdCapConfig)
		var configStr string
		if len(templateData) == 0 {
			configStr = jobs.EmptyStdCapConfig
		} else {
			// Parse and execute template with custom config
			tmpl, err := template.New("cronConfig").Parse(cronConfigTemplate)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse cron config template")
			}

			var configBuffer bytes.Buffer
			if err := tmpl.Execute(&configBuffer, templateData); err != nil {
				return nil, errors.Wrap(err, "failed to execute cron config template")
			}
			configStr = configBuffer.String()
		}

		for _, workerNode := range workflowNodeSet {
			nodeID, nodeIDErr := crenode.FindLabelValue(workerNode, crenode.NodeIDKey)
			if nodeIDErr != nil {
				return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
			}

			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.WorkerStandardCapability(nodeID, cre.CronCapability, cronBinaryPath, configStr, ""))
		}
	}

	return donToJobSpecs, nil
}
