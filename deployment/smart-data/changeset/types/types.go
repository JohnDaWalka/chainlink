package types

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type MCMSConfig struct {
	MinDelay time.Duration // delay for timelock worker to execute the transfers.
}

type SetForwarderConfig struct {
	ForwarderAddress common.Address
	ChainSelector    uint64
	DonID            uint32
	ConfigVersion    uint32
	F                uint8
	Signers          []common.Address
	McmsConfig       *MCMSConfig
}
