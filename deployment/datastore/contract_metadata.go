package datastore

import "errors"

var ErrContractMetadataNotFound = errors.New("contract metadata record not found")
var ErrContractMetadataExists = errors.New("contract metadata record already exists")

var _ Record[ContractMetadataKey, ContractMetadata[DefaultMetadata]] = ContractMetadata[DefaultMetadata]{}

// DefaultMetadata is a default implementation of the custom contract metadata
// that is used to store the contract metadata in the datastore.
type DefaultMetadata string

func (d DefaultMetadata) Clone() DefaultMetadata { return d }

type ContractMetadata[M Cloneable[M]] struct {
	ChainSelector uint64 `json:"chain_selector"`
	Address       string `json:"address"`
	Metadata      M      `json:"metadata"`
}

func (r ContractMetadata[M]) Clone() ContractMetadata[M] {
	return ContractMetadata[M]{
		ChainSelector: r.ChainSelector,
		Address:       r.Address,
		Metadata:      r.Metadata.Clone(),
	}
}

func (r ContractMetadata[M]) Key() ContractMetadataKey {
	return NewContractMetadataKey(r.ChainSelector, r.Address)
}
