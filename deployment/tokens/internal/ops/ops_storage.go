package ops

import (
	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment"
)

// OpAddAddrBookRecordDeps defines the dependencies to perform the OpAddAddrBookRecord.
type OpAddAddrBookRecordDeps struct {
	AddrBook deployment.AddressBook
}

// OpAddAddrBookRecordInput is the input to the OpAddAddrBookRecord operation.
type OpAddAddrBookRecordInput struct {
	ChainSelector uint64   `json:"chainSelector"`
	Address       string   `json:"address"`
	Type          string   `json:"type"`
	Version       string   `json:"version"`
	Labels        []string `json:"labels"`
}

// OpAddAddrBookRecordOutput is the output of the OpAddAddrBookRecord operation.
type OpAddAddrBookRecordOutput struct {
	ChainSelector  uint64 `json:"chainSelector"`
	Address        string `json:"address"`
	TypeAndVersion string `json:"typeAndVersion"`
}

// OpAddAddrBookRecord adds a new address record to the address book.
var OpAddAddrBookRecord = operations.NewOperation(
	"add-address-book-record",
	semver.MustParse("1.0.0"),
	"Adds an address record to address book",
	func(b operations.Bundle, deps OpAddAddrBookRecordDeps, in OpAddAddrBookRecordInput) (OpAddAddrBookRecordOutput, error) {
		out := OpAddAddrBookRecordOutput{}

		tv := deployment.NewTypeAndVersion(
			deployment.ContractType(in.Type),
			*semver.MustParse(in.Version),
		)

		for _, label := range in.Labels {
			tv.AddLabel(label)
		}

		if err := deps.AddrBook.Save(in.ChainSelector, in.Address, tv); err != nil {
			return out, err
		}

		return OpAddAddrBookRecordOutput{
			ChainSelector:  in.ChainSelector,
			Address:        in.Address,
			TypeAndVersion: tv.String(),
		}, nil
	})

// OpAddDatastoreAddrRefDeps defines the dependencies to perform the OpAddDatastoreAddrRef
// operation.
type OpAddDatastoreAddrRefDeps struct {
	Datastore datastore.MutableDataStore[
		datastore.DefaultMetadata,
		datastore.DefaultMetadata,
	]
}

// OpAddDatastoreAddrRefInput is the input to the OpAddDatastoreAddrRef operation.
type OpAddDatastoreAddrRefInput struct {
	ChainSelector uint64   `json:"chainSelector"`
	Address       string   `json:"address"`
	Qualifier     string   `json:"qualifier"`
	Type          string   `json:"type"`
	Version       string   `json:"version"`
	Labels        []string `json:"labels"`
}

// OpAddDatastoreAddrRefOutput is the output of the OpAddDatastoreAddrRef operation.
type OpAddDatastoreAddrRefOutput struct {
	ChainSelector uint64 `json:"chainSelector"`
	Address       string `json:"address"`
}

// OpAddDatastoreAddrRef adds a new address reference to the datastore.
var OpAddDatastoreAddrRef = operations.NewOperation(
	"add-datastore-address-reference",
	semver.MustParse("1.0.0"),
	"Adds an address reference to the datastore",
	func(b operations.Bundle, deps OpAddDatastoreAddrRefDeps, input OpAddDatastoreAddrRefInput) (OpAddDatastoreAddrRefOutput, error) {
		out := OpAddDatastoreAddrRefOutput{}

		if err := deps.Datastore.Addresses().Add(
			datastore.AddressRef{
				ChainSelector: input.ChainSelector,
				Address:       input.Address,
				Type:          datastore.ContractType(input.Type),
				Version:       semver.MustParse(input.Version),
				Qualifier:     input.Qualifier,
				Labels:        datastore.NewLabelSet(input.Labels...),
			},
		); err != nil {
			return out, err
		}

		return OpAddDatastoreAddrRefOutput{
			ChainSelector: input.ChainSelector,
			Address:       input.Address,
		}, nil
	})
