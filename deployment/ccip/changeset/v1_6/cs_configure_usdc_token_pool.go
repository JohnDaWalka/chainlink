package v1_6

import (
	"context"
	"fmt"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
)

var (
	ConfigUSDCTokenPool = cldf.CreateChangeSet(configUSDCTokenPoolLogic, configUSDCTokenPoolPrecondition)
)

// ConfigCCTPMessageTransmitterProxyInput defines all information required of the user to deploy a new CCTP message transmitter proxy contract.
type ConfigCCTPMessageTransmitterProxyInput struct {
}

func (i ConfigCCTPMessageTransmitterProxyInput) Validate(ctx context.Context, chain cldf_evm.Chain, state evm.CCIPChainState) error {

	return nil
}

// ConfigUSDCTokenPoolConfig defines the configuration for deploying CCTP message transmitter proxy contracts.
type ConfigUSDCTokenPoolConfig struct {
}

func configUSDCTokenPoolPrecondition(env cldf.Environment, c ConfigUSDCTokenPoolConfig) error {
	return nil
}

// ConfigUSDCTokenPoolChangeset deploys new USDC pools across multiple chains.
func configUSDCTokenPoolLogic(env cldf.Environment, c ConfigUSDCTokenPoolConfig) (cldf.ChangesetOutput, error) {
	if err := configUSDCTokenPoolPrecondition(env, c); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid ConfigUSDCTokenPoolConfig: %w", err)
	}
	newAddresses := cldf.NewMemoryAddressBook()

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	fmt.Println(state)
	/*
		for chainSelector, proxyConfig := range c.USDCProxies {
			chain, _, err := state.GetEVMChainState(env, chainSelector)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to get EVM chain state for chain selector %d: %w", chainSelector, err)
			}
			_, err = cldf.DeployContract(env.Logger, chain, newAddresses,
				func(chain cldf_evm.Chain) cldf.ContractDeploy[*cctp_message_transmitter_proxy.CCTPMessageTransmitterProxy] {
					proxyAddress, tx, proxy, err := cctp_message_transmitter_proxy.DeployCCTPMessageTransmitterProxy(
						chain.DeployerKey,          // auth
						chain.Client,               // backend
						proxyConfig.TokenMessenger, // tokenMessenger
					)
					return cldf.ContractDeploy[*cctp_message_transmitter_proxy.CCTPMessageTransmitterProxy]{
						Address:  proxyAddress,
						Contract: proxy,
						Tv:       cldf.NewTypeAndVersion("TODO: The proxy contract has a name", deployment.Version1_6_0),
						Tx:       tx,
						Err:      err,
					}
				},
			)
		}
	*/

	return cldf.ChangesetOutput{
		AddressBook: newAddresses, // TODO: this is deprecated, how do I use the DataStore instead?
	}, nil
}
