package cre

import (
	"encoding/json"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	evmreadcontracts "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evmread/contracts"
	workflowTypes "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evmread/types"

	keystonechangeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
)

func executeEVMReadTest(t *testing.T, in *envconfig.Config, envArtifact environment.EnvArtifact) {
	lggr := framework.L
	cldLogger := cldlogger.NewSingleFileLogger(t)

	workflowFileLocation := "./evmread/main.go"

	/*
		BUILD ENVIRONMENT FROM SAVED STATE
	*/
	fullCldEnvOutput, wrappedBlockchainOutputs, loadErr := environment.BuildFromSavedState(t.Context(), cldLogger, in, envArtifact)
	require.NoError(t, loadErr, "failed to load environment")

	homeChain := wrappedBlockchainOutputs[0]
	var wg sync.WaitGroup
	for _, bcOutput := range wrappedBlockchainOutputs {
		/*
			REGISTER WORKFLOW FOR EACH CHAIN
		*/
		workflowName := "evm-read-workflow-" + bcOutput.BlockchainOutput.ChainID

		workflowConfig := configureEVMReadWorkflow(t, lggr, bcOutput)
		marshaledCfg, err := json.Marshal(workflowConfig)
		require.NoError(t, err, "failed to marshal workflow config")
		workflowConfigFilePath, err := createConfigFile(t, workflowName, marshaledCfg)
		require.NoError(t, err, "failed to create workflow config file")

		lggr.Info().Msg("Proceeding to register workflow...")

		deployWorkflow(t, homeChain, fullCldEnvOutput,
			workflowConfigFilePath, workflowFileLocation, workflowName)

		wg.Add(1)
		forwarderAddr, err := crecontracts.FindAddressesForChain(
			fullCldEnvOutput.Environment.ExistingAddresses, //nolint:staticcheck // won't migrate now
			bcOutput.ChainSelector,
			keystonechangeset.KeystoneForwarder.String(),
		)
		require.NoError(t, err)
		// validate workflow execution
		go func(forwarderAddr common.Address, cfg workflowTypes.Config) {
			defer wg.Done()
			forwarderContract, err := forwarder.NewKeystoneForwarder(forwarderAddr, bcOutput.SethClient.Client)
			require.NoError(t, err)
			require.Eventually(t, func() bool {
				msgEmitterAddr := common.BytesToAddress(cfg.ContractAddress)
				msgs, err := forwarderContract.FilterReportProcessed(&bind.FilterOpts{
					Start:   pb.NewIntFromBigInt(cfg.ExpectedReceipt.BlockNumber).Uint64(),
					End:     nil,
					Context: t.Context(),
				}, []common.Address{msgEmitterAddr}, nil, nil)
				if err != nil {
					lggr.Error().Err(err).Msg("failed to filter messages emitted by contract")
					return false
				}

				for {
					if msgs.Error() != nil {
						lggr.Error().Err(msgs.Error()).Msg("error while reading messages emitted by contract")
						return false
					}
					if !msgs.Next() {
						return false
					}
					lggr.Info().Msgf("Workflow sucesfully executed on chain %s", bcOutput.BlockchainOutput.ChainID)
					return true
				}
			}, time.Minute*6, time.Second)
		}(forwarderAddr, workflowConfig)
	}

	wg.Wait()
}

func configureEVMReadWorkflow(t *testing.T, lggr zerolog.Logger, chain *cre.WrappedBlockchainOutput) workflowTypes.Config {
	lggr.Info().Msgf("Deploying message emitter for chain %s", chain.BlockchainOutput.ChainID)
	msgEmitterContractAddr, tx, msgEmitter, err := evmreadcontracts.DeployMessageEmitter(chain.SethClient.NewTXOpts(), chain.SethClient.Client)
	require.NoError(t, err)
	lggr.Info().Msgf("Deployed message emitter for chain %s at %s", chain.BlockchainOutput.ChainID, msgEmitterContractAddr.String())
	_, err = chain.SethClient.WaitMined(t.Context(), lggr, chain.SethClient.Client, tx)
	require.NoError(t, err)
	lggr.Printf("Emitting event to be picked up by workflow for chain %s", chain.BlockchainOutput.ChainID)
	emittingTx, err := msgEmitter.EmitMessage(chain.SethClient.NewTXOpts(), "Initial message to be read by workflow")
	require.NoError(t, err)
	emittingReceipt, err := chain.SethClient.WaitMined(t.Context(), lggr, chain.SethClient.Client, emittingTx)
	require.NoError(t, err)
	lggr.Info().Msgf("Updating nonces for chain %s", chain.BlockchainOutput.ChainID)
	// force update nonces to ensure the transfer works
	require.NoError(t, chain.SethClient.NonceManager.UpdateNonces())
	const expectedBalance = 10
	pk, err := crypto.GenerateKey()
	require.NoError(t, err)
	accountAddr := crypto.PubkeyToAddress(pk.PublicKey)
	lggr.Info().Msgf("Funding account %s for BalanceAt read test for chain %s", accountAddr.Hex(), chain.BlockchainOutput.ChainID)
	err = chain.SethClient.TransferETHFromKey(t.Context(), 0, accountAddr.Hex(), big.NewInt(expectedBalance), nil)
	require.NoError(t, err, "failed to transfer ETH to contract %s", msgEmitterContractAddr.String())

	return workflowTypes.Config{
		ContractAddress: msgEmitterContractAddr.Bytes(),
		ChainSelector:   chain.ChainSelector,
		AccountAddress:  accountAddr.Bytes(),
		ExpectedBalance: big.NewInt(expectedBalance),
		ExpectedReceipt: &evm.Receipt{
			Status:            emittingReceipt.Status,
			Logs:              make([]*evm.Log, len(emittingReceipt.Logs)), // workflow compares only number of logs, not their content
			TxHash:            emittingReceipt.TxHash.Bytes(),
			ContractAddress:   emittingReceipt.ContractAddress.Bytes(),
			GasUsed:           emittingReceipt.GasUsed,
			BlockHash:         emittingReceipt.BlockHash.Bytes(),
			BlockNumber:       pb.NewBigIntFromInt(emittingReceipt.BlockNumber),
			TxIndex:           uint64(emittingReceipt.TransactionIndex),
			EffectiveGasPrice: pb.NewBigIntFromInt(emittingReceipt.EffectiveGasPrice),
		},
		TxHash: emittingReceipt.TxHash.Bytes(),
		ExpectedTx: &evm.Transaction{
			Nonce:    emittingTx.Nonce(),
			Gas:      emittingTx.Gas(),
			To:       emittingTx.To().Bytes(),
			Data:     emittingTx.Data(),
			Hash:     emittingTx.Hash().Bytes(),
			Value:    pb.NewBigIntFromInt(emittingTx.Value()),
			GasPrice: pb.NewBigIntFromInt(emittingTx.GasPrice()),
		},
	}
}
