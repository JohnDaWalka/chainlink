package environment

import (
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
)

type StartBlockchainsOutput struct {
	Outputs     []*BlockchainOutput // TODO: This cannot be an output... It contains a private key, and we don't want to serialize it
	Blockchains map[uint64]chain.BlockChain
}

type StartBlockchainsDeps struct {
	logger          zerolog.Logger
	singeFileLogger *cldlogger.SingleFileLogger
}

var StartBlockchainsOp = operations.NewOperation[BlockchainsInput, StartBlockchainsOutput, StartBlockchainsDeps](
	"start-blockchains-op",
	semver.MustParse("1.0.0"),
	"Start Blockchains",
	func(b operations.Bundle, deps StartBlockchainsDeps, input BlockchainsInput) (StartBlockchainsOutput, error) {
		blockchainsOutput, bcOutErr := CreateBlockchains(deps.logger, input)
		if bcOutErr != nil {
			return StartBlockchainsOutput{}, pkgerrors.Wrap(bcOutErr, "failed to create blockchains")
		}

		var chainsConfigs []devenv.ChainConfig

		for _, bcOut := range blockchainsOutput {
			chainsConfigs = append(chainsConfigs, devenv.ChainConfig{
				ChainID:   strconv.FormatUint(bcOut.SethClient.Cfg.Network.ChainID, 10),
				ChainName: bcOut.SethClient.Cfg.Network.Name,
				ChainType: strings.ToUpper(bcOut.BlockchainOutput.Family),
				WSRPCs: []devenv.CribRPCs{{
					External: bcOut.BlockchainOutput.Nodes[0].ExternalWSUrl,
					Internal: bcOut.BlockchainOutput.Nodes[0].InternalWSUrl,
				}},
				HTTPRPCs: []devenv.CribRPCs{{
					External: bcOut.BlockchainOutput.Nodes[0].ExternalHTTPUrl,
					Internal: bcOut.BlockchainOutput.Nodes[0].InternalHTTPUrl,
				}},
				DeployerKey: bcOut.SethClient.NewTXOpts(seth.WithNonce(nil)), // set nonce to nil, so that it will be fetched from the RPC node
			})
		}

		allChains, _, allChainsErr := devenv.NewChains(deps.singeFileLogger, chainsConfigs)
		if allChainsErr != nil {
			return StartBlockchainsOutput{}, pkgerrors.Wrap(allChainsErr, "failed to create chains")
		}

		blockChains := map[uint64]chain.BlockChain{}
		for selector, ch := range allChains {
			blockChains[selector] = ch
		}

		return StartBlockchainsOutput{
			Outputs:     blockchainsOutput,
			Blockchains: blockChains,
		}, nil
	},
)
