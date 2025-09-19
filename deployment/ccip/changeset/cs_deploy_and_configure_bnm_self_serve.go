package changeset

import (
	"errors"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/burn_mint_erc677_helper"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/registry_module_owner_custom"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

// This changeset deploys a BurnMintERC677Helper token and configures it with a BurnMintTokenPool on a specified chain.
// It also registers the token as an admin in the RegistryModuleOwnerCustom, accepts the admin role in the TokenAdminRegistry,
// the deployer key will remain the owner of the token intended for token transfer testing in testnet and SHOULD not be used for customer facing operations.
var DeployAndConfigureBnMSelfServe = cldf.CreateChangeSet(deployAndConfigureBnMSelfServeLogic, deployAndConfigureBnMSelfServeValidation)

type DeployAndConfigureBnMSelfServeConfig struct {
	Selector                  uint64 `json:"selector"`
	TokenName                 string `json:"tokenName"`
	TokenSymbol               string `json:"tokenSymbol"`
	RegistryModuleOwnerCustom string `json:"registryModuleOwnerCustom"`
	TokenPoolConfig           v1_5_1.ConfigureTokenPoolContractsConfig
}

func deployAndConfigureBnMSelfServeValidation(e cldf.Environment, cfg DeployAndConfigureBnMSelfServeConfig) error {
	if err := cldf.IsValidChainSelector(cfg.Selector); err != nil {
		return fmt.Errorf("invalid chain selector: %d - %w", cfg.Selector, err)
	}

	chainName := e.BlockChains.EVMChains()[cfg.Selector].Name()
	// This workflow is only intended for testnet BnM Helper token for testing purposes
	if e.Name == "mainnet" || strings.Contains(chainName, "mainnet") {
		return errors.New("minting on LINK token is not allowed on Mainnet")
	}

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	if err := stateview.ValidateChain(e, state, cfg.Selector, nil); err != nil {
		return fmt.Errorf("failed to validate chain %d: %w", cfg.Selector, err)
	}

	chainState, ok := state.EVMChainState(cfg.Selector)
	if !ok {
		return fmt.Errorf("%d does not exist in state", cfg.Selector)
	}

	if registryModuleOwnerCustom := chainState.RegistryModules1_6; registryModuleOwnerCustom == nil {
		return fmt.Errorf("missing registry_module_owner_custom on %d", cfg.Selector)
	}

	if tokenAdminReg := chainState.TokenAdminRegistry; tokenAdminReg == nil {
		return fmt.Errorf("missing tokenAdminReg on %d", cfg.Selector)
	}

	if rmnProxy := chainState.RMNProxy; rmnProxy == nil {
		return fmt.Errorf("missing rmnProxy on %d", cfg.Selector)
	}

	if router := chainState.Router; router == nil {
		return fmt.Errorf("missing router on %d", cfg.Selector)
	}

	return nil
}

func deployAndConfigureBnMSelfServeLogic(e cldf.Environment, cfg DeployAndConfigureBnMSelfServeConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()
	finalCSOut := &cldf.ChangesetOutput{
		AddressBook: ab,
	}

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load existing onchain state: %w", err)
	}
	chain := e.BlockChains.EVMChains()[cfg.Selector]
	chainState, chainExists := state.EVMChainState(cfg.Selector)
	if !chainExists {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain with selector %d does not exist in state", cfg.Selector)
	}

	// Step 1: Deploy BnM Helper token With Drip
	bnmToken, _, err := deployTokenBnMHelperToken(e, state, ab, cfg)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy BnM Helper token: %w", err)
	}
	if err := e.ExistingAddresses.Merge(ab); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge address book for token %s: %w", cfg.TokenName, err)
	}

	// Step 2: Deploy BurnMintToken Pool for BnM token
	deployTPCfg := v1_5_1.DeployTokenPoolContractsConfig{
		TokenSymbol:  shared.CCIPBnMSymbol,
		IsTestRouter: false,
		NewPools: map[uint64]v1_5_1.DeployTokenPoolInput{
			cfg.Selector: {
				Type:               shared.BurnMintTokenPool,
				TokenAddress:       bnmToken.Address,
				LocalTokenDecimals: 18,
				AcceptLiquidity:    nil,
			},
		},
	}

	tokenPoolOutput, err := v1_5_1.DeployTokenPoolContractsChangeset(e, deployTPCfg)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy token pool for token %s: %w", cfg.TokenName, err)
	}
	if err := cldf.MergeChangesetOutput(e, finalCSOut, tokenPoolOutput); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge address book for token %s: %w", cfg.TokenName, err)
	}

	// Step 3: Configure the token pool for the BnM token
	output, err := v1_5_1.ConfigureTokenPoolContractsChangeset(e, cfg.TokenPoolConfig)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to configure token pool for token %s: %w", cfg.TokenName, err)
	}
	if err := cldf.MergeChangesetOutput(e, finalCSOut, output); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge changeset output after configuring token pool for token %s: %w", cfg.TokenName, err)
	}

	state, err = stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load existing onchain state: %w", err)
	}

	chainState, _ = state.EVMChainState(cfg.Selector)

	tokenAdminReg := chainState.TokenAdminRegistry
	var registryModuleOwnerCustom *registry_module_owner_custom.RegistryModuleOwnerCustom
	burnMintTokenPool, exists := v1_5_1.GetTokenPoolAddressFromSymbolTypeAndVersion(chainState, chain, shared.CCIPBnMSymbol, shared.BurnMintTokenPool, deployment.Version1_5_1)
	if !exists {
		return cldf.ChangesetOutput{}, fmt.Errorf("could not find token pool for BnM token with symbol %s on chain with selector %d", shared.CCIPBnMSymbol, cfg.Selector)
	}
	for _, registry := range chainState.RegistryModules1_6 {
		if cfg.RegistryModuleOwnerCustom != "" && registry.Address().Hex() == cfg.RegistryModuleOwnerCustom {
			registryModuleOwnerCustom = registry
			break
		}
	}
	if registryModuleOwnerCustom == nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("could not find registry_module_owner_custom with address %s on chain with selector %d", cfg.RegistryModuleOwnerCustom, cfg.Selector)
	}
	if tokenAdminReg == nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("could not find tokenAdminReg on chain with selector %d", cfg.Selector)
	}

	// Step 4: Register BnM Token as admin via Deployer key on RegistryModuleOwnerCustom
	tx, err := registryModuleOwnerCustom.RegisterAdminViaOwner(chain.DeployerKey, bnmToken.Address)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to register BnM token %s as admin: %w", cfg.TokenName, err)
	}
	_, err = chain.Confirm(tx)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm RegisterAdminViaOwner on RegistryModuleOwnerCustom: %w", err)
	}

	// Step 5: Accept Admin role on TokenAdminRegistry
	tx, err = tokenAdminReg.AcceptAdminRole(chain.DeployerKey, bnmToken.Address)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to Accept Admin role for BnM token %s as admin: %w", cfg.TokenName, err)
	}
	_, err = chain.Confirm(tx)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm AcceptAdminRole on TokenAdminRegistry: %w", err)
	}

	// Step 6: Set the BnM token pool on TokenAdminRegistry
	tx, err = tokenAdminReg.SetPool(chain.DeployerKey, bnmToken.Address, burnMintTokenPool)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to set BnM token %s pool on TokenAdminRegistry: %w", cfg.TokenName, err)
	}
	_, err = chain.Confirm(tx)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm SetPool on TokenAdminRegistry: %w", err)
	}

	// Step 7: Grant Mint role to the token pool on the BnM token
	tx, err = bnmToken.Contract.GrantMintRole(chain.DeployerKey, burnMintTokenPool)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to grant mint role to BnM token %s: %w", cfg.TokenName, err)
	}
	_, err = chain.Confirm(tx)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm grant mint role to BnM token %s: %w", cfg.TokenName, err)
	}

	return *finalCSOut, nil
}

func deployTokenBnMHelperToken(e cldf.Environment, state stateview.CCIPOnChainState, ab cldf.AddressBook, cfg DeployAndConfigureBnMSelfServeConfig) (cldf.ContractDeploy[*burn_mint_erc677_helper.BurnMintERC677Helper],
	cldf.ContractDeploy[*burn_mint_token_pool.BurnMintTokenPool],
	error) {
	chain := e.BlockChains.EVMChains()[cfg.Selector]

	token, err := cldf.DeployContract(e.Logger, chain, ab,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*burn_mint_erc677_helper.BurnMintERC677Helper] {
			tokenAddress, tx, token, err := burn_mint_erc677_helper.DeployBurnMintERC677Helper(
				chain.DeployerKey,
				chain.Client,
				cfg.TokenName,
				cfg.TokenSymbol,
			)

			return cldf.ContractDeploy[*burn_mint_erc677_helper.BurnMintERC677Helper]{
				Address:  tokenAddress,
				Contract: token,
				Tv:       cldf.NewTypeAndVersion(shared.BurnMintToken, deployment.Version1_0_0),
				Tx:       tx,
				Err:      err,
			}
		},
	)
	if err != nil {
		return cldf.ContractDeploy[*burn_mint_erc677_helper.BurnMintERC677Helper]{},
			cldf.ContractDeploy[*burn_mint_token_pool.BurnMintTokenPool]{},
			fmt.Errorf("failed to deploy BnM token: %w", err)
	}

	return *token, cldf.ContractDeploy[*burn_mint_token_pool.BurnMintTokenPool]{}, nil
}
