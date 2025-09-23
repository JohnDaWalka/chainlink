package node

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/sethvargo/go-retry"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	"github.com/smartcontractkit/chainlink/deployment/environment/web/sdk/client"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
)

const (
	NodeTypeKey            = "type"
	HostLabelKey           = "host"
	IndexKey               = "node_index"
	ExtraRolesKey          = "extra_roles"
	NodeIDKey              = "node_id"
	NodeOCRFamiliesKey     = "node_ocr_families"
	NodeOCR2KeyBundleIDKey = "ocr2_key_bundle_id"
	NodeP2PIDKey           = "p2p_id"
	NodeDKGRecipientKey    = "dkg_recipient_key"
	DONIDKey               = "don_id"
	EnvironmentKey         = "environment"
	ProductKey             = "product"
	DONNameKey             = "don_name"
)

// ocr2 keys depend on report's target chain family
func CreateNodeOCR2KeyBundleIDKey(chainFamily string) string {
	return NodeOCR2KeyBundleIDKey + "_" + chainFamily
}

func CreateNodeOCRFamiliesListValue(families []string) string {
	return strings.Join(families, ",")
}

func AddressKeyFromSelector(chainSelector uint64) string {
	return strconv.FormatUint(chainSelector, 10) + "_public_address"
}

func ExtractBundleKeysPerFamily(n *cre.NodeMetadata) (map[string]string, error) {
	keyBundlesFamilies, fErr := FindLabelValue(n, NodeOCRFamiliesKey)
	if fErr != nil {
		return nil, fmt.Errorf("failed to get ocr families bundle id from worker node labels: %w", fErr)
	}

	supportedFamilies := strings.Split(keyBundlesFamilies, ",")

	bundlesPerFamily := make(map[string]string)
	for _, family := range supportedFamilies {
		kBundle, kbErr := FindLabelValue(n, CreateNodeOCR2KeyBundleIDKey(family))
		if kbErr != nil {
			return nil, fmt.Errorf("failed to get ocr bundle id from worker node labels for family %s err: %w", family, kbErr)
		}
		bundlesPerFamily[family] = kBundle
	}

	return bundlesPerFamily, nil
}

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

func ToP2PID(node *cre.NodeMetadata, transformFn stringTransformer) (string, error) {
	for _, label := range node.Labels {
		if label.Key == NodeP2PIDKey {
			if label.Value == "" {
				return "", errors.New("p2p label value is empty for node")
			}
			return transformFn(label.Value), nil
		}
	}

	return "", errors.New("p2p label not found for node")
}

// copied from Bala's unmerged PR: https://github.com/smartcontractkit/chainlink/pull/15751
// TODO: remove this once the PR is merged and import his function
// IMPORTANT ADDITION: prefix to differentiate between the different DONs
func GetNodeInfo(nodeOut *ns.Output, prefix string, donID uint64, bootstrapNodeCount int) ([]devenv.NodeInfo, error) {
	var nodeInfo []devenv.NodeInfo
	for i := 1; i <= len(nodeOut.CLNodes); i++ {
		p2pURL, err := url.Parse(nodeOut.CLNodes[i-1].Node.InternalP2PUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to parse p2p url: %w", err)
		}

		info := devenv.NodeInfo{
			P2PPort: p2pURL.Port(),
			CLConfig: nodeclient.ChainlinkConfig{
				URL:        nodeOut.CLNodes[i-1].Node.ExternalURL,
				Email:      nodeOut.CLNodes[i-1].Node.APIAuthUser,
				Password:   nodeOut.CLNodes[i-1].Node.APIAuthPassword,
				InternalIP: nodeOut.CLNodes[i-1].Node.InternalIP,
			},
			Labels: map[string]string{
				"don-" + prefix: "true",
				ProductKey:      "keystone",
				EnvironmentKey:  "local",
				DONIDKey:        strconv.FormatUint(donID, 10),
				DONNameKey:      prefix,
			},
		}

		if i <= bootstrapNodeCount {
			info.IsBootstrap = true
			info.Name = fmt.Sprintf("%s_bootstrap-%d", prefix, i)
			info.Labels[NodeTypeKey] = cre.BootstrapNode
		} else {
			info.IsBootstrap = false
			info.Name = fmt.Sprintf("%s_node-%d", prefix, i)
			info.Labels[NodeTypeKey] = cre.WorkerNode
		}

		nodeInfo = append(nodeInfo, info)
	}
	return nodeInfo, nil
}

func FindOneWithLabel(nodes []*cre.NodeMetadata, wantedLabel *cre.Label, labelMatcherFn labelMatcherFn) (*cre.NodeMetadata, error) {
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
	return nil, fmt.Errorf("node with label %s=%s not found", wantedLabel.Key, wantedLabel.Value)
}

func FindManyWithLabel(nodes []*cre.NodeMetadata, wantedLabel *cre.Label, labelMatcherFn labelMatcherFn) ([]*cre.NodeMetadata, error) {
	if wantedLabel == nil {
		return nil, errors.New("label is nil")
	}

	var foundNodes []*cre.NodeMetadata

	for _, node := range nodes {
		for _, label := range node.Labels {
			if wantedLabel.Key == label.Key && labelMatcherFn(wantedLabel.Value, label.Value) {
				foundNodes = append(foundNodes, node)
			}
		}
	}

	return foundNodes, nil
}

func HasLabel(node *cre.NodeMetadata, labelKey string) bool {
	for _, label := range node.Labels {
		if label.Key == labelKey {
			return true
		}
	}
	return false
}

func FindLabelValue(node *cre.NodeMetadata, labelKey string) (string, error) {
	for _, label := range node.Labels {
		if label.Key == labelKey {
			if label.Value == "" {
				return "", fmt.Errorf("label %s found, but its value is empty", labelKey)
			}
			return label.Value, nil
		}
	}

	return "", fmt.Errorf("label %s not found", labelKey)
}

type labelMatcherFn func(first, second string) bool

func EqualLabels(first, second string) bool {
	return first == second
}

func LabelContains(first, second string) bool {
	return strings.Contains(first, second)
}

type SupportedChains struct {
	Family string
	ID     uint64
}

func NewNode(ctx context.Context, clNode *clnode.Output, capablities, roles []string, supportedChainSelectors []uint64, donName string, donID uint64, nodeIndex int) (*cre.Node, error) {
	gqlClient, err := client.NewWithContext(ctx, clNode.Node.ExternalURL, client.Credentials{
		Email:    clNode.Node.APIAuthUser,
		Password: clNode.Node.APIAuthPassword,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create node graphql client: %w", err)
	}
	chainlinkClient, err := clclient.NewChainlinkClient(&clclient.Config{
		URL:         clNode.Node.ExternalURL,
		Email:       clNode.Node.APIAuthUser,
		Password:    clNode.Node.APIAuthPassword,
		InternalIP:  clNode.Node.InternalIP,
		HTTPTimeout: ptr.Ptr(10 * time.Second),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create node rest client: %w", err)
	}

	// node Labels so that it's easier to query them
	jdLabels := []*ptypes.Label{
		{Key: "don-" + donName, Value: ptr.Ptr("true")},
		{Key: ProductKey, Value: ptr.Ptr("keystone")},
		{Key: EnvironmentKey, Value: ptr.Ptr("local")},
		{Key: DONIDKey, Value: ptr.Ptr(strconv.FormatUint(donID, 10))},
		{Key: DONNameKey, Value: ptr.Ptr(donName)},
	}

	node := &cre.Node{
		GQLClient:               gqlClient,
		RestClient:              chainlinkClient,
		Name:                    donName + "-node-" + strconv.Itoa(nodeIndex),
		JDLabels:                jdLabels,
		ChainsOcr2KeyBundlesID:  make(map[string]string),
		Capabilities:            capablities, // TODO this should actually be on DON-level, although not necessarily as bootstrap won't have any capabilities
		Roles:                   roles,
		SupportedChainSelectors: supportedChainSelectors, // TODO this should actually be on DON-level (probably)
	}
	if slices.Contains(roles, cre.BootstrapNode) {
		// create multi address for OCR2, applicable only for bootstrap nodes
		p2pURL, err := url.Parse(clNode.Node.InternalP2PUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to parse p2p url: %w", err)
		}
		node.MultiAddr = fmt.Sprintf("%s:%s", clNode.Node.InternalIP, p2pURL.Port())

		// no need to set admin address for bootstrap nodes, as there will be no payment
		node.AdminAddr = ""
		node.JDLabels = append(node.JDLabels, &ptypes.Label{
			Key:   devenv.LabelNodeTypeKey,
			Value: ptr.Ptr(devenv.LabelNodeTypeValueBootstrap),
		})
	} else {
		// multi address is not applicable for non-bootstrap nodes
		// explicitly set it to empty string to denote that
		node.MultiAddr = ""

		// set admin address for non-bootstrap nodes (capability registry requires non-null admin address; use arbitrary default value if node is not configured)
		node.AdminAddr = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"

		node.JDLabels = append(node.JDLabels, &ptypes.Label{
			Key:   devenv.LabelNodeTypeKey,
			Value: ptr.Ptr(devenv.LabelNodeTypeValuePlugin),
		})
	}

	node, fetchErr := FetchData(ctx, node)
	if fetchErr != nil {
		return nil, fmt.Errorf("failed to fetch node data: %w", fetchErr)
	}

	return node, nil
}

const LabelCSAKey = "csa_key"

// type Node struct {
// 	PeerID                 string            // p2p peer id fetched from the node
// 	NodeID                 string            // node id returned by job distributor after node is registered with it
// 	JDId                   string            // job distributor id returned by node after Job distributor is created in node
// 	Name                   string            // name of the node
// 	AccountAddr            map[string]string // chain id to node's account address mapping for supported chains
// 	ChainsOcr2KeyBundlesID map[string]string
// 	GQLClient              client.Client             // graphql client to interact with the node
// 	RestClient             *clclient.ChainlinkClient // rest client to interact with the node
// 	JDLabels               []*ptypes.Label           // labels with which the node is registered with the job distributor
// 	AdminAddr              string                    // admin address to send payments to, applicable only for non-bootstrap nodes
// 	MultiAddr              string                    // multi address denoting node's FQN (needed for deriving P2PBootstrappers in OCR), applicable only for bootstrap nodes
// 	CSAKey                 string                    // csa public key of the node

// 	Capabilities            []string
// 	Roles                   []string
// 	SupportedChainSelectors []uint64
// }

func FetchData(ctx context.Context, n *cre.Node) (*cre.Node, error) {
	for _, selector := range n.SupportedChainSelectors {
		if n.AccountAddr == nil {
			n.AccountAddr = make(map[string]string)
		}

		chainID, chainFamily, chErr := selectorToIDAndFamily(selector)
		if chErr != nil {
			return nil, fmt.Errorf("failed to get chain id and chain family from selector %d: %w", selector, chErr)
		}

		// TODO: maybe this should be done by write-evm/evm and write-solana capabilities?
		switch chainFamily {
		case chainselectors.FamilyEVM:
			accountAddr, accErr := n.GQLClient.FetchAccountAddress(ctx, chainID)
			if accErr != nil {
				return nil, fmt.Errorf("failed to fetch account address for node %s: %w", n.Name, accErr)
			}
			if accountAddr == nil {
				return nil, fmt.Errorf("no account address found for node %s", n.Name)
			}
			n.AccountAddr[chainID] = *accountAddr
		case chainselectors.FamilySolana:
			accounts, err := n.GQLClient.FetchKeys(ctx, chainFamily)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch account address for node %s and chain %s: %w", n.Name, chainFamily, err)
			}
			if len(accounts) == 0 {
				return nil, fmt.Errorf("failed to fetch account address for node %s and chain %s", n.Name, chainFamily)
			}
			n.AccountAddr[chainID] = accounts[0]
		default:
			return nil, fmt.Errorf("unsupported chain family %s", chainFamily)
		}

		// TODO: Maybe this should be part of the capability code? I.e. check if the node has these key bundles set and if not fetch them?
		// fetch OCR2 key bundle id only if the node has any capabilities that require OCR2
		// TODO: use don.RequiresOCR()
		if flags.HasFlagForAnyChain(n.Capabilities, cre.ConsensusCapability) || flags.HasFlagForAnyChain(n.Capabilities, cre.ConsensusCapabilityV2) || flags.HasFlagForAnyChain(n.Capabilities, cre.VaultCapability) || flags.HasFlagForAnyChain(n.Capabilities, cre.EVMCapability) {
			ocr2BundleId, err := n.GQLClient.FetchOCR2KeyBundleID(ctx, strings.ToUpper(chainFamily))
			if err != nil {
				return nil, fmt.Errorf("failed to fetch OCR2 key bundle id for node %s: %w", n.Name, err)
			}
			if ocr2BundleId == "" {
				return nil, fmt.Errorf("no OCR2 key bundle id found for node %s", n.Name)
			}

			n.ChainsOcr2KeyBundlesID[chainFamily] = ocr2BundleId
		}
	}

	// TODO: part of capability or role?
	// if the node is only a gateway and nothing else, then skip fetching peer id, since it doesn't have p2p
	if (slices.Contains(n.Roles, cre.GatewayNode) && len(n.Roles) > 1) || !slices.Contains(n.Roles, cre.GatewayNode) {
		peerID, err := n.GQLClient.FetchP2PPeerID(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch peer id for node %s: %w", n.Name, err)
		}

		if peerID == nil {
			return nil, fmt.Errorf("no peer id found for node %s", n.Name)
		}

		n.PeerID = *peerID
	}

	// Get the public CSA key of the node (every node has to have one)
	csaKeyRes, err := n.GQLClient.FetchCSAPublicKey(ctx)
	if err != nil {
		return nil, err
	}
	if csaKeyRes == nil {
		return nil, fmt.Errorf("no csa key found for node %s", n.Name)
	}
	n.CSAKey = strings.TrimPrefix(*csaKeyRes, "csa_")

	return n, nil
}

func selectorToIDAndFamily(selector uint64) (string, string, error) {
	chainID, chErr := chainselectors.ChainIdFromSelector(selector)
	if chErr != nil {
		return "", "", fmt.Errorf("failed to get chain id from selector %d: %w", selector, chErr)
	}

	chainFamily, familyErr := chainselectors.GetSelectorFamily(selector)
	if familyErr != nil {
		return "", "", fmt.Errorf("failed to get chain family from selector %d: %w", selector, familyErr)
	}

	chainIDStr := strconv.FormatUint(chainID, 10)

	return chainIDStr, chainFamily, nil
}

func CreateJDChainConfig(ctx context.Context, n *cre.Node, jd offchain.Client) error {
	for _, selector := range n.SupportedChainSelectors {
		if n.AccountAddr == nil {
			return fmt.Errorf("node %s account address map is nil", n.Name)
		}

		if n.ChainsOcr2KeyBundlesID == nil {
			return fmt.Errorf("node %s chains ocr2 key bundle id map is nil", n.Name)
		}

		chainID, chainFamily, chErr := selectorToIDAndFamily(selector)
		if chErr != nil {
			return fmt.Errorf("failed to get chain id and chain family from selector %d: %w", selector, chErr)
		}

		address, ok := n.AccountAddr[chainID]
		if !ok {
			return fmt.Errorf("no account address found for node %s and chain id %s", n.Name, chainID)
		}

		// retry twice with 5 seconds interval to create JobDistributorChainConfig
		err := retry.Do(ctx, retry.WithMaxDuration(10*time.Second, retry.NewConstant(3*time.Second)), func(ctx context.Context) error {
			// check the node chain config to see if this chain already exists
			nodeChainConfigs, err := jd.ListNodeChainConfigs(context.Background(), &nodev1.ListNodeChainConfigsRequest{
				Filter: &nodev1.ListNodeChainConfigsRequest_Filter{
					NodeIds: []string{n.NodeID},
				}})
			if err != nil {
				return retry.RetryableError(fmt.Errorf("failed to list node chain configs for node %s, retrying..: %w", n.Name, err))
			}
			if nodeChainConfigs != nil {
				for _, chainConfig := range nodeChainConfigs.ChainConfigs {
					if strings.EqualFold(chainConfig.Chain.Id, chainID) {
						return nil
					}
				}
			}

			// TODO probably capabilities should be able to modify this config according to their needs (like adding OCR2 keys)
			config := client.JobDistributorChainConfigInput{
				JobDistributorID: n.JDId,
				ChainID:          chainID,
				ChainType:        strings.ToUpper(chainFamily),
				AccountAddr:      address,
				AdminAddr:        n.AdminAddr,
			}

			if flags.HasFlagForAnyChain(n.Capabilities, cre.ConsensusCapability) || flags.HasFlagForAnyChain(n.Capabilities, cre.ConsensusCapabilityV2) || flags.HasFlagForAnyChain(n.Capabilities, cre.VaultCapability) || flags.HasFlagForAnyChain(n.Capabilities, cre.EVMCapability) {
				ocr2BundleId, ok := n.ChainsOcr2KeyBundlesID[chainFamily]
				if !ok {
					return fmt.Errorf("no ocr2 key bundle id found for node %s and chain family %s", n.Name, chainFamily)
				}

				config.Ocr2Enabled = true
				config.Ocr2KeyBundleID = ocr2BundleId
				config.Ocr2IsBootstrap = slices.Contains(n.Roles, cre.BootstrapNode)
				config.Ocr2Multiaddr = n.MultiAddr
				config.Ocr2P2PPeerID = n.PeerID
				config.Ocr2Plugins = `{"commit":true,"execute":true,"median":false,"mercury":false}` // TODO: should this be configurable? Do we even need this?
			}

			// JD silently fails to update nodeChainConfig. Therefore, we fetch the node config and
			// if it's not updated , throw an error
			_, err = n.GQLClient.CreateJobDistributorChainConfig(ctx, config)

			// todo: add a check if the chain config failed because of a duplicate in that case, should we update or return success?
			if err != nil {
				return fmt.Errorf("failed to create JD chain config (ID: %s, family: %s) for node %s: %w", chainID, chainFamily, n.Name, err)
			}

			return retry.RetryableError(errors.New("retrying CreateChainConfig in JD"))
		})

		if err != nil {
			return fmt.Errorf("failed to create JD chain config (ID: %s, family: %s) for node %s: %w", chainID, chainFamily, n.Name, err)
		}
	}

	return nil
}

// RegisterNodeToJobDistributor fetches the CSA public key of the node and registers the node with the job distributor
// it sets the node id returned by JobDistributor as a result of registration in the node struct
func RegisterNodeToJobDistributor(ctx context.Context, n cre.Node, jd jd.JobDistributor) (*cre.Node, error) {
	// Get the public key of the node
	// csaKeyRes, err := n.GQLClient.FetchCSAPublicKey(ctx)
	// if err != nil {
	// 	return err
	// }
	// if csaKeyRes == nil {
	// 	return fmt.Errorf("no csa key found for node %s", n.Name)
	// }
	// csaKey := strings.TrimPrefix(*csaKeyRes, "csa_")

	// // tag nodes with p2p_id for easy lookup
	// peerID, err := n.GQLClient.FetchP2PPeerID(ctx)
	// if err != nil {
	// 	return fmt.Errorf("failed to fetch peer id for node %s: %w", n.Name, err)
	// }
	// if peerID == nil {
	// 	return fmt.Errorf("no peer id found for node %s", n.Name)
	// }

	// not every node has p2p, e.g. gateway only nodes
	if n.PeerID != "" {
		n.JDLabels = append(n.JDLabels, &ptypes.Label{
			Key:   devenv.LabelNodeP2PIDKey,
			Value: &n.PeerID,
		})
	}

	// n.JDLabels = append(n.JDLabels, &ptypes.Label{
	// 	Key:   LabelCSAKey,
	// 	Value: &n.CSAKey,
	// })

	// register the node in the job distributor
	registerResponse, err := jd.RegisterNode(ctx, &nodev1.RegisterNodeRequest{
		PublicKey: n.CSAKey,
		Labels:    n.JDLabels,
		Name:      n.Name,
	})

	// node already registered, fetch it's id
	if err != nil && strings.Contains(err.Error(), "AlreadyExists") {
		nodesResponse, err := jd.ListNodes(ctx, &nodev1.ListNodesRequest{
			Filter: &nodev1.ListNodesRequest_Filter{
				// Selectors: []*ptypes.Selector{
				// 	{
				// 		Key:   LabelCSAKey,
				// 		Op:    ptypes.SelectorOp_EQ,
				// 		Value: &n.CSAKey,
				// 	},
				// },
				PublicKeys: []string{n.CSAKey},
			},
		})
		if err != nil {
			return nil, err
		}
		nodes := nodesResponse.GetNodes()
		if len(nodes) == 0 {
			return nil, fmt.Errorf("failed to find node: %v", n.Name)
		}
		n.NodeID = nodes[0].Id
		return &n, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to register node %s: %w", n.Name, err)
	}
	if registerResponse.GetNode().GetId() == "" {
		return nil, fmt.Errorf("no node id returned from job distributor for node %s", n.Name)
	}
	n.NodeID = registerResponse.GetNode().GetId()

	return &n, nil
}

// CreateJobDistributor fetches the keypairs from the job distributor and creates the job distributor in the node
// and returns the job distributor id
func CreateJobDistributor(ctx context.Context, n cre.Node, jd jd.JobDistributor) (string, error) {
	csaKey, csaErr := jd.GetCSAPublicKey(ctx)
	if csaErr != nil {
		return "", fmt.Errorf("failed to get csa public key from job distributor: %w", csaErr)
	}

	// create the job distributor in the node with the csa key
	resp, err := n.GQLClient.ListJobDistributors(ctx)
	if err != nil {
		return "", fmt.Errorf("could not list job distributors: %w", err)
	}

	if len(resp.FeedsManagers.Results) > 0 {
		for _, fm := range resp.FeedsManagers.Results {
			if fm.GetPublicKey() == csaKey {
				return fm.GetId(), nil
			}
		}
	}

	return n.GQLClient.CreateJobDistributor(ctx, client.JobDistributorInput{
		Name:      "Job Distributor",
		Uri:       jd.WSRPC,
		PublicKey: csaKey,
	})
}

// SetUpAndLinkJobDistributor sets up the job distributor in the node and registers the node with the job distributor
// it sets the job distributor id for node
func SetUpAndLinkJobDistributor(ctx context.Context, n cre.Node, jd jd.JobDistributor) (*cre.Node, error) {
	// register the node in the job distributor
	updatedNode, err := RegisterNodeToJobDistributor(ctx, n, jd)
	if err != nil {
		return nil, err
	}

	n = *updatedNode

	// now create the job distributor in the node
	id, err := CreateJobDistributor(ctx, n, jd)
	if err != nil &&
		!strings.Contains(err.Error(), "DuplicateFeedsManagerError") {
		return nil, fmt.Errorf("failed to create job distributor in node %s: %w", n.Name, err)
	}

	// wait for the node to connect to the job distributor
	err = retry.Do(ctx, retry.WithMaxDuration(1*time.Minute, retry.NewFibonacci(1*time.Second)), func(ctx context.Context) error {
		getRes, err := jd.GetNode(ctx, &nodev1.GetNodeRequest{
			Id: n.NodeID,
		})
		if err != nil {
			return retry.RetryableError(fmt.Errorf("failed to get node %s: %w", n.Name, err))
		}
		if getRes.GetNode() == nil {
			return fmt.Errorf("no node found for node id %s", n.NodeID)
		}
		if !getRes.GetNode().IsConnected {
			return retry.RetryableError(fmt.Errorf("node %s not connected to job distributor", n.Name))
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect node %s to job distributor: %w", n.Name, err)
	}

	n.JDId = id

	return &n, nil
}

func CancelProposalsByExternalJobID(ctx context.Context, n cre.Node, externalJobIDs []string) ([]string, error) {
	jd, err := n.GQLClient.GetJobDistributor(ctx, n.JDId)
	if err != nil {
		return nil, err
	}
	if jd.GetJobProposals() == nil {
		return nil, fmt.Errorf("no job proposals found for node %s", n.Name)
	}

	proposalIDs := []string{}
	for _, jp := range jd.JobProposals {
		if !slices.Contains(externalJobIDs, jp.ExternalJobID) {
			continue
		}

		proposalIDs = append(proposalIDs, jp.Id)
		spec, err := n.GQLClient.CancelJobProposalSpec(ctx, jp.Id)
		if err != nil {
			return nil, err
		}

		if spec == nil {
			return nil, fmt.Errorf("no job proposal spec found for id %s", jp.Id)
		}
	}

	return proposalIDs, nil
}

func ApproveProposals(ctx context.Context, n cre.Node, proposalIDs []string) error {
	for _, proposalID := range proposalIDs {
		spec, err := n.GQLClient.ApproveJobProposalSpec(ctx, proposalID, false)
		if err != nil {
			return err
		}
		if spec == nil {
			return fmt.Errorf("no job proposal spec found for id %s", proposalID)
		}
	}
	return nil
}

func AcceptJob(ctx context.Context, n cre.Node, spec string) error {
	// fetch JD to get the job proposals
	jd, err := n.GQLClient.GetJobDistributor(ctx, n.JDId)
	if err != nil {
		return err
	}
	if jd.GetJobProposals() == nil {
		return fmt.Errorf("no job proposals found for node %s", n.Name)
	}
	// locate the job proposal id for the given job spec
	var idToAccept string
	for _, jp := range jd.JobProposals {
		if jp.LatestSpec.Definition == spec {
			idToAccept = jp.Id
			break
		}
	}
	if idToAccept == "" {
		return fmt.Errorf("no job proposal found for job spec %s", spec)
	}
	approvedSpec, err := n.GQLClient.ApproveJobProposalSpec(ctx, idToAccept, false)
	if err != nil {
		return err
	}
	if approvedSpec == nil {
		return fmt.Errorf("no job proposal spec found for job id %s", idToAccept)
	}
	return nil
}
