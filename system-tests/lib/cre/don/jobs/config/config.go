package config

import (
	"maps"
	"strconv"

	"dario.cat/mergo"
	"github.com/pkg/errors"
)

// BuildFromTOML builds configuration from TOML config, requiring global config section
// Applies in order: global config (required) -> chain-specific config (optional)
func BuildFromTOML(config map[string]any, chainID int) (map[string]any, error) {
	result := make(map[string]any)

	// Require global capability config to be present
	globalConfig, ok := config["config"]
	if !ok {
		return nil, errors.New("global config section is required but not found")
	}

	globalConfigMap, ok := globalConfig.(map[string]any)
	if !ok {
		return nil, errors.New("global config section must be a map")
	}

	// Start with global config
	if err := mergo.Merge(&result, globalConfigMap, mergo.WithOverride); err != nil {
		return nil, errors.Wrap(err, "failed to merge global config")
	}

	// Then, apply chain-specific config if it exists
	if chainConfigs, ok := config["chain_configs"]; ok {
		if chainConfigsMap, ok := chainConfigs.(map[string]any); ok {
			chainIDStr := strconv.Itoa(chainID)
			if chainConfig, ok := chainConfigsMap[chainIDStr]; ok {
				if chainConfigMap, ok := chainConfig.(map[string]any); ok {
					if err := mergo.Merge(&result, chainConfigMap, mergo.WithOverride); err != nil {
						return nil, errors.Wrap(err, "failed to merge chain-specific config")
					}
				}
			}
		}
	}

	return result, nil
}

// BuildFromTOMLOptional builds configuration from TOML config where global config is optional
// Used for capabilities like cron that don't require config by default
func BuildFromTOMLOptional(config map[string]any) (map[string]any, error) {
	result := make(map[string]any)

	// Global config is optional
	if globalConfig, ok := config["config"]; ok {
		if globalConfigMap, ok := globalConfig.(map[string]any); ok {
			if err := mergo.Merge(&result, globalConfigMap, mergo.WithOverride); err != nil {
				return nil, errors.Wrap(err, "failed to merge global config")
			}
		}
	}

	return result, nil
}

// ApplyRuntimeFallbacks fills in any missing config values with runtime-generated fallbacks
func ApplyRuntimeFallbacks(userConfig map[string]any, runtimeFallbacks map[string]any) map[string]any {
	result := make(map[string]any)
	maps.Copy(result, userConfig)

	// Add runtime fallbacks only for keys not already specified by user
	for key, value := range runtimeFallbacks {
		if _, exists := result[key]; !exists {
			result[key] = value
		}
	}

	return result
}
