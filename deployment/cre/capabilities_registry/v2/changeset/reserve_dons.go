package changeset

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
)

var _ cldf.ChangeSetV2[ReserveDonsInput] = ReserveDons{}

// ReserveDonsInput must be JSON and YAML Serializable with no private fields
type ReserveDonsInput struct {
	Address       string `json:"address" yaml:"address"`
	ChainSelector uint64 `json:"chainSelector" yaml:"chainSelector"`
	N             int    `json:"n" yaml:"n"`
}

type ReserveDons struct{}

func (r ReserveDons) VerifyPreconditions(e cldf.Environment, config ReserveDonsInput) error {
	if config.Address == "" {
		return fmt.Errorf("address is not set")
	}
	if !common.IsHexAddress(config.Address) {
		return fmt.Errorf("address '%s' is not a valid hex address", config.Address)
	}
	if config.N <= 0 {
		return fmt.Errorf("N must be greater than 0, got %d", config.N)
	}
	if _, ok := e.BlockChains.EVMChains()[config.ChainSelector]; !ok {
		return fmt.Errorf("chain %d not found in environment", config.ChainSelector)
	}
	return nil
}

func (r ReserveDons) Apply(e cldf.Environment, config ReserveDonsInput) (cldf.ChangesetOutput, error) {
	// Generate globally unique DON names
	timestamp := time.Now().Unix()
	donNames := make([]string, config.N)
	for i := 0; i < config.N; i++ {
		donNames[i] = fmt.Sprintf("reserved-%d-%d", timestamp, i)
	}

	// Generate 4 fake P2P IDs (32-byte format)
	fakeP2PIDs := []string{
		"12D3KooWJzSEJxTfvA9S8tav5L8VdxkDZ9UpBuFk1vxWbDQJUjnR",
		"12D3KooWBhbBNNJfKu3xDhz3DfVxkjdpWKYHdLFyJGrP8Xw8Kn1A",
		"12D3KooWCdMKjesUMEz1Ph6JqCBkcBZQNYcKNP5GjEmqo6eHarqt",
		"12D3KooWFRgZxGnbwxfTwNhGzwjdgsrfyJBLT6YHhKNWB4qTK8qE",
	}

	// Convert P2P IDs to [32]byte format
	p2pIDsBytes := make([][32]byte, len(fakeP2PIDs))
	for i, id := range fakeP2PIDs {
		// For simplicity, we'll use the first 32 bytes of the string
		// In a real implementation, you'd properly decode the P2P ID
		copy(p2pIDsBytes[i][:], []byte(id)[:32])
	}

	// Create DON configurations
	dons := make([]capabilities_registry_v2.CapabilitiesRegistryNewDONParams, config.N)
	for i := 0; i < config.N; i++ {
		dons[i] = capabilities_registry_v2.CapabilitiesRegistryNewDONParams{
			Name:                     donNames[i],
			Nodes:                    p2pIDsBytes,
			F:                        1, // Minimum fault tolerance
			CapabilityConfigurations: nil,
		}
	}

	// Execute RegisterDons operation using the operations framework
	registerDonsReport, err := operations.ExecuteOperation(
		e.OperationsBundle,
		contracts.RegisterDons,
		contracts.RegisterDonsDeps{
			Env: &e,
		},
		contracts.RegisterDonsInput{
			Address:       config.Address,
			ChainSelector: config.ChainSelector,
			DONs:          dons,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to register reserved DONs: %w", err)
	}

	reports := make([]operations.Report[any, any], 0)
	reports = append(reports, registerDonsReport.ToGenericReport())

	return cldf.ChangesetOutput{
		Reports: reports,
	}, nil
}
