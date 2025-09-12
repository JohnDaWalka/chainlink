package logevent

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/triggers/logevent/logeventcap"
)

const ID = "log-event-trigger-%s-%s@1.0.0"

const defaultSendChannelBufferSize = 1000

// Log Event Trigger Capability Input
type Input struct {
}

// Log Event Trigger Capabilities Manager
// Manages different log event triggers using an underlying triggerStore
type TriggerService struct {
	services.StateMachine
	capabilities.CapabilityInfo
	capabilities.Validator[logeventcap.Config, Input, capabilities.TriggerResponse]
	lggr           logger.Logger
	triggers       CapabilitiesStore[logEventTrigger, capabilities.TriggerResponse]
	relayer        core.Relayer
	logEventConfig Config
	stopCh         services.StopChan
}

// Common capability level config across all workflows
type Config struct {
	ChainID        string `json:"chainId"`
	Network        string `json:"network"`
	LookbackBlocks uint64 `json:"lookbakBlocks"`
	PollPeriod     uint32 `json:"pollPeriod"`
	QueryCount     uint64 `json:"queryCount"`
}

func (config Config) Version(capabilityVersion string) string {
	return fmt.Sprintf(capabilityVersion, config.Network, config.ChainID)
}

var _ capabilities.TriggerCapability = (*TriggerService)(nil)
var _ services.Service = &TriggerService{}

// Creates a new Cron Trigger Service.
// Scheduling will commence on calling .Start()
func NewTriggerService(ctx context.Context,
	lggr logger.Logger,
	relayer core.Relayer,
	logEventConfig Config) (*TriggerService, error) {
	l := logger.Named(lggr, "LogEventTriggerCapabilityService")

	logEventStore := NewCapabilitiesStore[logEventTrigger, capabilities.TriggerResponse]()

	s := &TriggerService{
		lggr:           l,
		triggers:       logEventStore,
		relayer:        relayer,
		logEventConfig: logEventConfig,
		stopCh:         make(services.StopChan),
	}
	var err error
	s.CapabilityInfo, err = s.Info(ctx)
	if err != nil {
		return s, err
	}
	s.Validator = capabilities.NewValidator[logeventcap.Config, Input, capabilities.TriggerResponse](capabilities.ValidatorArgs{Info: s.CapabilityInfo})
	return s, nil
}

func (s *TriggerService) Info(ctx context.Context) (capabilities.CapabilityInfo, error) {
	return capabilities.NewCapabilityInfo(
		s.logEventConfig.Version(ID),
		capabilities.CapabilityTypeTrigger,
		"A trigger that listens for specific contract log events and starts a workflow run.",
	)
}

// Register a new trigger
// Can register triggers before the service is actively scheduling
func (s *TriggerService) RegisterTrigger(ctx context.Context,
	req capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	if req.Config == nil {
		return nil, errors.New("config is required to register a log event trigger")
	}

	// Check if this is a legacy single-log config and convert if needed
	convertedConfig, err := s.convertLegacyConfigIfNeeded(req.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to convert legacy config: %w", err)
	}

	// Now validate the config (either original multi-log or converted)
	reqConfig, err := s.ValidateConfig(convertedConfig)
	if err != nil {
		return nil, err
	}

	// Add log event trigger with Contract details to CapabilitiesStore
	var respCh chan capabilities.TriggerResponse
	ok := s.IfNotStopped(func() {
		respCh, err = s.triggers.InsertIfNotExists(req.TriggerID, func() (*logEventTrigger, chan capabilities.TriggerResponse, error) {
			l, ch, tErr := newLogEventTrigger(ctx, s.lggr, req.Metadata.WorkflowID, reqConfig, s.logEventConfig, s.relayer)
			if tErr != nil {
				return l, ch, tErr
			}
			tErr = l.Start(ctx)
			return l, ch, tErr
		})
	})
	if !ok {
		return nil, errors.New("cannot create new trigger since LogEventTriggerCapabilityService has been stopped")
	}
	if err != nil {
		return nil, fmt.Errorf("create new trigger failed %w", err)
	}
	s.lggr.Infow("RegisterTrigger", "triggerId", req.TriggerID, "WorkflowID", req.Metadata.WorkflowID)
	return respCh, nil
}

// convertLegacyConfigIfNeeded checks if the config is legacy format and converts it to multi-log format
func (s *TriggerService) convertLegacyConfigIfNeeded(config *values.Map) (*values.Map, error) {
	// Log the raw config structure for debugging
	var configMap map[string]interface{}
	if err := config.UnwrapTo(&configMap); err == nil {
		s.lggr.Infow("LEGACY_CONFIG_DEBUG: Raw config structure", "keys", getMapKeys(configMap), "rawConfig", configMap, "contractsValue", configMap["contracts"], "contractsType", fmt.Sprintf("%T", configMap["contracts"]))
	} else {
		s.lggr.Errorw("LEGACY_CONFIG_DEBUG: Failed to unwrap config for logging", "error", err)
	}

	// First, try to unmarshal as multi-log config to see if it's already in the right format
	var multiLogConfig logeventcap.Config
	err := config.UnwrapTo(&multiLogConfig)
	if err == nil {
		// Check if this is actually a valid multi-log config (has contracts)
		if len(multiLogConfig.Contracts) > 0 {
			s.lggr.Infow("LEGACY_CONFIG_DEBUG: Config is already in multi-log format, no conversion needed",
				"contracts", len(multiLogConfig.Contracts))
			return config, nil
		} else {
			// Check if there's an explicit contracts: null
			var configMap map[string]interface{}
			if err2 := config.UnwrapTo(&configMap); err2 == nil {
				contractsVal, hasContracts := configMap["contracts"]
				if hasContracts {
					s.lggr.Errorw("LEGACY_CONFIG_DEBUG: Config has explicit contracts field set to null - this should not happen!",
						"contractsValue", contractsVal, "contractsType", fmt.Sprintf("%T", contractsVal))
					return nil, fmt.Errorf("config has explicit contracts field set to null, this indicates a malformed workflow specification")
				}
			}
			s.lggr.Infow("LEGACY_CONFIG_DEBUG: Config unmarshaled as multi-log but has no contracts, treating as legacy",
				"contracts", len(multiLogConfig.Contracts))
		}
	} else {
		s.lggr.Infow("LEGACY_CONFIG_DEBUG: Config failed to unmarshal as multi-log format", "error", err)
	}

	s.lggr.Infow("LEGACY_CONFIG_DEBUG: Proceeding to check for legacy format")

	// Multi-log config failed, check if this looks like a legacy config
	legacyConfig, err := s.extractLegacyConfig(config)
	if err != nil {
		// Neither multi-log nor legacy format
		s.lggr.Errorw("LEGACY_CONFIG_DEBUG: Failed to extract legacy config", "error", err)
		return nil, fmt.Errorf("config is neither a valid multi-log trigger config nor a legacy single-log trigger config: %w", err)
	}

	// Convert legacy config to multi-log format
	s.lggr.Infow("LEGACY_CONFIG_DEBUG: Config appears to be legacy single-log trigger format, attempting conversion",
		"legacyConfig", legacyConfig)
	convertedConfig, err := s.convertLegacyToMultiLog(legacyConfig)
	if err != nil {
		s.lggr.Errorw("LEGACY_CONFIG_DEBUG: Failed to convert legacy to multi-log", "error", err)
		return nil, fmt.Errorf("failed to convert legacy config to multi-log format: %w", err)
	}
	s.lggr.Infow("LEGACY_CONFIG_DEBUG: Successfully converted legacy to multi-log",
		"convertedConfig", convertedConfig)

	// Get the original config to preserve all top-level fields
	var originalConfigMap map[string]interface{}
	if err := config.UnwrapTo(&originalConfigMap); err != nil {
		s.lggr.Errorw("LEGACY_CONFIG_DEBUG: Failed to unwrap original config for field preservation", "error", err)
		return nil, fmt.Errorf("failed to unwrap original config: %w", err)
	}

	s.lggr.Infow("LEGACY_CONFIG_DEBUG: Preserving top-level fields",
		"originalKeys", getMapKeys(originalConfigMap))

	// Create the final config with all original fields preserved
	finalConfig := make(map[string]interface{})

	// Preserve all the original top-level fields
	for key, value := range originalConfigMap {
		// Skip the legacy trigger config key (e.g., "detectEventTriggerConfig", "dtaSettlementFailedTriggerConfig")
		if strings.HasSuffix(key, "TriggerConfig") {
			s.lggr.Infow("LEGACY_CONFIG_DEBUG: Skipping legacy trigger config key", "key", key)
			continue
		}
		// Preserve all other fields
		finalConfig[key] = value
	}

	// Add the converted contracts array
	finalConfig["contracts"] = convertedConfig.Contracts

	s.lggr.Infow("LEGACY_CONFIG_DEBUG: Created final config",
		"finalConfigKeys", getMapKeys(finalConfig),
		"contractsLength", len(convertedConfig.Contracts),
		"finalConfig", finalConfig)

	// Convert the final config to a values.Map
	convertedConfigMap, err := values.WrapMap(finalConfig)
	if err != nil {
		s.lggr.Errorw("LEGACY_CONFIG_DEBUG: Failed to wrap final config", "error", err)
		return nil, fmt.Errorf("failed to wrap converted config: %w", err)
	}

	s.lggr.Infow("LEGACY_CONFIG_DEBUG: Successfully wrapped final config",
		"wrappedConfigMap", convertedConfigMap)

	s.lggr.Infow("LEGACY_CONFIG_DEBUG: CONVERSION COMPLETE - Successfully converted legacy config to multi-log format",
		"contracts", len(convertedConfig.Contracts),
		"contractNames", s.getContractNames(convertedConfig),
		"finalConfigKeys", getMapKeys(finalConfig),
		"hasContracts", finalConfig["contracts"] != nil)

	return convertedConfigMap, nil
}

// extractLegacyConfig looks for any legacy trigger config in the values.Map
func (s *TriggerService) extractLegacyConfig(config *values.Map) (*legacyTriggerConfig, error) {
	// Convert the values.Map to a generic map to look for any trigger config
	var configMap map[string]interface{}
	if err := config.UnwrapTo(&configMap); err != nil {
		return nil, fmt.Errorf("failed to unwrap config: %w", err)
	}

	// First, check if the top-level config itself is a legacy config
	if s.isValidLegacyTriggerConfig(configMap) {
		s.lggr.Infow("LEGACY_CONFIG_DEBUG: Found legacy config at top level")
		return s.parseLegacyTriggerConfig(configMap)
	}

	// Look for any key that ends with "TriggerConfig" and contains the expected structure
	for key, value := range configMap {
		if strings.HasSuffix(key, "TriggerConfig") {
			// Found a potential trigger config, check if it has the right structure
			if triggerConfig, ok := value.(map[string]interface{}); ok {
				if s.isValidLegacyTriggerConfig(triggerConfig) {
					s.lggr.Infow("LEGACY_CONFIG_DEBUG: Found legacy config in nested key", "key", key)
					return s.parseLegacyTriggerConfig(triggerConfig)
				}
			}
		}
	}

	return nil, fmt.Errorf("no valid legacy trigger config found")
}

// isValidLegacyTriggerConfig checks if a map has the structure of a legacy trigger config
func (s *TriggerService) isValidLegacyTriggerConfig(config map[string]interface{}) bool {
	requiredFields := []string{"contractName", "contractAddress", "contractEventName", "contractReaderConfig"}

	for _, field := range requiredFields {
		if _, exists := config[field]; !exists {
			return false
		}
	}

	return true
}

// parseLegacyTriggerConfig converts a generic map to a structured legacy config
func (s *TriggerService) parseLegacyTriggerConfig(config map[string]interface{}) (*legacyTriggerConfig, error) {
	// Extract the basic fields
	contractName, ok := config["contractName"].(string)
	if !ok {
		return nil, fmt.Errorf("contractName must be a string")
	}

	contractAddress, ok := config["contractAddress"].(string)
	if !ok {
		return nil, fmt.Errorf("contractAddress must be a string")
	}

	contractEventName, ok := config["contractEventName"].(string)
	if !ok {
		return nil, fmt.Errorf("contractEventName must be a string")
	}

	contractReaderConfig, ok := config["contractReaderConfig"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("contractReaderConfig must be an object")
	}

	return &legacyTriggerConfig{
		ContractName:         contractName,
		ContractAddress:      contractAddress,
		ContractEventName:    contractEventName,
		ContractReaderConfig: contractReaderConfig,
	}, nil
}

// convertLegacyToMultiLog converts a legacy trigger config to the new multi-log format
func (s *TriggerService) convertLegacyToMultiLog(legacy *legacyTriggerConfig) (*logeventcap.Config, error) {
	s.lggr.Infow("LEGACY_CONFIG_DEBUG: Starting conversion of legacy config",
		"contractName", legacy.ContractName,
		"contractAddress", legacy.ContractAddress,
		"contractEventName", legacy.ContractEventName,
		"contractReaderConfig", legacy.ContractReaderConfig)

	// Create the contract reader config structure
	contractReaderConfig := logeventcap.ConfigContractsElemContractReaderConfig{
		Contracts: make(map[string]interface{}),
	}

	// Extract the actual contract configuration from the legacy contractReaderConfig
	// Legacy format: contractReaderConfig.contracts.ContractName.{configs, contractABI, etc}
	// We need: contractReaderConfig.contracts.ContractName.{configs, contractABI, etc}
	if legacyContracts, ok := legacy.ContractReaderConfig["contracts"].(map[string]interface{}); ok {
		if contractConfig, ok := legacyContracts[legacy.ContractName]; ok {
			contractReaderConfig.Contracts[legacy.ContractName] = contractConfig
		} else {
			return nil, fmt.Errorf("contract %s not found in legacy contractReaderConfig.contracts", legacy.ContractName)
		}
	} else {
		return nil, fmt.Errorf("legacy contractReaderConfig missing 'contracts' field")
	}

	s.lggr.Infow("LEGACY_CONFIG_DEBUG: Created contract reader config",
		"contractReaderConfig", contractReaderConfig)

	// Create the contract element
	contractElem := logeventcap.ConfigContractsElem{
		ContractName:         legacy.ContractName,
		ContractAddress:      legacy.ContractAddress,
		ContractEventNames:   []string{legacy.ContractEventName},
		ContractReaderConfig: contractReaderConfig,
	}

	s.lggr.Infow("LEGACY_CONFIG_DEBUG: Created contract element",
		"contractElem", contractElem)

	// Create the multi-log config
	multiLogConfig := &logeventcap.Config{
		Contracts: []logeventcap.ConfigContractsElem{contractElem},
	}

	s.lggr.Infow("LEGACY_CONFIG_DEBUG: Created final multi-log config",
		"multiLogConfig", multiLogConfig,
		"contractsCount", len(multiLogConfig.Contracts))

	return multiLogConfig, nil
}

// getContractNames is a helper function to extract contract names for logging
func (s *TriggerService) getContractNames(config *logeventcap.Config) []string {
	var names []string
	for _, contract := range config.Contracts {
		names = append(names, contract.ContractName)
	}
	return names
}

// getMapKeys is a helper function to get keys from a map for logging
func getMapKeys(m map[string]interface{}) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// legacyTriggerConfig represents the structure of a legacy single-log trigger configuration
type legacyTriggerConfig struct {
	ContractName         string                 `json:"contractName"`
	ContractAddress      string                 `json:"contractAddress"`
	ContractEventName    string                 `json:"contractEventName"`
	ContractReaderConfig map[string]interface{} `json:"contractReaderConfig"`
}

func (s *TriggerService) UnregisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) error {
	trigger, ok := s.triggers.Read(req.TriggerID)
	if !ok {
		return fmt.Errorf("triggerId %s not found", req.TriggerID)
	}
	// Close callback channel and stop log event trigger listener
	err := trigger.Close()
	if err != nil {
		return fmt.Errorf("error closing trigger %s (chainID %s): %w", req.TriggerID, s.logEventConfig.ChainID, err)
	}
	// Remove from triggers context
	s.triggers.Delete(req.TriggerID)
	s.lggr.Infow("UnregisterTrigger", "triggerId", req.TriggerID, "WorkflowID", req.Metadata.WorkflowID)
	return nil
}

// Start the service.
func (s *TriggerService) Start(ctx context.Context) error {
	return s.StartOnce("LogEventTriggerCapabilityService", func() error {
		s.lggr.Info("Starting LogEventTriggerCapabilityService")
		return nil
	})
}

// Close stops the Service.
// After this call the Service cannot be started again,
// The service will need to be re-built to start scheduling again.
func (s *TriggerService) Close() error {
	return s.StopOnce("LogEventTriggerCapabilityService", func() error {
		s.lggr.Infow("Stopping LogEventTriggerCapabilityService")
		triggers := s.triggers.ReadAll()
		return services.MultiCloser(triggers).Close()
	})
}

func (s *TriggerService) HealthReport() map[string]error {
	return map[string]error{s.Name(): s.Healthy()}
}

func (s *TriggerService) Name() string {
	return s.lggr.Name()
}
