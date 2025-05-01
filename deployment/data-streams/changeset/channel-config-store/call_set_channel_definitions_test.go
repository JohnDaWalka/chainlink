package channel_config_store

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func TestCallSetChannelDefinitions(t *testing.T) {
	t.Parallel()

	e := testutil.NewMemoryEnv(t, false, 0)

	e, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			DeployChannelConfigStoreChangeset,
			DeployChannelConfigStoreConfig{
				ChainsToDeploy: []uint64{testutil.TestChain.Selector},
			},
		),
	)
	require.NoError(t, err)

	envDatastore, err := datastore.FromDefault[metadata.SerializedContractMetadata, datastore.DefaultMetadata](e.DataStore)
	require.NoError(t, err)

	record, err := envDatastore.Addresses().Get(
		datastore.NewAddressRefKey(testutil.TestChain.Selector, datastore.ContractType(types.ChannelConfigStore), &deployment.Version1_0_0, ""),
	)
	require.NoError(t, err)
	channelConfigStoreAddr := common.HexToAddress(record.Address)

	// Call the contract.
	callConf := SetChannelDefinitionsConfig{
		DefinitionsByChain: map[uint64]map[string]ChannelDefinition{
			testutil.TestChain.Selector: {
				channelConfigStoreAddr.String(): {
					ChannelConfigStore: channelConfigStoreAddr,
					DonID:              1,
					S3URL:              "https://s3.us-west-2.amazonaws.com/data-streams-channel-definitions.stage.cldev.sh/channel-definitions-staging-mainnet-5ce78acee5113c55f795984cccdaeb7b805653a1c1e2f9d0d1e3279a302f7966.json",
					Hash:               hexToByte32("5ce78acee5113c55f795984cccdaeb7b805653a1c1e2f9d0d1e3279a302f7966"),
				},
			},
		},
		MCMSConfig: nil,
	}
	_, err = CallSetChannelDefinitions(e, callConf)
	require.NoError(t, err)
}

func hexToByte32(s string) [32]byte {
	var b [32]byte
	copy(b[:], common.HexToAddress(s).Bytes())
	return b
}
