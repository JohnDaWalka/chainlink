package config

import (
	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/ethereum/go-ethereum/common"
	fee_quoter "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type AddTokenPoolConfig struct {
	// DeployAptosTokenConfig
	ChainSelector uint64
	TokenAddress  aptos.AccountAddress
	DeployAptosTokenPoolConfig
	TransferAdminRoleConfig
	TokenTransferFeeByRemoteChainConfig map[uint64]fee_quoter.TokenTransferFeeConfig
	AptosMCMS                           *proposalutils.TimelockConfig
	RemoteChainTokenPoolConfig
}

// // TODO: gather requirements for Aptos token deployment
// type DeployAptosTokenConfig struct {
// 	TokenDecimals uint8
// 	TokenSymbol   string
// }

type DeployAptosTokenPoolConfig struct {
	ChainSelector uint64
	PoolType      string // TODO: is there a standard or just string?
}

type TransferAdminRoleConfig struct {
	NewAdminAddress aptos.AccountAddress
}

type RemoteChainTokenPoolConfig struct {
	EVMRemoteConfigs map[uint64]EVMRemoteConfig
	MCMS             *proposalutils.TimelockConfig
	Metadata         string
}

// TODO: use this to get the correct token pool address
// type EVMRemoteConfig struct {
// 	TokenSymbol       changeset.TokenSymbol
// 	PoolType          cldf.ContractType
// 	PoolVersion       semver.Version
// 	RateLimiterConfig RateLimiterConfig
// 	OverrideConfig    bool
// }

type EVMRemoteConfig struct {
	TokenAddress     common.Address
	TokenPoolAddress common.Address
	RateLimiterConfig
}

type RateLimiterConfig struct {
	RemoteChainSelector uint64
	OutboundIsEnabled   bool
	OutboundCapacity    uint64
	OutboundRate        uint64
	InboundIsEnabled    bool
	InboundCapacity     uint64
	InboundRate         uint64
}
