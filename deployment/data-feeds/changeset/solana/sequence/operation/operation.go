package operation

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	solanaUtils "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	cldfsol "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	ks_cache "github.com/smartcontractkit/chainlink-solana/contracts/generated/data_feeds_cache"

	commonOps "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/operations"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
)

var Version1_0_0 = semver.MustParse("1.0.0")

var (
	InitCacheOp = operations.NewOperation(
		"init-cache-op",
		Version1_0_0,
		"Initialize DataFeeds Cache for Solana Chain",
		initCache,
	)
	DeployCacheOp = operations.NewOperation(
		"deploy-cache-op",
		Version1_0_0,
		"Deploys the DataFeeds Cache program for Solana Chain",
		commonOps.Deploy,
	)
	SetUpgradeAuthorityOp = operations.NewOperation(
		"set-upgrade-authority-op",
		Version1_0_0,
		"Sets Cache's upgrade authority for Solana Chain",
		setUpgradeAuthority,
	)
	ConfigureCacheDecimalReportOp = operations.NewOperation(
		"configure-cache-decimal-report-op",
		Version1_0_0,
		"Configure cache decimal report for Solana Chain",
		configureCacheDecimalReport,
	)
	InitCacheDecimalReportOp = operations.NewOperation(
		"init-cache-decimal-feed-op",
		Version1_0_0,
		"Initialize DataFeeds Cache Decimal Report for Solana Chain",
		initCacheDecimalReport,
	)
)

type (
	Deps struct {
		Env       cldf.Environment
		Chain     cldfsol.Chain
		Datastore datastore.DataStore
	}

	// For DataFeeds Cache initialization
	InitCacheInput struct {
		ProgramID solana.PublicKey
		ChainSel  uint64
	}

	InitCacheOutput struct {
		StatePubKey solana.PublicKey
	}

	SetUpgradeAuthorityInput struct {
		ChainSel            uint64
		ProgramID           string
		NewUpgradeAuthority string
		MCMS                *proposalutils.TimelockConfig // if set, assumes current upgrade authority is the timelock
	}

	SetUpgradeAuthorityOutput struct {
		Proposals []mcms.TimelockProposal // will be returned in case if timelock config is passed
	}

	ConfigureCacheDecimalReportInput struct {
		ChainSel             uint64
		Descriptions         [][32]uint8
		DataIDs              [][16]uint8
		MCMS                 *proposalutils.TimelockConfig // if set, assumes current owner is the timelock
		ProgramID            solana.PublicKey
		AllowedSender        []solana.PublicKey
		AllowedWorkflowOwner [][20]uint8
		AllowedWorkflowName  [][10]uint8
		FeedAdmin            solana.PublicKey
		State                solana.PublicKey
		Type                 cldf.ContractType
	}

	ConfigureCacheOutput struct {
		Proposals []mcms.TimelockProposal // will be returned in case if timelock config is passed
	}

	InitCacheDecimalReportInput struct {
		ChainSel  uint64
		Version   string
		Qualifier string
		MCMS      *proposalutils.TimelockConfig // if set, assumes current
		DataIDs   [][16]uint8
		FeedAdmin solana.PublicKey
		State     solana.PublicKey
		ProgramID solana.PublicKey
		Type      cldf.ContractType
	}
)

func initCache(b operations.Bundle, deps Deps, in InitCacheInput) (InitCacheOutput, error) {
	var out InitCacheOutput
	if ks_cache.ProgramID.IsZero() {
		ks_cache.SetProgramID(in.ProgramID)
	}

	stateKey, err := solana.NewRandomPrivateKey()
	if err != nil {
		return out, fmt.Errorf("failed to create random keys: %w", err)
	}

	adminStateKey, err := solana.NewRandomPrivateKey()
	if err != nil {
		return out, fmt.Errorf("failed to create random admin keys: %w", err)
	}

	instruction, err := ks_cache.NewInitializeInstruction([]solana.PublicKey{adminStateKey.PublicKey()}, deps.Chain.DeployerKey.PublicKey(), stateKey.PublicKey(), solana.SystemProgramID).ValidateAndBuild()
	if err != nil {
		return out, fmt.Errorf("failed to build and validate initialize instruction %w", err)
	}

	instructions := []solana.Instruction{instruction}
	if err = deps.Chain.Confirm(instructions, solanaUtils.AddSigners(stateKey)); err != nil {
		return out, errors.New("failed to confirm ")
	}

	out.StatePubKey = stateKey.PublicKey()

	return out, nil
}

func setUpgradeAuthority(b operations.Bundle, deps Deps, in SetUpgradeAuthorityInput) (SetUpgradeAuthorityOutput, error) {
	var out SetUpgradeAuthorityOutput

	programID, err := solana.PublicKeyFromBase58(in.ProgramID)
	if err != nil {
		return out, fmt.Errorf("failed parse programID: %w", err)
	}

	newAuthority, err := solana.PublicKeyFromBase58(in.NewUpgradeAuthority)
	if err != nil {
		return out, fmt.Errorf("failed parse upgrade authority: %w", err)
	}

	currentAuthority := deps.Chain.DeployerKey.PublicKey()
	if in.MCMS != nil {
		timelockSignerPDA, err := helpers.FetchTimelockSigner(deps.Datastore.Addresses().Filter(datastore.AddressRefByChainSelector(in.ChainSel)))
		if err != nil {
			return out, fmt.Errorf("failed to get timelock signer: %w", err)
		}
		currentAuthority = timelockSignerPDA
	}

	mcmsTxns := make([]mcmsTypes.Transaction, 0)

	ixn := helpers.SetUpgradeAuthority(&deps.Env, programID, currentAuthority, newAuthority, false)

	if in.MCMS == nil {
		if err := deps.Chain.Confirm([]solana.Instruction{ixn}); err != nil {
			return out, fmt.Errorf("failed to confirm instructions: %w", err)
		}

		return out, nil
	}

	// build MCMS proposal
	tx, err := helpers.BuildMCMSTxn(
		ixn,
		solana.BPFLoaderUpgradeableProgramID.String(),
		cldf.ContractType(solana.BPFLoaderUpgradeableProgramID.String()))
	if err != nil {
		return out, fmt.Errorf("failed to create transaction: %w", err)
	}
	mcmsTxns = append(mcmsTxns, *tx)

	proposal, err := helpers.BuildProposalsForTxns(
		deps.Env, in.ChainSel, "proposal to SetUpgradeAuthority in Solana", in.MCMS.MinDelay, mcmsTxns)
	if err != nil {
		return out, fmt.Errorf("failed to build proposal: %w", err)
	}
	out.Proposals = []mcms.TimelockProposal{*proposal}

	return out, nil
}

func initCacheDecimalReport(b operations.Bundle, deps Deps, in InitCacheDecimalReportInput) (ConfigureCacheOutput, error) {
	var out ConfigureCacheOutput
	var ixn *ks_cache.Instruction
	mcmsTxns := make([]mcmsTypes.Transaction, 0)
	if ks_cache.ProgramID.IsZero() {
		ks_cache.SetProgramID(in.ProgramID)
	}

	ixn, err := ks_cache.NewInitDecimalReportsInstruction(in.DataIDs, in.FeedAdmin, in.State, in.ProgramID).ValidateAndBuild()
	if err != nil {
		return out, fmt.Errorf("failed to build and validate initialize instruction %w", err)
	}

	if in.MCMS == nil {
		if err := deps.Chain.Confirm([]solana.Instruction{ixn}); err != nil {
			return out, fmt.Errorf("failed to confirm instructions: %w", err)
		}

		return out, nil
	}

	tx, err := helpers.BuildMCMSTxn(
		ixn,
		solana.BPFLoaderUpgradeableProgramID.String(),
		cldf.ContractType(solana.BPFLoaderUpgradeableProgramID.String()))
	if err != nil {
		return out, fmt.Errorf("failed to create transaction: %w", err)
	}

	mcmsTxns = append(mcmsTxns, *tx)

	proposal, err := helpers.BuildProposalsForTxns(
		deps.Env, in.ChainSel, "proposal to SetUpgradeAuthority in Solana", in.MCMS.MinDelay, mcmsTxns)
	if err != nil {
		return out, fmt.Errorf("failed to build proposal: %w", err)
	}
	out.Proposals = []mcms.TimelockProposal{*proposal}

	return out, nil
}

func configureCacheDecimalReport(b operations.Bundle, deps Deps, in ConfigureCacheDecimalReportInput) (ConfigureCacheOutput, error) {
	var out ConfigureCacheOutput
	var ixn *ks_cache.Instruction
	if ks_cache.ProgramID.IsZero() {
		ks_cache.SetProgramID(in.ProgramID)
	}

	workflowMetas := make([]ks_cache.WorkflowMetadata, len(in.AllowedSender))
	for i := range in.AllowedSender {
		workflowMetas[i] = ks_cache.WorkflowMetadata{
			AllowedSender:        in.AllowedSender[i],
			AllowedWorkflowOwner: in.AllowedWorkflowOwner[i],
			AllowedWorkflowName:  in.AllowedWorkflowName[i],
		}
	}

	mcmsTxns := make([]mcmsTypes.Transaction, 0)

	ixn, err := ks_cache.NewSetDecimalFeedConfigsInstruction(
		in.DataIDs,
		in.Descriptions,
		workflowMetas,
		in.FeedAdmin,
		in.State,
		in.ProgramID,
	).ValidateAndBuild()
	if err != nil {
		return out, fmt.Errorf("cant build init oracle instruction: %w", err)
	}

	if in.MCMS == nil {
		if err := deps.Chain.Confirm([]solana.Instruction{ixn}); err != nil {
			return out, fmt.Errorf("failed to confirm instructions: %w", err)
		}

		return out, nil
	}

	tx, err := helpers.BuildMCMSTxn(
		ixn,
		solana.BPFLoaderUpgradeableProgramID.String(),
		cldf.ContractType(solana.BPFLoaderUpgradeableProgramID.String()))
	if err != nil {
		return out, fmt.Errorf("failed to create transaction: %w", err)
	}

	mcmsTxns = append(mcmsTxns, *tx)

	proposal, err := helpers.BuildProposalsForTxns(
		deps.Env, in.ChainSel, "proposal to SetUpgradeAuthority in Solana", in.MCMS.MinDelay, mcmsTxns)
	if err != nil {
		return out, fmt.Errorf("failed to build proposal: %w", err)
	}
	out.Proposals = []mcms.TimelockProposal{*proposal}

	return out, nil
}
