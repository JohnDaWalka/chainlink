package mock

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	credon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

const flag = cre.MockCapability

type Mock struct{}

func (o *Mock) Flag() cre.CapabilityFlag {
	return flag
}

func (o *Mock) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	capabilities := []keystone_changeset.DONCapabilityWithConfig{{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName:   "mock",
			Version:        "1.0.0",
			CapabilityType: 0, // TRIGGER
		},
		Config: &capabilitiespb.CapabilityConfig{},
	}}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
	}, nil
}

const configTemplate = `"""
port={{.Port}}
{{- range .DefaultMocks }}
[[DefaultMocks]]
id = "{{ .Id }}"
description = "{{ .Description }}"
type = "{{ .Type }}"
{{- end }}
"""`

func (o *Mock) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	jobSpecs := cre.DonJobs{}

	capabilityConfig, ok := creEnv.CapabilityConfigs[flag]
	if !ok {
		return errors.Errorf("%s config not found in capabilities config. Make sure you have set it in the TOML config", flag)
	}

	command, cErr := standardcapability.GetCommand(capabilityConfig.BinaryPath, creEnv.Provider)
	if cErr != nil {
		return errors.Wrap(cErr, "failed to get command for cron capability")
	}

	var nodeSet cre.NodeSetWithCapabilityConfigs
	for _, ns := range dons.AsNodeSetWithChainCapabilities() {
		if ns.GetName() == don.Name {
			nodeSet = ns
			break
		}
	}
	if nodeSet == nil {
		return fmt.Errorf("could not find node set for Don named '%s'", don.Name)
	}

	templateData := envconfig.ResolveCapabilityConfigForDON(flag, capabilityConfig.Config, nodeSet.GetCapabilityConfigOverrides())
	tmpl, tmplErr := template.New(flag + "-config").Parse(configTemplate)
	if tmplErr != nil {
		return errors.Wrapf(tmplErr, "failed to parse %s config template", flag)
	}

	var configBuffer bytes.Buffer
	if err := tmpl.Execute(&configBuffer, templateData); err != nil {
		return errors.Wrapf(err, "failed to execute %s config template", flag)
	}
	configStr := configBuffer.String()

	if err := credon.ValidateTemplateSubstitution(configStr, flag); err != nil {
		return errors.Wrapf(err, "%s template validation failed", flag)
	}

	workerNodes, wErr := don.Workers()
	if wErr != nil {
		return errors.Wrap(wErr, "failed to find worker nodes")
	}

	for _, workerNode := range workerNodes {
		jobSpec := standardcapability.WorkerJobSpec(workerNode.JobDistributorDetails.NodeID, flag, command, configStr, "")
		jobSpec.Labels = []*ptypes.Label{{Key: cre.CapabilityLabelKey, Value: ptr.Ptr(flag)}}
		jobSpecs = append(jobSpecs, jobSpec)
	}

	// pass all dons, since some jobs might need to be created on multiple dons
	jobErr := jobs.Create(ctx, creEnv.CldfEnvironment.Offchain, dons, jobSpecs)
	if jobErr != nil {
		return fmt.Errorf("failed to create http action jobs for don %s: %w", don.Name, jobErr)
	}

	return nil
}
