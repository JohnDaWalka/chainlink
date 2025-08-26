package environment

import (
	"strings"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/nix"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

// logNodeSetConfiguration logs detailed configuration information for debugging
func logNodeSetConfiguration(lggr zerolog.Logger, nodeSetInput *cretypes.CapabilitiesAwareNodeSet) {
	lggr.Info().Msgf("=== Node Set Configuration: %s ===", nodeSetInput.Name)
	lggr.Info().Msgf("Capabilities: %v", nodeSetInput.Capabilities)
	lggr.Info().Msgf("DON Types: %v", nodeSetInput.DONTypes)
	lggr.Info().Msgf("Bootstrap Node Index: %d", nodeSetInput.BootstrapNodeIndex)
	lggr.Info().Msgf("Gateway Node Index: %d", nodeSetInput.GatewayNodeIndex)
	lggr.Info().Msgf("Number of nodes: %d", len(nodeSetInput.NodeSpecs))

	// Log the raw input configuration
	if nodeSetInput.Input != nil {
		lggr.Info().Msgf("Input configuration type: %T", nodeSetInput.Input)
		// Log any available fields from the input
		if nodeSetInput.Input.Name != "" {
			lggr.Info().Msgf("Input name: %s", nodeSetInput.Input.Name)
		}
	}

	lggr.Info().Msgf("=== End Node Set Configuration: %s ===", nodeSetInput.Name)
}

func StartDONs(
	lggr zerolog.Logger,
	nixShell *nix.Shell,
	topology *cretypes.Topology,
	infraType libtypes.InfraType,
	registryChainBlockchainOutput *blockchain.Output,
	capabilitiesAwareNodeSets []*cretypes.CapabilitiesAwareNodeSet,
) ([]*cretypes.WrappedNodeOutput, error) {
	startTime := time.Now()
	lggr.Info().Msgf("Starting %d DONs", len(capabilitiesAwareNodeSets))
	lggr.Info().Msgf("StartDONs: Infrastructure type: %s", infraType)
	lggr.Info().Msgf("StartDONs: Topology workflow DON ID: %d", topology.WorkflowDONID)

	if infraType == libtypes.CRIB {
		lggr.Info().Msg("StartDONs: Using CRIB infrastructure, saving node configs and secret overrides")
		deployCribDonsInput := &cretypes.DeployCribDonsInput{
			Topology:       topology,
			NodeSetInputs:  capabilitiesAwareNodeSets,
			NixShell:       nixShell,
			CribConfigsDir: cribConfigsDir,
		}

		var devspaceErr error
		capabilitiesAwareNodeSets, devspaceErr = crib.DeployDons(deployCribDonsInput)
		if devspaceErr != nil {
			lggr.Error().Err(devspaceErr).Msg("StartDONs: Failed to deploy Dons with devspace")
			return nil, pkgerrors.Wrap(devspaceErr, "failed to deploy Dons with devspace")
		}
		lggr.Info().Msg("StartDONs: Successfully deployed Dons with devspace")
	}

	nodeSetOutput := make([]*cretypes.WrappedNodeOutput, 0, len(capabilitiesAwareNodeSets))

	// TODO we could parallelize this as well in the future, but for single DON env this doesn't matter
	for i, nodeSetInput := range capabilitiesAwareNodeSets {
		lggr.Info().Msgf("StartDONs: Creating node set %d/%d: %s", i+1, len(capabilitiesAwareNodeSets), nodeSetInput.Name)

		// Log detailed configuration for debugging
		logNodeSetConfiguration(lggr, nodeSetInput)

		lggr.Info().Msgf("StartDONs: Node set %s has %d nodes", nodeSetInput.Name, len(nodeSetInput.NodeSpecs))
		lggr.Info().Msgf("StartDONs: Node set %s capabilities: %v", nodeSetInput.Name, nodeSetInput.Capabilities)
		lggr.Info().Msgf("StartDONs: Node set %s DON types: %v", nodeSetInput.Name, nodeSetInput.DONTypes)

		// Log node details for debugging
		for j, nodeSpec := range nodeSetInput.NodeSpecs {
			lggr.Info().Msgf("StartDONs: Node %d/%d in set %s: Image=%s, Dockerfile=%s, Context=%s",
				j+1, len(nodeSetInput.NodeSpecs), nodeSetInput.Name,
				nodeSpec.Node.Image, nodeSpec.Node.DockerFilePath, nodeSpec.Node.DockerContext)
		}

		lggr.Info().Msgf("StartDONs: About to call NewSharedDBNodeSet for %s", nodeSetInput.Name)
		nodeset, nodesetErr := ns.NewSharedDBNodeSet(nodeSetInput.Input, registryChainBlockchainOutput)
		if nodesetErr != nil {
			lggr.Error().Err(nodesetErr).Msgf("StartDONs: Failed to create node set named %s", nodeSetInput.Name)

			// Try to get more details about the error
			errorStr := nodesetErr.Error()
			lggr.Error().Msgf("StartDONs: Error details: %s", errorStr)

			// Check if it's a container startup error
			if strings.Contains(errorStr, "container exited with code 1") {
				lggr.Error().Msg("StartDONs: Container startup failed - this usually indicates a configuration or resource issue")
				lggr.Error().Msg("StartDONs: Check container logs and ensure all required resources are available")
			}

			// Check if it's a Docker image issue
			if strings.Contains(errorStr, "pull access denied") || strings.Contains(errorStr, "not found") {
				lggr.Error().Msg("StartDONs: Docker image issue detected - check if the image exists and is accessible")
			}

			return nil, pkgerrors.Wrapf(nodesetErr, "failed to create node set named %s", nodeSetInput.Name)
		}
		lggr.Info().Msgf("StartDONs: Successfully created node set %s", nodeSetInput.Name)

		nodeSetOutput = append(nodeSetOutput, &cretypes.WrappedNodeOutput{
			Output:       nodeset,
			NodeSetName:  nodeSetInput.Name,
			Capabilities: nodeSetInput.Capabilities,
		})
	}

	lggr.Info().Msgf("DONs started in %.2f seconds", time.Since(startTime).Seconds())
	lggr.Info().Msgf("StartDONs: Successfully created %d node set outputs", len(nodeSetOutput))

	return nodeSetOutput, nil
}
