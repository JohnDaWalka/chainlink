package changeset_test

import (
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonState "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/wsrpc/logger"
)

func TestDeployLinkToken(t *testing.T) {
	t.Parallel()
	changeset.DeployLinkTokenTest(t, 0)
}

func TestZKChain(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains:   0,
		ZKChains: 1,
	})

	chain := chainsel.TEST_90000051
	config := []uint64{chain.Selector}
	e, err := changeset.ApplyDeployLinkToken(t, e, config)
	require.NoError(t, err)

	addrs, err := e.ExistingAddresses.AddressesForChain(chain.Selector)
	require.NoError(t, err)
	state, err := commonState.MaybeLoadLinkTokenChainState(e.Chains[chain.Selector], addrs)
	require.NoError(t, err)
	view, err := state.GenerateLinkView()
	require.NoError(t, err)
	require.Equal(t, uint8(18), view.Decimals)
}
