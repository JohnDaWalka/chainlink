package forwarder

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	forwarder_wrapper "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder"
)

type ForwarderInstance struct {
	Address                       common.Address
	Contract                      *forwarder_wrapper.KeystoneForwarder
	sc                            *seth.Client
	ExistingHashedCapabilitiesIDs [][32]byte
}

func Deploy(sc *seth.Client) (*ForwarderInstance, error) {
	forwarderAddress, tx, forwarderContract, err := forwarder_wrapper.DeployKeystoneForwarder(
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
	return &ForwarderInstance{
		sc:       sc,
		Address:  forwarderAddress,
		Contract: forwarderContract,
	}, nil
}

func (i *ForwarderInstance) SetConfig(
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
