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
	"github.com/smartcontractkit/chainlink/deployment/environment/web/sdk/client"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
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

	if slices.Contains(roles, cre.WorkerNode) && slices.Contains(roles, cre.BootstrapNode) {
		return nil, fmt.Errorf("node cannot be both worker and bootstrap")
	}

	node, err = fetchBasicNodeData(ctx, node)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		var fetchErr error
		switch role {
		case cre.WorkerNode:
			// multi address is not applicable for non-bootstrap nodes
			// explicitly set it to empty string to denote that
			node.MultiAddr = ""

			// set admin address for non-bootstrap nodes (capability registry requires non-null admin address; use arbitrary default value if node is not configured)
			node.AdminAddr = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"

			node.JDLabels = append(node.JDLabels, &ptypes.Label{
				Key:   devenv.LabelNodeTypeKey,
				Value: ptr.Ptr(devenv.LabelNodeTypeValuePlugin),
			})

			node, fetchErr = fetchWorkerNodeData(ctx, node)
		case cre.BootstrapNode:
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

			node, fetchErr = fetchBootstrapNodeData(ctx, node)
		case cre.GatewayNode:
			// nothing to fetch for gateway node
		default:
			return nil, fmt.Errorf("unknown node role: %s", role)
		}

		if fetchErr != nil {
			return nil, fetchErr
		}
	}

	return node, nil
}

const LabelCSAKey = "csa_key"

func fetchBasicNodeData(ctx context.Context, n *cre.Node) (*cre.Node, error) {
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

func fetchWorkerNodeData(ctx context.Context, n *cre.Node) (*cre.Node, error) {
	n, err := fetchP2PID(ctx, n)
	if err != nil {
		return nil, err
	}

	n, err = fetchChainData(ctx, n)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func fetchBootstrapNodeData(ctx context.Context, n *cre.Node) (*cre.Node, error) {
	n, err := fetchP2PID(ctx, n)
	if err != nil {
		return nil, err
	}

	n, err = fetchChainData(ctx, n)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func fetchP2PID(ctx context.Context, n *cre.Node) (*cre.Node, error) {
	if n.PeerID != "" {
		return n, nil
	}

	peerID, err := n.GQLClient.FetchP2PPeerID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch peer id for node %s: %w", n.Name, err)
	}

	if peerID == nil {
		return nil, fmt.Errorf("no peer id found for node %s", n.Name)
	}

	n.PeerID = *peerID

	n.JDLabels = append(n.JDLabels, &ptypes.Label{
		Key:   devenv.LabelNodeP2PIDKey,
		Value: &n.PeerID,
	})

	return n, nil
}

func fetchChainData(ctx context.Context, n *cre.Node) (*cre.Node, error) {
	for _, selector := range n.SupportedChainSelectors {
		if n.AccountAddr == nil {
			n.AccountAddr = make(map[string]string)
		}

		chainID, chainFamily, chErr := selectorToIDAndFamily(selector)
		if chErr != nil {
			return nil, fmt.Errorf("failed to get chain id and chain family from selector %d: %w", selector, chErr)
		}

		if _, ok := n.AccountAddr[chainID]; ok {
			continue
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

		// We need to fetch this data for all nodes, even if they do not have OCR-based capabilities, because it is required in order to create the JobDistributorChainConfig in JD
		// which is later treated as data source for various changesets (e.g. they read P2PIDs from there)
		ocr2BundleId, err := n.GQLClient.FetchOCR2KeyBundleID(ctx, strings.ToUpper(chainFamily))
		if err != nil {
			return nil, fmt.Errorf("failed to fetch OCR2 key bundle id for node %s: %w", n.Name, err)
		}
		if ocr2BundleId == "" {
			return nil, fmt.Errorf("no OCR2 key bundle id found for node %s", n.Name)
		}

		n.ChainsOcr2KeyBundlesID[chainFamily] = ocr2BundleId
	}

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

			ocr2BundleId, ok := n.ChainsOcr2KeyBundlesID[chainFamily]
			if !ok {
				return fmt.Errorf("no ocr2 key bundle id found for node %s and chain family %s", n.Name, chainFamily)
			}

			// TODO probably capabilities should be able to modify this config according to their needs (like adding OCR2 keys)
			config := client.JobDistributorChainConfigInput{
				JobDistributorID: n.JDId,
				ChainID:          chainID,
				ChainType:        strings.ToUpper(chainFamily),
				AccountAddr:      address,
				AdminAddr:        n.AdminAddr,
				Ocr2Enabled:      true,
				Ocr2KeyBundleID:  ocr2BundleId,
				Ocr2IsBootstrap:  slices.Contains(n.Roles, cre.BootstrapNode),
				Ocr2Multiaddr:    n.MultiAddr,
				Ocr2P2PPeerID:    n.PeerID,
				// Ocr2Plugins:      `{"commit":true,"execute":true,"median":false,"mercury":false}` // TODO: should this be configurable? Do we even need this?
				Ocr2Plugins: `{}`, // TODO: should this be configurable? Do we even need this?
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
func RegisterNodeToJobDistributor(ctx context.Context, n cre.Node, jd *jd.JobDistributor) (*cre.Node, error) {
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
func CreateJobDistributor(ctx context.Context, n cre.Node, jd *jd.JobDistributor) (string, error) {
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
func SetUpAndLinkJobDistributor(ctx context.Context, n cre.Node, jd *jd.JobDistributor) (*cre.Node, error) {
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
