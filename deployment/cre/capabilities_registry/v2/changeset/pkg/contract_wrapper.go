package pkg

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
)

// TODO: we have to support pagination eventually
var (
	MaxCapabilities             = big.NewInt(128)
	MaxDONs                     = big.NewInt(32)
	MaxNodes                    = big.NewInt(256)
	MaxNOPs                     = big.NewInt(128)
	ErrPaginationNotImplemented = errors.New("pagination not implemented")
)

func GetCapabilities(opts *bind.CallOpts, capReg *capabilities_registry_v2.CapabilitiesRegistry) ([]capabilities_registry_v2.CapabilitiesRegistryCapabilityInfo, error) {
	caps, err := capReg.GetCapabilities(opts, big.NewInt(0), MaxCapabilities)
	if len(caps) >= int(MaxCapabilities.Int64()) {
		return nil, fmt.Errorf("too many capabilities, %w", ErrPaginationNotImplemented)
	}
	return caps, cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
}

func GetNodeOperators(opts *bind.CallOpts, capReg *capabilities_registry_v2.CapabilitiesRegistry) ([]capabilities_registry_v2.CapabilitiesRegistryNodeOperatorInfo, error) {
	nops, err := capReg.GetNodeOperators(opts, big.NewInt(0), MaxNOPs)
	if len(nops) >= int(MaxNOPs.Int64()) {
		return nil, fmt.Errorf("too many node operators, %w", ErrPaginationNotImplemented)
	}
	return nops, cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
}

func GetNodes(opts *bind.CallOpts, capReg *capabilities_registry_v2.CapabilitiesRegistry) ([]capabilities_registry_v2.INodeInfoProviderNodeInfo, error) {
	nodes, err := capReg.GetNodes(opts, big.NewInt(0), MaxNodes)
	if len(nodes) >= int(MaxNodes.Int64()) {
		return nil, fmt.Errorf("too many nodes, %w", ErrPaginationNotImplemented)
	}
	return nodes, cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
}

func GetDONs(opts *bind.CallOpts, capReg *capabilities_registry_v2.CapabilitiesRegistry) ([]capabilities_registry_v2.CapabilitiesRegistryDONInfo, error) {
	donsInfo, err := capReg.GetDONs(opts, big.NewInt(0), MaxDONs)
	if len(donsInfo) >= int(MaxDONs.Int64()) {
		return nil, fmt.Errorf("too many DONs, %w", ErrPaginationNotImplemented)
	}
	return donsInfo, cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
}

func GetDONsInFamily(opts *bind.CallOpts, capReg *capabilities_registry_v2.CapabilitiesRegistry, family string) ([]*big.Int, error) {
	donsInfo, err := capReg.GetDONsInFamily(opts, family, big.NewInt(0), MaxDONs)
	if len(donsInfo) >= int(MaxDONs.Int64()) {
		return nil, fmt.Errorf("too many DONs in family %s, %w", family, ErrPaginationNotImplemented)
	}
	return donsInfo, cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
}
