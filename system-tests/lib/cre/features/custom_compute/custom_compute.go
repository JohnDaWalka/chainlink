package customcompute

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	coregateway "github.com/smartcontractkit/chainlink/v2/core/services/gateway"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	credon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/gateway"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
)

const flag = cre.CustomComputeCapability

type CustomCompute struct{}

func (o *CustomCompute) Flag() cre.CapabilityFlag {
	return flag
}

func (o *CustomCompute) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	// use registry chain, because that is the chain we used when generating gateway connector part of node config (check below)
	registryChainID, chErr := chainselectors.ChainIdFromSelector(creEnv.RegistryChainSelector)
	if chErr != nil {
		return nil, errors.Wrapf(chErr, "failed to get chain ID from selector %d", creEnv.RegistryChainSelector)
	}

	// add 'web-api' handler to gateway config (future jobspec)
	// add gateway connector to to node TOML config, so that node can route http requests to the gateway
	handlerConfig, confErr := gateway.HandlerConfig(coregateway.WebAPICapabilitiesType)
	if confErr != nil {
		return nil, errors.Wrapf(confErr, "failed to get %s handler config for don %s", coregateway.WebAPICapabilitiesType, don.Name)
	}
	hErr := gateway.AddHandlers(*don, registryChainID, topology.GatewayJobConfigs, []config.Handler{handlerConfig})
	if hErr != nil {
		return nil, errors.Wrapf(hErr, "failed to add gateway handlers to gateway config (jobspec) for don %s ", don.Name)
	}

	cErr := gateway.AddConnectors(don, registryChainID, *topology.GatewayConnectors)
	if cErr != nil {
		return nil, errors.Wrapf(cErr, "failed to add gateway connectors to node's TOML config in for don %s", don.Name)
	}

	capabilities := []keystone_changeset.DONCapabilityWithConfig{{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName:   "custom-compute",
			Version:        "1.0.0",
			CapabilityType: 1, // ACTION
		},
		Config: &capabilitiespb.CapabilityConfig{},
	}}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
	}, nil
}

const configTemplate = `"""
NumWorkers = {{.NumWorkers}}
[rateLimiter]
globalRPS = {{.GlobalRPS}}
globalBurst = {{.GlobalBurst}}
perSenderRPS = {{.PerSenderRPS}}
perSenderBurst = {{.PerSenderBurst}}
"""`

func (o *CustomCompute) PostEnvStartup(
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
	var nodeSet cre.NodeSetWithCapabilityConfigs
	for _, ns := range dons.AsNodeSetWithChainCapabilities() {
		if ns.GetName() == don.Name {
			nodeSet = ns
			break
		}
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

	// Create job specs for each worker node
	for _, workerNode := range workerNodes {
		jobSpec := standardcapability.WorkerJobSpec(workerNode.JobDistributorDetails.NodeID, flag, "__builtin_custom-compute-action", configStr, "")
		jobSpec.Labels = []*ptypes.Label{{Key: cre.CapabilityLabelKey, Value: ptr.Ptr(flag)}}
		jobSpecs = append(jobSpecs, jobSpec)
	}

	// pass all dons, since some jobs might need to be created on multiple ones
	jobErr := jobs.Create(ctx, creEnv.CldfEnvironment.Offchain, dons, jobSpecs)
	if jobErr != nil {
		return fmt.Errorf("failed to create http action jobs for don %s: %w", don.Name, jobErr)
	}

	return nil
}
