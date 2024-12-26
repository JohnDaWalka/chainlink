package changeset

import (
	"context"
	"fmt"

	bin "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
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
	chainState, _ := state.SolChains[chain.Selector]
	ccipRouterProgram := chainState.CcipRouter
	ccipReceiverProgram := chainState.CcipReceiver
	tokenPoolProgram := chainState.TokenPool
	if ccipRouterProgram.IsZero() {
		panic("ccipRouter is not set")
	} else {
		// Fetch account info for the program ID
		account, _ := chain.Client.GetAccountInfoWithOpts(context.Background(), ccipRouterProgram, &rpc.GetAccountInfoOpts{
			Commitment: rpc.CommitmentConfirmed,
		})
		programFile := "ccipRouter.so"           //TODO
		keypairPath := "ccipRouter-keypair.json" // TODO
		programKeyPair := "ccipRouter-keypair"   // TODO
		// Deploy program if it doesn't exist
		if account != nil && account.Value.Executable {
			lggr.Info("ccipRouter exists and is executable.")
		} else {
			lggr.Info("Program does not exist or is not executable.")
			// Deploy the program
			programID, err := deployment.DeploySolProgramCLI(programFile, keypairPath, programKeyPair)
			if err != nil {
				lggr.Fatalf("Failed to deploy program: %v", err)
			}
			// Verify the program ID (simple check for non-empty string)
			if programID == "" {
				lggr.Fatalf("Program ID is empty")
			}

			lggr.Infof("programID %s", programID)
		}

		// program should exist by now (either already deployed, or deployed and waited for confirmation)
		ccip_router.SetProgramID(ccipRouterProgram)

		// wallet keys
		privateKey, _ := ag_solanago.PrivateKeyFromSolanaKeygenFile(keypairPath)
		publicKey := privateKey.PublicKey()

		// this is a PDA that gets initialised when you call init on the programID
		RouterConfigPDA, _, _ := ag_solanago.FindProgramAddress([][]byte{[]byte("config")}, ccipRouterProgram)
		RouterStatePDA, _, _ := ag_solanago.FindProgramAddress([][]byte{[]byte("state")}, ccipRouterProgram)
		ExternalExecutionConfigPDA, _, _ := ag_solanago.FindProgramAddress([][]byte{[]byte("external_execution_config")}, ccipRouterProgram)
		ExternalTokenPoolsSignerPDA, _, _ := ag_solanago.FindProgramAddress([][]byte{[]byte("external_token_pools_signer")}, ccipRouterProgram)

		// check if the PDA is already initialised
		data, err := chain.Client.GetAccountInfoWithOpts(context.Background(), ccipRouterProgram, &rpc.GetAccountInfoOpts{
			Commitment: rpc.CommitmentConfirmed,
		})
		if err != nil {
			lggr.Fatalf("Failed to get account info: %v", err)
		}

		var programData struct {
			DataType uint32
			Address  ag_solanago.PublicKey
		}
		err = bin.UnmarshalBorsh(&programData, data.Bytes())
		if err != nil {
			lggr.Fatalf("Failed to unmarshal data: %v", err)
		}
		instruction, err := ccip_router.NewInitializeInstruction(
			chain.Selector, // chain selector
			bin.Uint128{},  // default gas limit
			true,           // allow out of order execution
			int64(1800),    // period to wait before allowing manual execution. 30 mins
			RouterConfigPDA,
			RouterStatePDA,
			publicKey,
			ag_solanago.SystemProgramID,
			ccipRouterProgram,
			programData.Address,
			ExternalExecutionConfigPDA,
			ExternalTokenPoolsSignerPDA,
		).ValidateAndBuild()
		_, err = common.SendAndConfirm(context.Background(), chain.Client, []ag_solanago.Instruction{instruction}, privateKey, rpc.CommitmentConfirmed)
		if err != nil {
			lggr.Fatalf("Failed to send and confirm: %v", err)
		}
	}
	if ccipReceiverProgram.IsZero() {
		panic("ccipReceiver is not set")
	} else {
		// ReceiverTargetAccountPDA, _, _ = ag_solanago.FindProgramAddress([][]byte{[]byte("counter")}, ccipReceiver)
		// ReceiverExternalExecutionConfigPDA, _, _ = ag_solanago.FindProgramAddress([][]byte{[]byte("external_execution_config")}, CcipReceiverProgram)
		// TODO
	}
	if tokenPoolProgram.IsZero() {

	} else {
		// TODO
	}
	return nil
}
