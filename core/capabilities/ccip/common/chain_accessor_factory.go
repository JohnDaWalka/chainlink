package common

import (
	"github.com/smartcontractkit/chainlink-ccip/pkg/contractreader"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// TODO: Better to define here or in chainlink-ccip?
type ChainAccessorFactory interface {
	NewChainAccessor(
		lggr logger.Logger,
		chainSelector cciptypes.ChainSelector,
		contractReader contractreader.Extended,
		contractWriter types.ContractWriter,
		addrCodec cciptypes.AddressCodec,
	) (cciptypes.ChainAccessor, error)
}
