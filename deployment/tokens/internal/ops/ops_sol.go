package ops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment"
)

const (
	// TokenDecimalsSolana is the number of decimals for the Solana LINK token.
	tokenDecimalsSolana = 9
)

// OpSolDeployLinkToken deploys a LINK token contract on Solana.
type OpSolDeployLinkTokenDeps struct {
	Client      *solRpc.Client
	ConfirmFunc func(instructions []solana.Instruction, opts ...solCommonUtil.TxModifier) error
}

// OpSolDeployLinkTokenInput represents the input parameters for the OpSolDeployLinkToken operation.
// The chainSelector and chainName fields are used to identify the target chain for deployment
// and to generate a unique cache key for the report. These fields do not directly affect the
// deployment logic.
type OpSolDeployLinkTokenInput struct {
	// ChainSelector is the unique identifier for the chain where the operation will be executed.
	// It is used as part of the unique cache key for the report and in logging.
	ChainSelector uint64 `json:"chainSelector"`
	// ChainName is the human friendly name of the chain. This is only used for logging.
	ChainName string `json:"chainName"`
	// TokenAdminPublicKey is the public key of the admin account for the token.
	TokenAdminPublicKey solana.PublicKey `json:"tokenAdminPublicKey"`
}

// OpSolDeployLinkTokenOutput represents the output of the OpSolDeployLinkToken operation.
type OpSolDeployLinkTokenOutput struct {
	// MintPublicKey is represents the token's address
	MintPublicKey solana.PublicKey `json:"mintPublicKey"`
	Type          string           `json:"type"`
	Version       string           `json:"version"`
}

// OpSolDeployLinkToken is an operation that deploys the LINK token contract on the Solana
// blockchain.
var OpSolDeployLinkToken = operations.NewOperation(
	"sol-deploy-link-token",
	semver.MustParse("1.0.0"),
	"Deploy Solana LINK Token Contract",
	func(b operations.Bundle, deps OpSolDeployLinkTokenDeps, input OpSolDeployLinkTokenInput) (OpSolDeployLinkTokenOutput, error) {
		out := OpSolDeployLinkTokenOutput{}

		// Generate the publicKey of the new token mint
		mint, _ := solana.NewRandomPrivateKey()
		mintPublicKey := mint.PublicKey() // This is the token address

		// Create the token
		instructions, err := solTokenUtil.CreateToken(
			b.GetContext(),
			solana.Token2022ProgramID,
			mintPublicKey,
			input.TokenAdminPublicKey,
			tokenDecimalsSolana,
			deps.Client,
			deployment.SolDefaultCommitment,
		)
		if err != nil {
			b.Logger.Errorw("Failed to generate instructions for link token deployment",
				"chainSelector", input.ChainSelector,
				"chainName", input.ChainName,
				"err", err,
			)

			return out, fmt.Errorf("Failed to generate instructions for link token deployment: %w", err)
		}

		// Confirm the transaction
		if err = deps.ConfirmFunc(instructions, solCommonUtil.AddSigners(mint)); err != nil {
			b.Logger.Errorw("Failed to confirm instructions for link token deployment",
				"chainSelector", input.ChainSelector,
				"chainName", input.ChainName,
				"err", err,
			)

			return out, err
		}

		return OpSolDeployLinkTokenOutput{
			MintPublicKey: mintPublicKey,
			Type:          LinkTokenTypeAndVersion1.Type.String(),
			Version:       LinkTokenTypeAndVersion1.Version.String(),
		}, nil
	})
