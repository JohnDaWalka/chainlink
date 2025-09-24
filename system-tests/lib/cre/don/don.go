package don

import (
	"context"
	"fmt"
	"slices"
	"strconv"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/credentials/insecure"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

func CreateJobs(ctx context.Context, testLogger zerolog.Logger, input cre.CreateJobsInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "input validation failed")
	}

	for _, don := range input.DonTopology.DonsWithMetadata {
		if jobSpecs, ok := input.DonToJobSpecs[don.ID]; ok {
			createErr := jobs.Create(ctx, input.CldEnv.Offchain, don.DON.Nodes, jobSpecs)
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
	topology := &cre.Topology{}
	donsWithMetadata := make([]*cre.DonMetadata, len(nodeSetInput))

	for i := range nodeSetInput {
		flags, err := flags.NodeSetFlags(nodeSetInput[i])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get flags for nodeset %s", nodeSetInput[i].Name)
		}

		donsWithMetadata[i] = &cre.DonMetadata{
			ID:              libc.MustSafeUint64FromInt(i + 1), // optimistically set the id to the that which the capabilities registry will assign it
			Flags:           flags,
			NodesMetadata:   make([]*cre.NodeMetadata, len(nodeSetInput[i].NodeSpecs)),
			Name:            nodeSetInput[i].Name,
			SupportedChains: nodeSetInput[i].SupportedChains,
		}
	}

	for donIdx, donMetadata := range donsWithMetadata {
		for nodeIdx := range donMetadata.NodesMetadata {
			nodeWithLabels := cre.NodeMetadata{}
			nodeType := cre.WorkerNode
			if nodeSetInput[donIdx].BootstrapNodeIndex != -1 && nodeIdx == nodeSetInput[donIdx].BootstrapNodeIndex {
				nodeType = cre.BootstrapNode
			}
			nodeWithLabels.Labels = append(nodeWithLabels.Labels, &cre.Label{
				Key:   node.NodeTypeKey,
				Value: nodeType,
			})

			// TODO think whether it would make sense for infraInput to also hold functions that resolve hostnames for various infra and node types
			// and use it with some default, so that we can easily modify it with little effort
			internalHost := InternalHost(nodeIdx, nodeType, donMetadata.Name, infraInput)

			if flags.HasFlag(donMetadata.Flags, cre.GatewayDON) {
				if nodeSetInput[donIdx].GatewayNodeIndex != -1 && nodeIdx == nodeSetInput[donIdx].GatewayNodeIndex {
					nodeWithLabels.Labels = append(nodeWithLabels.Labels, &cre.Label{
						Key:   node.ExtraRolesKey,
						Value: cre.GatewayNode,
					})

					gatewayInternalHost := InternalGatewayHost(nodeIdx, nodeType, donMetadata.Name, infraInput)

					if topology.GatewayConnectorOutput == nil {
						topology.GatewayConnectorOutput = &cre.GatewayConnectorOutput{
							Configurations: make([]*cre.GatewayConfiguration, 0),
						}
					}

					topology.GatewayConnectorOutput.Configurations = append(topology.GatewayConnectorOutput.Configurations, &cre.GatewayConfiguration{
						Outgoing: cre.Outgoing{
							Path: "/node",
							Port: GatewayOutgoingPort,
							Host: gatewayInternalHost,
						},
						Incoming: cre.Incoming{
							Protocol:     "http",
							Path:         "/",
							InternalPort: GatewayIncomingPort,
							ExternalPort: ExternalGatewayPort(infraInput),
							Host:         ExternalGatewayHost(nodeIdx, nodeType, donMetadata.Name, infraInput),
						},
						AuthGatewayID: "cre-gateway",
						// do not set gateway connector dons, they will be resolved automatically
					})
				}
			}

			nodeWithLabels.Labels = append(nodeWithLabels.Labels, &cre.Label{
				Key:   node.IndexKey,
				Value: strconv.Itoa(nodeIdx),
			})

			nodeWithLabels.Labels = append(nodeWithLabels.Labels, &cre.Label{
				Key:   node.HostLabelKey,
				Value: internalHost,
			})

			donsWithMetadata[donIdx].NodesMetadata[nodeIdx] = &nodeWithLabels
		}
	}

	maybeID, err := flags.OneDonMetadataWithFlag(donsWithMetadata, cre.WorkflowDON)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get workflow DON ID")
	}

	topology.DonsMetadata = donsWithMetadata
	topology.WorkflowDONID = maybeID.ID
	topology.HomeChainSelector = homeChainSelector

	return topology, nil
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
func LinkToJobDistributor(ctx context.Context, input *cre.LinkDonsToJDInput) (*cldf.Environment, []*cre.DON, error) {
	if input == nil {
		return nil, nil, errors.New("input is nil")
	}
	if err := input.Validate(); err != nil {
		return nil, nil, errors.Wrap(err, "input validation failed")
	}

	dons := make([]*cre.DON, len(input.NodeSetOutput))

	jdConfig := jd.JDConfig{
		GRPC:  input.JdOutput.ExternalGRPCUrl,
		WSRPC: input.JdOutput.InternalWSRPCUrl,
		Creds: insecure.NewCredentials(),
	}

	jdClient, jdErr := jd.NewJDClient(jdConfig)
	if jdErr != nil {
		return nil, nil, errors.Wrap(jdErr, "failed to create JD client")
	}

	for idx, nodeOutput := range input.NodeSetOutput {
		supportedChainSelectors := make([]uint64, 0)
		for _, bcOut := range input.BlockchainOutputs {
			if len(input.Topology.DonsMetadata[idx].SupportedChains) > 0 && !slices.Contains(input.Topology.DonsMetadata[idx].SupportedChains, bcOut.ChainID) {
				continue
			}

			supportedChainSelectors = append(supportedChainSelectors, bcOut.ChainSelector)
		}

		don, donErr := NewDON(ctx, nodeOutput.CLNodes, nodeOutput.Capabilities, input.Topology.DonsMetadata[idx], supportedChainSelectors)
		if donErr != nil {
			return nil, nil, fmt.Errorf("failed to create registered DON: %w", donErr)
		}

		var regErr error
		dons[idx], regErr = configureJD(ctx, don, jdClient)
		if regErr != nil {
			return nil, nil, fmt.Errorf("failed to configure JD for DON: %w", regErr)
		}
	}

	var nodeIDs []string
	for _, don := range dons {
		nodeIDs = append(nodeIDs, don.NodeIds()...)
	}

	dons = addOCRKeyLabelsToNodeMetadata(dons, input.Topology)

	input.CldfEnvironment.Offchain = jdClient
	input.CldfEnvironment.NodeIDs = nodeIDs

	return input.CldfEnvironment, dons, nil
}

func NewDON(ctx context.Context, clNodes []*clnode.Output, capabilities []string, donMetadata *cre.DonMetadata, supportedChainSelectors []uint64) (*cre.DON, error) {
	don := &cre.DON{
		Nodes: make([]cre.Node, 0),
	}

	for idx, clNode := range clNodes {
		nodeRoles := make([]string, 0)

		nodeType, typeErr := node.FindLabelValue(donMetadata.NodesMetadata[idx], node.NodeTypeKey)
		if typeErr != nil {
			return nil, errors.Wrapf(typeErr, "failed to find node type for node index %d in DON %s", idx, donMetadata.Name)
		}
		nodeRoles = append(nodeRoles, nodeType)

		if node.HasLabel(donMetadata.NodesMetadata[idx], node.ExtraRolesKey) {
			role, rErr := node.FindLabelValue(donMetadata.NodesMetadata[idx], node.ExtraRolesKey)
			if rErr != nil {
				return nil, errors.Wrapf(rErr, "failed to find extra role for node index %d in DON %s", idx, donMetadata.Name)
			}
			nodeRoles = append(nodeRoles, role)
		}

		node, err := node.NewNode(ctx, clNode, capabilities, nodeRoles, supportedChainSelectors, donMetadata.Name, donMetadata.ID, idx)
		if err != nil {
			return nil, fmt.Errorf("failed to create node %d: %w", idx, err)
		}

		don.Nodes = append(don.Nodes, *node)
	}
	return don, nil
}

func CreateSupportedJobDistributorChains(ctx context.Context, don *cre.DON, jd offchain.Client) error {
	g := new(errgroup.Group)
	for i := range don.Nodes {
		i := i
		g.Go(func() error {
			n := &don.Nodes[i]
			if err := node.CreateJDChainConfig(ctx, n, jd); err != nil {
				return err
			}
			don.Nodes[i] = *n
			return nil
		})
	}
	return g.Wait()
}

func configureJD(ctx context.Context, don *cre.DON, jdClient *jd.JobDistributor) (*cre.DON, error) {
	// todo parallelize each node
	for idx, n := range don.Nodes {
		for _, role := range n.Roles {
			switch role {
			case cre.BootstrapNode, cre.WorkerNode:
				updatedNode, linkErr := node.SetUpAndLinkJobDistributor(ctx, n, jdClient)
				if linkErr != nil {
					return nil, fmt.Errorf("failed to set up job distributor in node %s: %w", n.Name, linkErr)
				}

				if err := node.CreateJDChainConfig(ctx, updatedNode, jdClient); err != nil {
					return nil, err
				}
				don.Nodes[idx] = *updatedNode
			case cre.GatewayNode:
				// nothing to do for gateway nodes
			default:
				return nil, fmt.Errorf("unknown node role: %s", role)
			}
		}
	}

	return don, nil
}

func FindSupportedChainSelectors(donMetadata *cre.DonMetadata, blockchainOutputs []*cre.WrappedBlockchainOutput) ([]uint64, error) {
	selectors := make([]uint64, 0)

	for _, bcOut := range blockchainOutputs {
		if len(donMetadata.SupportedChains) > 0 && !slices.Contains(donMetadata.SupportedChains, bcOut.ChainID) {
			continue
		}

		selectors = append(selectors, bcOut.ChainSelector)
	}

	return selectors, nil
}

// func findSupportedChainsForDON(donMetadata *cre.DonMetadata, blockchainOutputs []*cre.WrappedBlockchainOutput) ([]devenv.ChainConfig, error) {
// 	chains := make([]devenv.ChainConfig, 0)

// 	for chainSelector, bcOut := range blockchainOutputs {
// 		if len(donMetadata.SupportedChains) > 0 && !slices.Contains(donMetadata.SupportedChains, bcOut.ChainID) {
// 			continue
// 		}

// 		cfg, cfgErr := cre.ChainConfigFromWrapped(bcOut)
// 		if cfgErr != nil {
// 			return nil, errors.Wrapf(cfgErr, "failed to build chain config for chain selector %d", chainSelector)
// 		}

// 		chains = append(chains, cfg)
// 	}

// 	return chains, nil
// }

func addOCRKeyLabelsToNodeMetadata(dons []*cre.DON, topology *cre.Topology) []*cre.DON {
	for i, don := range dons {
		for j, donNode := range topology.DonsMetadata[i].NodesMetadata {
			// required for job proposals, because they need to include the ID of the node in Job Distributor
			donNode.Labels = append(donNode.Labels, &cre.Label{
				Key:   node.NodeIDKey,
				Value: don.NodeIds()[j],
			})

			ocrSupportedFamilies := make([]string, 0)
			for family, key := range don.Nodes[j].ChainsOcr2KeyBundlesID {
				donNode.Labels = append(donNode.Labels, &cre.Label{
					Key:   node.CreateNodeOCR2KeyBundleIDKey(family),
					Value: key,
				})
				ocrSupportedFamilies = append(ocrSupportedFamilies, family)
			}

			donNode.Labels = append(donNode.Labels, &cre.Label{
				Key:   node.NodeOCRFamiliesKey,
				Value: node.CreateNodeOCRFamiliesListValue(ocrSupportedFamilies),
			})
		}
	}

	return dons
}
