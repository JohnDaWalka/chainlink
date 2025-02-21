package config

import chain_selectors "github.com/smartcontractkit/chain-selectors"

type EthKeyConfig interface {
	//Address() string
	ChainDetails() chain_selectors.ChainDetails
	// JSON must be a valid JSON string conforming to the
	// the geth keystore format.
	JSON() string
	// Password is the password used to encrypt the key.
	Password() string
}
