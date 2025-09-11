package lane

import (
	"math/big"
)

type ChainDefinition struct {
	// ConnectionConfig holds configuration for connection.
	ConnectionConfig `json:"connectionConfig"`
	// Selector is the chain selector of this chain.
	Selector uint64 `json:"selector"`
	// GasPrice defines the USD price (18 decimals) per unit gas for this chain as a destination.
	GasPrice *big.Int `json:"gasPrice"`
	// TokenPrices define the USD price (18 decimals) per 1e18 of the smallest token denomination for various tokens on this chain.
	TokenPrices map[string]*big.Int `json:"tokenPrices"`
	// FeeQuoterDestChainConfig is the configuration on a fee quoter for this chain as a destination.
	FeeQuoterDestChainConfig FeeQuoterDestChainConfig `json:"feeQuoterDestChainConfig"`
}

type ConnectionConfig struct {
	// RMNVerificationDisabled is true if we do not want the RMN to bless messages FROM this chain.
	RMNVerificationDisabled bool `json:"rmnVerificationDisabled"`
	// AllowListEnabled is true if we want an allowlist to dictate who can send messages TO this chain.
	AllowListEnabled bool `json:"allowListEnabled"`
}

type FeeQuoterDestChainConfig struct {
	IsEnabled                         bool
	MaxNumberOfTokensPerMsg           uint16
	MaxDataBytes                      uint32
	MaxPerMsgGasLimit                 uint32
	DestGasOverhead                   uint32
	DestGasPerPayloadByteBase         uint8
	DestGasPerPayloadByteHigh         uint8
	DestGasPerPayloadByteThreshold    uint16
	DestDataAvailabilityOverheadGas   uint32
	DestGasPerDataAvailabilityByte    uint16
	DestDataAvailabilityMultiplierBps uint16
	ChainFamilySelector               [4]byte
	EnforceOutOfOrder                 bool
	DefaultTokenFeeUSDCents           uint16
	DefaultTokenDestGasOverhead       uint32
	DefaultTxGasLimit                 uint32
	GasMultiplierWeiPerEth            uint64
	GasPriceStalenessThreshold        uint32
	NetworkFeeUSDCents                uint32
}
