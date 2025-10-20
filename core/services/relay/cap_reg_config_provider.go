package relay

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"strconv"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var _ ocrtypes.ContractConfigTracker = &CapRegConfigProvider{}

// CapRegConfigProvider subscribes to the registrySyncer for on-chain changes in the Capability Registry.
// Parses config from on-chain Capability Config.
type CapRegConfigProvider struct {
	services.StateMachine
	lggr logger.Logger

	lastSyncedBlockHeight string
	localConfig           ocrtypes.ContractConfig
	donID                 registrysyncer.DonID
	capability            string
}

func NewCapRegConfigProvider(ctx context.Context, lggr logger.Logger, donID uint32, capability string) (*CapRegConfigProvider, error) {
	return newCapRegConfigProvider(ctx, lggr, donID, capability)
}

func newCapRegConfigProvider(ctx context.Context, lggr logger.Logger, donID uint32, capability string) (*CapRegConfigProvider, error) {
	cp := &CapRegConfigProvider{
		lggr:       logger.Named(lggr, "ConfigPoller"),
		donID:      registrysyncer.DonID(donID),
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
var _ registrysyncer.Listener = &CapRegConfigProvider{}

func (cp *CapRegConfigProvider) OnNewRegistry(ctx context.Context, registry *registrysyncer.LocalRegistry) error {
	if registry == nil {
		return errors.New("registry is nil")
	}

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

	cp.lastSyncedBlockHeight = registry.LastSyncedBlockHeight

	// This config is on-chain in the Capability Registry
	newOnChainConfig := capConfig.Config

	/* TODO: We also want to handle changes to these configs:
	don.ID
	don.F
	...
	*/

	if !bytes.Equal(newOnChainConfig, cp.localConfig.OnchainConfig) {
		cp.lggr.Infow("capability config updated", "donID", cp.donID, "capability", cp.capability)
		cp.localConfig.OnchainConfig = newOnChainConfig
	}
	return nil
}

// LatestConfigDetails returns the latest config details from the logs
func (cp *CapRegConfigProvider) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	blockHeight, err := cp.LatestBlockHeight(ctx)
	if err != nil {
		return 0, ocrtypes.ConfigDigest{}, err
	}
	// TODO: Implement Config Digest...
	return blockHeight, ocrtypes.ConfigDigest{}, errors.New("Unimplemented")
}

// LatestConfig returns the latest config from the logs on a certain block
func (cp *CapRegConfigProvider) LatestConfig(ctx context.Context, changedInBlock uint64) (ocrtypes.ContractConfig, error) {
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
func (cp *CapRegConfigProvider) LatestBlockHeight(_ context.Context) (blockHeight uint64, err error) {
	blockHeight, err = strconv.ParseUint(cp.lastSyncedBlockHeight, 10, 64)
	if err != nil {
		return 0, err
	}
	return blockHeight, err
}

func (cp *CapRegConfigProvider) Notify() <-chan struct{} {
	return nil
}
