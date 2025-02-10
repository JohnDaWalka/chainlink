package addressbook

import (
	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink/deployment"
)

// TODO: Addin a PK index.
// TODO consider concurrent access constraints / locking semantics

// NewSimpleInMemoryAddressStore Creates an in-memory implementation of MutableAddressStore. This is a non-threadsafe
// store, and should not be accessed concurrently.
func NewSimpleInMemoryAddressStore() MutableAddressStore {
	return NewSimpleInMemoryStore[AddressKey, AddressRecord, AddressFilterQuery](
		func(s *Store[AddressKey, AddressRecord, AddressFilterQuery]) AddressFilterQuery {
			return &simpleInMemoryAddressFilterQuery{store: s}
		},
		func() AddressRecord { return AddressRecord{} },
	)
}

type simpleInMemoryAddressFilterQuery struct {
	store *Store[AddressKey, AddressRecord, AddressFilterQuery]
}

func (s *simpleInMemoryAddressFilterQuery) Fetch() ([]AddressRecord, error) {
	//TODO implement me
	panic("implement me")
}

func (s *simpleInMemoryAddressFilterQuery) Id(key AddressKey) (AddressRecord, error) {
	return (*s.store).fetchById(key)
}

func (s *simpleInMemoryAddressFilterQuery) Chain(u uint64) AddressFilterQuery {
	//TODO implement me
	panic("implement me")
}

func (s *simpleInMemoryAddressFilterQuery) Type(contractType deployment.ContractType) AddressFilterQuery {
	//TODO implement me
	panic("implement me")
}

func (s *simpleInMemoryAddressFilterQuery) Version(version semver.Version) AddressFilterQuery {
	//TODO implement me
	panic("implement me")
}

func (s *simpleInMemoryAddressFilterQuery) QualifierEquals(s2 string) AddressFilterQuery {
	//TODO implement me
	panic("implement me")
}

func (s *simpleInMemoryAddressFilterQuery) QualifierMatches(s2 string) AddressFilterQuery {
	//TODO implement me
	panic("implement me")
}
