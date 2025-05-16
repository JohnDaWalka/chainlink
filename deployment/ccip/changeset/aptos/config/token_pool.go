package config

import (
	"math/big"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/ethereum/go-ethereum/common"
	fee_quoter "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
)

type AddTokenPoolConfig struct {
	// DeployAptosTokenConfig
	ChainSelector                       uint64
	TokenAddress                        aptos.AccountAddress
	TokenObjAddress                     aptos.AccountAddress
	TokenSymbol                         changeset.TokenSymbol
	PoolType                            deployment.ContractType
	TokenTransferFeeByRemoteChainConfig map[uint64]fee_quoter.TokenTransferFeeConfig
	EVMRemoteConfigs                    map[uint64]EVMRemoteConfig
	TokenParams                         TokenParams
}

type TokenParams struct {
	MaxSupply *big.Int
	Name      string
	Symbol    string
	Decimals  byte
	Icon      string
	Project   string
}

// // TODO: gather requirements for Aptos token deployment
// type DeployAptosTokenConfig struct {
// 	TokenDecimals uint8
// 	TokenSymbol   string
// }

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
