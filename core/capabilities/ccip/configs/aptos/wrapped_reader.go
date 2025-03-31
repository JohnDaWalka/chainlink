package aptosconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	//"errors"
	//"iter"

	"github.com/smartcontractkit/chainlink-aptos/relayer/chainreader"
	"github.com/smartcontractkit/chainlink-aptos/relayer/codec"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
)

func NewWrappedChainReader(logger logger.Logger, cr types.ContractReader) types.ContractReader {
	fmt.Printf("DEBUG: wrappedChainReader.NewWrappedChainReader\n")
	config, _ := GetChainReaderConfig()
	return &wrappedChainReader{logger: logger, cr: cr, config: config, moduleAddresses: map[string]string{}}
}

type wrappedChainReader struct {
	services.Service
	types.UnimplementedContractReader
	logger          logger.Logger
	cr              types.ContractReader
	config          chainreader.ChainReaderConfig
	moduleAddresses map[string]string
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

func (a *wrappedChainReader) GetLatestValue(ctx context.Context, readIdentifier string, confidenceLevel primitives.ConfidenceLevel, params, returnVal any) error {
	readComponents := strings.Split(readIdentifier, "-")
	if len(readComponents) != 3 {
		return fmt.Errorf("invalid read identifier: %s", readIdentifier)
	}

	_, contractName, _ := readComponents[0], readComponents[1], readComponents[2]

	_, ok := a.moduleAddresses[contractName]
	if !ok {
		return fmt.Errorf("no such contract: %s", contractName)
	}

	convertedResult := []byte{}

	jsonParamBytes, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal params: %+w", err)
	}

	// we always bind before calling query functions, because the LOOP plugin may have restarted.
	err = a.cr.Bind(ctx, a.getBindings())
	if err != nil {
		return fmt.Errorf("failed to re-bind before GetLatestValue: %w", err)
	}

	err = a.cr.GetLatestValue(ctx, readIdentifier, confidenceLevel, jsonParamBytes, &convertedResult)
	if err != nil {
		return fmt.Errorf("failed to call GetLatestValue over LOOP: %w", err)
	}

	fmt.Printf("DEBUG: wrappedChainReader.GetLatestValue convertedResult: %s\n", string(convertedResult))

	err = a.decodeGLVReturnValue(readIdentifier, convertedResult, returnVal)
	if err != nil {
		return fmt.Errorf("failed to decode GetLatestValue return value: %w", err)
	}

	return nil
}

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

	// we always bind before calling query functions, because the LOOP plugin may have restarted.
	err := a.cr.Bind(ctx, a.getBindings())
	if err != nil {
		return nil, fmt.Errorf("failed to re-bind before BatchGetLatestValues: %w", err)
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
				err := a.decodeGLVReturnValue(result.ReadName, *resultPointer, read.ReturnVal)
				if err != nil {
					resultError = fmt.Errorf("failed to decode BatchGetLatestValue return value: %w", err)
				}
			}
			convertedResult.SetResult(read.ReturnVal, resultError)
			convertedBatch = append(convertedBatch, convertedResult)
		}
		convertedResult[contract] = convertedBatch
	}

	return convertedResult, nil
}

func (a *wrappedChainReader) QueryKey(ctx context.Context, contract types.BoundContract, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]types.Sequence, error) {
	err := a.cr.Bind(ctx, a.getBindings())
	if err != nil {
		return nil, fmt.Errorf("failed to re-bind before BatchGetLatestValues: %w", err)
	}

	// the relay wrapper calls getContractEncodedType which defaults to returning a map[string]any
	// https://github.com/smartcontractkit/chainlink-common/blob/fe3ec4466fb5adfffd8fc77eef1cef67c4a918cc/pkg/loop/internal/relayer/pluginprovider/contractreader/contract_reader.go#L1033
	// in ccipChainReader.ExecutedMessages, it's a primitive

	convertedExpressions := []query.Expression{}
	for _, expr := range filter.Expressions {
		convertedExpressions = append(convertedExpressions, query.Expression{
			Comparator: expr.Comparator,
			Value:      expr.Value,
			Operator:   expr.Operator,
		})
	}

	convertedFilter := query.KeyFilter{
		Key:         filter.Key,
		Expressions: convertedExpressions,
	}

	sequences, err := a.cr.QueryKey(ctx, contract, filter, limitAndSort, &[]byte{})
	if err != nil {
		return nil, fmt.Errorf("failed to call QueryKey over LOOP: %w", err)
	}

	for i, sequence := range sequences {
		jsonBytes := sequence.Data.(*[]byte)
		jsonData := map[string]any{}
		err := json.Unmarshal(*jsonBytes, &jsonData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal LOOP sourced JSON event data (`%s`): %w", string(*jsonBytes), err)
		}

		eventData := reflect.New(reflect.TypeOf(sequenceDataType).Elem()).Interface()
		err = codec.DecodeAptosJsonValue(jsonData, &eventData)
		if err != nil {
			return nil, fmt.Errorf("failed to decode LOOP sourced event data (`%s`) into an Aptos value: %+w", string(*jsonBytes), err)
		}

		sequences[i].Data = eventData
	}

	return sequences, nil
}

func (a *wrappedChainReader) Bind(ctx context.Context, bindings []types.BoundContract) error {
	for _, binding := range bindings {
		a.moduleAddresses[binding.Name] = binding.Address
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

	// we ignore unbind errors, because if the LOOP plugin restarted, the binding would not exist.
	_ = a.cr.Unbind(ctx, bindings)

	return nil
}

func (a *wrappedChainReader) getBindings() []types.BoundContract {
	bindings := []types.BoundContract{}

	for name, address := range a.moduleAddresses {
		bindings = append(bindings, types.BoundContract{
			Address: address,
			Name:    name,
		})
	}

	return bindings
}

func (a *wrappedChainReader) decodeGLVReturnValue(label string, jsonBytes []byte, returnVal any) error {
	var unmarshalledData []any
	err := json.Unmarshal(jsonBytes, &unmarshalledData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %s GetLatestValue result (`%s`): %w", label, string(jsonBytes), err)
	}

	var unwrappedData any
	if len(unmarshalledData) == 1 {
		unwrappedData = unmarshalledData[0]
	} else {
		unwrappedData = unmarshalledData
	}

	err = codec.DecodeAptosJsonValue(unwrappedData, returnVal)
	if err != nil {
		return fmt.Errorf("failed to decode %s GetLatestValue JSON value (`%s`) to %T: %w", label, string(jsonBytes), returnVal, err)
	}

	return nil
}
