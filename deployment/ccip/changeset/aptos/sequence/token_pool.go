package sequence

import (
	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

// Deploy Token Pool sequence input
type DeployTokenPoolSeqInput struct {
	TokenAddress aptos.AccountAddress
	TokenAdmin   aptos.AccountAddress
	PoolType     string
}

type DeployTokenPoolSeqOutput struct {
	TokenPoolAddress     aptos.AccountAddress
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
	// 1 - Deploy token (if not deployed)
	// TODO: Deploy token

	// 2 - Deploy token pool (if not deployed)
	deployInput := operation.DeployTokenPoolInput{
		MCMSAddress:  deps.OnChainState.MCMSAddress,
		PoolType:     in.PoolType,
		TokenAddress: in.TokenAddress,
	}
	deployReport, err := operations.ExecuteOperation(b, operation.DeployTokenPoolOp, deps, deployInput)
	if err != nil {
		return DeployTokenPoolSeqOutput{}, err
	}
	mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployReport.Output.MCMSOps)...)

	// 3 - Transfer admin role
	txs := []mcmstypes.Transaction{}
	transferInput := operation.TransferAdminRoleInput{
		Token:    in.TokenAddress,
		NewAdmin: in.TokenAdmin,
	}
	transferReport, err := operations.ExecuteOperation(b, operation.TransferAdminRoleOp, deps, transferInput)
	if err != nil {
		return DeployTokenPoolSeqOutput{}, err
	}
	txs = append(txs, transferReport.Output)

	// 4 - Accept admin role
	acceptReport, err := operations.ExecuteOperation(b, operation.AcceptAdminRoleOp, deps, in.TokenAddress)
	if err != nil {
		return DeployTokenPoolSeqOutput{}, err
	}
	txs = append(txs, acceptReport.Output)

	// 5 - Set Pool (token admin registry)
	setPoolInput := operation.SetPoolInput{
		TokenAddress: in.TokenAddress,
		PoolAddress:  deployReport.Output.TokenPoolAddress,
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
		TokenPoolAddress:     deployReport.Output.TokenPoolAddress,
		CCIPTokenPoolAddress: deployReport.Output.CCIPTokenPoolAddress,
		MCMSOperations:       mcmsOperations,
	}, nil
}

// Connect Token Pool sequence input
type ConnectTokenPoolSeqInput struct {
	TokenPoolAddress aptos.AccountAddress
	PoolType         string // TODO: should be a typed const
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
