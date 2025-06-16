package operations

import (
	"github.com/Masterminds/semver/v3"
	common "github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type DeployKeystoneContractsSequenceDeps struct {
	// Any dep (for example, the logger, datastore and the env)
	lggr common.Logger
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
		ocr3DeployReport, err := operations.ExecuteOperation(b, DeployOCR3Op, DeployOCR3OpDeps{}, DeployOCR3OpInput{})
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}

		capabilitiesRegistryDeployReport, err := operations.ExecuteOperation(b, DeployCapabilityRegistryOp, DeployCapabilityRegistryOpDeps{}, DeployCapabilityRegistryInput{})
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}

		workflowRegistryDeployReport, err := operations.ExecuteOperation(b, DeployWorkflowRegistryOp, DeployWorkflowRegistryOpDeps{}, DeployWorkflowRegistryInput{})
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}

		keystoneForwarderDeployReport, err := operations.ExecuteOperation(b, DeployKeystoneForwarderOp, DeployForwarderOpDeps{}, DeployForwarderOpInput{})
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
