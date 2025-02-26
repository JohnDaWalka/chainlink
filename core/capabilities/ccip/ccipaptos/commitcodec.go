package ccipaptos

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/ccip_aptos_utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
)

var (
	aptosUtilsABI = abihelpers.MustParseABI(ccip_aptos_utils.AptosUtilsABI)
)

// CommitPluginCodecV1 is a codec for encoding and decoding commit plugin reports.
// Compatible with ccip::offramp version 1.6.0
type CommitPluginCodecV1 struct{}

func NewCommitPluginCodecV1() *CommitPluginCodecV1 {
	return &CommitPluginCodecV1{}
}

func (c *CommitPluginCodecV1) Encode(ctx context.Context, report cciptypes.CommitPluginReport) ([]byte, error) {
	isBlessed := make(map[cciptypes.ChainSelector]bool)
	for _, root := range report.BlessedMerkleRoots {
		isBlessed[root.ChainSel] = true
	}

	blessedMerkleRoots := make([]ccip_aptos_utils.AptosUtilsMerkleRoot, 0, len(report.BlessedMerkleRoots))
	unblessedMerkleRoots := make([]ccip_aptos_utils.AptosUtilsMerkleRoot, 0, len(report.UnblessedMerkleRoots))

	for _, root := range append(report.BlessedMerkleRoots, report.UnblessedMerkleRoots...) {
		imr := ccip_aptos_utils.AptosUtilsMerkleRoot{
			SourceChainSelector: uint64(root.ChainSel),
			OnRampAddress:       root.OnRampAddress,
			MinSequenceNumber:   uint64(root.SeqNumsRange.Start()),
			MaxSequenceNumber:   uint64(root.SeqNumsRange.End()),
			MerkleRoot:          root.MerkleRoot,
		}
		if isBl, ok := isBlessed[root.ChainSel]; ok && isBl {
			blessedMerkleRoots = append(blessedMerkleRoots, imr)
		} else {
			unblessedMerkleRoots = append(unblessedMerkleRoots, imr)
		}
	}

	rmnSignatures := make([]ccip_aptos_utils.AptosUtilsRMNSignature, 0, len(report.RMNSignatures))
	for _, sig := range report.RMNSignatures {
		rmnSignatures = append(rmnSignatures, ccip_aptos_utils.AptosUtilsRMNSignature{
			R: sig.R,
			S: sig.S,
		})
	}

	tokenPriceUpdates := make([]ccip_aptos_utils.AptosUtilsTokenPriceUpdate, 0, len(report.PriceUpdates.TokenPriceUpdates))
	for _, update := range report.PriceUpdates.TokenPriceUpdates {
		if !addressIsValid(string(update.TokenID)) {
			return nil, fmt.Errorf("invalid token address: %s", update.TokenID)
		}
		if update.Price.IsEmpty() {
			return nil, fmt.Errorf("empty price for token: %s", update.TokenID)
		}
		sourceToken, err := addressStringToBytes32(string(update.TokenID))
		if err != nil {
			return nil, fmt.Errorf("failed to convert token address to bytes32: %w", err)
		}
		tokenPriceUpdates = append(tokenPriceUpdates, ccip_aptos_utils.AptosUtilsTokenPriceUpdate{
			SourceToken: sourceToken,
			UsdPerToken: update.Price.Int,
		})
	}

	gasPriceUpdates := make([]ccip_aptos_utils.AptosUtilsGasPriceUpdate, 0, len(report.PriceUpdates.GasPriceUpdates))
	for _, update := range report.PriceUpdates.GasPriceUpdates {
		if update.GasPrice.IsEmpty() {
			return nil, fmt.Errorf("empty gas price for chain: %d", update.ChainSel)
		}

		gasPriceUpdates = append(gasPriceUpdates, ccip_aptos_utils.AptosUtilsGasPriceUpdate{
			DestChainSelector: uint64(update.ChainSel),
			UsdPerUnitGas:     update.GasPrice.Int,
		})
	}

	priceUpdates := ccip_aptos_utils.AptosUtilsPriceUpdates{
		TokenPriceUpdates: tokenPriceUpdates,
		GasPriceUpdates:   gasPriceUpdates,
	}

	commitReport := &ccip_aptos_utils.AptosUtilsCommitReport{
		PriceUpdates:         priceUpdates,
		BlessedMerkleRoots:   blessedMerkleRoots,
		UnblessedMerkleRoots: unblessedMerkleRoots,
		RmnSignatures:        rmnSignatures,
	}

	packed, err := aptosUtilsABI.Pack("exposeCommitReport", commitReport)
	if err != nil {
		return nil, fmt.Errorf("failed to pack commit report: %w", err)
	}

	return packed[4:], nil
}

func (c *CommitPluginCodecV1) Decode(ctx context.Context, bytes []byte) (cciptypes.CommitPluginReport, error) {
	method, ok := aptosUtilsABI.Methods["exposeCommitReport"]
	if !ok {
		return cciptypes.CommitPluginReport{}, errors.New("missing method exposeCommitReport")
	}

	unpacked, err := method.Inputs.Unpack(bytes)
	if err != nil {
		return cciptypes.CommitPluginReport{}, fmt.Errorf("failed to unpack commit report: %w", err)
	}
	if len(unpacked) != 1 {
		return cciptypes.CommitPluginReport{}, fmt.Errorf("expected 1 argument, got %d", len(unpacked))
	}

	commitReport := *abi.ConvertType(unpacked[0], new(ccip_aptos_utils.AptosUtilsCommitReport)).(*ccip_aptos_utils.AptosUtilsCommitReport)

	blessedMerkleRoots := make([]cciptypes.MerkleRootChain, 0, len(commitReport.BlessedMerkleRoots))
	for _, root := range commitReport.BlessedMerkleRoots {
		mrc := cciptypes.MerkleRootChain{
			ChainSel:      cciptypes.ChainSelector(root.SourceChainSelector),
			OnRampAddress: root.OnRampAddress,
			SeqNumsRange: cciptypes.NewSeqNumRange(
				cciptypes.SeqNum(root.MinSequenceNumber),
				cciptypes.SeqNum(root.MaxSequenceNumber),
			),
			MerkleRoot: root.MerkleRoot,
		}
		blessedMerkleRoots = append(blessedMerkleRoots, mrc)
	}

	unblessedMerkleRoots := make([]cciptypes.MerkleRootChain, 0, len(commitReport.UnblessedMerkleRoots))
	for _, root := range commitReport.UnblessedMerkleRoots {
		mrc := cciptypes.MerkleRootChain{
			ChainSel:      cciptypes.ChainSelector(root.SourceChainSelector),
			OnRampAddress: root.OnRampAddress,
			SeqNumsRange: cciptypes.NewSeqNumRange(
				cciptypes.SeqNum(root.MinSequenceNumber),
				cciptypes.SeqNum(root.MaxSequenceNumber),
			),
			MerkleRoot: root.MerkleRoot,
		}
		unblessedMerkleRoots = append(unblessedMerkleRoots, mrc)
	}

	tokenPriceUpdates := make([]cciptypes.TokenPrice, 0, len(commitReport.PriceUpdates.TokenPriceUpdates))
	for _, update := range commitReport.PriceUpdates.TokenPriceUpdates {
		sourceTokenStr, err := addressBytesToString(update.SourceToken[:])
		if err != nil {
			return cciptypes.CommitPluginReport{}, fmt.Errorf("failed to convert token address %v to string: %w", update.SourceToken, err)
		}
		tokenPriceUpdates = append(tokenPriceUpdates, cciptypes.TokenPrice{
			TokenID: cciptypes.UnknownEncodedAddress(sourceTokenStr),
			Price:   cciptypes.NewBigInt(big.NewInt(0).Set(update.UsdPerToken)),
		})
	}

	gasPriceUpdates := make([]cciptypes.GasPriceChain, 0, len(commitReport.PriceUpdates.GasPriceUpdates))
	for _, update := range commitReport.PriceUpdates.GasPriceUpdates {
		gasPriceUpdates = append(gasPriceUpdates, cciptypes.GasPriceChain{
			GasPrice: cciptypes.NewBigInt(big.NewInt(0).Set(update.UsdPerUnitGas)),
			ChainSel: cciptypes.ChainSelector(update.DestChainSelector),
		})
	}

	rmnSignatures := make([]cciptypes.RMNECDSASignature, 0, len(commitReport.RmnSignatures))
	for _, sig := range commitReport.RmnSignatures {
		rmnSignatures = append(rmnSignatures, cciptypes.RMNECDSASignature{
			R: sig.R,
			S: sig.S,
		})
	}

	return cciptypes.CommitPluginReport{
		BlessedMerkleRoots:   blessedMerkleRoots,
		UnblessedMerkleRoots: unblessedMerkleRoots,
		PriceUpdates: cciptypes.PriceUpdates{
			TokenPriceUpdates: tokenPriceUpdates,
			GasPriceUpdates:   gasPriceUpdates,
		},
		RMNSignatures: rmnSignatures,
	}, nil
}

// Ensure CommitPluginCodec implements the CommitPluginCodec interface
var _ cciptypes.CommitPluginCodec = (*CommitPluginCodecV1)(nil)
