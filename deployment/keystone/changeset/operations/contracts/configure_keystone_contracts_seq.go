package contracts

import (
	"errors"
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

type ConfigureKeystoneContractsSeqDeps struct {
	Lggr logger.Logger
	Env  *cldf.Environment
}

type ConfigureKeystoneContractsSeqInput struct {
	RegistryChainSel uint64
	Dons             []internal.DonCapabilities // externally sourced based on the environment
	OCR3Config       *internal.OracleConfig     // TODO: probably should be a map of don to config; but currently we only have one wf don therefore one config
}

func (c ConfigureKeystoneContractsSeqInput) Validate() error {
	if c.OCR3Config == nil {
		return errors.New("OCR3Config is nil")
	}
	for _, don := range c.Dons {
		if err := don.Validate(); err != nil {
			return fmt.Errorf("don validation failed for '%s': %w", don.Name, err)
		}
	}
	_, ok := chainsel.ChainBySelector(c.RegistryChainSel)
	if !ok {
		return fmt.Errorf("chain %d not found in environment", c.RegistryChainSel)
	}
	return nil
}

type ConfigureKeystoneContractsSeqOutput struct {
}
