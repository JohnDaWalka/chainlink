package griddle

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

const testYAML = `version: 1

metadata:
  account: test-account
  project: devenv-test
  service: devenv-test
  owner: team-name
  contact: test@abc.xyz

deployTemplates:
  blockchain-anvil:
    instances:
      - name: anvil-1337
        chart: app-template
        version: 4.3.0
        repository: https://bjw-s-labs.github.io/helm-charts
        localRepositoryName: bjw-s
        config:
          - deploy/config/base/values/bjw-s/blockchain.anvil.base.yaml
          - deploy/config/base/values/bjw-s/blockchain.anvil.1337.yaml
      - name: anvil-2337
        chart: app-template
        version: 4.3.0
        repository: https://bjw-s-labs.github.io/helm-charts
        localRepositoryName: bjw-s
        config:
          - deploy/config/base/values/bjw-s/blockchain.anvil.base.yaml
          - deploy/config/base/values/bjw-s/blockchain.anvil.2337.yaml
  blockchain-anvil-deps:
    instances:
      - name: anvil-1337
        chart: app-template
        version: 4.3.0
        repository: https://bjw-s-labs.github.io/helm-charts
        localRepositoryName: bjw-s
        config:
          - deploy/config/base/values/bjw-s/blockchain.anvil.base.yaml
          - deploy/config/base/values/bjw-s/blockchain.anvil.1337.yaml
      - name: anvil-2337
        chart: app-template
        version: 4.3.0
        repository: https://bjw-s-labs.github.io/helm-charts
        localRepositoryName: bjw-s
        dependsOn:
          - name: anvil-1337
        config:
          - deploy/config/base/values/bjw-s/blockchain.anvil.base.yaml
          - deploy/config/base/values/bjw-s/blockchain.anvil.2337.yaml
      - name: anvil-31337
        chart: app-template
        version: 4.3.0
        repository: https://bjw-s-labs.github.io/helm-charts
        localRepositoryName: bjw-s
        dependsOn:
          - name: anvil-1337
          - name: anvil-2337
        config:
          - deploy/config/base/values/bjw-s/blockchain.anvil.base.yaml
          - deploy/config/base/values/bjw-s/blockchain.anvil.31337.yaml
`

func TestGriddleConfigUnmarshall(t *testing.T) {
	var config GriddleConfig
	err := yaml.Unmarshal([]byte(testYAML), &config)
	require.NoError(t, err)

	// Test metadata
	assert.Equal(t, 1, config.Version)
	assert.Equal(t, "test-account", config.Metadata.Account)
	assert.Equal(t, "devenv-test", config.Metadata.Project)
	assert.Equal(t, "devenv-test", config.Metadata.Service)
	assert.Equal(t, "team-name", config.Metadata.Owner)
	assert.Equal(t, "test@abc.xyz", config.Metadata.Contact)

	// Test deploy templates
	assert.Len(t, config.DeployTemplates, 2)

	// Test blockchain-anvil template
	blockchainAnvil := config.DeployTemplates["blockchain-anvil"]
	assert.Len(t, blockchainAnvil.Instances, 2)

	anvil1337 := blockchainAnvil.Instances[0]
	assert.Equal(t, "anvil-1337", anvil1337.Name)
	assert.Equal(t, "app-template", anvil1337.Chart)
	assert.Equal(t, "4.3.0", anvil1337.Version)
	assert.Equal(t, "https://bjw-s-labs.github.io/helm-charts", anvil1337.Repository)
	assert.Equal(t, "bjw-s", anvil1337.LocalRepositoryName)
	assert.Len(t, anvil1337.Config, 2)
	assert.Empty(t, anvil1337.DependsOn)

	// Test blockchain-anvil-deps template with dependencies
	blockchainAnvilDeps := config.DeployTemplates["blockchain-anvil-deps"]
	assert.Len(t, blockchainAnvilDeps.Instances, 3)

	// Test instance with single dependency
	anvil2337Deps := blockchainAnvilDeps.Instances[1]
	assert.Equal(t, "anvil-2337", anvil2337Deps.Name)
	assert.Len(t, anvil2337Deps.DependsOn, 1)
	assert.Equal(t, "anvil-1337", anvil2337Deps.DependsOn[0].Name)

	// Test instance with multiple dependencies
	anvil31337 := blockchainAnvilDeps.Instances[2]
	assert.Equal(t, "anvil-31337", anvil31337.Name)
	assert.Len(t, anvil31337.DependsOn, 2)
	assert.Equal(t, "anvil-1337", anvil31337.DependsOn[0].Name)
	assert.Equal(t, "anvil-2337", anvil31337.DependsOn[1].Name)
}

func TestGriddleConfigMarshall(t *testing.T) {
	// Create a config struct
	config := GriddleConfig{
		Version: 1,
		Metadata: Metadata{
			Account: "test-account",
			Project: "devenv-test",
			Service: "devenv-test",
			Owner:   "team-name",
			Contact: "test@abc.xyz",
		},
		DeployTemplates: map[string]DeployTemplate{
			"blockchain-anvil": {
				Instances: []Instance{
					{
						Name:                "anvil-1337",
						Chart:               "app-template",
						Version:             "4.3.0",
						Repository:          "https://bjw-s-labs.github.io/helm-charts",
						LocalRepositoryName: "bjw-s",
						Config: []string{
							"deploy/config/base/values/bjw-s/blockchain.anvil.base.yaml",
							"deploy/config/base/values/bjw-s/blockchain.anvil.1337.yaml",
						},
					},
				},
			},
		},
	}

	// Marshal to YAML
	yamlData, err := yaml.Marshal(&config)
	require.NoError(t, err)

	// Unmarshal back to verify round-trip works
	var roundTripConfig GriddleConfig
	err = yaml.Unmarshal(yamlData, &roundTripConfig)
	require.NoError(t, err)

	assert.Equal(t, config.Version, roundTripConfig.Version)
	assert.Equal(t, config.Metadata.Account, roundTripConfig.Metadata.Account)
	assert.Equal(t, "anvil-1337", roundTripConfig.DeployTemplates["blockchain-anvil"].Instances[0].Name)
}
