package chainlink

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWasmConfigTest(t *testing.T) {
	opts := GeneralConfigOpts{
		ConfigStrings:  []string{fullTOML},
		SecretsStrings: []string{secretsFullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	wcfg := cfg.Wasm()

	require.Equal(t, "test/root/dir", wcfg.SerialisedModulesDir())
}
