package deployment

import (
	"errors"

	"github.com/Masterminds/semver/v3"
)

var (
	ErrAddressNotFound = errors.New("no such address with these qualifiers")
)

// AddressBookV2 represents a read-only view over an address book, and provides a simple Read() function, which then
// allows for various means of access. The full list of address records can be obtained, a filtered list can be
// obtained, or a single address (if it is present) can be obtained by way of a compound key.
type AddressBookV2 interface {
	Addresses() AddressResult
	Metadata() MetadataResult
}

type MutableAddressBookV2 interface {
	AddressBookV2

	// Add inserts a new record into the address book, throwing an error if a record with the same composite primary
	// key is already present.
	Add(record AddressRecord) error

	// AddOrUpdate behaves like Add where there is not already a record with the same composite primary key as the
	// supplied record, otherwise it behaves like an update.
	AddOrUpdate(record AddressRecord) error

	// Update edits an existing record whose fields match the primary key elements of the supplied AddressRecord, with
	// the non-primary-key values of the supplied AddressRecord.
	Update(record AddressRecord) error

	// Delete deletes record whose primary key elements match the supplied AddressRecord, returning an error if no
	// such record exists to be deleted
	Delete(record AddressRecord) error

	AddOrUpdateMetadata(chain AddressOnChain, metadata string)
}

type AddressKey interface {
	Chain() uint64
	Type() ContractType
	Version() semver.Version
	Qualifier() string // optional (by way of an empty-string)
}

// AddressOnChain represents a key made up of the Chain and Address, which can be used to index into metadata
// about that specific address.
type AddressOnChain interface {
	Chain() uint64   // Chain selector
	Address() string // Address value itself.
}

type Fetcher[C any] interface {
	Fetch() ([]C, error)
}

type AddressResult interface {
	Fetcher[AddressRecord]
	FilteredBy() AddressFilterQuery
	ById(AddressKey) (AddressRecord, error)
}

type AddressFilterQuery interface {
	Fetcher[AddressRecord]
	Chain(uint64) AddressFilterQuery
	Type(ContractType) AddressFilterQuery
	Version(version semver.Version) AddressFilterQuery
	QualifierEquals(string) AddressFilterQuery
	QualifierMatches(string) AddressFilterQuery
}

type MetadataResult interface {
	Fetcher[MetadataRecord]
	FilteredBy() MetadataFilterQuery
	ById(AddressOnChain) (MetadataRecord, error)
	MappedTo(...AddressRecord) (map[AddressKey]MetadataRecord, error)
}

type MetadataFilterQuery interface {
	Fetcher[MetadataRecord]
	Chain(uint64) AddressFilterQuery
	Address(string) AddressFilterQuery
}

func foo(ab AddressBookV2) {
	addresses, _ := ab.Addresses().Fetch()
	for _, address := range addresses {
		println(address.Labels)
	}

	addresses, err := ab.Addresses().FilteredBy().Chain(1).QualifierEquals("foo").Fetch()
	if err != nil {
		// do stuff
	}
	for _, address := range addresses {
		println(address.Labels)
	}

	key := NewAddressKey(1, "blah", Version1_0_0, "")
	rec, err := ab.Addresses().ById(key)
	if errors.Is(err, ErrAddressNotFound) {
		println("Address could not be found for %v", key)
	} else if err != nil {
		// panic
	}
	println("Found record %v", rec)

	metadata := rec.Metadata()
	if metadata != nil {
		// do stuff with metadata
	}

	ab.SetMetadata(rec.key(), metadata)

	metadatas, _ := ab.Metadata().FilteredBy().Chain(3).Fetch()

	metadata, _ := ab.Metadata().ById(NewAddressOnChain(3, "0xABCD0123"))

	metadataMap, _ := ab.Metadata().MappedTo(rec)
	metadata = metadataMap[rec.Key()]

	key := NewAddressKey(1, "blah", Version1_0_0, "")
	rec, err := ab.Addresses().ById(key)
	if errors.Is(err, ErrAddressNotFound) {
		println("Address could not be found for %v", key)
	} else if err != nil {
		// panic
	}
	println("Found record %v", rec)

}

type AddressRecord struct {
	Chain     uint64         // composite primary key element
	Type      ContractType   // composite primary key segment
	Version   semver.Version // composite primary key element
	Qualifier string         // composite primary key element
	Address   string
	Labels    LabelSet
}

func (rec AddressRecord) Key() AddressKey {
	return addressKeyImpl{
		chain:        rec.Chain,
		contractType: rec.Type,
		version:      rec.Version,
		qualifier:    rec.Qualifier,
	}
}

type MetadataRecord struct {
	Chain    uint64 // composite primary key element
	Address  string // composite primary key element
	Metadata string // json probably
}

func (rec MetadataRecord) Key() AddressOnChain {
	return addressOnChainImpl{
		chain:   rec.Chain,
		address: rec.Address,
	}
}

type addressOnChainImpl struct {
	chain   uint64 // Chain selector
	address string // Address value itself.
}

func (a addressOnChainImpl) Chain() uint64   { return a.chain }
func (a addressOnChainImpl) Address() string { return a.address }

func NewAddressOnChain(chain uint64, address string) AddressOnChain {
	return addressOnChainImpl{
		chain:   chain,
		address: address,
	}
}

type addressKeyImpl struct {
	chain        uint64         // composite primary key element
	contractType ContractType   // composite primary key segment
	version      semver.Version // composite primary key element
	qualifier    string         // composite primary key element
}

func (a addressKeyImpl) Chain() uint64           { return a.chain }
func (a addressKeyImpl) Type() ContractType      { return a.contractType }
func (a addressKeyImpl) Version() semver.Version { return a.version }
func (a addressKeyImpl) Qualifier() string       { return a.qualifier }

func NewAddressKey(chain uint64, contractType ContractType, version semver.Version, qualifier string) AddressKey {
	return addressKeyImpl{
		chain:        chain,
		contractType: contractType,
		version:      version,
		qualifier:    qualifier,
	}
}
