package ccipaptos

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/bcs"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/wsrpc/logger"
)

// ExecutePluginCodecV1 is a codec for encoding and decoding execute plugin reports.
// Compatible with ccip_offramp::offramp version 1.6.0
type ExecutePluginCodecV1 struct {
	extraDataCodec ccipcommon.ExtraDataCodec
}

func NewExecutePluginCodecV1(extraDataCodec ccipcommon.ExtraDataCodec) *ExecutePluginCodecV1 {
	return &ExecutePluginCodecV1{
		extraDataCodec: extraDataCodec,
	}
}

func pad32(b []byte) []byte {
	if len(b) > 32 {
		return b[len(b)-32:]
	}
	if len(b) == 32 {
		return b
	}
	out := make([]byte, 32)
	copy(out[32-len(b):], b)
	return out
}

func mustAddr32(s string) []byte {
	h := strings.TrimPrefix(s, "0x")
	if len(h)%2 == 1 {
		h = "0" + h
	}
	bb, err := hex.DecodeString(h)
	if err != nil {
		panic(fmt.Errorf("invalid hex address %q: %w", s, err))
	}
	return pad32(bb)
}

func (e *ExecutePluginCodecV1) Encode(ctx context.Context, report cciptypes.ExecutePluginReport) ([]byte, error) {
	lggr, err := logger.New()
	if err != nil {
		return nil, err
	}

	lggr.Info("ENCODING SUI REPORT: ", report.ChainReports)

	if len(report.ChainReports) == 0 {
		return nil, nil
	}
	if len(report.ChainReports) != 1 {
		return nil, fmt.Errorf("ExecutePluginCodecV1 expects exactly one ChainReport, found %d", len(report.ChainReports))
	}

	chainReport := report.ChainReports[0]
	if len(chainReport.Messages) != 1 {
		return nil, fmt.Errorf("only single report message expected, got %d", len(chainReport.Messages))
	}
	if len(chainReport.OffchainTokenData) != 1 {
		return nil, fmt.Errorf("only single group of offchain token data expected, got %d", len(chainReport.OffchainTokenData))
	}

	message := chainReport.Messages[0]

	// normalize vectors to non-nil so ULEB128(0) will be emitted
	var offchainTokenData [][]byte
	if len(chainReport.OffchainTokenData) > 0 {
		offchainTokenData = chainReport.OffchainTokenData[0]
	}
	if offchainTokenData == nil {
		offchainTokenData = [][]byte{}
	}
	if message.TokenAmounts == nil {
		message.TokenAmounts = []cciptypes.RampTokenAmount{}
	}
	if chainReport.Proofs == nil {
		chainReport.Proofs = []cciptypes.Bytes32{}
	}

	s := &bcs.Serializer{}

	// 1) source_chain_selector: u64
	s.U64(uint64(chainReport.SourceChainSelector))

	// --- Message Header ---
	// 2) message_id: fixed_vector<u8>(32)
	if len(message.Header.MessageID) != 32 {
		return nil, fmt.Errorf("invalid message ID length: expected 32, got %d", len(message.Header.MessageID))
	}
	s.FixedBytes(message.Header.MessageID[:])

	// 3) header_source_chain_selector: u64
	s.U64(uint64(message.Header.SourceChainSelector))

	// 4) dest_chain_selector: u64
	s.U64(uint64(message.Header.DestChainSelector))

	// 5) sequence_number: u64
	s.U64(uint64(message.Header.SequenceNumber))

	// 6) nonce: u64
	s.U64(message.Header.Nonce)

	// 7) sender: vector<u8>
	s.WriteBytes(message.Sender)

	// 8) data: vector<u8>
	s.WriteBytes(message.Data)

	// 9) receiver (Sui address): write raw 32 bytes
	receiverBytes := mustAddr32(message.Receiver.String())
	s.FixedBytes(receiverBytes)

	// 10) gas_limit: u256 (from ExtraArgs via codec)
	lggr.Infow("Initializing plugin config",
		"extraDataCodecType", fmt.Sprintf("%T", e.extraDataCodec),
	)

	decodedExtraArgsMap, err := e.extraDataCodec.DecodeExtraArgs(message.ExtraArgs, chainReport.SourceChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ExtraArgs: %w", err)
	}
	gasLimit, err := parseExtraDataMap(decodedExtraArgsMap)
	if err != nil {
		return nil, fmt.Errorf("failed to extract gas limit: %w", err)
	}
	s.U256(*gasLimit)
	lggr.Infow("Extracted gasLimit from ExtraArgs", "gasLimit", gasLimit.String())

	// 11) token_amounts: vector<Any2SuiTokenTransfer>
	//    Write ULEB128 length manually, then each item
	s.Uleb128(uint32(len(message.TokenAmounts)))
	for _, ta := range message.TokenAmounts {
		// 11a) source_pool_address: vector<u8>
		s.WriteBytes(ta.SourcePoolAddress)

		// 11b) dest_token_address: address (32 raw bytes)
		dst := mustAddr32(ta.DestTokenAddress.String())
		s.FixedBytes(dst)

		// 11c) dest_gas_amount: u32 (decoded from DestExecData)
		destExecDataDecodedMap, err2 := e.extraDataCodec.DecodeTokenAmountDestExecData(ta.DestExecData, chainReport.SourceChainSelector)
		if err2 != nil {
			return nil, fmt.Errorf("failed to decode DestExecData: %w", err2)
		}
		destGasAmount, err3 := extractDestGasAmountFromMap(destExecDataDecodedMap)
		if err3 != nil {
			return nil, fmt.Errorf("failed to extract dest gas amount: %w", err3)
		}
		s.U32(destGasAmount)

		// 11d) extra_data: vector<u8>
		s.WriteBytes(ta.ExtraData)

		// 11e) amount: u256
		if ta.Amount.Int == nil {
			return nil, fmt.Errorf("token amount is nil")
		}
		s.U256(*ta.Amount.Int)
	}

	// 12) offchain_token_data: vector<vector<u8>>
	s.Uleb128(uint32(len(offchainTokenData)))
	for _, od := range offchainTokenData {
		s.WriteBytes(od)
	}

	// 13) proofs: vector<fixed_vector<u8, 32>>
	s.Uleb128(uint32(len(chainReport.Proofs)))
	for _, p := range chainReport.Proofs {
		if len(p) != 32 {
			return nil, fmt.Errorf("invalid proof length: expected 32, got %d", len(p))
		}
		s.FixedBytes(p[:])
	}

	if s.Error() != nil {
		return nil, fmt.Errorf("BCS serialization failed: %w", s.Error())
	}

	out := s.ToBytes()
	lggr.Infow("SERIALIZED SUI REPORT IN BCS FORMAT",
		"bytesLen", len(out),
		"tail", tailHex(out, 32),
	)
	return out, nil
}

// tailHex just helps logs
func tailHex(b []byte, n int) string {
	if len(b) <= n {
		return fmt.Sprintf("% x", b)
	}
	return fmt.Sprintf("% x", b[len(b)-n:])
}

func (e *ExecutePluginCodecV1) Decode(ctx context.Context, encodedReport []byte) (cciptypes.ExecutePluginReport, error) {
	des := bcs.NewDeserializer(encodedReport)
	report := cciptypes.ExecutePluginReport{}
	var chainReport cciptypes.ExecutePluginReportSingleChain
	var message cciptypes.Message

	// 1. source_chain_selector: u64
	chainReport.SourceChainSelector = cciptypes.ChainSelector(des.U64())
	if des.Error() != nil {
		return report, fmt.Errorf("failed to deserialize source_chain_selector: %w", des.Error())
	}

	// --- Start Message Header ---
	// 2. message_id: fixed_vector_u8(32)
	messageIDBytes := des.ReadFixedBytes(32)
	if des.Error() != nil {
		return report, fmt.Errorf("failed to deserialize message_id: %w", des.Error())
	}
	copy(message.Header.MessageID[:], messageIDBytes)

	// 3. header_source_chain_selector: u64
	message.Header.SourceChainSelector = cciptypes.ChainSelector(des.U64())
	if des.Error() != nil {
		return report, fmt.Errorf("failed to deserialize header_source_chain_selector: %w", des.Error())
	}

	// 4. dest_chain_selector: u64
	message.Header.DestChainSelector = cciptypes.ChainSelector(des.U64())
	if des.Error() != nil {
		return report, fmt.Errorf("failed to deserialize dest_chain_selector: %w", des.Error())
	}

	// 5. sequence_number: u64
	message.Header.SequenceNumber = cciptypes.SeqNum(des.U64())
	if des.Error() != nil {
		return report, fmt.Errorf("failed to deserialize sequence_number: %w", des.Error())
	}

	// 6. nonce: u64
	message.Header.Nonce = des.U64()
	if des.Error() != nil {
		return report, fmt.Errorf("failed to deserialize nonce: %w", des.Error())
	}

	// --- End Message Header ---

	// 7. sender: vector<u8>
	message.Sender = des.ReadBytes()
	if des.Error() != nil {
		return report, fmt.Errorf("failed to deserialize sender: %w", des.Error())
	}

	// 8. data: vector<u8>
	message.Data = des.ReadBytes()
	if des.Error() != nil {
		return report, fmt.Errorf("failed to deserialize data: %w", des.Error())
	}

	// 9. receiver: address
	var receiverAddr aptos.AccountAddress
	des.Struct(&receiverAddr)
	if des.Error() != nil {
		return report, fmt.Errorf("failed to deserialize receiver: %w", des.Error())
	}
	message.Receiver = receiverAddr[:]

	// 10. gas_limit: u256
	_ = des.U256()
	if des.Error() != nil {
		return report, fmt.Errorf("failed to deserialize gas_limit: %w", des.Error())
	}

	// 11. token_amounts: vector<Any2AptosTokenTransfer>
	message.TokenAmounts = bcs.DeserializeSequenceWithFunction(des, func(des *bcs.Deserializer, item *cciptypes.RampTokenAmount) {
		// 11a. source_pool_address: vector<u8>
		item.SourcePoolAddress = des.ReadBytes()
		if des.Error() != nil {
			return // Error handled by caller
		}

		// 11b. dest_token_address: address
		var destTokenAddr aptos.AccountAddress
		des.Struct(&destTokenAddr)
		if des.Error() != nil {
			return // Error handled by caller
		}
		item.DestTokenAddress = destTokenAddr[:]

		// 11c. dest_gas_amount: u32
		destGasAmount := des.U32()
		if des.Error() != nil {
			return // Error handled by caller
		}
		// Encode dest gas amount back into DestExecData
		destData, err := bcs.SerializeU32(destGasAmount)
		if err != nil {
			des.SetError(fmt.Errorf("abi encode dest gas amount: %w", err))
			return
		}
		item.DestExecData = destData

		// 11d. extra_data: vector<u8>
		item.ExtraData = des.ReadBytes()
		if des.Error() != nil {
			return // Error handled by caller
		}

		// 11e. amount: u256
		amountU256 := des.U256()
		if des.Error() != nil {
			return // Error handled by caller
		}
		item.Amount = cciptypes.NewBigInt(&amountU256)
	})
	if des.Error() != nil {
		return report, fmt.Errorf("failed to deserialize token_amounts: %w", des.Error())
	}

	// 12. offchain_token_data: vector<vector<u8>>
	offchainTokenDataGroup := bcs.DeserializeSequenceWithFunction(des, func(des *bcs.Deserializer, item *[]byte) {
		*item = des.ReadBytes()
	})
	if des.Error() != nil {
		return report, fmt.Errorf("failed to deserialize offchain_token_data: %w", des.Error())
	}
	// Wrap it in the expected [][][]byte structure
	chainReport.OffchainTokenData = [][][]byte{offchainTokenDataGroup}

	// 13. proofs: vector<fixed_vector_u8(32)>
	proofsBytes := bcs.DeserializeSequenceWithFunction(des, func(des *bcs.Deserializer, item *[]byte) {
		*item = des.ReadFixedBytes(32)
	})
	if des.Error() != nil {
		return report, fmt.Errorf("failed to deserialize proofs: %w", des.Error())
	}
	// Convert [][]byte to [][32]byte
	chainReport.Proofs = make([]cciptypes.Bytes32, len(proofsBytes))
	for i, proofB := range proofsBytes {
		if len(proofB) != 32 {
			// This shouldn't happen if ReadFixedBytes worked correctly
			return report, fmt.Errorf("internal error: deserialized proof %d has length %d, expected 32", i, len(proofB))
		}
		copy(chainReport.Proofs[i][:], proofB)
	}

	// Check if all bytes were consumed
	if des.Remaining() > 0 {
		return report, fmt.Errorf("unexpected remaining bytes after decoding: %d", des.Remaining())
	}

	// Set empty fields
	message.Header.MsgHash = cciptypes.Bytes32{}
	message.Header.OnRamp = cciptypes.UnknownAddress{}
	message.FeeToken = cciptypes.UnknownAddress{}
	message.ExtraArgs = cciptypes.Bytes{}
	message.FeeTokenAmount = cciptypes.BigInt{}

	// Assemble the final report
	chainReport.Messages = []cciptypes.Message{message}
	// ProofFlagBits is not part of the Aptos report, initialize it empty/zero.
	chainReport.ProofFlagBits = cciptypes.NewBigInt(big.NewInt(0))
	report.ChainReports = []cciptypes.ExecutePluginReportSingleChain{chainReport}

	return report, nil
}

// Ensure ExecutePluginCodec implements the ExecutePluginCodec interface
var _ cciptypes.ExecutePluginCodec = (*ExecutePluginCodecV1)(nil)
