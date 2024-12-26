package changeset

import (
	"fmt"

	ag_solanago "github.com/gagliardetto/solana-go"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink/deployment"
)

// TODO: Solana re-write
// common.Address used
// LoadOnchainState used which wont work for solana
// deployPrerequisiteContracts needs to be re-written for solana
// basically everything

// DeployPrerequisites deploys the pre-requisite contracts for CCIP
// pre-requisite contracts are the contracts which can be reused from previous versions of CCIP
// Or the contracts which are already deployed on the chain ( for example, tokens, feeds, etc)
// Caller should update the environment's address book with the returned addresses.
func DeployPrerequisitesSolana(env deployment.Environment, cfg DeployPrerequisiteConfig) (deployment.ChangesetOutput, error) {
	err := cfg.Validate()
	if err != nil {
		return deployment.ChangesetOutput{}, errors.Wrapf(deployment.ErrInvalidConfig, "%v", err)
	}
	ab := deployment.NewMemoryAddressBook()
	err = deployPrerequisiteChainContractsSolana(env, ab, cfg)
	if err != nil {
		env.Logger.Errorw("Failed to deploy prerequisite contracts", "err", err, "addressBook", ab)
		return deployment.ChangesetOutput{
			AddressBook: ab,
		}, fmt.Errorf("failed to deploy prerequisite contracts: %w", err)
	}
	return deployment.ChangesetOutput{
		Proposals:   []timelock.MCMSWithTimelockProposal{},
		AddressBook: ab,
		JobSpecs:    nil,
	}, nil
}

func deployPrerequisiteChainContractsSolana(e deployment.Environment, ab deployment.AddressBook, cfg DeployPrerequisiteConfig) error {
	state, err := LoadOnchainStateSolana(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err")
		return err
	}
	deployGrp := errgroup.Group{}
	for _, c := range cfg.Configs {
		chain := e.SolChains[c.ChainSelector]
		deployGrp.Go(func() error {
			err := deployPrerequisiteContractsSolana(e, ab, state, chain)
			if err != nil {
				e.Logger.Errorw("Failed to deploy prerequisite contracts", "chain", chain.String(), "err", err)
				return err
			}
			return nil
		})
	}
	return deployGrp.Wait()
}

// deployPrerequisiteContracts deploys the contracts that can be ported from previous CCIP version to the new one.
// This is only required for staging and test environments where the contracts are not already deployed.
func deployPrerequisiteContractsSolana(e deployment.Environment, ab deployment.AddressBook, state CCIPOnChainState, chain deployment.SolChain) error {
	lggr := e.Logger
	chainState, chainExists := state.SolChains[chain.Selector]
	var ccipRouter *ag_solanago.PublicKey
	var ccipReceiver *ag_solanago.PublicKey
	var tokenPool *ag_solanago.PublicKey
	if chainExists {
		ccipRouter = chainState.CcipRouter
		ccipReceiver = chainState.CcipReceiver
		tokenPool = chainState.TokenPool
	}
	if ccipRouter == nil {

	} else {
		lggr.Infow("ccipRouter already deployed", "chain", chain.String(), "addr", ccipRouter)
	}
	if ccipReceiver == nil {

	} else {
		lggr.Infow("ccipReceiver already deployed", "chain", chain.String(), "addr", ccipReceiver)
	}
	if tokenPool == nil {

	} else {
		lggr.Infow("tokenPool already deployed", "chain", chain.String(), "addr", tokenPool)
	}
	return nil
}
