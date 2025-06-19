package globals

import cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

type ConfigType string

const (
	MaxOracles                          = 31
	KeystoneForwarder cldf.ContractType = "KeystoneForwarder" // https://github.com/smartcontractkit/chainlink/blob/50c1b3dbf31bd145b312739b08967600a5c67f30/contracts/src/v0.8/keystone/KeystoneForwarder.sol#L90
)
