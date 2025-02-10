package addressbook

import (
	"slices"
)

// TODO: Addin a PK index.
// TODO consider concurrent access constraints / locking semantics

// ListBackedMemoryStore is a MutableStore which is implemented by an in-memory list (slice) structure.
// The interface provides convenience functions important for implementation of concrete stores and for testing.
type ListBackedMemoryStore[K comparable, R Record[K, R], F QueryFilter[K, R]] interface {
	MutableStore[K, R, F]

	// indexOf supplies the index
	indexOf(key K) int

	internalData() *[]R
}

type simpleInMemoryStore[K comparable, R Record[K, R], F QueryFilter[K, R]] struct {
	data                []R
	filterProvider      func(*ListBackedMemoryStore[K, R, F]) F
	emptyRecordProvider func() R
}

// NewSimpleInMemoryStore Creates an in-memory implementation of MutableAddressStore. This is a non-threadsafe
// store, and should not be accessed concurrently.
func NewSimpleInMemoryStore[K comparable, R Record[K, R], F QueryFilter[K, R]](
	filterProvider func(*ListBackedMemoryStore[K, R, F]) F,
	emptyRecordProvider func() R,
) MutableStore[K, R, F] {
	return internalInMemoryStoreWithData([]R{}, filterProvider, emptyRecordProvider)
}

// NewSimpleInMemoryStoreWithData Creates an in-memory implementation of MutableAddressStore, priming the store with
// the supplied slice of records. This is a non-threadsafe store, and should not be accessed concurrently.
func internalInMemoryStoreWithData[K comparable, R Record[K, R], F QueryFilter[K, R]](
	data []R,
	filterProvider func(*ListBackedMemoryStore[K, R, F]) F,
	emptyRecordProvider func() R,
) ListBackedMemoryStore[K, R, F] {
	return &simpleInMemoryStore[K, R, F]{
		data:                data,
		filterProvider:      filterProvider,
		emptyRecordProvider: emptyRecordProvider,
	}
}

// indexOf does a look up and returns the index of the record that has that key.
func (s *simpleInMemoryStore[K, R, F]) indexOf(key K) int {
	for index, record := range s.data {
		if record.Key() == key {
			return index
		}
	}
	return -1
}

// internalData returns a reference to the underlying slice. Used in testing.
func (s *simpleInMemoryStore[K, R, F]) internalData() *[]R { return &s.data }

// fetchById returns a single record, or ErrRecordNotFound if the record cannot be found with that key.
func (s *simpleInMemoryStore[K, R, F]) fetchById(key K) (R, error) {
	index := s.indexOf(key)
	if index >= 0 {
		// TODO: This is gross. There has to be a cleaner way to do this without double-de-referencing.
		return s.data[index].Clone(), nil
	}
	return s.emptyRecordProvider(), ErrRecordNotFound
}

func (s *simpleInMemoryStore[K, R, F]) Add(record R) error {
	index := s.indexOf(record.Key())
	switch {
	case index >= 0:
		return ErrRecordExists
	default:
		s.data = append(s.data, record)
		return nil
	}
}

func (s *simpleInMemoryStore[K, R, F]) AddOrUpdate(record R) error {
	index := s.indexOf(record.Key())
	switch {
	case index >= 0:
		s.data[index] = record
		return nil
	default:
		s.data = append(s.data, record)
		return nil
	}
}

func (s *simpleInMemoryStore[K, R, F]) Update(record R) error {
	index := s.indexOf(record.Key())
	switch {
	case index >= 0:
		slices.Replace(s.data, index, index+1, record)
		return nil
	default:
		return ErrRecordNotFound
	}
}

func (s *simpleInMemoryStore[K, R, F]) Delete(record R) error {
	index := s.indexOf(record.Key())
	if index < 0 {
		return ErrRecordNotFound
	}
	if !record.Equals(s.data[index]) {
		return ErrRecordsNotEqual
	}
	s.data[index] = s.data[len(s.data)-1] // overwrite deleted element with the last one.
	s.data = s.data[:len(s.data)-1]       // snip snip
	return nil
}

func (s *simpleInMemoryStore[K, R, F]) DeleteById(key K) error {
	index := s.indexOf(key)
	if index < 0 {
		return ErrRecordNotFound
	}
	s.data[index] = s.data[len(s.data)-1] // overwrite deleted element with the last one.
	s.data = s.data[:len(s.data)-1]       // snip snip
	return nil
}

func (s *simpleInMemoryStore[K, R, F]) Fetch() ([]R, error) {
	var records []R
	for _, r := range s.data {
		records = append(records, r.Clone())
	}
	return records, nil
}

func (s *simpleInMemoryStore[K, R, F]) By() F {
	var internalStore ListBackedMemoryStore[K, R, F] = s // Hint to the compiler by explicitly matching to the interface
	return s.filterProvider(&internalStore)
}

// AbstractListBackedQueryFilter is a struct which contains the basic infrastructure for doing query filters on the
// underlying list data
type AbstractListBackedQueryFilter[K comparable, R Record[K, R], F AbstractListBackedQueryFilter[K, R, F]] struct {
	store   *ListBackedMemoryStore[K, R, F]
	filters []func(element R) bool
}

func (f AbstractListBackedQueryFilter[K, R, F]) internalData() *[]R {
	return (*f.store).internalData()
}

func FetchImpl[K comparable, R Record[K, R], F AbstractListBackedQueryFilter[K, R, F]](f F) ([]R, error) {
	(*f).
	filtered := make([]R, 0, len(data))
	for _, rec := range data {
		if filterAll(rec, t.filters) {
			filtered = append(filtered, rec)
		}
	}
	return filtered, nil
}

func (f AbstractListBackedQueryFilter[K, R, F]) Id(key K) (R, error) {
	return (*f.store).fetchById(key)
}
