package contracts

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/datastore"

	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"
	workflow_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper"
)

type KeystoneContract interface {
	Address() common.Address
}

func GetContractV1[T KeystoneContract](store datastore.AddressRefStore, addressRef datastore.AddressRef, chain deployment.Chain) (*T, error) {
	r, err := store.Get(addressRef.Key())
	if err != nil {
		return nil, err
	}

	handlers := map[deployment.ContractType]func(string) (*T, error){
		CapabilitiesRegistry: func(addr string) (*T, error) {
			c, err := capabilities_registry.NewCapabilitiesRegistry(common.HexToAddress(addr), chain.Client)
			if err != nil {
				return nil, fmt.Errorf("failed to create capability registry contract from address %s: %w", addr, err)
			}
			return any(c).(*T), nil
		},
		KeystoneForwarder: func(addr string) (*T, error) {
			c, err := forwarder.NewKeystoneForwarder(common.HexToAddress(addr), chain.Client)
			if err != nil {
				return nil, fmt.Errorf("failed to create forwarder contract from address %s: %w", addr, err)
			}
			return any(c).(*T), nil
		},
		OCR3Capability: func(addr string) (*T, error) {
			c, err := ocr3_capability.NewOCR3Capability(common.HexToAddress(addr), chain.Client)
			if err != nil {
				return nil, fmt.Errorf("failed to create OCR3 capability contract from address %s: %w", addr, err)
			}
			return any(c).(*T), nil
		},
		WorkflowRegistry: func(addr string) (*T, error) {
			c, err := workflow_registry.NewWorkflowRegistry(common.HexToAddress(addr), chain.Client)
			if err != nil {
				return nil, fmt.Errorf("failed to create workflow registry contract from address %s: %w", addr, err)
			}
			return any(c).(*T), nil
		},
	}

	var contractType deployment.ContractType
	// Determine contract type based on T
	switch any(*new(T)).(type) {
	case *forwarder.KeystoneForwarder:
		contractType = KeystoneForwarder
	case *capabilities_registry.CapabilitiesRegistry:
		contractType = CapabilitiesRegistry
	case *ocr3_capability.OCR3Capability:
		contractType = OCR3Capability
	case *workflow_registry.WorkflowRegistry:
		contractType = WorkflowRegistry
	default:
		return nil, fmt.Errorf("unsupported contract type %T", *new(T))
	}

	handler, exists := handlers[contractType]
	if !exists {
		return nil, fmt.Errorf("unsupported contract type: %s", contractType)
	}

	return handler(r.Address)
}

func GetContractV2[T Ownable](store datastore.AddressRefStore, addressRef datastore.AddressRef, chain deployment.Chain) (*OwnedContract[T], error) {
	r, err := store.Get(addressRef.Key())
	if err != nil {
		return nil, err
	}

	// TODO: we need to refactor `GetOwnedContract` to use the datastore instead of the address book
	return GetOwnedContract[T](nil, chain, r.Address)
}
