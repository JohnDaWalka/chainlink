package changeset

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	capocr3types "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/view"
	common_v1_0 "github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

type KeystoneChainView struct {
	CapabilityRegistry map[string]common_v1_0.CapabilityRegistryView `json:"capabilityRegistry,omitempty"`
	// OCRContracts is a map of OCR3 contract addresses to their configuration view
	OCRContracts     map[string]OCR3ConfigView                   `json:"ocrContracts,omitempty"`
	WorkflowRegistry map[string]common_v1_0.WorkflowRegistryView `json:"workflowRegistry,omitempty"`
	Forwarders       map[string][]ForwarderView                  `json:"forwarders,omitempty"`
}

type OCR3ConfigView struct {
	Signers               []string            `json:"signers"`
	Transmitters          []ocr2types.Account `json:"transmitters"`
	F                     uint8               `json:"f"`
	OnchainConfig         []byte              `json:"onchainConfig"`
	OffchainConfigVersion uint64              `json:"offchainConfigVersion"`
	OffchainConfig        OracleConfig        `json:"offchainConfig"`
}

type ForwarderView struct {
	DonID         uint32   `json:"donId"`
	ConfigVersion uint32   `json:"configVersion"`
	F             uint8    `json:"f"`
	Signers       []string `json:"signers"`
	TxHash        string   `json:"txHash,omitempty"`
	BlockNumber   uint64   `json:"blockNumber,omitempty"`
}

var (
	ErrOCR3NotConfigured      = errors.New("OCR3 not configured")
	ErrForwarderNotConfigured = errors.New("forwarder not configured")
)

// GenerateKeystoneChainView is a view of the keystone chain
// It is best-effort, logs errors and generates the views in parallel.
func GenerateKeystoneChainView(
	ctx context.Context,
	lggr logger.Logger,
	chain cldf.Chain,
	prevView KeystoneChainView,
	contracts viewContracts,
) (KeystoneChainView, error) {
	out := NewKeystoneChainView()
	var outMu sync.Mutex
	var allErrs error
	var wg sync.WaitGroup
	errCh := make(chan error, 4) // We are generating 4 views concurrently

	// Check if context is already done before starting work
	select {
	case <-ctx.Done():
		return out, ctx.Err()
	default:
		// Continue processing
	}

	if contracts.CapabilitiesRegistry != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for addr, capabilitiesRegistry := range contracts.CapabilitiesRegistry {
				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				default:
					cr := capabilitiesRegistry
					addrCopy := addr
					capRegView, err := common_v1_0.GenerateCapabilityRegistryView(cr)
					if err != nil {
						lggr.Warnf("failed to generate capability registry view for address %s: %v", addrCopy, err)
						errCh <- err
					}
					outMu.Lock()
					out.CapabilityRegistry[addrCopy.String()] = capRegView
					outMu.Unlock()
				}
			}
		}()
	}

	if contracts.OCR3 != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for addr, ocr3Cap := range contracts.OCR3 {
				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				default:
					oc := *ocr3Cap
					addrCopy := addr
					ocrView, err := GenerateOCR3ConfigView(ctx, oc)
					if err != nil {
						// don't block view on single OCR3 not being configured
						if errors.Is(err, ErrOCR3NotConfigured) {
							lggr.Warnf("ocr3 not configured for address %s", addrCopy)
						} else {
							lggr.Errorf("failed to generate OCR3 config view for address %s: %v", addrCopy, err)
							errCh <- err
						}
						continue
					}
					outMu.Lock()
					out.OCRContracts[addrCopy.String()] = ocrView
					outMu.Unlock()
				}
			}
		}()
	}

	// Process the workflow registry and print if WorkflowRegistryError errors.
	if contracts.WorkflowRegistry != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for addr, workflowRegistry := range contracts.WorkflowRegistry {
				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				default:
					wr := workflowRegistry
					addrCopy := addr
					wrView, wrErrs := common_v1_0.GenerateWorkflowRegistryView(wr)
					for _, err := range wrErrs {
						lggr.Errorf("WorkflowRegistry error for address %s: %v", addrCopy, err)
						errCh <- err
					}
					outMu.Lock()
					out.WorkflowRegistry[addrCopy.String()] = wrView
					outMu.Unlock()
				}
			}
		}()
	}

	if contracts.Forwarder != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, fwr := range contracts.Forwarder {
				fwrCopy := fwr
				fwrAddr := fwrCopy.Address().String()
				var prevViews []ForwarderView
				if prevView.Forwarders != nil {
					pv, ok := prevView.Forwarders[fwrAddr]
					if !ok {
						prevViews = []ForwarderView{}
					} else {
						prevViews = pv
					}
				} else {
					prevViews = []ForwarderView{}
				}

				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				default:
					fwrView, fwrErr := GenerateForwarderView(ctx, chain.Client, fwrCopy, prevViews)
					if fwrErr != nil {
						// don't block view on single forwarder not being configured
						switch {
						case errors.Is(fwrErr, ErrForwarderNotConfigured):
							lggr.Warnf("forwarder not configured for address %s", fwrCopy.Address())
						case errors.Is(fwrErr, context.Canceled), errors.Is(fwrErr, context.DeadlineExceeded):
							lggr.Warnf("forwarder view generation cancelled for address %s", fwrCopy.Address())
							errCh <- fwrErr
						default:
							lggr.Errorf("failed to generate forwarder view for address %s: %v", fwrCopy.Address(), fwrErr)
							errCh <- fwrErr
						}
					} else {
						outMu.Lock()
						out.Forwarders[fwrAddr] = fwrView
						outMu.Unlock()
					}
				}
			}
		}()
	}

	wg.Wait()
	close(errCh)

	var errList []error
	// Collect all errors
	for err := range errCh {
		errList = append(errList, err)
	}
	allErrs = errors.Join(errList...)

	return out, allErrs
}

func GenerateOCR3ConfigView(ctx context.Context, ocr3Cap ocr3_capability.OCR3Capability) (OCR3ConfigView, error) {
	details, err := ocr3Cap.LatestConfigDetails(nil)
	if err != nil {
		return OCR3ConfigView{}, err
	}

	blockNumber := uint64(details.BlockNumber)
	configIterator, err := ocr3Cap.FilterConfigSet(&bind.FilterOpts{
		Start:   blockNumber,
		End:     &blockNumber,
		Context: ctx,
	})
	if err != nil {
		return OCR3ConfigView{}, err
	}
	var config *ocr3_capability.OCR3CapabilityConfigSet
	for configIterator.Next() {
		// We wait for the iterator to receive an event
		if configIterator.Event == nil {
			return OCR3ConfigView{}, ErrOCR3NotConfigured
		}

		config = configIterator.Event
	}
	if config == nil {
		return OCR3ConfigView{}, ErrOCR3NotConfigured
	}

	var signers []ocr2types.OnchainPublicKey
	var readableSigners []string
	for _, s := range config.Signers {
		signers = append(signers, s)
		readableSigners = append(readableSigners, hex.EncodeToString(s))
	}
	var transmitters []ocr2types.Account
	for _, t := range config.Transmitters {
		transmitters = append(transmitters, ocr2types.Account(t.String()))
	}
	// `PublicConfigFromContractConfig` returns the `ocr2types.PublicConfig` that contains all the `OracleConfig` fields we need, including the
	// report plugin config.
	publicConfig, err := ocr3confighelper.PublicConfigFromContractConfig(true, ocr2types.ContractConfig{
		ConfigDigest:          config.ConfigDigest,
		ConfigCount:           config.ConfigCount,
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     config.F,
		OnchainConfig:         nil, // empty onChain config, currently we always use a nil onchain config when calling SetConfig
		OffchainConfigVersion: config.OffchainConfigVersion,
		OffchainConfig:        config.OffchainConfig,
	})
	if err != nil {
		return OCR3ConfigView{}, err
	}
	var cfg capocr3types.ReportingPluginConfig
	if err = proto.Unmarshal(publicConfig.ReportingPluginConfig, &cfg); err != nil {
		return OCR3ConfigView{}, err
	}
	oracleConfig := OracleConfig{
		MaxQueryLengthBytes:       cfg.MaxQueryLengthBytes,
		MaxObservationLengthBytes: cfg.MaxObservationLengthBytes,
		MaxReportLengthBytes:      cfg.MaxReportLengthBytes,
		MaxOutcomeLengthBytes:     cfg.MaxOutcomeLengthBytes,
		MaxReportCount:            cfg.MaxReportCount,
		MaxBatchSize:              cfg.MaxBatchSize,
		OutcomePruningThreshold:   cfg.OutcomePruningThreshold,
		RequestTimeout:            cfg.RequestTimeout.AsDuration(),
		UniqueReports:             true, // This is hardcoded to true in the OCR3 contract

		DeltaProgressMillis:               millisecondsToUint32(publicConfig.DeltaProgress),
		DeltaResendMillis:                 millisecondsToUint32(publicConfig.DeltaResend),
		DeltaInitialMillis:                millisecondsToUint32(publicConfig.DeltaInitial),
		DeltaRoundMillis:                  millisecondsToUint32(publicConfig.DeltaRound),
		DeltaGraceMillis:                  millisecondsToUint32(publicConfig.DeltaGrace),
		DeltaCertifiedCommitRequestMillis: millisecondsToUint32(publicConfig.DeltaCertifiedCommitRequest),
		DeltaStageMillis:                  millisecondsToUint32(publicConfig.DeltaStage),
		MaxRoundsPerEpoch:                 publicConfig.RMax,
		TransmissionSchedule:              publicConfig.S,

		MaxDurationQueryMillis:          millisecondsToUint32(publicConfig.MaxDurationQuery),
		MaxDurationObservationMillis:    millisecondsToUint32(publicConfig.MaxDurationObservation),
		MaxDurationShouldAcceptMillis:   millisecondsToUint32(publicConfig.MaxDurationShouldAcceptAttestedReport),
		MaxDurationShouldTransmitMillis: millisecondsToUint32(publicConfig.MaxDurationShouldTransmitAcceptedReport),

		MaxFaultyOracles: publicConfig.F,
	}

	return OCR3ConfigView{
		Signers:               readableSigners,
		Transmitters:          transmitters,
		F:                     config.F,
		OnchainConfig:         nil, // empty onChain config
		OffchainConfigVersion: config.OffchainConfigVersion,
		OffchainConfig:        oracleConfig,
	}, nil
}

func GenerateForwarderView(ctx context.Context, client cldf.OnchainClient, f *forwarder.KeystoneForwarder, prevViews []ForwarderView) ([]ForwarderView, error) {
	startBlock := uint64(0)

	// Track seen transactions to avoid duplicates
	seenTxs := make(map[string]struct{})
	var seenTxsMu sync.Mutex

	// Initialize with previous views
	forwarderViews := make([]ForwarderView, 0, len(prevViews))
	if len(prevViews) > 0 {
		// Sort `prevViews` by block number in ascending order, we make sure the last item has the highest block number
		sort.Slice(prevViews, func(i, j int) bool {
			return prevViews[i].BlockNumber < prevViews[j].BlockNumber
		})

		// Add previous events to the seen map and result list
		for _, prevView := range prevViews {
			if prevView.TxHash != "" {
				seenTxs[prevView.TxHash] = struct{}{}
			}
			forwarderViews = append(forwarderViews, prevView)
		}

		startBlock = prevViews[len(prevViews)-1].BlockNumber + 1
	} else {
		// If we don't have previous views, we will start from the deployment block number
		// which is stored in the forwarder's type and version labels.
		var deploymentBlock uint64
		lblPrefix := internal.DeploymentBlockLabel + ": "
		tvStr, err := f.TypeAndVersion(nil)
		if err != nil {
			return nil, fmt.Errorf("error getting TypeAndVersion for forwarder: %w", err)
		}
		tv, err := cldf.TypeAndVersionFromString(tvStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
		}
		for lbl := range tv.Labels {
			if strings.HasPrefix(lbl, lblPrefix) {
				// Extract the block number part after the prefix
				blockStr := strings.TrimPrefix(lbl, lblPrefix)
				blockNum, err := strconv.ParseUint(blockStr, 10, 64)
				if err == nil {
					deploymentBlock = blockNum
					break
				}
			}
		}

		if deploymentBlock > 0 {
			startBlock = deploymentBlock
		}
	}

	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting latest block header: %w", err)
	}
	currentBlock := header.Number.Uint64()

	// Early exit if no new blocks to process
	if startBlock > currentBlock {
		return forwarderViews, nil
	}

	// Create a context with timeout for the pagination process
	paginationCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	const pageSize = uint64(100000)
	const maxConcurrentWorkers = 5 // Limit concurrent RPC calls

	// Calculate the number of pages
	totalPages := (currentBlock - startBlock + pageSize) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}

	// Create work channels
	type pageRange struct {
		start, end uint64
	}
	workCh := make(chan pageRange, totalPages)
	resultCh := make(chan []ForwarderView, totalPages)
	errorCh := make(chan error, totalPages)

	// Fill a work channel with page ranges
	for pageStart := startBlock; pageStart <= currentBlock; pageStart += pageSize {
		pageEnd := pageStart + pageSize - 1
		if pageEnd > currentBlock {
			pageEnd = currentBlock
		}
		workCh <- pageRange{start: pageStart, end: pageEnd}
	}
	close(workCh)

	// Start workers (limited to maxConcurrentWorkers)
	var wg sync.WaitGroup
	workerCount := int(math.Min(float64(totalPages), float64(maxConcurrentWorkers)))
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for page := range workCh {
				// Check context before starting work
				if paginationCtx.Err() != nil {
					return
				}

				views, fvErr := processForwarderViewBlockRange(paginationCtx, f, page.start, page.end, seenTxs, &seenTxsMu)
				if fvErr != nil {
					errorCh <- fvErr
					continue
				}

				if len(views) > 0 {
					resultCh <- views
				}
			}
		}()
	}

	// Wait for all workers and close channels
	go func() {
		wg.Wait()
		close(resultCh)
		close(errorCh)
	}()

	// Collect results and errors
	var newViews []ForwarderView
	var partialErrors []error

	for views := range resultCh {
		newViews = append(newViews, views...)
	}

	for err := range errorCh {
		partialErrors = append(partialErrors, err)
	}

	// Add new views to result
	forwarderViews = append(forwarderViews, newViews...)

	// Return error if no data and not continuing previous views
	if len(forwarderViews) == 0 && len(prevViews) == 0 {
		return nil, ErrForwarderNotConfigured
	}

	// Sort results by block number
	sort.Slice(forwarderViews, func(i, j int) bool {
		return forwarderViews[i].BlockNumber < forwarderViews[j].BlockNumber
	})

	// Return partial results with error if applicable
	if len(partialErrors) > 0 {
		return forwarderViews, fmt.Errorf("partial results with errors: %w", errors.Join(partialErrors...))
	}

	return forwarderViews, nil
}

func processForwarderViewBlockRange(ctx context.Context, f *forwarder.KeystoneForwarder, startBlock, endBlock uint64,
	seenTxs map[string]struct{}, seenTxsMu *sync.Mutex) ([]ForwarderView, error) {
	const maxRetries = 3
	var views []ForwarderView

	// Try with retries for transient errors
	var configIterator *forwarder.KeystoneForwarderConfigSetIterator
	var iterErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(1<<attempt) * 200 * time.Millisecond):
				// Continue with retry
			}
		}

		// Check if the context is still valid
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		// Create request context with timeout for this specific page
		pageCtx, pageCancel := context.WithTimeout(ctx, 30*time.Second)

		configIterator, iterErr = f.FilterConfigSet(&bind.FilterOpts{
			Start:   startBlock,
			End:     &endBlock,
			Context: pageCtx,
		}, nil, nil)

		pageCancel() // Cancel the page context

		if iterErr == nil || errors.Is(iterErr, context.Canceled) {
			break // Success or intended cancellation
		}
	}

	if iterErr != nil {
		return nil, fmt.Errorf("error filtering blocks %d-%d: %w", startBlock, endBlock, iterErr)
	}

	// Process events
	func() {
		defer configIterator.Close()

		for configIterator.Next() {
			// Check for context cancellation
			if ctx.Err() != nil {
				return
			}

			event := configIterator.Event
			if event == nil {
				continue
			}

			txHash := event.Raw.TxHash.String()

			// Thread-safe check of seen transactions
			seenTxsMu.Lock()
			seen := false
			if _, exists := seenTxs[txHash]; exists {
				seen = true
			} else {
				seenTxs[txHash] = struct{}{}
			}
			seenTxsMu.Unlock()

			if seen {
				continue
			}

			var readableSigners []string
			for _, s := range event.Signers {
				readableSigners = append(readableSigners, s.String())
			}

			views = append(views, ForwarderView{
				DonID:         event.DonId,
				ConfigVersion: event.ConfigVersion,
				F:             event.F,
				Signers:       readableSigners,
				BlockNumber:   event.Raw.BlockNumber,
				TxHash:        txHash,
			})
		}

		if err := configIterator.Error(); err != nil {
			iterErr = fmt.Errorf("iterator error for blocks %d-%d: %w", startBlock, endBlock, err)
		}
	}()

	if iterErr != nil {
		return views, iterErr
	}

	return views, nil
}

func millisecondsToUint32(dur time.Duration) uint32 {
	ms := dur.Milliseconds()
	if ms > int64(math.MaxUint32) {
		return math.MaxUint32
	}
	//nolint:gosec // disable G115 as it is practically impossible to overflow here
	return uint32(ms)
}

func NewKeystoneChainView() KeystoneChainView {
	return KeystoneChainView{
		CapabilityRegistry: make(map[string]common_v1_0.CapabilityRegistryView),
		OCRContracts:       make(map[string]OCR3ConfigView),
		WorkflowRegistry:   make(map[string]common_v1_0.WorkflowRegistryView),
		Forwarders:         make(map[string][]ForwarderView),
	}
}

type KeystoneView struct {
	Chains map[string]KeystoneChainView `json:"chains,omitempty"`
	Nops   map[string]view.NopView      `json:"nops,omitempty"`
}

func (v KeystoneView) MarshalJSON() ([]byte, error) {
	// Alias to avoid recursive calls
	type Alias KeystoneView
	return json.MarshalIndent(&struct{ Alias }{Alias: Alias(v)}, "", " ")
}
