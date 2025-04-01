package datastore

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/require"
)

func TestMemoryDataStore_Merge(t *testing.T) {
	// Create two mutable MemoryDataStores
	dataStore1 := NewMemoryDataStore[DefaultMetadata]()
	dataStore2 := NewMemoryDataStore[DefaultMetadata]()

	// Add some data to the second data store
	err := dataStore2.Addresses().AddOrUpdate(AddressRef{
		Address:   "0x123",
		Type:      "type1",
		Version:   semver.MustParse("1.0.0"),
		Qualifier: "qualifier1",
	})
	require.NoError(t, err, "Adding data to dataStore2 should not fail")

	// Merge dataStore2 into dataStore1
	err = dataStore1.Merge(dataStore2.Seal())
	require.NoError(t, err, "Merging dataStore2 into dataStore1 should not fail")

	// Verify that dataStore1 contains the merged data
	addressRefs, err := dataStore1.Addresses().Fetch()
	require.NoError(t, err, "Fetching addresses from dataStore1 should not fail")
	require.Len(t, addressRefs, 1, "dataStore1 should contain 1 address after merge")
	require.Equal(t, "0x123", addressRefs[0].Address, "Merged address should match")
}
