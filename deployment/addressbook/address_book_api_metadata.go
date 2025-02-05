package addressbook

type MetadataStore interface {
	Store[MetadataKey, MetadataRecord, MetadataFilterQuery]
	Fetch() ([]MetadataRecord, error) // from Store->Fetcher
	By() MetadataFilterQuery

	// MappedTo returns a map of AddressKey->MetadataRecord, as a convenience, to allow for bulk fetching
	// of metadata.
	Associate(...AddressRecord) (map[AddressKey]MetadataRecord, error)
}

type MutableMetadataStore interface {
	MetadataStore
	MutableStore[MetadataKey, MetadataRecord, MetadataFilterQuery]

	Add(record MetadataRecord) error         // from MutableStore
	AddOrUpdate(record MetadataRecord) error // from MutableStore
	Update(record MetadataRecord) error      // from MutableStore
	Delete(record MetadataRecord) error      // from MutableStore
}

type MetadataFilterQuery interface {
	QueryFilter[MetadataKey, MetadataRecord]
	Fetch() ([]MetadataRecord, error)         // from QueryFilter->Fetcher
	ById(MetadataKey) (MetadataRecord, error) // from QueryFilter->Fetcher
	Chain(uint64) AddressFilterQuery
	Address(string) AddressFilterQuery
}

var _ Record[MetadataKey, MetadataRecord] = MetadataRecord{} // Just to make sure.

type MetadataRecord struct {
	Chain    uint64 // composite primary key element
	Address  string // composite primary key element
	Metadata string // json probably
}

func (rec MetadataRecord) Clone() MetadataRecord {
	return MetadataRecord{Chain: rec.Chain, Address: rec.Address, Metadata: rec.Metadata}
}

func (rec MetadataRecord) Key() MetadataKey {
	return metadataKeyImpl{
		chain:   rec.Chain,
		address: rec.Address,
	}
}

// MetadataKey

// MetadataKey represents a key made up of the Chain and Address, which can be used to index into metadata
// about that specific address.
type MetadataKey interface {
	Chain() uint64   // Chain selector
	Address() string // Address value itself.
}

type metadataKeyImpl struct {
	chain   uint64 // Chain selector
	address string // Address value itself.
}

func (a metadataKeyImpl) Chain() uint64   { return a.chain }
func (a metadataKeyImpl) Address() string { return a.address }

func NewMetadataKey(chain uint64, address string) MetadataKey {
	return metadataKeyImpl{
		chain:   chain,
		address: address,
	}
}
