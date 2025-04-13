package ocrimpls

import (
	"context"
	"encoding/hex"

	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// Aptos config tracker implementation. Aptos stores transmitters as ed25519 public keys instead of addresses, so no
// address codec is required.
type aptosConfigTracker struct {
	cfg cctypes.OCR3ConfigWithMeta
}

func NewAptosConfigTracker(cfg cctypes.OCR3ConfigWithMeta) *aptosConfigTracker {
	return &aptosConfigTracker{cfg: cfg}
}

func (c *aptosConfigTracker) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	return 0, nil
}

func (c *aptosConfigTracker) LatestConfig(ctx context.Context, changedInBlock uint64) (types.ContractConfig, error) {
	return c.contractConfig(), nil
}

func (c *aptosConfigTracker) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest types.ConfigDigest, err error) {
	return 0, c.cfg.ConfigDigest, nil
}

func (c *aptosConfigTracker) Notify() <-chan struct{} {
	return nil
}

func (c *aptosConfigTracker) contractConfig() types.ContractConfig {
	var signers [][]byte
	var transmitters [][]byte
	for _, node := range c.cfg.Config.Nodes {
		signers = append(signers, node.SignerKey)
		transmitters = append(transmitters, node.TransmitterKey)
	}

	return types.ContractConfig{
		ConfigDigest:          c.cfg.ConfigDigest,
		ConfigCount:           uint64(c.cfg.Version),
		Signers:               c.toOnchainPublicKeys(signers),
		Transmitters:          c.toOCRAccounts(transmitters),
		F:                     c.cfg.Config.FRoleDON,
		OnchainConfig:         []byte{},
		OffchainConfigVersion: c.cfg.Config.OffchainConfigVersion,
		OffchainConfig:        c.cfg.Config.OffchainConfig,
	}
}

func (c *aptosConfigTracker) PublicConfig() (ocr3confighelper.PublicConfig, error) {
	return ocr3confighelper.PublicConfigFromContractConfig(false, c.contractConfig())
}

func (c *aptosConfigTracker) toOnchainPublicKeys(signers [][]byte) []types.OnchainPublicKey {
	keys := make([]types.OnchainPublicKey, len(signers))
	for i, signer := range signers {
		keys[i] = types.OnchainPublicKey(signer)
	}
	return keys
}

func (c *aptosConfigTracker) toOCRAccounts(transmitters [][]byte) []types.Account {
	accounts := make([]types.Account, len(transmitters))
	for i, transmitter := range transmitters {
		s := hex.EncodeToString(transmitter)
		accounts[i] = types.Account(s)
		continue
	}
	return accounts
}

var _ types.ContractConfigTracker = (*aptosConfigTracker)(nil)
