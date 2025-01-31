package changeset

import (
	"github.com/smartcontractkit/chainlink/deployment"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

var _ deployment.ChangeSet[InitialContractsCfg] = ConfigureInitialContractsChangeset

type InitialContractsCfg struct {
	RegistryChainSel uint64
	Dons             []kslib.DonCapabilities
	OCR3Config       *kslib.OracleConfig
}

func ConfigureInitialContractsChangeset(e deployment.Environment, cfg InitialContractsCfg) (deployment.ChangesetOutput, error) {

	return deployment.ChangesetOutput{}, nil

}
