package metadata

import (
	"testing"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/stretchr/testify/require"
)

func TestEnvMetadata(t *testing.T) {
	ds := datastore.NewMemoryDataStore[
		SerializedContractMetadata,
		DataStreamsMetadata,
	]()

	metaData, err := ds.EnvMetadataStore.Get()

	metaData.Metadata.DONs = []DonMetadata{
		{
			ID:                  "don1",
			ConfiguratorAddress: "0x1234567890abcdef",
			OffchainConfig: OffchainConfig{
				DeltaGrace:   "1",
				DeltaInitial: "1",
			},
			Streams: []int{1, 2},
		},
	}

	err = ds.EnvMetadataStore.Set(metaData)
	require.NoError(t, err)

	metadata, err := ds.EnvMetadataStore.Get()
	require.NoError(t, err)

	don, err := metadata.Metadata.GetDonById("don1")
	require.NoError(t, err)
	require.Equal(t, 2, len(don.Streams))

}
