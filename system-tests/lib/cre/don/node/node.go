package node

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"

	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

const (
	HostLabelKey            = "host"
	IndexKey                = "node_index"
	EthAddressKey           = "eth_address"
	ExtraRolesKey           = "extra_roles"
	NodeIDKeyType           = "node_id"
	NodeOCR2KeyBundleIDType = "ocr2_key_bundle_id"
)

type stringTransformer func(string) string

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

func ToP2PID(node *types.NodeMetadata, transformFn stringTransformer) (string, error) {
	for _, label := range node.Labels {
		if label.Key == devenv.NodeLabelP2PIDType {
			if label.Value == nil {
				return "", fmt.Errorf("p2p label value is nil for node %s", node)
			}
			return transformFn(*label.Value), nil
		}
	}

	return "", fmt.Errorf("p2p label not found for node %s", node)
}

// copied from Bala's unmerged PR: https://github.com/smartcontractkit/chainlink/pull/15751
// TODO: remove this once the PR is merged and import his function
// IMPORTANT ADDITION: prefix to differentiate between the different DONs
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
				Labels: map[string]string{
					HostLabelKey:            nodeOut.CLNodes[i-1].Node.ContainerName,
					IndexKey:                strconv.Itoa(i - 1),
					devenv.NodeLabelKeyType: devenv.NodeLabelValueBootstrap,
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
				Labels: map[string]string{
					HostLabelKey:            nodeOut.CLNodes[i-1].Node.ContainerName,
					IndexKey:                strconv.Itoa(i - 1),
					devenv.NodeLabelKeyType: devenv.NodeLabelValuePlugin,
				},
			})
		}
	}
	return nodeInfo, nil
}

func FindOneWithLabel(nodes []*types.NodeMetadata, wantedLabel *ptypes.Label, labelMatcherFn labelMatcherFn) (*types.NodeMetadata, error) {
	if wantedLabel == nil {
		return nil, errors.New("label is nil")
	}
	for _, node := range nodes {
		for _, label := range node.Labels {
			if wantedLabel.Key == label.Key && labelMatcherFn(wantedLabel.Value, label.Value) {
				return node, nil
			}
		}
	}
	return nil, fmt.Errorf("node with label %s=%s not found", wantedLabel.Key, *wantedLabel.Value)
}

func FindManyWithLabel(nodes []*types.NodeMetadata, wantedLabel *ptypes.Label, labelMatcherFn labelMatcherFn) ([]*types.NodeMetadata, error) {
	if wantedLabel == nil {
		return nil, errors.New("label is nil")
	}

	var foundNodes []*types.NodeMetadata

	for _, node := range nodes {
		for _, label := range node.Labels {
			if wantedLabel.Key == label.Key && labelMatcherFn(wantedLabel.Value, label.Value) {
				foundNodes = append(foundNodes, node)
			}
		}
	}

	if len(foundNodes) == 0 {
		return nil, fmt.Errorf("node with label %s=%s not found", wantedLabel.Key, *wantedLabel.Value)
	}

	return foundNodes, nil
}

func FindLabelValue(node *types.NodeMetadata, labelKey string) (string, error) {
	for _, label := range node.Labels {
		if label.Key == labelKey {
			if label.Value == nil {
				return "", fmt.Errorf("label %s found, but its value is nil", labelKey)
			}
			if *label.Value == "" {
				return "", fmt.Errorf("label %s found, but its value is empty", labelKey)
			}
			return *label.Value, nil
		}
	}

	return "", fmt.Errorf("label %s not found", labelKey)
}

type labelMatcherFn func(first, second *string) bool

func EqualLabels(first, second *string) bool {
	if first == nil && second == nil {
		return true
	}
	if first == nil || second == nil {
		return false
	}
	return *first == *second
}

func LabelContains(first, second *string) bool {
	if first == nil && second == nil {
		return true
	}
	if first == nil || second == nil {
		return false
	}

	return strings.Contains(*first, *second)
}
