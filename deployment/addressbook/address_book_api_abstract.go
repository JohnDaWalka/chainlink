package addressbook

// Fetcher provides an abstract Fetch() method which is used to complete a read query from a Store.
type Fetcher[R any] interface {
	// Fetch returns a slice of records, depending on the context of this fetcher. For instance, it could be the
	// entire data set, or it could be a filtered subset.
	//
	// The slice will be a new slice (not a reference to an existing slice) and the records returned should be copies
	// of the records in storage. Altering these slices and records should have no effect on the underlying data.
	Fetch() ([]R, error)
}

type Cloneable[R any] interface {
	// Clone returns a semi-deep copy of the type, calling clone() on any cloneable data, and making a shallow copy
	// of any slice or map. (in theory it should go further, but navigating not-cloneable references to cloneable
	// things is beyond the needs nad scope of this interface.
	Clone() R
}

// PrimaryKeyHolder is a type that can extract a key that is considered a unique identifier for itself.
type PrimaryKeyHolder[K comparable] interface {
	Key() K
}

// Record is a data element stored in the Store, and returned from a Fetcher or QueryFilter
type Record[K comparable, R PrimaryKeyHolder[K]] interface {
	Cloneable[R]
	PrimaryKeyHolder[K]
}

type QueryFilter[K comparable, R Record[K, R]] interface {
	Fetcher[R]

	// Id represents a fetch of a data element via its unique key.
	Id(K) (R, error)
}

type Store[K comparable, R Record[K, R], F QueryFilter[K, R]] interface {
	Fetcher[R]
	By() F
}

type MutableStore[K comparable, R Record[K, R], F QueryFilter[K, R]] interface {
	Store[K, R, F]
	// Add inserts a new record into the MutableStore.
	Add(record R) error

	// AddOrUpdate behaves like Add where there is not already a record with the same composite primary key as the
	// supplied record, otherwise it behaves like an update.
	AddOrUpdate(record R) error

	// Update edits an existing record whose fields match the primary key elements of the supplied AddressRecord, with
	// the non-primary-key values of the supplied AddressRecord.
	Update(record R) error

	// Delete deletes record whose primary key elements match the supplied AddressRecord, returning an error if no
	// such record exists to be deleted
	Delete(record R) error
}
