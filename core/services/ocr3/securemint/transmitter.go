package securemint

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/por_mock_ocr3plugin/por"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	coretypes "github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr3/securemint/config"
)

const (
	defaultCapabilityName        = "securemint-trigger"
	defaultCapabilityVersion     = "1.0.0"
	defaultTickerResolutionMs    = 1000
	defaultSendChannelBufferSize = 1000
)

type Transmitter interface {
	ocr3types.ContractTransmitter[por.ChainSelector]
	services.Service
}

type TransmitterConfig struct {
	Logger               logger.Logger                  `json:"-"`
	CapabilitiesRegistry coretypes.CapabilitiesRegistry `json:"-"`
	DonID                uint32                         `json:"-"`

	TriggerCapabilityName        string `json:"triggerCapabilityName"`
	TriggerCapabilityVersion     string `json:"triggerCapabilityVersion"`
	TriggerTickerMinResolutionMs int    `json:"triggerTickerMinResolutionMs"`
	TriggerSendChannelBufferSize int    `json:"triggerSendChannelBufferSize"`
}

var _ Transmitter = &transmitter{}
var _ capabilities.TriggerCapability = &transmitter{}

type transmitter struct {
	services.Service
	eng *services.Engine
	capabilities.CapabilityInfo

	config      TransmitterConfig
	fromAccount ocr2types.Account
	registry    coretypes.CapabilitiesRegistry

	subscribers map[string]*subscriber
	mu          sync.Mutex
}

type subscriber struct {
	ch         chan<- capabilities.TriggerResponse
	workflowID string
	config     config.SecureMintTriggerConfig
}

func (c TransmitterConfig) NewTransmitter() (*transmitter, error) {
	return c.newTransmitter(c.Logger)
}

func (c TransmitterConfig) newTransmitter(lggr logger.Logger) (*transmitter, error) {
	t := &transmitter{
		config:      c,
		fromAccount: ocr2types.Account(lggr.Name() + strconv.FormatUint(uint64(c.DonID), 10)),
		registry:    c.CapabilitiesRegistry,
		subscribers: make(map[string]*subscriber),
	}
	if t.config.TriggerCapabilityName == "" {
		t.config.TriggerCapabilityName = defaultCapabilityName
	}
	if t.config.TriggerCapabilityVersion == "" {
		t.config.TriggerCapabilityVersion = defaultCapabilityVersion
	}
	if t.config.TriggerTickerMinResolutionMs == 0 {
		t.config.TriggerTickerMinResolutionMs = defaultTickerResolutionMs
	}
	if t.config.TriggerSendChannelBufferSize == 0 {
		t.config.TriggerSendChannelBufferSize = defaultSendChannelBufferSize
	}

	capInfo, err := capabilities.NewCapabilityInfo(
		// TODO(CAPPL-645): add labels
		t.config.TriggerCapabilityName+"@"+t.config.TriggerCapabilityVersion,
		capabilities.CapabilityTypeTrigger,
		"Secure Mint Trigger",
	)
	if err != nil {
		return nil, err
	}
	t.CapabilityInfo = capInfo

	t.Service, t.eng = services.Config{
		Name:  "SecureMintTransmitter",
		Start: t.start,
		Close: t.close,
	}.NewServiceEngine(lggr)

	return t, nil
}

func (t *transmitter) start(ctx context.Context) error {
	return t.registry.Add(ctx, t)
}

func (t *transmitter) close() error {
	return t.registry.Remove(context.Background(), t.CapabilityInfo.ID)
}

func (t *transmitter) FromAccount(context.Context) (ocr2types.Account, error) {
	return t.fromAccount, nil
}

func (t *transmitter) Transmit(
	ctx context.Context,
	cd ocr2types.ConfigDigest,
	seqNr uint64,
	report ocr3types.ReportWithInfo[por.ChainSelector],
	sigs []types.AttributedOnchainSignature,
) error {
	// Process the secure mint report and convert it to a trigger event
	capSigs := make([]capabilities.OCRAttributedOnchainSignature, len(sigs))
	for i, sig := range sigs {
		capSigs[i] = capabilities.OCRAttributedOnchainSignature{
			Signer:    uint32(sig.Signer),
			Signature: sig.Signature,
		}
	}
	ev := &capabilities.OCRTriggerEvent{
		ConfigDigest: cd[:],
		SeqNr:        seqNr,
		Report:       report.Report,
		Sigs:         capSigs,
	}
	return t.processNewEvent(ctx, ev)
}

func (t *transmitter) processNewEvent(ctx context.Context, event *capabilities.OCRTriggerEvent) error {
	// unmarshal signed report to extract timestamp and eventID
	p := &capabilitiespb.OCRTriggerReport{}
	err := proto.Unmarshal(event.Report, p)
	if err != nil {
		return fmt.Errorf("failed to unmarshal OCRTriggerReport: %w", err)
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	o, err := event.ToMap()
	if err != nil {
		return fmt.Errorf("failed to convert OCRTriggerEvent to map: %w", err)
	}
	capResponse := capabilities.TriggerResponse{
		Event: capabilities.TriggerEvent{
			TriggerType: t.CapabilityInfo.ID,
			ID:          p.EventID,
			Outputs:     o,
		},
	}

	t.eng.Debugw("ProcessReport pushing event", "eventID", p.EventID)
	nIncludedSubscribers := 0
	for _, sub := range t.subscribers {
		// include this subscriber (no frequency limiting as requested)
		select {
		case sub.ch <- capResponse:
		case <-ctx.Done():
			t.eng.Error("context done, dropping event")
			return ctx.Err()
		default:
			// drop event if channel is full - processNewEvent() should be non-blocking
			t.eng.Errorw("subscriber channel full, dropping event", "eventID", p.EventID, "workflowID", sub.workflowID)
		}
		nIncludedSubscribers++
	}
	t.eng.Debugw("ProcessReport done", "eventID", p.EventID, "nIncludedSubscribers", nIncludedSubscribers)
	return nil
}

func (t *transmitter) RegisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	config, err := validateConfig(req.Config, &t.config)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	if _, ok := t.subscribers[req.TriggerID]; ok {
		return nil, fmt.Errorf("triggerId %s already registered", t.ID)
	}

	ch := make(chan capabilities.TriggerResponse, defaultSendChannelBufferSize)
	t.subscribers[req.TriggerID] =
		&subscriber{
			ch:         ch,
			workflowID: req.Metadata.WorkflowID,
			config:     *config,
		}
	return ch, nil
}

func validateConfig(registerConfig *values.Map, capabilityConfig *TransmitterConfig) (*config.SecureMintTriggerConfig, error) {
	cfg := &config.SecureMintTriggerConfig{}
	if err := registerConfig.UnwrapTo(cfg); err != nil {
		return nil, err
	}
	if int64(cfg.MaxFrequencyMs)%int64(capabilityConfig.TriggerTickerMinResolutionMs) != 0 { //nolint:gosec // disable G115
		return nil, fmt.Errorf("MaxFrequencyMs must be a multiple of %d", capabilityConfig.TriggerTickerMinResolutionMs)
	}
	return cfg, nil
}

func (t *transmitter) UnregisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	subscriber, ok := t.subscribers[req.TriggerID]
	if !ok {
		return fmt.Errorf("triggerId %s not registered", t.ID)
	}
	close(subscriber.ch)
	delete(t.subscribers, req.TriggerID)
	return nil
}
