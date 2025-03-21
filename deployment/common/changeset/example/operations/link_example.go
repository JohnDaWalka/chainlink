package example

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/operations"
)

var _ deployment.ChangeSetV2[SqDeployLinkInput] = LinkExampleChangeset{}

type SqDeployLinkInput struct {
	MintAmount *big.Int
	Amount     *big.Int
	To         common.Address
	chainID    uint64
}

type SqDeployLinkOutput struct {
	Address common.Address
}

type EthereumDeps struct {
	Auth  *bind.TransactOpts
	Chain deployment.Chain
	AB    deployment.AddressBook
}

type LinkExampleChangeset struct{}

func (l LinkExampleChangeset) VerifyPreconditions(e deployment.Environment, config SqDeployLinkInput) error {
	return nil
}

func (l LinkExampleChangeset) Apply(e deployment.Environment, config SqDeployLinkInput) (deployment.ChangesetOutput, error) {
	auth := e.Chains[config.chainID].DeployerKey
	ab := deployment.NewMemoryAddressBook()
	deps := EthereumDeps{
		Auth:  auth,
		Chain: e.Chains[config.chainID],
		AB:    ab,
	}

	seqReport, err := operations.ExecuteSequence(e.OperationsBundle, LinkExampleSequence, deps, config)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	return deployment.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReport.ExecutionReports,
	}, nil
}
