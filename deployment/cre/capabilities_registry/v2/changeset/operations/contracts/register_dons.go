package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
)

type RegisterDonsDeps struct {
	Env *cldf.Environment
}

type RegisterDonsInput struct {
	Address       string
	ChainSelector uint64
	DONs          []capabilities_registry_v2.CapabilitiesRegistryNewDONParams
}

type RegisterDonsOutput struct {
	DONs []capabilities_registry_v2.CapabilitiesRegistryDONInfo
}

// RegisterDons is an operation that registers DONs in the V2 Capabilities Registry contract.
var RegisterDons = operations.NewOperation[RegisterDonsInput, RegisterDonsOutput, RegisterDonsDeps](
	"register-dons-op",
	semver.MustParse("1.0.0"),
	"Register DONs in Capabilities Registry",
	func(b operations.Bundle, deps RegisterDonsDeps, input RegisterDonsInput) (RegisterDonsOutput, error) {
		// Get the target chain
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return RegisterDonsOutput{}, fmt.Errorf("chain not found for selector %d", input.ChainSelector)
		}

		capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(common.HexToAddress(input.Address), chain.Client)
		if err != nil {
			return RegisterDonsOutput{}, fmt.Errorf("failed to create CapabilitiesRegistry: %w", err)
		}
		/*
			// Get the CapabilitiesRegistryTransactor contract
			capabilityRegistryTransactor, err := capabilities_registry_v2.NewCapabilitiesRegistryTransactor(
				common.HexToAddress(input.Address),
				chain.Client,
			)
			if err != nil {
				return RegisterDonsOutput{}, fmt.Errorf("failed to create CapabilitiesRegistryTransactor: %w", err)
			}
		*/
		tx, err := capReg.AddDONs(chain.DeployerKey, input.DONs)
		if err != nil {
			err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
			return RegisterDonsOutput{}, fmt.Errorf("failed to call AddDONs: %w", err)
		}

		_, err = chain.Confirm(tx)
		if err != nil {
			return RegisterDonsOutput{}, fmt.Errorf("failed to confirm AddDONs transaction %s: %w", tx.Hash().String(), err)
		}

		resp := RegisterDonsOutput{}

		// Get the CapabilitiesRegistryCaller contract
		var donInfo []capabilities_registry_v2.CapabilitiesRegistryDONInfo
		for _, d := range input.DONs {
			di, err := capReg.GetDONByName(&bind.CallOpts{}, d.Name)
			if err != nil {
				return RegisterDonsOutput{}, fmt.Errorf("failed to get DONs: %w", err)
			}
			allDons, err := capReg.GetDONs(&bind.CallOpts{})
			if err != nil {
				return RegisterDonsOutput{}, fmt.Errorf("failed to get all DONs: %w", err)
			}
			fmt.Printf("All DONs after registering %s: %+v\n", d.Name, allDons)
			n, _ := capReg.GetNextDONId(&bind.CallOpts{})
			fmt.Printf("Next DON ID: %d\n", n)
			donInfo = append(donInfo, di)
		}

		resp.DONs = donInfo

		return resp, nil
	},
)
