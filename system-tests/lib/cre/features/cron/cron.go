package cron

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	ptypes "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	standardcapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
)

const flag = cre.CronCapability

type Cron struct{}

func (c *Cron) Flag() cre.CapabilityFlag {
	return flag
}

func (c *Cron) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	capabilities := []keystone_changeset.DONCapabilityWithConfig{{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName:   "cron-trigger",
			Version:        "1.0.0",
			CapabilityType: 0, // TRIGGER
		},
		Config: &capabilitiespb.CapabilityConfig{},
	}}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
	}, nil
}

func (c *Cron) PostEnvStartup(
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

	workerNodes, wErr := don.Workers()
	if wErr != nil {
		return errors.Wrap(wErr, "failed to find worker nodes")
	}

	for _, workerNode := range workerNodes {
		jobSpec := standardcapability.WorkerJobSpec(workerNode.JobDistributorDetails.NodeID, flag, command, `""`, "")
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
