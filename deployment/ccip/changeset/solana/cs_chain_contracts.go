package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
)

var _ deployment.ChangeSet[v1_6.SetOCR3OffRampConfig] = SetOCR3ConfigSolana
var _ deployment.ChangeSet[AddRemoteChainToSolanaConfig] = AddRemoteChainToSolana
var _ deployment.ChangeSet[BillingTokenConfig] = AddBillingTokenChangeset
var _ deployment.ChangeSet[BillingTokenForRemoteChainConfig] = AddBillingTokenForRemoteChain
var _ deployment.ChangeSet[RegisterTokenAdminRegistryConfig] = RegisterTokenAdminRegistry
var _ deployment.ChangeSet[TransferAdminRoleTokenAdminRegistryConfig] = TransferAdminRoleTokenAdminRegistry
var _ deployment.ChangeSet[AcceptAdminRoleTokenAdminRegistryConfig] = AcceptAdminRoleTokenAdminRegistry
var _ deployment.ChangeSet[SetFeeAggregatorConfig] = SetFeeAggregator

// HELPER FUNCTIONS
// GetTokenProgramID returns the program ID for the given token program name
func GetTokenProgramID(programName deployment.ContractType) (solana.PublicKey, error) {
	tokenPrograms := map[deployment.ContractType]solana.PublicKey{
		ccipChangeset.SPLTokens:     solana.TokenProgramID,
		ccipChangeset.SPL2022Tokens: solana.Token2022ProgramID,
	}

	programID, ok := tokenPrograms[programName]
	if !ok {
		return solana.PublicKey{}, fmt.Errorf("invalid token program: %s. Must be one of: %s, %s", programName, ccipChangeset.SPLTokens, ccipChangeset.SPL2022Tokens)
	}
	return programID, nil
}

func commonValidation(e deployment.Environment, selector uint64, tokenPubKey solana.PublicKey) error {
	chain, ok := e.SolChains[selector]
	if !ok {
		return fmt.Errorf("chain selector %d not found in environment", selector)
	}
	state, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState, chainExists := state.SolChains[selector]
	if !chainExists {
		return fmt.Errorf("chain %s not found in existing state, deploy the link token first", chain.String())
	}
	if tokenPubKey.Equals(chainState.LinkToken) || tokenPubKey.Equals(chainState.WSOL) {
		return nil
	}
	exists := false
	allTokens := chainState.SPL2022Tokens
	allTokens = append(allTokens, chainState.SPLTokens...)
	for _, token := range allTokens {
		if token.Equals(tokenPubKey) {
			exists = true
			break
		}
	}
	if !exists {
		return fmt.Errorf("token %s not found in existing state, deploy the token first", tokenPubKey.String())
	}
	return nil
}

// TEST ROUTER ? -> can easily check for test router here
func validateRouterConfig(chain deployment.SolChain, chainState ccipChangeset.SolCCIPChainState) error {
	if chainState.Router.IsZero() {
		return fmt.Errorf("router not found in existing state, deploy the router first for chain %d", chain.Selector)
	}
	// addressing errcheck in the next PR
	var routerConfigAccount solRouter.Config
	err := chain.GetAccountDataBorshInto(context.Background(), chainState.RouterConfigPDA, &routerConfigAccount)
	if err != nil {
		return fmt.Errorf("router config not found in existing state, initialize the router first %d", chain.Selector)
	}
	return nil
}

func validateFeeQuoterConfig(chain deployment.SolChain, chainState ccipChangeset.SolCCIPChainState) error {
	if chainState.FeeQuoter.IsZero() {
		return fmt.Errorf("fee quoter not found in existing state, deploy the fee quoter first for chain %d", chain.Selector)
	}
	var fqConfig solFeeQuoter.Config
	feeQuoterConfigPDA, _, _ := solState.FindFqConfigPDA(chainState.FeeQuoter)
	err := chain.GetAccountDataBorshInto(context.Background(), feeQuoterConfigPDA, &fqConfig)
	if err != nil {
		return fmt.Errorf("fee quoter config not found in existing state, initialize the fee quoter first %d", chain.Selector)
	}
	return nil
}

func validateOffRampConfig(chain deployment.SolChain, chainState ccipChangeset.SolCCIPChainState) error {
	if chainState.OffRamp.IsZero() {
		return fmt.Errorf("offramp not found in existing state, deploy the offramp first for chain %d", chain.Selector)
	}
	var offRampConfig solOffRamp.Config
	offRampConfigPDA, _, _ := solState.FindOfframpConfigPDA(chainState.OffRamp)
	err := chain.GetAccountDataBorshInto(context.Background(), offRampConfigPDA, &offRampConfig)
	if err != nil {
		return fmt.Errorf("offramp config not found in existing state, initialize the offramp first %d", chain.Selector)
	}
	return nil
}

type SetFeeAggregatorConfig struct {
	ChainSelector uint64
	FeeAggregator string
}

func (cfg SetFeeAggregatorConfig) Validate(e deployment.Environment) error {
	state, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState, chainExists := state.SolChains[cfg.ChainSelector]
	if !chainExists {
		return fmt.Errorf("chain %d not found in existing state", cfg.ChainSelector)
	}
	chain := e.SolChains[cfg.ChainSelector]

	if err := validateRouterConfig(chain, chainState); err != nil {
		return err
	}

	// Validate fee aggregator address is valid
	if _, err := solana.PublicKeyFromBase58(cfg.FeeAggregator); err != nil {
		return fmt.Errorf("invalid fee aggregator address: %w", err)
	}

	if chainState.FeeAggregator.Equals(solana.MustPublicKeyFromBase58(cfg.FeeAggregator)) {
		return fmt.Errorf("fee aggregator %s is already set on chain %d", cfg.FeeAggregator, cfg.ChainSelector)
	}

	return nil
}

func SetFeeAggregator(e deployment.Environment, cfg SetFeeAggregatorConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]

	feeAggregatorPubKey := solana.MustPublicKeyFromBase58(cfg.FeeAggregator)
	routerConfigPDA, _, _ := solState.FindConfigPDA(chainState.Router)

	solRouter.SetProgramID(chainState.Router)
	instruction, err := solRouter.NewUpdateFeeAggregatorInstruction(
		feeAggregatorPubKey,
		routerConfigPDA,
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if err := chain.Confirm([]solana.Instruction{instruction}); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	newAddresses := deployment.NewMemoryAddressBook()
	err = newAddresses.Save(cfg.ChainSelector, cfg.FeeAggregator, deployment.NewTypeAndVersion(ccipChangeset.FeeAggregator, deployment.Version1_0_0))
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to save address: %w", err)
	}

	e.Logger.Infow("Set new fee aggregator", "chain", chain.String(), "fee_aggregator", feeAggregatorPubKey.String())
	return deployment.ChangesetOutput{
		AddressBook: newAddresses,
	}, nil
}
