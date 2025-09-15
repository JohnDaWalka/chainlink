package v1_6_3_dev

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	fqSui "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/fee_quoter"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
)

type DeployFeeQInput struct {
	Chain         uint64
	Params        FeeQuoterParamsSui
	LinkAddr      common.Address
	WethAddr      common.Address
	PriceUpdaters []common.Address
}

type ApplyTokenTransferFeeConfigUpdatesConfigPerChain struct {
	TokenTransferFeeConfigs       []fqSui.FeeQuoterTokenTransferFeeConfigArgs
	TokenTransferFeeConfigsRemove []fqSui.FeeQuoterTokenTransferFeeConfigRemoveArgs
}

type ApplyFeeTokensUpdatesInput struct {
	FeeTokensToAdd    []common.Address
	FeeTokensToRemove []common.Address
}

var (
	DeploySuiSupportedFeeQuoterOp = opsutil.NewEVMDeployOperation(
		"DeploySuiSupportedFeeQuoterOp",
		semver.MustParse("1.6.3-dev"),
		"Deploys FeeQuoter 1.6.3-dev contract on the specified evm chain",
		cldf.NewTypeAndVersion(shared.FeeQuoter, deployment.Version1_6_3Dev),
		opsutil.VMDeployers[DeployFeeQInput]{
			DeployEVM: func(opts *bind.TransactOpts, backend bind.ContractBackend, input DeployFeeQInput) (common.Address, *types.Transaction, error) {
				addr, tx, _, err := fqSui.DeployFeeQuoter(opts, backend,
					fqSui.FeeQuoterStaticConfig{
						MaxFeeJuelsPerMsg:            input.Params.MaxFeeJuelsPerMsg,
						LinkToken:                    input.LinkAddr,
						TokenPriceStalenessThreshold: input.Params.TokenPriceStalenessThreshold,
					},
					input.PriceUpdaters,
					[]common.Address{input.WethAddr, input.LinkAddr}, // fee tokens
					input.Params.TokenPriceFeedUpdates,
					input.Params.TokenTransferFeeConfigArgs,
					append([]fqSui.FeeQuoterPremiumMultiplierWeiPerEthArgs{
						{
							PremiumMultiplierWeiPerEth: input.Params.LinkPremiumMultiplierWeiPerEth,
							Token:                      input.LinkAddr,
						},
						{
							PremiumMultiplierWeiPerEth: input.Params.WethPremiumMultiplierWeiPerEth,
							Token:                      input.WethAddr,
						},
					}, input.Params.MorePremiumMultiplierWeiPerEth...),
					input.Params.DestChainConfigArgs,
				)
				return addr, tx, err
			},
			DeployZksyncVM: func(opts *accounts.TransactOpts, client *clients.Client, wallet *accounts.Wallet, backend bind.ContractBackend, input DeployFeeQInput) (common.Address, error) {
				addr, _, _, err := fqSui.DeployFeeQuoterZk(opts, client, wallet, backend,
					fqSui.FeeQuoterStaticConfig{
						MaxFeeJuelsPerMsg:            input.Params.MaxFeeJuelsPerMsg,
						LinkToken:                    input.LinkAddr,
						TokenPriceStalenessThreshold: input.Params.TokenPriceStalenessThreshold,
					},
					input.PriceUpdaters,
					[]common.Address{input.WethAddr, input.LinkAddr}, // fee tokens
					input.Params.TokenPriceFeedUpdates,
					input.Params.TokenTransferFeeConfigArgs,
					append([]fqSui.FeeQuoterPremiumMultiplierWeiPerEthArgs{
						{
							PremiumMultiplierWeiPerEth: input.Params.LinkPremiumMultiplierWeiPerEth,
							Token:                      input.LinkAddr,
						},
						{
							PremiumMultiplierWeiPerEth: input.Params.WethPremiumMultiplierWeiPerEth,
							Token:                      input.WethAddr,
						},
					}, input.Params.MorePremiumMultiplierWeiPerEth...),
					input.Params.DestChainConfigArgs)
				return addr, err
			},
		})
)

type FeeQuoterParamsSui struct {
	MaxFeeJuelsPerMsg              *big.Int
	TokenPriceStalenessThreshold   uint32
	LinkPremiumMultiplierWeiPerEth uint64
	WethPremiumMultiplierWeiPerEth uint64
	MorePremiumMultiplierWeiPerEth []fqSui.FeeQuoterPremiumMultiplierWeiPerEthArgs
	TokenPriceFeedUpdates          []fqSui.FeeQuoterTokenPriceFeedUpdate
	TokenTransferFeeConfigArgs     []fqSui.FeeQuoterTokenTransferFeeConfigArgs
	DestChainConfigArgs            []fqSui.FeeQuoterDestChainConfigArgs
}

func (c FeeQuoterParamsSui) Validate() error {
	if c.MaxFeeJuelsPerMsg == nil {
		return errors.New("MaxFeeJuelsPerMsg is nil")
	}
	if c.MaxFeeJuelsPerMsg.Cmp(big.NewInt(0)) <= 0 {
		return errors.New("MaxFeeJuelsPerMsg must be positive")
	}
	if c.TokenPriceStalenessThreshold == 0 {
		return errors.New("TokenPriceStalenessThreshold can't be 0")
	}
	return nil
}

func DefaultFeeQuoterParams() FeeQuoterParamsSui {
	return FeeQuoterParamsSui{
		MaxFeeJuelsPerMsg:              big.NewInt(0).Mul(big.NewInt(2e2), big.NewInt(1e18)),
		TokenPriceStalenessThreshold:   uint32(24 * 60 * 60),
		LinkPremiumMultiplierWeiPerEth: 9e17, // 0.9 ETH
		WethPremiumMultiplierWeiPerEth: 1e18, // 1.0 ETH
		TokenPriceFeedUpdates:          []fqSui.FeeQuoterTokenPriceFeedUpdate{},
		TokenTransferFeeConfigArgs:     []fqSui.FeeQuoterTokenTransferFeeConfigArgs{},
		MorePremiumMultiplierWeiPerEth: []fqSui.FeeQuoterPremiumMultiplierWeiPerEthArgs{},
		DestChainConfigArgs:            []fqSui.FeeQuoterDestChainConfigArgs{},
	}
}

const (
	// https://github.com/smartcontractkit/chainlink/blob/1423e2581e8640d9e5cd06f745c6067bb2893af2/contracts/src/v0.8/ccip/libraries/Internal.sol#L275-L279
	/*
				```Solidity
					// bytes4(keccak256("CCIP ChainFamilySelector EVM"))
					bytes4 public constant CHAIN_FAMILY_SELECTOR_EVM = 0x2812d52c;
					// bytes4(keccak256("CCIP ChainFamilySelector SVM"));
		  		bytes4 public constant CHAIN_FAMILY_SELECTOR_SVM = 0x1e10bdc4;
				```
	*/
	EVMFamilySelector   = "2812d52c"
	SVMFamilySelector   = "1e10bdc4"
	AptosFamilySelector = "ac77ffec"
)

func DefaultFeeQuoterDestChainConfig(configEnabled bool, destChainSelector ...uint64) fqSui.FeeQuoterDestChainConfig {
	familySelector, _ := hex.DecodeString(EVMFamilySelector) // evm
	if len(destChainSelector) > 0 {
		destFamily, _ := chain_selectors.GetSelectorFamily(destChainSelector[0])
		switch destFamily {
		case chain_selectors.FamilySolana:
			familySelector, _ = hex.DecodeString(SVMFamilySelector) // solana
		case chain_selectors.FamilyAptos:
			familySelector, _ = hex.DecodeString(AptosFamilySelector) // aptos
		}
	}
	return fqSui.FeeQuoterDestChainConfig{
		IsEnabled:                         configEnabled,
		MaxNumberOfTokensPerMsg:           10,
		MaxDataBytes:                      30_000,
		MaxPerMsgGasLimit:                 3_000_000, // TODO: this needs to be updated based on RMN sig verification per chain?! 220/250K
		DestGasOverhead:                   ccipevm.DestGasOverhead,
		DefaultTokenFeeUSDCents:           25,
		DestGasPerPayloadByteBase:         ccipevm.CalldataGasPerByteBase,
		DestGasPerPayloadByteHigh:         ccipevm.CalldataGasPerByteHigh,
		DestGasPerPayloadByteThreshold:    ccipevm.CalldataGasPerByteThreshold,
		DestDataAvailabilityOverheadGas:   100,
		DestGasPerDataAvailabilityByte:    16,
		DestDataAvailabilityMultiplierBps: 1,
		DefaultTokenDestGasOverhead:       90_000,
		DefaultTxGasLimit:                 200_000,
		GasMultiplierWeiPerEth:            11e17, // Gas multiplier in wei per eth is scaled by 1e18, so 11e17 is 1.1 = 110%
		NetworkFeeUSDCents:                10,
		ChainFamilySelector:               [4]byte(familySelector),
		GasPriceStalenessThreshold:        90000,
	}
}
