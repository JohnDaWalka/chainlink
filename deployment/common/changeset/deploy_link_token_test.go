package changeset_test

import (
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/wsrpc/logger"
)

func TestDeployLinkToken(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	solChains := 0
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains:    1,
		SolChains: solChains,
	})
	changeset.DeployLinkTokenTest(t, e, solChains)
}

func TestDeployLinkTokenZK(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains:    0,
		SolChains: 0,
		ZkChains:  1,
	})
	changeset.DeployLinkTokenTest(t, e, 0)
}
