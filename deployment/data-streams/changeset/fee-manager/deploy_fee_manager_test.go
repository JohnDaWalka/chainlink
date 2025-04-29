package fee_manager

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func TestDeployFeeManager(t *testing.T) {
	testEnv := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{ShouldDeployMCMS: true})

	// Need the Link Token
	e, err := commonChangesets.Apply(t, testEnv.Environment, nil,
		commonChangesets.Configure(
			deployment.CreateLegacyChangeSet(commonChangesets.DeployLinkToken),
			[]uint64{testutil.TestChain.Selector},
		),
	)
	require.NoError(t, err)

	addresses, err := e.ExistingAddresses.AddressesForChain(testutil.TestChain.Selector)
	require.NoError(t, err)

	chain := e.Chains[testutil.TestChain.Selector]
	linkState, err := commonstate.MaybeLoadLinkTokenChainState(chain, addresses)
	require.NoError(t, err)

	nativeAddr := common.HexToAddress("0x3e5e9111ae8eb78fe1cc3bb8915d5d461f3ef9a9")
	verifierProxyAddr := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e")
	rewardManagerAddr := common.HexToAddress("0x0fd8b81e3d1143ec7f1ce474827ab93c43523ea2")

	cc := DeployFeeManagerConfig{
		ChainsToDeploy: map[uint64]DeployFeeManager{testutil.TestChain.Selector: {
			LinkTokenAddress:     linkState.LinkToken.Address(),
			NativeTokenAddress:   nativeAddr,
			VerifierProxyAddress: verifierProxyAddr,
			RewardManagerAddress: rewardManagerAddr,
		}},
		Ownership: types.OwnershipSettings{
			ShouldTransfer: true,
			MCMSProposalConfig: &proposalutils.TimelockConfig{
				MinDelay: 0,
			},
		},
	}

	resp, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(DeployFeeManagerChangeset, cc),
	)

	require.NoError(t, err)

	envDatastore, err := datastore.FromDefault[
		metadata.SerializedContractMetadata,
		datastore.DefaultMetadata,
	](resp.DataStore)
	require.NoError(t, err)

	// Verify Contract Is Deployed
	record, err := envDatastore.Addresses().Get(
		datastore.NewAddressRefKey(testutil.TestChain.Selector, datastore.ContractType(types.FeeManager), &deployment.Version0_5_0, ""),
	)
	require.NoError(t, err)
	require.NotNil(t, record)

	// Store address for other tests
	t.Run("VerifyMetadata", func(t *testing.T) {
		cm, err := envDatastore.ContractMetadata().Get(
			datastore.NewContractMetadataKey(record.ChainSelector, record.Address),
		)
		require.NoError(t, err)

		md, err := cm.Metadata.ToFeeManagerMetadata()
		require.NoError(t, err)

		expectedMetadata := metadata.FeeManagerMetadata{
			RewardManagerAddress: rewardManagerAddr.String(),
			VerifierProxyAddress: verifierProxyAddr.String(),
			FeeTokens: []metadata.FeeToken{
				{
					TokenType: metadata.Link,
					Address:   linkState.LinkToken.Address().String(),
				},
				{
					TokenType: metadata.Native,
					Address:   nativeAddr.String(),
				},
			},
		}

		assert.Equal(t, expectedMetadata, md)
	})

	t.Run("VerifyOwnershipTransferred", func(t *testing.T) {
		chain := e.Chains[testutil.TestChain.Selector]
		owner, _, err := commonChangesets.LoadOwnableContract(common.HexToAddress(record.Address), chain.Client)
		require.NoError(t, err)
		assert.Equal(t, testEnv.Timelocks[testutil.TestChain.Selector].Timelock.Address(), owner)
	})
}
