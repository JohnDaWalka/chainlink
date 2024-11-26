package crib

import (
	"context"
	"errors"

	chainsel "github.com/smartcontractkit/chain-selectors"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/deployment/ccip"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

// DeployHomeChainContracts deploys the home chain contracts so that the chainlink nodes can be started with the CR address in Capabilities.ExternalRegistry
func DeployHomeChainContracts(lggr logger.Logger, envConfig devenv.EnvironmentConfig, homeChainSel uint64) (deployment.CapabilityRegistryConfig, deployment.AddressBook, error) {
	chains, err := devenv.NewChains(lggr, envConfig.Chains)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}

	ab := deployment.NewMemoryAddressBook()
	capReg, err := ccipdeployment.DeployCapReg(lggr, ab, chains[homeChainSel])
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}
	evmChainID, err := chainsel.ChainIdFromSelector(homeChainSel)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}
	return deployment.CapabilityRegistryConfig{
		NetworkType: relay.NetworkEVM,
		EVMChainID:  evmChainID,
		Contract:    capReg.Address,
	}, ab, nil
}

func DeployCCIPAndAddLanes(lggr logger.Logger, envCfg devenv.EnvironmentConfig, homeChainSel, feedChainSel uint64, ab deployment.AddressBook) (DeployCCIPOutput, error) {
	e, _, err := devenv.NewEnvironment(context.Background(), lggr, envCfg)
	if err != nil {
		return DeployCCIPOutput{}, err
	}
	if e == nil {
		return DeployCCIPOutput{}, errors.New("environment is nil")
	}

	_, err = ccipdeployment.DeployFeeds(lggr, ab, e.Chains[feedChainSel])
	if err != nil {
		return DeployCCIPOutput{AddressBook: e.ExistingAddresses}, err
	}
	err = ccipdeployment.DeployFeeTokensToChains(lggr, ab, e.Chains)
	if err != nil {
		return DeployCCIPOutput{AddressBook: e.ExistingAddresses}, err
	}
	e.ExistingAddresses = ab
	tenv := ccipdeployment.DeployedEnv{
		Env:          *e,
		HomeChainSel: homeChainSel,
		FeedChainSel: feedChainSel,
	}

	state, err := ccipdeployment.LoadOnchainState(tenv.Env)
	if err != nil {
		return DeployCCIPOutput{AddressBook: e.ExistingAddresses}, err
	}
	if state.Chains[tenv.HomeChainSel].LinkToken == nil {
		return DeployCCIPOutput{AddressBook: e.ExistingAddresses}, errors.New("link token not deployed")
	}

	feeds := state.Chains[tenv.FeedChainSel].USDFeeds
	tokenConfig := ccipdeployment.NewTestTokenConfig(feeds)
	mcmsCfg, err := ccipdeployment.NewTestMCMSConfig(tenv.Env)
	if err != nil {
		return DeployCCIPOutput{e.ExistingAddresses}, err
	}
	output, err := changeset.InitialDeploy(tenv.Env, ccipdeployment.DeployCCIPContractConfig{
		HomeChainSel:   tenv.HomeChainSel,
		FeedChainSel:   tenv.FeedChainSel,
		ChainsToDeploy: tenv.Env.AllChainSelectors(),
		TokenConfig:    tokenConfig,
		MCMSConfig:     mcmsCfg,
		OCRSecrets:     deployment.XXXGenerateTestOCRSecrets(),
	})
	if err != nil {
		return DeployCCIPOutput{AddressBook: e.ExistingAddresses}, err
	}
	err = tenv.Env.ExistingAddresses.Merge(output.AddressBook)
	if err != nil {
		return DeployCCIPOutput{AddressBook: e.ExistingAddresses}, err
	}
	// Get new state after migration.
	state, err = ccipdeployment.LoadOnchainState(tenv.Env)
	if err != nil {
		return DeployCCIPOutput{AddressBook: e.ExistingAddresses}, err
	}

	// Apply the jobs.
	for nodeID, jobs := range output.JobSpecs {
		for _, job := range jobs {
			// Note these auto-accept
			_, err := tenv.Env.Offchain.ProposeJob(context.Background(),
				&jobv1.ProposeJobRequest{
					NodeId: nodeID,
					Spec:   job,
				})
			if err != nil {
				return DeployCCIPOutput{AddressBook: e.ExistingAddresses}, err
			}
		}
	}

	// Add all lanes
	err = ccipdeployment.AddLanesForAll(tenv.Env, state)
	if err != nil {
		return DeployCCIPOutput{AddressBook: e.ExistingAddresses}, err
	}
	err = tenv.Env.ExistingAddresses.Merge(output.AddressBook)
	return DeployCCIPOutput{AddressBook: e.ExistingAddresses}, err
}
