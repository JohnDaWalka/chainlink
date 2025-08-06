package readcontract

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

const readContractConfigTemplate = `'{"chainId":{{.ChainID}},"network":"{{.NetworkFamily}}"}'`

// buildReadContractRuntimeFallbacks creates runtime-generated fallback values for any keys not specified in TOML
func buildReadContractRuntimeFallbacks(chainID int, networkFamily string) map[string]any {
	return map[string]any{
		"ChainID":       chainID,
		"NetworkFamily": networkFamily,
	}
}

var ReadContractJobSpecFactoryFn = func(chainID int, networkFamily, readContractBinaryPath string, tomlConfig map[string]any) cre.JobSpecFactoryFn {
	return func(input *cre.JobSpecFactoryInput) (cre.DonsToJobSpecs, error) {
		return GenerateJobSpecs(input.DonTopology, chainID, networkFamily, readContractBinaryPath, tomlConfig)
	}
}

var ReadContractJobName = func(chainID int) string {
	return fmt.Sprintf("read-contract-%d", chainID)
}

func GenerateJobSpecs(donTopology *cre.DonTopology, chainID int, networkFamily, readContractBinaryPath string, tomlConfig map[string]any) (cre.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(cre.DonsToJobSpecs)

	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		if !flags.HasFlag(donWithMetadata.Flags, cre.ReadContractCapability) {
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

			// Build user configuration from TOML (global config is required)
			userConfig, err := config.BuildFromTOML(tomlConfig, chainID)
			if err != nil {
				return nil, errors.Wrap(err, "failed to build config from TOML")
			}

			// Build runtime fallbacks for any missing values
			runtimeFallbacks := buildReadContractRuntimeFallbacks(chainID, networkFamily)

			// Apply runtime fallbacks only for keys not specified by user
			templateData := config.ApplyRuntimeFallbacks(userConfig, runtimeFallbacks)

			// Parse and execute template
			tmpl, err := template.New("readContractConfig").Parse(readContractConfigTemplate)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse read contract config template")
			}

			var configBuffer bytes.Buffer
			if err := tmpl.Execute(&configBuffer, templateData); err != nil {
				return nil, errors.Wrap(err, "failed to execute read contract config template")
			}
			configStr := configBuffer.String()

			jobSpec := libjobs.WorkerStandardCapability(nodeID, ReadContractJobName(chainID), readContractBinaryPath, configStr, "")

			if _, ok := donToJobSpecs[donWithMetadata.ID]; !ok {
				donToJobSpecs[donWithMetadata.ID] = make(cre.DonJobs, 0)
			}

			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobSpec)
		}
	}

	return donToJobSpecs, nil
}
