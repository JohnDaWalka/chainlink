package seqs

import (
	"github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/tokens/internal/ops"
)

// SeqDeployTokensDeps contains the dependencies for the SeqDeployTokens sequence.
type SeqDeployTokensDeps struct {
	EVMChains map[uint64]deployment.Chain
	SolChains map[uint64]deployment.SolChain
	AddrBook  deployment.AddressBook
	Datastore datastore.MutableDataStore[
		datastore.DefaultMetadata,
		datastore.DefaultMetadata,
	]
}

// SeqDeployTokensInput is the input to the SeqDeployTokens sequence.
type SeqDeployTokensInput struct {
	// ChainSelectors are the chain selectors of the chains to which the Link Token contract
	ChainSelectors []uint64 `json:"chainSelectors"`

	// Qualifier is a string that will be used to tag the deployed contracts in the address book and datastore.
	Qualifier string `json:"qualifier"`

	// Labels are a list of labels that will be used to tag the deployed contracts in the address book and datastore.
	Labels []string `json:"labels"`
}

// SeqDeployTokensOutput is the output of the SeqDeployTokens sequence.
type SeqDeployTokensOutput struct {
	// Addresses are the addresses of the deployed Link Token contracts.
	Addresses []string `json:"address"`
}

// SeqDeployTokens is a sequence that deploys LINK token contracts across multiple chains.
var SeqDeployTokens = operations.NewSequence(
	"seq-deploy-tokens",
	semver.MustParse("1.0.0"),
	"Deploy LINK token contracts across multiple chains",
	func(b operations.Bundle, deps SeqDeployTokensDeps, input SeqDeployTokensInput) (SeqDeployTokensOutput, error) {
		out := SeqDeployTokensOutput{
			Addresses: make([]string, 0),
		}

		for _, csel := range input.ChainSelectors {
			fam, err := chainsel.GetSelectorFamily(csel)
			if err != nil {
				return out, err
			}

			switch fam {
			case chainsel.FamilyEVM:
				chain := deps.EVMChains[csel]

				// Deploy the link token
				deployReport, err := operations.ExecuteOperation(b, ops.OpEVMDeployLinkToken,
					ops.OpEVMDeployLinkTokenDeps{
						Auth:        chain.DeployerKey,
						Backend:     chain.Client,
						ConfirmFunc: chain.Confirm,
					},
					ops.OpEVMDeployLinkTokenInput{
						ChainSelector: csel,
						ChainName:     chain.Name(),
					},
				)
				if err != nil {
					return out, err
				}

				_, err = operations.ExecuteSequence(b, SeqPersistAddress,
					SeqPersistAddressDeps{
						AddrBook:  deps.AddrBook,
						Datastore: deps.Datastore,
					},
					SeqPersistAddressInput{
						ChainSelector: csel,
						Address:       deployReport.Output.Address.String(),
						Type:          deployReport.Output.Type,
						Version:       deployReport.Output.Version,
						Qualifier:     input.Qualifier,
						Labels:        input.Labels,
					},
				)
				if err != nil {
					return out, err
				}

				out.Addresses = append(out.Addresses, deployReport.Output.Address.String())
			case chainsel.FamilySolana:
				chain := deps.SolChains[csel]

				deployReport, err := operations.ExecuteOperation(b, ops.OpSolDeployLinkToken,
					ops.OpSolDeployLinkTokenDeps{
						Client:      chain.Client,
						ConfirmFunc: chain.Confirm,
					},
					ops.OpSolDeployLinkTokenInput{
						ChainSelector:       csel,
						ChainName:           chain.Name(),
						TokenAdminPublicKey: chain.DeployerKey.PublicKey(),
					},
				)
				if err != nil {
					return out, err
				}

				_, err = operations.ExecuteSequence(b, SeqPersistAddress,
					SeqPersistAddressDeps{
						AddrBook:  deps.AddrBook,
						Datastore: deps.Datastore,
					},
					SeqPersistAddressInput{
						ChainSelector: csel,
						Address:       deployReport.Output.MintPublicKey.String(),
						Type:          deployReport.Output.Type,
						Version:       deployReport.Output.Version,
						Qualifier:     input.Qualifier,
						Labels:        input.Labels,
					},
				)
				if err != nil {
					return out, err
				}

				out.Addresses = append(out.Addresses, deployReport.Output.MintPublicKey.String())
			}
		}

		return out, nil
	},
)
