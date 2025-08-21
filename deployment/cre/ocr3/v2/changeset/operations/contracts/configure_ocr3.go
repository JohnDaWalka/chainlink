package contracts

import (
	"errors"
	"fmt"
	"io"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	kchangeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	kinternal "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

type ConfigureOCR3Deps struct {
	Env                  *cldf.Environment
	WriteGeneratedConfig io.Writer
	Registry             *capabilities_registry_v2.CapabilitiesRegistry
}

type ConfigureOCR3Input struct {
	ContractAddress  *common.Address
	RegistryChainSel uint64
	DONs             []ConfigureKeystoneDON
	Config           *kinternal.OracleConfig
	DryRun           bool

	MCMSConfig *kchangeset.MCMSConfig
}

func (i ConfigureOCR3Input) UseMCMS() bool {
	return i.MCMSConfig != nil
}

type ConfigureOCR3OpOutput struct {
	MCMSTimelockProposals []mcms.TimelockProposal
}

var ConfigureOCR3 = operations.NewOperation[ConfigureOCR3Input, ConfigureOCR3OpOutput, ConfigureOCR3Deps](
	"configure-ocr3-op",
	semver.MustParse("1.0.0"),
	"Configure OCR3 Contract",
	func(b operations.Bundle, deps ConfigureOCR3Deps, input ConfigureOCR3Input) (ConfigureOCR3OpOutput, error) {
		if input.ContractAddress == nil {
			return ConfigureOCR3OpOutput{}, errors.New("ContractAddress is required")
		}

		var nodeIDs []string
		for _, don := range input.DONs {
			donConfig := RegisteredDonConfig{
				NodeIDs:          don.NodeIDs,
				Name:             don.Name,
				RegistryChainSel: input.RegistryChainSel,
				Registry:         deps.Registry,
			}
			d, err := NewRegisteredDon(*deps.Env, donConfig)
			if err != nil {
				return ConfigureOCR3OpOutput{}, fmt.Errorf("configure-ocr3-op failed: failed to create registered DON %s: %w", don.Name, err)
			}

			// We double-check that the DON accepts workflows...
			if d.Info.AcceptsWorkflows {
				for _, node := range d.Nodes {
					nodeIDs = append(nodeIDs, node.NodeID)
				}
			}
		}

		resp, err := kchangeset.ConfigureOCR3Contract(*deps.Env, kchangeset.ConfigureOCR3Config{
			ChainSel:             input.RegistryChainSel,
			NodeIDs:              nodeIDs,
			Address:              input.ContractAddress,
			OCR3Config:           input.Config,
			DryRun:               input.DryRun,
			WriteGeneratedConfig: deps.WriteGeneratedConfig,
			MCMSConfig:           input.MCMSConfig,
		})
		if err != nil {
			return ConfigureOCR3OpOutput{}, fmt.Errorf("configure-ocr3-op failed: %w", err)
		}

		return ConfigureOCR3OpOutput{MCMSTimelockProposals: resp.MCMSTimelockProposals}, nil
	},
)
