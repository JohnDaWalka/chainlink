package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	mcmslib "github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/deployment/smart-data/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/smart-data/changeset/types"
)

var SetForwarderConfigChangeset = cldf.CreateChangeSet(setForwarderConfigLogic, setForwarderConfigPrecondition)

func setForwarderConfigPrecondition(env cldf.Environment, c types.SetForwarderConfig) error {
	_, ok := env.BlockChains.EVMChains()[c.ChainSelector]
	if !ok {
		return fmt.Errorf("chain not found in env %d", c.ChainSelector)
	}

	if c.F == 0 {
		return errors.New("f tolerance must be positive")
	}
	if len(c.Signers) > globals.MaxOracles {
		return errors.New("signers excess")
	}

	if len(c.Signers) <= int(c.F)*3 {
		return errors.New("insufficient signers")
	}

	seen := make(map[common.Address]struct{})
	for _, signer := range c.Signers {
		if signer == utils.ZeroAddress {
			return errors.New("invalid signer zero address")
		}

		if _, exists := seen[signer]; exists {
			return fmt.Errorf("duplicated signer: %s", signer.String())
		}
		seen[signer] = struct{}{}
	}

	if c.McmsConfig != nil {
		if err := ValidateMCMSAddresses(env.ExistingAddresses, c.ChainSelector); err != nil {
			return err
		}
	}

	return ValidateForwarderForChain(env, c.ChainSelector, c.ForwarderAddress)
}

func setForwarderConfigLogic(env cldf.Environment, c types.SetForwarderConfig) (cldf.ChangesetOutput, error) {
	evmChains := env.BlockChains.EVMChains()
	chain := evmChains[c.ChainSelector]

	contractSetsResp, err := GetContractSets(env.Logger, &GetContractSetsRequest{
		Chains:      evmChains,
		AddressBook: env.ExistingAddresses,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get contract sets: %w", err)
	}

	chainSets := contractSetsResp.ContractSets[c.ChainSelector]
	contract := chainSets.Forwarder[c.ForwarderAddress]

	txOpt := chain.DeployerKey
	if c.McmsConfig != nil {
		txOpt = cldf.SimTransactOpts()
	}

	tx, err := contract.SetConfig(txOpt, c.DonID, c.ConfigVersion, c.F, c.Signers)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed building SetForwarderConfig txs: %w", err)
	}

	if c.McmsConfig != nil {
		proposals := MultiChainProposalConfig{
			c.ChainSelector: []ProposalData{
				{
					contract: contract.Address().Hex(),
					tx:       tx,
				},
			},
		}
		proposal, err := BuildMultiChainProposals(env, "proposal to set config on a forwarder", proposals, c.McmsConfig.MinDelay)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
	}

	if _, err := cldf.ConfirmIfNoError(chain, tx, err); err != nil {
		if tx != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm transaction: %s, %w", tx.Hash().String(), err)
		}

		return cldf.ChangesetOutput{}, fmt.Errorf("failed to submit transaction: %w", err)
	}

	return cldf.ChangesetOutput{}, nil
}
