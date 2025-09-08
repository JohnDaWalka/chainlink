package sui

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	ccip_ops "github.com/smartcontractkit/chainlink-sui/ops/ccip"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

const (
	mockMCMSAddress = "0x3f20aa841a0eb5c038775bdb868924770df1ce377cc0013b3ba4ac9fd69a4f90"
	mockAddress     = "0x13a9f1a109368730f2e355d831ba8fbf5942fb82321863d55de54cb4ebe5d18f"
	mockLinkAddress = "0xa"

	sepChainSelector     = 18395503381733958356
	sepMockOnRampAddress = "0x0BF3dE8c5D3e8A2B34D2BEeB17ABfCeBaf363A59"
)

func GetMockChainContractParams(t *testing.T, e cldf.Environment, chainSelector uint64) (ChainContractParams, error) {

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return ChainContractParams{}, err
	}

	linkTokenObjectMetadataId := state.SuiChains[chainSelector].LinkTokenCoinMetadataId

	return ChainContractParams{
		DestChainSelector: 909606746561742123,
		FeeQuoterParams: ccip_ops.InitFeeQuoterInput{ //1112246792648961560000000
			MaxFeeJuelsPerMsg:             "200000000000000000000", // 1000 LINK (101 224 679 264 896 156)
			TokenPriceStalenessThreshold:  1000000,
			LinkTokenCoinMetadataObjectId: linkTokenObjectMetadataId,           // CoinMetadataObjectId
			FeeTokens:                     []string{linkTokenObjectMetadataId}, // CoinMetadataObjectId
		},
		OffRampParams: OffRampParams{
			ChainSelector:                    chainSelector,
			PermissionlessExecutionThreshold: uint32(60 * 60 * 8),
			IsRMNVerificationDisabled:        []bool{true},
			SourceChainSelectors:             []uint64{sepChainSelector},
			SourceChainIsEnabled:             []bool{true},
			SourceChainsOnRamp:               [][]byte{common.HexToAddress(sepMockOnRampAddress).Bytes()},
		},
		OnRampParams: OnRampParams{
			ChainSelector:  chainSelector,
			AllowlistAdmin: mockAddress,
			FeeAggregator:  mockAddress,
		},
	}, nil
}
