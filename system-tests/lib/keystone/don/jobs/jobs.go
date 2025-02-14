package jobs

import (
	"context"
	"sync"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	keystoneflags "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/flags"
	types "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/types"
)

var SupportedJobs = []types.JobDescription{
	{Flag: types.OCR3Capability, NodeType: types.BootstrapNode},
	{Flag: types.WorkflowDON, NodeType: types.BootstrapNode},
	{Flag: types.CustomComputeCapability, NodeType: types.BootstrapNode},
	{Flag: types.CronCapability, NodeType: types.WorkerNode},
	{Flag: types.CustomComputeCapability, NodeType: types.WorkerNode},
	{Flag: types.OCR3Capability, NodeType: types.WorkerNode},
}

func Create(t *testing.T, offChainClient deployment.OffchainClient, don *devenv.DON, flags []string, jobSpecs types.DonJobs) {
	errCh := make(chan error, calculateJobCount(jobSpecs))

	var wg sync.WaitGroup

	for _, jobDesc := range SupportedJobs {
		if keystoneflags.HasFlag(flags, jobDesc.Flag) {
			if jobReqs, ok := jobSpecs[jobDesc]; ok {
				for _, jobReq := range jobReqs {
					wg.Add(1)
					go func(jobReq *jobv1.ProposeJobRequest) {
						defer wg.Done()
						_, err := offChainClient.ProposeJob(context.Background(), jobReq)
						if err != nil {
							errCh <- errors.Wrapf(err, "failed to propose job for node %s", jobReq.NodeId)
						}
					}(jobReq)
				}
			}
		}
	}

	wg.Wait()
	close(errCh)

	errFound := false
	for err := range errCh {
		errFound = true
		//nolint:testifylint // we want to assert here to catch all errors
		assert.NoError(t, err, "job creation/acception failed")
	}

	require.False(t, errFound, "failed to create at least one job for DON")
}

func calculateJobCount(jobSpecs types.DonJobs) int {
	count := 0
	for _, jobSpec := range jobSpecs {
		count += len(jobSpec)
	}

	return count
}
