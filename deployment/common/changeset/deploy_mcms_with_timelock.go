package changeset

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/internal/mcmsnew/evm"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

var _ deployment.ChangeSet[map[uint64]types.MCMSWithTimelockConfig] = DeployMCMSWithTimelock
var _ deployment.ChangeSet[map[uint64]types.MCMSWithTimelockConfigV2] = DeployMCMSWithTimelockV2

func DeployMCMSWithTimelock(e deployment.Environment, cfgByChain map[uint64]types.MCMSWithTimelockConfig) (deployment.ChangesetOutput, error) {
    newAddresses := deployment.NewMemoryAddressBook()
    err := internal.DeployMCMSWithTimelockContractsBatch(
        e.Logger, e.Chains, newAddresses, cfgByChain,
    )
    if err != nil {
        return deployment.ChangesetOutput{AddressBook: newAddresses}, err
    }
    return deployment.ChangesetOutput{AddressBook: newAddresses}, nil
}

func DeployMCMSWithTimelockV2(
	e deployment.Environment, cfgByChain map[uint64]types.MCMSWithTimelockConfigV2,
) (deployment.ChangesetOutput, error) {
	newAddresses := deployment.NewMemoryAddressBook()

	for chainSelector, config := range cfgByChain {
		family, _ := chainselectors.GetSelectorFamily(chainSelector)
		switch family {
		case chainselectors.FamilyEVM:
			chain, found := e.Chains[chainSelector]
			if !found {
				err := fmt.Errorf("unable to find chain for selector %d", chainSelector)
				return deployment.ChangesetOutput{AddressBook: newAddresses}, err
			}

			_, err := evm.DeployMCMSWithTimelockContractsEVM(e.Logger, chain, newAddresses, config)
			if err != nil {
				return deployment.ChangesetOutput{AddressBook: newAddresses}, err
			}

		case chainselectors.FamilySolana:
			chain, found := e.SolChains[chainSelector]
			if !found {
				err := fmt.Errorf("unable to find chain for selector %d", chainSelector)
				return deployment.ChangesetOutput{AddressBook: newAddresses}, err
			}

			state, err := LoadOnchainStateSolana(e)
			if err != nil {
				e.Logger.Errorw("Failed to load existing onchain state", "err", err)
				return deployment.ChangesetOutput{AddressBook: newAddresses}, err
			}

			_, err = DeployMCMSWithTimelockContractsSolana(e, state, chain, newAddresses, config)
			if err != nil {
				return deployment.ChangesetOutput{AddressBook: newAddresses}, err
			}

		default:
			return deployment.ChangesetOutput{AddressBook: newAddresses}, fmt.Errorf("unsupported chain family: %s", family)
		}
	}

	return deployment.ChangesetOutput{AddressBook: newAddresses}, nil
}

func ValidateOwnership(ctx context.Context, mcms bool, deployerKey, timelock common.Address, contract Ownable) error {
	owner, err := contract.Owner(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to get owner: %w", err)
	}
	if mcms && owner != timelock {
		return fmt.Errorf("%s not owned by deployer key", contract.Address())
	} else if !mcms && owner != deployerKey {
		return fmt.Errorf("%s not owned by deployer key", contract.Address())
	}
	return nil
}

// TODO: SOLANA_CCIP
func ValidateOwnershipSolana(ctx context.Context, mcms bool, deployerKey, timelock, ccipRouter solana.PublicKey) error {
	return nil
}
