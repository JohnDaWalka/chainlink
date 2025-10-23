package readcontract

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

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

const flag = cre.ReadContractCapability

type ReadContract struct{}

func (o *ReadContract) Flag() cre.CapabilityFlag {
	return flag
}

func (o *ReadContract) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	capabilities := []keystone_changeset.DONCapabilityWithConfig{}
	for _, chainID := range don.CapabilitiesAwareNodeSet().GetChainCapabilityConfigs()[flag].EnabledChains {
		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   fmt.Sprintf("read-contract-evm-%d", chainID),
				Version:        "1.0.0",
				CapabilityType: 1, // ACTION
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
	}, nil
}

const configTemplate = `'{"chainId":{{.ChainID}},"network":"{{.NetworkFamily}}"}'`

func (o *ReadContract) PostEnvStartup(
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

	chainCapConfig, ok := nodeSet.GetChainCapabilityConfigs()[flag]
	if !ok || chainCapConfig == nil {
		return fmt.Errorf("could not find chain capability config for '%s' in don '%s'", flag, don.Name)
	}

	for _, chainID := range chainCapConfig.EnabledChains {
		_, templateData, tErr := envconfig.ResolveCapabilityForChain(flag, nodeSet.GetChainCapabilityConfigs(), capabilityConfig.Config, chainID)
		if tErr != nil {
			return errors.Wrapf(tErr, "failed to resolve capability config for chain %d", chainID)
		}
		templateData["ChainID"] = chainID

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
			jobSpec := standardcapability.WorkerJobSpec(workerNode.JobDistributorDetails.NodeID, fmt.Sprintf("%s-%d", flag, chainID), command, configStr, "")
			jobSpec.Labels = []*ptypes.Label{{Key: cre.CapabilityLabelKey, Value: ptr.Ptr(flag)}}
			jobSpecs = append(jobSpecs, jobSpec)
		}
	}

	// pass all dons, since some jobs might need to be created on multiple ones
	jobErr := jobs.Create(ctx, creEnv.CldfEnvironment.Offchain, dons, jobSpecs)
	if jobErr != nil {
		return fmt.Errorf("failed to create http action jobs for don %s: %w", don.Name, jobErr)
	}

	return nil
}
