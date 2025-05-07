package operation

import (
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
)

type AptosDeps struct {
	AB               *deployment.AddressBookMap
	AptosChain       deployment.AptosChain
	CCIPOnChainState changeset.CCIPOnChainState
}
