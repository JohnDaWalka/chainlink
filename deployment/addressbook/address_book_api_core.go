package addressbook

import (
	"errors"
)

var (
	ErrRecordNotFound  = errors.New("no such record can be found with the supplied key")
	ErrRecordExists    = errors.New("a record with the supplied key already exists")
	ErrRecordsNotEqual = errors.New("compared records were not strictly equal")
)

// addressBookV2 represents an abstract generic base type which can support some variations. In particular, this is
// used to allow for a ReadableAddressBook, and a MutableAddressBook, which have overlapping APIs.
type genericAddressBookV2[A AddressStore, M MetadataStore] interface {
	Addresses() A
	Metadata() M
}

// ReadableAddressBookV2 represents a read-only view over an address book, and provides a simple Read() function, which then
// allows for various means of access. The full list of address records can be obtained, a filtered list can be
// obtained, or a single address (if it is present) can be obtained by way of a compound key.
type AddressBook interface {
	genericAddressBookV2[AddressStore, MetadataStore]
	Addresses() AddressStore
	Metadata() MetadataStore
}

type MutableAddressBook interface {
	genericAddressBookV2[MutableAddressStore, MutableMetadataStore]
	Addresses() MutableAddressStore
	Metadata() MutableMetadataStore
}
