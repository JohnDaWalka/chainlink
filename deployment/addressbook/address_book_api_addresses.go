package addressbook

import (
	"github.com/Masterminds/semver/v3"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink/deployment"
)

// Note: Many method declarations are redundant here, as they are implied by the generics. They are replicated here
// for readability. The compiler requires them to be specified correctly here, but the specification here is the
// more concrete types.

// AddressStore is a Store for AddressRecord instances.
type AddressStore interface {
	Store[AddressKey, AddressRecord, AddressFilterQuery]

	Fetch() ([]AddressRecord, error) // from Store->Fetcher
	By() AddressFilterQuery          // from Store
}

// MutableAddressStore provides write capabilities to the AddressStore, a portion of the overall AddressBook.
type MutableAddressStore interface {
	AddressStore
	MutableStore[AddressKey, AddressRecord, AddressFilterQuery]

	Add(record AddressRecord) error         // from MutableStore
	AddOrUpdate(record AddressRecord) error // from MutableStore
	Update(record AddressRecord) error      // from MutableStore
	Delete(record AddressRecord) error      // from MutableStore
}

// AddressFilterQuery provides a query interface suited to an AddressRecord and its AddressKey elements.
type AddressFilterQuery interface {
	QueryFilter[AddressKey, AddressRecord]

	Fetch() ([]AddressRecord, error)      // from QueryFilter->Fetcher
	Id(AddressKey) (AddressRecord, error) // From QueryFilter

	Chain(uint64) AddressFilterQuery
	Type(deployment.ContractType) AddressFilterQuery
	Version(version semver.Version) AddressFilterQuery
	QualifierEquals(string) AddressFilterQuery
	QualifierMatches(string) AddressFilterQuery
}

var _ Record[AddressKey, AddressRecord] = AddressRecord{} // Just to make sure.

type AddressRecord struct {
	Chain     uint64                  // composite primary key element
	Type      deployment.ContractType // composite primary key segment
	Version   semver.Version          // composite primary key element
	Qualifier string                  // composite primary key element
	Address   string
	Labels    deployment.LabelSet
}

func (rec AddressRecord) Clone() AddressRecord {
	return AddressRecord{
		Chain:     rec.Chain,
		Type:      rec.Type,
		Version:   rec.Version,
		Qualifier: rec.Qualifier,
		Address:   rec.Address,
		Labels:    maps.Clone(rec.Labels),
	}
}

func (rec AddressRecord) Key() AddressKey {
	return addressKeyImpl{
		chain:        rec.Chain,
		contractType: rec.Type,
		version:      rec.Version,
		qualifier:    rec.Qualifier,
	}
}

// AddressKey provides a composite primary key for an address book record, which can be used to look up that record,
// or otherwise reference a given address. The key is hierarchically composed of a Chain, a Type (string-alike),
// a Version (semver), and a Qualifier. Only one record can exist with the same data.
type AddressKey interface {
	Chain() uint64
	Type() deployment.ContractType
	Version() semver.Version
	Qualifier() string // optional (by way of an empty-string)
}

type addressKeyImpl struct {
	chain        uint64                  // composite primary key element
	contractType deployment.ContractType // composite primary key segment
	version      semver.Version          // composite primary key element
	qualifier    string                  // composite primary key element
}

func (a addressKeyImpl) Chain() uint64                 { return a.chain }
func (a addressKeyImpl) Type() deployment.ContractType { return a.contractType }
func (a addressKeyImpl) Version() semver.Version       { return a.version }
func (a addressKeyImpl) Qualifier() string             { return a.qualifier }

func NewAddressKey(chain uint64, contractType deployment.ContractType, version semver.Version, qualifier string) AddressKey {
	return addressKeyImpl{
		chain:        chain,
		contractType: contractType,
		version:      version,
		qualifier:    qualifier,
	}
}
