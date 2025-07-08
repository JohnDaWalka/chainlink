package v1_6

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	// TODO: New token pool contract should be imported from the latest version
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/cctp_message_transmitter_proxy"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
)

var (
	DeployCCTPMessageTransmitterProxyNew = cldf.CreateChangeSet(deployCCTPMessageTransmitterProxyContractLogic, deployCCTPMessageTransmitterProxyContractPrecondition)
)

// DeployCCTPMessageTransmitterProxyInput defines all information required of the user to deploy a new CCTP message transmitter proxy contract.
type DeployCCTPMessageTransmitterProxyInput struct {
	// TokenMessenger is the address of the USDC token messenger contract.
	TokenMessenger common.Address
}

func (i DeployCCTPMessageTransmitterProxyInput) Validate(ctx context.Context, chain cldf_evm.Chain, state evm.CCIPChainState) error {

	return nil
}

// DeployCCTPMessageTransmitterProxyContractConfig
type DeployCCTPMessageTransmitterProxyContractConfig struct {
	USDCProxies map[uint64]DeployCCTPMessageTransmitterProxyInput
}

func deployCCTPMessageTransmitterProxyContractPrecondition(env cldf.Environment, c DeployCCTPMessageTransmitterProxyContractConfig) error {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	for chainSelector, proxyConfig := range c.USDCProxies {
		chain, chainState, err := state.GetEVMChainState(env, chainSelector)
		if err != nil {
			return fmt.Errorf("failed to get EVM chain state for chain selector %d: %w", chainSelector, err)
		}
		err = proxyConfig.Validate(env.GetContext(), chain, chainState)
		if err != nil {
			return fmt.Errorf("failed to validate USDC token pool config for chain selector %d: %w", chainSelector, err)
		}
	}
	return nil
}

// DeployCCTPMessageTransmitterProxyContractChangeset deploys new USDC pools across multiple chains.
func deployCCTPMessageTransmitterProxyContractLogic(env cldf.Environment, c DeployCCTPMessageTransmitterProxyContractConfig) (cldf.ChangesetOutput, error) {
	if err := deployCCTPMessageTransmitterProxyContractPrecondition(env, c); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid DeployCCTPMessageTransmitterProxyContractConfig: %w", err)
	}
	newAddresses := cldf.NewMemoryAddressBook()

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

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

	return cldf.ChangesetOutput{
		AddressBook: newAddresses, // TODO: this is deprecated, how do I use the DataStore instead?
	}, nil
}
