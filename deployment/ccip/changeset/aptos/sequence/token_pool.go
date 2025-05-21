package sequence

import (
	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

// Deploy Token Pool sequence input
type DeployTokenPoolSeqInput struct {
	TokenParams      config.TokenParams
	TokenAddress     aptos.AccountAddress
	TokenObjAddress  aptos.AccountAddress
	TokenPoolAddress aptos.AccountAddress
	TokenAdmin       aptos.AccountAddress
	PoolType         deployment.ContractType
}

type DeployTokenPoolSeqOutput struct {
	TokenPoolAddress     aptos.AccountAddress
	TokenAddress         aptos.AccountAddress
	TokenObjAddress      aptos.AccountAddress
	CCIPTokenPoolAddress aptos.AccountAddress
	MCMSOperations       []mcmstypes.BatchOperation
}

var DeployAptosTokenPoolSequence = operations.NewSequence(
	"deploy-aptos-token-pool",
	operation.Version1_0_0,
	"Deploys token and token pool and configures",
	deployAptosTokenPoolSequence,
)

func deployAptosTokenPoolSequence(b operations.Bundle, deps operation.AptosDeps, in DeployTokenPoolSeqInput) (DeployTokenPoolSeqOutput, error) {
	mcmsOperations := []mcmstypes.BatchOperation{}
	// 0 - Cleanup staging area
	cleanupInput := operation.CleanupStagingAreaInput{
		MCMSAddress: deps.OnChainState.MCMSAddress,
	}
	cleanupReport, err := operations.ExecuteOperation(b, operation.CleanupStagingAreaOp, deps, cleanupInput)
	if err != nil {
		return DeployTokenPoolSeqOutput{}, err
	}
	if len(cleanupReport.Output.Transactions) > 0 {
		mcmsOperations = append(mcmsOperations, cleanupReport.Output)
	}

	tokenAddress := in.TokenAddress
	tokenObjAddres := in.TokenObjAddress
	// 1 - Deploy token (if not deployed)
	if in.TokenAddress == (aptos.AccountAddress{}) {
		deployTInput := operation.DeployTokenInput{
			MaxSupply: in.TokenParams.MaxSupply,
			Name:      in.TokenParams.Name,
			Symbol:    in.TokenParams.Symbol,
			Decimals:  in.TokenParams.Decimals,
			Icon:      in.TokenParams.Icon,
			Project:   in.TokenParams.Project,
		}
		deployTReport, err := operations.ExecuteOperation(b, operation.DeployTokenOp, deps, deployTInput)
		if err != nil {
			return DeployTokenPoolSeqOutput{}, err
		}
		mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployTReport.Output.MCMSOps)...)
		tokenAddress = deployTReport.Output.TokenAddress
		tokenObjAddres = deployTReport.Output.TokenObjAddress
	}

	tokenPoolAddress := in.TokenPoolAddress
	// 2 - Deploy token pool (if not deployed)
	if in.TokenPoolAddress == (aptos.AccountAddress{}) {
		deployInput := operation.DeployTokenPoolInput{
			PoolType:        in.PoolType,
			TokenAddress:    tokenAddress,
			TokenObjAddress: tokenObjAddres,
		}
		deployReport, err := operations.ExecuteOperation(b, operation.DeployTokenPoolOp, deps, deployInput)
		if err != nil {
			return DeployTokenPoolSeqOutput{}, err
		}
		mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployReport.Output.MCMSOps)...)
		tokenPoolAddress = deployReport.Output.TokenPoolAddress
	}

	txs := []mcmstypes.Transaction{}

	// signer of transfer_admin_role is MCMS, but initial owner is the token pool
	// // 3 - Transfer admin role
	// transferInput := operation.TransferAdminRoleInput{
	// 	Token:    tokenAddress,
	// 	NewAdmin: tokenObjAddres,
	// }
	// transferReport, err := operations.ExecuteOperation(b, operation.TransferAdminRoleOp, deps, transferInput)
	// if err != nil {
	// 	return DeployTokenPoolSeqOutput{}, err
	// }
	// txs = append(txs, transferReport.Output)

	// // 4 - Accept admin role
	// acceptReport, err := operations.ExecuteOperation(b, operation.AcceptAdminRoleOp, deps, tokenAddress)
	// if err != nil {
	// 	return DeployTokenPoolSeqOutput{}, err
	// }
	// txs = append(txs, acceptReport.Output)

	// 5 - Set Pool (token admin registry)
	setPoolInput := operation.SetPoolInput{
		TokenAddress: tokenAddress,
		PoolAddress:  tokenPoolAddress,
	}

	setPoolReport, err := operations.ExecuteOperation(b, operation.SetPoolOp, deps, setPoolInput)
	if err != nil {
		return DeployTokenPoolSeqOutput{}, err
	}
	txs = append(txs, setPoolReport.Output)

	mcmsOperations = append(mcmsOperations, mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(deps.AptosChain.Selector),
		Transactions:  txs,
	})
	return DeployTokenPoolSeqOutput{
		TokenPoolAddress:     tokenPoolAddress,
		TokenAddress:         tokenAddress,
		TokenObjAddress:      tokenObjAddres,
		CCIPTokenPoolAddress: tokenPoolAddress,
		MCMSOperations:       mcmsOperations,
	}, nil
}

// Connect Token Pool sequence input
type ConnectTokenPoolSeqInput struct {
	TokenPoolAddress aptos.AccountAddress
	PoolType         deployment.ContractType // TODO: should be a typed const
	RemotePools      map[uint64]operation.RemotePool
}

var ConnectTokenPoolSequence = operations.NewSequence(
	"connect-aptos-evm-token-pools",
	operation.Version1_0_0,
	"Connects EVM<>Aptos lanes token pools",
	connectTokenPoolSequence,
)

func connectTokenPoolSequence(b operations.Bundle, deps operation.AptosDeps, in ConnectTokenPoolSeqInput) (mcmstypes.BatchOperation, error) {
	var txs []mcmstypes.Transaction

	// 10 - SetupTokenPoolForRemoteChain (apply_chain_updates/add_remote_pool)
	setRemoteInput := operation.SetupTokenPoolInput{
		TokenPoolAddress: in.TokenPoolAddress,
		PoolType:         in.PoolType,
		RemotePools:      in.RemotePools,
	}
	setRemoteReport, err := operations.ExecuteOperation(b, operation.SetupTokenPoolOp, deps, setRemoteInput)
	if err != nil {
		return mcmstypes.BatchOperation{}, err
	}
	txs = append(txs, setRemoteReport.Output...)

	// 11 - AddTokenTransferFeeForRemoteChain (fee quoter)
	addTTFInput := operation.AddTokenTransferFeeInput{}
	addTTFReport, err := operations.ExecuteOperation(b, operation.AddTokenTransferFeeOp, deps, addTTFInput)
	if err != nil {
		return mcmstypes.BatchOperation{}, err
	}
	txs = append(txs, addTTFReport.Output...)

	return mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(deps.AptosChain.Selector),
		Transactions:  txs,
	}, nil
}
