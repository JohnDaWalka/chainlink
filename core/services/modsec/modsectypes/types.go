package modsectypes

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/verifier_events"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/verifier_proxy"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	txmgrcommon "github.com/smartcontractkit/chainlink-framework/chains/txmgr"
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
	fromAddress    common.Address
}

func NewEVMTransmitter(lggr logger.Logger, txm txmgr.TxManager, offRampAddress common.Address, fromAddress common.Address) Transmitter {
	return &evmTransmitter{lggr: lggr, txm: txm, offRampAddress: offRampAddress, fromAddress: fromAddress}
}

func (t *evmTransmitter) Transmit(ctx context.Context, payload StorageValuePayload) error {
	// pack the function call for executeMessage
	method, ok := verifierEventsABI.Methods["executeMessage"]
	if !ok {
		return fmt.Errorf("executeMessage method not found")
	}

	idempotencyKey := uuid.New().String()

	// the executeMessage function just takes a single argument which is the any2EVM message.
	// its already abi-encoded in the payload, so don't need to re-encode it again, just prepend
	// the method ID.
	calldata := append(method.ID, payload.ABIEncodedMessageData...)

	etx, err := t.txm.CreateTransaction(ctx, types.TxRequest[common.Address, common.Hash]{
		IdempotencyKey: &idempotencyKey,
		Strategy:       txmgrcommon.NewSendEveryStrategy(),
		ToAddress:      t.offRampAddress,
		EncodedPayload: calldata,
		FeeLimit:       1e6, // gas limit
		FromAddress:    t.fromAddress,
	})
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	t.lggr.Infow("created transaction", "tx", etx)

	return nil
}

type Transmitter interface {
	Transmit(ctx context.Context, payload StorageValuePayload) error
}

type CCIPMessageSent interface {
	MessageID() [32]byte
	// EncodedEVM2EVM returns the abi-encoded message.
	// TODO: this is evm-specific, needs to be more generic.
	EncodedEVM2EVM() ([]byte, error)
}

type CCIPMessageParser interface {
	// TODO: EVM-specific, needs to be more generic.
	ParseEVM2AnyMessageSent(log gethtypes.Log) (CCIPMessageSent, error)
}

type evmCCIPMessageSent struct {
	ccipMessageSent verifier_proxy.VerifierProxyCCIPMessageSent
}

func (m *evmCCIPMessageSent) MessageID() [32]byte {
	return m.ccipMessageSent.Message.Header.MessageId
}

func EVM2AnyToAny2EVM(m verifier_proxy.VerifierProxyCCIPMessageSent) (verifier_events.InternalAny2EVMMultiProofMessage, error) {
	var tokenAmounts []verifier_events.InternalAny2EVMMultiProofTokenTransfer
	for _, tokenAmount := range m.Message.TokenAmounts {
		tokenAmounts = append(tokenAmounts, verifier_events.InternalAny2EVMMultiProofTokenTransfer{
			SourcePoolAddress: common.LeftPadBytes(tokenAmount.SourcePoolAddress.Bytes(), 32),
			DestTokenAddress:  common.BytesToAddress(tokenAmount.DestTokenAddress),
			ExtraData:         tokenAmount.ExtraData,
			Amount:            tokenAmount.Amount,
		})
	}

	var requiredVerifiers []verifier_events.InternalRequiredVerifier
	for _, requiredVerifier := range m.Message.RequiredVerifiers {
		requiredVerifiers = append(requiredVerifiers, verifier_events.InternalRequiredVerifier{
			VerifierId: requiredVerifier.VerifierId,
			Payload:    requiredVerifier.Payload,
		})
	}

	// translate the EVM2Any to Any2EVM and abi-encode it.
	any2EVM := verifier_events.InternalAny2EVMMultiProofMessage{
		Header: verifier_events.InternalHeader{
			MessageId:           m.Message.Header.MessageId,
			SourceChainSelector: m.Message.Header.SourceChainSelector,
			DestChainSelector:   m.Message.Header.DestChainSelector,
			SequenceNumber:      m.Message.Header.SequenceNumber,
		},
		Sender:            common.LeftPadBytes(m.Message.Sender.Bytes(), 32),
		Data:              m.Message.Data,
		Receiver:          common.BytesToAddress(m.Message.Receiver),
		GasLimit:          200_000, // TODO: parse from extraArgs
		TokenAmounts:      tokenAmounts,
		RequiredVerifiers: requiredVerifiers,
	}

	return any2EVM, nil
}

func (m *evmCCIPMessageSent) EncodedEVM2EVM() ([]byte, error) {
	any2EVM, err := EVM2AnyToAny2EVM(m.ccipMessageSent)
	if err != nil {
		return nil, err
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

func (p *evmMessageParser) ParseEVM2AnyMessageSent(log gethtypes.Log) (CCIPMessageSent, error) {
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

type StorageValuePayload struct {
	// ABIEncodedMessageData is the abi-encoded message data of the Any2EVM message.
	ABIEncodedMessageData hexutil.Bytes `json:"abiEncodedMessageData"`
	// MessageHash is the hash of the ABIEncodedMessageData.
	MessageHash hexutil.Bytes `json:"messageHash"`
	// Signature is the signature of the MessageHash.
	Signature hexutil.Bytes `json:"signature"`
}
