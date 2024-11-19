package changeset

import (
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/deployment/ccip"
)

var _ deployment.ChangeSet[ccipdeployment.InitialAddChainConfig] = InitialAddChain

func InitialAddChain(env deployment.Environment, c ccipdeployment.InitialAddChainConfig) (deployment.ChangesetOutput, error) {
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
