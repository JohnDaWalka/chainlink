package ton

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// TODO: This is to test the implementation of Ton chains in memory environment
// To be deleted after changesets tests are added
func TestTonMemoryEnv(t *testing.T) {
	lggr := logger.TestLogger(t)
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		TonChains: 1,
	})
	tonChainSelectors := env.AllChainSelectorsTon()
	require.Len(t, tonChainSelectors, 1)
	require.NotEqual(t, 0, env.TonChains[0].Selector)
}

// TODO: This is to test the implementation of Ton chains in memory environment
// To be deleted after changesets tests are added
func TestTonHelperMemoryEnv(t *testing.T) {
	depEvn, testEnv := testhelpers.NewMemoryEnvironment(
		t,
		testhelpers.WithTonChains(1),
		testhelpers.WithNoJobsAndContracts(), // currently not supporting jobs and contracts
	)
	tonChainSelectors := depEvn.Env.AllChainSelectorsTon()
	require.Len(t, tonChainSelectors, 1)
	tonChainSelectors2 := testEnv.DeployedEnvironment().Env.AllChainSelectorsTon()
	require.Len(t, tonChainSelectors2, 1)
}
