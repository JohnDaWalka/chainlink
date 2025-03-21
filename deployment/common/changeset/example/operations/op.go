package example

import (
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/operations"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/link_token"
)

var DeployLinkOp = operations.NewOperation(
	"deploy-link-token",
	semver.MustParse("1.0.0"),
	"Deploy Contract Operation",
	func(b operations.Bundle, deps EthereumDeps, input operations.EmptyInput) (*deployment.ContractDeploy[*link_token.LinkToken], error) {
		linkToken, err := deployment.DeployContract[*link_token.LinkToken](b.Logger, deps.Chain, deps.AB,
			func(chain deployment.Chain) deployment.ContractDeploy[*link_token.LinkToken] {
				linkTokenAddr, tx, linkToken, err2 := link_token.DeployLinkToken(
					chain.DeployerKey,
					chain.Client,
				)
				return deployment.ContractDeploy[*link_token.LinkToken]{
					Address:  linkTokenAddr,
					Contract: linkToken,
					Tx:       tx,
					Tv:       deployment.NewTypeAndVersion(types.LinkToken, deployment.Version1_0_0),
					Err:      err2,
				}
			})
		if err != nil {
			b.Logger.Errorw("Failed to deploy link token", "chain", deps.Chain.String(), "err", err)
			return linkToken, err
		}
		return linkToken, nil
	})

type GrantMintRoleConfig struct {
	Contract *link_token.LinkToken
	To       common.Address
}

var GrantMintOp = operations.NewOperation(
	"grant-mint-role",
	semver.MustParse("1.0.0"),
	"Grant Mint Role Operation",
	func(b operations.Bundle, deps EthereumDeps, input GrantMintRoleConfig) (any, error) {
		tx, err := input.Contract.GrantMintRole(deps.Auth, input.To)
		_, err = deployment.ConfirmIfNoError(deps.Chain, tx, err)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

type MintLinkConfig struct {
	Amount   *big.Int
	To       common.Address
	Contract *link_token.LinkToken
}

var MintLinkOp = operations.NewOperation(
	"mint-link-token",
	semver.MustParse("1.0.0"),
	"Mint LINK Operation",
	func(ctx operations.Bundle, deps EthereumDeps, input MintLinkConfig) (any, error) {
		tx, err := input.Contract.Mint(deps.Auth, input.To, input.Amount)
		_, err = deployment.ConfirmIfNoError(deps.Chain, tx, err)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
