package capabilities_registry

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	cr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
)

type CapabilitiesRegistryInstance struct {
	Address                       common.Address
	Contract                      *cr.CapabilitiesRegistry
	sc                            *seth.Client
	ExistingHashedCapabilitiesIDs [][32]byte
}

func Deploy(sc *seth.Client) (*CapabilitiesRegistryInstance, error) {
	capabilitiesRegistryAddress, tx, capabilitiesRegistryContract, err := cr.DeployCapabilitiesRegistry(
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
	return &CapabilitiesRegistryInstance{
		sc:       sc,
		Address:  capabilitiesRegistryAddress,
		Contract: capabilitiesRegistryContract,
	}, nil
}

func (cr *CapabilitiesRegistryInstance) AddCapabilities(capabilities []cr.CapabilitiesRegistryCapability) error {
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
