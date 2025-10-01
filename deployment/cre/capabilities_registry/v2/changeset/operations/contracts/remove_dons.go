package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	mcmslib "github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
)

type RemoveDONsDeps struct {
	Env                  *cldf.Environment
	CapabilitiesRegistry *capabilities_registry_v2.CapabilitiesRegistry

	MCMSContracts *commonchangeset.MCMSWithTimelockState // Required if MCMSConfig is not nil
}

type RemoveDONsInput struct {
	Domain        string
	Zone          string
	DONNames      []string
	ChainSelector uint64

	MCMSConfig *ocr3.MCMSConfig
}

type RemoveDONsOutput struct {
	Proposals []mcmslib.TimelockProposal
}

var RemoveDONs = operations.NewOperation[RemoveDONsInput, RemoveDONsOutput, RemoveDONsDeps](
	"remove-dons-op",
	semver.MustParse("1.0.0"),
	"Remove DONs from Capabilities Registry",
	func(b operations.Bundle, deps RemoveDONsDeps, input RemoveDONsInput) (RemoveDONsOutput, error) {
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return RemoveDONsOutput{}, cldf.ErrChainNotFound
		}

		strategy, err := strategies.CreateStrategy(
			chain,
			*deps.Env,
			input.MCMSConfig,
			deps.MCMSContracts,
			deps.CapabilitiesRegistry.Address(),
			RemoveDONsDescription,
		)
		if err != nil {
			return RemoveDONsOutput{}, fmt.Errorf("failed to create strategy: %w", err)
		}

		proposals, err := strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := deps.CapabilitiesRegistry.RemoveDONsByName(opts, input.DONNames)
			if err != nil {
				err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
				return nil, fmt.Errorf("failed to call RemoveDONsByName: %w", err)
			}

			// For direct execution, we can confirm and get the updated DON info
			if input.MCMSConfig == nil {
				// Confirm transaction
				_, err = chain.Confirm(tx)
				if err != nil {
					return nil, fmt.Errorf("failed to confirm RemoveDONsByName transaction %s: %w", tx.Hash().String(), err)
				}

				ctx := b.GetContext()
				_, err = bind.WaitMined(ctx, chain.Client, tx)
				if err != nil {
					return nil, fmt.Errorf("failed to mine RemoveDONsByName transaction %s: %w", tx.Hash().String(), err)
				}
			}

			return tx, nil
		})
		if err != nil {
			return RemoveDONsOutput{}, fmt.Errorf("failed to remove DONs: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for RemoveDONs '%s' on chain %d", input.DONNames, input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully remove DONs '%s' on chain %d", input.DONNames, input.ChainSelector)
		}

		return RemoveDONsOutput{
			Proposals: proposals,
		}, nil
	},
)
