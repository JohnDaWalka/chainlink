package changeset

import (
	"github.com/smartcontractkit/chainlink/deployment"

	ccipdeployment "github.com/smartcontractkit/chainlink/deployment/ccip"
)

var _ deployment.ChangeSet[ccipdeployment.InitialAddChainConfig] = InitialDeploy

func InitialDeploy(env deployment.Environment, c ccipdeployment.InitialAddChainConfig) (deployment.ChangesetOutput, error) {
	newAddresses := deployment.NewMemoryAddressBook()
	err := ccipdeployment.InitialAddChain(env, newAddresses, c)
	if err != nil {
		env.Logger.Errorw("Failed to deploy initial chain", "err", err, "addressBook", newAddresses)
		return deployment.ChangesetOutput{
			AddressBook: newAddresses,
		}, err
	}
	return deployment.ChangesetOutput{
		AddressBook: newAddresses,
	}, nil
}
