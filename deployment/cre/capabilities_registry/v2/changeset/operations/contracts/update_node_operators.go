package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	mcmslib "github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
)

type UpdateNopsDeps struct {
	Env           *cldf.Environment
	MCMSContracts *commonchangeset.MCMSWithTimelockState // Required if MCMSConfig is not nil
}

type UpdateNopsInput struct {
	Address       string
	ChainSelector uint64
	Nops          map[uint32]capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams
	MCMSConfig    *ocr3.MCMSConfig
}

type UpdateNopsOutput struct {
	Proposals []mcmslib.TimelockProposal
}

// UpdateNops is an operation that updates node operators in the V2 Capabilities Registry contract.
var UpdateNops = operations.NewOperation[UpdateNopsInput, UpdateNopsOutput, UpdateNopsDeps](
	"update-nops-op",
	semver.MustParse("1.0.0"),
	"Update Node Operators in Capabilities Registry",
	func(b operations.Bundle, deps UpdateNopsDeps, input UpdateNopsInput) (UpdateNopsOutput, error) {
		// Get the target chain
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return UpdateNopsOutput{}, fmt.Errorf("chain not found for selector %d", input.ChainSelector)
		}

		// Get the CapabilitiesRegistryTransactor contract
		capabilityRegistryTransactor, err := capabilities_registry_v2.NewCapabilitiesRegistry(
			common.HexToAddress(input.Address),
			chain.Client,
		)
		if err != nil {
			return UpdateNopsOutput{}, fmt.Errorf("failed to create CapabilitiesRegistryTransactor: %w", err)
		}

		// Prepare ids and params slices
		var ids []uint32
		var params []capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams
		for id, param := range input.Nops {
			ids = append(ids, id)
			params = append(params, param)
		}

		// Create the appropriate strategy
		strategy, err := strategies.CreateStrategy(
			chain,
			*deps.Env,
			input.MCMSConfig,
			deps.MCMSContracts,
			common.HexToAddress(input.Address),
			"Update Node Operators",
		)
		if err != nil {
			return UpdateNopsOutput{}, fmt.Errorf("failed to create strategy: %w", err)
		}

		// Execute the transaction using the strategy
		proposals, err := strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := capabilityRegistryTransactor.UpdateNodeOperators(opts, ids, params)
			if err != nil {
				err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
				return nil, fmt.Errorf("failed to call UpdateNodeOperators: %w", err)
			}

			// For direct execution, confirm transaction
			if input.MCMSConfig == nil {
				_, err = chain.Confirm(tx)
				if err != nil {
					return nil, fmt.Errorf("failed to confirm UpdateNodeOperators transaction %s: %w", tx.Hash().String(), err)
				}
			}

			return tx, nil
		})
		if err != nil {
			return UpdateNopsOutput{}, fmt.Errorf("failed to execute UpdateNodeOperators: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for UpdateNops on chain %d", input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully updated %d node operators on chain %d", len(ids), input.ChainSelector)
		}

		return UpdateNopsOutput{
			Proposals: proposals,
		}, nil
	},
)
