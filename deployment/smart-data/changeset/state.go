package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	"github.com/smartcontractkit/chainlink/deployment/smart-data/changeset/globals"
)

type GetContractSetsRequest struct {
	Chains      map[uint64]cldf_evm.Chain
	AddressBook cldf.AddressBook

	// Labels indicates the label set that a contract must include to be considered as a member
	// of the returned contract set.  By default, an empty label set implies that only contracts without
	// labels will be considered.  Otherwise, all labels must be on the contract (e.g., "label1" AND "label2").
	Labels []string
}

type GetContractSetsResponse struct {
	ContractSets map[uint64]ContractSet
}

// ContractSet is a set of contracts for a single chain
// It is a mirror of changeset.ContractSet, and acts an an adapter to the internal package
type ContractSet struct {
	commonchangeset.MCMSWithTimelockState
	Forwarder map[common.Address]*forwarder.KeystoneForwarder
}

func GetContractSets(lggr logger.Logger, req *GetContractSetsRequest) (*GetContractSetsResponse, error) {
	resp := &GetContractSetsResponse{
		ContractSets: make(map[uint64]ContractSet),
	}
	for id, chain := range req.Chains {
		addrs, err := req.AddressBook.AddressesForChain(id)
		if err != nil {
			return nil, fmt.Errorf("failed to get addresses for chain %d: %w", id, err)
		}

		// Forwarder addresses now have informative labels, but we don't want them to be ignored if no labels are provided for filtering.
		// If labels are provided, just filter by those.
		forwarderAddrs := make(map[string]cldf.TypeAndVersion)
		if len(req.Labels) == 0 {
			for addr, tv := range addrs {
				if tv.Type == globals.KeystoneForwarder {
					forwarderAddrs[addr] = tv
				}
			}
		}

		filtered := deployment.LabeledAddresses(addrs).And(req.Labels...)

		for addr, tv := range forwarderAddrs {
			filtered[addr] = tv
		}

		cs, err := loadContractSet(lggr, chain, filtered)
		if err != nil {
			return nil, fmt.Errorf("failed to load contract set for chain %d: %w", id, err)
		}
		resp.ContractSets[id] = *cs
	}
	return resp, nil
}

// loadContractSet loads the MCMS state and then sets the SmartData contract state.
func loadContractSet(
	lggr logger.Logger,
	chain cldf_evm.Chain,
	addresses map[string]cldf.TypeAndVersion,
) (*ContractSet, error) {
	var out ContractSet
	mcmsWithTimelock, err := commonchangeset.MaybeLoadMCMSWithTimelockChainState(chain, addresses)
	if err != nil {
		return nil, fmt.Errorf("failed to load mcms contract: %w", err)
	}
	out.MCMSWithTimelockState = *mcmsWithTimelock

	if err := setContracts(lggr, addresses, chain.Client, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

// setContracts sets the SmartData contract state. Other contracts are ignored.
func setContracts(
	lggr logger.Logger,
	addresses map[string]cldf.TypeAndVersion,
	client cldf_evm.OnchainClient,
	set *ContractSet,
) error {
	for addr, tv := range addresses {
		// todo handle versions
		switch tv.Type {
		case globals.KeystoneForwarder:
			c, err := forwarder.NewKeystoneForwarder(common.HexToAddress(addr), client)
			if err != nil {
				return fmt.Errorf("failed to create forwarder contract from address %s: %w", addr, err)
			}
			if set.Forwarder == nil {
				set.Forwarder = make(map[common.Address]*forwarder.KeystoneForwarder)
			}

			set.Forwarder[common.HexToAddress(addr)] = c
		default:
			// do nothing, non-exhaustive
			lggr.Warnf("skipping contract of type : %s", tv.Type)
		}
	}
	return nil
}
