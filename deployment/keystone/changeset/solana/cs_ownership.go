package solana

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// ContractConfig defines the configuration for a contract ownership transfer
type ContractConfig struct {
	ContractType datastore.ContractType
	StateType    datastore.ContractType
	OperationID  string
	Description  string
}

// TransferOwnershipRequest represents a generic ownership transfer request
type TransferOwnershipRequest struct {
	ChainSel                    uint64
	CurrentOwner, ProposedOwner solana.PublicKey
	Version                     string
	Qualifier                   string
	MCMSCfg                     proposalutils.TimelockConfig
	ContractConfig              ContractConfig
}

// TransferOwnershipForwarderRequest wraps the generic request for forwarder contracts
type TransferOwnershipForwarderRequest struct {
	ChainSel                    uint64
	CurrentOwner, ProposedOwner solana.PublicKey
	Version                     string
	Qualifier                   string
	MCMSCfg                     proposalutils.TimelockConfig
}

// genericTransferOwnership handles the common ownership transfer logic
func GenericTransferOwnership(env cldf.Environment, req *TransferOwnershipRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	version := semver.MustParse(req.Version)
	
	// Build address references
	contractStateRef := datastore.NewAddressRefKey(req.ChainSel, req.ContractConfig.StateType, version, req.Qualifier)
	contractRef := datastore.NewAddressRefKey(req.ChainSel, req.ContractConfig.ContractType, version, req.Qualifier)

	// Get contract addresses
	contract, err := env.DataStore.Addresses().Get(contractRef)
	if err != nil {
		return out, fmt.Errorf("failed to get contract address: %w", err)
	}
	
	contractState, err := env.DataStore.Addresses().Get(contractStateRef)
	if err != nil {
		return out, fmt.Errorf("failed to get contract state address: %w", err)
	}

	// Load MCMS state
	mcmsState, err := state.MaybeLoadMCMSWithTimelockChainStateSolanaV2(
		env.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(req.ChainSel)))
	if err != nil {
		return out, err
	}

	solChain := env.BlockChains.SolanaChains()[req.ChainSel]

	// Execute the transfer operation
	execOut, err := operations.ExecuteOperation(env.OperationsBundle,
		operations.NewOperation(
			req.ContractConfig.OperationID,
			version,
			req.ContractConfig.Description,
			commonchangeset.TransferToTimelockSolanaOp,
		),
		commonchangeset.Deps{
			Env:   env,
			State: mcmsState,
			Chain: solChain,
		},
		commonchangeset.TransferToTimelockInput{
			Contract: commonchangeset.OwnableContract{
				Type:      cldf.ContractType(req.ContractConfig.ContractType),
				ProgramID: solana.MustPublicKeyFromBase58(contract.Address),
				OwnerPDA:  solana.MustPublicKeyFromBase58(contractState.Address),
			},
			MCMSCfg: req.MCMSCfg,
		},
	)
	if err != nil {
		return out, err
	}

	// Build proposal maps
	timelocks := map[uint64]string{}
	proposers := map[uint64]string{}
	inspectors := map[uint64]mcmssdk.Inspector{}

	inspectors[req.ChainSel] = mcmssolanasdk.NewInspector(solChain.Client)
	timelocks[req.ChainSel] = mcmssolanasdk.ContractAddress(mcmsState.TimelockProgram, mcmssolanasdk.PDASeed(mcmsState.TimelockSeed))
	proposers[req.ChainSel] = mcmssolanasdk.ContractAddress(mcmsState.McmProgram, mcmssolanasdk.PDASeed(mcmsState.ProposerMcmSeed))

	// Create timelock proposal
	proposal, err := proposalutils.BuildProposalFromBatchesV2(env, timelocks, proposers, inspectors,
		execOut.Output.Batches, fmt.Sprintf("proposal to transfer ownership of %s to timelock", req.ContractConfig.ContractType), req.MCMSCfg)
	if err != nil {
		return out, fmt.Errorf("failed to build proposal: %w", err)
	}
	env.Logger.Debugw("created timelock proposal", "# batches", len(execOut.Output.Batches))

	out.MCMSTimelockProposals = []mcms.TimelockProposal{*proposal}
	return out, nil
}

// genericVerifyPreconditions handles the common precondition verification logic
func GenericVerifyPreconditions(env cldf.Environment, chainSel uint64, version, qualifier string, contractType datastore.ContractType) error {
	// Validate version
	if _, err := semver.NewVersion(version); err != nil {
		return err
	}

	// Check if chain exists
	if _, ok := env.BlockChains.SolanaChains()[chainSel]; !ok {
		return fmt.Errorf("solana chain not found for chain selector %d", chainSel)
	}

	// Verify contract exists
	v := semver.MustParse(version)
	contractKey := datastore.NewAddressRefKey(chainSel, contractType, v, qualifier)
	if _, err := env.DataStore.Addresses().Get(contractKey); err != nil {
		return fmt.Errorf("failed to get %s for chain selector %d: %w", contractType, chainSel, err)
	}

	return nil
}

// TransferOwnershipForwarder implementation
var _ cldf.ChangeSetV2[*TransferOwnershipForwarderRequest] = TransferOwnershipForwarder{}

type TransferOwnershipForwarder struct{}

func (cs TransferOwnershipForwarder) VerifyPreconditions(env cldf.Environment, req *TransferOwnershipForwarderRequest) error {
	return GenericVerifyPreconditions(env, req.ChainSel, req.Version, req.Qualifier, "ForwarderContract")
}

func (cs TransferOwnershipForwarder) Apply(env cldf.Environment, req *TransferOwnershipForwarderRequest) (cldf.ChangesetOutput, error) {
	genericReq := &TransferOwnershipRequest{
		ChainSel:      req.ChainSel,
		CurrentOwner:  req.CurrentOwner,
		ProposedOwner: req.ProposedOwner,
		Version:       req.Version,
		Qualifier:     req.Qualifier,
		MCMSCfg:       req.MCMSCfg,
		ContractConfig: ContractConfig{
			ContractType: "ForwarderContract",
			StateType:    "ForwarderState",
			OperationID:  "transfer-ownership-forwarder",
			Description:  "transfers ownership of forwarder to mcms",
		},
	}
	return GenericTransferOwnership(env, genericReq)
}

