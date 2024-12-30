package ccipsolana

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipsolana/ccip_router"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	agbinary "github.com/gagliardetto/binary"
)

// CommitPluginCodecV1 is a codec for encoding and decoding commit plugin reports.
// Compatible with:
// - "OffRamp 1.6.0-dev"
type CommitPluginCodecV1 struct{}

func NewCommitPluginCodecV1() *CommitPluginCodecV1 {
	return &CommitPluginCodecV1{}
}

func (c *CommitPluginCodecV1) Encode(ctx context.Context, report cciptypes.CommitPluginReport) ([]byte, error) {
	var buf bytes.Buffer
	encoder := agbinary.NewBorshEncoder(&buf)

	mr := ccip_router.MerkleRoot{
		SourceChainSelector: uint64(report.MerkleRoots[0].ChainSel),
		OnRampAddress:       report.MerkleRoots[0].OnRampAddress,
		MinSeqNr:            uint64(report.MerkleRoots[0].SeqNumsRange.Start()),
		MaxSeqNr:            uint64(report.MerkleRoots[0].SeqNumsRange.End()),
		MerkleRoot:          report.MerkleRoots[0].MerkleRoot,
	}

	tpu := make([]ccip_router.TokenPriceUpdate, 0, len(report.PriceUpdates.TokenPriceUpdates))
	for _, update := range report.PriceUpdates.TokenPriceUpdates {
		b, err := update.Price.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("error marshaling token price: %v", err)
		}

		if len(b) > 28 {
			return nil, errors.New("token price is too large")
		}

		tpu = append(tpu, ccip_router.TokenPriceUpdate{
			SourceToken: solana.MPK(string(update.TokenID)),
			UsdPerToken: [28]uint8(b),
		})
	}

	gpu := make([]ccip_router.GasPriceUpdate, 0, len(report.PriceUpdates.GasPriceUpdates))
	for _, update := range report.PriceUpdates.GasPriceUpdates {
		b, err := update.GasPrice.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("error marshaling gas price update: %w", err)
		}

		if len(b) > 28 {
			return nil, errors.New("error marshaling gas price update: gas price is too large")
		}

		gpu = append(gpu, ccip_router.GasPriceUpdate{
			DestChainSelector: uint64(update.ChainSel),
			UsdPerUnitGas:     [28]uint8(b),
		})
	}

	commit := ccip_router.CommitInput{
		MerkleRoot: mr,
		PriceUpdates: ccip_router.PriceUpdates{
			TokenPriceUpdates: tpu,
			GasPriceUpdates:   gpu,
		},
	}

	err := commit.MarshalWithEncoder(encoder)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *CommitPluginCodecV1) Decode(ctx context.Context, bytes []byte) (cciptypes.CommitPluginReport, error) {
	decoder := agbinary.NewBorshDecoder(bytes)
	commitReport := ccip_router.CommitInput{}
	err := commitReport.UnmarshalWithDecoder(decoder)
	if err != nil {
		return cciptypes.CommitPluginReport{}, err
	}

	merkleRoots := []cciptypes.MerkleRootChain{
		{
			ChainSel:      cciptypes.ChainSelector(commitReport.MerkleRoot.SourceChainSelector),
			OnRampAddress: commitReport.MerkleRoot.OnRampAddress,
			SeqNumsRange: cciptypes.NewSeqNumRange(
				cciptypes.SeqNum(commitReport.MerkleRoot.MinSeqNr),
				cciptypes.SeqNum(commitReport.MerkleRoot.MaxSeqNr),
			),
			MerkleRoot: commitReport.MerkleRoot.MerkleRoot,
		},
	}

	tokenPriceUpdates := make([]cciptypes.TokenPrice, 0, len(commitReport.PriceUpdates.TokenPriceUpdates))
	for _, update := range commitReport.PriceUpdates.TokenPriceUpdates {
		price := cciptypes.BigInt{}
		err = price.UnmarshalJSON(update.UsdPerToken[:])
		if err != nil {
			return cciptypes.CommitPluginReport{}, err
		}
		tokenPriceUpdates = append(tokenPriceUpdates, cciptypes.TokenPrice{
			TokenID: cciptypes.UnknownEncodedAddress(update.SourceToken.String()),
			Price:   price,
		})
	}

	gasPriceUpdates := make([]cciptypes.GasPriceChain, 0, len(commitReport.PriceUpdates.GasPriceUpdates))
	for _, update := range commitReport.PriceUpdates.GasPriceUpdates {
		price := cciptypes.BigInt{}
		err = price.UnmarshalJSON(update.UsdPerUnitGas[:])
		if err != nil {
			return cciptypes.CommitPluginReport{}, err
		}
		gasPriceUpdates = append(gasPriceUpdates, cciptypes.GasPriceChain{
			GasPrice: price,
			ChainSel: cciptypes.ChainSelector(update.DestChainSelector),
		})
	}

	return cciptypes.CommitPluginReport{
		MerkleRoots: merkleRoots,
		PriceUpdates: cciptypes.PriceUpdates{
			TokenPriceUpdates: tokenPriceUpdates,
			GasPriceUpdates:   gasPriceUpdates,
		},
	}, nil
}

// Ensure CommitPluginCodec implements the CommitPluginCodec interface
var _ cciptypes.CommitPluginCodec = (*CommitPluginCodecV1)(nil)
