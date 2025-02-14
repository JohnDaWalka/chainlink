package node

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"

	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

type StringTransformer func(string) string

func NoOpTransformFn(value string) string {
	return value
}

func KeyExtractingTransformFn(value string) string {
	parts := strings.Split(value, "_")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return value
}

func ToP2PID(node devenv.Node, transformFn StringTransformer) (string, error) {
	for _, label := range node.Labels() {
		if label.Key == devenv.NodeLabelP2PIDType {
			if label.Value == nil {
				return "", fmt.Errorf("p2p label value is nil for node %s", node.Name)
			}
			return transformFn(*label.Value), nil
		}
	}

	return "", fmt.Errorf("p2p label not found for node %s", node.Name)
}

// copied from Bala's unmerged PR: https://github.com/smartcontractkit/chainlink/pull/15751
// TODO: remove this once the PR is merged and import his function
// IMPORTANT ADDITION:  prefix to differentiate between the different DONs
func GetNodeInfo(nodeOut *ns.Output, prefix string, bootstrapNodeCount int) ([]devenv.NodeInfo, error) {
	var nodeInfo []devenv.NodeInfo
	for i := 1; i <= len(nodeOut.CLNodes); i++ {
		p2pURL, err := url.Parse(nodeOut.CLNodes[i-1].Node.DockerP2PUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to parse p2p url: %w", err)
		}
		if i <= bootstrapNodeCount {
			nodeInfo = append(nodeInfo, devenv.NodeInfo{
				IsBootstrap: true,
				Name:        fmt.Sprintf("%s_bootstrap-%d", prefix, i),
				P2PPort:     p2pURL.Port(),
				CLConfig: nodeclient.ChainlinkConfig{
					URL:        nodeOut.CLNodes[i-1].Node.HostURL,
					Email:      nodeOut.CLNodes[i-1].Node.APIAuthUser,
					Password:   nodeOut.CLNodes[i-1].Node.APIAuthPassword,
					InternalIP: nodeOut.CLNodes[i-1].Node.InternalIP,
				},
			})
		} else {
			nodeInfo = append(nodeInfo, devenv.NodeInfo{
				IsBootstrap: false,
				Name:        fmt.Sprintf("%s_node-%d", prefix, i),
				P2PPort:     p2pURL.Port(),
				CLConfig: nodeclient.ChainlinkConfig{
					URL:        nodeOut.CLNodes[i-1].Node.HostURL,
					Email:      nodeOut.CLNodes[i-1].Node.APIAuthUser,
					Password:   nodeOut.CLNodes[i-1].Node.APIAuthPassword,
					InternalIP: nodeOut.CLNodes[i-1].Node.InternalIP,
				},
			})
		}
	}
	return nodeInfo, nil
}
