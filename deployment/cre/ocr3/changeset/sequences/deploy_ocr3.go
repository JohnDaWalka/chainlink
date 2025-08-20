package sequences

import (
	"github.com/Masterminds/semver/v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3/changeset/operations/contracts"
)

type DeployOCR3Deps struct {
	Env *cldf.Environment
}

type DeployOCR3Input struct {
	RegistryChainSel uint64
	Qualifier        string
}

func (c DeployOCR3Input) Validate() error {
	return nil
}

type DeployOCR3Output struct {
	ChainSelector uint64
	Address       string
	Type          string
	Version       string
	Labels        []string
}

var DeployOCR3 = operations.NewSequence(
	"deploy-ocr3",
	semver.MustParse("1.0.0"),
	"Deploys the OCR3 contract",
	func(b operations.Bundle, deps DeployOCR3Deps, input DeployOCR3Input) (DeployOCR3Output, error) {
		report, err := operations.ExecuteOperation(b, contracts.DeployOCR3, contracts.DeployOCR3Deps{Env: deps.Env}, contracts.DeployOCR3Input{
			ChainSelector: input.RegistryChainSel,
			Qualifier:     input.Qualifier,
		})
		if err != nil {
			return DeployOCR3Output{}, err
		}

		return DeployOCR3Output{
			ChainSelector: report.Output.ChainSelector,
			Address:       report.Output.Address,
			Type:          report.Output.Type,
			Version:       report.Output.Version,
			Labels:        report.Output.Labels,
		}, nil
	},
)
