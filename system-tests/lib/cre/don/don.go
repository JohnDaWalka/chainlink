package don

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
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
