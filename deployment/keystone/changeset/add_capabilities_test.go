package changeset_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func TestAddCapabilities(t *testing.T) {
	t.Parallel()

	t.Run("no mcms", func(t *testing.T) {
		te := test.SetupTestEnv(t, test.TestConfig{
			WFDonConfig:     test.DonConfig{N: 4},
			AssetDonConfig:  test.DonConfig{N: 4},
			WriterDonConfig: test.DonConfig{N: 4},
			NumChains:       1,
		})

		cfg := changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "test-cap",
				Version:        "0.0.1",
				CapabilityType: 1,
			},
			Config: &capabilitiespb.CapabilityConfig{},
		}

		csOut, err := changeset.AddCapabilities(te.Env, &changeset.AddCapabilitiesRequest{
			RegistryChainSel: te.RegistrySelector,
			DonCapabilities:  map[string][]changeset.DONCapabilityWithConfig{"anything": {cfg}},
		})
		require.NoError(t, err)
		require.Empty(t, csOut.Proposals)
		require.Nil(t, csOut.AddressBook)
		reg := te.CapabilitiesRegistry()
		wantID, err := reg.GetHashedCapabilityId(nil, "test-cap", "0.0.1")
		require.NoError(t, err)
		info, err := reg.GetCapability(nil, wantID)
		require.NoError(t, err)
		assert.Equal(t, uint8(1), info.CapabilityType)
	})

	t.Run("with mcms", func(t *testing.T) {
		te := test.SetupTestEnv(t, test.TestConfig{
			WFDonConfig:     test.DonConfig{N: 4},
			AssetDonConfig:  test.DonConfig{N: 4},
			WriterDonConfig: test.DonConfig{N: 4},
			NumChains:       1,
			UseMCMS:         true,
		})

		cfg := changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "test-cap",
				Version:        "0.0.1",
				CapabilityType: 1,
			},
			Config: &capabilitiespb.CapabilityConfig{},
		}

		req := &changeset.AddCapabilitiesRequest{
			RegistryChainSel: te.RegistrySelector,
			DonCapabilities:  map[string][]changeset.DONCapabilityWithConfig{"anything": {cfg}},
			MCMSConfig:       &changeset.MCMSConfig{MinDuration: 0},
		}
		csOut, err := changeset.AddCapabilities(te.Env, req)
		require.NoError(t, err)
		require.Len(t, csOut.Proposals, 1)
		require.Nil(t, csOut.AddressBook)

		// now apply the changeset such that the proposal is signed and execed
		contracts := te.ContractSets()[te.RegistrySelector]
		timelockContracts := map[uint64]*proposalutils.TimelockExecutionContracts{
			te.RegistrySelector: {
				Timelock:  contracts.Timelock,
				CallProxy: contracts.CallProxy,
			},
		}
		_, err = commonchangeset.ApplyChangesets(t, te.Env, timelockContracts, []commonchangeset.ChangesetApplication{
			{
				Changeset: commonchangeset.WrapChangeSet(changeset.AddCapabilities),
				Config:    req,
			},
		})
		require.NoError(t, err)

		reg := te.CapabilitiesRegistry()
		wantID, err := reg.GetHashedCapabilityId(nil, "test-cap", "0.0.1")
		require.NoError(t, err)
		info, err := reg.GetCapability(nil, wantID)
		require.NoError(t, err)
		assert.Equal(t, uint8(1), info.CapabilityType)

	})
}

// validateUpdate checks reads nodes from the registry and checks they have the expected updates
func validateAddCapability(t *testing.T, te test.TestEnv, expected map[p2pkey.PeerID]changeset.NodeUpdate) {
	registry := te.ContractSets()[te.RegistrySelector].CapabilitiesRegistry
	wfP2PIDs := p2pIDs(t, maps.Keys(te.WFNodes))
	nodes, err := registry.GetNodesByP2PIds(nil, wfP2PIDs)
	require.NoError(t, err)
	require.Len(t, nodes, len(wfP2PIDs))
	for _, node := range nodes {
		// only check the fields that were updated
		assert.Equal(t, expected[node.P2pId].EncryptionPublicKey, hex.EncodeToString(node.EncryptionPublicKey[:]))
		assert.Equal(t, expected[node.P2pId].Signer, node.Signer)
	}
}
