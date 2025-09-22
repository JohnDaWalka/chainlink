package v1_6_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/siloed_lock_release_token_pool"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

func TestSiloedLockReleaseTokenPoolUpdateDesignations(t *testing.T) {
	t.Parallel()

	e, _ := testhelpers.NewMemoryEnvironment(t)
	evmSelectors := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	chain1, chain2 := evmSelectors[0], evmSelectors[1]

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	_, tx, siloedTokenPool, err := siloed_lock_release_token_pool.DeploySiloedLockReleaseTokenPool(
		e.Env.BlockChains.EVMChains()[chain1].DeployerKey,
		e.Env.BlockChains.EVMChains()[chain1].Client,
		state.Chains[chain1].LinkToken.Address(),
		18,
		[]common.Address{},
		state.Chains[chain1].RMNProxy.Address(),
		state.Chains[chain1].Router.Address(),
	)

	require.NoError(t, err)
	_, err = e.Env.BlockChains.EVMChains()[chain1].Confirm(tx)
	require.NoError(t, err)

	err = e.Env.ExistingAddresses.Save(chain1, siloedTokenPool.Address().Hex(), cldf.NewTypeAndVersion(shared.SiloedLockReleaseTokenPool, deployment.Version1_6_0))
	require.NoError(t, err)

	_, tx, siloedTokenPoolChain2, err := siloed_lock_release_token_pool.DeploySiloedLockReleaseTokenPool(
		e.Env.BlockChains.EVMChains()[chain2].DeployerKey,
		e.Env.BlockChains.EVMChains()[chain2].Client,
		state.Chains[chain2].LinkToken.Address(),
		18,
		[]common.Address{},
		state.Chains[chain2].RMNProxy.Address(),
		state.Chains[chain2].Router.Address(),
	)

	require.NoError(t, err)
	_, err = e.Env.BlockChains.EVMChains()[chain2].Confirm(tx)
	require.NoError(t, err)

	err = e.Env.ExistingAddresses.Save(chain2, siloedTokenPoolChain2.Address().Hex(), cldf.NewTypeAndVersion(shared.SiloedLockReleaseTokenPool, deployment.Version1_6_0))
	require.NoError(t, err)
	aType := shared.SiloedLockReleaseTokenPool
	bType := shared.SiloedLockReleaseTokenPool
	e.Env, err = commonchangeset.Apply(t, e.Env,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_5_1.ConfigureTokenPoolContractsChangeset),
			v1_5_1.ConfigureTokenPoolContractsConfig{
				TokenSymbol: shared.TokenSymbol("LINK"),
				MCMS:        nil,
				PoolUpdates: map[uint64]v1_5_1.TokenPoolConfig{
					chain1: {
						Type:    aType,
						Version: deployment.Version1_6_0,
						ChainUpdates: v1_5_1.RateLimiterPerChain{
							chain2: v1_5_1.RateLimiterConfig{
								Inbound: token_pool.RateLimiterConfig{
									IsEnabled: false,
									Rate:      big.NewInt(0),
									Capacity:  big.NewInt(0),
								},
								Outbound: token_pool.RateLimiterConfig{
									IsEnabled: false,
									Rate:      big.NewInt(0),
									Capacity:  big.NewInt(0),
								},
							},
						},
					},
					chain2: {
						Type:    bType,
						Version: deployment.Version1_6_0,
						ChainUpdates: v1_5_1.RateLimiterPerChain{
							chain1: v1_5_1.RateLimiterConfig{
								Inbound: token_pool.RateLimiterConfig{
									IsEnabled: false,
									Rate:      big.NewInt(0),
									Capacity:  big.NewInt(0),
								},
								Outbound: token_pool.RateLimiterConfig{
									IsEnabled: false,
									Rate:      big.NewInt(0),
									Capacity:  big.NewInt(0),
								},
							},
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)

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

	err = cfg.Validate(e.Env)
	require.NoError(t, err)

	_, err = v1_6.SiloedLockReleaseTokenPoolUpdateDesignations(e.Env, cfg)
	require.NoError(t, err)
}

func TestSiloedLockReleaseTokenPoolSetRebalancer(t *testing.T) {
	t.Parallel()

	e, _ := testhelpers.NewMemoryEnvironment(t)
	evmSelectors := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	chain1 := evmSelectors[0]
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	_, tx, siloedTokenPool, err := siloed_lock_release_token_pool.DeploySiloedLockReleaseTokenPool(
		e.Env.BlockChains.EVMChains()[chain1].DeployerKey,
		e.Env.BlockChains.EVMChains()[chain1].Client,
		state.Chains[chain1].LinkToken.Address(),
		18,
		[]common.Address{},
		state.Chains[chain1].RMNProxy.Address(),
		state.Chains[chain1].Router.Address(),
	)

	require.NoError(t, err)
	_, err = e.Env.BlockChains.EVMChains()[chain1].Confirm(tx)
	require.NoError(t, err)

	err = e.Env.ExistingAddresses.Save(chain1, siloedTokenPool.Address().Hex(), cldf.NewTypeAndVersion(shared.SiloedLockReleaseTokenPool, deployment.Version1_6_0))
	require.NoError(t, err)

	cfg := v1_6.SiloedLockReleaseTokenPoolSetRebalancerChangesetConfig{
		Tokens: map[uint64]map[shared.TokenSymbol]common.Address{
			chain1: {
				shared.TokenSymbol("LINK"): state.Chains[chain1].LinkToken.Address(),
			},
		},
		MCMS: &proposalutils.TimelockConfig{},
	}

	err = cfg.Validate(e.Env)
	require.NoError(t, err)

	_, err = v1_6.SiloedLockReleaseTokenPoolSetRebalancer(e.Env, cfg)
	require.NoError(t, err)
}
