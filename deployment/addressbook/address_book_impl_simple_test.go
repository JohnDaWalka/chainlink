package addressbook

import (
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

//
// Define test implementation classes. Namely the key type, the record type, and the query filter type.
//

type testKey struct {
	id string
}

var _ Record[testKey, testRecord] = testRecord{}

type testRecord struct {
	id    string
	value string
}

func (t testRecord) Equals(other any) bool {
	o, ok := other.(testRecord)
	return ok &&
		t.id == o.id &&
		t.value == o.value
}

func (t testRecord) Clone() testRecord { return testRecord{id: t.id, value: t.value} }

func (t testRecord) Key() testKey { return testKey{id: t.id} }

var _ QueryFilter[testKey, testRecord] = testQueryFilter{}

type testQueryFilter AbstractListBackedQueryFilter[testKey, testRecord, testQueryFilter]

func (t testQueryFilter) Fetch() ([]testRecord, error) {
	return FetchImpl[testKey, testRecord, testQueryFilter](t)
	data := *(*t.store).internalData()
	filtered := make([]testRecord, 0, len(data))
	for _, rec := range data {
		if filterAll(rec, t.filters) {
			filtered = append(filtered, rec)
		}
	}
	return filtered, nil
}

func (t testQueryFilter) Id(key testKey) (testRecord, error) { return (*t.store).fetchById(key) }

func (t testQueryFilter) ValueContains(substr string) testQueryFilter {
	t.filters = append(t.filters, func(element testRecord) bool {
		return strings.Contains(element.value, substr)
	})
	return t
}

func newTestAddresStore(data ...testRecord) ListBackedMemoryStore[testKey, testRecord, testQueryFilter] {
	return internalInMemoryStoreWithData[testKey, testRecord, testQueryFilter](
		data,
		func(s *ListBackedMemoryStore[testKey, testRecord, testQueryFilter]) testQueryFilter {
			return testQueryFilter{store: s}
		},
		func() testRecord { return testRecord{} },
	)
}

func Test_Store_indexOf(t *testing.T) {
	rec := testRecord{
		id:    "A",
		value: "blah",
	}
	store := newTestAddresStore(rec)
	index := store.indexOf(rec.Key())
	require.Equal(t, 0, index, "Expected record at position 0")
	require.Equal(t, rec, (*store.internalData())[index], "Should be the same object.")
}
func Test_Store_Add(t *testing.T) {
	store := newTestAddresStore()
	rec := testRecord{
		id:    "A",
		value: "blah",
	}
	err := store.Add(rec)
	require.NoError(t, err)
	index := store.indexOf(rec.Key())
	require.Equal(t, 0, index, "Expected record at position 0")
	require.Equal(t, rec, (*store.internalData())[index], "Should be the same object.")
}

func Test_Store_Add_DupeFail(t *testing.T) {
	store := newTestAddresStore()
	rec1 := testRecord{
		id:    "A",
		value: "blah",
	}
	rec2 := testRecord{
		id:    "A",
		value: "foo",
	}
	err := store.Add(rec1)
	require.NoError(t, err)

	// Add in a new record, but with the same key.
	err = store.Add(rec2)
	require.Error(t, ErrRecordExists)
}

func Test_Store_Update(t *testing.T) {
	rec := testRecord{
		id:    "A",
		value: "blah",
	}
	store := newTestAddresStore(rec)

	// just check that this is as expected.
	require.Equal(t, rec, (*store.internalData())[0], "Should be the same object.")

	rec2 := testRecord{
		id:    "A",
		value: "foo",
	}
	err := store.Update(rec2)
	require.NoError(t, err)
	require.Equal(t, "foo", (*store.internalData())[0].value, "Should be the same value.")
}

func Test_Store_Update_NoRecord(t *testing.T) {
	store := newTestAddresStore()
	require.Equal(t, 0, len(*store.internalData()), "Expected no data in store")

	rec := testRecord{
		id:    "A",
		value: "blah",
	}
	err := store.Update(rec)
	require.Error(t, err, ErrRecordNotFound)
	require.Equal(t, 0, len(*store.internalData()), "Expected no data in store")
}

func Test_Store_AddOrUpdate(t *testing.T) {
	store := newTestAddresStore()
	rec1 := testRecord{
		id:    "A",
		value: "blah",
	}
	rec2 := testRecord{
		id:    "A",
		value: "foo",
	}
	require.Equal(t, 0, len(*store.internalData()))

	err := store.AddOrUpdate(rec1)
	require.NoError(t, err)
	require.Equal(t, 1, len(*store.internalData()))
	require.Equal(t, "blah", (*store.internalData())[0].value)

	err = store.AddOrUpdate(rec2)
	require.NoError(t, err)
	require.Equal(t, 1, len(*store.internalData()))
	require.Equal(t, "foo", (*store.internalData())[0].value)
}

func Test_Store_Delete(t *testing.T) {
	rec1 := testRecord{
		id:    "A",
		value: "blah",
	}
	rec2 := testRecord{
		id:    "B",
		value: "foo",
	}
	store := newTestAddresStore(rec1, rec2)
	require.Equal(t, 2, len(*store.internalData()))

	err := store.Delete(rec2)
	require.NoError(t, err)
	require.Equal(t, 1, len(*store.internalData()))
	require.Equal(t, "blah", (*store.internalData())[0].value)

	_, err = store.fetchById(testKey{"B"})
	require.Error(t, err)
}

func Test_Store_Delete_FailNotEqual(t *testing.T) {
	rec1 := testRecord{
		id:    "A",
		value: "blah",
	}
	rec2 := testRecord{
		id:    "B",
		value: "foo",
	}
	rec3 := testRecord{
		id:    "B",
		value: "bar",
	}
	store := newTestAddresStore(rec1, rec2)
	require.Equal(t, 2, len(*store.internalData()))

	err := store.Delete(rec3)
	require.ErrorIs(t, err, ErrRecordsNotEqual)
	require.Equal(t, 2, len(*store.internalData()))
}

func Test_Store_Query_ById(t *testing.T) {
	rec1 := testRecord{
		id:    "A",
		value: "blah",
	}
	rec2 := testRecord{
		id:    "B",
		value: "foo",
	}
	rec3 := testRecord{
		id:    "C",
		value: "bar",
	}
	store := newTestAddresStore(rec1, rec2, rec3)
	require.Equal(t, 3, len(*store.internalData()))

	rec, err := store.By().Id(testKey{"B"})
	require.NoError(t, err, ErrRecordsNotEqual)
	require.Equal(t, "foo", rec.value)
}

func Test_Store_Query_ById_NotFound(t *testing.T) {
	rec1 := testRecord{
		id:    "A",
		value: "blah",
	}
	rec2 := testRecord{
		id:    "B",
		value: "foo",
	}
	rec3 := testRecord{
		id:    "C",
		value: "bar",
	}
	store := newTestAddresStore(rec1, rec2, rec3)
	require.Equal(t, 3, len(*store.internalData()))

	_, err := store.By().Id(testKey{"D"})
	require.ErrorIs(t, err, ErrRecordNotFound)
}

func Test_Store_Query_ByCustomFilter(t *testing.T) {
	rec1 := testRecord{
		id:    "A",
		value: "blah",
	}
	rec2 := testRecord{
		id:    "B",
		value: "foo",
	}
	rec3 := testRecord{
		id:    "C",
		value: "bar",
	}
	rec4 := testRecord{
		id:    "D",
		value: "baz",
	}
	store := newTestAddresStore(rec1, rec2, rec3, rec4)
	require.Equal(t, 4, len(*store.internalData()))

	result, err := store.By().ValueContains("ba").Fetch()
	require.NoError(t, err)
	values := transform(result, func(e testRecord) string { return e.value })
	require.True(t, slices.Contains(values, "bar"))
	require.True(t, slices.Contains(values, "baz"))
	require.False(t, slices.Contains(values, "blah"))
	require.False(t, slices.Contains(values, "foo"))
}
