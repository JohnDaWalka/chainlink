package ccip_attestation

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var (
	EVMSignerRegistrySetNewSignerAddressesChangeset = cldf.CreateChangeSet(signerRegistrySetNewSignerAddressesLogic, signerRegistrySetNewSignerAddressesPrecondition)
)

type SetNewSignerAddressesConfig struct {
	// MCMS defines the delay to use for Timelock (if absent, the changeset will attempt to use the deployer key).
	MCMS *proposalutils.TimelockConfig
	// UpdatesByChain maps chain selector -> (existing signer -> new signer) for per-chain updates.
	SignersAddress []SignerRegistryAddress
}

type SignerRegistryAddress struct {
	ExistingSigner common.Address
	NewSigner      common.Address
}

func signerRegistrySetNewSignerAddressesPrecondition(env cldf.Environment, config SetNewSignerAddressesConfig) error {
	if len(config.SignersAddress) == 0 {
		return errors.New("no signer updates provided")
	}
	selector := BaseSepoliaSelector
	if env.Name == "mainnet" {
		selector = BaseMainnetSelector
	}
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState := state.Chains[selector]
	existingSigners := make(map[common.Address]bool)
	// TODO: Check if SignerRegistrySigners is populated in stateview
	for _, signer := range chainState.SignerRegistrySigners {
		existingSigners[signer.EvmAddress] = true

		if signer.NewEVMAddress != utils.ZeroAddress {
			existingSigners[signer.NewEVMAddress] = true
		}
	}

	for _, signer := range config.SignersAddress {
		if !existingSigners[signer.ExistingSigner] {
			return fmt.Errorf("address %s is not a registered signer on chain selector %d", signer.ExistingSigner.Hex(), selector)
		}
		if signer.NewSigner != utils.ZeroAddress {
			if existingSigners[signer.NewSigner] {
				return fmt.Errorf("new address %s is already a signer or pending new address on chain selector %d", signer.NewSigner.Hex(), selector)
			}
		}
	}

	return nil
}

func signerRegistrySetNewSignerAddressesLogic(env cldf.Environment, config SetNewSignerAddressesConfig) (cldf.ChangesetOutput, error) {
	selector := BaseSepoliaSelector
	if env.Name == "mainnet" {
		selector = BaseMainnetSelector
	}
	// Load onchain state to get MCMS addresses if needed
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	deployerGroup := deployergroup.NewDeployerGroup(env, state, config.MCMS).WithDeploymentContext("configure signer registry with new signer addresses")
	chainState := state.Chains[selector]
	opts, err := deployerGroup.GetDeployer(selector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer for chain selector %d: %w", selector, err)
	}

	signerRegistry := chainState.SignerRegistry
	// Prepare arrays for the contract call
	var existingAddresses []common.Address
	var newAddresses []common.Address
	for _, signer := range config.SignersAddress {
		existingAddresses = append(existingAddresses, signer.ExistingSigner)
		newAddresses = append(newAddresses, signer.NewSigner)
	}
	if len(existingAddresses) == 0 || len(newAddresses) == 0 {
		return cldf.ChangesetOutput{}, errors.New("no signer updates provided")
	}
	env.Logger.Infof("Setting new signer address. Existing addresses %v and new addresses %v", existingAddresses, newAddresses)
	_, err = signerRegistry.SetNewSignerAddresses(opts, existingAddresses, newAddresses)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to set new signer addresses on chain selector %d: %w", selector, err)
	}

	return deployerGroup.Enact()
}
