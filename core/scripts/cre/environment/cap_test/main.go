package main

import (
	"fmt"
	"log"
	"os"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

type Config struct {
	AdditionalCapabilities map[string]CapabilityConfig `toml:"additional_capabilities"` // capability name -> capability config
}

type CapabilityConfig struct {
	BinaryPath   string         `toml:"binary_path"`
	Config       string         `toml:"config"`
	Chains       string         `toml:"chains"`
	ChainConfigs map[string]any `toml:"chain_configs"`
}

func main() {
	_ = os.Setenv("CTF_CONFIGS", "cap_config.toml")
	capConfig, err := framework.Load[Config](nil)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	for name, capConfig := range capConfig.AdditionalCapabilities {
		fmt.Printf("Capability: %s\n", name)
		fmt.Printf("  BinaryPath: %s\n", capConfig.BinaryPath)
		fmt.Printf("  Config: %s\n", capConfig.Config)
		fmt.Printf("  Chains: %s\n", capConfig.Chains)
		for chainID, chainConfig := range capConfig.ChainConfigs {
			fmt.Printf("  Chain %s:\n", chainID)
			fmt.Printf("    Config: %+v\n", chainConfig)
		}
		fmt.Println()
	}
}
