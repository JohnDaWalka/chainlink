package ccip_attestation_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	ccip_attestation "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/ccip-attestation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	signer_registry "github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/signer-registry"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Test selector for non-Base chains
var otherSelector uint64 = ccip_attestation.BaseMainnetSelector + 1

// Helper function to create test environment with specified number of chains
func makeTestEnvironment(t *testing.T, numChains int) cldf.Environment {
	lggr := logger.TestLogger(t)
	return memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: numChains,
	})
}

// Helper function to find signer registry address in address book
func findSignerRegistryAddress(e cldf.Environment, selector uint64) (common.Address, bool) {
	addresses, err := e.ExistingAddresses.AddressesForChain(selector)
	if err != nil {
		return common.Address{}, false
	}

	for addr, tv := range addresses {
		if tv.Type == shared.EVMSignerRegistry && tv.Version == deployment.Version1_0_0 {
			return common.HexToAddress(addr), true
		}
	}
	return common.Address{}, false
}

// Helper function to bind to deployed signer registry
func bindSignerRegistry(t *testing.T, e cldf.Environment, selector uint64) *signer_registry.SignerRegistry {
	addr, found := findSignerRegistryAddress(e, selector)
	require.True(t, found, "signer registry not found for chain %d", selector)

	chain := e.BlockChains.EVMChains()[selector]
	registry, err := signer_registry.NewSignerRegistry(addr, chain.Client)
	require.NoError(t, err)
	return registry
}

// Helper function to create test signers
func makeSigners(n int) []signer_registry.ISignerRegistrySigner {
	signers := make([]signer_registry.ISignerRegistrySigner, n)
	for i := 0; i < n; i++ {
		signers[i] = signer_registry.ISignerRegistrySigner{
			EvmAddress: utils.RandomAddress(),
			// Alternate between zero and non-zero NewEVMAddress
			NewEVMAddress: func() common.Address {
				if i%2 == 0 {
					return utils.ZeroAddress
				}
				return utils.RandomAddress()
			}(),
		}
	}
	return signers
}

func TestEVMSignerRegistry_Preconditions(t *testing.T) {
	t.Parallel()

	// Create a minimal environment for precondition tests (chains won't be Base chains)
	e := makeTestEnvironment(t, 1)

	tests := []struct {
		name        string
		config      ccip_attestation.SignerRegistryChangesetConfig
		expectedErr string
	}{
		{
			name: "MaxSigners mismatch",
			config: ccip_attestation.SignerRegistryChangesetConfig{
				MaxSigners: ccip_attestation.MaxSigners - 1,
				Signers:    []signer_registry.ISignerRegistrySigner{},
			},
			expectedErr: "max signers must be",
		},
		{
			name: "Too many signers",
			config: ccip_attestation.SignerRegistryChangesetConfig{
				MaxSigners: ccip_attestation.MaxSigners,
				Signers:    makeSigners(ccip_attestation.MaxSigners + 1),
			},
			expectedErr: "too many signers",
		},
		{
			name: "Zero evm address",
			config: ccip_attestation.SignerRegistryChangesetConfig{
				MaxSigners: ccip_attestation.MaxSigners,
				Signers: []signer_registry.ISignerRegistrySigner{
					{EvmAddress: utils.ZeroAddress, NewEVMAddress: utils.RandomAddress()},
				},
			},
			expectedErr: "has zero evm address",
		},
		{
			name: "Same evm and new address",
			config: ccip_attestation.SignerRegistryChangesetConfig{
				MaxSigners: ccip_attestation.MaxSigners,
				Signers: []signer_registry.ISignerRegistrySigner{
					{EvmAddress: utils.RandomAddress(), NewEVMAddress: utils.ZeroAddress},
				},
			},
			expectedErr: "", // Will be set dynamically in test
		},
		{
			name: "Duplicate EvmAddress",
			config: func() ccip_attestation.SignerRegistryChangesetConfig {
				addr := utils.RandomAddress()
				return ccip_attestation.SignerRegistryChangesetConfig{
					MaxSigners: ccip_attestation.MaxSigners,
					Signers: []signer_registry.ISignerRegistrySigner{
						{EvmAddress: addr, NewEVMAddress: utils.RandomAddress()},
						{EvmAddress: addr, NewEVMAddress: utils.ZeroAddress},
					},
				}
			}(),
			expectedErr: "duplicate signer evm address",
		},
		{
			name: "Duplicate non-zero NewEVMAddress",
			config: func() ccip_attestation.SignerRegistryChangesetConfig {
				newAddr := utils.RandomAddress()
				return ccip_attestation.SignerRegistryChangesetConfig{
					MaxSigners: ccip_attestation.MaxSigners,
					Signers: []signer_registry.ISignerRegistrySigner{
						{EvmAddress: utils.RandomAddress(), NewEVMAddress: newAddr},
						{EvmAddress: utils.RandomAddress(), NewEVMAddress: newAddr},
					},
				}
			}(),
			expectedErr: "duplicate signer new EVM address",
		},
		{
			name: "EvmAddress equals another's NewEVMAddress",
			config: func() ccip_attestation.SignerRegistryChangesetConfig {
				addrB := utils.RandomAddress()
				return ccip_attestation.SignerRegistryChangesetConfig{
					MaxSigners: ccip_attestation.MaxSigners,
					Signers: []signer_registry.ISignerRegistrySigner{
						{EvmAddress: utils.RandomAddress(), NewEVMAddress: addrB},
						{EvmAddress: addrB, NewEVMAddress: utils.RandomAddress()},
					},
				}
			}(),
			expectedErr: "duplicate", // This should catch either evm or new address duplicate
		},
		{
			name: "Valid config with multiple zero new addresses",
			config: ccip_attestation.SignerRegistryChangesetConfig{
				MaxSigners: ccip_attestation.MaxSigners,
				Signers: []signer_registry.ISignerRegistrySigner{
					{EvmAddress: utils.RandomAddress(), NewEVMAddress: utils.ZeroAddress},
					{EvmAddress: utils.RandomAddress(), NewEVMAddress: utils.ZeroAddress},
				},
			},
			expectedErr: "", // Should succeed
		},
		{
			name: "Valid config with max signers",
			config: ccip_attestation.SignerRegistryChangesetConfig{
				MaxSigners: ccip_attestation.MaxSigners,
				Signers:    makeSigners(ccip_attestation.MaxSigners),
			},
			expectedErr: "", // Should succeed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Special handling for "Same evm and new address" test
			if tt.name == "Same evm and new address" {
				addr := utils.RandomAddress()
				tt.config.Signers = []signer_registry.ISignerRegistrySigner{
					{EvmAddress: addr, NewEVMAddress: addr},
				}
				tt.expectedErr = "has the same evm address and new evm address"
			}

			_, err := commonchangeset.Apply(t, e,
				commonchangeset.Configure(ccip_attestation.EVMSignerRegistryDeploymentChangeset, tt.config))

			if tt.expectedErr != "" {
				require.ErrorContains(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestEVMSignerRegistry_DeploysOnlyOnBaseChains(t *testing.T) {
	t.Parallel()

	// Create environment with 3 chains (none will be Base chains in test environment)
	e := makeTestEnvironment(t, 3)
	selectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))

	// Create config with test signers
	config := ccip_attestation.SignerRegistryChangesetConfig{
		MaxSigners: ccip_attestation.MaxSigners,
		Signers: []signer_registry.ISignerRegistrySigner{
			{EvmAddress: utils.RandomAddress(), NewEVMAddress: utils.ZeroAddress},
			{EvmAddress: utils.RandomAddress(), NewEVMAddress: utils.RandomAddress()},
		},
	}

	// Apply changeset
	e, err := commonchangeset.Apply(t, e,
		commonchangeset.Configure(ccip_attestation.EVMSignerRegistryDeploymentChangeset, config))
	require.NoError(t, err)

	// Load onchain state
	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)

	// Since test chains won't have Base selectors, no registries should be deployed
	for _, selector := range selectors {
		// Check address book
		_, found := findSignerRegistryAddress(e, selector)
		require.False(t, found, "signer registry should not be deployed on non-base chain %d", selector)

		// Check stateview
		chainState := state.Chains[selector]
		require.Len(t, chainState.SignerRegistrySigners, 0)
	}
}
