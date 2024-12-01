package ccipevm

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	commoncodec "github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/codec"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

const commitReportABI = `[{"components":[{"name":"chainSel","type":"uint64","internalType":"uint64"},{"name":"onRampAddress","type":"bytes","internalType":"bytes"},{"name":"seqNumsRange","type":"uint64[2]","internalType":"uint64[2]"},{"name":"merkleRoot","type":"bytes32","internalType":"bytes32"}],"name":"merkleRoots","type":"tuple[]","internalType":"struct MerkleRootChain[]"},{"components":[{"components":[{"name":"tokenID","type":"string","internalType":"string"},{"name":"price","type":"uint256","internalType":"uint256"}],"name":"tokenPriceUpdates","type":"tuple[]","internalType":"struct TokenPrice[]"},{"components":[{"name":"chainSel","type":"uint64","internalType":"uint64"},{"name":"gasPrice","type":"uint256","internalType":"uint256"}],"name":"gasPriceUpdates","type":"tuple[]","internalType":"struct GasPriceChain[]"}],"name":"priceUpdates","type":"tuple","internalType":"struct PriceUpdates"},{"components":[{"name":"r","type":"bytes32","internalType":"bytes32"},{"name":"s","type":"bytes32","internalType":"bytes32"}],"name":"rmnSignatures","type":"tuple[]"}]`

var commitCodecConfig = types.CodecConfig{
	Configs: map[string]types.ChainCodecConfig{
		"CommitPluginReport": {
			TypeABI: commitReportABI,
			ModifierConfigs: commoncodec.ModifiersConfig{
				&commoncodec.WrapperModifierConfig{Fields: map[string]string{
					"PriceUpdates.TokenPriceUpdates.Price":  "Int",
					"PriceUpdates.GasPriceUpdates.GasPrice": "Int",
				}},
			},
		},
	},
}

// CommitPluginCodecV2 is a codec for encoding and decoding commit plugin reports using generic evm codec
type CommitPluginCodecV2 struct{}

func NewCommitPluginCodecV2() *CommitPluginCodecV2 {
	return &CommitPluginCodecV2{}
}

func validateReport(report cciptypes.CommitPluginReport) error {
	for _, update := range report.PriceUpdates.TokenPriceUpdates {
		if !common.IsHexAddress(string(update.TokenID)) {
			return fmt.Errorf("invalid token address: %s", update.TokenID)
		}
		if update.Price.IsEmpty() {
			return fmt.Errorf("empty price for token: %s", update.TokenID)
		}
	}

	for _, update := range report.PriceUpdates.GasPriceUpdates {
		if update.GasPrice.IsEmpty() {
			return fmt.Errorf("empty gas price for chain: %d", update.ChainSel)
		}
	}

	return nil
}

func (c *CommitPluginCodecV2) Encode(ctx context.Context, report cciptypes.CommitPluginReport) ([]byte, error) {
	if err := validateReport(report); err != nil {
		return nil, err
	}

	cd, err := codec.NewCodec(commitCodecConfig)
	if err != nil {
		return nil, err
	}

	return cd.Encode(ctx, report, "CommitPluginReport")
}

func (c *CommitPluginCodecV2) Decode(ctx context.Context, bytes []byte) (cciptypes.CommitPluginReport, error) {
	report := cciptypes.CommitPluginReport{}
	cd, err := codec.NewCodec(commitCodecConfig)
	if err != nil {
		return report, err
	}

	err = cd.Decode(ctx, bytes, &report, "CommitPluginReport")
	if err != nil {
		return report, err
	}

	return report, nil
}

// Ensure CommitPluginCodec implements the CommitPluginCodec interface
var _ cciptypes.CommitPluginCodec = (*CommitPluginCodecV2)(nil)
