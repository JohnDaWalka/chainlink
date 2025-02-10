package chainlink

import "github.com/smartcontractkit/chainlink/v2/core/config/toml"

type wasmConfig struct {
	w toml.Wasm
}

func (t *wasmConfig) SerialisedModulesDir() string {
	return *t.w.SerialisedModulesDir
}
