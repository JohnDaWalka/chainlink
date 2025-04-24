package tokens

import (
	"fmt"
	"maps"
	"slices"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/tokens/internal/ops"
	"github.com/smartcontractkit/chainlink/deployment/tokens/internal/seqs"
)

// DeployLinkTokenInput contains the selectors of the chains to which the Link Token contract
// should be deployed.
type DeployLinkTokenInput struct {
	ChainSelectors []uint64
	Qualifier      string
	Labels         []string
}

var _ deployment.ChangeSetV2[DeployLinkTokenInput] = DeployLinkToken{}

// DeployLinkToken deploys Link Token contracts to the specified chains in the
// DeployLinkTokenInput. The qualifiers and labels will be used to tag all deployed contracts in
// the address book and datastore.
//
// Supports deploying EVM and Solana Link Token contracts.
type DeployLinkToken struct{}

// VerifyPreconditions ensures that all listed chain selectors are valid and available in the
// environment chains.
func (DeployLinkToken) VerifyPreconditions(
	e deployment.Environment, input DeployLinkTokenInput,
) error {
	// Validate there are no duplicate chain selectors
	seen := make(map[uint64]bool)
	for _, selector := range input.ChainSelectors {
		if seen[selector] {
			return fmt.Errorf("duplicate chain selector found: %d", selector)
		}
		seen[selector] = true
	}

	// Validate there is no existing Link Token contract in the address book
	for _, csel := range input.ChainSelectors {
		addrBookByChain, err := e.ExistingAddresses.AddressesForChain(csel)
		if err != nil {
			// If the chain selector is not found in the address book, it means that
			// no contract record exists. This is not an error and we can continue.
			continue
		}

		if incl := slices.ContainsFunc(
			slices.Collect(maps.Values(addrBookByChain)),
			func(tv deployment.TypeAndVersion) bool {
				return tv.Equal(ops.LinkTokenTypeAndVersion1)
			},
		); incl {
			return fmt.Errorf("link token contract already exists for chain selector %d in address book", csel)
		}
	}

	// Validate there is no existing token contract in the datastore
	for _, csel := range input.ChainSelectors {
		if _, err := e.DataStore.Addresses().Get(datastore.NewAddressRefKey(
			csel,
			datastore.ContractType(ops.LinkTokenTypeAndVersion1.Type.String()),
			&ops.LinkTokenTypeAndVersion1.Version,
			input.Qualifier,
		)); err == nil {
			return fmt.Errorf("link token contract already exists for chain selector %d in datastore", csel)
		}
	}

	return deployment.ValidateSelectorsInEnvironment(e, input.ChainSelectors)
}

// Apply executes the SeqDeployTokens sequence to deploy Link Token contracts to the specified chains.
func (DeployLinkToken) Apply(
	e deployment.Environment, input DeployLinkTokenInput,
) (deployment.ChangesetOutput, error) {
	var (
		out = deployment.ChangesetOutput{
			AddressBook: deployment.NewMemoryAddressBook(),
			DataStore:   datastore.NewMemoryDataStore[datastore.DefaultMetadata, datastore.DefaultMetadata](),
		}

		seqDeps = seqs.SeqDeployTokensDeps{
			EVMChains: e.Chains,
			AddrBook:  out.AddressBook,
			Datastore: out.DataStore,
		}
		seqInput = seqs.SeqDeployTokensInput{
			ChainSelectors: input.ChainSelectors,
			Qualifier:      input.Qualifier,
			Labels:         input.Labels,
		}
	)

	seqReport, err := operations.ExecuteSequence(
		e.OperationsBundle, seqs.SeqDeployTokens, seqDeps, seqInput,
	)

	out.Reports = seqReport.ExecutionReports

	return out, err
}
