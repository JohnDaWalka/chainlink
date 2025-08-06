package gateway

import (
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/vault"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

const (
	// Note: these values should correspond to the first part of the method name
	// in order to use them with the multi handler.
	// Eg. handler type = functions, method name = functions.do_something
	FunctionsHandlerType HandlerType = "functions"
	DummyHandlerType     HandlerType = "dummy"
	HTTPCapabilityType   HandlerType = "http-capabilities" // Handler for v1.0 HTTP capabilities for NoDAG workflows
	VaultHandlerType     HandlerType = "vault"
)

type handlerFactory struct {
	legacyChains legacyevm.LegacyChainContainer
	ds           sqlutil.DataSource
	lggr         logger.Logger
	httpClient   network.HTTPClient
}

var _ HandlerFactory = (*handlerFactory)(nil)

func NewHandlerFactory(legacyChains legacyevm.LegacyChainContainer, ds sqlutil.DataSource, httpClient network.HTTPClient, lggr logger.Logger) HandlerFactory {
	return &handlerFactory{
		legacyChains,
		ds,
		lggr,
		httpClient,
	}
}

func (hf *handlerFactory) NewHandler(handlerType HandlerType, handlerConfig json.RawMessage, donConfig *config.DONConfig, don handlers.DON) (handlers.Handler, error) {
	switch handlerType {
	case FunctionsHandlerType:
		return functions.NewFunctionsHandlerFromConfig(handlerConfig, donConfig, don, hf.legacyChains, hf.ds, hf.lggr)
	case DummyHandlerType:
		return handlers.NewDummyHandler(donConfig, don, hf.lggr)
	case HandlerType(capabilities.HandlerType):
		return capabilities.NewHandler(handlerConfig, donConfig, don, hf.httpClient, hf.lggr)
	case HTTPCapabilityType:
		return v2.NewGatewayHandler(handlerConfig, donConfig, don, hf.httpClient, hf.lggr)
	case HandlerType(vault.HandlerType):
		return vault.NewHandler(handlerConfig, donConfig, don, hf.lggr)
	default:
		return nil, fmt.Errorf("unsupported handler type %s", handlerType)
	}
}
