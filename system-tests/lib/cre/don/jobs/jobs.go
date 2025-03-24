package jobs

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	types "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

func Create(offChainClient deployment.OffchainClient, don *devenv.DON, flags []string, jobSpecs types.DonJobs) error {
	if len(jobSpecs) == 0 {
		return nil
	}

	errCh := make(chan error, calculateJobCount(jobSpecs))
	var wg sync.WaitGroup
	sem := make(chan struct{}, 5)

	for jobDesc, jobReqs := range jobSpecs {
		for _, jobReq := range jobReqs {
			wg.Add(1)
			sem <- struct{}{}
			go func(jobReq *jobv1.ProposeJobRequest) {
				defer wg.Done()
				defer func() { <-sem }()
				timeout := time.Second * 60
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()
				_, err := offChainClient.ProposeJob(ctx, jobReq)
				if err != nil {
					errCh <- errors.Wrapf(err, "failed to propose job %s for node %s", jobDesc.Flag, jobReq.NodeId)
				}
				err = ctx.Err()
				if err != nil {
					errCh <- errors.Wrapf(err, "timed out after %s proposing job %s for node %s", timeout.String(), jobDesc.Flag, jobReq.NodeId)
				}
			}(jobReq)
		}
	}

	wg.Wait()
	close(errCh)

	var finalErr error
	for err := range errCh {
		if finalErr == nil {
			finalErr = err
		} else {
			finalErr = errors.Wrap(finalErr, err.Error())
		}
	}

	if finalErr != nil {
		return errors.Wrap(finalErr, "failed to create at least one job for DON")
	}

	return nil
}

func calculateJobCount(jobSpecs types.DonJobs) int {
	count := 0
	for _, jobSpec := range jobSpecs {
		count += len(jobSpec)
	}

	return count
}
