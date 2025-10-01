package evm

import (
	"bytes"
	"context"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"

	"github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	evmRelayTypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

var _ evmRelayTypes.ConfigPoller = &CapRegConfigPoller{}

// CapRegConfigPoller subscribes to the registrySyncer for on-chain changes in the Capability Registry.
// Parses config from on-chain Capability Config.
type CapRegConfigPoller struct {
	services.StateMachine
	lggr logger.Logger

	localConfig ocrtypes.ContractConfig
	donID       registrysyncer.DonID
	capability  string
}

func NewCapRegConfigPoller(ctx context.Context, lggr logger.Logger, donID registrysyncer.DonID, capability string) (*CapRegConfigPoller, error) {
	return newCapRegConfigPoller(ctx, lggr, donID, capability)
}

func newCapRegConfigPoller(ctx context.Context, lggr logger.Logger, donID registrysyncer.DonID, capability string) (*CapRegConfigPoller, error) {
	cp := &CapRegConfigPoller{
		lggr:       logger.Named(lggr, "ConfigPoller"),
		donID:      donID,
		capability: capability,
		localConfig: ocrtypes.ContractConfig{
			// TODO: Set initial config?
			// TODO: Do we get all this from the cap-reg? What about offChainConfig?
			ConfigDigest:          ocrtypes.ConfigDigest{},
			ConfigCount:           0,
			Signers:               nil,
			Transmitters:          nil,
			F:                     0,
			OnchainConfig:         nil,
			OffchainConfigVersion: 0,
			OffchainConfig:        nil,
		},
	}

	return cp, nil
}

// Subscribes to registry syncer for config changes
var _ registrysyncer.Listener = &CapRegConfigPoller{}

func (cp *CapRegConfigPoller) OnNewRegistry(ctx context.Context, registry *registrysyncer.LocalRegistry) error {
	don, ok := registry.IDsToDONs[cp.donID]
	if !ok {
		cp.lggr.Warnw("DON not found in registry", "donID", cp.donID)
		return nil
	}

	capConfig, ok := don.CapabilityConfigurations[cp.capability]
	if !ok {
		cp.lggr.Warnw("Capability not found for DON", "donID", cp.donID, "capability", cp.capability)
		return nil
	}

	// This config is on-chain in the Capability Registry
	newOnChainConfig := capConfig.Config

	don.ID
	don.
		don.F

	// TODO: We also care if the DON config was updated
	if !bytes.Equal(newOnChainConfig, cp.localConfig.OnchainConfig) {
		cp.lggr.Infow("capability config updated", "donID", cp.donID, "capability", cp.capability)
		cp.localConfig.OnchainConfig = newOnChainConfig
	}
	return nil
}

func (cp *CapRegConfigPoller) Start() {}

func (cp *CapRegConfigPoller) Close() error {
	return nil
}

// Notify noop method
func (cp *CapRegConfigPoller) Notify() <-chan struct{} {
	return nil
}

// Replay abstracts the logpoller.LogPoller Replay() implementation
func (cp *CapRegConfigPoller) Replay(ctx context.Context, fromBlock int64) error {
	return errors.New("Unimplemented")
}

// LatestConfigDetails returns the latest config details from the logs
func (cp *CapRegConfigPoller) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	// TODO: Do we need this?
	return 0, ocrtypes.ConfigDigest{}, errors.New("Unimplemented")
}

// LatestConfig returns the latest config from the logs on a certain block
func (cp *CapRegConfigPoller) LatestConfig(ctx context.Context, changedInBlock uint64) (ocrtypes.ContractConfig, error) {
	latestConfigSet := ocrtypes.ContractConfig{
		ConfigDigest:          cp.localConfig.ConfigDigest,
		ConfigCount:           cp.localConfig.ConfigCount,
		Signers:               cp.localConfig.Signers,
		Transmitters:          cp.localConfig.Transmitters,
		F:                     cp.localConfig.F,
		OnchainConfig:         cp.localConfig.OnchainConfig,
		OffchainConfigVersion: cp.localConfig.OffchainConfigVersion,
		OffchainConfig:        cp.localConfig.OffchainConfig,
	}
	cp.lggr.Infow("LatestConfig", "latestConfig", latestConfigSet)
	return latestConfigSet, nil
}

// LatestBlockHeight returns the latest block height from the logs
func (cp *CapRegConfigPoller) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	// TODO: There's no reason for this poller to get the block height?
	// TODO: We can return it if this gets called for some reason however.
	/*
		latest, err := cp.destChainLogPoller.LatestBlock(ctx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return 0, nil
			}
			return 0, err
		}
		return uint64(latest.BlockNumber), nil
	*/
	return 0, errors.New("Unimplemented")
}
