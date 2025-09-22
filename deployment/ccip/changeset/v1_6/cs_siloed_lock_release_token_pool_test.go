package v1_6_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/burn_mint_erc20"
	"github.com/smartcontractkit/chainlink/deployment"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/siloed_lock_release_token_pool"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

func deployTokenAndPoolPrerequisites(
	t *testing.T,
	e cldf.Environment,
	chain cldf_evm.Chain,
	addressBook cldf.AddressBook,
) cldf.Environment {
	siloedToken, err := cldf.DeployContract(e.Logger, chain, addressBook,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*burn_mint_erc20.BurnMintERC20] {
			tokenAddress, tx, token, err := burn_mint_erc20.DeployBurnMintERC20(
				chain.DeployerKey,
				chain.Client,
				"TESTSILOED",
				"TESTSILOED",
				6,
				big.NewInt(0),
				big.NewInt(0),
			)
			return cldf.ContractDeploy[*burn_mint_erc20.BurnMintERC20]{
				Address:  tokenAddress,
				Contract: token,
				Tv:       cldf.NewTypeAndVersion(shared.BurnMintERC20Token, deployment.Version1_0_0),
				Tx:       tx,
				Err:      err,
			}
		},
	)
	require.NoError(t, err)

	_, err = cldf.DeployContract(e.Logger, chain, addressBook,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*siloed_lock_release_token_pool.SiloedLockReleaseTokenPool] {
			poolAddress, tx, token, err := siloed_lock_release_token_pool.DeploySiloedLockReleaseTokenPool(
				chain.DeployerKey,
				chain.Client,
				siloedToken.Address,
				6,
				[]common.Address{},
				common.HexToAddress("0x1234567890123456789012345678901234567890"),
				common.HexToAddress("0x1234567890123456789012345678901234567890"),
			)
			return cldf.ContractDeploy[*siloed_lock_release_token_pool.SiloedLockReleaseTokenPool]{
				Address:  poolAddress,
				Contract: token,
				Tv:       cldf.NewTypeAndVersion(shared.SiloedLockReleaseTokenPool, deployment.Version1_6_1),
				Tx:       tx,
				Err:      err,
			}
		},
	)
	require.NoError(t, err)
	return e
}

func TestSiloedLockReleaseTokenPoolUpdateDesignations(t *testing.T) {
	t.Parallel()

	e, _ := testhelpers.NewMemoryEnvironment(t)
	evmSelectors := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	chain1, chain2 := evmSelectors[0], evmSelectors[1]

	addressBook := cldf.NewMemoryAddressBook()
	modifiedEnv := deployTokenAndPoolPrerequisites(t, e.Env, e.Env.BlockChains.EVMChains()[chain1], addressBook)

	cfg := v1_6.SiloedLockReleaseTokenPoolUpdateDesignationsChangesetConfig{
		Tokens: map[uint64]map[shared.TokenSymbol]v1_6.SiloedLockReleaseTokenPoolUpdateDesignationsConfig{
			chain1: {
				shared.TokenSymbol("TESTSILOED"): {
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

	err := cfg.Validate(modifiedEnv)
	require.NoError(t, err)

	_, err = v1_6.SiloedLockReleaseTokenPoolUpdateDesignations(modifiedEnv, cfg)
	require.NoError(t, err)
}

func TestSiloedLockReleaseTokenPoolSetRebalancer(t *testing.T) {
	t.Parallel()

	e, _ := testhelpers.NewMemoryEnvironment(t)
	evmSelectors := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	chain1 := evmSelectors[0]

	addressBook := cldf.NewMemoryAddressBook()
	modifiedEnv := deployTokenAndPoolPrerequisites(t, e.Env, e.Env.BlockChains.EVMChains()[chain1], addressBook)

	cfg := v1_6.SiloedLockReleaseTokenPoolSetRebalancerChangesetConfig{
		Tokens: map[uint64]map[shared.TokenSymbol]common.Address{
			chain1: {
				shared.TokenSymbol("TESTSILOED"): common.HexToAddress("0x9876543210987654321098765432109876543210"),
			},
		},
		MCMS: &proposalutils.TimelockConfig{},
	}

	err := cfg.Validate(modifiedEnv)
	require.NoError(t, err)

	_, err = v1_6.SiloedLockReleaseTokenPoolSetRebalancer(modifiedEnv, cfg)
	require.NoError(t, err)
}
