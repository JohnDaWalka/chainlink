package v0_5_0

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/txutil"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/configurator"
	"github.com/smartcontractkit/chainlink/deployment"
)

var SetProductionConfigChangeset = deployment.CreateChangeSet(setProductionConfigLogic, setProductionConfigPrecondition)

// SetProductionConfigConfig contains the parameters needed to set a production config.
type SetProductionConfigConfig struct {
	ConfigurationsByChain map[uint64][]SetProductionConfig
	MCMSConfig            *types.MCMSConfig
}

// SetProductionConfig groups the parameters for a production config call.
type SetProductionConfig struct {
	ConfiguratorAddress   common.Address
	ConfigID              [32]byte
	Signers               [][]byte
	OffchainTransmitters  [][32]byte
	F                     uint8
	OnchainConfig         []byte
	OffchainConfigVersion uint64
	OffchainConfig        []byte
}

func (pc SetProductionConfig) GetContractAddress() common.Address { return pc.ConfiguratorAddress }

func (cfg SetProductionConfigConfig) Validate() error {
	if len(cfg.ConfigurationsByChain) == 0 {
		return errors.New("ConfigurationsByChain cannot be empty")
	}
	return nil
}

func setProductionConfigPrecondition(_ deployment.Environment, cfg SetProductionConfigConfig) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid DeployConfiguratorConfig: %w", err)
	}

	return nil
}

func setProductionConfigLogic(e deployment.Environment, cfg SetProductionConfigConfig) (deployment.ChangesetOutput, error) {
	txs, err := txutil.GetTxs(
		e,
		types.Configurator.String(),
		cfg.ConfigurationsByChain,
		LoadConfigurator,
		doSetProductionConfig,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed building SetProductionConfig txs: %w", err)
	}

	return mcmsutil.ExecuteOrPropose(e, txs, cfg.MCMSConfig, "SetProductionConfig proposal")
}

func doSetProductionConfig(
	c *configurator.Configurator,
	prodCfg SetProductionConfig,
) (*ethTypes.Transaction, error) {
	return c.SetProductionConfig(deployment.SimTransactOpts(),
		prodCfg.ConfigID,
		prodCfg.Signers,
		prodCfg.OffchainTransmitters,
		prodCfg.F,
		prodCfg.OnchainConfig,
		prodCfg.OffchainConfigVersion,
		prodCfg.OffchainConfig)
}
