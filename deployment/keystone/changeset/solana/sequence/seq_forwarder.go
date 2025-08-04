package sequence

import (
	"fmt"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	commonOps "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/operations"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana/sequence/operation"
	"github.com/smartcontractkit/mcms"
)

var (
	DeployForwarderSeq = operations.NewSequence(
		"deploy-forwarder-seq",
		operation.Version1_0_0,
		"Deploys forwarder contract and initializes it",
		deployForwarder,
	)
	ConfigureForwarderSeq = operations.NewSequence(
		"configure-forwarder-seq",
		operation.Version1_0_0,
		"Configure forwarder program",
		configureForwarder,
	)
)

type (
	DeployForwarderSeqInput struct {
		ChainSel     uint64
		ProgramName  string
		Overallocate bool
	}

	DeployForwarderSeqOutput struct {
		ProgramID solana.PublicKey
		State     solana.PublicKey
	}

	ConfigureForwarderSeqInput struct {
		WFDonName string
		// workflow don node ids in the offchain client. Used to fetch and derive the signer keys
		WFNodeIDs        []string
		RegistryChainSel uint64
		Chains           map[uint64]struct{}
		Qualifier        string
		Version          string
	}

	ConfigureForwarderSeqOutput struct {
		Proposals []mcms.TimelockProposal
	}
)

const KeystoneForwarderProgramSize = 5 * 1024 * 1024

func configureForwarder(b operations.Bundle, deps operation.Deps, in ConfigureForwarderSeqInput) (ConfigureForwarderSeqOutput, error) {
	var out ConfigureForwarderSeqOutput

	return out, nil
}

func deployForwarder(b operations.Bundle, deps operation.Deps, in DeployForwarderSeqInput) (DeployForwarderSeqOutput, error) {
	var out DeployForwarderSeqOutput

	// 1. Deploy
	deployOut, err := operations.ExecuteOperation(b, operation.DeployForwarderOp, commonOps.Deps{Chain: deps.Chain}, commonOps.DeployInput{
		ProgramName:  in.ProgramName,
		Overallocate: in.Overallocate,
		Size:         KeystoneForwarderProgramSize,
		ChainSel:     in.ChainSel,
	})

	if err != nil {
		return DeployForwarderSeqOutput{}, err
	}
	out.ProgramID = deployOut.Output.ProgramID

	// 2. Initialize
	initOut, err := operations.ExecuteOperation(b, operation.InitForwarderOp, deps, operation.InitForwarderInput{
		ProgramID: out.ProgramID,
		ChainSel:  in.ChainSel,
	})

	if err != nil {
		return DeployForwarderSeqOutput{}, err
	}
	out.State = initOut.Output.StatePubKey

	fmt.Println("deployed forwarder programID ", out.ProgramID, " stateID ", out.State)

	return out, nil
}
