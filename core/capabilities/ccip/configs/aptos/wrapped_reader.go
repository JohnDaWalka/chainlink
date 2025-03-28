package aptosconfig

import (
	"context"
	"encoding/json"
	"fmt"
	//"strings"
	//"errors"
	//"iter"

	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink-aptos/relayer/chainreader"
	//"github.com/smartcontractkit/chainlink-aptos/relayer/codec"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
)

func NewWrappedChainReader(logger logger.Logger, cr types.ContractReader) types.ContractReader {
	config, _ := GetChainReaderConfig()
	return &wrappedChainReader{logger: logger, cr: cr, config: config, moduleAddresses: map[string]aptos.AccountAddress{}}
}

type wrappedChainReader struct {
	services.Service
	types.UnimplementedContractReader
	logger          logger.Logger
	cr              types.ContractReader
	config          chainreader.ChainReaderConfig
	moduleAddresses map[string]aptos.AccountAddress
}

func (a *wrappedChainReader) Name() string {
	return a.cr.Name()
}

func (a *wrappedChainReader) Ready() error {
	return a.cr.Ready()
}

func (a *wrappedChainReader) HealthReport() map[string]error {
	return a.cr.HealthReport()
}

func (a *wrappedChainReader) Start(ctx context.Context) error {
	return a.cr.Start(ctx)
}

func (a *wrappedChainReader) Close() error {
	return a.cr.Close()
}

func (a *wrappedChainReader) getFunctionConfig(_address string, contractName string, method string) (*chainreader.ChainReaderFunction, error) {
	// Source the read configuration, by contract name
	address, ok := a.moduleAddresses[contractName]
	if !ok {
		return &chainreader.ChainReaderFunction{}, fmt.Errorf("no bound address for module %s", contractName)
	}

	// Notice: the address in the readIdentifier should match the bound address, by contract name
	if address.String() != _address {
		return &chainreader.ChainReaderFunction{}, fmt.Errorf("bound address %s for module %s does not match read address %s", address, contractName, _address)
	}

	moduleConfig, ok := a.config.Modules[contractName]
	if !ok {
		return &chainreader.ChainReaderFunction{}, fmt.Errorf("no such contract: %s", contractName)
	}

	if moduleConfig.Functions == nil {
		return &chainreader.ChainReaderFunction{}, fmt.Errorf("no functions for contract: %s", contractName)
	}

	functionConfig, ok := moduleConfig.Functions[method]
	if !ok {
		return &chainreader.ChainReaderFunction{}, fmt.Errorf("no such method: %s", method)
	}

	return functionConfig, nil
}

func (a *wrappedChainReader) GetLatestValue(ctx context.Context, readIdentifier string, confidenceLevel primitives.ConfidenceLevel, params, returnVal any) error {
	convertedResult := []byte{}

	jsonParamBytes, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal params: %+w", err)
	}

	err = a.cr.GetLatestValue(ctx, readIdentifier, confidenceLevel, jsonParamBytes, &convertedResult)
	if err != nil {
		return err
	}

	err = json.Unmarshal(convertedResult, returnVal)
	if err != nil {
		return err
	}

	return nil
}

//func (c *wrappedChainReader) GetLatestValueWithHeadData(ctx context.Context, readIdentifier string, confidenceLevel primitives.ConfidenceLevel, params, retVal any) (*types.Head, error) {
//return nil, errors.New("TODO")
//}

func (a *wrappedChainReader) BatchGetLatestValues(ctx context.Context, request types.BatchGetLatestValuesRequest) (types.BatchGetLatestValuesResult, error) {
	convertedRequest := types.BatchGetLatestValuesRequest{}
	for contract, requestBatch := range request {
		convertedBatch := []types.BatchRead{}
		for _, read := range requestBatch {
			jsonParamBytes, err := json.Marshal(read.Params)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal params: %+w", err)
			}
			convertedBatch = append(convertedBatch, types.BatchRead{
				ReadName:  read.ReadName,
				Params:    jsonParamBytes,
				ReturnVal: &[]byte{},
			})
		}
		convertedRequest[contract] = convertedBatch
	}

	result, err := a.cr.BatchGetLatestValues(ctx, convertedRequest)
	if err != nil {
		return nil, err
	}

	convertedResult := types.BatchGetLatestValuesResult{}
	for contract, resultBatch := range result {
		requestBatch := request[contract]
		convertedBatch := []types.BatchReadResult{}
		for i, result := range resultBatch {
			read := requestBatch[i]
			resultValue, resultError := result.GetResult()
			convertedResult := types.BatchReadResult{ReadName: result.ReadName}
			if resultError == nil {
				resultPointer := resultValue.(*[]byte)
				if err != nil {
					return nil, err
				}
				err = json.Unmarshal(*resultPointer, read.ReturnVal)
				convertedResult.SetResult(read.ReturnVal, err)
			}
			convertedBatch = append(convertedBatch, convertedResult)
		}
		convertedResult[contract] = convertedBatch
	}

	return convertedResult, nil
}

func (a *wrappedChainReader) QueryKey(ctx context.Context, contract types.BoundContract, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]types.Sequence, error) {
	// TODO
	return a.cr.QueryKey(ctx, contract, filter, limitAndSort, sequenceDataType)
}

//func (a *wrappedChainReader) QueryKeys(ctx context.Context, filters []types.ContractKeyFilter, limitAndSort query.LimitAndSort) (iter.Seq2[string, types.Sequence], error) {

//return nil, errors.New("TODO")
//}

func (a *wrappedChainReader) Bind(ctx context.Context, bindings []types.BoundContract) error {
	newBindings := map[string]aptos.AccountAddress{}
	for _, binding := range bindings {
		moduleAddress := &aptos.AccountAddress{}
		err := moduleAddress.ParseStringRelaxed(binding.Address)
		if err != nil {
			return fmt.Errorf("failed to convert module address %s: %+w", binding.Address, err)
		}
		newBindings[binding.Name] = *moduleAddress
	}

	for name, address := range newBindings {
		a.moduleAddresses[name] = address
	}

	return a.cr.Bind(ctx, bindings)
}

func (a *wrappedChainReader) Unbind(ctx context.Context, bindings []types.BoundContract) error {
	for _, binding := range bindings {
		key := binding.Name
		if _, ok := a.moduleAddresses[key]; ok {
			delete(a.moduleAddresses, key)
		} else {
			return fmt.Errorf("no such binding: %s", key)
		}
	}
	return a.cr.Unbind(ctx, bindings)
}
