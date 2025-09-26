package don

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

var (
	configToml = `[Capabilities]
[Capabilities.ExternalRegistry]
Address = '0x49fd2BE640DB2910c2fAb69bB8531Ab6E76127ff'
ChainID = '1337'
NetworkID = 'evm'

[[EVM]]
ChainID = '1337'

[[EVM.Nodes]]
HTTPURL = 'https://crib-local-geth-1337-http.devenv.test'
Name = 'anvil'
WSURL = 'wss://crib-local-geth-1337-ws.devenv.test'`

	secretsToml = `[EVM]
[[EVM.Keys]]
JSON = '{"address":"bf04453133a8cb5384a2c1d4f26cd3edf7a88aec"}'
Password = ''
ID = 1337

[P2PKey]
JSON = '{"keyType":"P2P","publicKey":"8b160c162eb9a17edf539f87a9810a90a663451743875618061ef36b8c494b36"}'
Password = ''`
)

func TestValuesGenerating(t *testing.T) {
	t.Parallel()

	bootNodeConfig := NewBootNodeValuesConfig().
		SetImage("image:latest").
		SetAppInstanceName("cre-bt-0").
		SetConfig(configToml).
		SetSecretsOverride(secretsToml).
		Build()

	node1Config := NewNodeValuesConfig().
		SetImage("image:latest").
		SetAppInstanceName("cre-node-0").
		SetConfig(configToml).
		SetSecretsOverride(secretsToml).
		Build()

	config := NewChainlinkClusterValuesConfig().
		SetNodes([]NodeValuesConfig{bootNodeConfig, node1Config}).
		Build()

	// Marshal to YAML
	yamlData, err := yaml.Marshal(config)
	assert.NoError(t, err)

	// Assert the full YAML content using go-snaps
	snaps.MatchStandaloneYAML(t, string(yamlData))

	assert.NotNil(t, config)
	assert.Equal(t, 2, config.NodeCount)
	assert.Equal(t, 1, config.BootNodeCount)
}
