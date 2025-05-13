package operation

import (
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
)

type AptosDeps struct {
	AB         *deployment.AddressBookMap
	AptosChain deployment.AptosChain
	// TODO: Refactor this?
	Env              deployment.Environment
	OnChainState     changeset.AptosCCIPChainState
	CCIPOnChainState changeset.CCIPOnChainState
}
