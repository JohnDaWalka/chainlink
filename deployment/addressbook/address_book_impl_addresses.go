package addressbook

import (
	"slices"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink/deployment"
)

var _ AddressStore = &simpleInMemoryAddressStore{}
var _ MutableAddressStore = &simpleInMemoryAddressStore{}

// TODO: Addin a PK index.
// TODO consider concurrent access constraints / locking semantics

type simpleInMemoryAddressStore struct {
	data []AddressRecord
}

// NewSimpleInMemoryAddressStore Creates an in-memory implementation of MutableAddressStore
func NewSimpleInMemoryAddressStore() MutableAddressStore {
	return &simpleInMemoryAddressStore{data: []AddressRecord{}}
}

// indexOf does a look up and returns the index of the record that has that key.
func (s *simpleInMemoryAddressStore) indexOf(key AddressKey) int {
	for index, record := range s.data {
		if record.Key() == key {
			return index
		}
	}
	return -1
}

func (s *simpleInMemoryAddressStore) Add(record AddressRecord) error {
	index := s.indexOf(record.Key())
	switch {
	case index >= 0:
		return ErrRecordExists
	default:
		s.data = append(s.data, record)
		return nil
	}
}

func (s *simpleInMemoryAddressStore) AddOrUpdate(record AddressRecord) error {
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

func (s *simpleInMemoryAddressStore) Update(record AddressRecord) error {
	index := s.indexOf(record.Key())
	switch {
	case index >= 0:
		slices.Replace(s.data, index, index+1, record)
		return nil
	default:
		return ErrRecordNotFound
	}
}

func (s *simpleInMemoryAddressStore) Delete(record AddressRecord) error {
	//TODO implement me
	panic("implement me")
}

func (s *simpleInMemoryAddressStore) Fetch() ([]AddressRecord, error) {
	var records []AddressRecord
	for _, r := range s.data {
		records = append(records, r.Clone())
	}
	return records, nil
}

func (s *simpleInMemoryAddressStore) By() AddressFilterQuery {
	return &simpleInMemoryFilterQuery{store: s}
}

type simpleInMemoryFilterQuery struct {
	store *simpleInMemoryAddressStore
}

func (s *simpleInMemoryFilterQuery) Fetch() ([]AddressRecord, error) {
	//TODO implement me
	panic("implement me")
}

func (s *simpleInMemoryFilterQuery) Id(key AddressKey) (AddressRecord, error) {
	index := s.store.indexOf(key)
	if index >= 0 {
		return s.store.data[index].Clone(), nil
	}
	return AddressRecord{}, ErrRecordNotFound
}

func (s *simpleInMemoryFilterQuery) Chain(u uint64) AddressFilterQuery {
	//TODO implement me
	panic("implement me")
}

func (s *simpleInMemoryFilterQuery) Type(contractType deployment.ContractType) AddressFilterQuery {
	//TODO implement me
	panic("implement me")
}

func (s *simpleInMemoryFilterQuery) Version(version semver.Version) AddressFilterQuery {
	//TODO implement me
	panic("implement me")
}

func (s *simpleInMemoryFilterQuery) QualifierEquals(s2 string) AddressFilterQuery {
	//TODO implement me
	panic("implement me")
}

func (s *simpleInMemoryFilterQuery) QualifierMatches(s2 string) AddressFilterQuery {
	//TODO implement me
	panic("implement me")
}
