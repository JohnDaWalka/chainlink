package node

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	//TODO replace with CTF?

	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/devenv"

	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

var (
	NodeTypeKey            = cre.NodeTypeKey
	NodeIDKey              = cre.NodeIDKey
	NodeOCR2KeyBundleIDKey = cre.NodeOCR2KeyBundleIDKey
	NodeOCRFamiliesKey     = cre.NodeOCRFamiliesKey
	DONIDKey               = cre.DONIDKey
	EnvironmentKey         = cre.EnvironmentKey
	ProductKey             = cre.ProductKey
	DONNameKey             = cre.DONNameKey
)

// ocr2 keys depend on report's target chain family
func CreateNodeOCR2KeyBundleIDKey(chainFamily string) string {
	return NodeOCR2KeyBundleIDKey + "_" + chainFamily
}

func CreateNodeOCRFamiliesListValue(families []string) string {
	return strings.Join(families, ",")
}

func ExtractBundleKeysPerFamily(n *cre.NodeMetadata) (map[string]string, error) {
	keyBundlesFamilies, fErr := FindLabelValue(n, cre.NodeOCRFamiliesKey)
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

// type Node struct {
// 	Name                  string
// 	IDs                   NodeIDs
// 	Keys                  *secrets.NodeKeys
// 	Addresses             Addresses
// 	JobDistributorDetails *JobDistributorDetails
// 	Roles                 []string

// 	Clients NodeClients
// 	DON     don.DON // to easily get parent info
// }

// func (n *Node) HasRole(role string) bool {
// 	for _, r := range n.Roles {
// 		if strings.EqualFold(r, role) {
// 			return true
// 		}
// 	}

// 	return false
// }

// func NewNode(ctx context.Context, name string, nodeMetadata *cre.NodeMetadata, clNode *clnode.Output) (*Node, error) {
// 	gqlClient, gqErr := client.NewWithContext(ctx, clNode.Node.ExternalURL, client.Credentials{
// 		Email:    clNode.Node.APIAuthUser,
// 		Password: clNode.Node.APIAuthPassword,
// 	})
// 	if gqErr != nil {
// 		return nil, fmt.Errorf("failed to create node graphql client: %w", gqErr)
// 	}

// 	chainlinkClient, cErr := clclient.NewChainlinkClient(&clclient.ChainlinkConfig{
// 		URL:         clNode.Node.ExternalURL,
// 		Email:       clNode.Node.APIAuthUser,
// 		Password:    clNode.Node.APIAuthPassword,
// 		InternalIP:  clNode.Node.InternalIP,
// 		HTTPTimeout: ptr.Ptr(10 * time.Second),
// 	}, framework.L)
// 	if cErr != nil {
// 		return nil, fmt.Errorf("failed to create node rest client: %w", cErr)
// 	}

// 	node := &Node{
// 		Clients: NodeClients{
// 			GQLClient:  gqlClient,
// 			RestClient: chainlinkClient,
// 		},
// 		Name:  name,
// 		Keys:  nodeMetadata.Keys,
// 		Roles: nodeMetadata.Roles,
// 	}

// 	for _, role := range node.Roles {
// 		switch role {
// 		case cre.WorkerNode:
// 			// multi address is not applicable for non-bootstrap nodes
// 			// explicitly set it to empty string to denote that
// 			node.Addresses.MultiAddress = ""

// 			// set admin address for non-bootstrap nodes (capability registry requires non-null admin address; use arbitrary default value if node is not configured)
// 			node.Addresses.AdminAddress = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
// 		case cre.BootstrapNode:
// 			// create multi address for OCR2, applicable only for bootstrap nodes
// 			p2pURL, err := url.Parse(clNode.Node.InternalP2PUrl)
// 			if err != nil {
// 				return nil, fmt.Errorf("failed to parse p2p url: %w", err)
// 			}
// 			node.Addresses.MultiAddress = fmt.Sprintf("%s:%s", clNode.Node.InternalIP, p2pURL.Port())

// 			// no need to set admin address for bootstrap nodes, as there will be no payment
// 			node.Addresses.AdminAddress = ""
// 		case cre.GatewayNode:
// 			// no specific data to set for gateway nodes yet
// 		default:
// 			return nil, fmt.Errorf("unknown node role: %s", role)
// 		}
// 	}

// 	return node, nil
// }

// type JobDistributorDetails struct {
// 	NodeID string // node id returned by JD after node is registered with it
// 	JDID   string // JD id returned by node after Job distributor is created in the node
// }

// // Do we need to store public address per chain or is it enough to store keys, so that we can derive public address when needed?
// type Addresses struct {
// 	AdminAddress string
// 	MultiAddress string
// }

// // type Keys struct {
// // 	CSA           string
// // 	EVM           map[uint64][]string
// // 	OCR2BundleIDs map[uint64][]string
// // 	Solana        []string
// // 	P2P           []string
// // }

// type NodeIDs struct {
// 	PeerID string
// }

// type NodeClients struct {
// 	GQLClient  client.Client             // graphql client to interact with the node
// 	RestClient *clclient.ChainlinkClient // rest client to interact with the node
// }

// type JDChainConfigInput struct {
// 	ChainID   string
// 	ChainType string
// }

// func (n *Node) CreateJDChainConfigs(ctx context.Context, chains []JDChainConfigInput, jd *devenv.JobDistributor) error {
// 	for _, chain := range chains {
// 		var account string

// 		switch strings.ToLower(chain.ChainType) {
// 		case chainselectors.FamilyEVM, chainselectors.FamilyTron:
// 			chainIDUint64, parseErr := strconv.ParseUint(chain.ChainID, 10, 64)
// 			if parseErr != nil {
// 				return fmt.Errorf("failed to parse chain id %s: %w", chain.ChainID, parseErr)
// 			}

// 			if chainIDUint64 == 0 {
// 				return fmt.Errorf("invalid chain id: %s", chain.ChainID)
// 			}

// 			evmKey, ok := n.Keys.EVM[chainIDUint64]
// 			if !ok {
// 				var fetchErr error
// 				accountAddr, fetchErr := n.Clients.GQLClient.FetchAccountAddress(ctx, chain.ChainID)
// 				if fetchErr != nil {
// 					return fmt.Errorf("failed to fetch account address for node %s: %w", n.Name, fetchErr)
// 				}
// 				if accountAddr == nil {
// 					return fmt.Errorf("no account address found for node %s", n.Name)
// 				}
// 				account = *accountAddr
// 			} else {
// 				account = evmKey.PublicAddress.Hex()
// 			}
// 			// accountAddr, err := n.gqlClient.FetchAccountAddress(ctx, chain.ChainID)
// 			// if err != nil {
// 			// 	return fmt.Errorf("failed to fetch account address for node %s: %w", n.Name, err)
// 			// }
// 			// if accountAddr == nil {
// 			// 	return fmt.Errorf("no account address found for node %s", n.Name)
// 			// }
// 			// n.AccountAddr[chain.ChainID] = *accountAddr
// 			// account = *accountAddr
// 		case chainselectors.FamilySolana:
// 			solKey, ok := n.Keys.Solana[chain.ChainID]
// 			if !ok {
// 				accounts, fetchErr := n.Clients.GQLClient.FetchKeys(ctx, chain.ChainType)
// 				if fetchErr != nil {
// 					return fmt.Errorf("failed to fetch account address for node %s and chain %s: %w", n.Name, chain.ChainType, fetchErr)
// 				}
// 				if len(accounts) == 0 {
// 					return fmt.Errorf("failed to fetch account address for node %s and chain %s", n.Name, chain.ChainType)
// 				}
// 				account = accounts[0]
// 			} else {
// 				account = solKey.PublicAddress.String()
// 			}
// 		case chainselectors.FamilyAptos:
// 			accounts, err := n.Clients.GQLClient.FetchKeys(ctx, chain.ChainType)
// 			if err != nil {
// 				return fmt.Errorf("failed to fetch account address for node %s and chain %s: %w", n.Name, chain.ChainType, err)
// 			}
// 			if len(accounts) == 0 {
// 				return fmt.Errorf("failed to fetch account address for node %s and chain %s", n.Name, chain.ChainType)
// 			}
// 			// TODO store that or not?
// 			// n.AccountAddr[chain.ChainID] = accounts[0]
// 			account = accounts[0]
// 		default:
// 			return fmt.Errorf("unsupported chainType %v", chain.ChainType)
// 		}

// 		// peerID, err := n.gqlClient.FetchP2PPeerID(ctx)
// 		// if err != nil {
// 		// 	return fmt.Errorf("failed to fetch peer id for node %s: %w", n.Name, err)
// 		// }
// 		// if peerID == nil {
// 		// 	return fmt.Errorf("no peer id found for node %s", n.Name)
// 		// }

// 		chainType := chain.ChainType
// 		if strings.EqualFold(chain.ChainType, blockchain.FamilyTron) {
// 			chainType = strings.ToUpper(blockchain.FamilyEVM)
// 		}
// 		ocr2BundleID, createErr := n.Clients.GQLClient.FetchOCR2KeyBundleID(ctx, chainType)
// 		if createErr != nil {
// 			return fmt.Errorf("failed to fetch OCR2 key bundle id for node %s: %w", n.Name, createErr)
// 		}
// 		if ocr2BundleID == "" {
// 			return fmt.Errorf("no OCR2 key bundle id found for node %s", n.Name)
// 		}

// 		n.Keys.OCR2BundleIDs[strings.ToLower(chainType)] = ocr2BundleID

// 		// fetch node labels to know if the node is bootstrap or plugin
// 		// if multi address is set, then it's a bootstrap node
// 		// isBootstrap := n.HasRole(cre.BootstrapNode)
// 		// for _, label := range n.labels {
// 		// 	if label.Key == LabelNodeTypeKey && value(label.Value) == LabelNodeTypeValueBootstrap {
// 		// 		isBootstrap = true
// 		// 		break
// 		// 	}
// 		// }

// 		// retry twice with 5 seconds interval to create JobDistributorChainConfig
// 		retryErr := retry.Do(ctx, retry.WithMaxDuration(10*time.Second, retry.NewConstant(3*time.Second)), func(ctx context.Context) error {
// 			// check the node chain config to see if this chain already exists
// 			nodeChainConfigs, listErr := jd.ListNodeChainConfigs(context.Background(), &nodev1.ListNodeChainConfigsRequest{
// 				Filter: &nodev1.ListNodeChainConfigsRequest_Filter{
// 					NodeIds: []string{n.JobDistributorDetails.NodeID},
// 				}})
// 			if listErr != nil {
// 				return retry.RetryableError(fmt.Errorf("failed to list node chain configs for node %s, retrying..: %w", n.Name, listErr))
// 			}
// 			if nodeChainConfigs != nil {
// 				for _, chainConfig := range nodeChainConfigs.ChainConfigs {
// 					if chainConfig.Chain.Id == chain.ChainID {
// 						return nil
// 					}
// 				}
// 			}

// 			// we need to create JD chain config for each chain, because later on changestes ask the node for that chain data
// 			// each node needs to have OCR2 enabled, because p2pIDs are used by some contracts to identify nodes (e.g. capability registry)
// 			_, createErr = n.Clients.GQLClient.CreateJobDistributorChainConfig(ctx, client.JobDistributorChainConfigInput{
// 				JobDistributorID: n.JobDistributorDetails.JDID,
// 				ChainID:          chain.ChainID,
// 				ChainType:        chainType,
// 				AccountAddr:      account,
// 				AdminAddr:        n.Addresses.AdminAddress,
// 				Ocr2Enabled:      true,
// 				Ocr2IsBootstrap:  n.HasRole(cre.BootstrapNode),
// 				Ocr2Multiaddr:    n.Addresses.MultiAddress,
// 				Ocr2P2PPeerID:    n.Keys.P2PKey.PeerID.String(),
// 				Ocr2KeyBundleID:  ocr2BundleID,
// 				Ocr2Plugins:      `{}`,
// 			})
// 			// TODO: add a check if the chain config failed because of a duplicate in that case, should we update or return success?
// 			if createErr != nil {
// 				return fmt.Errorf("failed to create JD chain configuration for node %s: %w", n.Name, createErr)
// 			}

// 			// JD silently fails to update nodeChainConfig. Therefore, we fetch the node config and
// 			// if it's not updated , throw an error
// 			return retry.RetryableError(errors.New("retrying CreateChainConfig in JD"))
// 		})

// 		if retryErr != nil {
// 			return fmt.Errorf("failed to create  create JD chain configuration for node %s: %w", n.Name, retryErr)
// 		}
// 	}
// 	return nil
// }

// // AcceptJob accepts the job proposal for the given job proposal spec
// func (n *Node) AcceptJob(ctx context.Context, spec string) error {
// 	// fetch JD to get the job proposals
// 	jd, err := n.Clients.GQLClient.GetJobDistributor(ctx, n.JobDistributorDetails.JDID)
// 	if err != nil {
// 		return err
// 	}
// 	if jd.GetJobProposals() == nil {
// 		return fmt.Errorf("no job proposals found for node %s", n.Name)
// 	}
// 	// locate the job proposal id for the given job spec
// 	var idToAccept string
// 	for _, jp := range jd.JobProposals {
// 		if jp.LatestSpec.Definition == spec {
// 			idToAccept = jp.Id
// 			break
// 		}
// 	}
// 	if idToAccept == "" {
// 		return fmt.Errorf("no job proposal found for job spec %s", spec)
// 	}
// 	approvedSpec, err := n.Clients.GQLClient.ApproveJobProposalSpec(ctx, idToAccept, false)
// 	if err != nil {
// 		return err
// 	}
// 	if approvedSpec == nil {
// 		return fmt.Errorf("no job proposal spec found for job id %s", idToAccept)
// 	}
// 	return nil
// }

// // RegisterNodeToJobDistributor fetches the CSA public key of the node and registers the node with the job distributor
// // it sets the node id returned by JobDistributor as a result of registration in the node struct
// func (n *Node) RegisterNodeToJobDistributor(ctx context.Context, jd *devenv.JobDistributor, labels []*ptypes.Label) error {
// 	// Get the public key of the node
// 	csaKeyRes, err := n.Clients.GQLClient.FetchCSAPublicKey(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	if csaKeyRes == nil {
// 		return fmt.Errorf("no csa key found for node %s", n.Name)
// 	}

// 	n.Keys.CSAKey = crypto.CSAKey{
// 		Key: *csaKeyRes,
// 	}

// 	// // tag nodes with p2p_id for easy lookup
// 	// peerID, err := n.gqlClient.FetchP2PPeerID(ctx)
// 	// if err != nil {
// 	// 	return fmt.Errorf("failed to fetch peer id for node %s: %w", n.Name, err)
// 	// }
// 	// if peerID == nil {
// 	// 	return fmt.Errorf("no peer id found for node %s", n.Name)
// 	// }
// 	labels = append(labels, &ptypes.Label{
// 		Key:   devenv.LabelNodeP2PIDKey,
// 		Value: ptr.Ptr(n.Keys.P2PKey.PeerID.String()),
// 	})

// 	// register the node in the job distributor
// 	registerResponse, err := jd.RegisterNode(ctx, &nodev1.RegisterNodeRequest{
// 		PublicKey: n.Keys.CSAKey.CleansedKey(),
// 		Labels:    labels,
// 		Name:      n.Name,
// 	})

// 	if n.JobDistributorDetails == nil {
// 		n.JobDistributorDetails = &JobDistributorDetails{}
// 	}

// 	// node already registered, fetch it's id
// 	if err != nil && strings.Contains(err.Error(), "AlreadyExists") {
// 		nodesResponse, listErr := jd.ListNodes(ctx, &nodev1.ListNodesRequest{
// 			Filter: &nodev1.ListNodesRequest_Filter{
// 				Selectors: []*ptypes.Selector{
// 					{
// 						Key:   devenv.LabelNodeP2PIDKey,
// 						Op:    ptypes.SelectorOp_EQ,
// 						Value: ptr.Ptr(n.Keys.P2PKey.PeerID.String()),
// 					},
// 				},
// 			},
// 		})
// 		if listErr != nil {
// 			return listErr
// 		}
// 		nodes := nodesResponse.GetNodes()
// 		if len(nodes) == 0 {
// 			return fmt.Errorf("failed to find node: %v", n.Name)
// 		}
// 		n.JobDistributorDetails.NodeID = nodes[0].Id
// 		return nil
// 	} else if err != nil {
// 		return fmt.Errorf("failed to register node %s: %w", n.Name, err)
// 	}
// 	if registerResponse.GetNode().GetId() == "" {
// 		return fmt.Errorf("no node id returned from job distributor for node %s", n.Name)
// 	}
// 	n.JobDistributorDetails.NodeID = registerResponse.GetNode().GetId()

// 	return nil
// }

// // CreateJobDistributor fetches the keypairs from the job distributor and creates the job distributor in the node
// // and returns the job distributor id
// func (n *Node) CreateJobDistributor(ctx context.Context, jd *devenv.JobDistributor) (string, error) {
// 	// Get the keypairs from the job distributor
// 	csaKey, err := jd.GetCSAPublicKey(ctx)
// 	if err != nil {
// 		return "", err
// 	}
// 	// create the job distributor in the node with the csa key
// 	resp, err := n.Clients.GQLClient.ListJobDistributors(ctx)
// 	if err != nil {
// 		return "", fmt.Errorf("could not list job distributors: %w", err)
// 	}
// 	if len(resp.FeedsManagers.Results) > 0 {
// 		for _, fm := range resp.FeedsManagers.Results {
// 			if fm.GetPublicKey() == csaKey {
// 				return fm.GetId(), nil
// 			}
// 		}
// 	}
// 	return n.Clients.GQLClient.CreateJobDistributor(ctx, client.JobDistributorInput{
// 		Name:      "Job Distributor",
// 		Uri:       jd.WSRPC,
// 		PublicKey: csaKey,
// 	})
// }

// // SetUpAndLinkJobDistributor sets up the job distributor in the node and registers the node with the job distributor
// // it sets the job distributor id for node
// func (n *Node) SetUpAndLinkJobDistributor(ctx context.Context, jd *devenv.JobDistributor, labels []*ptypes.Label) error {
// 	// register the node in the job distributor
// 	err := n.RegisterNodeToJobDistributor(ctx, jd, labels)
// 	if err != nil {
// 		return err
// 	}
// 	// now create the job distributor in the node
// 	id, err := n.CreateJobDistributor(ctx, jd)
// 	if err != nil &&
// 		!strings.Contains(err.Error(), "DuplicateFeedsManagerError") {
// 		return fmt.Errorf("failed to create job distributor in node %s: %w", n.Name, err)
// 	}
// 	// wait for the node to connect to the job distributor
// 	err = retry.Do(ctx, retry.WithMaxDuration(1*time.Minute, retry.NewFibonacci(1*time.Second)), func(ctx context.Context) error {
// 		getRes, getErr := jd.GetNode(ctx, &nodev1.GetNodeRequest{
// 			Id: n.JobDistributorDetails.NodeID,
// 		})
// 		if getErr != nil {
// 			return retry.RetryableError(fmt.Errorf("failed to get node %s: %w", n.Name, getErr))
// 		}
// 		if getRes.GetNode() == nil {
// 			return fmt.Errorf("no node found for node id %s", n.JobDistributorDetails.NodeID)
// 		}
// 		if !getRes.GetNode().IsConnected {
// 			return retry.RetryableError(fmt.Errorf("node %s not connected to job distributor", n.Name))
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to connect node %s to job distributor: %w", n.Name, err)
// 	}
// 	n.JobDistributorDetails.JDID = id
// 	return nil
// }

// // func ptr[T any](v T) *T {
// // 	return &v
// // }

// // func value[T any](v *T) T {
// // 	zero := new(T)
// // 	if v == nil {
// // 		return *zero
// // 	}
// // 	return *v
// // }

// func (n *Node) CancelProposalsByExternalJobID(ctx context.Context, externalJobIDs []string) ([]string, error) {
// 	jd, err := n.Clients.GQLClient.GetJobDistributor(ctx, n.JobDistributorDetails.JDID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if jd.GetJobProposals() == nil {
// 		return nil, fmt.Errorf("no job proposals found for node %s", n.Name)
// 	}

// 	proposalIDs := []string{}
// 	for _, jp := range jd.JobProposals {
// 		if !slices.Contains(externalJobIDs, jp.ExternalJobID) {
// 			continue
// 		}

// 		proposalIDs = append(proposalIDs, jp.Id)
// 		spec, err := n.Clients.GQLClient.CancelJobProposalSpec(ctx, jp.Id)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if spec == nil {
// 			return nil, fmt.Errorf("no job proposal spec found for id %s", jp.Id)
// 		}
// 	}

// 	return proposalIDs, nil
// }

// func (n *Node) ApproveProposals(ctx context.Context, proposalIDs []string) error {
// 	for _, proposalID := range proposalIDs {
// 		spec, err := n.Clients.GQLClient.ApproveJobProposalSpec(ctx, proposalID, false)
// 		if err != nil {
// 			return err
// 		}
// 		if spec == nil {
// 			return fmt.Errorf("no job proposal spec found for id %s", proposalID)
// 		}
// 	}
// 	return nil
// }
