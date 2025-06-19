package config

import (
	"encoding/json"

	gw_net "github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

type GatewayConfig struct {
	UserServerConfig        gw_net.HTTPServerConfig
	NodeServerConfig        gw_net.WebSocketServerConfig
	ConnectionManagerConfig ConnectionManagerConfig
	// HTTPClientConfig is configuration for outbound HTTP calls to external endpoints
	HTTPClientConfig gw_net.HTTPClientConfig
	Dons             []DONConfig
}

type ConnectionManagerConfig struct {
	AuthGatewayId             string
	AuthTimestampToleranceSec uint32
	AuthChallengeLen          uint32
	HeartbeatIntervalSec      uint32
}

const (
	// LegacyMessageType is for JSON-RPC methods with api.Message as params.
	// See: [api/message.go](../api/message.go). DonID is used to route the request to the appropriate handler.
	LegacyMessageType = "Legacy"

	// CustomParamsMessageType is for JSON-RPC methods with custom params.
	// HandlerName is used to route the request to the appropriate handler.
	CustomParamsMessageType = "CustomParams"
)

type DONConfig struct {
	DonId         string
	HandlerName   string
	HandlerConfig json.RawMessage
	Members       []NodeConfig
	F             int
	MessageType   string
}

type NodeConfig struct {
	Name    string
	Address string
}
