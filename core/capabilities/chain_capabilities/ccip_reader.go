package chain_capabilities

import (
	"context"
	"fmt"
	"iter"
	"sync"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type CCIPReader struct {
	readMapper
	cr   commontypes.ContractReader
	lggr logger.Logger
	// key is chain family, value is
	Addresses sync.Mutex
	commontypes.UnimplementedContractReader
}

type readMapper struct {
	family  string
	chainID string
	cc      commontypes.ChainCapabilities

	// only one of these is initialised per CCIPReader readMapper depending on family
	solanaContract map[string]solanaContract
	evmContract    map[string]evmContract
}

type solanaGetLatestValueReadFunc func(reader commontypes.SolanaChainReader, address string, confidenceLevel primitives.ConfidenceLevel, params, returnVal any) error
type evmGetLatestValueReadFunc func(reader commontypes.EVMChainReader, address string, confidenceLevel primitives.ConfidenceLevel, params, returnVal any) error

type solanaContract struct {
	address             string
	getLatestValueReads map[string]solanaGetLatestValueReadFunc
}

type evmContract struct {
	address             string
	getLatestValueReads map[string]evmGetLatestValueReadFunc
}

func (r *readMapper) getLatestValue(contractName, methodName string, confidenceLevel primitives.ConfidenceLevel, params, returnVal any) error {
	switch r.family {
	case chain_selectors.FamilyEVM:
		contract := r.evmContract[contractName]
		return contract.getLatestValueReads[methodName](r.cc, contract.address, confidenceLevel, params, returnVal)
	case chain_selectors.FamilySolana:
		contract := r.solanaContract[contractName]
		return contract.getLatestValueReads[methodName](r.cc, contract.address, confidenceLevel, params, returnVal)
	default:
		return fmt.Errorf("unsupported chain family %s", r.family)
	}
}

func (r *readMapper) bind(contractName, address string) error {
	switch r.family {
	case chain_selectors.FamilyEVM:
		contract, exists := r.evmContract[contractName]
		if !exists {
			return fmt.Errorf("contract %s not found", contractName)
		}
		contract.address = address
		r.evmContract[contractName] = contract
	case chain_selectors.FamilySolana:
		contract, exists := r.solanaContract[contractName]
		if !exists {
			return fmt.Errorf("contract %s not found", contractName)
		}
		contract.address = address
		r.solanaContract[contractName] = contract
	default:
		return fmt.Errorf("unsupported chain family %s", r.family)
	}

	return nil
}

var _ commontypes.ContractReader = (*CCIPReader)(nil) //or

// NewCCIPReader should be used as a shim that allows CCIP to have a complete pass through to Contract Reader or to implement chain aware CCIP reads using the chain capabilities.
func NewCCIPReader(contractReader commontypes.ContractReader, cc commontypes.ChainCapabilities, family, chainID string) *CCIPReader {

	rm := readMapper{
		family:         family,
		chainID:        chainID,
		cc:             cc,
		solanaContract: make(map[string]solanaContract),
		evmContract:    make(map[string]evmContract),
	}

	rm.solanaContract[consts.MethodNameFeeQuoterGetTokenPrices].getLatestValueReads[consts.MethodNameFeeQuoterGetTokenPrices] = readTokenPricesSolana
	//rm.solanaContract[consts.MethodNameFeeQuoterGetTokenPrices].getLatestValueReads[consts.MethodNameFeeQuoterGetTokenPrices] = readTokenPricesEVM

	return &CCIPReader{
		cr: contractReader,
	}
}

func (c *CCIPReader) Start(ctx context.Context) error {
	return c.cr.Start(ctx)
}

func (c *CCIPReader) Close() error {
	return c.cr.Close()
}

func (c *CCIPReader) Ready() error {
	return c.cr.Ready()
}

func (c *CCIPReader) HealthReport() map[string]error {
	return c.cr.HealthReport()
}

func (c *CCIPReader) Name() string {
	return c.Name()
}

// TODO Bind should bind addresses to contract names for
func (c *CCIPReader) Bind(ctx context.Context, bindings []commontypes.BoundContract) error {
	// TODO implement a switch case logic for choosing whether to do a pass through to contract reader
	c.Addresses.Lock()
	for _, binding := range bindings {
		if err := c.readMapper.bind(binding.Name, binding.Address); err != nil {
			return err
		}

	}
	c.Addresses.Unlock()

	return c.Bind(ctx, bindings)
}

func (c *CCIPReader) Unbind(ctx context.Context, bindings []commontypes.BoundContract) error {
	// TODO implement a switch case logic for choosing whether to do a pass through to contract reader
	//c.Addresses.Lock()
	//for _, binding := range bindings {
	//	delete(c.ContractNameToAddress, binding.Name)
	//}
	//c.Addresses.Unlock()

	return c.Unbind(ctx, bindings)
}

func (c *CCIPReader) GetLatestValue(ctx context.Context, readIdentifier string, confidenceLevel primitives.ConfidenceLevel, params, returnVal any) error {
	// TODO implement a switch case logic for choosing whether to do a pass through to contract reader
	//if !passThrough
	// TODO implement logic for parsing readIdentifier
	//		return c.readMapper.getLatestValue("", readIdentifier, confidenceLevel, params, returnVal)
	// } else {
	return c.cr.GetLatestValue(ctx, readIdentifier, confidenceLevel, params, returnVal)
	// }
}

func (c *CCIPReader) GetLatestValueWithHeadData(ctx context.Context, readIdentifier string, confidenceLevel primitives.ConfidenceLevel, params, returnVal any) (*commontypes.Head, error) {
	return c.cr.GetLatestValueWithHeadData(ctx, readIdentifier, confidenceLevel, params, returnVal)
}

func (c *CCIPReader) BatchGetLatestValues(ctx context.Context, request commontypes.BatchGetLatestValuesRequest) (commontypes.BatchGetLatestValuesResult, error) {
	return c.BatchGetLatestValues(ctx, request)
}

func (c *CCIPReader) QueryKey(ctx context.Context, contract commontypes.BoundContract, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]commontypes.Sequence, error) {
	return c.cr.QueryKey(ctx, contract, filter, limitAndSort, sequenceDataType)
}

func (c *CCIPReader) QueryKeys(ctx context.Context, filters []commontypes.ContractKeyFilter, limitAndSort query.LimitAndSort) (iter.Seq2[string, commontypes.Sequence], error) {
	return c.QueryKeys(ctx, filters, limitAndSort)
}
