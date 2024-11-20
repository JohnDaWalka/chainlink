package changeset

import (
	"fmt"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/deployment/ccip"
)

var InitialAddChainCS deployment.ChangeSet[ccipdeployment.InitialAddChainConfig] = InitialAddChain

// InitialAddChain enables new chains as destination for CCIP
// It performs the following steps:
// - AddChainConfig + AddDON (candidate->primary promotion i.e. init) on the home chain
// - SetOCR3Config on the remote chain
// InitialAddChain assumes that the home chain is already enabled and all CCIP contracts are already deployed.
func InitialAddChain(env deployment.Environment, c ccipdeployment.InitialAddChainConfig) (deployment.ChangesetOutput, error) {
	if err := c.Validate(); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid InitialAddChainConfig: %w", err)
	}
	newAddresses := deployment.NewMemoryAddressBook()
	err := ccipdeployment.InitialAddChain(env, newAddresses, c)
	if err != nil {
		env.Logger.Errorw("Failed to deploy CCIP contracts", "err", err, "newAddresses", newAddresses)
		return deployment.ChangesetOutput{AddressBook: newAddresses}, deployment.MaybeDataErr(err)
	}
	return deployment.ChangesetOutput{
		Proposals:   []timelock.MCMSWithTimelockProposal{},
		AddressBook: newAddresses,
		JobSpecs:    nil,
	}, nil
}
