package changeset_test

import (
	"testing"
	"time"

	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestMakeTxOptsAndHandlerForContract_UseMCMS(t *testing.T) {
	e, selectorA, _, tokens, _ := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	chain := e.Chains[selectorA]
	tokenAddress := tokens[selectorA].Address
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)
	chainState := state.Chains[selectorA]

	opts, handle, err := changeset.MakeTxOptsAndHandlerForContract(
		chainState.TokenAdminRegistry.Address(),
		chain,
		&changeset.MCMSConfig{
			MinDelay: 0 * time.Second,
		},
	)
	require.NoError(t, err)
	require.Equal(t, utils.ZeroAddress, opts.From)

	tx, err := chainState.TokenAdminRegistry.ProposeAdministrator(opts, tokenAddress, chain.DeployerKey.From)
	require.NoError(t, err)
	op, err := handle(tx)
	require.NoError(t, err)
	require.NotNil(t, op)
}

func TestMakeTxOptsAndHandlerForContract_UseDeployer(t *testing.T) {
	e, selectorA, _, tokens, _ := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), false)

	chain := e.Chains[selectorA]
	tokenAddress := tokens[selectorA].Address
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)
	chainState := state.Chains[selectorA]

	opts, handle, err := changeset.MakeTxOptsAndHandlerForContract(
		chainState.TokenAdminRegistry.Address(),
		chain,
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, chain.DeployerKey.From, opts.From)

	tx, err := chainState.TokenAdminRegistry.ProposeAdministrator(opts, tokenAddress, chain.DeployerKey.From)
	require.NoError(t, err)
	op, err := handle(tx)
	require.NoError(t, err)
	require.Nil(t, op)

	config, err := chainState.TokenAdminRegistry.GetTokenConfig(nil, tokenAddress)
	require.Equal(t, chain.DeployerKey.From, config.PendingAdministrator)
}
