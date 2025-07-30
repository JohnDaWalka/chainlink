package modsecverifier

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/modsec/modsecstorage"
	"github.com/smartcontractkit/chainlink/v2/core/services/modsec/modsectypes"
	"github.com/smartcontractkit/chainlink/v2/core/utils/crypto"
)

const (
	// TODO: make configurable?
	lookback = time.Hour * 24
)

var _ job.ServiceCtx = &verifier{}

// verifier is a service that monitors the source chain for CCIPMessageSent events
// and pushes attestations to them to offchain storage.
type verifier struct {
	lggr          logger.Logger
	lp            logpoller.LogPoller
	wg            sync.WaitGroup
	runCtx        context.Context
	runCtxCancel  context.CancelFunc
	eventSig      string
	onRampAddress string
	storage       modsecstorage.Storage
	parser        modsectypes.CCIPMessageParser
	signerKey     ethkey.KeyV2

	latestHeadMu     sync.RWMutex
	latestHeadNumber uint64
}

func New(
	lggr logger.Logger,
	lp logpoller.LogPoller,
	eventSig string,
	onRampAddress string,
	storage modsecstorage.Storage,
	signerKey ethkey.KeyV2,
) *verifier {
	return &verifier{
		lggr:          lggr,
		lp:            lp,
		eventSig:      eventSig,
		onRampAddress: onRampAddress,
		storage:       storage,
		parser:        modsectypes.NewEVMCCIPMessageParser(),
		signerKey:     signerKey,
	}
}

func (v *verifier) run(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	var (
		lastProcessedBlock int64
		startingUp         = true
	)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// check if the filter is registered in the log poller, if not then register it
			if err := v.registerFilterIfNeeded(ctx); err != nil {
				v.lggr.Errorw("verifier failed to register filter in log poller, retrying later", "err", err)
				continue
			}

			if startingUp {
				var err error
				lastProcessedBlock, err = v.initializeLastProcessedBlock(ctx)
				if err != nil {
					v.lggr.Errorw("verifier failed to initialize last processed block, retrying later", "err", err)
					continue
				}
				startingUp = false
				v.lggr.Infow("verifier initialized last processed block", "lastProcessedBlock", lastProcessedBlock)
			}

			pending, err := v.pollLogs(ctx, lastProcessedBlock)
			if err != nil {
				v.lggr.Errorw("verifier failed to poll logs, retrying later", "err", err)
				continue
			}

			// process pending messages and insert any verifications into offchain storage
			v.processPendingMessages(ctx, pending)

			lastProcessedBlock, err = v.updateLastProcessedBlock(ctx, lastProcessedBlock)
			if err != nil {
				v.lggr.Errorw("verifier failed to update last processed block, continuing anyway", "err", err)
			} else {
				v.lggr.Infow("verifier updated last processed block", "lastProcessedBlock", lastProcessedBlock)
			}
		}
	}
}

type storageValuePayload struct {
	MessageData hexutil.Bytes `json:"messageData"`
	Signature   hexutil.Bytes `json:"signature"`
}

func (v *verifier) processPendingMessages(ctx context.Context, pending []logpoller.Log) {
	// sign messages and put them into offchain storage
	for _, message := range pending {
		parsedMessage, err := v.parser.ParseCCIPMessageSent(message.ToGethLog())
		if err != nil {
			v.lggr.Errorw("verifier failed to parse CCIPMessageSent, skipping message", "err", err)
			continue
		}

		// TODO: signature logic, for now just keccak256 the messageID and sign that.
		messageID := parsedMessage.MessageID()
		messageData, err := crypto.Keccak256(messageID[:])
		if err != nil {
			v.lggr.Errorw("verifier failed to keccak256 messageID, skipping message", "err", err)
			continue
		}
		signature, err := v.signerKey.Sign(messageData)
		if err != nil {
			v.lggr.Errorw("verifier failed to sign message, skipping message", "err", err)
			continue
		}

		payload := storageValuePayload{
			MessageData: messageData,
			Signature:   signature,
		}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			v.lggr.Errorw("verifier failed to marshal payload, skipping message", "err", err)
			continue
		}
		if err := v.storage.Set(ctx, hexutil.Encode(messageID[:]), payloadBytes); err != nil {
			v.lggr.Errorw("verifier failed to set payload in offchain storage, skipping message", "err", err)
			continue
		}

		v.lggr.Infow("verifier signed message and put in offchain storage",
			"messageID", hexutil.Encode(messageID[:]),
			"messageBlockNumber", message.BlockNumber,
			"messageTxHash", message.TxHash,
			"messageBlockHash", message.BlockHash,
			"signerKey", v.signerKey.EIP55Address.Hex(),
			"storagePayload", string(payloadBytes),
		)
	}
}

// pollLogs polls the log poller for new finalized CCIPMessageSent events.
func (v *verifier) pollLogs(ctx context.Context, lastProcessedBlock int64) ([]logpoller.Log, error) {
	latestBlock, err := v.lp.LatestBlock(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block: %w", err)
	}

	lggr := v.lggr.With(
		"latestBlock", latestBlock.BlockNumber,
		"latestFinalizedBlock", latestBlock.FinalizedBlockNumber,
	)
	lggr.Infow("verifier polling logs")
	defer func() {
		lggr.Infow("verifier done polling logs")
	}()

	// We don't specify confs because each request can have a different conf above
	// the minimum. So we do all conf handling in getConfirmedAt.
	logs, err := v.lp.LogsWithSigs(
		ctx,
		lastProcessedBlock,
		latestBlock.BlockNumber,
		[]common.Hash{common.HexToHash(v.eventSig)},
		common.HexToAddress(v.onRampAddress),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs: %w", err)
	}

	unverifiedMessages := v.getUnverifiedMessages(ctx, logs)
	if len(unverifiedMessages) == 0 {
		v.lggr.Debugw("verifier found no unverified messages")
		return nil, nil
	} else {
		v.lggr.Infow("verifier found unverified messages",
			"unverifiedMessages", len(unverifiedMessages))
	}

	return v.handleUnverifiedMessages(ctx, unverifiedMessages)
}

func (v *verifier) handleUnverifiedMessages(ctx context.Context, unverifiedMessages []logpoller.Log) (pending []logpoller.Log, err error) {
	latestBlock, err := v.lp.LatestBlock(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block: %w", err)
	}

	for _, message := range unverifiedMessages {
		parsedMessage, err := v.parser.ParseCCIPMessageSent(message.ToGethLog())
		if err != nil {
			v.lggr.Errorw("verifier failed to parse CCIPMessageSent, skipping message", "err", err)
			continue
		}

		// check if the message is finalized
		// TODO: should be able to handle FTF somewhere here.
		if message.BlockNumber > int64(latestBlock.FinalizedBlockNumber) {
			v.lggr.Debugw("verifier skipping unfinalized message",
				"messageBlockNumber", message.BlockNumber,
				"latestFinalizedBlock", latestBlock.FinalizedBlockNumber,
			)
			continue
		}

		messageID := parsedMessage.MessageID()
		if _, err := v.storage.Get(ctx, hexutil.Encode(messageID[:])); err == nil {
			v.lggr.Debugw("verifier skipping already verified message", "messageID", hexutil.Encode(messageID[:]))
			continue
		}

		v.lggr.Infow("verifier found unverified message",
			"messageID", hexutil.Encode(messageID[:]),
			"messageBlockNumber", message.BlockNumber,
			"messageTxHash", message.TxHash,
			"messageBlockHash", message.BlockHash,
			"latestFinalizedBlock", latestBlock.FinalizedBlockNumber,
		)

		pending = append(pending, message)
	}

	return pending, nil
}

// updateLastProcessedBlock returns the block number of the earliest as-of-yet unverified message.
// It uses offchain storage to determine whether a message has been verified already.
func (v *verifier) updateLastProcessedBlock(ctx context.Context, currLastProcessedBlock int64) (int64, error) {
	latestBlock, err := v.lp.LatestBlock(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest block: %w", err)
	}

	lggr := v.lggr.With(
		"currLastProcessedBlock", currLastProcessedBlock,
		"latestBlock", latestBlock.BlockNumber,
		"latestFinalizedBlock", latestBlock.FinalizedBlockNumber,
	)
	lggr.Infow("verifier updating last processed block")
	defer func() {
		lggr.Infow("verifier done updating last processed block")
	}()

	messages, err := v.lp.LogsWithSigs(
		ctx,
		currLastProcessedBlock,
		latestBlock.FinalizedBlockNumber,
		[]common.Hash{common.HexToHash(v.eventSig)},
		common.HexToAddress(v.onRampAddress),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get logs: %w", err)
	}

	unverifiedMessages := v.getUnverifiedMessages(ctx, messages)
	var earliestUnverifiedBlock = latestBlock.FinalizedBlockNumber
	for _, message := range unverifiedMessages {
		if message.BlockNumber < earliestUnverifiedBlock {
			earliestUnverifiedBlock = int64(message.BlockNumber)
		}
	}

	return earliestUnverifiedBlock, nil
}

// initializeLastProcessedBlock initializes the last processed block for the verifier
// in order to efficiently query for logs in the log poller.
// It uses offchain storage to determine whether a message has been verified already.
func (v *verifier) initializeLastProcessedBlock(ctx context.Context) (int64, error) {
	latestBlock, err := v.lp.LatestBlock(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest block: %w", err)
	}

	fromTimestamp := time.Now().UTC().Add(-lookback)
	lggr := v.lggr.With(
		"latestBlock", latestBlock.BlockNumber,
		"latestFinalizedBlock", latestBlock.FinalizedBlockNumber,
		"fromTimestamp", fromTimestamp,
	)
	lggr.Infow("verifier initializing last processed block")
	defer func() {
		lggr.Infow("verifier done initializing last processed block")
	}()

	// get messages from the from timestamp
	messages, err := v.lp.LogsCreatedAfter(
		ctx,
		common.HexToHash(v.eventSig),
		common.HexToAddress(v.onRampAddress),
		fromTimestamp,
		evmtypes.Finalized, // TODO: finalized only for now.
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get messages: %w", err)
	}

	// determine which messages have been verified already by querying
	// offchain storage for each message.
	// TODO: make a batch call?
	unverifiedMessages := v.getUnverifiedMessages(ctx, messages)
	// find message block of earliest unverified message
	// even if this block is > latest finalized, we use latest finalized as earliest unprocessed
	// because re-orgs can occur on any unfinalized block.
	var earliestUnverifiedBlock = latestBlock.FinalizedBlockNumber
	for _, message := range unverifiedMessages {
		if message.BlockNumber < earliestUnverifiedBlock {
			earliestUnverifiedBlock = int64(message.BlockNumber)
		}
	}

	return earliestUnverifiedBlock, nil
}

func (v *verifier) getUnverifiedMessages(ctx context.Context, messages []logpoller.Log) (unverifiedMessages []logpoller.Log) {
	for _, message := range messages {
		parsedMessage, err := v.parser.ParseCCIPMessageSent(message.ToGethLog())
		if err != nil {
			v.lggr.Errorw("verifier failed to parse CCIPMessageSent, skipping message", "err", err)
			continue
		}

		// TODO: can make a batch call here.
		messageID := parsedMessage.MessageID()
		if _, err := v.storage.Get(ctx, hexutil.Encode(messageID[:])); err == nil {
			v.lggr.Debugw("verifier skipping already verified message", "messageID", hexutil.Encode(messageID[:]))
			continue
		}

		unverifiedMessages = append(unverifiedMessages, message)
	}

	return unverifiedMessages
}

func (v *verifier) registerFilterIfNeeded(ctx context.Context) error {
	filters := v.lp.GetFilters()
	for _, filter := range filters {
		if len(filter.EventSigs) == 1 && filter.EventSigs[0] == common.HexToHash(v.eventSig) {
			return nil
		}
	}

	name := fmt.Sprintf("modsec-verifier-%s", v.onRampAddress)
	v.lggr.Infow("verifier registering filter in log poller",
		"eventSig", v.eventSig,
		"onRampAddress", v.onRampAddress,
		"name", name,
	)
	return v.lp.RegisterFilter(ctx, logpoller.Filter{
		EventSigs: []common.Hash{common.HexToHash(v.eventSig)},
		Addresses: []common.Address{common.HexToAddress(v.onRampAddress)},
		Name:      name,
	})
}

// Close implements job.ServiceCtx.
func (v *verifier) Close() error {
	v.runCtxCancel()
	v.wg.Wait()
	return nil
}

// Start implements job.ServiceCtx.
func (v *verifier) Start(context.Context) error {
	v.wg.Add(1)
	v.runCtx, v.runCtxCancel = context.WithCancel(context.Background())
	go func() {
		defer v.wg.Done()
		v.run(v.runCtx)
	}()
	return nil
}
