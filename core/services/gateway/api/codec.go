package api

import "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"

// Codec implements (de)serialization of Message objects.
type Codec interface {
	DecodeRequest(msgBytes []byte) (*gateway.Message, error)

	EncodeRequest(msg *gateway.Message) ([]byte, error)

	DecodeResponse(msgBytes []byte) (*gateway.Message, error)

	EncodeResponse(msg *gateway.Message) ([]byte, error)

	EncodeNewErrorResponse(id string, code int, message string, data []byte) ([]byte, error)
}
