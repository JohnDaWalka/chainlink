package don

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

func CreateJobs(ctx context.Context, testLogger zerolog.Logger, input cre.CreateJobsInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "input validation failed")
	}

	for _, don := range input.DonTopology.Dons.List() {
		if jobSpecs, ok := input.DonToJobSpecs[don.ID]; ok {
			createErr := jobs.Create(ctx, input.CldEnv.Offchain, *input.DonTopology, jobSpecs)
			if createErr != nil {
				return errors.Wrapf(createErr, "failed to create jobs for DON %d", don.ID)
			}
		} else {
			testLogger.Warn().Msgf("No job specs found for DON %d", don.ID)
		}
	}

	return nil
}

func AnyDonHasCapability(dons []*cre.DON, capability cre.CapabilityFlag) bool {
	for _, don := range dons {
		if flags.HasFlagForAnyChain(don.Flags, capability) {
			return true
		}
	}

	return false
}

func NodeNeedsAnyGateway(nodeFlags []cre.CapabilityFlag) bool {
	return flags.HasFlag(nodeFlags, cre.CustomComputeCapability) ||
		flags.HasFlag(nodeFlags, cre.WebAPITriggerCapability) ||
		flags.HasFlag(nodeFlags, cre.WebAPITargetCapability) ||
		flags.HasFlag(nodeFlags, cre.VaultCapability) ||
		flags.HasFlag(nodeFlags, cre.HTTPActionCapability) ||
		flags.HasFlag(nodeFlags, cre.HTTPTriggerCapability)
}

func NodeNeedsWebAPIGateway(nodeFlags []cre.CapabilityFlag) bool {
	return flags.HasFlag(nodeFlags, cre.CustomComputeCapability) ||
		flags.HasFlag(nodeFlags, cre.WebAPITriggerCapability) ||
		flags.HasFlag(nodeFlags, cre.WebAPITargetCapability)
}
