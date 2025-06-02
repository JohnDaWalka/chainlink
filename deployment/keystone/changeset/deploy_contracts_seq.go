package changeset

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	pkgerrors "github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	"golang.org/x/sync/errgroup"
)

// DeployOCR3Op is an operation that deploys the OCR3 contract.
var DeployOCR3Op = operations.NewOperation[operations.EmptyInput, internal.DeployResponse, DeployKeystoneContractsSequenceDeps](
	"deploy-ocr3-op",
	semver.MustParse("1.0.0"),
	"Deploy OCR3 Contract",
	func(b operations.Bundle, deps DeployKeystoneContractsSequenceDeps, input operations.EmptyInput) (internal.DeployResponse, error) {
		// Here we would implement the logic to deploy the OCR3 contract.
		return internal.DeployResponse{}, nil
	},
)

// DeployCapabilityRegistryOp is an operation that deploys the Capability Registry contract.
var DeployCapabilityRegistryOp = operations.NewOperation[operations.EmptyInput, internal.DeployResponse, DeployKeystoneContractsSequenceDeps](
	"deploy-capability-registry-op",
	semver.MustParse("1.0.0"),
	"Deploy CapabilityRegistry Contract",
	func(b operations.Bundle, deps DeployKeystoneContractsSequenceDeps, input operations.EmptyInput) (internal.DeployResponse, error) {
		// Here we would implement the logic to deploy the OCR3 contract.
		return internal.DeployResponse{}, nil
	},
)

// DeployWorkflowRegistryOp is an operation that deploys the Workflow Registry contract.
var DeployWorkflowRegistryOp = operations.NewOperation[operations.EmptyInput, internal.DeployResponse, DeployKeystoneContractsSequenceDeps](
	"deploy-workflow-registry-op",
	semver.MustParse("1.0.0"),
	"Deploy WorkflowRegistry Contract",
	func(b operations.Bundle, deps DeployKeystoneContractsSequenceDeps, input operations.EmptyInput) (internal.DeployResponse, error) {
		// Here we would implement the logic to deploy the Workflow Registry contract.
		return internal.DeployResponse{}, nil
	},
)

// DeployKeystoneForwarderOp is an operation that deploys the Keystone Forwarder contract.
var DeployKeystoneForwarderOp = operations.NewOperation[operations.EmptyInput, internal.DeployResponse, DeployKeystoneContractsSequenceDeps](
	"deploy-keystone-forwarder-op",
	semver.MustParse("1.0.0"),
	"Deploy KeystoneForwarder Contract",
	func(b operations.Bundle, deps DeployKeystoneContractsSequenceDeps, input operations.EmptyInput) (internal.DeployResponse, error) {
		// Here we would implement the logic to deploy the Keystone Forwarder contract.
		return internal.DeployResponse{}, nil
	},
)

type DeployKeystoneForwardersInput struct {
	Targets []uint64 // The target chains for the Keystone Forwarders
}

type DeployKeystoneForwardersOutput struct {
	AddressBook deployment.AddressBook // The address book containing the deployed Keystone Forwarders
}

var DeployKeystoneForwardersSequence = operations.NewSequence[DeployKeystoneForwardersInput, DeployKeystoneForwardersOutput, DeployKeystoneContractsSequenceDeps](
	"deploy-keystone-forwarders-seq",
	semver.MustParse("1.0.0"),
	"Deploy Keystone Forwarders",
	func(b operations.Bundle, deps DeployKeystoneContractsSequenceDeps, input DeployKeystoneForwardersInput) (DeployKeystoneForwardersOutput, error) {
		ab := deployment.NewMemoryAddressBook()
		contractErrGroup := &errgroup.Group{}
		for _, target := range input.Targets {
			fmt.Println(target)
			contractErrGroup.Go(func() error {
				// For each target, we would deploy the Keystone Forwarder.
				// This is a placeholder for the actual deployment logic.
				// TODO: we would pass here the target as an input to the operation.
				r, err := operations.ExecuteOperation(b, DeployKeystoneForwarderOp, deps, operations.EmptyInput{})
				if err != nil {
					return err
				}
				err = ab.Save(target, r.Output.Address.String(), r.Output.Tv)
				if err != nil {
					return pkgerrors.Wrapf(err, "failed to save Keystone Forwarder address for target %d", target)
				}

				return nil
			})
		}
		if err := contractErrGroup.Wait(); err != nil {
			return DeployKeystoneForwardersOutput{AddressBook: ab}, pkgerrors.Wrap(err, "failed to deploy Keystone contracts")
		}
		return DeployKeystoneForwardersOutput{AddressBook: ab}, nil
	},
)

type DeployKeystoneContractsSequenceDeps struct {
	// Any dep (for example, the logger, datastore and the env)
}

// inputs and outputs have to be serializable, and must not contain sensitive data

type DeployKeystoneContractsSequenceInput struct {
	HomeChainSelector uint64
}

type DeployKeystoneContractsSequenceOutput struct {
	// Not sure if we can serialize the address book without modifications, but whatever is returned needs to be serializable.
	// This could also be the address datastore instead.
	AddressBook deployment.AddressBook
}

// DeployKeystoneContractsSequence is a sequence that deploys the Keystone contracts (OCR3, Capabilities Registry, Workflow Registry, Keystone Forwarder).
var DeployKeystoneContractsSequence = operations.NewSequence[DeployKeystoneContractsSequenceInput, DeployKeystoneContractsSequenceOutput, DeployKeystoneContractsSequenceDeps](
	"deploy-keystone-contracts-seq",
	semver.MustParse("1.0.0"),
	"Deploy Keystone Contracts (OCR3, Capabilities Registry, Workflow Registry, Keystone Forwarder)",
	func(b operations.Bundle, deps DeployKeystoneContractsSequenceDeps, input DeployKeystoneContractsSequenceInput) (output DeployKeystoneContractsSequenceOutput, err error) {
		// TODO: add logging
		// TODO: make it parallelizable

		ab := deployment.NewMemoryAddressBook()

		// Here we would execute one operation per step.
		// For example, deploying the OCR3 contract, then the capabilities registry, etc.
		ocr3DeployReport, err := operations.ExecuteOperation(b, DeployOCR3Op, deps, operations.EmptyInput{})
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}

		capabilitiesRegistryDeployReport, err := operations.ExecuteOperation(b, DeployCapabilityRegistryOp, deps, operations.EmptyInput{})
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}

		workflowRegistryDeployReport, err := operations.ExecuteOperation(b, DeployWorkflowRegistryOp, deps, operations.EmptyInput{})
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}

		keystoneForwarderDeployReport, err := operations.ExecuteOperation(b, DeployKeystoneForwarderOp, deps, operations.EmptyInput{})
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}

		err = ab.Save(input.HomeChainSelector, ocr3DeployReport.Output.Address.String(), ocr3DeployReport.Output.Tv)
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}
		err = ab.Save(input.HomeChainSelector, capabilitiesRegistryDeployReport.Output.Address.String(), capabilitiesRegistryDeployReport.Output.Tv)
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}
		err = ab.Save(input.HomeChainSelector, workflowRegistryDeployReport.Output.Address.String(), workflowRegistryDeployReport.Output.Tv)
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}
		err = ab.Save(input.HomeChainSelector, keystoneForwarderDeployReport.Output.Address.String(), keystoneForwarderDeployReport.Output.Tv)
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}

		// Here we would collect the addresses of the deployed contracts and return them.
		return DeployKeystoneContractsSequenceOutput{
			AddressBook: ab,
		}, nil
	},
)
