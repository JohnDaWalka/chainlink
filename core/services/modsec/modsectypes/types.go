package modsectypes

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/verifier_events"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/verifier_proxy"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	"github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
)

var (
	verifierEventsABI = abihelpers.MustParseABI(verifier_events.VerifierEventsABI)
)

type evmTransmitter struct {
	lggr           logger.Logger
	txm            txmgr.TxManager
	offRampAddress common.Address
}

func NewEVMTransmitter(lggr logger.Logger, txm txmgr.TxManager, offRampAddress common.Address) Transmitter {
	return &evmTransmitter{lggr: lggr, txm: txm, offRampAddress: offRampAddress}
}

func (t *evmTransmitter) Transmit(ctx context.Context, payload TransmissionPayload) error {
	etx, err := t.txm.CreateTransaction(ctx, types.TxRequest[common.Address, common.Hash]{
		ToAddress:      t.offRampAddress,
		EncodedPayload: payload.MessageData(),
		FeeLimit:       1e6, // gas limit
	})
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	t.lggr.Infow("created transaction", "tx", etx)

	return nil
}

// TransmissionPayload is the payload that is transmitted to the destination chain.
// It contains the message data and the signatures for the message being transmitted.
// The kind of signature depends on the destination chain family.
// For example, EVM would use ECDSA, Solana would use ed25519, etc.
type TransmissionPayload interface {
	// MessageData contains the message encoded for the destination chain family.
	// If the destination chain is EVM, this would be an abi-encoded Any2EVMMessage.
	// If the destination chain is Solana, this would be a borsh-encoded Any2SolanaMessage.
	// etc.
	MessageData() []byte
	// Signatures contains the signatures for the message being transmitted.
	// The kind of signature depends on the destination chain family.
	// For example, EVM would use ECDSA, Solana would use ed25519, etc.
	Signatures() [][]byte
}

type Transmitter interface {
	Transmit(ctx context.Context, payload TransmissionPayload) error
}

type CCIPMessageSent interface {
	MessageID() [32]byte
	// Encoded returns the abi-encoded message.
	// TODO: this is evm-specific, needs to be more generic.
	Encoded() ([]byte, error)
}

type CCIPMessageParser interface {
	// TODO: EVM-specific, needs to be more generic.
	ParseCCIPMessageSent(log gethtypes.Log) (CCIPMessageSent, error)
}

type evmCCIPMessageSent struct {
	ccipMessageSent verifier_proxy.VerifierProxyCCIPMessageSent
}

func (m *evmCCIPMessageSent) MessageID() [32]byte {
	return m.ccipMessageSent.Message.Header.MessageId
}

func (m *evmCCIPMessageSent) Encoded() ([]byte, error) {
	var tokenAmounts []verifier_events.InternalAny2EVMMultiProofTokenTransfer
	for _, tokenAmount := range m.ccipMessageSent.Message.TokenAmounts {
		tokenAmounts = append(tokenAmounts, verifier_events.InternalAny2EVMMultiProofTokenTransfer{
			SourcePoolAddress: common.LeftPadBytes(tokenAmount.SourcePoolAddress.Bytes(), 32),
			DestTokenAddress:  common.BytesToAddress(tokenAmount.DestTokenAddress),
			ExtraData:         tokenAmount.ExtraData,
			Amount:            tokenAmount.Amount,
		})
	}

	var requiredVerifiers []verifier_events.InternalRequiredVerifier
	for _, requiredVerifier := range m.ccipMessageSent.Message.RequiredVerifiers {
		requiredVerifiers = append(requiredVerifiers, verifier_events.InternalRequiredVerifier{
			VerifierId: requiredVerifier.VerifierId,
			Payload:    requiredVerifier.Payload,
		})
	}

	// translate the EVM2Any to Any2EVM and abi-encode it.
	any2EVM := verifier_events.InternalAny2EVMMultiProofMessage{
		Header: verifier_events.InternalHeader{
			MessageId:           m.ccipMessageSent.Message.Header.MessageId,
			SourceChainSelector: m.ccipMessageSent.Message.Header.SourceChainSelector,
			DestChainSelector:   m.ccipMessageSent.Message.Header.DestChainSelector,
			SequenceNumber:      m.ccipMessageSent.Message.Header.SequenceNumber,
		},
		Sender:            common.LeftPadBytes(m.ccipMessageSent.Message.Sender.Bytes(), 32),
		Data:              m.ccipMessageSent.Message.Data,
		Receiver:          common.BytesToAddress(m.ccipMessageSent.Message.Receiver),
		GasLimit:          200_000, // TODO: parse from extraArgs
		TokenAmounts:      tokenAmounts,
		RequiredVerifiers: requiredVerifiers,
	}

	method, ok := verifierEventsABI.Methods["exposeAny2EVMMessage"]
	if !ok {
		return nil, fmt.Errorf("exposeAny2EVMMessage method not found")
	}

	return method.Inputs.Pack(any2EVM)
}

type evmMessageParser struct {
}

func NewEVMCCIPMessageParser() CCIPMessageParser {
	return &evmMessageParser{}
}

func (p *evmMessageParser) ParseCCIPMessageSent(log gethtypes.Log) (CCIPMessageSent, error) {
	// Don't actually need to call the contract, just need to parse the log.
	verifierOnramp, err := verifier_proxy.NewVerifierProxy(log.Address, nil)
	if err != nil {
		return nil, err
	}

	parsedLog, err := verifierOnramp.ParseCCIPMessageSent(log)
	if err != nil {
		return nil, err
	}

	return &evmCCIPMessageSent{
		ccipMessageSent: *parsedLog,
	}, nil
}
