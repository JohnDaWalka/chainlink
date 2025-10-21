package contracts

import (
	"fmt"
	"math/big"

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

type RegisterNopsDeps struct {
	Env           *cldf.Environment
	MCMSContracts *commonchangeset.MCMSWithTimelockState // Required if MCMSConfig is not nil
}

type RegisterNopsInput struct {
	Address       string
	ChainSelector uint64
	Nops          []capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams
	MCMSConfig    *ocr3.MCMSConfig
}

type RegisterNopsOutput struct {
	Nops      []*capabilities_registry_v2.CapabilitiesRegistryNodeOperatorAdded
	Proposals []mcmslib.TimelockProposal
}

// RegisterNops is an operation that registers node operators in the V2 Capabilities Registry contract.
var RegisterNops = operations.NewOperation[RegisterNopsInput, RegisterNopsOutput, RegisterNopsDeps](
	"register-nops-op",
	semver.MustParse("1.0.0"),
	"Register Node Operators in Capabilities Registry",
	func(b operations.Bundle, deps RegisterNopsDeps, input RegisterNopsInput) (RegisterNopsOutput, error) {
		// Get the target chain
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return RegisterNopsOutput{}, fmt.Errorf("chain not found for selector %d", input.ChainSelector)
		}

		// Get the CapabilitiesRegistryTransactor contract
		capabilityRegistryTransactor, err := capabilities_registry_v2.NewCapabilitiesRegistry(
			common.HexToAddress(input.Address),
			chain.Client,
		)
		if err != nil {
			return RegisterNopsOutput{}, fmt.Errorf("failed to create CapabilitiesRegistryTransactor: %w", err)
		}

		newNops, err := dedupNodeOperators(capabilityRegistryTransactor, input.Nops)
		if err != nil {
			return RegisterNopsOutput{}, fmt.Errorf("failed to deduplicate node operators: %w", err)
		}

		// Create the appropriate strategy
		strategy, err := strategies.CreateStrategy(
			chain,
			*deps.Env,
			input.MCMSConfig,
			deps.MCMSContracts,
			common.HexToAddress(input.Address),
			RegisterNopsDescription,
		)
		if err != nil {
			return RegisterNopsOutput{}, fmt.Errorf("failed to create strategy: %w", err)
		}

		var resultNops []*capabilities_registry_v2.CapabilitiesRegistryNodeOperatorAdded

		// Execute the transaction using the strategy
		proposals, err := strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := capabilityRegistryTransactor.AddNodeOperators(opts, newNops)
			if err != nil {
				err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
				return nil, fmt.Errorf("failed to call AddNodeOperators: %w", err)
			}

			// For direct execution, we can get the receipt and parse logs
			if input.MCMSConfig == nil {
				// Confirm transaction and get receipt
				_, err = chain.Confirm(tx)
				if err != nil {
					return nil, fmt.Errorf("failed to confirm AddNodeOperators transaction %s: %w", tx.Hash().String(), err)
				}

				ctx := b.GetContext()
				receipt, err := bind.WaitMined(ctx, chain.Client, tx)
				if err != nil {
					return nil, fmt.Errorf("failed to mine AddNodeOperators transaction %s: %w", tx.Hash().String(), err)
				}

				// Get the CapabilitiesRegistryFilterer contract for parsing logs
				capabilityRegistryFilterer, err := capabilities_registry_v2.NewCapabilitiesRegistryFilterer(
					common.HexToAddress(input.Address),
					chain.Client,
				)
				if err != nil {
					return nil, fmt.Errorf("failed to create CapabilitiesRegistryFilterer: %w", err)
				}

				// Parse the logs to get the added node operators
				resultNops = make([]*capabilities_registry_v2.CapabilitiesRegistryNodeOperatorAdded, 0, len(receipt.Logs))
				for i, log := range receipt.Logs {
					if log == nil {
						continue
					}

					o, err := capabilityRegistryFilterer.ParseNodeOperatorAdded(*log)
					if err != nil {
						return nil, fmt.Errorf("failed to parse log %d for operator added: %w", i, err)
					}
					resultNops = append(resultNops, o)
				}
			}

			return tx, nil
		})
		if err != nil {
			return RegisterNopsOutput{}, fmt.Errorf("failed to execute AddNodeOperators: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for RegisterNops on chain %d", input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully registered %d node operators on chain %d", len(resultNops), input.ChainSelector)
		}

		return RegisterNopsOutput{
			Nops:      resultNops,
			Proposals: proposals,
		}, nil
	},
)

func dedupNodeOperators(c capabilities_registry_v2.CapabilitiesRegistryInterface, nops []capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams) ([]capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams, error) {
	m := make(map[capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams]struct{})
	for i := range nops {
		m[nops[i]] = struct{}{}
	}
	r, err := c.GetNodeOperators(nil, big.NewInt(1), big.NewInt(256)) // TODO pagination
	if err != nil {
		return nil, fmt.Errorf("failed to get existing node operators: %w", err)
	}

	if len(r) == 0 {
		return nops, nil
	}
	if len(r) >= 256 {
		return nil, fmt.Errorf("too many existing node operators, pagination not implemented")
	}
	for _, nop := range r {
		delete(m, capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams{Admin: nop.Admin, Name: nop.Name})
	}

	// maintain order of input slice
	if len(m) == 0 {
		return []capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams{}, nil
	}
	var nopsToRegister []capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams
	for _, nop := range nops {
		if _, ok := m[nop]; ok {
			nopsToRegister = append(nopsToRegister, nop)
		}
	}
	return nopsToRegister, nil
}
