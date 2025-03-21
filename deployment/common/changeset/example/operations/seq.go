package example

import (
	"math/big"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink/deployment/operations"
)

var LinkExampleSequence = operations.NewSequence(
	"link-example-sequence",
	semver.MustParse("1.0.0"),
	"Deploy LINK token contract, grants mint and mints some amount to same address",
	func(b operations.Bundle, deps EthereumDeps, input SqDeployLinkInput) (SqDeployLinkOutput, error) {
		linkDeployReport, err := operations.ExecuteOperation(b, DeployLinkOp, deps, operations.EmptyInput{},
			// showcasing how to disable retry, by default operations will retry up to 10 times with exponential backoff
			operations.WithRetryConfig[operations.EmptyInput, EthereumDeps](
				operations.RetryConfig[operations.EmptyInput, EthereumDeps]{
					DisableRetry: true,
				}),
		)
		if err != nil {
			return SqDeployLinkOutput{}, err
		}

		grantMintConfig := GrantMintRoleConfig{
			Contract: linkDeployReport.Output.Contract,
			To:       deps.Auth.From,
		}
		_, err = operations.ExecuteOperation(b, GrantMintOp, deps, grantMintConfig)
		if err != nil {
			return SqDeployLinkOutput{}, err
		}

		mintConfig := MintLinkConfig{
			Contract: linkDeployReport.Output.Contract,
			Amount:   input.MintAmount,
			To:       input.To,
		}
		_, err = operations.ExecuteOperation(b, MintLinkOp, deps, mintConfig,
			// showcasing how to update input before retrying the operation
			// with a nonsense example of updating the amount to 500.
			// this is useful for scenarios like updating the gas limit
			operations.WithRetryConfig[MintLinkConfig, EthereumDeps](
				operations.RetryConfig[MintLinkConfig, EthereumDeps]{
					InputHook: func(input MintLinkConfig, deps EthereumDeps) MintLinkConfig {
						input.Amount = big.NewInt(500)
						return input
					},
				}),
		)
		if err != nil {
			return SqDeployLinkOutput{}, err
		}

		return SqDeployLinkOutput{Address: linkDeployReport.Output.Address}, nil
	},
)
