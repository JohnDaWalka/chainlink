package datastore

// ContractMetadataKey is an interface that represents a key for ContractMetadata records.
// It is used to uniquely identify a record in the ContractMetadataStore.
type ContractMetadataKey interface {
	Comparable[ContractMetadataKey]

	// Address returns the address of the contract on the chain.
	Address() string
	// ChainSelector returns the chain-selector of the chain where the contract is deployed.
	ChainSelector() uint64
}

// contractMetadataKey implements the ContractMetadataKey interface.
var _ ContractMetadataKey = contractMetadataKey{}

// contractMetadataKey is a struct that implements the ContractMetadataKey interface.
// It is used to uniquely identify a record in the ContractMetadataStore.
type contractMetadataKey struct {
	chain   uint64
	address string
}

// ChainSelector returns the chain-selector of the chain where the contract is deployed.
func (c contractMetadataKey) ChainSelector() uint64 { return c.chain }

// Address returns the address of the contract on the chain.
func (c contractMetadataKey) Address() string { return c.address }

// Equals returns true if the two ContractMetadataKey instances are equal, false otherwise.
func (c contractMetadataKey) Equals(other ContractMetadataKey) bool {
	return c.chain == other.ChainSelector() &&
		c.address == other.Address()
}

// NewContractMetadataKey creates a new ContractMetadataKey instance.
func NewContractMetadataKey(chain uint64, address string) ContractMetadataKey {
	return contractMetadataKey{
		chain:   chain,
		address: address,
	}
}
