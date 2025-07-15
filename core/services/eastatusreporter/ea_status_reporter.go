package eastatusreporter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/recovery"
)

// EAStatusResponse represents the response schema from EA status endpoint
type EAStatusResponse struct {
	Adapter struct {
		Name          string `json:"name"`
		Version       string `json:"version"`
		UptimeSeconds int64  `json:"uptimeSeconds"`
	} `json:"adapter"`
	Endpoints []struct {
		Name       string   `json:"name"`
		Aliases    []string `json:"aliases"`
		Transports []string `json:"transports"`
	} `json:"endpoints"`
	Configuration []struct {
		Name               string      `json:"name"`
		Value              interface{} `json:"value"`
		Type               string      `json:"type"`
		Description        string      `json:"description"`
		Required           bool        `json:"required"`
		Default            interface{} `json:"default"`
		CustomSetting      bool        `json:"customSetting"`
		EnvDefaultOverride interface{} `json:"envDefaultOverride"`
	} `json:"configuration"`
	Runtime struct {
		NodeVersion  string `json:"nodeVersion"`
		Platform     string `json:"platform"`
		Architecture string `json:"architecture"`
		Hostname     string `json:"hostname"`
	} `json:"runtime"`
	Metrics struct {
		Enabled  bool    `json:"enabled"`
		Port     *int    `json:"port,omitempty"`
		Endpoint *string `json:"endpoint,omitempty"`
	} `json:"metrics"`
}

// Service polls EA status and pushes them to Beholder
type Service struct {
	services.StateMachine

	config     config.EAStatusReporter
	bridgeORM  bridges.ORM
	httpClient *http.Client
	emitter    custmsg.MessageEmitter
	lggr       logger.Logger

	// Service management
	chStop services.StopChan
	wg     sync.WaitGroup

	// Tracking state
	bridgeURLs map[string]string
	mu         sync.RWMutex
}

const ServiceName = "EAStatusReporter"

// NewService creates a new EA Status Reporter Service
func NewEaStatusReporter(
	config config.EAStatusReporter,
	bridgeORM bridges.ORM,
	httpClient *http.Client,
	emitter custmsg.MessageEmitter,
	lggr logger.Logger,
) *Service {
	return &Service{
		config:     config,
		bridgeORM:  bridgeORM,
		httpClient: httpClient,
		emitter:    emitter,
		lggr:       lggr.Named(ServiceName),
		chStop:     make(services.StopChan),
		bridgeURLs: make(map[string]string),
	}
}

// Start starts the EA Status Reporter Service
func (s *Service) Start(ctx context.Context) error {
	return s.StartOnce(ServiceName, func() error {
		if !s.config.Enabled() {
			s.lggr.Info("EA Status Reporter Service is disabled")
			return nil
		}

		s.lggr.Info("Starting EA Status Reporter Service")

		// Start periodic polling
		s.wg.Add(1)
		go s.pollLoop()

		return nil
	})
}

// Close stops the EA Status Reporter Service
func (s *Service) Close() error {
	return s.StopOnce(ServiceName, func() error {
		s.lggr.Info("Stopping " + ServiceName)
		close(s.chStop)
		s.wg.Wait()
		return nil
	})
}

// Name returns the service name
func (s *Service) Name() string {
	return s.lggr.Name()
}

// HealthReport returns the service health
func (s *Service) HealthReport() map[string]error {
	return map[string]error{s.Name(): s.Healthy()}
}

// pollLoop runs the main polling loop
func (s *Service) pollLoop() {
	defer s.wg.Done()

	interval := s.config.PollingInterval()
	if interval <= 0 {
		interval = 5 * time.Minute // Default to 5 minutes as specified
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	s.lggr.Infow("Started EA Status Reporter polling loop", "interval", interval)

	for {
		select {
		case <-ticker.C:
			ctx, cancel := s.chStop.CtxWithTimeout(30 * time.Second)
			s.pollAllBridges(ctx)
			cancel()

		case <-s.chStop:
			return
		}
	}
}

// refreshBridgeURLs loads all bridges from the database
func (s *Service) refreshBridgeURLs(ctx context.Context) error {
	// Get all bridges from database
	bridges, _, err := s.bridgeORM.BridgeTypes(ctx, 0, 1000)
	if err != nil {
		return fmt.Errorf("failed to fetch bridges: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Clear existing URLs
	s.bridgeURLs = make(map[string]string)

	// Add all bridges
	for _, bridge := range bridges {
		s.bridgeURLs[string(bridge.Name)] = bridge.URL.String()
	}

	s.lggr.Infow("Refreshed bridge URLs for EA Status Reporter", "count", len(bridges))
	return nil
}

// pollAllBridges polls all registered bridges
func (s *Service) pollAllBridges(ctx context.Context) {
	// Refresh bridge list on every poll to catch new/deleted bridges
	if err := s.refreshBridgeURLs(ctx); err != nil {
		s.lggr.Warnw("Failed to refresh bridge URLs", "error", err)
		return
	}

	s.mu.RLock()
	bridgeURLs := make(map[string]string)
	for name, url := range s.bridgeURLs {
		bridgeURLs[name] = url
	}
	s.mu.RUnlock()

	if len(bridgeURLs) == 0 {
		s.lggr.Debug("No bridges configured for EA Status Reporter polling")
		return
	}

	s.lggr.Debugw("Polling EA Status Reporter for all bridges", "count", len(bridgeURLs))

	// Poll each bridge concurrently
	var wg sync.WaitGroup
	for bridgeName, bridgeURL := range bridgeURLs {
		wg.Add(1)
		go recovery.WrapRecover(s.lggr, func() {
			defer wg.Done()
			s.pollBridge(ctx, bridgeName, bridgeURL)
		})
	}

	wg.Wait()
}

// pollBridge polls a single bridge's status endpoint
func (s *Service) pollBridge(ctx context.Context, bridgeName string, bridgeURL string) {
	// Parse bridge URL and construct status endpoint
	parsedURL, err := url.Parse(bridgeURL)
	if err != nil {
		s.lggr.Warnw("Failed to parse bridge URL", "bridge", bridgeName, "url", bridgeURL, "error", err)
		return
	}

	// Construct status endpoint URL (bridge::8080/status)
	statusURL := &url.URL{
		Scheme: parsedURL.Scheme,
		Host:   parsedURL.Host + ":8080",
		Path:   s.config.StatusPath(),
	}

	// Make HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", statusURL.String(), nil)
	if err != nil {
		s.lggr.Warnw("Failed to create request for EA Status Reporter status", "bridge", bridgeName, "url", statusURL.String(), "error", err)
		return
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.lggr.Warnw("Failed to fetch EA Status Reporter status", "bridge", bridgeName, "url", statusURL.String(), "error", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.lggr.Warnw("EA Status Reporter status endpoint returned non-200 status", "bridge", bridgeName, "url", statusURL.String(), "status", resp.StatusCode)
		return
	}

	// Parse response
	var status EAStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		s.lggr.Warnw("Failed to decode EA Status Reporter status response", "bridge", bridgeName, "url", statusURL.String(), "error", err)
		return
	}

	s.lggr.Debugw("Successfully fetched EA Status Reporter status", "bridge", bridgeName, "adapter", status.Adapter.Name, "version", status.Adapter.Version)

	// Emit telemetry to Beholder
	s.emitEAStatus(ctx, bridgeName, status)
}

// emitEAStatus sends EA Status Reporter data to Beholder
func (s *Service) emitEAStatus(ctx context.Context, bridgeName string, status EAStatusResponse) {
	// Create human-readable message
	message := fmt.Sprintf("EA Status - Bridge: %s, Adapter: %s, Version: %s",
		bridgeName,
		status.Adapter.Name,
		status.Adapter.Version,
	)

	// Build emitter with structured labels
	emitter := s.emitter.With(
		"bridge_name", bridgeName,
		"adapter_name", status.Adapter.Name,
		"adapter_version", status.Adapter.Version,
		"adapter_uptime_seconds", fmt.Sprintf("%d", status.Adapter.UptimeSeconds),
		"runtime_platform", status.Runtime.Platform,
		"runtime_architecture", status.Runtime.Architecture,
		"runtime_node_version", status.Runtime.NodeVersion,
		"runtime_hostname", status.Runtime.Hostname,
		"metrics_enabled", fmt.Sprintf("%t", status.Metrics.Enabled),
	)

	// Add endpoints information as structured data
	if len(status.Endpoints) > 0 {
		endpointsJSON, err := json.Marshal(status.Endpoints)
		if err == nil {
			emitter = emitter.With("endpoints", string(endpointsJSON))
		}
	}

	// Add configuration information as structured data
	if len(status.Configuration) > 0 {
		configJSON, err := json.Marshal(status.Configuration)
		if err == nil {
			emitter = emitter.With("configuration", string(configJSON))
		}
	}

	// Emit to Beholder
	if err := emitter.Emit(ctx, message); err != nil {
		s.lggr.Warnw("Failed to emit EA Status Reporter data to Beholder", "bridge", bridgeName, "error", err)
		return
	}

	s.lggr.Debugw("Successfully emitted EA Status Reporter data to Beholder",
		"bridge", bridgeName,
		"adapter", status.Adapter.Name,
		"version", status.Adapter.Version,
	)
}
