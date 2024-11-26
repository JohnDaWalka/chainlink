package forwarder

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
)

type instance struct {
	Address                       common.Address
	Contract                      *KeystoneForwarder
	sc                            *seth.Client
	ExistingHashedCapabilitiesIDs [][32]byte
}

func Deploy(sc *seth.Client) (*instance, error) {
	forwarderAddress, tx, forwarderContract, err := DeployKeystoneForwarder(
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

	fmt.Printf("🚀 Deployed \033[1mforwarder\033[0m contract at \033[1m%s\033[0m\n", forwarderAddress)
	return &instance{
		sc:       sc,
		Address:  forwarderAddress,
		Contract: forwarderContract,
	}, nil
}

func (i *instance) SetConfig(
	donID uint32,
	configVersion uint32,
	f uint8,
	signers []common.Address,
) error {
	tx, err := i.Contract.SetConfig(
		i.sc.NewTXOpts(),
		donID,
		configVersion,
		f,
		signers,
	)
	if err != nil {
		return err
	}

	_, err = bind.WaitMined(context.Background(), i.sc.Client, tx)
	if err != nil {
		return err
	}

	return nil
}
