package v2

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/aggregation"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	sdkpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"
	billing "github.com/smartcontractkit/chainlink-protos/billing/go"
	protoevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"

	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/events"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/metering"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils/safe"
)

type Engine struct {
	services.Service
	srvcEng *services.Engine

	cfg          *EngineConfig
	lggr         logger.Logger
	loggerLabels map[string]string
	localNode    capabilities.Node

	// registration ID -> trigger capability
	triggers map[string]*triggerCapability
	// used to separate registration and unregistration phases
	triggersRegMu sync.Mutex

	allTriggerEventsQueueCh chan enqueuedTriggerEvent
	executionsSemaphore     chan struct{}
	capCallsSemaphore       *semaphore[*sdkpb.CapabilityResponse]

	meterReports *metering.Reports

	metrics *monitoring.WorkflowsMetricLabeler

	// Beholder logs are handled asynchronously.
	// To avoid losing logs, we implement a two-step graceful shutdown.
	emitQueue              chan func(context.Context)
	beholderEmitSoftStopCh services.StopChan
	beholderEmitHardStopCh services.StopChan
	beholderEmitDone       chan struct{}
}

type triggerCapability struct {
	capabilities.TriggerCapability
	payload *anypb.Any
}

type enqueuedTriggerEvent struct {
	triggerCapID string
	triggerIndex int
	timestamp    time.Time
	event        capabilities.TriggerResponse
}

func NewEngine(cfg *EngineConfig) (*Engine, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	em, err := monitoring.InitMonitoringResources()
	if err != nil {
		return nil, fmt.Errorf("could not initialize monitoring resources: %w", err)
	}

	// LocalNode() is expected to be non-blocking at this stage (i.e. the registry is already synced)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	localNode, err := cfg.CapRegistry.LocalNode(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get local node state: %w", err)
	}

	engine := &Engine{
		cfg:                     cfg,
		localNode:               localNode,
		triggers:                make(map[string]*triggerCapability),
		allTriggerEventsQueueCh: make(chan enqueuedTriggerEvent, cfg.LocalLimits.TriggerEventQueueSize),
		executionsSemaphore:     make(chan struct{}, cfg.LocalLimits.MaxConcurrentWorkflowExecutions),
		capCallsSemaphore:       NewSemaphore[*sdkpb.CapabilityResponse](cfg.LocalLimits.MaxConcurrentCapabilityCallsPerWorkflow),
		emitQueue:               make(chan func(context.Context), cfg.LocalLimits.BeholderEmitQueueSize),
		beholderEmitSoftStopCh:  make(chan struct{}),
		beholderEmitHardStopCh:  make(chan struct{}),
		beholderEmitDone:        make(chan struct{}),
	}

	labels := []any{
		platform.KeyWorkflowID, cfg.WorkflowID,
		platform.KeyWorkflowOwner, cfg.WorkflowOwner,
		platform.KeyWorkflowName, cfg.WorkflowName.String(),
		platform.KeyWorkflowVersion, platform.ValueWorkflowVersionV2,
		platform.KeyDonID, strconv.Itoa(int(localNode.WorkflowDON.ID)),
		platform.KeyDonF, strconv.Itoa(int(localNode.WorkflowDON.F)),
		platform.KeyDonN, strconv.Itoa(len(localNode.WorkflowDON.Members)),
		platform.KeyDonQ, strconv.Itoa(aggregation.ByzantineQuorum(
			len(localNode.WorkflowDON.Members),
			int(localNode.WorkflowDON.F),
		)),
		platform.KeyP2PID, localNode.PeerID.String(),
	}

	beholderLogger := custmsg.NewBeholderLogger(cfg.Lggr, cfg.BeholderEmitter, engine).Named("WorkflowEngine").With(labels...)
	metricsLabeler := monitoring.NewWorkflowsMetricLabeler(metrics.NewLabeler(), em).With(
		platform.KeyWorkflowID, cfg.WorkflowID,
		platform.KeyWorkflowOwner, cfg.WorkflowOwner,
		platform.KeyWorkflowName, cfg.WorkflowName.String())
	labelsMap := make(map[string]string, len(labels)/2)
	for i := 0; i < len(labels); i += 2 {
		labelsMap[labels[i].(string)] = labels[i+1].(string)
	}

	if cfg.DebugMode {
		beholderLogger.Errorw("WARNING: Debug mode is enabled, this is not suitable for production")
	}

	if cfg.SecretsFetcher == nil {
		cfg.SecretsFetcher = NewSecretsFetcher(
			metricsLabeler,
			cfg.CapRegistry,
			beholderLogger,
			NewSemaphore[[]*sdkpb.SecretResponse](cfg.LocalLimits.MaxConcurrentSecretsCallsPerWorkflow),
			cfg.WorkflowOwner,
			cfg.WorkflowName.String(),
			func(shares []string) (string, error) {
				return "", errors.New("decryption not implemented in v2 engine")
			},
		)
	}

	engine.lggr = beholderLogger
	engine.loggerLabels = labelsMap
	engine.meterReports = metering.NewReports(cfg.BillingClient, cfg.WorkflowOwner, cfg.WorkflowID, beholderLogger, labelsMap, metricsLabeler, cfg.WorkflowRegistryAddress, cfg.WorkflowRegistryChainSelector)
	engine.metrics = metricsLabeler

	engine.Service, engine.srvcEng = services.Config{
		Name:  "WorkflowEngineV2",
		Start: engine.start,
		Close: engine.close,
	}.NewServiceEngine(beholderLogger)
	return engine, nil
}

func (e *Engine) start(ctx context.Context) error {
	e.cfg.Module.Start()
	ctx = context.WithoutCancel(ctx)
	ctx = contexts.WithCRE(ctx, contexts.CRE{Owner: e.cfg.WorkflowOwner, Workflow: e.cfg.WorkflowID}) // TODO org?
	e.srvcEng.GoCtx(ctx, e.heartbeatLoop)
	e.srvcEng.GoCtx(ctx, e.init)
	e.srvcEng.GoCtx(ctx, e.handleAllTriggerEvents)
	// launch outside srvcEng because we want this to outlive all other goroutines
	go e.beholderAsyncEmitLoop(ctx)
	return nil
}

func (e *Engine) init(ctx context.Context) {
	// apply global engine instance limits
	// TODO(CAPPL-794): consider moving this outside of the engine, into the Syncer
	err := e.cfg.GlobalLimits.Use(ctx, 1)
	if err != nil {
		var errLimited limits.ErrorResourceLimited[int]
		if errors.As(err, &errLimited) {
			switch errLimited.Scope {
			case settings.ScopeOwner:
				e.lggr.Info("Per owner workflow count limit reached", "err", err)
				e.metrics.IncrementWorkflowLimitPerOwnerCounter(ctx)
				e.cfg.Hooks.OnInitialized(types.ErrPerOwnerWorkflowCountLimitReached)
			case settings.ScopeGlobal:
				e.lggr.Info("Global workflow count limit reached", "err", err)
				e.metrics.IncrementWorkflowLimitGlobalCounter(ctx)
				e.cfg.Hooks.OnInitialized(types.ErrGlobalWorkflowCountLimitReached)
			default:
				e.lggr.Errorw("Workflow count limit reached for unexpected scope", "scope", errLimited.Scope, "err", err)
				e.cfg.Hooks.OnInitialized(err)
			}
		} else {
			e.cfg.Hooks.OnInitialized(err)
		}
		return
	}

	err = e.runTriggerSubscriptionPhase(ctx)
	if err != nil {
		e.lggr.Errorw("Workflow Engine initialization failed", "err", err)
		e.cfg.Hooks.OnInitialized(err)
		return
	}

	e.lggr.Info("Workflow Engine initialized")
	e.metrics.IncrementWorkflowInitializationCounter(ctx)
	e.cfg.Hooks.OnInitialized(nil)
}

func (e *Engine) runTriggerSubscriptionPhase(ctx context.Context) error {
	// call into the workflow to get trigger subscriptions
	subCtx, subCancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(e.cfg.LocalLimits.TriggerSubscriptionRequestTimeoutMs))
	defer subCancel()

	result, err := e.cfg.Module.Execute(subCtx, &sdkpb.ExecuteRequest{
		Request:         &sdkpb.ExecuteRequest_Subscribe{},
		MaxResponseSize: uint64(e.cfg.LocalLimits.ModuleExecuteMaxResponseSizeBytes),
		Config:          e.cfg.WorkflowConfig,
	}, NewDisallowedExecutionHelper(e.lggr, e, TimeProvider{}, e.cfg.SecretsFetcher))
	if err != nil {
		return fmt.Errorf("failed to execute subscribe: %w", err)
	}
	if result.GetError() != "" {
		return fmt.Errorf("failed to execute subscribe: %s", result.GetError())
	}
	subs := result.GetTriggerSubscriptions()
	if subs == nil {
		return errors.New("subscribe result is nil")
	}
	if len(subs.Subscriptions) > int(e.cfg.LocalLimits.MaxTriggerSubscriptions) {
		return fmt.Errorf("too many trigger subscriptions: %d", len(subs.Subscriptions))
	}

	// check if all requested triggers exist in the registry
	triggers := make([]capabilities.TriggerCapability, 0, len(subs.Subscriptions))
	for _, sub := range subs.Subscriptions {
		triggerCap, err := e.cfg.CapRegistry.GetTrigger(ctx, sub.Id)
		if err != nil {
			return fmt.Errorf("trigger capability not found: %w", err)
		}
		triggers = append(triggers, triggerCap)
	}

	// register to all triggers
	regCtx, regCancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(e.cfg.LocalLimits.TriggerAllRegistrationsTimeoutMs))
	defer regCancel()
	e.triggersRegMu.Lock()
	defer e.triggersRegMu.Unlock()
	eventChans := make([]<-chan capabilities.TriggerResponse, len(subs.Subscriptions))
	triggerCapIDs := make([]string, len(subs.Subscriptions))
	for i, sub := range subs.Subscriptions {
		triggerCap := triggers[i]
		registrationID := fmt.Sprintf("trigger_reg_%s_%d", e.cfg.WorkflowID, i)
		e.lggr.Debugw("Registering trigger", "triggerID", sub.Id, "method", sub.Method)
		triggerEventCh, err := triggerCap.RegisterTrigger(regCtx, capabilities.TriggerRegistrationRequest{
			TriggerID: registrationID,
			Metadata: capabilities.RequestMetadata{
				WorkflowID:               e.cfg.WorkflowID,
				WorkflowOwner:            e.cfg.WorkflowOwner,
				WorkflowName:             e.cfg.WorkflowName.Hex(),
				DecodedWorkflowName:      e.cfg.WorkflowName.String(),
				WorkflowDonID:            e.localNode.WorkflowDON.ID,
				WorkflowDonConfigVersion: e.localNode.WorkflowDON.ConfigVersion,
				ReferenceID:              fmt.Sprintf("trigger_%d", i),
				// no WorkflowExecutionID needed (or available at this stage)
			},
			Payload: sub.Payload,
			Method:  sub.Method,
			// no Config needed - NoDAG uses Payload
		})
		if err != nil {
			e.lggr.Errorw("One of trigger registrations failed - reverting all", "triggerID", sub.Id, "err", err)
			e.metrics.With(platform.KeyTriggerID, sub.Id).IncrementRegisterTriggerFailureCounter(ctx)
			e.unregisterAllTriggers(ctx)
			return fmt.Errorf("failed to register trigger: %w", err)
		}
		e.triggers[registrationID] = &triggerCapability{
			TriggerCapability: triggerCap,
			payload:           sub.Payload,
		}
		eventChans[i] = triggerEventCh
		triggerCapIDs[i] = sub.Id
	}

	// start listening for trigger events only if all registrations succeeded
	for idx, triggerEventCh := range eventChans {
		e.srvcEng.Go(func(srvcCtx context.Context) {
			for {
				select {
				case <-srvcCtx.Done():
					return
				case event, isOpen := <-triggerEventCh:
					if !isOpen {
						return
					}
					if event.Err != nil {
						e.lggr.Errorw("Received a trigger event with error, dropping", "triggerID", subs.Subscriptions[idx].Id, "err", event.Err)
						e.metrics.With(platform.KeyTriggerID, subs.Subscriptions[idx].Id).IncrementWorkflowTriggerEventErrorCounter(srvcCtx)
						continue
					}
					select {
					case e.allTriggerEventsQueueCh <- enqueuedTriggerEvent{
						triggerCapID: subs.Subscriptions[idx].Id,
						triggerIndex: idx,
						timestamp:    e.cfg.Clock.Now(),
						event:        event,
					}:
					default: // queue full, drop the event
						e.lggr.Errorw("Trigger event queue is full, dropping event", "triggerID", subs.Subscriptions[idx].Id, "triggerIndex", idx)
						e.metrics.With(platform.KeyTriggerID, subs.Subscriptions[idx].Id).IncrementWorkflowTriggerEventQueueFullCounter(srvcCtx)
					}
				}
			}
		})
	}
	e.lggr.Infow("All triggers registered successfully", "numTriggers", len(subs.Subscriptions), "triggerIDs", triggerCapIDs)
	e.metrics.IncrementWorkflowRegisteredCounter(ctx)
	e.cfg.Hooks.OnSubscribedToTriggers(triggerCapIDs)
	return nil
}

func (e *Engine) handleAllTriggerEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case queueHead, isOpen := <-e.allTriggerEventsQueueCh:
			if !isOpen {
				return
			}
			eventAge := queueHead.timestamp.Sub(e.cfg.Clock.Now())
			if eventAge > time.Duration(e.cfg.LocalLimits.TriggerEventMaxAgeMs)*time.Millisecond {
				e.lggr.Warnw("Trigger event is too old, skipping execution", "triggerID", queueHead.triggerCapID, "eventID", queueHead.event.Event.ID, "eventAgeMs", eventAge.Milliseconds())
				continue
			}
			select {
			case e.executionsSemaphore <- struct{}{}: // block if too many concurrent workflow executions
				e.srvcEng.Go(func(srvcCtx context.Context) {
					e.startExecution(srvcCtx, queueHead)
					<-e.executionsSemaphore
				})
			case <-ctx.Done():
				return
			}
		}
	}
}

// startExecution initiates a new workflow execution, blocking until completed
func (e *Engine) startExecution(ctx context.Context, wrappedTriggerEvent enqueuedTriggerEvent) {
	triggerEvent := wrappedTriggerEvent.event.Event
	executionID, err := types.GenerateExecutionID(e.cfg.WorkflowID, triggerEvent.ID)
	if err != nil {
		e.lggr.Errorw("Failed to generate execution ID", "err", err, "triggerID", wrappedTriggerEvent.triggerCapID)
		return
	}

	// TODO(CAPPL-911): add rate-limiting

	meteringReport, meteringErr := e.meterReports.Start(ctx, executionID)
	if meteringErr != nil {
		e.lggr.Errorw("could start metering workflow execution. continuing without metering", "err", meteringErr)
	}

	isMetering := meteringErr == nil
	if isMetering {
		mrErr := meteringReport.Reserve(ctx)
		if mrErr != nil {
			e.lggr.Errorw("could not reserve metering", "err", mrErr)
			return
		}

		e.deductStandardBalances(meteringReport)
	}

	execCtx, execCancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(e.cfg.LocalLimits.WorkflowExecutionTimeoutMs))
	defer execCancel()
	executionLogger := logger.With(e.lggr, "executionID", executionID, "triggerID", wrappedTriggerEvent.triggerCapID, "triggerIndex", wrappedTriggerEvent.triggerIndex)

	tid, err := safe.IntToUint64(wrappedTriggerEvent.triggerIndex)
	if err != nil {
		executionLogger.Errorw("Failed to convert trigger index to uint64", "err", err)
		return
	}

	startTime := e.cfg.Clock.Now()
	executionLogger.Infow("Workflow execution starting ...")
	_ = events.EmitExecutionStartedEvent(ctx, e.loggerLabels, triggerEvent.ID, executionID)
	var executionStatus string // store.StatusStarted

	result, err := e.cfg.Module.Execute(execCtx, &sdkpb.ExecuteRequest{
		Request: &sdkpb.ExecuteRequest_Trigger{
			Trigger: &sdkpb.Trigger{
				Id:      tid,
				Payload: triggerEvent.Payload,
			},
		},
		MaxResponseSize: uint64(e.cfg.LocalLimits.ModuleExecuteMaxResponseSizeBytes),
		Config:          e.cfg.WorkflowConfig,
	}, &ExecutionHelper{Engine: e, WorkflowExecutionID: executionID, SecretsFetcher: e.cfg.SecretsFetcher})

	endTime := e.cfg.Clock.Now()
	executionDuration := endTime.Sub(startTime)

	if isMetering {
		computeUnit := billing.ResourceType_name[int32(billing.ResourceType_RESOURCE_TYPE_COMPUTE)]
		mrErr := meteringReport.Settle(computeUnit,
			[]capabilities.MeteringNodeDetail{{
				Peer2PeerID: e.localNode.PeerID.String(),
				SpendUnit:   computeUnit,
				SpendValue:  strconv.Itoa(int(executionDuration.Milliseconds())),
			}})
		if mrErr != nil {
			e.lggr.Errorw("could not set metering for compute", "err", mrErr)
		}
		mrErr = e.meterReports.End(ctx, executionID)
		if mrErr != nil {
			e.lggr.Errorw("could not end metering report", "err", mrErr)
		}
	}

	if err != nil {
		executionStatus = store.StatusErrored
		if errors.Is(err, context.DeadlineExceeded) {
			executionStatus = store.StatusTimeout
		}
		executionLogger.Errorw("Workflow execution failed", "err", err, "status", executionStatus, "durationMs", executionDuration.Milliseconds())
		_ = events.EmitExecutionFinishedEvent(ctx, e.loggerLabels, executionStatus, executionID)
		e.cfg.Hooks.OnExecutionFinished(executionID, executionStatus)
		return
	}

	if e.cfg.DebugMode {
		e.lggr.Debugw("User workflow execution result", "result", result.GetValue(), "err", result.GetError())
	}

	executionStatus = store.StatusCompleted
	executionLogger.Infow("Workflow execution finished successfully", "durationMs", executionDuration.Milliseconds())
	_ = events.EmitExecutionFinishedEvent(ctx, e.loggerLabels, executionStatus, executionID)

	e.cfg.Hooks.OnResultReceived(result)
	e.cfg.Hooks.OnExecutionFinished(executionID, executionStatus)
}

func (e *Engine) close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(e.cfg.LocalLimits.ShutdownTimeoutMs))
	defer cancel()
	ctx = contexts.WithCRE(ctx, contexts.CRE{Owner: e.cfg.WorkflowOwner, Workflow: e.cfg.WorkflowID}) // TODO org?
	e.triggersRegMu.Lock()
	e.unregisterAllTriggers(ctx)
	e.triggersRegMu.Unlock()

	e.cfg.Module.Close()

	// reset metering mode metric so that a positive value does not persist
	e.metrics.UpdateWorkflowMeteringModeGauge(ctx, false)

	// all executions are finished, we can shut down beholder emit goroutine now
	close(e.beholderEmitSoftStopCh)
	select {
	case <-e.beholderEmitDone:
	case <-ctx.Done():
		close(e.beholderEmitHardStopCh)
		<-e.beholderEmitDone
	}

	return e.cfg.GlobalLimits.Free(ctx, 1)
}

// NOTE: needs to be called under the triggersRegMu lock
func (e *Engine) unregisterAllTriggers(ctx context.Context) {
	for registrationID, trigger := range e.triggers {
		err := trigger.UnregisterTrigger(ctx, capabilities.TriggerRegistrationRequest{
			TriggerID: registrationID,
			Metadata: capabilities.RequestMetadata{
				WorkflowID:    e.cfg.WorkflowID,
				WorkflowDonID: e.localNode.WorkflowDON.ID,
			},
			Payload: trigger.payload,
		})
		if err != nil {
			e.lggr.Errorw("Failed to unregister trigger", "registrationId", registrationID, "err", err)
		}
	}
	e.triggers = make(map[string]*triggerCapability)
	e.lggr.Infow("All triggers unregistered", "numTriggers", len(e.triggers))
	e.metrics.IncrementWorkflowUnregisteredCounter(ctx)
}

func (e *Engine) heartbeatLoop(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(e.cfg.LocalLimits.HeartbeatFrequencyMs) * time.Millisecond)
	defer ticker.Stop()
	e.lggr.Info("Starting heartbeat loop")
	e.metrics.EngineHeartbeatGauge(ctx, 1)

	for {
		select {
		case <-ctx.Done():
			e.metrics.EngineHeartbeatGauge(ctx, 0)
			e.lggr.Info("Shutting down heartbeat")
			return
		case <-ticker.C:
			e.lggr.Debugw("Engine heartbeat tick", "time", e.cfg.Clock.Now().Format(time.RFC3339))
			e.metrics.IncrementEngineHeartbeatCounter(ctx)
		}
	}
}

func (e *Engine) deductStandardBalances(meteringReport *metering.Report) {
	// V2Engine runs the entirety of a module's execution as compute. Ensure that the max execution time can run.
	// Add an extra second of metering padding for context cancel propagation
	ctxCancelPadding := (time.Millisecond * 1000).Milliseconds()
	compMs := decimal.NewFromInt(int64(e.cfg.LocalLimits.WorkflowExecutionTimeoutMs) + ctxCancelPadding)
	computeUnit := billing.ResourceType_RESOURCE_TYPE_COMPUTE.String()

	if _, err := meteringReport.Deduct(
		computeUnit,
		metering.ByResource(computeUnit, compMs),
	); err != nil {
		e.lggr.Errorw("could not deduct balance for capability request", "capReq", "standard-deduction-compute", "err", err)
	}
}

func (e *Engine) enqueueUserLog(msg string, count int, workflowExecutionID string) error {
	logLine := &protoevents.LogLine{
		NodeTimestamp: time.Now().Format(time.RFC3339Nano),
		Message:       msg,
	}

	if e.cfg.DebugMode {
		e.lggr.Debugf("User log: <<<%s>>>, local node timestamp: %s", logLine.Message, logLine.NodeTimestamp)
	}
	if count > int(e.cfg.LocalLimits.MaxUserLogEventsPerExecution) {
		e.lggr.Warnw("Max user log events per execution reached, dropping event", "maxEvents", e.cfg.LocalLimits.MaxUserLogEventsPerExecution)
		return fmt.Errorf("max user log events per execution reached: %d", e.cfg.LocalLimits.MaxUserLogEventsPerExecution)
	}
	if len(logLine.Message) > int(e.cfg.LocalLimits.MaxUserLogLineLength) {
		logLine.Message = logLine.Message[:e.cfg.LocalLimits.MaxUserLogLineLength] + " ...(truncated)"
	}

	e.Enqueue(func(ctx context.Context) {
		if err := events.EmitUserLogs(ctx, e.loggerLabels, []*protoevents.LogLine{logLine}, workflowExecutionID); err != nil {
			e.lggr.Errorw("Failed to emit user logs", "err", err)
		}
	})
	return nil
}

func (e *Engine) beholderAsyncEmitLoop(ctx context.Context) {
	defer func() { e.beholderEmitDone <- struct{}{} }()
	ctx, cancel := e.beholderEmitHardStopCh.Ctx(context.WithoutCancel(ctx))
	defer cancel()
	e.lggr.Info("Starting Beholder async emit loop")
	for {
		select {
		case <-e.beholderEmitSoftStopCh:
			e.lggr.Info("Shutting down Beholder async emit loop - engine shutting down")
			// drain emit queue until hard stop
			for {
				select {
				case fn, ok := <-e.emitQueue:
					if !ok {
						return
					}
					fn(ctx)
				case <-e.beholderEmitHardStopCh:
					e.lggr.Info("Shutting down Beholder async emit loop - hard stop")
					return
				default:
					return
				}
			}
		case fn, ok := <-e.emitQueue:
			if !ok {
				// should never happen
				e.cfg.Lggr.Error("Shutting down Beholder async emit loop - queue channel closed")
				return
			}
			fn(ctx)
		}
	}
}

// Implementation of the taskProcessor interface for BeholderLoger
func (e *Engine) Enqueue(fn func(context.Context)) {
	select {
	case e.emitQueue <- fn:
	default:
		// IMPORTANT: Use the raw stderr logger, not BeholderLogger (e.lggr)
		// Consider adding a metric.
		e.cfg.Lggr.Errorw("emit queue is full, dropping log", "queueLen", len(e.emitQueue))
	}
}
