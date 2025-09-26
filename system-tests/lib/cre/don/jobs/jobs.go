package jobs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.uber.org/ratelimit"
	"golang.org/x/sync/errgroup"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/topology"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

type (
	DonJobs        = []*jobv1.ProposeJobRequest
	DonsToJobSpecs = map[uint64]DonJobs
	JobSpecFn      = func(input *JobSpecInput) (DonsToJobSpecs, error)
)

type JobSpecInput struct {
	CldEnvironment            *cldf.Environment
	BlockchainOutput          *blockchain.Output
	DonTopology               *topology.DonTopology
	InfraInput                infra.Provider
	CapabilityConfigs         map[string]cre.CapabilityConfig
	Capabilities              []types.InstallableCapability
	CapabilitiesAwareNodeSets []*cre.CapabilitiesAwareNodeSet
}

type CreateJobsInput struct {
	CldEnv        *cldf.Environment
	DonTopology   *topology.DonTopology
	DonToJobSpecs DonsToJobSpecs
}

func (c *CreateJobsInput) Validate() error {
	if c.CldEnv == nil {
		return errors.New("chainlink deployment env not set")
	}
	if c.DonTopology == nil {
		return errors.New("don topology not set")
	}
	if len(c.DonTopology.Dons.List()) == 0 {
		return errors.New("topology dons not set")
	}
	if len(c.DonToJobSpecs) == 0 {
		return errors.New("don to job specs not set")
	}

	return nil
}

func CreateJobs(ctx context.Context, testLogger zerolog.Logger, input CreateJobsInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "input validation failed")
	}

	for _, donMetadata := range input.DonTopology.ToDonMetadata() {
		if jobSpecs, ok := input.DonToJobSpecs[donMetadata.ID]; ok {
			createErr := Create(ctx, input.CldEnv.Offchain, jobSpecs)
			if createErr != nil {
				return errors.Wrapf(createErr, "failed to create jobs for DON %d", donMetadata.ID)
			}
		} else {
			testLogger.Warn().Msgf("No job specs found for DON %d", donMetadata.ID)
		}
	}

	return nil
}

func Create(ctx context.Context, offChainClient cldf_offchain.Client, jobSpecs cre.DonJobs) error {
	if len(jobSpecs) == 0 {
		return nil
	}

	eg := &errgroup.Group{}
	jobRateLimit := ratelimit.New(5)

	for _, jobReq := range jobSpecs {
		eg.Go(func() error {
			jobRateLimit.Take()
			timeout := time.Second * 60
			ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			_, err := offChainClient.ProposeJob(ctxWithTimeout, jobReq)
			if err != nil {
				// Workflow specs get auto approved
				// TODO: Narrow down scope by checking type == workflow
				if strings.Contains(err.Error(), "cannot approve an approved spec") {
					return nil
				}
				fmt.Println("Failed jobspec proposal:")
				fmt.Println(jobReq)
				return errors.Wrapf(err, "failed to propose job for node %s", jobReq.NodeId)
			}
			if ctx.Err() != nil {
				return errors.Wrapf(err, "timed out after %s proposing job for node %s", timeout.String(), jobReq.NodeId)
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return errors.Wrap(err, "failed to create at least one job for DON")
	}

	return nil
}
