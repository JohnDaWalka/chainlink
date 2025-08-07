package mock

import (
	"bytes"
	"path/filepath"
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

const flag = cre.MockCapability
const mockConfigTemplate = `"""port={{.Port}}"""`

var MockJobSpecFactoryFn = func(input *cre.JobSpecFactoryInput) (cre.DonsToJobSpecs, error) {
	return generateJobSpecs(
		input.DonTopology,
		*input.InfraInput,
		input.AdditionalCapabilities,
	)
}

func generateJobSpecs(donTopology *cre.DonTopology, infraInput infra.Input, capabilitiesConfig cre.AdditionalCapabilitiesConfigs) (cre.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(cre.DonsToJobSpecs)

	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		if !flags.HasFlag(donWithMetadata.Flags, cre.MockCapability) {
			continue
		}

		mockConfig, ok := capabilitiesConfig[flag]
		if !ok {
			return nil, errors.New("mock config not found in capabilities config")
		}

		containerPath, pathErr := crecapabilities.DefaultContainerDirectory(infraInput.Type)
		if pathErr != nil {
			return nil, errors.Wrapf(pathErr, "failed to get default container directory for infra type %s", infraInput.Type)
		}

		mockBinaryPath := filepath.Join(containerPath, filepath.Base(mockConfig.BinaryPath))

		workflowNodeSet, err := libnode.FindManyWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: libnode.NodeTypeKey, Value: cre.WorkerNode}, libnode.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		// Build user configuration from TOML (optional for cron)
		globalConfig, err := jobs.BuildGlobalConfigFromTOML(mockConfig.Config)
		if err != nil {
			return nil, errors.Wrap(err, "failed to build config from TOML")
		}

		// Parse and execute template
		tmpl, err := template.New("mockConfig").Parse(mockConfigTemplate)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse mock config template")
		}

		var configBuffer bytes.Buffer
		if err := tmpl.Execute(&configBuffer, globalConfig); err != nil {
			return nil, errors.Wrap(err, "failed to execute mock config template")
		}
		configStr := configBuffer.String()

		for _, workerNode := range workflowNodeSet {
			nodeID, nodeIDErr := libnode.FindLabelValue(workerNode, libnode.NodeIDKey)
			if nodeIDErr != nil {
				return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
			}

			jobSpec := libjobs.WorkerStandardCapability(nodeID, "mock-cap", mockBinaryPath, configStr, "")

			if _, ok := donToJobSpecs[donWithMetadata.ID]; !ok {
				donToJobSpecs[donWithMetadata.ID] = make(cre.DonJobs, 0)
			}

			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobSpec)
		}
	}

	return donToJobSpecs, nil
}
