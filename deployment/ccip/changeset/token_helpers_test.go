package changeset_test

import (
	"testing"

	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestGetAllTokenPoolsWithSymbol(t *testing.T) {
	t.Parallel()

	l := logger.TestLogger(t)

	e, selectorA, _, tokens, _ := testhelpers.SetupTwoChainEnvironmentWithTokens(t, l, true)
	acceptLiquidity := false

	tests := []struct {
		Msg   string
		Input changeset.DeployTokenPoolInput
	}{
		{
			Msg: "Add BurnMint",
			Input: changeset.DeployTokenPoolInput{
				TokenAddress:       tokens[selectorA].Address,
				Type:               changeset.BurnMintTokenPool,
				LocalTokenDecimals: testhelpers.LocalTokenDecimals,
				AllowList:          []common.Address{},
			},
		},
		{
			Msg: "Add BurnWithFromMint",
			Input: changeset.DeployTokenPoolInput{
				TokenAddress:       tokens[selectorA].Address,
				Type:               changeset.BurnWithFromMintTokenPool,
				LocalTokenDecimals: testhelpers.LocalTokenDecimals,
				AllowList:          []common.Address{},
			},
		},
		{
			Msg: "Add BurnFromMint",
			Input: changeset.DeployTokenPoolInput{
				TokenAddress:       tokens[selectorA].Address,
				Type:               changeset.BurnFromMintTokenPool,
				LocalTokenDecimals: testhelpers.LocalTokenDecimals,
				AllowList:          []common.Address{},
			},
		},
		{
			Msg: "Add LockRelease",
			Input: changeset.DeployTokenPoolInput{
				TokenAddress:       tokens[selectorA].Address,
				Type:               changeset.LockReleaseTokenPool,
				LocalTokenDecimals: testhelpers.LocalTokenDecimals,
				AllowList:          []common.Address{},
				AcceptLiquidity:    &acceptLiquidity,
			},
		},
	}

	for i, test := range tests {
		chain := e.Chains[selectorA]
		state, err := changeset.LoadOnchainState(e)
		require.NoError(t, err)
		chainState := state.Chains[selectorA]
		addressBook := deployment.NewMemoryAddressBook()

		if i == 0 {
			tokenPools, err := changeset.GetAllTokenPoolsWithSymbol(state.Chains[selectorA], chain.Client, testhelpers.TestTokenSymbol)
			require.NoError(t, err)
			require.Empty(t, tokenPools)
		}

		_, err = changeset.DeployTokenPool(l, chain, chainState, addressBook, test.Input)
		require.NoError(t, err)

		err = e.ExistingAddresses.Merge(addressBook)
		require.NoError(t, err)

		state, err = changeset.LoadOnchainState(e)
		require.NoError(t, err)
		chainState = state.Chains[selectorA]

		tokenPools, err := changeset.GetAllTokenPoolsWithSymbol(chainState, chain.Client, testhelpers.TestTokenSymbol)
		require.NoError(t, err)
		require.Len(t, tokenPools, i+1)
	}
}

func TestGetTokenPoolWithSymbolAndAddress(t *testing.T) {
	t.Parallel()

	l := logger.TestLogger(t)

	e, selectorA, _, tokens, _ := testhelpers.SetupTwoChainEnvironmentWithTokens(t, l, true)

	chain := e.Chains[selectorA]
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)
	chainState := state.Chains[selectorA]
	addressBook := deployment.NewMemoryAddressBook()

	_, err = changeset.DeployTokenPool(l, chain, chainState, addressBook, changeset.DeployTokenPoolInput{
		TokenAddress:       tokens[selectorA].Address,
		Type:               changeset.BurnMintTokenPool,
		LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		AllowList:          []common.Address{},
	})
	require.NoError(t, err)

	err = e.ExistingAddresses.Merge(addressBook)
	require.NoError(t, err)

	state, err = changeset.LoadOnchainState(e)
	require.NoError(t, err)
	chainState = state.Chains[selectorA]

	poolAddress := chainState.BurnMintTokenPools[testhelpers.TestTokenSymbol][0].Address()
	wrongPoolAddress := utils.RandomAddress()

	// Get correct pool
	tokenPool, err := changeset.GetTokenPoolWithSymbolAndAddress(chainState, chain, testhelpers.TestTokenSymbol, poolAddress)
	require.NoError(t, err)
	require.NotNil(t, tokenPool)

	// Get wrong pool
	tokenPool, err = changeset.GetTokenPoolWithSymbolAndAddress(chainState, chain, testhelpers.TestTokenSymbol, wrongPoolAddress)
	require.ErrorContains(t, err, fmt.Sprintf("no token pool found with symbol %s and address %s on chain %s", testhelpers.TestTokenSymbol, wrongPoolAddress, chain.String()))
	require.Nil(t, tokenPool)
}
