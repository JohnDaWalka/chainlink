package don

import (
	"context"
	"slices"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

func CreateJobs(ctx context.Context, testLogger zerolog.Logger, input cre.CreateJobsInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "input validation failed")
	}

	for _, don := range input.DonTopology.DonsWithMetadata {
		if jobSpecs, ok := input.DonToJobSpecs[don.ID]; ok {
			createErr := jobs.Create(ctx, input.CldEnv.Offchain, jobSpecs)
			if createErr != nil {
				return errors.Wrapf(createErr, "failed to create jobs for DON %d", don.ID)
			}
		} else {
			testLogger.Warn().Msgf("No job specs found for DON %d", don.ID)
		}
	}

	return nil
}

func ValidateTopology(nodeSetInput []*cre.CapabilitiesAwareNodeSet, infraInput infra.Input) error {
	if len(nodeSetInput) == 0 {
		return errors.New("at least one nodeset is required")
	}

	hasAtLeastOneBootstrapNode := false
	for _, nodeSet := range nodeSetInput {
		if nodeSet.BootstrapNodeIndex != -1 {
			hasAtLeastOneBootstrapNode = true
			break
		}
	}

	if !hasAtLeastOneBootstrapNode {
		return errors.New("at least one nodeSet must have a bootstrap node")
	}

	workflowDONHasBootstrapNode := false
	for _, nodeSet := range nodeSetInput {
		if nodeSet.BootstrapNodeIndex != -1 && slices.Contains(nodeSet.DONTypes, cre.WorkflowDON) {
			workflowDONHasBootstrapNode = true
			break
		}
	}

	if !workflowDONHasBootstrapNode {
		return errors.New("due to the limitations of our implementation, workflow DON must always have a bootstrap node")
	}

	isGatewayRequired := false
	for _, nodeSet := range nodeSetInput {
		if NodeNeedsAnyGateway(nodeSet.ComputedCapabilities) {
			isGatewayRequired = true
			break
		}
	}

	if !isGatewayRequired {
		return nil
	}

	anyDONHasGatewayConfigured := false
	for _, nodeSet := range nodeSetInput {
		if isGatewayRequired {
			if flags.HasFlag(nodeSet.DONTypes, cre.GatewayDON) && nodeSet.GatewayNodeIndex != -1 {
				anyDONHasGatewayConfigured = true
				break
			}
		}
	}

	if !anyDONHasGatewayConfigured {
		return errors.New("at least one DON must be configured with gateway DON type and have a gateway node index set, because at least one DON requires gateway due to its capabilities")
	}

	return nil
}

func BuildTopology(nodeSetInput []*cre.CapabilitiesAwareNodeSet, infraInput infra.Input, homeChainSelector uint64) (*cre.Topology, error) {
	return cre.NewTopology(nodeSetInput, infraInput, homeChainSelector)
}

func AnyDonHasCapability(donMetadata []*cre.DonMetadata, capability cre.CapabilityFlag) bool {
	for _, don := range donMetadata {
		if flags.HasFlagForAnyChain(don.Flags, capability) {
			return true
		}
	}

	return false
}

func NodeNeedsAnyGateway(nodeFlags []cre.CapabilityFlag) bool {
	return flags.HasFlag(nodeFlags, cre.CustomComputeCapability) ||
		flags.HasFlag(nodeFlags, cre.WebAPITriggerCapability) ||
		flags.HasFlag(nodeFlags, cre.WebAPITargetCapability) ||
		flags.HasFlag(nodeFlags, cre.VaultCapability) ||
		flags.HasFlag(nodeFlags, cre.HTTPActionCapability) ||
		flags.HasFlag(nodeFlags, cre.HTTPTriggerCapability)
}

func NodeNeedsWebAPIGateway(nodeFlags []cre.CapabilityFlag) bool {
	return flags.HasFlag(nodeFlags, cre.CustomComputeCapability) ||
		flags.HasFlag(nodeFlags, cre.WebAPITriggerCapability) ||
		flags.HasFlag(nodeFlags, cre.WebAPITargetCapability)
}
