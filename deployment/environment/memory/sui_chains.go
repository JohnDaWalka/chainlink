package memory

import (
	"testing"

	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	suichain "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf_sui_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui/provider"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func getTestSuiChainSelectors() []uint64 {
	// TODO: CTF to support different chain ids, need to investigate if it's possible (thru node config.yaml?)
	return []uint64{chainsel.SUI_LOCALNET.Selector}
}

func GenerateChainsSui(t *testing.T, numChains int) []cldf_chain.BlockChain {
	testSuiChainSelectors := getTestSuiChainSelectors()
	if len(testSuiChainSelectors) < numChains {
		t.Fatalf("not enough test sui chain selectors available")
	}
	chains := make([]cldf_chain.BlockChain, 0, numChains)
	for i := range numChains {
		selector := testSuiChainSelectors[i]

		c, err := cldf_sui_provider.NewCTFChainProvider(t, selector,
			cldf_sui_provider.CTFChainProviderConfig{
				Once: once,
			},
		).Initialize(t.Context())
		require.NoError(t, err)

		chains = append(chains, c)
	}
	t.Logf("Created %d Sui chains: %+v", len(chains), chains)
	return chains
}

func createSuiChainConfig(chainID string, chain suichain.Chain) chainlink.RawConfig {
	chainConfig := chainlink.RawConfig{}

	chainConfig["Enabled"] = true
	chainConfig["ChainID"] = chainID
	chainConfig["NetworkName"] = "sui-localnet"
	chainConfig["NetworkNameFull"] = "sui-localnet"
	chainConfig["Nodes"] = []any{
		map[string]any{
			"Name": "primary",
			"URL":  chain.URL,
		},
	}

	return chainConfig
}
