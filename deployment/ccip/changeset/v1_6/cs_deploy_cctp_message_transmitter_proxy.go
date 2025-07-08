package v1_6

import (
	"context"
	"fmt"

	// TODO: New token pool contract should be imported from the latest version
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/cctp_message_transmitter_proxy"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
)

var (
	DeployCCTPMessageTransmitterProxyNew = cldf.CreateChangeSet(deployCCTPMessageTransmitterProxyContractLogic, deployCCTPMessageTransmitterProxyContractPrecondition)
)

// DeployCCTPMessageTransmitterProxyInput defines all information required of the user to deploy a new CCTP message transmitter proxy contract.
type DeployCCTPMessageTransmitterProxyInput struct {
}

func (i DeployCCTPMessageTransmitterProxyInput) Validate(ctx context.Context, chain cldf_evm.Chain, state evm.CCIPChainState) error {

	return nil
}

// DeployCCTPMessageTransmitterProxyContractConfig
type DeployCCTPMessageTransmitterProxyContractConfig struct {
}

func deployCCTPMessageTransmitterProxyContractPrecondition(env cldf.Environment, c DeployCCTPMessageTransmitterProxyContractConfig) error {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
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

	fmt.Println(state)

	return cldf.ChangesetOutput{
		AddressBook: newAddresses, // TODO: this is deprecated, how do I use the DataStore instead?
	}, nil
}
