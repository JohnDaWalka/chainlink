package common

import (
	"context"
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/pkg/contractreader"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

// ChainReaderProviderOpts is a struct that contains the parameters for GetChainReader.
type ChainReaderProviderOpts struct {
	Lggr            logger.Logger
	Relayer         loop.Relayer
	ChainID         string
	DestChainID     string
	HomeChainID     string
	Ofc             OffChainConfig
	ChainSelector   cciptypes.ChainSelector
	ChainFamily     string
	DestChainFamily string
	Transmitters    map[types.RelayID][]string
}

// ChainWriterProviderOpts is a struct that contains the parameters for GetChainWriter.
type ChainWriterProviderOpts struct {
	ChainID                        string
	Relayer                        loop.Relayer
	Transmitters                   map[types.RelayID][]string
	ExecBatchGasLimit              uint64
	ChainFamily                    string
	OfframpProgramAddress          []byte
	SolanaChainWriterConfigVersion *string
}

// ChainRWProvider is an interface that defines the methods to get a ContractReader and a ContractWriter.
type ChainRWProvider interface {
	GetChainReader(ctx context.Context, params ChainReaderProviderOpts) (types.ContractReader, error)
	GetChainWriter(ctx context.Context, params ChainWriterProviderOpts) (types.ContractWriter, error)
}

// MultiChainRW is a struct that implements the ChainRWProvider interface for all chains.
type MultiChainRW struct {
	cwProviderMap map[string]ChainRWProvider
}

// NewCRCW is a constructor for MultiChainRW.
func NewCRCW(cwProviderMap map[string]ChainRWProvider) MultiChainRW {
	return MultiChainRW{
		cwProviderMap: cwProviderMap,
	}
}

// GetChainReader returns a new ContractReader base on relay chain family.
func (c *MultiChainRW) GetChainReader(ctx context.Context, params ChainReaderProviderOpts) (types.ContractReader, error) {
	provider, exist := c.cwProviderMap[params.ChainFamily]
	if !exist {
		return nil, fmt.Errorf("unsupported chain family %s", params.ChainFamily)
	}

	return provider.GetChainReader(ctx, params)
}

// GetChainWriter returns a new ContractWriter based on relay chain family.
func (c *MultiChainRW) GetChainWriter(ctx context.Context, params ChainWriterProviderOpts) (types.ContractWriter, error) {
	provider, exist := c.cwProviderMap[params.ChainFamily]
	if !exist {
		return nil, fmt.Errorf("unsupported chain family %s", params.ChainFamily)
	}

	return provider.GetChainWriter(ctx, params)
}

func WrapContractReaderInObservedExtended(
	lggr logger.Logger,
	contractReader types.ContractReader,
	chainSelector cciptypes.ChainSelector,
) (contractreader.Extended, error) {
	chainFamily, err1 := chainsel.GetSelectorFamily(uint64(chainSelector))
	if err1 != nil {
		return nil, fmt.Errorf("failed to get chain family from selector: %w", err1)
	}
	chainID, err1 := chainsel.GetChainIDFromSelector(uint64(chainSelector))
	if err1 != nil {
		return nil, fmt.Errorf("failed to get chain id from selector: %w", err1)
	}
	// NewExtendedContractReader() protects against double wrapping an extended reader.
	reader := contractreader.NewExtendedContractReader(
		contractreader.NewObserverReader(contractReader, lggr, chainFamily, chainID),
	)
	if reader == nil {
		return nil, fmt.Errorf("failed to create extended contract reader for chain selector %d", chainSelector)
	}
	return reader, nil
}
