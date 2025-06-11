package sui

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pattonkan/sui-go/sui"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	ccip_ops "github.com/smartcontractkit/chainlink-sui/ops/ccip"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

const (
	mockMCMSAddress = "0x3f20aa841a0eb5c038775bdb868924770df1ce377cc0013b3ba4ac9fd69a4f90"
	mockAddress     = "0x13a9f1a109368730f2e355d831ba8fbf5942fb82321863d55de54cb4ebe5d18f"
	mockLinkAddress = "0xa"

	sepChainSelector     = 11155111
	sepMockOnRampAddress = "0x0BF3dE8c5D3e8A2B34D2BEeB17ABfCeBaf363A59"
)

func GetMockChainContractParams(t *testing.T, e cldf.Environment, chainSelector uint64) (ChainContractParams, error) {
	mockParsedAddress := sui.MustAddressFromHex(mockAddress)
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return ChainContractParams{}, err
	}

	linkTokenObjectMetadataId := state.SuiChains[chainSelector].LinkTokenCoinMetadataId.String()

	return ChainContractParams{
		FeeQuoterParams: ccip_ops.InitFeeQuoterInput{
			MaxFeeJuelsPerMsg:             "100000000",
			TokenPriceStalenessThreshold:  1000000,
			LinkTokenCoinMetadataObjectId: linkTokenObjectMetadataId,           // CoinMetadataObjectId
			FeeTokens:                     []string{linkTokenObjectMetadataId}, // CoinMetadataObjectId
		},
		OffRampParams: OffRampParams{
			ChainSelector:                    chainSelector,
			PermissionlessExecutionThreshold: uint32(60 * 60 * 8),
			IsRMNVerificationDisabled:        []bool{false},
			SourceChainSelectors:             []uint64{sepChainSelector},
			SourceChainIsEnabled:             []bool{true},
			SourceChainsOnRamp:               [][]byte{common.HexToAddress(sepMockOnRampAddress).Bytes()},
		},
		OnRampParams: OnRampParams{
			ChainSelector:  chainSelector,
			AllowlistAdmin: *mockParsedAddress,
			FeeAggregator:  *mockParsedAddress,
		},
	}, nil
}
