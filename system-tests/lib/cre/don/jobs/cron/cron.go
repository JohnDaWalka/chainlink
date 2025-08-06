package cron

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/config"
	crenode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

// Cron capability uses empty config by default, but can be overridden via TOML
const cronConfigTemplate = `""` // Empty config by default

// buildCronRuntimeFallbacks creates empty runtime fallbacks (cron needs no defaults)
func buildCronRuntimeFallbacks() map[string]any {
	return map[string]any{} // Empty by default
}

var CronJobSpecFactoryFn = func(cronBinaryPath string, config map[string]any) cre.JobSpecFactoryFn {
	return func(input *cre.JobSpecFactoryInput) (cre.DonsToJobSpecs, error) {
		return GenerateJobSpecs(
			input.DonTopology,
			cronBinaryPath,
			config,
		)
	}
}

func GenerateJobSpecs(donTopology *cre.DonTopology, cronBinaryPath string, tomlConfig map[string]any) (cre.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(cre.DonsToJobSpecs)

	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		if !flags.HasFlag(donWithMetadata.Flags, cre.CronCapability) {
			continue
		}
		workflowNodeSet, err := crenode.FindManyWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: crenode.NodeTypeKey, Value: cre.WorkerNode}, crenode.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		for _, workerNode := range workflowNodeSet {
			nodeID, nodeIDErr := crenode.FindLabelValue(workerNode, crenode.NodeIDKey)
			if nodeIDErr != nil {
				return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
			}

			// Build user configuration from TOML (optional for cron)
			userConfig, err := config.BuildFromTOMLOptional(tomlConfig)
			if err != nil {
				return nil, errors.Wrap(err, "failed to build config from TOML")
			}

			// Build runtime fallbacks for any missing values (empty for cron)
			runtimeFallbacks := buildCronRuntimeFallbacks()

			// Apply runtime fallbacks only for keys not specified by user
			templateData := config.ApplyRuntimeFallbacks(userConfig, runtimeFallbacks)

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

			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.WorkerStandardCapability(nodeID, cre.CronCapability, cronBinaryPath, configStr, ""))
		}
	}

	return donToJobSpecs, nil
}
