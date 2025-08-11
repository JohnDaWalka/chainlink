package jobs

import (
	"maps"
	"regexp"
	"strconv"
	"strings"

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

// ApplyRuntimeValues fills in any missing config values with runtime-generated values
func ApplyRuntimeValues(userConfig map[string]any, runtimeValues map[string]any) map[string]any {
	result := make(map[string]any)
	maps.Copy(result, userConfig)

	// Merge runtime fallbacks without overriding existing user values
	// By default, mergo.Merge won't override existing keys (no WithOverride flag)
	mergo.Merge(&result, runtimeValues)

	return result
}

// ValidateTemplateSubstitution checks that all template variables have been properly substituted
// Returns an error if any {{.Variable}} patterns are found in the rendered output
//
// This function helps catch configuration issues early by ensuring all template variables
// have been provided and substituted. Common causes of unsubstituted variables:
// - Missing fields in templateData map
// - Typos in template variable names
// - Missing values in runtimeValues
//
// Usage:
//
//	configStr := configBuffer.String()
//	if err := ValidateTemplateSubstitution(configStr, "capability-name"); err != nil {
//	    return nil, errors.Wrap(err, "template validation failed")
//	}
func ValidateTemplateSubstitution(rendered string, templateName string) error {
	// Regex to find unsubstituted template variables like {{.Variable}}
	templateVarRegex := regexp.MustCompile(`\{\{\s*\.[A-Za-z_][A-Za-z0-9_]*\s*\}\}`)

	matches := templateVarRegex.FindAllString(rendered, -1)
	if len(matches) > 0 {
		return errors.Errorf("template '%s' contains unsubstituted variables: %s",
			templateName, strings.Join(matches, ", "))
	}

	return nil
}
