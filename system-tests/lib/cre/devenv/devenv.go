package devenv

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"

	clclient "github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
)

// All label keys:
// * must be non-empty,
// * must be 63 characters or less,
// * must begin and end with an alphanumeric character ([a-z0-9A-Z]),
// * could contain dashes (-), underscores (_), dots (.), and alphanumerics between.
//
// All label values:
// * must be 63 characters or less (can be empty),
// * unless empty, must begin and end with an alphanumeric character ([a-z0-9A-Z]),
// * could contain dashes (-), underscores (_), dots (.), and alphanumerics between.
//
// Source: https://github.com/smartcontractkit/job-distributor/blob/main/pkg/entities/labels.go
const (
	LabelNodeTypeKey            = "type"
	LabelNodeTypeValueBootstrap = "bootstrap"
	LabelNodeTypeValuePlugin    = "plugin"

	LabelNodeP2PIDKey = "p2p_id"

	LabelJobTypeKey         = "jobType"
	LabelJobTypeValueLLO    = "llo"
	LabelJobTypeValueStream = "stream"

	LabelEnvironmentKey = "environment"
	LabelStreamIDKey    = "streamID"

	LabelProductKey = "product"
)

// NodeInfo holds the information required to create a node
type NodeInfo struct {
	CLConfig      clclient.ChainlinkConfig // config to connect to chainlink node via API
	P2PPort       string                   // port for P2P communication
	IsBootstrap   bool                     // denotes if the node is a bootstrap node
	Name          string                   // name of the node, used to identify the node, helpful in logs
	AdminAddr     string                   // admin address to send payments to, applicable only for non-bootstrap nodes
	MultiAddr     string                   // multi address denoting node's FQN (needed for deriving P2PBootstrappers in OCR), applicable only for bootstrap nodes
	Labels        map[string]string        // labels to use when registering the node with job distributor
	ContainerName string                   // name of Docker container
}

// type DON struct {
// 	Nodes []Node
// }

// func (don *DON) PluginNodes() []Node {
// 	var pluginNodes []Node
// 	for _, node := range don.Nodes {
// 		for _, label := range node.labels {
// 			if label.Key == LabelNodeTypeKey && value(label.Value) == LabelNodeTypeValuePlugin {
// 				pluginNodes = append(pluginNodes, node)
// 			}
// 		}
// 	}
// 	return pluginNodes
// }

// func (don *DON) JDNodeIDs() []string {
// 	nodeIDs := []string{}
// 	for _, node := range don.Nodes {
// 		nodeIDs = append(nodeIDs, node.JobDistributorDetails.NodeID)
// 	}
// 	return nodeIDs
// }

// func (don *DON) CreateSupportedChains(ctx context.Context, chains []ChainConfig, jd *JobDistributor) error {
// 	g := new(errgroup.Group)
// 	for i := range don.Nodes {
// 		g.Go(func() error {
// 			node := &don.Nodes[i]
// 			var jdChains []JDChainConfigInput
// 			for _, chain := range chains {
// 				jdChains = append(jdChains, JDChainConfigInput{
// 					ChainID:   chain.ChainID,
// 					ChainType: chain.ChainType,
// 				})
// 			}
// 			if err1 := node.CreateJDChainConfigs(ctx, jdChains, jd); err1 != nil {
// 				return err1
// 			}
// 			don.Nodes[i] = *node
// 			return nil
// 		})
// 	}
// 	return g.Wait()
// }

// NewRegisteredDON creates a DON with the given node info, registers the nodes with the job distributor
// and sets up the job distributor in the nodes
// func NewRegisteredDON(ctx context.Context, nodeInfo []NodeInfo, jd JobDistributor) (*DON, error) {
// 	don := &DON{
// 		Nodes: make([]Node, 0),
// 	}
// 	for i, info := range nodeInfo {
// 		if info.Name == "" {
// 			info.Name = fmt.Sprintf("node-%d", i)
// 		}
// 		node, err := NewNodeWithContext(ctx, info)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to create node %d: %w", i, err)
// 		}

// 		if info.IsBootstrap {
// 			// create multi address for OCR2, applicable only for bootstrap nodes
// 			if info.MultiAddr == "" {
// 				node.multiAddr = fmt.Sprintf("%s:%s", info.CLConfig.InternalIP, info.P2PPort)
// 			} else {
// 				node.multiAddr = info.MultiAddr
// 			}
// 			// no need to set admin address for bootstrap nodes, as there will be no payment
// 			node.adminAddr = ""
// 			node.labels = append(node.labels, &ptypes.Label{
// 				Key:   LabelNodeTypeKey,
// 				Value: ptr(LabelNodeTypeValueBootstrap),
// 			})
// 		} else {
// 			// multi address is not applicable for non-bootstrap nodes
// 			// explicitly set it to empty string to denote that
// 			node.multiAddr = ""

// 			// set admin address for non-bootstrap nodes
// 			node.adminAddr = info.AdminAddr

// 			// capability registry requires non-null admin address; use arbitrary default value if node is not configured
// 			if info.AdminAddr == "" {
// 				node.adminAddr = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
// 			}

// 			node.labels = append(node.labels, &ptypes.Label{
// 				Key:   LabelNodeTypeKey,
// 				Value: ptr(LabelNodeTypeValuePlugin),
// 			})

// 			for key, val := range info.Labels {
// 				node.labels = append(node.labels, &ptypes.Label{
// 					Key:   key,
// 					Value: ptr(val),
// 				})
// 			}
// 		}
// 		// Set up Job distributor in node and register node with the job distributor
// 		err = node.SetUpAndLinkJobDistributor(ctx, jd)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to set up job distributor in node %s: %w", info.Name, err)
// 		}

// 		don.Nodes = append(don.Nodes, *node)
// 	}
// 	return don, nil
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

// func NewDON(ctx context.Context, donMetadata *cre.DonMetadata, nodeSetOut *cre.WrappedNodeOutput, supportedChains []ChainConfig, jd *JobDistributor) (*DON, error) {
// 	don := &DON{
// 		Nodes: make([]Node, 0),
// 	}
// 	for idx, nodeMetadata := range donMetadata.NodesMetadata {
// 		node, err := NewNode(ctx, fmt.Sprintf("%s-node%d", donMetadata.Name, idx), nodeMetadata, nodeSetOut.CLNodes[idx])
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to create node %d: %w", idx, err)
// 		}

// 		labels := make([]*ptypes.Label, 0)

// 		for _, role := range node.Roles {
// 			switch role {
// 			case cre.WorkerNode:
// 				labels = append(labels, &ptypes.Label{
// 					Key:   LabelNodeTypeKey,
// 					Value: ptr.Ptr(LabelNodeTypeValuePlugin),
// 				})
// 			case cre.BootstrapNode:
// 				labels = append(labels, &ptypes.Label{
// 					Key:   LabelNodeTypeKey,
// 					Value: ptr.Ptr(LabelNodeTypeValueBootstrap),
// 				})
// 			case cre.GatewayNode:
// 				// no specific data to set for gateway nodes yet
// 			default:
// 				return nil, fmt.Errorf("unknown node role: %s", role)
// 			}
// 		}

// 		// Set up Job distributor in node and register node with the job distributor
// 		setupErr := node.SetUpAndLinkJobDistributor(ctx, jd, labels)
// 		if setupErr != nil {
// 			return nil, fmt.Errorf("failed to set up job distributor in node %s: %w", node.Name, setupErr)
// 		}

// 		for _, role := range node.Roles {
// 			switch role {
// 			case cre.WorkerNode, cre.BootstrapNode:
// 				if err := don.CreateSupportedChains(ctx, supportedChains, jd); err != nil {
// 					return nil, fmt.Errorf("failed to create supported chains: %w", err)
// 				}
// 			case cre.GatewayNode:
// 				// no chains configuration needed for gateway nodes
// 			default:
// 				return nil, fmt.Errorf("unknown node role: %s", role)
// 			}
// 		}

// 		don.Nodes = append(don.Nodes, *node)
// 	}
// 	return don, nil
// }

// func NewNodeWithContext(ctx context.Context, nodeInfo NodeInfo) (*Node, error) {
// 	gqlClient, err := client.NewWithContext(ctx, nodeInfo.CLConfig.URL, client.Credentials{
// 		Email:    nodeInfo.CLConfig.Email,
// 		Password: nodeInfo.CLConfig.Password,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create node graphql client: %w", err)
// 	}
// 	chainlinkClient, err := clclient.NewChainlinkClient(&nodeInfo.CLConfig, zerolog.Logger{})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create node rest client: %w", err)
// 	}
// 	// node Labels so that it's easier to query them
// 	labels := make([]*ptypes.Label, 0)
// 	for key, value := range nodeInfo.Labels {
// 		labels = append(labels, &ptypes.Label{
// 			Key:   key,
// 			Value: &value,
// 		})
// 	}
// 	return &Node{
// 		gqlClient:              gqlClient,
// 		restClient:             chainlinkClient,
// 		Name:                   nodeInfo.Name,
// 		adminAddr:              nodeInfo.AdminAddr,
// 		multiAddr:              nodeInfo.MultiAddr,
// 		labels:                 labels,
// 		ChainsOcr2KeyBundlesID: make(map[string]string),
// 	}, nil
// }

// type Node struct {
// 	NodeID                 string            // node id returned by job distributor after node is registered with it
// 	JDId                   string            // job distributor id returned by node after Job distributor is created in node
// 	Name                   string            // name of the node
// 	AccountAddr            map[string]string // chain id to node's account address mapping for supported chains
// 	ChainsOcr2KeyBundlesID map[string]string
// 	gqlClient              client.Client             // graphql client to interact with the node
// 	restClient             *clclient.ChainlinkClient // rest client to interact with the node
// 	labels                 []*ptypes.Label           // labels with which the node is registered with the job distributor
// 	adminAddr              string                    // admin address to send payments to, applicable only for non-bootstrap nodes
// 	multiAddr              string                    // multi address denoting node's FQN (needed for deriving P2PBootstrappers in OCR), applicable only for bootstrap nodes
// }

// type DON struct {
// 	Name string
// 	ID   uint32

// 	Nodes []Node

// 	Capabilities    []cre.CapabilityFlag
// 	Roles           []string // workflow, capability, gateway
// 	SupportedChains []uint64 // chain selector... optionally? to indicate, whether each node should connect to every chain in the environment or only some
// }

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

// func (n *Node) CreateJDChainConfigs(ctx context.Context, chains []JDChainConfigInput, jd *JobDistributor) error {
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
// func (n *Node) RegisterNodeToJobDistributor(ctx context.Context, jd *JobDistributor, labels []*ptypes.Label) error {
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
// 		Key:   LabelNodeP2PIDKey,
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
// 						Key:   LabelNodeP2PIDKey,
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
// func (n *Node) CreateJobDistributor(ctx context.Context, jd *JobDistributor) (string, error) {
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
// func (n *Node) SetUpAndLinkJobDistributor(ctx context.Context, jd *JobDistributor, labels []*ptypes.Label) error {
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

// func value[T any](v *T) T {
// 	zero := new(T)
// 	if v == nil {
// 		return *zero
// 	}
// 	return *v
// }

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

type JDConfig struct {
	GRPC     string
	WSRPC    string
	Creds    credentials.TransportCredentials
	Auth     oauth2.TokenSource
	NodeInfo []NodeInfo
}

// JobDistributor implements the OffchainClient interface in CLDF and wraps the CLDF JD client and add DON functionality.
// The CLDF JD client does not have the DON functionality, so we wrap it here.
type JobDistributor struct {
	*jd.JobDistributor
	don *don.DON
}

// NewJDClient creates a new Job Distributor client with the provided configuration.
func NewJDClient(ctx context.Context, cfg JDConfig) (cldf_offchain.Client, error) {
	jdConfig := jd.JDConfig{
		GRPC:  cfg.GRPC,
		WSRPC: cfg.WSRPC,
		Creds: cfg.Creds,
		Auth:  cfg.Auth,
	}
	jdClient, err := jd.NewJDClient(jdConfig)
	if err != nil {
		return nil, err
	}
	donJDClient := &JobDistributor{
		JobDistributor: jdClient,
	}
	// if len(cfg.NodeInfo) > 0 {
	// 	donJDClient.don, err = NewRegisteredDON(ctx, cfg.NodeInfo, *donJDClient)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to create registered DON: %w", err)
	// 	}
	// }
	return donJDClient, err
}

// ProposeJob proposes jobs through the jobService and accepts the proposed job on selected node based on ProposeJobRequest.NodeId
func (jd JobDistributor) ProposeJob(ctx context.Context, in *jobv1.ProposeJobRequest, opts ...grpc.CallOption) (*jobv1.ProposeJobResponse, error) {
	res, err := jd.JobDistributor.ProposeJob(ctx, in, opts...)
	if err != nil {
		return nil, err
	}

	if jd.don == nil || len(jd.don.Nodes) == 0 {
		return res, nil
	}
	for _, node := range jd.don.Nodes {
		if node.JobDistributorDetails.NodeID != in.NodeId {
			continue
		}
		// TODO : is there a way to accept the job with proposal id?
		if err := node.AcceptJob(ctx, res.Proposal.Spec); err != nil {
			return nil, fmt.Errorf("failed to accept job. err: %w", err)
		}
	}
	return res, nil
}

// ChainConfigFromWrapped converts a single wrapped chain into a devenv.ChainConfig.
func ChainConfigFromWrapped(w *cre.WrappedBlockchainOutput) (ChainConfig, error) {
	if w == nil || w.BlockchainOutput == nil || len(w.BlockchainOutput.Nodes) == 0 {
		return ChainConfig{}, errors.New("invalid wrapped blockchain output")
	}
	n := w.BlockchainOutput.Nodes[0]

	cfg := ChainConfig{
		WSRPCs: []CribRPCs{{
			External: n.ExternalWSUrl, Internal: n.InternalWSUrl,
		}},
		HTTPRPCs: []CribRPCs{{
			External: n.ExternalHTTPUrl, Internal: n.InternalHTTPUrl,
		}},
	}

	cfg.ChainType = strings.ToUpper(w.BlockchainOutput.Family)

	// Solana
	if w.SolChain != nil {
		cfg.ChainID = w.SolChain.ChainID
		cfg.SolDeployerKey = w.SolChain.PrivateKey
		cfg.SolArtifactDir = w.SolChain.ArtifactsDir
		return cfg, nil
	}

	if strings.EqualFold(cfg.ChainType, blockchain.FamilyTron) {
		cfg.ChainID = strconv.FormatUint(w.ChainID, 10)
		privateKey, err := ethcrypto.HexToECDSA(w.DeployerPrivateKey)
		if err != nil {
			return ChainConfig{}, errors.Wrap(err, "failed to parse private key for Tron")
		}

		deployerKey, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(conversions.MustSafeInt64(w.ChainID)))
		if err != nil {
			return ChainConfig{}, errors.Wrap(err, "failed to create transactor for Tron")
		}
		cfg.DeployerKey = deployerKey
		return cfg, nil
	}

	// EVM
	if w.SethClient == nil {
		return ChainConfig{}, fmt.Errorf("blockchain output evm family without SethClient for chainID %d", w.ChainID)
	}

	cfg.ChainID = strconv.FormatUint(w.ChainID, 10)
	cfg.ChainName = w.SethClient.Cfg.Network.Name
	// ensure nonce fetched from chain at use time
	cfg.DeployerKey = w.SethClient.NewTXOpts(seth.WithNonce(nil))

	return cfg, nil
}
