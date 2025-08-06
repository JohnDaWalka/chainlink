package mock

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	libjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/config"
	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

const mockConfigTemplate = `"""port={{.Port}}"""`

// buildMockRuntimeFallbacks creates runtime-generated fallback values for any keys not specified in TOML
func buildMockRuntimeFallbacks(port int) map[string]any {
	return map[string]any{
		"Port": port,
	}
}

var MockJobSpecFactoryFn = func(port int, tomlConfig map[string]any) cre.JobSpecFactoryFn {
	return func(input *cre.JobSpecFactoryInput) (cre.DonsToJobSpecs, error) {
		return GenerateJobSpecs(input.DonTopology, port, tomlConfig)
	}
}

var MockJobName = func(chainID int) string {
	return fmt.Sprintf("mock-%d", chainID)
}

func GenerateJobSpecs(donTopology *cre.DonTopology, port int, tomlConfig map[string]any) (cre.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(cre.DonsToJobSpecs)

	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		if !flags.HasFlag(donWithMetadata.Flags, cre.MockCapability) {
			continue
		}

		workflowNodeSet, err := libnode.FindManyWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: libnode.NodeTypeKey, Value: cre.WorkerNode}, libnode.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		for _, workerNode := range workflowNodeSet {
			nodeID, nodeIDErr := libnode.FindLabelValue(workerNode, libnode.NodeIDKey)
			if nodeIDErr != nil {
				return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
			}

			// Build user configuration from TOML (global config is optional)
			userConfig, err := config.BuildFromTOMLOptional(tomlConfig)
			if err != nil {
				return nil, errors.Wrap(err, "failed to build config from TOML")
			}

			// Build runtime fallbacks for any missing values
			runtimeFallbacks := buildMockRuntimeFallbacks(port)

			// Apply runtime fallbacks only for keys not specified by user
			templateData := config.ApplyRuntimeFallbacks(userConfig, runtimeFallbacks)

			// Parse and execute template
			tmpl, err := template.New("mockConfig").Parse(mockConfigTemplate)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse mock config template")
			}

			var configBuffer bytes.Buffer
			if err := tmpl.Execute(&configBuffer, templateData); err != nil {
				return nil, errors.Wrap(err, "failed to execute mock config template")
			}
			configStr := configBuffer.String()

			jobSpec := libjobs.WorkerStandardCapability(nodeID, "mock-cap", "mock", configStr, "")

			if _, ok := donToJobSpecs[donWithMetadata.ID]; !ok {
				donToJobSpecs[donWithMetadata.ID] = make(cre.DonJobs, 0)
			}

			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobSpec)
		}
	}

	return donToJobSpecs, nil
}
