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
	"github.com/smartcontractkit/chainlink/v2/core/services/eastatusreporter/events"
)

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
}

const (
	ServiceName        = "EAStatusReporter"
	bridgePollPageSize = 1_000
)

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

	ticker := time.NewTicker(s.config.PollingInterval())
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Create a context with polling interval timeout
			ctx, cancel := s.chStop.CtxWithTimeout(s.config.PollingInterval())
			s.pollAllBridges(ctx)
			cancel()

		case <-s.chStop:
			return
		}
	}
}

// pollAllBridges polls all registered bridges using pagination
func (s *Service) pollAllBridges(ctx context.Context) {
	var allBridges []bridges.BridgeType
	var offset = 0

	// Paginate through all bridges
	for {
		bridgeList, _, err := s.bridgeORM.BridgeTypes(ctx, offset, bridgePollPageSize)
		if err != nil {
			s.lggr.Debugw("Failed to fetch bridges", "error", err, "offset", offset)
			return
		}

		allBridges = append(allBridges, bridgeList...)

		// If we got fewer than pageSize bridges, we've reached the end
		if len(bridgeList) < bridgePollPageSize {
			break
		}

		offset += bridgePollPageSize
	}

	if len(allBridges) == 0 {
		s.lggr.Debug("No bridges configured for EA Status Reporter polling")
		return
	}

	s.lggr.Debugw("Polling EA Status Reporter for all bridges", "count", len(allBridges))

	// Poll each bridge concurrently and wait for completion
	var wg sync.WaitGroup
	for _, bridge := range allBridges {
		wg.Add(1)
		bridgeName := string(bridge.Name)
		bridgeURL := bridge.URL.String()
		go func(name, url string) {
			defer wg.Done()
			s.pollBridge(ctx, name, url)
		}(bridgeName, bridgeURL)
	}

	wg.Wait()
}

// pollBridge polls a single bridge's status endpoint
func (s *Service) pollBridge(ctx context.Context, bridgeName string, bridgeURL string) {
	s.lggr.Debugw("Polling bridge", "bridge", bridgeName, "url", bridgeURL)

	// Parse bridge URL and construct status endpoint
	parsedURL, err := url.Parse(bridgeURL)
	if err != nil {
		s.lggr.Debugw("Failed to parse bridge URL", "bridge", bridgeName, "url", bridgeURL, "error", err)
		return
	}

	// Construct status endpoint URL (bridge::8080/status)
	statusURL := &url.URL{
		Scheme: parsedURL.Scheme,
		Host:   parsedURL.Host,
		Path:   s.config.StatusPath(),
	}

	// Make HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", statusURL.String(), nil)
	if err != nil {
		s.lggr.Debugw("Failed to create request for EA Status Reporter status", "bridge", bridgeName, "url", statusURL.String(), "error", err)
		return
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.lggr.Debugw("Failed to fetch EA Status Reporter status", "bridge", bridgeName, "url", statusURL.String(), "error", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.lggr.Debugw("EA Status Reporter status endpoint returned non-200 status", "bridge", bridgeName, "url", statusURL.String(), "status", resp.StatusCode)
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
	// Convert runtime info
	runtime := &events.RuntimeInfo{
		NodeVersion:  status.Runtime.NodeVersion,
		Platform:     status.Runtime.Platform,
		Architecture: status.Runtime.Architecture,
		Hostname:     status.Runtime.Hostname,
	}

	// Convert metrics info
	metrics := &events.MetricsInfo{
		Enabled: status.Metrics.Enabled,
	}

	// Convert endpoints
	endpointsProto := make([]*events.EndpointInfo, len(status.Endpoints))
	for i, endpoint := range status.Endpoints {
		endpointsProto[i] = &events.EndpointInfo{
			Name:       endpoint.Name,
			Aliases:    endpoint.Aliases,
			Transports: endpoint.Transports,
		}
	}

	// Convert configuration including values
	configProto := make([]*events.ConfigurationItem, len(status.Configuration))
	for i, config := range status.Configuration {
		configProto[i] = &events.ConfigurationItem{
			Name:               config.Name,
			Value:              fmt.Sprintf("%v", config.Value),
			Type:               config.Type,
			Description:        config.Description,
			Required:           config.Required,
			DefaultValue:       fmt.Sprintf("%v", config.Default),
			CustomSetting:      config.CustomSetting,
			EnvDefaultOverride: fmt.Sprintf("%v", config.EnvDefaultOverride),
		}
	}

	// Create the protobuf event
	event := &events.EAStatusEvent{
		BridgeName:           bridgeName,
		AdapterName:          status.Adapter.Name,
		AdapterVersion:       status.Adapter.Version,
		AdapterUptimeSeconds: status.Adapter.UptimeSeconds,
		DefaultEndpoint:      status.DefaultEndpoint,
		Runtime:              runtime,
		Metrics:              metrics,
		Endpoints:            endpointsProto,
		Configuration:        configProto,
	}

	// Emit the protobuf event through the configured emitter
	if err := events.EmitEAStatusEvent(ctx, s.emitter, event); err != nil {
		s.lggr.Warnw("Failed to emit EA Status Reporter protobuf data to Beholder", "bridge", bridgeName, "error", err)
		return
	}

	s.lggr.Debugw("Successfully emitted EA Status Reporter protobuf data to Beholder",
		"bridge", bridgeName,
		"adapter", status.Adapter.Name,
		"version", status.Adapter.Version,
	)
}
