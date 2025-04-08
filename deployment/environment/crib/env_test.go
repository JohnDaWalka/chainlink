package crib

import (
	"testing"

	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldProvideEnvironmentConfig(t *testing.T) {
	t.Parallel()
	singeFileLogger := cldlogger.NewSingleFileLogger(t)
	env := NewCRIBEnvFromStateDir(singeFileLogger, "testdata/lanes-deployed-state")
	config, err := env.GetConfig("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")

	require.NoError(t, err)
	require.NotNil(t, config)
	assert.NotEmpty(t, config.NodeIDs)
	assert.NotNil(t, config.AddressBook)
	assert.NotEmpty(t, config.Chains)

	assert.Len(t, config.BlockchainOutputs, 2)
	assert.Equal(t, "evm", config.BlockchainOutputs["1337"].Family)
	assert.Len(t, config.BlockchainOutputs["1337"].Nodes, 1)

	assert.Equal(t, "2337", config.BlockchainOutputs["2337"].ChainID)

	assert.NotNil(t, config.NodesetOutput)
	assert.NotNil(t, config.NodesetOutput.Output)
	assert.NotEmpty(t, config.NodesetOutput.Output.CLNodes)
	assert.Len(t, config.NodesetOutput.Output.CLNodes, 5)
	assert.NotNil(t, config.NodesetOutput.CLNodes[0].Node)
	assert.NotEmpty(t, config.NodesetOutput.CLNodes[0].Node.ExternalURL)
	assert.Equal(t, "http://crib-local-ccip-bt-0.local:80", config.NodesetOutput.Output.CLNodes[0].Node.ExternalURL)

	assert.NotNil(t, config.JDOutput)
	assert.NotEmpty(t, config.JDOutput.ExternalGRPCUrl)
	assert.NotEmpty(t, config.JDOutput.InternalWSRPCUrl)
}
