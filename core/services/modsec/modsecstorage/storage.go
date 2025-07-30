package modsecstorage

import "context"

// Storage abstracts the storage of ccip verifier data.
// The value can be differently formatted per verifier.
// For the commit verifier, the key will be the message ID, emitted on source,
// and the value will be:
//
//	{
//	  "message_data": {... fully encoded message data, including verifiers specified ...},
//	  "proofs": {
//	    "evm": "<ecdsa signature>",
//	    "solana": "<solana signature>",
//	    ... etc. for all supported chain families.
//	  }
//	}
//
// stored as a JSON blob.
type Storage interface {
	// Set sets the value for a given key.
	// If the key already exists, the value is overwritten.
	// Errors are only returned if the storage is unavailable or other similar transport errors.
	Set(ctx context.Context, key string, value []byte) error

	// Get returns the value for a given key.
	// If the key is not found, an error is returned.
	// Errors are also returned if the storage is unavailable or other similar transport errors.
	Get(ctx context.Context, key string) ([]byte, error)

	// GetMany returns a map of found key-value pairs.
	// If a key is not found, it is not included in the map.
	// Errors are also returned if the storage is unavailable or other similar transport errors.
	GetMany(ctx context.Context, keys []string) (map[string][]byte, error)

	// GetAll returns all key-value pairs in storage.
	// TODO: very unlikely this will remain in the final interface, just putting it here
	// to get an E2E test going while we figure out efficient ways to query KV storage
	// w/out source chain events being ingested by the executor.
	// Errors are also returned if the storage is unavailable or other similar transport errors.
	GetAll(ctx context.Context) (map[string][]byte, error)
}
