package chain_capabilities

import (
	"context"
	"iter"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
)

type CCIPReader struct {
	family  string
	chainID string
	cr      commontypes.ContractReader
	cc      commontypes.ChainCapabilities
	commontypes.UnimplementedContractReader
}

var _ commontypes.ContractReader = (*CCIPReader)(nil) //or

func NewCCIPReader(contractReader commontypes.ContractReader, family, chainID string) *CCIPReader {
	return &CCIPReader{
		cr:      contractReader,
		family:  family,
		chainID: chainID,
	}
}

func (c CCIPReader) Start(ctx context.Context) error {
	return c.cr.Start(ctx)
}

func (c CCIPReader) Close() error {
	return c.cr.Close()
}

func (c CCIPReader) Ready() error {
	return c.cr.Ready()
}

func (c CCIPReader) HealthReport() map[string]error {
	return c.cr.HealthReport()
}

func (c CCIPReader) Name() string {
	return c.Name()
}

func (c CCIPReader) Bind(ctx context.Context, bindings []commontypes.BoundContract) error {
	return c.Bind(ctx, bindings)
}

func (c CCIPReader) Unbind(ctx context.Context, bindings []commontypes.BoundContract) error {
	return c.Unbind(ctx, bindings)
}

func (c CCIPReader) GetLatestValue(ctx context.Context, readIdentifier string, confidenceLevel primitives.ConfidenceLevel, params, returnVal any) error {
	return c.cr.GetLatestValue(ctx, readIdentifier, confidenceLevel, params, returnVal)
}

func (c CCIPReader) GetLatestValueWithHeadData(ctx context.Context, readIdentifier string, confidenceLevel primitives.ConfidenceLevel, params, returnVal any) (*commontypes.Head, error) {
	return c.cr.GetLatestValueWithHeadData(ctx, readIdentifier, confidenceLevel, params, returnVal)
}

func (c CCIPReader) BatchGetLatestValues(ctx context.Context, request commontypes.BatchGetLatestValuesRequest) (commontypes.BatchGetLatestValuesResult, error) {
	return c.BatchGetLatestValues(ctx, request)
}

func (c CCIPReader) QueryKey(ctx context.Context, contract commontypes.BoundContract, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]commontypes.Sequence, error) {
	return c.cr.QueryKey(ctx, contract, filter, limitAndSort, sequenceDataType)
}

func (c CCIPReader) QueryKeys(ctx context.Context, filters []commontypes.ContractKeyFilter, limitAndSort query.LimitAndSort) (iter.Seq2[string, commontypes.Sequence], error) {
	return c.QueryKeys(ctx, filters, limitAndSort)
}
