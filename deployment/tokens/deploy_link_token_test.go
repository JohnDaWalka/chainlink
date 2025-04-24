package tokens

import (
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/tokens/internal/ops"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func Test_DeployLinkToken_VerifyPreconditions(t *testing.T) {
	t.Parallel()

	var (
		csel uint64 = 16015286601757825753 // Sepolia
	)

	tests := []struct {
		name       string
		beforeFunc func(t *testing.T, e *deployment.Environment)
		input      DeployLinkTokenInput
		wantErr    string
	}{
		{
			name: "valid input",
			beforeFunc: func(t *testing.T, e *deployment.Environment) {
				// Inject the environment chain. The actual value does not matter for validation.
				e.Chains = map[uint64]deployment.Chain{
					csel: {},
				}

				// Inject empty address book and datastore
				e.ExistingAddresses = deployment.NewMemoryAddressBook()
				e.DataStore = datastore.NewMemoryDataStore[
					datastore.DefaultMetadata, datastore.DefaultMetadata,
				]().Seal()
			},
			input: DeployLinkTokenInput{
				ChainSelectors: []uint64{csel},
			},
		},
		{
			name: "fail: duplicate chain selectors",
			input: DeployLinkTokenInput{
				ChainSelectors: []uint64{1, 1},
			},
			wantErr: "duplicate chain selector found",
		},
		{
			name: "fail: link token contracts exists in address book",
			beforeFunc: func(t *testing.T, e *deployment.Environment) {
				t.Helper()

				e.ExistingAddresses = deployment.NewMemoryAddressBook()
				err := e.ExistingAddresses.Save(csel,
					"0xeC91988D7dD84d8adE801b739172ad15c860A700",
					ops.LinkTokenTypeAndVersion1,
				)
				require.NoError(t, err)
			},
			input: DeployLinkTokenInput{
				ChainSelectors: []uint64{csel},
			},
			wantErr: "link token contract already exists for chain selector",
		},
		{
			name: "fail: link token contract exists in datastore",
			beforeFunc: func(t *testing.T, e *deployment.Environment) {
				t.Helper()

				// Insert the selector with no addresses to pass address book check
				e.ExistingAddresses = deployment.NewMemoryAddressBookFromMap(
					map[uint64]map[string]deployment.TypeAndVersion{
						csel: {},
					},
				)

				ds := datastore.NewMemoryDataStore[datastore.DefaultMetadata, datastore.DefaultMetadata]()
				err := ds.Addresses().Add(datastore.AddressRef{
					ChainSelector: csel,
					Address:       "0xeC91988D7dD84d8adE801b739172ad15c860A700",
					Type:          datastore.ContractType(ops.LinkTokenTypeAndVersion1.Type.String()),
					Version:       &ops.LinkTokenTypeAndVersion1.Version,
				})
				require.NoError(t, err)

				e.DataStore = ds.Seal()
			},
			input: DeployLinkTokenInput{
				ChainSelectors: []uint64{csel},
			},
			wantErr: "link token contract already exists for chain selector",
		},
		{
			name: "fail: chain selector not found in environment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				cs  = DeployLinkToken{}
				env = &deployment.Environment{}
			)

			if tt.beforeFunc != nil {
				tt.beforeFunc(t, env)
			}

			err := cs.VerifyPreconditions(*env, tt.input)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_DeployLinkToken_Apply(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		giveFunc func(e deployment.Environment) DeployLinkTokenInput
	}{
		{
			name: "valid input",
			giveFunc: func(e deployment.Environment) DeployLinkTokenInput {
				csels := e.AllChainSelectorsAllFamilies()

				return DeployLinkTokenInput{
					ChainSelectors: csels,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			lggr := logger.Test(t)

			e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
				Chains:    1,
				SolChains: 0, // TODO: Add Solana chain test
			})

			var (
				cs = DeployLinkToken{}
			)

			got, err := cs.Apply(e, tt.giveFunc(e))
			require.NoError(t, err)

			// Check that the address book has the link token contract for each chain
			for _, csel := range e.AllChainSelectors() {
				addrBookByChain, err := got.AddressBook.AddressesForChain(csel)
				require.NoError(t, err)
				require.NotEmpty(t, addrBookByChain)
				require.Len(t, addrBookByChain, 1)
			}

			addrRefs, err := got.DataStore.Addresses().Fetch()
			require.NoError(t, err)
			require.Len(t, addrRefs, 1)
		})
	}
}
