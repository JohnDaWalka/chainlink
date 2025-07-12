package sui

import (
	"math/big"

	ccip_ops "github.com/smartcontractkit/chainlink-sui/ops/ccip"
)

type DeploySuiChainConfig struct {
	ContractParamsPerChain map[uint64]ChainContractParams
}

type MintSuiTokenConfig struct {
	ChainSelector  uint64
	TokenPackageId string
	TreasuryCapId  string
	Amount         uint64
}

type DeployLinkTokenConfig struct {
	ChainSelector uint64
}

type UpdateSuiPriceConfig struct {
	ChainSelector       uint64
	DestChainSelector   []uint64
	CCIPPackageId       string
	CCIPObjectRef       string
	SourceTokenMetadata []string
	SourceUsdPerToken   []*big.Int
	GasUsdPerUnitGas    []*big.Int
}

// ChainContractParams stores configuration to call initialize in CCIP contracts
type ChainContractParams struct {
	DestChainSelector uint64
	FeeQuoterParams   ccip_ops.InitFeeQuoterInput
	OffRampParams     OffRampParams
	OnRampParams      OnRampParams
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
	AllowlistAdmin string
	FeeAggregator  string
}
