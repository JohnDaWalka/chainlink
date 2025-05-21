package operation

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/aptos-labs/aptos-go-sdk"
	link_token "github.com/smartcontractkit/chainlink-aptos/bindings/link-token"
	mcmsbind "github.com/smartcontractkit/chainlink-aptos/bindings/mcms"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	"github.com/smartcontractkit/mcms/types"
)

// AcceptTokenOwnershipOp ...
var AcceptTokenOwnershipOp = operations.NewOperation(
	"accept-token-ownership-op",
	Version1_0_0,
	"Accept token ownership",
	acceptTokenOwnership,
)

func acceptTokenOwnership(b operations.Bundle, deps AptosDeps, token aptos.AccountAddress) (types.Transaction, error) {
	// Bind CCIP Package
	tokenBind := link_token.Bind(token, deps.AptosChain.Client)
	linkBind := tokenBind.LinkToken()
	encoder := linkBind.Encoder()
	moduleInfo, function, _, args, err := encoder.AcceptOwnership()
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode ApplyDestChainConfigUpdates: %w", err)
	}
	additionalFields := aptosmcms.AdditionalFields{
		PackageName: moduleInfo.PackageName,
		ModuleName:  moduleInfo.ModuleName,
		Function:    function,
	}
	afBytes, err := json.Marshal(additionalFields)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to marshal additional fields: %w", err)
	}
	return types.Transaction{
		To:               token.StringLong(),
		Data:             aptosmcms.ArgsToData(args),
		AdditionalFields: afBytes,
	}, nil
}

type DeployTokenInput struct {
	MaxSupply *big.Int
	Name      string
	Symbol    string
	Decimals  byte
	Icon      string
	Project   string
}

type DeployTokenOutput struct {
	TokenObjAddress aptos.AccountAddress
	TokenAddress    aptos.AccountAddress
	MCMSOps         []types.Operation
}

// DeployTokenOp ...
var DeployTokenOp = operations.NewOperation(
	"deploy-token-op",
	Version1_0_0,
	"deploy token",
	deployTokenOp,
)

func deployTokenOp(b operations.Bundle, deps AptosDeps, in DeployTokenInput) (DeployTokenOutput, error) {
	var mcmsOps []types.Operation
	// TODO: this should be a sequence, not a single operation
	mcmsContract := mcmsbind.Bind(deps.OnChainState.MCMSAddress, deps.AptosChain.Client)

	// Deploy LINK token
	linkTokenSeed := "BNM_TOKEN_8"
	linkTokenObjectAddress, err := mcmsContract.MCMSRegistry().GetNewCodeObjectAddress(nil, []byte(linkTokenSeed))
	if err != nil {
		return DeployTokenOutput{}, fmt.Errorf("failed to GetNewCodeObjectAddress: %w", err)
	}

	linkTokenStateAddress := linkTokenObjectAddress.NamedObjectAddress([]byte("link::link_token::token_state"))
	linkTokenMetadataAddress := linkTokenStateAddress.NamedObjectAddress([]byte(in.Symbol))
	fmt.Printf("LINK Token Metadata address: %v\n", linkTokenMetadataAddress.StringLong())

	linkTokenPayload, err := link_token.Compile(linkTokenObjectAddress)
	if err != nil {
		return DeployTokenOutput{}, fmt.Errorf("failed to compile LINK token: %w", err)
	}
	ops, err := utils.CreateChunksAndStage(linkTokenPayload, mcmsContract, deps.AptosChain.Selector, linkTokenSeed, nil)
	if err != nil {
		return DeployTokenOutput{}, fmt.Errorf("failed to create chunks for token pool: %w", err)
	}
	mcmsOps = append(mcmsOps, ops...)

	// Deploy LINK MCMS Registrar
	mcmsRegistrarPayload, err := link_token.CompileMCMSRegistrar(linkTokenObjectAddress, deps.OnChainState.MCMSAddress, true)
	if err != nil {
		return DeployTokenOutput{}, fmt.Errorf("failed to compile LINK token: %w", err)
	}
	ops, err = utils.CreateChunksAndStage(mcmsRegistrarPayload, mcmsContract, deps.AptosChain.Selector, "", &linkTokenObjectAddress)
	if err != nil {
		return DeployTokenOutput{}, fmt.Errorf("failed to create chunks for token pool: %w", err)
	}
	mcmsOps = append(mcmsOps, ops...)

	// Initialize LINK token
	boundLinkToken := link_token.Bind(linkTokenObjectAddress, deps.AptosChain.Client)
	moduleInfo, function, _, args, err := boundLinkToken.LinkToken().Encoder().Initialize(
		&in.MaxSupply,
		in.Name,
		in.Symbol,
		in.Decimals,
		in.Icon,
		in.Project,
	)
	// Create MCMS operation
	additionalFields := aptosmcms.AdditionalFields{
		PackageName: moduleInfo.PackageName,
		ModuleName:  moduleInfo.ModuleName,
		Function:    function,
	}
	afBytes, err := json.Marshal(additionalFields)
	if err != nil {
		return DeployTokenOutput{}, fmt.Errorf("failed to marshal additional fields: %w", err)
	}

	mcmsOps = append(mcmsOps, types.Operation{
		ChainSelector: types.ChainSelector(deps.AptosChain.Selector),
		Transaction: types.Transaction{
			To:               linkTokenObjectAddress.StringLong(),
			Data:             aptosmcms.ArgsToData(args),
			AdditionalFields: afBytes,
		},
	})
	return DeployTokenOutput{
		TokenObjAddress: linkTokenObjectAddress,
		TokenAddress:    linkTokenMetadataAddress,
		MCMSOps:         mcmsOps,
	}, nil
}
