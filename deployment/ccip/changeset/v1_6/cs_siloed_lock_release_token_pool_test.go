package v1_6_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/siloed_lock_release_token_pool"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

func TestSiloedLockReleaseTokenPoolUpdateDesignations(t *testing.T) {
	t.Parallel()

	e, _ := testhelpers.NewMemoryEnvironment(t)
	evmSelectors := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	chain1, chain2 := evmSelectors[0], evmSelectors[1]

	cfg := v1_6.SiloedLockReleaseTokenPoolUpdateDesignationsChangesetConfig{
		Tokens: map[uint64]map[shared.TokenSymbol]v1_6.SiloedLockReleaseTokenPoolUpdateDesignationsConfig{
			chain1: {
				shared.TokenSymbol("LINK"): {
					Removes: []uint64{},
					Adds: []siloed_lock_release_token_pool.SiloedLockReleaseTokenPoolSiloConfigUpdate{
						{
							RemoteChainSelector: chain2,
							Rebalancer:          common.HexToAddress("0x1234567890123456789012345678901234567890"),
						},
					},
				},
			},
		},
		MCMS: &proposalutils.TimelockConfig{},
	}

	err := cfg.Validate(e.Env)
	require.NoError(t, err)

	_, err = v1_6.SiloedLockReleaseTokenPoolUpdateDesignations(e.Env, cfg)
	require.NoError(t, err)
}

func TestSiloedLockReleaseTokenPoolSetRebalancer(t *testing.T) {
	t.Parallel()

	e, _ := testhelpers.NewMemoryEnvironment(t)
	evmSelectors := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	chain1 := evmSelectors[0]

	cfg := v1_6.SiloedLockReleaseTokenPoolSetRebalancerChangesetConfig{
		Tokens: map[uint64]map[shared.TokenSymbol]common.Address{
			chain1: {
				shared.TokenSymbol("LINK"): common.HexToAddress("0x9876543210987654321098765432109876543210"),
			},
		},
		MCMS: &proposalutils.TimelockConfig{},
	}

	err := cfg.Validate(e.Env)
	require.NoError(t, err)

	_, err = v1_6.SiloedLockReleaseTokenPoolSetRebalancer(e.Env, cfg)
	require.NoError(t, err)
}
