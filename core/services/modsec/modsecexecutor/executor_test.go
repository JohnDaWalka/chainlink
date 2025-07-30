package modsecexecutor

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/verifier_events"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/verifier_proxy"
	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/modsec/modsecstorage"
	"github.com/smartcontractkit/chainlink/v2/core/services/modsec/modsectypes"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/test-go/testify/require"
)

const (
	sourceChainSelector = 1
	destChainSelector   = 2
)

var (
	verifierEventsABI = abihelpers.MustParseABI(verifier_events.VerifierEventsABI)
)

type simTransmitter struct {
	simBackend     *simulated.Backend
	verifierEvents *verifier_events.VerifierEvents
	auth           *bind.TransactOpts
}

func (t *simTransmitter) Transmit(ctx context.Context, payload modsectypes.StorageValuePayload) error {
	method, ok := verifierEventsABI.Methods["executeMessage"]
	if !ok {
		return fmt.Errorf("executeMessage method not found")
	}

	// decode the payload into the struct, we have to call the func directly in this case.
	ifaces, err := method.Inputs.Unpack(payload.ABIEncodedMessageData)
	if err != nil {
		return fmt.Errorf("failed to unpack abi-encoded data: %w", err)
	}

	if len(ifaces) != 1 {
		return fmt.Errorf("expected 1 interface, got %d", len(ifaces))
	}

	any2EVM := *abi.ConvertType(ifaces[0], new(verifier_events.InternalAny2EVMMultiProofMessage)).(*verifier_events.InternalAny2EVMMultiProofMessage)

	_, err = t.verifierEvents.ExecuteMessage(t.auth, any2EVM)
	if err != nil {
		return fmt.Errorf("failed to execute message: %w", err)
	}

	t.simBackend.Commit()

	return nil
}

func setupSimulatedBackendAndAuth(t *testing.T) (*simulated.Backend, client.Client, *bind.TransactOpts, *verifier_events.VerifierEvents) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	blnc, ok := big.NewInt(0).SetString("999999999999999999999999999999999999", 10)
	require.True(t, ok)

	alloc := map[common.Address]ethtypes.Account{crypto.PubkeyToAddress(privateKey.PublicKey): {Balance: blnc}}
	simulatedBackend := simulated.NewBackend(alloc, simulated.WithBlockGasLimit(8000000))

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	require.NoError(t, err)
	auth.GasLimit = uint64(6000000)

	simClient := client.NewSimulatedBackendClient(t, simulatedBackend, big.NewInt(1337))

	// adjust time so the queries by the log poller make sense.
	head, err := simClient.HeadByNumber(t.Context(), nil)
	require.NoError(t, err)
	simulatedBackend.AdjustTime(time.Since(head.Timestamp) - time.Hour)
	simulatedBackend.Commit()

	addr, _, _, err := verifier_events.DeployVerifierEvents(auth, simClient)
	require.NoError(t, err)
	simulatedBackend.Commit()

	verifierEvents, err := verifier_events.NewVerifierEvents(addr, simClient)
	require.NoError(t, err)

	return simulatedBackend, simClient, auth, verifierEvents
}

func verifierEventsToVerifierProxyCCIPMessageSent(input verifier_events.InternalEVM2AnyCommitVerifierMessage) verifier_proxy.VerifierProxyCCIPMessageSent {
	var tokenAmounts []verifier_proxy.InternalEVMTokenTransfer
	for _, tokenAmount := range input.TokenAmounts {
		tokenAmounts = append(tokenAmounts, verifier_proxy.InternalEVMTokenTransfer{
			SourcePoolAddress: tokenAmount.SourcePoolAddress,
			DestTokenAddress:  tokenAmount.DestTokenAddress,
			ExtraData:         tokenAmount.ExtraData,
			Amount:            tokenAmount.Amount,
		})
	}

	var requiredVerifiers []verifier_proxy.InternalRequiredVerifier
	for _, requiredVerifier := range input.RequiredVerifiers {
		requiredVerifiers = append(requiredVerifiers, verifier_proxy.InternalRequiredVerifier{
			VerifierId: requiredVerifier.VerifierId,
			Payload:    requiredVerifier.Payload,
		})
	}

	output := verifier_proxy.VerifierProxyCCIPMessageSent{
		Message: verifier_proxy.InternalEVM2AnyCommitVerifierMessage{
			Header: verifier_proxy.InternalHeader{
				MessageId:           input.Header.MessageId,
				SourceChainSelector: input.Header.SourceChainSelector,
				DestChainSelector:   input.Header.DestChainSelector,
				SequenceNumber:      input.Header.SequenceNumber,
			},
			Sender:             input.Sender,
			Data:               input.Data,
			Receiver:           input.Receiver,
			DestChainExtraArgs: input.DestChainExtraArgs,
			VerifierExtraArgs:  input.VerifierExtraArgs,
			FeeToken:           input.FeeToken,
			FeeTokenAmount:     input.FeeTokenAmount,
			FeeValueJuels:      input.FeeValueJuels,
			TokenAmounts:       tokenAmounts,
			RequiredVerifiers:  requiredVerifiers,
		},
	}

	return output
}

func Test_loop(t *testing.T) {
	const numMessages = 10

	simBackend, _, auth, verifierEvents := setupSimulatedBackendAndAuth(t)
	defer simBackend.Close()

	transmitter := &simTransmitter{
		simBackend:     simBackend,
		verifierEvents: verifierEvents,
		auth:           auth,
	}

	storage := modsecstorage.NewTestStorage()

	executor := New(
		logger.TestLogger(t),
		transmitter,
		verifierEvents.Address().Hex(),
		verifierEvents,
		storage,
	)

	// set some messages in the storage so we can execute them.
	for seqNr := range numMessages {
		messageID := [32]byte{byte(seqNr)}

		msg := verifier_events.InternalEVM2AnyCommitVerifierMessage{
			Header: verifier_events.InternalHeader{
				MessageId:           messageID,
				SourceChainSelector: sourceChainSelector,
				DestChainSelector:   destChainSelector,
				SequenceNumber:      uint64(seqNr),
			},
			Sender:             common.HexToAddress("0x1234"),
			Data:               nil,
			Receiver:           common.LeftPadBytes(common.HexToAddress("0x5678").Bytes(), 32),
			DestChainExtraArgs: nil,
			VerifierExtraArgs:  nil,
			FeeToken:           common.HexToAddress("0x0"),
			FeeTokenAmount:     big.NewInt(13371337),
			FeeValueJuels:      big.NewInt(13371337),
			TokenAmounts:       nil,
			RequiredVerifiers:  nil,
		}
		any2EVM, err := modsectypes.EVM2AnyToAny2EVM(verifierEventsToVerifierProxyCCIPMessageSent(msg))
		require.NoError(t, err)
		abiEncoded, err := verifierEventsABI.Methods["executeMessage"].Inputs.Pack(any2EVM)
		require.NoError(t, err)

		payload := modsectypes.StorageValuePayload{
			MessageHash:           crypto.Keccak256(abiEncoded),
			Signature:             fmt.Appendf(nil, "mock signature of message %d", seqNr),
			ABIEncodedMessageData: abiEncoded,
		}
		jsonPayload, err := json.Marshal(payload)
		require.NoError(t, err)

		storage.Set(t.Context(), hexutil.Encode(messageID[:]), jsonPayload)
	}

	require.NoError(t, executor.loop(t.Context()))

	// check execution states
	numExecuted, err := verifierEvents.SNumMessagesExecuted(&bind.CallOpts{Context: t.Context()})
	require.NoError(t, err)
	require.Equal(t, numMessages, int(numExecuted))

	// running the loop again should do nothing
	require.NoError(t, executor.loop(t.Context()))

	// check re-executions, should be 0.
	numReExecuted, err := verifierEvents.SNumMessagesReExecuted(&bind.CallOpts{Context: t.Context()})
	require.NoError(t, err)
	require.Equal(t, 0, int(numReExecuted))
}
