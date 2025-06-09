package sui

import "github.com/pattonkan/sui-go/sui"

type DeploySuiChainConfig struct {
	ContractParamsPerChain map[uint64]ChainContractParams
}

// ChainContractParams stores configuration to call initialize in CCIP contracts
type ChainContractParams struct {
	FeeQuoterParams FeeQuoterParams
	OffRampParams   OffRampParams
	OnRampParams    OnRampParams
}

type FeeQuoterParams struct {
	MaxFeeJuelsPerMsg            uint64
	LinkToken                    sui.Address
	TokenPriceStalenessThreshold uint64
	FeeTokens                    []sui.Address
}

type OffRampParams struct {
	ChainSelector                    uint64
	PermissionlessExecutionThreshold uint32
	IsRMNVerificationDisabled        []bool
	SourceChainSelectors             []uint64
	SourceChainIsEnabled             []bool
	SourceChainsOnRamp               [][]byte
}

type OnRampParams struct {
	ChainSelector  uint64
	AllowlistAdmin sui.Address
	FeeAggregator  sui.Address
}
