package writesolana

import "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"

type Config struct {
}

func GetGenerateConfig(in Config) func(types.GenerateConfigsInput) (types.NodeIndexToConfigOverride, error) {
	return nil
}
