package crecli

import (
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-yaml/yaml"
	"github.com/stretchr/testify/require"
)

const (
	CRECLISettingsFileName = ".cre-cli-settings.yaml"
)

var (
	// absolute path to the CRE CLI binary
	CRECLICommand string
)

type Settings struct {
	DevPlatform  DevPlatform  `yaml:"dev-platform"`
	UserWorkflow UserWorkflow `yaml:"user-workflow"`
	Logging      Logging      `yaml:"logging"`
	McmsConfig   McmsConfig   `yaml:"mcms-config"`
	Contracts    Contracts    `yaml:"contracts"`
	Rpcs         []RPC        `yaml:"rpcs"`
}

type DevPlatform struct {
	CapabilitiesRegistryAddress string `yaml:"capabilities-registry-contract-address"`
	DonID                       uint32 `yaml:"don-id"`
	WorkflowRegistryAddress     string `yaml:"workflow-registry-contract-address"`
}

type UserWorkflow struct {
	WorkflowOwnerAddress string `yaml:"workflow-owner-address"`
}

type Logging struct {
	SethConfigPath string `yaml:"seth-config-path"`
}

type McmsConfig struct {
	ProposalsDirectory string `yaml:"proposals-directory"`
}

type Contracts struct {
	ContractRegistry []ContractRegistry `yaml:"registries"`
}

type ContractRegistry struct {
	Name          string `yaml:"name"`
	Address       string `yaml:"address"`
	ChainSelector uint64 `yaml:"chain-selector"`
}

type RPC struct {
	ChainSelector uint64 `yaml:"chain-selector"`
	URL           string `yaml:"url"`
}

type PoRWorkflowConfig struct {
	FeedID          string `json:"feed_id"`
	URL             string `json:"url"`
	ConsumerAddress string `json:"consumer_address"`
}

func PrepareCRECLISettingsFile(t *testing.T, workflowOwner, capRegAddr, workflowRegistryAddr common.Address, donID uint32, chainSelector uint64, rpcHTTPURL string) *os.File {
	settingsFile, err := os.CreateTemp("", CRECLISettingsFileName)
	require.NoError(t, err, "failed to create CRE CLI settings file")

	settings := Settings{
		DevPlatform: DevPlatform{
			CapabilitiesRegistryAddress: capRegAddr.Hex(),
			DonID:                       donID,
			WorkflowRegistryAddress:     workflowRegistryAddr.Hex(),
		},
		UserWorkflow: UserWorkflow{
			WorkflowOwnerAddress: workflowOwner.Hex(),
		},
		Logging: Logging{},
		McmsConfig: McmsConfig{
			ProposalsDirectory: "./",
		},
		Contracts: Contracts{
			ContractRegistry: []ContractRegistry{
				{
					Name:          "CapabilitiesRegistry",
					Address:       capRegAddr.Hex(),
					ChainSelector: chainSelector,
				},
				{
					Name:          "WorkflowRegistry",
					Address:       workflowRegistryAddr.Hex(),
					ChainSelector: chainSelector,
				},
			},
		},
		Rpcs: []RPC{
			{
				ChainSelector: chainSelector,
				URL:           rpcHTTPURL,
			},
		},
	}

	settingsMarshalled, err := yaml.Marshal(settings)
	require.NoError(t, err, "failed to marshal CRE CLI settings")

	_, err = settingsFile.Write(settingsMarshalled)
	require.NoError(t, err, "failed to write %s settings file", CRECLISettingsFileName)

	return settingsFile
}
