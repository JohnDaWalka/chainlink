package changeset

import (
	"fmt"
	"math/big"

	binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	accessControllerBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/access_controller"
	mcmBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/mcm"
	timelockBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/timelock"
	solanaUtils "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

var (
	MCM              deployment.ContractType = "MCM"
	Timelock         deployment.ContractType = "Timelock"
	AccessController deployment.ContractType = "AccessController"
)

func DeployMCMSWithTimelockContractsSolana(
	e deployment.Environment,
	state MCMSSolanaState,
	chain deployment.SolChain,
	addressBook deployment.AddressBook,
	config types.MCMSWithTimelockConfigV2,
) (any, error) { // FIXME: define return type
	err := deployMCMSSolana(e, state, chain, addressBook, config)
	if err != nil {
		return nil, fmt.Errorf("unable to deploy mcms contract: %w", err)
	}

	err = deployAccessControllerSolana(e, state, chain, addressBook, config)
	if err != nil {
		return nil, fmt.Errorf("unable to deploy access controller contract: %w", err)
	}

	err = deployTimelockSolana(e, state, chain, addressBook, config)
	if err != nil {
		return nil, fmt.Errorf("unable to deploy timelock contract: %w", err)
	}

	err = setupRolesAndOwnership(e, state, chain, addressBook)
	if err != nil {
		return nil, fmt.Errorf("unable to setup mcms and timelock roles and ownership: %w", err)
	}

	return nil, nil
}

func deployMCMSSolana(
	e deployment.Environment, state MCMSSolanaState, chain deployment.SolChain, addressBook deployment.AddressBook,
	_ types.MCMSWithTimelockConfigV2,
) error {
	var mcmProgram solana.PublicKey
	if state.MCM.IsZero() {
		programID, err := chain.DeployProgram(e.Logger, "mcm")
		if err != nil {
			return fmt.Errorf("unable to deploy mcm program: %w", err)
		}

		typeAndVersion := deployment.NewTypeAndVersion(MCM, deployment.Version1_0_0)
		e.Logger.Infow("Deployed contract", "Contract", typeAndVersion.String(), "addr", programID, "chain", chain.String())

		mcmProgram = solana.MustPublicKeyFromBase58(programID)

		err = addressBook.Save(chain.Selector, programID, typeAndVersion)
		if err != nil {
			return fmt.Errorf("unable to save address: %w", err)
		}
	} else {
		e.Logger.Infow("Using existing MCM program", "addr", state.MCM.String())
		mcmProgram = state.MCM
	}

	mcmBindings.SetProgramID(mcmProgram)

	err := initializeMCM(e, chain, mcmProgram)
	if err != nil {
		return fmt.Errorf("unable to initialize mcm: %w", err)
	}

	// FIXME: review if we need to setup an "AddressLookupTable".

	return nil
}

func deployTimelockSolana(
	e deployment.Environment, state MCMSSolanaState, chain deployment.SolChain, addressBook deployment.AddressBook,
	config types.MCMSWithTimelockConfigV2,
) error {
	var timelockProgram solana.PublicKey
	if state.Timelock.IsZero() {
		programID, err := chain.DeployProgram(e.Logger, "timelock")
		if err != nil {
			return fmt.Errorf("unable to deploy timelock program: %w", err)
		}

		typeAndVersion := deployment.NewTypeAndVersion(Timelock, deployment.Version1_0_0)
		e.Logger.Infow("Deployed contract", "Contract", typeAndVersion.String(), "addr", programID, "chain", chain.String())

		timelockProgram = solana.MustPublicKeyFromBase58(programID)

		err = addressBook.Save(chain.Selector, programID, typeAndVersion)
		if err != nil {
			return fmt.Errorf("unable to save address: %w", err)
		}
	} else {
		e.Logger.Infow("Using existing Timelock program", "addr", state.Timelock.String())
		timelockProgram = state.Timelock
	}

	timelockBindings.SetProgramID(timelockProgram)

	err := initializeTimelock(e, chain, timelockProgram, config.TimelockMinDelay)
	if err != nil {
		return fmt.Errorf("unable to initialize timelock: %w", err)
	}

	// FIXME: review if we need to setup an "AddressLookupTable".

	return nil
}

func deployAccessControllerSolana(
	e deployment.Environment, state MCMSSolanaState, chain deployment.SolChain, addressBook deployment.AddressBook,
	_ types.MCMSWithTimelockConfigV2,
) error {
	var accessControllerProgram solana.PublicKey
	if state.AccessController.IsZero() {
		programID, err := chain.DeployProgram(e.Logger, "access_controller")
		if err != nil {
			return fmt.Errorf("unable to deploy access controller program: %w", err)
		}

		typeAndVersion := deployment.NewTypeAndVersion(AccessController, deployment.Version1_0_0)
		e.Logger.Infow("Deployed contract", "Contract", typeAndVersion.String(), "addr", programID, "chain", chain.String())

		accessControllerProgram = solana.MustPublicKeyFromBase58(programID)

		err = addressBook.Save(chain.Selector, programID, typeAndVersion)
		if err != nil {
			return fmt.Errorf("unable to save address: %w", err)
		}
	} else {
		e.Logger.Infow("Using existing AccessController program", "addr", state.AccessController.String())
		accessControllerProgram = state.AccessController
	}

	accessControllerBindings.SetProgramID(accessControllerProgram)

	err := initializeAccessController(e, chain, accessControllerProgram)
	if err != nil {
		return fmt.Errorf("unable to initialize timelock: %w", err)
	}

	// FIXME: review if we need to setup an "AddressLookupTable".

	return nil
}

func setupRolesAndOwnership(
	e deployment.Environment, state MCMSSolanaState, chain deployment.SolChain, addressBook deployment.AddressBook,
) error {
	return fmt.Errorf("unimplemented")
}

func initializeMCM(e deployment.Environment, chain deployment.SolChain, mcmProgram solana.PublicKey) error {
	multisigID := [32]byte{} // FIXME: where should this come from?

	var mcmConfig mcmBindings.MultisigConfig
	err := chain.GetAccountDataBorshInto(e.GetContext(), GetMCMConfigPDA(mcmProgram, multisigID), &mcmConfig)
	if err == nil {
		e.Logger.Infow("MCM already initialized, skipping initialization", "chain", chain.String())
		return nil
	}

	var programData struct {
		DataType uint32
		Address  solana.PublicKey
	}
	opts := &rpc.GetAccountInfoOpts{Commitment: rpc.CommitmentConfirmed}

	data, err := chain.Client.GetAccountInfoWithOpts(e.GetContext(), mcmProgram, opts)
	if err != nil {
		return fmt.Errorf("unable to get mcm program account info: %w", err)
	}
	err = binary.UnmarshalBorsh(&programData, data.Bytes())
	if err != nil {
		return fmt.Errorf("unable to unmarshal program data: %w", err)
	}

	instruction, err := mcmBindings.NewInitializeInstruction(
		chain.Selector,
		multisigID,
		GetMCMConfigPDA(mcmProgram, multisigID),
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
		mcmProgram,
		programData.Address,
		GetMCMRootMetadataPDA(mcmProgram, multisigID),
		GetMCMExpiringRootAndOpCountPDA(mcmProgram, multisigID),
	).ValidateAndBuild()
	if err != nil {
		return fmt.Errorf("unable to build instruction: %w", err)
	}

	err = chain.Confirm([]solana.Instruction{instruction})
	if err != nil {
		return fmt.Errorf("unable to confirm instructions: %w", err)
	}

	return nil
}

func initializeTimelock(
	e deployment.Environment, chain deployment.SolChain, timelockProgram solana.PublicKey, minDelay *big.Int,
) error {
	timelockID := [32]byte{} // FIXME: where should this come from?

	if minDelay == nil {
		minDelay = big.NewInt(0)
	}

	var timelockConfig timelockBindings.Config
	err := chain.GetAccountDataBorshInto(e.GetContext(), GetTimelockConfigPDA(timelockProgram, timelockID),
		&timelockConfig)
	if err == nil {
		e.Logger.Infow("Timelock already initialized, skipping initialization", "chain", chain.String())
		return nil
	}

	var programData struct {
		DataType uint32
		Address  solana.PublicKey
	}
	opts := &rpc.GetAccountInfoOpts{Commitment: rpc.CommitmentConfirmed}

	data, err := chain.Client.GetAccountInfoWithOpts(e.GetContext(), timelockProgram, opts)
	if err != nil {
		return fmt.Errorf("unable to get timelock program account info: %w", err)
	}
	err = binary.UnmarshalBorsh(&programData, data.Bytes())
	if err != nil {
		return fmt.Errorf("unable to unmarshal program data: %w", err)
	}

	// FIXME: retrieve from chain state
	var accessControllerProgram solana.PublicKey
	var proposerRoleAccessController solana.PublicKey
	var executorRoleAccessController solana.PublicKey
	var cancellerRoleAccessController solana.PublicKey
	var bypasserRoleAccessController solana.PublicKey

	instruction, err := timelockBindings.NewInitializeInstruction(
		timelockID,
		minDelay.Uint64(), // minDelay,
		GetTimelockConfigPDA(timelockProgram, timelockID),
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
		timelockProgram,
		programData.Address,
		accessControllerProgram,
		proposerRoleAccessController,
		executorRoleAccessController,
		cancellerRoleAccessController,
		bypasserRoleAccessController,
	).ValidateAndBuild()
	if err != nil {
		return fmt.Errorf("unable to build instruction: %w", err)
	}

	err = chain.Confirm([]solana.Instruction{instruction})
	if err != nil {
		return fmt.Errorf("unable to confirm instructions: %w", err)
	}

	return nil
}

func initializeAccessController(
	e deployment.Environment, chain deployment.SolChain, accessControllerProgram solana.PublicKey,
) error {
	// discriminator + owner + proposed owner + access_list (64 max addresses + length)
	dataSize := uint64(8 + 32 + 32 + ((32 * 64) + 8))
	rentExemption, err := chain.Client.GetMinimumBalanceForRentExemption(e.GetContext(), dataSize, rpc.CommitmentConfirmed)
	if err != nil {
		return fmt.Errorf("unable to get minimum balance for rent exemption: %w", err)
	}

	accessControllerAccount, err := solana.NewRandomPrivateKey() // FIXME: what should we do with the priv. key? store in the state???
	if err != nil {
		return fmt.Errorf("unable to generate new random private key: %w", err)
	}

	createAccountInstruction, err := system.NewCreateAccountInstruction(rentExemption, dataSize,
		accessControllerProgram, chain.DeployerKey.PublicKey(), accessControllerAccount.PublicKey()).ValidateAndBuild()
	if err != nil {
		return fmt.Errorf("unable to create CreateAccount instruction: %w", err)
	}

	initializeInstruction, err := accessControllerBindings.NewInitializeInstruction(
		accessControllerAccount.PublicKey(),
		chain.DeployerKey.PublicKey(),
	).ValidateAndBuild()
	if err != nil {
		return fmt.Errorf("unable to build instruction: %w", err)
	}

	instructions := []solana.Instruction{createAccountInstruction, initializeInstruction}
	err = chain.Confirm(instructions, solanaUtils.AddSigners(accessControllerAccount))
	if err != nil {
		return fmt.Errorf("unable to confirm instructions: %w", err)
	}

	return nil
}
