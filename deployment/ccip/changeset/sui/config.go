package sui

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ccip_ops "github.com/smartcontractkit/chainlink-sui/ops/ccip"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
)

type SuiChainDefinition struct {
	// ConnectionConfig holds configuration for connection.
	v1_6.ConnectionConfig `json:"connectionConfig"`
	// Selector is the chain selector of this chain.
	Selector uint64 `json:"selector"`
	// GasPrice defines the USD price (18 decimals) per unit gas for this chain as a destination.
	GasPrice *big.Int `json:"gasPrice"`
}

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

type DeployDummyRecieverConfig struct {
	ChainSelector uint64
	CCIPPackageId string
	McmsPackageId string
	McmsOwner     string
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

type ApplyFeeTokenUpdateConfig struct {
	ChainSelector     uint64
	CCIPPackageId     string
	StateObjectId     string
	OwnerCapObjectId  string
	FeeTokensToRemove []string
	FeeTokensToAdd    []string
}

type ApplyPremiumMultiplierWeiPerEthConfig struct {
	ChainSelector              uint64
	CCIPPackageId              string
	StateObjectId              string
	OwnerCapObjectId           string
	Tokens                     []string
	PremiumMultiplierWeiPerEth []uint64
}

// ChainContractParams stores configuration to call initialize in CCIP contracts
type ChainContractParams struct {
	DestChainSelector uint64
	FeeQuoterParams   ccip_ops.InitFeeQuoterInput
	OffRampParams     OffRampParams
	OnRampParams      OnRampParams
}

type DeploySuiBurnMintTpConfig struct {
	ChainSelector       uint64
	RemoteChainSelector uint64
	EVMToken            common.Address
	EVMTokenPool        common.Address
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
