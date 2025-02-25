package types

import "github.com/ethereum/go-ethereum/common"

type EVMKeysToChains struct {
	ChainSelector uint64
	ChainName     string
}

type EVMKeys struct {
	EncryptedJSONs  [][]byte
	PublicAddresses []common.Address
	Password        string
	Chains          []EVMKeysToChains
}

type P2PKeys struct {
	EncryptedJSONs [][]byte
	PeerIDs        []string
	PublicHexKeys  []string
	PrivateKeys    []string
	Password       string
}
