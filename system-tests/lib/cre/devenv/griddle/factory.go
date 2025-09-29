package griddle

import "github.com/smartcontractkit/chainlink/system-tests/lib/infra"

func ConfigFromInfraInput(infraInput infra.Provider) GriddleConfig {
	metadata := infraInput.GriddleDevenvInput.GriddleMetadata

	// Create a config struct
	config := GriddleConfig{
		Version: 1,
		Metadata: Metadata{
			Account: metadata.Account,
			Project: metadata.Project,
			Service: metadata.Service,
			Owner:   metadata.Owner,
			Contact: metadata.Contact,
		},
	}

	return config
}

// ConfigFromInfraInputWithTemplate creates a GriddleConfig from the given infra input and adds a custom deploy template.
func ConfigFromInfraInputWithTemplate(infraInput infra.Provider, templateName string, template DeployTemplate) GriddleConfig {
	config := ConfigFromInfraInput(infraInput)

	config.DeployTemplates = map[string]DeployTemplate{
		templateName: template,
	}

	return config
}
