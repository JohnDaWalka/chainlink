package metadata

import (
	"testing"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/stretchr/testify/require"
)

func TestContractMetadata(t *testing.T) {
	ds := datastore.NewMemoryDataStore[
		SerializedContractMetadata,
		datastore.DefaultMetadata,
	]()

	verifierMetadata := VerifierMetadata{
		Active:               true,
		VerifierProxyAddress: "0x123567",
	}

	// Serialize it
	serialized, err := NewVerifierMetadata(verifierMetadata)
	require.NoError(t, err)

	// Create a ContractMetadata with the serialized data
	contractMetadata := datastore.ContractMetadata[SerializedContractMetadata]{
		Address:       "0x0001",
		ChainSelector: 1,
		Metadata:      serialized,
	}

	err = ds.ContractMetadata().Upsert(contractMetadata)
	if err != nil {
		// handle error
	}

	// Later, retrieve and deserialize
	retrievedSerialized := contractMetadata.Metadata
	//retrievedVerifier, err := retrievedSerialized.ToVerifierMetadata()
	//if err != nil {
	//	// handle error
	//}

	//vpm := verification.VerifierProxyMetadata{
	//	AccessControllerAddress: "0x123",
	//	FeeManagerAddress:       "0x456",
	//	Verifiers:               []string{"0x789", "0xabc"},
	//}
	//
	//ds.ContractMetadataStore.Upsert(vpm)
	//
	//previousMetadata, err := ds.EnvMetadataStore.Get()
	//
	//previousMetadata.Metadata.Data = "streams_metadata"
	////previousMetadata.Metadata.DONs = []DonMetadata{
	////	{
	////		ID:                  "don1",
	////		ConfiguratorAddress: "0x1234567890abcdef",
	////		OffchainConfig: OffchainConfig{
	////			DeltaGrace:   "0x1234567890abcdef",
	////			DeltaInitial: "0x1234567890abcdef",
	////		},
	////		Streams: []int{1, 2},
	////	},
	////}
	//
	//err = ds.EnvMetadataStore.Set(previousMetadata)
	//require.NoError(t, err)
	//
	//metadata, err := ds.EnvMetadataStore.Get()
	//require.NoError(t, err)
	//
	//don, err := metadata.Metadata.GetDonById("don1")
	//require.NoError(t, err)
	//
	//fmt.Println("DON:", don)
	// End test
}

//func TestEnvMetadata(t *testing.T) {
//	//ds := datastore.NewMemoryDataStore[
//	//	SerializedContractMetadata,
//	//	datastore.DefaultMetadata,
//	//]()
//	//
//
//	//previousMetadata, err := ds.EnvMetadataStore.Get()
//	//
//	//previousMetadata.Metadata.Data = "streams_metadata"
//	////previousMetadata.Metadata.DONs = []DonMetadata{
//	////	{
//	////		ID:                  "don1",
//	////		ConfiguratorAddress: "0x1234567890abcdef",
//	////		OffchainConfig: OffchainConfig{
//	////			DeltaGrace:   "0x1234567890abcdef",
//	////			DeltaInitial: "0x1234567890abcdef",
//	////		},
//	////		Streams: []int{1, 2},
//	////	},
//	////}
//	//
//	//err = ds.EnvMetadataStore.Set(previousMetadata)
//	//require.NoError(t, err)
//	//
//	//metadata, err := ds.EnvMetadataStore.Get()
//	//require.NoError(t, err)
//	//
//	//don, err := metadata.Metadata.GetDonById("don1")
//	//require.NoError(t, err)
//	//
//	//fmt.Println("DON:", don)
//	// End test
//}
