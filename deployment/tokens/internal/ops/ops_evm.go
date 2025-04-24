package ops

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	link_token "github.com/smartcontractkit/chainlink-testing-framework/seth/contracts/bind/link"
	"github.com/smartcontractkit/chainlink/deployment"
)

var (
	// LinkToken is the burn/mint link token which is now used in all new deployments.
	// https://github.com/smartcontractkit/chainlink/blob/develop/core/gethwrappers/shared/generated/link_token/link_token.go#L34
	LinkTokenTypeAndVersion1 = deployment.NewTypeAndVersion(
		"LinkToken",
		*semver.MustParse("1.0.0"),
	)
)

// OpEVMDeployLinkTokenDeps defines the dependencies to perform the OpEVMDeployLinkToken
// operation.
type OpEVMDeployLinkTokenDeps struct {
	Auth        *bind.TransactOpts
	Backend     bind.ContractBackend
	ConfirmFunc func(tx *types.Transaction) (uint64, error)
}

// OpEVMDeployLinkTokenInput represents the input parameters for the OpEVMDeployLinkToken operation.
// The chainSelector and chainName fields are used to identify the target chain for deployment
// and to generate a unique cache key for the report. These fields do not directly affect the
// deployment logic.
type OpEVMDeployLinkTokenInput struct {
	// ChainSelector is the unique identifier for the chain where the operation will be executed.
	// It is used as part of the unique cache key for the report and in logging.
	ChainSelector uint64 `json:"chainSelector"`
	// ChainName is the human friendly name of the chain. This is only used for logging.
	ChainName string `json:"chainName"`
}

// OpEvmDeployLinkTokenOutput represents the output of the OpEVMDeployLinkToken operation.
type OpEvmDeployLinkTokenOutput struct {
	Address common.Address `json:"address"`
	Type    string         `json:"type"`
	Version string         `json:"version"`
}

// OpEVMDeployLinkToken is an operation that deploys the LINK token contract
// on an EVM-compatible blockchain.
var OpEVMDeployLinkToken = operations.NewOperation(
	"evm-deploy-link-token",
	semver.MustParse("1.0.0"),
	"Deploy EVM LINK Token Contract",
	func(b operations.Bundle, deps OpEVMDeployLinkTokenDeps, input OpEVMDeployLinkTokenInput) (OpEvmDeployLinkTokenOutput, error) {
		out := OpEvmDeployLinkTokenOutput{}

		// Deploy the link token
		addr, tx, _, err := link_token.DeployLinkToken(
			deps.Auth,
			deps.Backend,
		)
		if err != nil {
			b.Logger.Errorw("Failed to deploy link token",
				"chainSelector", input.ChainSelector,
				"chainName", input.ChainName,
				"err", err,
			)

			return out, err
		}

		// Confirm the transaction
		if _, err = deps.ConfirmFunc(tx); err != nil {
			b.Logger.Errorw("Failed to confirm deployment",
				"chainSelector", input.ChainSelector,
				"chainName", input.ChainName,
				"contractAddr", addr.String(),
				"err", err,
			)

			return out, err
		}

		return OpEvmDeployLinkTokenOutput{
			Address: addr,
			Type:    LinkTokenTypeAndVersion1.Type.String(),
			Version: LinkTokenTypeAndVersion1.Version.String(),
		}, nil
	})
