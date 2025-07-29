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
	Set(ctx context.Context, key string, value []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
}
