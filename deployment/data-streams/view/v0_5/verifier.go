package v0_5

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	dsutil "github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/interfaces"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_v0_5_0"
)

// VerifierState represents a single verifier configuration state
type VerifierState struct {
	ConfigDigest            string   `json:"configDigest"`
	LatestConfigBlockNumber uint32   `json:"latestConfigBlockNumber"`
	IsActive                bool     `json:"isActive"`
	F                       uint8    `json:"f"`
	Signers                 []string `json:"signers"`
}

// VerifierView represents a simplified view of the verifier contract state
type VerifierView struct {
	Configs        map[string]*VerifierState `json:"configs"`
	TypeAndVersion string                    `json:"typeAndVersion,omitempty"`
	Owner          common.Address            `json:"owner,omitempty"`
}

// Ensure VerifierView implements the ContractView interface
var _ interfaces.ContractView = (*VerifierView)(nil)

// NewVerifierView creates a new empty VerifierView
func NewVerifierView() *VerifierView {
	return &VerifierView{
		Configs: make(map[string]*VerifierState),
	}
}

// SerializeView serializes the VerifierView to JSON
func (v *VerifierView) SerializeView() (string, error) {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal verifier view: %w", err)
	}
	return string(bytes), nil
}

// GetVerifierState returns the VerifierState for a specific configDigest
func (v *VerifierView) GetVerifierState(configDigest string) (*VerifierState, error) {
	state, ok := v.Configs[configDigest]
	if !ok {
		return nil, fmt.Errorf("configDigest %s not found", configDigest)
	}
	return state, nil
}

// VerifierViewBuilder builds views for the Verifier contract
type VerifierViewBuilder struct{}

type VerifierContext struct {
	FromBlock uint64
	ToBlock   *uint64
}

// BuildView builds a view of the Verifier contract state from logs and calls
func (b *VerifierViewBuilder) BuildView(ctx context.Context, verifier *verifier_v0_5_0.Verifier, chainParams VerifierContext) (VerifierView, error) {
	view := NewVerifierView()

	// Get contract owner
	owner, err := verifier.Owner(&bind.CallOpts{Context: ctx})
	if err != nil {
		return VerifierView{}, fmt.Errorf("failed to get contract owner: %w", err)
	}
	view.Owner = owner

	// Define the filter options
	filterOpts := &bind.FilterOpts{
		Start:   chainParams.FromBlock,
		End:     chainParams.ToBlock,
		Context: ctx,
	}

	// Process all ConfigSet events
	if err := b.processConfigSetEvents(filterOpts, verifier, view); err != nil {
		return VerifierView{}, err
	}

	// Process all ConfigUpdated events
	if err := b.processConfigUpdatedEvents(filterOpts, verifier, view); err != nil {
		return VerifierView{}, err
	}

	// Process all ConfigActivated events
	if err := b.processConfigActivatedEvents(filterOpts, verifier, view); err != nil {
		return VerifierView{}, err
	}

	// Process all ConfigDeactivated events
	if err := b.processConfigDeactivatedEvents(filterOpts, verifier, view); err != nil {
		return VerifierView{}, err
	}

	return *view, nil
}

// processConfigSetEvents processes ConfigSet events
func (b *VerifierViewBuilder) processConfigSetEvents(opts *bind.FilterOpts, verifier *verifier_v0_5_0.Verifier, view *VerifierView) error {
	iter, err := verifier.FilterConfigSet(opts, nil)
	if err != nil {
		return fmt.Errorf("failed to filter ConfigSet events: %w", err)
	}
	defer iter.Close()

	for iter.Next() {
		event := iter.Event
		configDigestHex := dsutil.HexEncodeBytes(event.ConfigDigest[:])

		// Get block number for this config
		blockNumber, err := verifier.LatestConfigDetails(&bind.CallOpts{}, event.ConfigDigest)
		if err != nil {
			return fmt.Errorf("failed to get latest config details: %w", err)
		}

		state := &VerifierState{
			ConfigDigest:            configDigestHex,
			LatestConfigBlockNumber: blockNumber,
			IsActive:                true, // New configs are active by default
			F:                       event.F,
			Signers:                 make([]string, 0, len(event.Signers)),
		}

		// Add signers
		for _, signer := range event.Signers {
			state.Signers = append(state.Signers, signer.String())
		}

		view.Configs[configDigestHex] = state
	}

	return nil
}

// processConfigUpdatedEvents processes ConfigUpdated events
func (b *VerifierViewBuilder) processConfigUpdatedEvents(opts *bind.FilterOpts, verifier *verifier_v0_5_0.Verifier, view *VerifierView) error {
	iter, err := verifier.FilterConfigUpdated(opts, nil)
	if err != nil {
		return fmt.Errorf("failed to filter ConfigUpdated events: %w", err)
	}
	defer iter.Close()

	for iter.Next() {
		event := iter.Event
		configDigestHex := dsutil.HexEncodeBytes(event.ConfigDigest[:])

		// Skip if this configDigest doesn't exist yet
		state, exists := view.Configs[configDigestHex]
		if !exists {
			// This is unexpected, but we'll create a new state for this config
			state = &VerifierState{
				ConfigDigest: configDigestHex,
				IsActive:     true,
				Signers:      make([]string, 0),
			}

			// Get block number for this config
			blockNumber, err := verifier.LatestConfigDetails(&bind.CallOpts{}, event.ConfigDigest)
			if err != nil {
				return fmt.Errorf("failed to get latest config details: %w", err)
			}
			state.LatestConfigBlockNumber = blockNumber
			view.Configs[configDigestHex] = state
		}

		// Update signers with new set
		state.Signers = make([]string, 0, len(event.NewSigners))

		// Add new signers
		for _, signer := range event.NewSigners {
			state.Signers = append(state.Signers, signer.String())
		}
	}

	return nil
}

// processConfigActivatedEvents processes ConfigActivated events
func (b *VerifierViewBuilder) processConfigActivatedEvents(opts *bind.FilterOpts, verifier *verifier_v0_5_0.Verifier, view *VerifierView) error {
	iter, err := verifier.FilterConfigActivated(opts, nil)
	if err != nil {
		return fmt.Errorf("failed to filter ConfigActivated events: %w", err)
	}
	defer iter.Close()

	for iter.Next() {
		event := iter.Event
		configDigestHex := dsutil.HexEncodeBytes(event.ConfigDigest[:])

		// Skip if this configDigest doesn't exist yet
		state, exists := view.Configs[configDigestHex]
		if !exists {
			continue
		}

		state.IsActive = true
	}

	return nil
}

// processConfigDeactivatedEvents processes ConfigDeactivated events
func (b *VerifierViewBuilder) processConfigDeactivatedEvents(opts *bind.FilterOpts, verifier *verifier_v0_5_0.Verifier, view *VerifierView) error {
	iter, err := verifier.FilterConfigDeactivated(opts, nil)
	if err != nil {
		return fmt.Errorf("failed to filter ConfigDeactivated events: %w", err)
	}
	defer iter.Close()

	for iter.Next() {
		event := iter.Event
		configDigestHex := dsutil.HexEncodeBytes(event.ConfigDigest[:])

		// Skip if this configDigest doesn't exist yet
		state, exists := view.Configs[configDigestHex]
		if !exists {
			continue
		}

		state.IsActive = false
	}

	return nil
}
