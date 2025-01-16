package changeset

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestValidateDeployTokenPoolContractsConfig(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 2,
	})

	tests := []struct {
		Msg         string
		TokenSymbol TokenSymbol
		Input       DeployTokenPoolContractsConfig
		ErrStr      string
	}{
		{
			Msg:    "Token symbol is missing",
			Input:  DeployTokenPoolContractsConfig{},
			ErrStr: "token symbol must be defined",
		},
		{
			Msg: "Chain selector is not valid",
			Input: DeployTokenPoolContractsConfig{
				TokenSymbol: "TEST",
				NewPools: map[uint64]DeployTokenPoolInput{
					0: DeployTokenPoolInput{},
				},
			},
			ErrStr: "failed to validate chain selector 0",
		},
		{
			Msg: "Chain selector doesn't exist in environment",
			Input: DeployTokenPoolContractsConfig{
				TokenSymbol: "TEST",
				NewPools: map[uint64]DeployTokenPoolInput{
					5009297550715157269: DeployTokenPoolInput{},
				},
			},
			ErrStr: "chain with selector 5009297550715157269 does not exist in environment",
		},
		{
			Msg: "Router contract is missing from chain",
			Input: DeployTokenPoolContractsConfig{
				TokenSymbol: "TEST",
				NewPools: map[uint64]DeployTokenPoolInput{
					e.AllChainSelectors()[0]: DeployTokenPoolInput{},
				},
			},
			ErrStr: fmt.Sprintf("missing router on %d", e.AllChainSelectors()[0]),
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			err := test.Input.Validate(e)
			require.Contains(t, err.Error(), test.ErrStr)
		})
	}
}

func TestValidateDeployTokenPoolInput(t *testing.T) {
	t.Parallel()

	e, selectorA, _, tokens, _ := setupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)
	acceptLiquidity := false
	invalidAddress := utils.RandomAddress()

	e = deployTestTokenPools(t, e, map[uint64]DeployTokenPoolInput{
		selectorA: {
			Type:               BurnMintTokenPool,
			TokenAddress:       tokens[selectorA].Address,
			LocalTokenDecimals: 18,
		},
	}, true)

	tests := []struct {
		Msg    string
		Symbol TokenSymbol
		Input  DeployTokenPoolInput
		ErrStr string
	}{
		{
			Msg:    "Token address is missing",
			Input:  DeployTokenPoolInput{},
			ErrStr: "token address must be defined",
		},
		{
			Msg: "Token pool type is missing",
			Input: DeployTokenPoolInput{
				TokenAddress: invalidAddress,
			},
			ErrStr: "type must be defined",
		},
		{
			Msg: "Token pool type is invalid",
			Input: DeployTokenPoolInput{
				TokenAddress: invalidAddress,
				Type:         deployment.ContractType("InvalidTokenPool"),
			},
			ErrStr: "requested token pool type InvalidTokenPool is unknown",
		},
		{
			Msg: "Token address is invalid",
			Input: DeployTokenPoolInput{
				Type:         BurnMintTokenPool,
				TokenAddress: invalidAddress,
			},
			ErrStr: fmt.Sprintf("failed to fetch symbol from token with address %s", invalidAddress),
		},
		{
			Msg:    "Token symbol mismatch",
			Symbol: "WRONG",
			Input: DeployTokenPoolInput{
				Type:         BurnMintTokenPool,
				TokenAddress: tokens[selectorA].Address,
			},
			ErrStr: fmt.Sprintf("symbol of token with address %s (%s) does not match expected symbol (WRONG)", tokens[selectorA].Address, testTokenSymbol),
		},
		{
			Msg:    "Token decimal mismatch",
			Symbol: testTokenSymbol,
			Input: DeployTokenPoolInput{
				Type:               BurnMintTokenPool,
				TokenAddress:       tokens[selectorA].Address,
				LocalTokenDecimals: 17,
			},
			ErrStr: fmt.Sprintf("decimals of token with address %s (%d) does not match localTokenDecimals (17)", tokens[selectorA].Address, localTokenDecimals),
		},
		{
			Msg:    "Accept liquidity should be defined",
			Symbol: testTokenSymbol,
			Input: DeployTokenPoolInput{
				Type:               LockReleaseTokenPool,
				TokenAddress:       tokens[selectorA].Address,
				LocalTokenDecimals: 18,
			},
			ErrStr: "accept liquidity must be defined for lock release pools",
		},
		{
			Msg:    "Accept liquidity should be omitted",
			Symbol: testTokenSymbol,
			Input: DeployTokenPoolInput{
				Type:               BurnMintTokenPool,
				TokenAddress:       tokens[selectorA].Address,
				LocalTokenDecimals: 18,
				AcceptLiquidity:    &acceptLiquidity,
			},
			ErrStr: "accept liquidity must be nil for burn mint pools",
		},
		{
			Msg:    "Token pool already exists",
			Symbol: testTokenSymbol,
			Input: DeployTokenPoolInput{
				Type:               BurnMintTokenPool,
				TokenAddress:       tokens[selectorA].Address,
				LocalTokenDecimals: 18,
			},
			ErrStr: fmt.Sprintf("token pool already exists for %s on %d (use forceDeployment to bypass)", testTokenSymbol, selectorA),
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			state, err := LoadOnchainState(e)
			require.NoError(t, err)

			err = test.Input.Validate(e.Chains[selectorA], state.Chains[selectorA], test.Symbol)
			require.Contains(t, err.Error(), test.ErrStr)
		})
	}
}

func TestDeployTokenPool(t *testing.T) {
	t.Parallel()

	e, selectorA, _, tokens, _ := setupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)
	acceptLiquidity := false

	tests := []struct {
		Msg   string
		Input DeployTokenPoolInput
	}{
		{
			Msg: "BurnMint",
			Input: DeployTokenPoolInput{
				TokenAddress:       tokens[selectorA].Address,
				Type:               BurnMintTokenPool,
				LocalTokenDecimals: localTokenDecimals,
				AllowList:          []common.Address{},
			},
		},
		{
			Msg: "BurnWithFromMint",
			Input: DeployTokenPoolInput{
				TokenAddress:       tokens[selectorA].Address,
				Type:               BurnWithFromMintTokenPool,
				LocalTokenDecimals: localTokenDecimals,
				AllowList:          []common.Address{},
			},
		},
		{
			Msg: "BurnFromMint",
			Input: DeployTokenPoolInput{
				TokenAddress:       tokens[selectorA].Address,
				Type:               BurnFromMintTokenPool,
				LocalTokenDecimals: localTokenDecimals,
				AllowList:          []common.Address{},
			},
		},
		{
			Msg: "LockRelease",
			Input: DeployTokenPoolInput{
				TokenAddress:       tokens[selectorA].Address,
				Type:               LockReleaseTokenPool,
				LocalTokenDecimals: localTokenDecimals,
				AllowList:          []common.Address{},
				AcceptLiquidity:    &acceptLiquidity,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			addressBook := deployment.NewMemoryAddressBook()
			state, err := LoadOnchainState(e)
			require.NoError(t, err)

			_, err = deployTokenPool(
				e.Logger,
				e.Chains[selectorA],
				state.Chains[selectorA],
				addressBook,
				test.Input,
			)
			require.NoError(t, err)

			err = e.ExistingAddresses.Merge(addressBook)
			require.NoError(t, err)

			state, err = LoadOnchainState(e)
			require.NoError(t, err)

			switch test.Input.Type {
			case BurnMintTokenPool:
				_, ok := state.Chains[selectorA].BurnMintTokenPools[testTokenSymbol]
				require.True(t, ok)
			case LockReleaseTokenPool:
				_, ok := state.Chains[selectorA].LockReleaseTokenPools[testTokenSymbol]
				require.True(t, ok)
			case BurnWithFromMintTokenPool:
				_, ok := state.Chains[selectorA].BurnWithFromMintTokenPools[testTokenSymbol]
				require.True(t, ok)
			case BurnFromMintTokenPool:
				_, ok := state.Chains[selectorA].BurnFromMintTokenPools[testTokenSymbol]
				require.True(t, ok)
			}
		})
	}
}

func TestDeployTokenPoolContracts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Msg             string
		Redeploy        bool
		ForceDeployment bool
		ErrStr          string
	}{
		{
			Msg: "Deploy once",
		},
		{
			Msg:             "Redeploy but don't force redeployment",
			Redeploy:        true,
			ForceDeployment: false,
			ErrStr:          "token pool already exists for TEST",
		},
		{
			Msg:             "Redeploy with force",
			Redeploy:        true,
			ForceDeployment: true,
		},
	}

	for _, test := range tests {

		e, selectorA, _, tokens, timelockContracts := setupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)
		changesetApplication := commonchangeset.ChangesetApplication{
			Changeset: commonchangeset.WrapChangeSet(DeployTokenPoolContracts),
			Config: DeployTokenPoolContractsConfig{
				TokenSymbol: testTokenSymbol,
				NewPools: map[uint64]DeployTokenPoolInput{
					selectorA: {
						TokenAddress:       tokens[selectorA].Address,
						Type:               BurnMintTokenPool,
						LocalTokenDecimals: localTokenDecimals,
						AllowList:          []common.Address{},
						ForceDeployment:    test.ForceDeployment,
					},
				},
			},
		}

		// Initial deployment
		e, err := commonchangeset.ApplyChangesets(t, e, timelockContracts, []commonchangeset.ChangesetApplication{
			changesetApplication,
		})
		require.NoError(t, err)

		state, err := LoadOnchainState(e)
		require.NoError(t, err)

		burnMintTokenPools, ok := state.Chains[selectorA].BurnMintTokenPools[testTokenSymbol]
		require.True(t, ok)
		require.Len(t, burnMintTokenPools, 1)
		owner, err := burnMintTokenPools[0].Owner(nil)
		require.NoError(t, err)
		require.Equal(t, e.Chains[selectorA].DeployerKey.From, owner)

		// Redeployment
		if test.Redeploy {
			e, err = commonchangeset.ApplyChangesets(t, e, timelockContracts, []commonchangeset.ChangesetApplication{
				changesetApplication,
			})
			if test.ErrStr != "" {
				require.ErrorContains(t, err, fmt.Sprintf("token pool already exists for TEST on %d (use forceDeployment to bypass)", selectorA))
			} else {
				state, err = LoadOnchainState(e)
				require.NoError(t, err)

				burnMintTokenPools, ok = state.Chains[selectorA].BurnMintTokenPools[testTokenSymbol]
				require.True(t, ok)
				require.Len(t, burnMintTokenPools, 2)
			}
		}
	}
}
