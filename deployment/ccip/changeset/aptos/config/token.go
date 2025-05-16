package config

import "github.com/aptos-labs/aptos-go-sdk"

type AcceptTokenOwnershipConfig struct {
	ChainSelector uint64
	TokenAddress  aptos.AccountAddress
}
