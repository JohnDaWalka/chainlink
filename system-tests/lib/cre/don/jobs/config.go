package jobs

import (
	"maps"
	"strconv"

	"dario.cat/mergo"
	"github.com/pkg/errors"
)

// BuildConfigFromTOML builds configuration from TOML config, requiring global config section
// Applies in order: global config (required) -> chain-specific config (optional)
func BuildConfigFromTOML(globalConfig, config map[string]any, chainID int) (map[string]any, error) {
	result := make(map[string]any)

	// Start with global config
	if err := mergo.Merge(&result, globalConfig, mergo.WithOverride); err != nil {
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

// BuildGlobalConfigFromTOML builds global configuration from TOML
func BuildGlobalConfigFromTOML(config map[string]any) (map[string]any, error) {
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

	// Merge runtime fallbacks without overriding existing user values
	// By default, mergo.Merge won't override existing keys (no WithOverride flag)
	mergo.Merge(&result, runtimeFallbacks)

	return result
}
