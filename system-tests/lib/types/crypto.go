package types

import "github.com/ethereum/go-ethereum/common"

type EVMKeys struct {
	EncryptedJSONs  [][]byte
	PublicAddresses []common.Address
	Password        string
	ChainSelector   uint64
}

type P2PKeys struct {
	EncryptedJSONs [][]byte
	PeerIDs        []string
	PublicHexKeys  []string
	PrivateKeys    []string
	Password       string
}
