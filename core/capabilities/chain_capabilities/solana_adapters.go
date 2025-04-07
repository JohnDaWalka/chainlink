package chain_capabilities

import (
	"context"
	"fmt"
	"math/big"
	"reflect"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/plugintypes"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
)

func readTokenPricesSolana(reader commontypes.SolanaChainReader, address string, _ primitives.ConfidenceLevel, params, returnVal any) error {
	tokens, ok := params.(plugintypes.GetFeeQuoterTokenUpdatesParamsType)
	if !ok {
		return fmt.Errorf("read token prices params type should be %q, but got %q", reflect.TypeOf(plugintypes.GetFeeQuoterTokenUpdatesParamsType{}), reflect.TypeOf(params))
	}

	tokenUpdates, ok := returnVal.(plugintypes.GetFeeQuoterTokenUpdatesResponse)
	if !ok {
		return fmt.Errorf("read token prices return type should be %q, but got %q", reflect.TypeOf(plugintypes.GetFeeQuoterTokenUpdatesResponse{}), reflect.TypeOf(returnVal))
	}

	var pdaAddresses [][32]byte
	for _, token := range tokens.Tokens {
		tokenAddr := solana.PublicKeyFromBytes(token[:])
		if !tokenAddr.IsOnCurve() || tokenAddr.IsZero() {
			return fmt.Errorf("read token prices return invalid token address %v (off-curve or zero)", tokenAddr)
		}

		programAddress, err := solana.PublicKeyFromBase58(address)
		if err != nil {
			return fmt.Errorf("read token prices (could not parse program address %q)", address)
		}

		pdaAddress, _, err := solana.FindProgramAddress([][]byte{[]byte("fee_billing_token_config"), tokenAddr.Bytes()}, programAddress)
		if err != nil {
			return fmt.Errorf("read token prices (failed to find PDA for token %v)", tokenAddr)
		}

		pdaAddresses = append(pdaAddresses, pdaAddress)
	}

	if len(pdaAddresses) == 0 {
		return nil
	}

	// TODO wire context
	accountsRes, err := reader.GetMultipleAccountData(context.Background(), pdaAddresses)
	if err != nil {
		return err
	}

	for _, accRes := range accountsRes {
		var wrapper fee_quoter.BillingTokenConfigWrapper

		if len(accRes.Data) > 0 {
			if err = wrapper.UnmarshalWithDecoder(bin.NewBorshDecoder(accRes.Data)); err != nil {
				return err
			}
		}

		returnVal = append(tokenUpdates, plugintypes.TimestampedUnixBig{
			Timestamp: uint32(wrapper.Config.UsdPerToken.Timestamp),
			Value:     big.NewInt(0).SetBytes(wrapper.Config.UsdPerToken.Value[:])})
		timeStampedUnixBig := plugintypes.TimestampedUnixBig{}
		timeStampedUnixBig.Timestamp = uint32(wrapper.Config.UsdPerToken.Timestamp)
		timeStampedUnixBig.Value = big.NewInt(0).SetBytes(wrapper.Config.UsdPerToken.Value[:])
	}

	return nil
}
