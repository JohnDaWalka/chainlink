package capabilities_registry

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
)

type capabilitiesRegistry struct {
	Address                       common.Address
	Contract                      *CapabilitiesRegistry
	sc                            *seth.Client
	ExistingHashedCapabilitiesIDs [][32]byte
}

func Deploy(sc *seth.Client) (*capabilitiesRegistry, error) {
	capabilitiesRegistryAddress, tx, capabilitiesRegistryContract, err := DeployCapabilitiesRegistry(
		sc.NewTXOpts(),
		sc.Client,
	)
	if err != nil {
		return nil, err
	}

	_, err = bind.WaitMined(context.Background(), sc.Client, tx)
	if err != nil {
		return nil, err
	}

	fmt.Printf("🚀 Deployed \033[1mcapabilities_registry\033[0m contract at \033[1m%s\033[0m\n", capabilitiesRegistryAddress)
	return &capabilitiesRegistry{
		sc:       sc,
		Address:  capabilitiesRegistryAddress,
		Contract: capabilitiesRegistryContract,
	}, nil
}

func (cr *capabilitiesRegistry) AddCapabilities(capabilities []CapabilitiesRegistryCapability) error {
	tx, err := cr.Contract.AddCapabilities(
		cr.sc.NewTXOpts(),
		capabilities,
	)
	if err != nil {
		return err
	}

	_, err = bind.WaitMined(context.Background(), cr.sc.Client, tx)
	if err != nil {
		return err
	}

	for _, capability := range capabilities {
		hashedCapabilityID, err := cr.Contract.GetHashedCapabilityId(
			cr.sc.NewCallOpts(),
			capability.LabelledName,
			capability.Version,
		)
		if err != nil {
			return err
		}
		cr.ExistingHashedCapabilitiesIDs = append(cr.ExistingHashedCapabilitiesIDs, hashedCapabilityID)
	}

	return nil
}
