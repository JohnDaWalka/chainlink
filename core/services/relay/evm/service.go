package evm

import (
	"context"
	"fmt"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-evm/pkg/service"
)

// TODO there should be an evm capability config that can override the relayer defaults
func NewEVMService(_ context.Context, relayer *Relayer) (capabilities.ExecutableAndTriggerCapability, error) {
	evmRelayer, err := relayer.EVM()
	if err != nil {
		return nil, err
	}

	opts := service.EVMCapabilityOpts{
		ID:         GenerateEVMCapabilityName(relayer.chain.ID().Uint64()),
		EVMService: evmRelayer,
	}

	return service.NewEVMCapablity(opts), nil
}

func GenerateEVMCapabilityName(chainID uint64) string {
	id := fmt.Sprintf("evm_service_%v@1.0.0", chainID)
	chainName, err := chainselectors.NameFromChainId(chainID)
	if err == nil {
		id = fmt.Sprintf("evm_service_%v@1.0.0", chainName)
	}

	return id
}
