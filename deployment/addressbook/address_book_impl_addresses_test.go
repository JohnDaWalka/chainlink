package addressbook

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
)

func Test_AddressStore_indexOf(t *testing.T) {
	rec := AddressRecord{
		Chain:     0,
		Type:      "a",
		Version:   deployment.Version1_1_0,
		Qualifier: "b",
		Address:   "0x1234abcd",
		Labels:    deployment.LabelSet{},
	}
	store := simpleInMemoryAddressStore{
		data: []AddressRecord{rec},
	}
	index := store.indexOf(rec.Key())
	require.Equal(t, 0, index, "Expected record at position 0")
	require.Equal(t, rec, store.data[index], "Should be the same object.")
}
func Test_AddressStore_Add(t *testing.T) {
	store := simpleInMemoryAddressStore{}
	rec := AddressRecord{
		Chain:     0,
		Type:      "a",
		Version:   deployment.Version1_1_0,
		Qualifier: "b",
		Address:   "0x1234abcd",
		Labels:    deployment.LabelSet{},
	}
	err := store.Add(rec)
	require.NoError(t, err)
	index := store.indexOf(rec.Key())
	require.Equal(t, 0, index, "Expected record at position 0")
	require.Equal(t, rec, store.data[index], "Should be the same object.")
}

func Test_AddressStore_Add_DupeFail(t *testing.T) {
	store := simpleInMemoryAddressStore{}
	rec1 := AddressRecord{
		Chain:     0,
		Type:      "a",
		Version:   deployment.Version1_1_0,
		Qualifier: "q",
		Address:   "0x1234abcd",
		Labels:    deployment.LabelSet{},
	}
	rec2 := AddressRecord{
		Chain:     0,
		Type:      "a",
		Version:   deployment.Version1_1_0,
		Qualifier: "q",
		Address:   "0xdeadbeef",
		Labels:    deployment.LabelSet{},
	}
	err := store.Add(rec1)
	require.NoError(t, err)

	// Add in a new record, but with the same key.
	err = store.Add(rec2)
	require.Error(t, ErrRecordExists)
}

func Test_AddressStore_Update(t *testing.T) {
	rec := AddressRecord{
		Chain:     0,
		Type:      "a",
		Version:   deployment.Version1_1_0,
		Qualifier: "q",
		Address:   "0x1234abcd",
		Labels:    deployment.LabelSet{},
	}
	store := simpleInMemoryAddressStore{
		data: []AddressRecord{rec},
	}
	// just check that this is as expected.
	require.Equal(t, rec, store.data[0], "Should be the same object.")

	rec2 := AddressRecord{
		Chain:     0,
		Type:      "a",
		Version:   deployment.Version1_1_0,
		Qualifier: "q",
		Address:   "0xdeadbeef",
		Labels:    deployment.LabelSet{},
	}
	err := store.Update(rec2)
	require.NoError(t, err)
	require.Equal(t, rec2, store.data[0], "Should be the same object.")
}

func Test_AddressStore_Update_NoRecord(t *testing.T) {
	store := simpleInMemoryAddressStore{}
	require.Equal(t, 0, len(store.data), "Expected no data in store")

	rec := AddressRecord{
		Chain:     0,
		Type:      "a",
		Version:   deployment.Version1_1_0,
		Qualifier: "q",
		Address:   "0xdeadbeef",
		Labels:    deployment.LabelSet{},
	}
	err := store.Update(rec)
	require.Error(t, err, ErrRecordNotFound)
	require.Equal(t, 0, len(store.data), "Expected no data in store")
}

func Test_AddressStore_AddOrUpdate(t *testing.T) {
	store := simpleInMemoryAddressStore{}
	rec1 := AddressRecord{
		Chain:     0,
		Type:      "a",
		Version:   deployment.Version1_1_0,
		Qualifier: "q",
		Address:   "0x1234abcd",
		Labels:    deployment.LabelSet{},
	}
	rec2 := AddressRecord{
		Chain:     0,
		Type:      "a",
		Version:   deployment.Version1_1_0,
		Qualifier: "q",
		Address:   "0xdeadbeef",
		Labels:    deployment.LabelSet{},
	}
	require.Equal(t, 0, len(store.data))

	err := store.AddOrUpdate(rec1)
	require.NoError(t, err)
	require.Equal(t, 1, len(store.data))
	require.Equal(t, "0x1234abcd", store.data[0].Address)

	err = store.AddOrUpdate(rec2)
	require.NoError(t, err)
	require.Equal(t, 1, len(store.data))
	require.Equal(t, "0xdeadbeef", store.data[0].Address)
}
