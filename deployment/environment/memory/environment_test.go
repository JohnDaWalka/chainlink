package memory_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestNewZKChains(t *testing.T) {
	numChains := 3
	chains := memory.NewZKChains(t, numChains)
	for i := 0; i < numChains; i++ {
		chainId := chainsel.TEST_90000051.EvmChainID + uint64(i)
		sel, err := chainsel.SelectorFromChainId(chainId)
		chain := chains[sel]
		balance, err := chain.Client.BalanceAt(context.Background(), chain.DeployerKey.From, nil)
		require.NoError(t, err)
		require.Positive(t, balance.Cmp(big.NewInt(0)))
	}
}
