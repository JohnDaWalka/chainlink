package ccip

import (
	"fmt"
	"golang.org/x/exp/maps"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

func Test_CCIPMessaging_Sui2EVM(t *testing.T) {
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithSuiChains(1),
	)

	evmChainSelectors := maps.Keys(e.Env.Chains)
	suiChains, err := e.Env.BlockChains.SuiChains()
	require.NoError(t, err)
	suiChainSelectors := maps.Keys(suiChains)
	require.Equal(t, len(suiChainSelectors), 1)

	fmt.Println("EVM: ", evmChainSelectors)
	fmt.Println("Sui: ", suiChainSelectors)
}
