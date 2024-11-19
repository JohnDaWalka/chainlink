package changeset

import (
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/deployment/ccip"
)

var _ deployment.ChangeSet[InitialAddChainConfig] = InitialAddChain

func InitialAddChain(env deployment.Environment, c InitialAddChainConfig) (deployment.ChangesetOutput, error) {
	newAddresses := deployment.NewMemoryAddressBook()
	err := ccipdeployment.DeployChainContractsForChains(env, newAddresses, c.HomeChainSelector, c.ChainSelectors, c.MCMSCfg)
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

type InitialAddChainConfig struct {
	ChainSelectors    []uint64
	HomeChainSelector uint64
	FeedChainSel      uint64
	TokenConfig       TokenConfig
	// I believe it makes sense to have the same signers across all chains
	// since that's the point MCMS.
	MCMSConfig MCMSConfig
	USDCConfig USDCConfig
	// For setting OCR configuration
	OCRSecrets deployment.OCRSecrets
}
