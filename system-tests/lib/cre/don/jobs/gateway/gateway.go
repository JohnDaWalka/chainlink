package gateway

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/naoina/toml"
	"github.com/pkg/errors"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

var GatewayJobSpecFactoryFn = func(extraAllowedPorts []int, extraAllowedIPs, extraAllowedIPsCIDR []string) cre.JobSpecFactoryFn {
	return func(input *cre.JobSpecFactoryInput) (cre.DonsToJobSpecs, error) {
		return GenerateJobSpecs(
			input.DonTopology,
			extraAllowedPorts,
			extraAllowedIPs,
			extraAllowedIPsCIDR,
			input.DonTopology.GatewayConnectorOutput,
		)
	}
}

func GenerateJobSpecs(donTopology *cre.DonTopology, extraAllowedPorts []int, extraAllowedIPs, extraAllowedIPsCIDR []string, gatewayConnectorOutput *cre.GatewayConnectorOutput) (cre.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}

	donToJobSpecs := make(cre.DonsToJobSpecs)

	// if we don't have a gateway connector output, we don't need to create any job specs
	if gatewayConnectorOutput == nil {
		return donToJobSpecs, nil
	}

	// First create the DON entries for the gateway job spec.
	dons := []DON{}
	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		workflowNodeSet, err := node.FindManyWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: node.NodeTypeKey, Value: cre.WorkerNode}, node.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		ethAddresses := make([]string, len(workflowNodeSet))
		var ethAddressErr error
		for i, n := range workflowNodeSet {
			ethAddresses[i], ethAddressErr = node.FindLabelValue(n, node.AddressKeyFromSelector(donTopology.HomeChainSelector))
			if ethAddressErr != nil {
				return nil, errors.Wrap(ethAddressErr, "failed to get eth address from labels")
			}
		}

		members := []DONMember{}
		for nodeIdx, ethAddress := range ethAddresses {
			members = append(members, DONMember{
				Address: ethAddress,
				Name:    fmt.Sprintf("DON %d - Node %d", donWithMetadata.ID, nodeIdx),
			})
		}

		// If it's a workflow DON or needs custom compute or access to web capabilities, add a DON entry
		// with the web-api-capabilities handler.
		if flags.HasFlag(donWithMetadata.Flags, cre.WorkflowDON) || don.NodeNeedsGateway(donWithMetadata.Flags) {
			don := DON{
				DonId:       fmt.Sprintf("%d", donWithMetadata.ID),
				F:           1,
				HandlerName: "web-api-capabilities",
				Members:     members,
				HandlerConfig: HandlerConfig{
					MaxAllowedMessageAgeSec: 1000,
					NodeRateLimiter: NodeRateLimiterConfig{
						GlobalBurst:    10,
						GlobalRPS:      50,
						PerSenderBurst: 10,
						PerSenderRPS:   10,
					},
				},
			}
			dons = append(dons, don)
		}

		// Otherwise, check if it's a vault capability DON; if so we add a vault handler.
		if flags.HasFlag(donWithMetadata.Flags, cre.VaultCapability) {
			don := DON{
				DonId:       "vault",
				F:           1,
				HandlerName: "vault",
				Members:     members,
				HandlerConfig: HandlerConfig{
					MaxAllowedMessageAgeSec: 1000,
					NodeRateLimiter: NodeRateLimiterConfig{
						GlobalBurst:    10,
						GlobalRPS:      50,
						PerSenderBurst: 10,
						PerSenderRPS:   10,
					},
				},
			}
			dons = append(dons, don)
		}
	}

	if len(dons) == 0 {
		return donToJobSpecs, nil
	}

	// For each gateway node, let's create a job spec.
	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		// create job specs for the gateway node
		if !flags.HasFlag(donWithMetadata.Flags, cre.GatewayDON) {
			continue
		}

		gatewayNode, nodeErr := node.FindOneWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: node.ExtraRolesKey, Value: cre.GatewayNode}, node.LabelContains)
		if nodeErr != nil {
			return nil, errors.Wrap(nodeErr, "failed to find bootstrap node")
		}

		gatewayNodeID, gatewayErr := node.FindLabelValue(gatewayNode, node.NodeIDKey)
		if gatewayErr != nil {
			return nil, errors.Wrap(gatewayErr, "failed to get gateway node id from labels")
		}

		ports := []int{80, 443}
		if len(extraAllowedPorts) != 0 {
			ports = append(ports, extraAllowedPorts...)
		}

		httpConfig := HTTPClientConfig{
			MaxResponseBytes: 100000000,
			AllowedPorts:     ports,
			AllowedIps:       extraAllowedIPs,
			AllowedIPsCIDR:   extraAllowedIPsCIDR,
		}

		jobSpec := GatewayJobSpec{
			Type:              "gateway",
			SchemaVersion:     1,
			ExternalJobID:     uuid.New().String(),
			Name:              cre.GatewayJobName,
			ForwardingAllowed: false,
			ConnectionManagerConfig: ConnectionManagerConfig{
				AuthChallengeLen:          10,
				AuthGatewayId:             "por_gateway",
				AuthTimestampToleranceSec: 5,
				HeartbeatIntervalSec:      20,
			},
			DONs: dons,
			NodeServerConfig: NodeServerConfig{
				HandshakeTimeoutMillis: 1000,
				MaxRequestBytes:        100000,
				Path:                   "", // TODO
				Port:                   0,  // TODO
				ReadTimeoutMillis:      1000,
				RequestTimeoutMillis:   10000,
				WriteTimeoutMillis:     1000,
			},
			UserServerConfig: UserServerConfig{
				ContentTypeHeader:    "application/jsonrpc",
				MaxRequestBytes:      100000,
				Path:                 "", // TODO
				Port:                 0,  // TODO
				ReadTimeoutMillis:    1000,
				RequestTimeoutMillis: 10000,
				WriteTimeoutMillis:   1000,
				CORSEnabled:          false,
				CORSAllowedOrigins:   []string{},
			},
			HTTPClientConfig: httpConfig,
		}

		tomlSpec, err := toml.Marshal(&jobSpec)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal job spec to TOML")
		}

		donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], &jobv1.ProposeJobRequest{
			NodeId: gatewayNodeID,
			Spec:   tomlSpec,
		})
	}

	return donToJobSpecs, nil
}
