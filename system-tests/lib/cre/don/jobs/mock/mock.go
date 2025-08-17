package mock

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crenode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

var MockJobSpecFactoryFn = func(mockBinaryPath string) cre.JobSpecFactoryFn {
	return func(input *cre.JobSpecFactoryInput) (cre.DonsToJobSpecs, error) {
		return GenerateJobSpecs(
			input.DonTopology,
			mockBinaryPath,
		)
	}
}

func GenerateJobSpecs(donTopology *cre.DonTopology, mockBinaryPath string) (cre.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(cre.DonsToJobSpecs)

	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		if !flags.HasFlag(donWithMetadata.Flags, cre.MockCapability) {
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

			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], MockCapabilitiesJob(nodeID, "mock", []*MockCapabilities{
				{
					Name:        "mock",
					Version:     "1.0.0",
					Type:        "trigger",
					Description: "mock",
				},
			}))
		}
	}

	return donToJobSpecs, nil
}

type MockCapabilities struct {
	Name        string `toml:"name"`
	Version     string `toml:"version"`
	Type        string `toml:"type"`
	Description string `toml:"description"`
}

func MockCapabilitiesJob(nodeID, binaryPath string, mocks []*MockCapabilities) *jobv1.ProposeJobRequest {
	jobTemplate := `type = "standardcapabilities"
			schemaVersion = 1
			externalJobID = "{{ .JobID }}"
			name = "mock-capability"
			forwardingAllowed = false
			command = "{{ .BinaryPath }}"
			config = """
				port=7777
		{{ range $index, $m := .Mocks }}
 		  [[DefaultMocks]]
				id="{{ $m.ID }}"
				description="{{ $m.Description }}"
				type="{{ $m.Type }}"
 		{{- end }}
			"""`
	tmpl, err := template.New("mock-job").Parse(jobTemplate)

	if err != nil {
		panic(err)
	}
	mockJobsData := make([]map[string]string, 0)
	for _, m := range mocks {
		mockJobsData = append(mockJobsData, map[string]string{
			"ID":          m.Name + "@" + m.Version,
			"Description": m.Description,
			"Type":        m.Type,
		})
	}

	jobUUID := uuid.NewString()
	var renderedTemplate bytes.Buffer
	err = tmpl.Execute(&renderedTemplate, map[string]interface{}{
		"JobID":      jobUUID,
		"ShortID":    jobUUID[0:8],
		"BinaryPath": binaryPath,
		"Mocks":      mockJobsData,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("template", renderedTemplate.String())
	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec:   renderedTemplate.String(),
	}
}

