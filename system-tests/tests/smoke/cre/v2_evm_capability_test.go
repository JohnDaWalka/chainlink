package cre

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evmread/config"

	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"

	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	evmreadcontracts "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evmread/contracts"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"

	keystonechangeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

func executeEVMReadTest(t *testing.T, testEnv *TestEnvironment) {
	enabledChains := map[string]struct{}{}
	for _, nodeSet := range testEnv.Config.NodeSets {
		require.NoError(t, nodeSet.ParseChainCapabilities())
		if nodeSet.ChainCapabilities == nil {
			continue
		}

		for _, chainID := range nodeSet.ChainCapabilities[cre.EVMCapability].EnabledChains {
			strChainID := strconv.FormatUint(chainID, 10)
			enabledChains[strChainID] = struct{}{}
		}
	}
	require.NotEmpty(t, enabledChains, "No chains enabled for EVM read workflow test")
	const workflowFileLocation = "./evmread/main.go"
	lggr := framework.L
	var wg sync.WaitGroup
	var successfulWorkflowRuns atomic.Int32
	for _, bcOutput := range testEnv.WrappedBlockchainOutputs {
		if _, ok := enabledChains[bcOutput.BlockchainOutput.ChainID]; !ok {
			lggr.Info().Msgf("Skipping chain %s as it is not enabled for EVM read workflow test", bcOutput.BlockchainOutput.ChainID)
			continue
		}
		workflowName := "evm-read-workflow-" + bcOutput.BlockchainOutput.ChainID

		workflowConfig := configureEVMReadWorkflow(t, lggr, bcOutput)

		lggr.Info().Msg("Proceeding to register workflow...")
		compileAndDeployWorkflow(t, testEnv, lggr, fmt.Sprintf("evmreadtest-%d", bcOutput.ChainID), &workflowConfig, workflowFileLocation)

		wg.Add(1)
		forwarderAddr := common.HexToAddress(testEnv.FullCldEnvOutput.Environment.DataStore.Addresses().Filter(
			datastore.AddressRefByChainSelector(bcOutput.ChainSelector),
			datastore.AddressRefByType(datastore.ContractType(keystonechangeset.KeystoneForwarder.String())),
		)[0].Address)

		// validate workflow execution
		go func(bcOutput *cre.WrappedBlockchainOutput) {
			defer wg.Done()
			err := validateWorkflowExecution(t, lggr, bcOutput, workflowName, forwarderAddr, workflowConfig)
			if err != nil {
				lggr.Error().Msgf("Workflow %s execution failed on chain %s: %v", workflowName, bcOutput.BlockchainOutput.ChainID, err)
				return
			}

			lggr.Info().Msgf("Workflow %s executed successfully on chain %s", workflowName, bcOutput.BlockchainOutput.ChainID)
			successfulWorkflowRuns.Add(1)
		}(bcOutput)
	}

	wg.Wait()
	require.Equal(t, len(enabledChains), int(successfulWorkflowRuns.Load()), "Not all workflows executed successfully")
}

func validateWorkflowExecution(t *testing.T, lggr zerolog.Logger, bcOutput *cre.WrappedBlockchainOutput, workflowName string, forwarderAddr common.Address, cfg config.Config) error {
	forwarderContract, err := forwarder.NewKeystoneForwarder(forwarderAddr, bcOutput.SethClient.Client)
	if err != nil {
		return fmt.Errorf("failed to create forwarder contract instance: %w", err)
	}
	msgEmitterAddr := common.BytesToAddress(cfg.ContractAddress)
	isWorkflowFinished := func(ctx context.Context) (bool, error) {
		iter, err := forwarderContract.FilterReportProcessed(&bind.FilterOpts{
			Start:   cfg.ExpectedReceipt.BlockNumber.Uint64(),
			End:     nil,
			Context: ctx,
		}, []common.Address{msgEmitterAddr}, nil, nil)
		if err != nil {
			return false, fmt.Errorf("failed to filter forwarder: %w", err)
		}

		if iter.Error() != nil {
			return false, fmt.Errorf("error while filtering forwarder: %w", iter.Error())
		}

		return iter.Next(), nil
	}
	ctx, cancel := context.WithTimeout(t.Context(), testutils.WaitTimeout(t))
	defer cancel()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			lggr.Info().Msgf("Checking if workflow %s executed on chain %s", workflowName, bcOutput.BlockchainOutput.ChainID)
			ok, err := isWorkflowFinished(t.Context())
			if err != nil {
				lggr.Error().Msgf("Error checking workflow execution: %v", err)
				continue
			}

			if ok {
				lggr.Info().Msgf("Workflow %s executed successfully on chain %s", workflowName, bcOutput.BlockchainOutput.ChainID)
				return nil
			}
		case <-ctx.Done():
			return fmt.Errorf("workflow %s did not execute on chain %s within the timeout", workflowName, bcOutput.BlockchainOutput.ChainID)
		}
	}
}

func configureEVMReadWorkflow(t *testing.T, lggr zerolog.Logger, chain *cre.WrappedBlockchainOutput) config.Config {
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
	marshalledTx, err := emittingTx.MarshalBinary()
	require.NoError(t, err)
	return config.Config{
		ContractAddress:  msgEmitterContractAddr.Bytes(),
		ChainSelector:    chain.ChainSelector,
		AccountAddress:   accountAddr.Bytes(),
		ExpectedBalance:  big.NewInt(expectedBalance),
		ExpectedReceipt:  emittingReceipt,
		TxHash:           emittingReceipt.TxHash.Bytes(),
		ExpectedBinaryTx: marshalledTx,
	}
}
