package operation

import (
	"encoding/json"
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/token_pool"
	mcmsbind "github.com/smartcontractkit/chainlink-aptos/bindings/mcms"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	"github.com/smartcontractkit/mcms/types"
)

type DeployTokenPoolInput struct {
	PoolType        deployment.ContractType
	TokenAddress    aptos.AccountAddress
	TokenObjAddress aptos.AccountAddress
}

type DeployTokenPoolOutput struct {
	MCMSOps              []types.Operation
	CCIPTokenPoolAddress aptos.AccountAddress
	TokenPoolAddress     aptos.AccountAddress
}

// DeployTokenPoolOp operation to update OffRamp source configurations
var DeployTokenPoolOp = operations.NewOperation(
	"deploy-token-pool-op",
	Version1_0_0,
	"Deploy Aptos token pook",
	deployTokenPool,
)

func deployTokenPool(b operations.Bundle, deps AptosDeps, in DeployTokenPoolInput) (DeployTokenPoolOutput, error) {
	// TODO: check if token pool already deployed
	var mcmsOps []types.Operation
	// Bind MCMS Package
	mcmsContract := mcmsbind.Bind(deps.OnChainState.MCMSAddress, deps.AptosChain.Client)

	// seed := "token_pool_package"
	// poolPackageAddress, err := mcmsContract.MCMSRegistry().GetNewCodeObjectAddress(nil, []byte(seed))
	// if err != nil {
	// 	return DeployTokenPoolOutput{}, fmt.Errorf("failed to GetNewCodeObjectAddress: %w", err)
	// }

	// Deploy token pool package
	// TODO: for now, deploy token pool to the same address as the token
	payload, err := token_pool.Compile(in.TokenObjAddress, deps.OnChainState.CCIPAddress, mcmsContract.Address())
	if err != nil {
		return DeployTokenPoolOutput{}, fmt.Errorf("failed to compile token pool: %w", err)
	}
	ops, err := utils.CreateChunksAndStage(payload, mcmsContract, deps.AptosChain.Selector, "", &in.TokenObjAddress)
	if err != nil {
		return DeployTokenPoolOutput{}, fmt.Errorf("failed to create chunks for token pool: %w", err)
	}
	mcmsOps = append(mcmsOps, ops...)

	ccipOwnerAddr, err := mcmsContract.MCMSRegistry().GetRegisteredOwnerAddress(nil, deps.OnChainState.CCIPAddress)
	if err != nil {
		return DeployTokenPoolOutput{}, fmt.Errorf("failed to get preexisting code object owner address: %w", err)
	}

	switch in.PoolType {
	case changeset.BurnMintTokenPool:
		payload, err = burn_mint_token_pool.Compile(
			in.TokenObjAddress,
			deps.OnChainState.CCIPAddress,
			deps.OnChainState.MCMSAddress,
			in.TokenObjAddress,
			in.TokenAddress,
			ccipOwnerAddr,
			true,
		)
		if err != nil {
			return DeployTokenPoolOutput{}, fmt.Errorf("failed to compile token pool: %w", err)
		}
		ops, err := utils.CreateChunksAndStage(payload, mcmsContract, deps.AptosChain.Selector, "", &in.TokenObjAddress)
		if err != nil {
			return DeployTokenPoolOutput{}, fmt.Errorf("failed to create chunks for token pool: %w", err)
		}
		mcmsOps = append(mcmsOps, ops...)
	case changeset.LockReleaseTokenPool:
		payload, err = lock_release_token_pool.Compile(
			in.TokenObjAddress,
			deps.OnChainState.CCIPAddress,
			deps.OnChainState.MCMSAddress,
			in.TokenObjAddress,
			in.TokenAddress,
			deps.OnChainState.MCMSAddress,
			true,
		)
		if err != nil {
			return DeployTokenPoolOutput{}, fmt.Errorf("failed to compile token pool: %w", err)
		}
		ops, err := utils.CreateChunksAndStage(payload, mcmsContract, deps.AptosChain.Selector, "", &in.TokenObjAddress)
		if err != nil {
			return DeployTokenPoolOutput{}, fmt.Errorf("failed to create chunks for token pool: %w", err)
		}
		mcmsOps = append(mcmsOps, ops...)

	default:
		return DeployTokenPoolOutput{}, fmt.Errorf("invalid token pool type: %s", in.PoolType)
	}
	return DeployTokenPoolOutput{
		MCMSOps:              mcmsOps,
		CCIPTokenPoolAddress: in.TokenObjAddress,
		TokenPoolAddress:     in.TokenObjAddress,
	}, nil
}

type SetupTokenPoolInput struct {
	TokenPoolAddress aptos.AccountAddress
	PoolType         deployment.ContractType
	RemotePools      map[uint64]RemotePool
}

type RemotePool struct {
	RemotePoolAddress  []byte
	RemoteTokenAddress []byte
	config.RateLimiterConfig
}

// SetupTokenPoolOp ...
var SetupTokenPoolOp = operations.NewOperation(
	"setup-token-pool-op",
	Version1_0_0,
	"Setup Aptos token pook",
	setupTokenPool,
)

func setupTokenPool(b operations.Bundle, deps AptosDeps, in SetupTokenPoolInput) ([]types.Transaction, error) {
	txs := []types.Transaction{}
	switch in.PoolType {
	case changeset.BurnMintTokenPool:
		bnmBind := burn_mint_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)
		var remoteChainSelectors []uint64
		var remotePoolAddresses [][][]byte
		var remoteTokenAddresses [][]byte
		var outboundIsEnableds []bool
		var outboundCapacities []uint64
		var outboundRates []uint64
		var inboundIsEnableds []bool
		var inboundCapacities []uint64
		var inboundRates []uint64

		for remoteSel, remotePool := range in.RemotePools {
			tx, err := toTransaction(bnmBind.BurnMintTokenPool().Encoder().AddRemotePool(
				remoteSel,
				remotePool.RemotePoolAddress,
			))
			if err != nil {
				return nil, fmt.Errorf("failed to encode AddRemotePool for chains: %w", err)
			}
			txs = append(txs, tx)
			remoteChainSelectors = append(remoteChainSelectors, remoteSel)
			remotePoolAddresses = append(remotePoolAddresses, [][]byte{remotePool.RemotePoolAddress})
			remoteTokenAddresses = append(remoteTokenAddresses, remotePool.RemoteTokenAddress)
			outboundIsEnableds = append(outboundIsEnableds, remotePool.OutboundIsEnabled)
			outboundCapacities = append(outboundCapacities, remotePool.OutboundCapacity)
			outboundRates = append(outboundRates, remotePool.OutboundRate)
			inboundIsEnableds = append(inboundIsEnableds, remotePool.InboundIsEnabled)
			inboundCapacities = append(inboundCapacities, remotePool.InboundCapacity)
			inboundRates = append(inboundRates, remotePool.InboundRate)
		}

		tx, err := toTransaction(bnmBind.BurnMintTokenPool().Encoder().ApplyChainUpdates(
			[]uint64{},
			remoteChainSelectors,
			remotePoolAddresses,
			remoteTokenAddresses,
		))
		if err != nil {
			return nil, fmt.Errorf("failed to encode ApplyChainUpdates for chains: %w", err)
		}
		txs = append(txs, tx)

		tx, err = toTransaction(bnmBind.BurnMintTokenPool().Encoder().SetChainRateLimiterConfigs(
			remoteChainSelectors,
			outboundIsEnableds,
			outboundCapacities,
			outboundRates,
			inboundIsEnableds,
			inboundCapacities,
			inboundRates,
		))
		if err != nil {
			return nil, fmt.Errorf("failed to encode SetChainRateLimiterConfigs for chains: %w", err)
		}
		txs = append(txs, tx)

	case changeset.LockReleaseTokenPool:
		lrBind := lock_release_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)

		var remoteChainSelectors []uint64
		var remotePoolAddresses [][][]byte
		var remoteTokenAddresses [][]byte
		var outboundIsEnableds []bool
		var outboundCapacities []uint64
		var outboundRates []uint64
		var inboundIsEnableds []bool
		var inboundCapacities []uint64
		var inboundRates []uint64

		for remoteSel, remotePool := range in.RemotePools {
			tx, err := toTransaction(lrBind.LockReleaseTokenPool().Encoder().AddRemotePool(
				remoteSel,
				remotePool.RemotePoolAddress,
			))
			if err != nil {
				return nil, fmt.Errorf("failed to encode AddRemotePool for chains: %w", err)
			}
			txs = append(txs, tx)
			remoteChainSelectors = append(remoteChainSelectors, remoteSel)
			remotePoolAddresses = append(remotePoolAddresses, [][]byte{remotePool.RemotePoolAddress})
			remoteTokenAddresses = append(remoteTokenAddresses, remotePool.RemoteTokenAddress)
			outboundIsEnableds = append(outboundIsEnableds, remotePool.OutboundIsEnabled)
			outboundCapacities = append(outboundCapacities, remotePool.OutboundCapacity)
			outboundRates = append(outboundRates, remotePool.OutboundRate)
			inboundIsEnableds = append(inboundIsEnableds, remotePool.InboundIsEnabled)
			inboundCapacities = append(inboundCapacities, remotePool.InboundCapacity)
			inboundRates = append(inboundRates, remotePool.InboundRate)
		}

		tx, err := toTransaction(lrBind.LockReleaseTokenPool().Encoder().ApplyChainUpdates(
			[]uint64{},
			remoteChainSelectors,
			remotePoolAddresses,
			remoteTokenAddresses,
		))
		if err != nil {
			return nil, fmt.Errorf("failed to encode ApplyChainUpdates for chains: %w", err)
		}
		txs = append(txs, tx)

		tx, err = toTransaction(lrBind.LockReleaseTokenPool().Encoder().SetChainRateLimiterConfigs(
			remoteChainSelectors,
			outboundIsEnableds,
			outboundCapacities,
			outboundRates,
			inboundIsEnableds,
			inboundCapacities,
			inboundRates,
		))
		if err != nil {
			return nil, fmt.Errorf("failed to encode SetChainRateLimiterConfigs for chains: %w", err)
		}
		txs = append(txs, tx)
	}

	return txs, nil
}

func toTransaction(moduleInfo bind.ModuleInformation, function string, _ []aptos.TypeTag, args [][]byte, err error) (types.Transaction, error) {
	if err != nil {
		return types.Transaction{}, err
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
		To:               moduleInfo.Address.StringLong(),
		Data:             aptosmcms.ArgsToData(args),
		AdditionalFields: afBytes,
	}, nil
}
