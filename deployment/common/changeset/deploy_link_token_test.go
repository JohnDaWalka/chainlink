package changeset_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestDeployLinkToken(t *testing.T) {
	t.Parallel()
	changeset.DeployLinkTokenTest(t, memory.MemoryEnvironmentConfig{
		Chains:    1,
		SolChains: 0,
		ZkChains:  0,
	})
}

func TestDeployLinkTokenZk(t *testing.T) {
	t.Parallel()
	changeset.DeployLinkTokenTest(t, memory.MemoryEnvironmentConfig{
		Chains:    0,
		SolChains: 0,
		ZkChains:  1,
	})
}
