package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestDeployLinkToken(t *testing.T) {
	t.Parallel()
	env := memory.NewMemoryEnvironment(t, logger.Test(t), zapcore.DebugLevel, memory.MemoryEnvironmentConfig{
		Nodes:  1,
		Chains: 2,
	})
	chain1 := env.AllChainSelectors()[0]
	chain2 := env.AllChainSelectors()[1]

	resp, err := changeset.DeployLinkToken(env, changeset.DeployLinkTokenConfig{
		LinkTokenByChain: map[uint64]deployment.ContractType{
			chain1: types.StaticLinkToken,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	// LinkToken should be deployed on chain 1
	addrs, err := resp.AddressBook.AddressesForChain(chain1)
	require.NoError(t, err)
	require.Len(t, addrs, 1)

	// nothing on chain 2
	oaddrs, _ := resp.AddressBook.AddressesForChain(chain2)
	assert.Len(t, oaddrs, 0)
}
