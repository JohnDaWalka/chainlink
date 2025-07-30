package modsectypes

import (
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/verifier_proxy"
)

type CCIPMessageSent interface {
	MessageID() [32]byte
}

type CCIPMessageParser interface {
	// TODO: EVM-specific, needs to be more generic.
	ParseCCIPMessageSent(log gethtypes.Log) (CCIPMessageSent, error)
}

type evmCCIPMessageSent struct {
	ccipMessageSent verifier_proxy.VerifierProxyCCIPMessageSent
}

func (m *evmCCIPMessageSent) MessageID() [32]byte {
	return m.ccipMessageSent.Message.Header.MessageId
}

type evmMessageParser struct {
}

func NewEVMCCIPMessageParser() CCIPMessageParser {
	return &evmMessageParser{}
}

func (p *evmMessageParser) ParseCCIPMessageSent(log gethtypes.Log) (CCIPMessageSent, error) {
	// Don't actually need to call the contract, just need to parse the log.
	verifierOnramp, err := verifier_proxy.NewVerifierProxy(log.Address, nil)
	if err != nil {
		return nil, err
	}

	parsedLog, err := verifierOnramp.ParseCCIPMessageSent(log)
	if err != nil {
		return nil, err
	}

	return &evmCCIPMessageSent{
		ccipMessageSent: *parsedLog,
	}, nil
}
