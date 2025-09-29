package cre

import (
	"context"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sethvargo/go-retry"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	clclient "github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	"github.com/smartcontractkit/chainlink/deployment/environment/web/sdk/client"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/secrets"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

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

type DON struct {
	Name string
	ID   uint64

	Nodes []*Node

	Capabilities    []CapabilityFlag
	Roles           []string // workflow, capability, gateway
	SupportedChains []uint64 // chain selector... optionally? to indicate, whether each node should connect to every chain in the environment or only some
}

func (m *DON) ContainsBootstrapNode() bool {
	for _, node := range m.Nodes {
		if slices.Contains(node.Roles, BootstrapNode) {
			return true
		}
	}

	return false
}

// Currently only one bootstrap node is supported.
func (m *DON) BootstrapNode() (*Node, error) {
	if !m.ContainsBootstrapNode() {
		return nil, errors.New("don does not contain a bootstrap node")
	}

	for _, node := range m.Nodes {
		if slices.Contains(node.Roles, BootstrapNode) {
			return node, nil
		}
	}

	return nil, errors.New("no bootstrap node found in don")
}

func (m *DON) WorkerNodes() ([]*Node, error) {
	workers := make([]*Node, 0)
	for _, node := range m.Nodes {
		if slices.Contains(node.Roles, WorkerNode) {
			workers = append(workers, node)
		}
	}

	if len(workers) == 0 {
		return nil, errors.New("don does not contain any worker nodes")
	}

	return workers, nil
}

func (don *DON) JDNodeIDs() []string {
	nodeIDs := []string{}
	for _, n := range don.Nodes {
		nodeIDs = append(nodeIDs, n.JobDistributorDetails.NodeID)
	}
	return nodeIDs
}

func NewDON(ctx context.Context, donMetadata *DonMetadata, ctfNodes []*clnode.Output, supportedChains []ChainConfig, jd *JobDistributor) (*DON, error) {
	don := &DON{
		Nodes: make([]*Node, 0),
	}
	for idx, nodeMetadata := range donMetadata.NodesMetadata {
		node, err := NewNode(ctx, fmt.Sprintf("%s-node%d", donMetadata.Name, idx), nodeMetadata, ctfNodes[idx])
		if err != nil {
			return nil, fmt.Errorf("failed to create node %d: %w", idx, err)
		}

		labels := make([]*ptypes.Label, 0)

		for _, role := range node.Roles {
			switch role {
			case WorkerNode:
				labels = append(labels, &ptypes.Label{
					Key:   LabelNodeTypeKey,
					Value: ptr.Ptr(LabelNodeTypeValuePlugin),
				})
			case BootstrapNode:
				labels = append(labels, &ptypes.Label{
					Key:   LabelNodeTypeKey,
					Value: ptr.Ptr(LabelNodeTypeValueBootstrap),
				})
			case GatewayNode:
				// no specific data to set for gateway nodes yet
			default:
				return nil, fmt.Errorf("unknown node role: %s", role)
			}
		}

		// Set up Job distributor in node and register node with the job distributor
		setupErr := node.SetUpAndLinkJobDistributor(ctx, jd, labels)
		if setupErr != nil {
			return nil, fmt.Errorf("failed to set up job distributor in node %s: %w", node.Name, setupErr)
		}

		for _, role := range node.Roles {
			switch role {
			case WorkerNode, BootstrapNode:
				if err := don.CreateSupportedChains(ctx, supportedChains, jd); err != nil {
					return nil, fmt.Errorf("failed to create supported chains: %w", err)
				}
			case GatewayNode:
				// no chains configuration needed for gateway nodes
			default:
				return nil, fmt.Errorf("unknown node role: %s", role)
			}
		}

		don.Nodes = append(don.Nodes, node)
	}
	return don, nil
}

func (don *DON) CreateSupportedChains(ctx context.Context, chains []ChainConfig, jd *JobDistributor) error {
	g := new(errgroup.Group)
	for i := range don.Nodes {
		g.Go(func() error {
			n := don.Nodes[i]
			var jdChains []JDChainConfigInput
			for _, chain := range chains {
				jdChains = append(jdChains, JDChainConfigInput{
					ChainID:   chain.ChainID,
					ChainType: chain.ChainType,
				})
			}
			if err1 := n.CreateJDChainConfigs(ctx, jdChains, jd); err1 != nil {
				return err1
			}
			don.Nodes[i] = n
			return nil
		})
	}
	return g.Wait()
}

type Node struct {
	Name                  string
	Host                  string
	IDs                   NodeIDs
	Keys                  *secrets.NodeKeys
	Addresses             Addresses
	JobDistributorDetails *JobDistributorDetails
	Roles                 []string

	Clients NodeClients
	DON     DON // to easily get parent info
}

func (n *Node) HasRole(role string) bool {
	for _, r := range n.Roles {
		if strings.EqualFold(r, role) {
			return true
		}
	}

	return false
}

func NewNode(ctx context.Context, name string, nodeMetadata *NodeMetadata, ctfNode *clnode.Output) (*Node, error) {
	gqlClient, gqErr := client.NewWithContext(ctx, ctfNode.Node.ExternalURL, client.Credentials{
		Email:    ctfNode.Node.APIAuthUser,
		Password: ctfNode.Node.APIAuthPassword,
	})
	if gqErr != nil {
		return nil, fmt.Errorf("failed to create node graphql client: %w", gqErr)
	}

	chainlinkClient, cErr := clclient.NewChainlinkClient(&clclient.ChainlinkConfig{
		URL:         ctfNode.Node.ExternalURL,
		Email:       ctfNode.Node.APIAuthUser,
		Password:    ctfNode.Node.APIAuthPassword,
		InternalIP:  ctfNode.Node.InternalIP,
		HTTPTimeout: ptr.Ptr(10 * time.Second),
	}, framework.L)
	if cErr != nil {
		return nil, fmt.Errorf("failed to create node rest client: %w", cErr)
	}

	node := &Node{
		Clients: NodeClients{
			GQLClient:  gqlClient,
			RestClient: chainlinkClient,
		},
		Name:  name,
		Keys:  nodeMetadata.Keys,
		Roles: nodeMetadata.Roles,
		Host:  nodeMetadata.Host,
	}

	for _, role := range node.Roles {
		switch role {
		case WorkerNode:
			// multi address is not applicable for non-bootstrap nodes
			// explicitly set it to empty string to denote that
			node.Addresses.MultiAddress = ""

			// set admin address for non-bootstrap nodes (capability registry requires non-null admin address; use arbitrary default value if node is not configured)
			node.Addresses.AdminAddress = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
		case BootstrapNode:
			// create multi address for OCR2, applicable only for bootstrap nodes
			p2pURL, err := url.Parse(ctfNode.Node.InternalP2PUrl)
			if err != nil {
				return nil, fmt.Errorf("failed to parse p2p url: %w", err)
			}
			node.Addresses.MultiAddress = fmt.Sprintf("%s:%s", ctfNode.Node.InternalIP, p2pURL.Port())

			// no need to set admin address for bootstrap nodes, as there will be no payment
			node.Addresses.AdminAddress = ""
		case GatewayNode:
			// no specific data to set for gateway nodes yet
		default:
			return nil, fmt.Errorf("unknown node role: %s", role)
		}
	}

	return node, nil
}

type JobDistributorDetails struct {
	NodeID string // node id returned by JD after node is registered with it
	JDID   string // JD id returned by node after Job distributor is created in the node
}

// Do we need to store public address per chain or is it enough to store keys, so that we can derive public address when needed?
type Addresses struct {
	AdminAddress string
	MultiAddress string
}

type NodeIDs struct {
	PeerID string
}

type NodeClients struct {
	GQLClient  client.Client             // graphql client to interact with the node
	RestClient *clclient.ChainlinkClient // rest client to interact with the node
}

type JDChainConfigInput struct {
	ChainID   string
	ChainType string
}

func (n *Node) CreateJDChainConfigs(ctx context.Context, chains []JDChainConfigInput, jd *JobDistributor) error {
	for _, chain := range chains {
		var account string

		switch strings.ToLower(chain.ChainType) {
		case chainselectors.FamilyEVM, chainselectors.FamilyTron:
			chainIDUint64, parseErr := strconv.ParseUint(chain.ChainID, 10, 64)
			if parseErr != nil {
				return fmt.Errorf("failed to parse chain id %s: %w", chain.ChainID, parseErr)
			}

			if chainIDUint64 == 0 {
				return fmt.Errorf("invalid chain id: %s", chain.ChainID)
			}

			evmKey, ok := n.Keys.EVM[chainIDUint64]
			if !ok {
				var fetchErr error
				accountAddr, fetchErr := n.Clients.GQLClient.FetchAccountAddress(ctx, chain.ChainID)
				if fetchErr != nil {
					return fmt.Errorf("failed to fetch account address for node %s: %w", n.Name, fetchErr)
				}
				if accountAddr == nil {
					return fmt.Errorf("no account address found for node %s", n.Name)
				}
				account = *accountAddr
			} else {
				account = evmKey.PublicAddress.Hex()
			}
			// accountAddr, err := n.gqlClient.FetchAccountAddress(ctx, chain.ChainID)
			// if err != nil {
			// 	return fmt.Errorf("failed to fetch account address for node %s: %w", n.Name, err)
			// }
			// if accountAddr == nil {
			// 	return fmt.Errorf("no account address found for node %s", n.Name)
			// }
			// n.AccountAddr[chain.ChainID] = *accountAddr
			// account = *accountAddr
		case chainselectors.FamilySolana:
			solKey, ok := n.Keys.Solana[chain.ChainID]
			if !ok {
				accounts, fetchErr := n.Clients.GQLClient.FetchKeys(ctx, chain.ChainType)
				if fetchErr != nil {
					return fmt.Errorf("failed to fetch account address for node %s and chain %s: %w", n.Name, chain.ChainType, fetchErr)
				}
				if len(accounts) == 0 {
					return fmt.Errorf("failed to fetch account address for node %s and chain %s", n.Name, chain.ChainType)
				}
				account = accounts[0]
			} else {
				account = solKey.PublicAddress.String()
			}
		case chainselectors.FamilyAptos:
			accounts, err := n.Clients.GQLClient.FetchKeys(ctx, chain.ChainType)
			if err != nil {
				return fmt.Errorf("failed to fetch account address for node %s and chain %s: %w", n.Name, chain.ChainType, err)
			}
			if len(accounts) == 0 {
				return fmt.Errorf("failed to fetch account address for node %s and chain %s", n.Name, chain.ChainType)
			}
			// TODO store that or not?
			// n.AccountAddr[chain.ChainID] = accounts[0]
			account = accounts[0]
		default:
			return fmt.Errorf("unsupported chainType %v", chain.ChainType)
		}

		// peerID, err := n.gqlClient.FetchP2PPeerID(ctx)
		// if err != nil {
		// 	return fmt.Errorf("failed to fetch peer id for node %s: %w", n.Name, err)
		// }
		// if peerID == nil {
		// 	return fmt.Errorf("no peer id found for node %s", n.Name)
		// }

		chainType := chain.ChainType
		if strings.EqualFold(chain.ChainType, blockchain.FamilyTron) {
			chainType = strings.ToUpper(blockchain.FamilyEVM)
		}
		ocr2BundleID, createErr := n.Clients.GQLClient.FetchOCR2KeyBundleID(ctx, chainType)
		if createErr != nil {
			return fmt.Errorf("failed to fetch OCR2 key bundle id for node %s: %w", n.Name, createErr)
		}
		if ocr2BundleID == "" {
			return fmt.Errorf("no OCR2 key bundle id found for node %s", n.Name)
		}

		if n.Keys.OCR2BundleIDs == nil {
			n.Keys.OCR2BundleIDs = make(map[string]string)
		}

		n.Keys.OCR2BundleIDs[strings.ToLower(chainType)] = ocr2BundleID

		// fetch node labels to know if the node is bootstrap or plugin
		// if multi address is set, then it's a bootstrap node
		// isBootstrap := n.HasRole(BootstrapNode)
		// for _, label := range n.labels {
		// 	if label.Key == LabelNodeTypeKey && value(label.Value) == LabelNodeTypeValueBootstrap {
		// 		isBootstrap = true
		// 		break
		// 	}
		// }

		// retry twice with 5 seconds interval to create JobDistributorChainConfig
		retryErr := retry.Do(ctx, retry.WithMaxDuration(10*time.Second, retry.NewConstant(3*time.Second)), func(ctx context.Context) error {
			// check the node chain config to see if this chain already exists
			nodeChainConfigs, listErr := jd.ListNodeChainConfigs(context.Background(), &nodev1.ListNodeChainConfigsRequest{
				Filter: &nodev1.ListNodeChainConfigsRequest_Filter{
					NodeIds: []string{n.JobDistributorDetails.NodeID},
				}})
			if listErr != nil {
				return retry.RetryableError(fmt.Errorf("failed to list node chain configs for node %s, retrying..: %w", n.Name, listErr))
			}
			if nodeChainConfigs != nil {
				for _, chainConfig := range nodeChainConfigs.ChainConfigs {
					if chainConfig.Chain.Id == chain.ChainID {
						return nil
					}
				}
			}

			// we need to create JD chain config for each chain, because later on changestes ask the node for that chain data
			// each node needs to have OCR2 enabled, because p2pIDs are used by some contracts to identify nodes (e.g. capability registry)
			_, createErr = n.Clients.GQLClient.CreateJobDistributorChainConfig(ctx, client.JobDistributorChainConfigInput{
				JobDistributorID: n.JobDistributorDetails.JDID,
				ChainID:          chain.ChainID,
				ChainType:        chainType,
				AccountAddr:      account,
				AdminAddr:        n.Addresses.AdminAddress,
				Ocr2Enabled:      true,
				Ocr2IsBootstrap:  n.HasRole(BootstrapNode),
				Ocr2Multiaddr:    n.Addresses.MultiAddress,
				Ocr2P2PPeerID:    n.Keys.P2PKey.PeerID.String(),
				Ocr2KeyBundleID:  ocr2BundleID,
				Ocr2Plugins:      `{}`,
			})
			// TODO: add a check if the chain config failed because of a duplicate in that case, should we update or return success?
			if createErr != nil {
				return fmt.Errorf("failed to create JD chain configuration for node %s: %w", n.Name, createErr)
			}

			// JD silently fails to update nodeChainConfig. Therefore, we fetch the node config and
			// if it's not updated , throw an error
			return retry.RetryableError(errors.New("retrying CreateChainConfig in JD"))
		})

		if retryErr != nil {
			return fmt.Errorf("failed to create  create JD chain configuration for node %s: %w", n.Name, retryErr)
		}
	}
	return nil
}

// AcceptJob accepts the job proposal for the given job proposal spec
func (n *Node) AcceptJob(ctx context.Context, spec string) error {
	// fetch JD to get the job proposals
	jd, err := n.Clients.GQLClient.GetJobDistributor(ctx, n.JobDistributorDetails.JDID)
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
	approvedSpec, err := n.Clients.GQLClient.ApproveJobProposalSpec(ctx, idToAccept, false)
	if err != nil {
		return err
	}
	if approvedSpec == nil {
		return fmt.Errorf("no job proposal spec found for job id %s", idToAccept)
	}
	return nil
}

// RegisterNodeToJobDistributor fetches the CSA public key of the node and registers the node with the job distributor
// it sets the node id returned by JobDistributor as a result of registration in the node struct
func (n *Node) RegisterNodeToJobDistributor(ctx context.Context, jd *JobDistributor, labels []*ptypes.Label) error {
	// Get the public key of the node
	csaKeyRes, err := n.Clients.GQLClient.FetchCSAPublicKey(ctx)
	if err != nil {
		return err
	}
	if csaKeyRes == nil {
		return fmt.Errorf("no csa key found for node %s", n.Name)
	}

	n.Keys.CSAKey = crypto.CSAKey{
		Key: *csaKeyRes,
	}

	// // tag nodes with p2p_id for easy lookup
	// peerID, err := n.gqlClient.FetchP2PPeerID(ctx)
	// if err != nil {
	// 	return fmt.Errorf("failed to fetch peer id for node %s: %w", n.Name, err)
	// }
	// if peerID == nil {
	// 	return fmt.Errorf("no peer id found for node %s", n.Name)
	// }
	labels = append(labels, &ptypes.Label{
		Key:   LabelNodeP2PIDKey,
		Value: ptr.Ptr(n.Keys.P2PKey.PeerID.String()),
	})

	// register the node in the job distributor
	registerResponse, err := jd.RegisterNode(ctx, &nodev1.RegisterNodeRequest{
		PublicKey: n.Keys.CSAKey.CleansedKey(),
		Labels:    labels,
		Name:      n.Name,
	})

	if n.JobDistributorDetails == nil {
		n.JobDistributorDetails = &JobDistributorDetails{}
	}

	// node already registered, fetch it's id
	if err != nil && strings.Contains(err.Error(), "AlreadyExists") {
		nodesResponse, listErr := jd.ListNodes(ctx, &nodev1.ListNodesRequest{
			Filter: &nodev1.ListNodesRequest_Filter{
				Selectors: []*ptypes.Selector{
					{
						Key:   LabelNodeP2PIDKey,
						Op:    ptypes.SelectorOp_EQ,
						Value: ptr.Ptr(n.Keys.P2PKey.PeerID.String()),
					},
				},
			},
		})
		if listErr != nil {
			return listErr
		}
		nodes := nodesResponse.GetNodes()
		if len(nodes) == 0 {
			return fmt.Errorf("failed to find node: %v", n.Name)
		}
		n.JobDistributorDetails.NodeID = nodes[0].Id
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to register node %s: %w", n.Name, err)
	}
	if registerResponse.GetNode().GetId() == "" {
		return fmt.Errorf("no node id returned from job distributor for node %s", n.Name)
	}
	n.JobDistributorDetails.NodeID = registerResponse.GetNode().GetId()

	return nil
}

// CreateJobDistributor fetches the keypairs from the job distributor and creates the job distributor in the node
// and returns the job distributor id
func (n *Node) CreateJobDistributor(ctx context.Context, jd *JobDistributor) (string, error) {
	// Get the keypairs from the job distributor
	csaKey, err := jd.GetCSAPublicKey(ctx)
	if err != nil {
		return "", err
	}
	// create the job distributor in the node with the csa key
	resp, err := n.Clients.GQLClient.ListJobDistributors(ctx)
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
	return n.Clients.GQLClient.CreateJobDistributor(ctx, client.JobDistributorInput{
		Name:      "Job Distributor",
		Uri:       jd.WSRPC,
		PublicKey: csaKey,
	})
}

// SetUpAndLinkJobDistributor sets up the job distributor in the node and registers the node with the job distributor
// it sets the job distributor id for node
func (n *Node) SetUpAndLinkJobDistributor(ctx context.Context, jd *JobDistributor, labels []*ptypes.Label) error {
	// register the node in the job distributor
	err := n.RegisterNodeToJobDistributor(ctx, jd, labels)
	if err != nil {
		return err
	}
	// now create the job distributor in the node
	id, err := n.CreateJobDistributor(ctx, jd)
	if err != nil &&
		!strings.Contains(err.Error(), "DuplicateFeedsManagerError") {
		return fmt.Errorf("failed to create job distributor in node %s: %w", n.Name, err)
	}
	// wait for the node to connect to the job distributor
	err = retry.Do(ctx, retry.WithMaxDuration(1*time.Minute, retry.NewFibonacci(1*time.Second)), func(ctx context.Context) error {
		getRes, getErr := jd.GetNode(ctx, &nodev1.GetNodeRequest{
			Id: n.JobDistributorDetails.NodeID,
		})
		if getErr != nil {
			return retry.RetryableError(fmt.Errorf("failed to get node %s: %w", n.Name, getErr))
		}
		if getRes.GetNode() == nil {
			return fmt.Errorf("no node found for node id %s", n.JobDistributorDetails.NodeID)
		}
		if !getRes.GetNode().IsConnected {
			return retry.RetryableError(fmt.Errorf("node %s not connected to job distributor", n.Name))
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to connect node %s to job distributor: %w", n.Name, err)
	}
	n.JobDistributorDetails.JDID = id
	return nil
}

// func ptr[T any](v T) *T {
// 	return &v
// }

// func value[T any](v *T) T {
// 	zero := new(T)
// 	if v == nil {
// 		return *zero
// 	}
// 	return *v
// }

func (n *Node) CancelProposalsByExternalJobID(ctx context.Context, externalJobIDs []string) ([]string, error) {
	jd, err := n.Clients.GQLClient.GetJobDistributor(ctx, n.JobDistributorDetails.JDID)
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
		spec, err := n.Clients.GQLClient.CancelJobProposalSpec(ctx, jp.Id)
		if err != nil {
			return nil, err
		}

		if spec == nil {
			return nil, fmt.Errorf("no job proposal spec found for id %s", jp.Id)
		}
	}

	return proposalIDs, nil
}

func (n *Node) ApproveProposals(ctx context.Context, proposalIDs []string) error {
	for _, proposalID := range proposalIDs {
		spec, err := n.Clients.GQLClient.ApproveJobProposalSpec(ctx, proposalID, false)
		if err != nil {
			return err
		}
		if spec == nil {
			return fmt.Errorf("no job proposal spec found for id %s", proposalID)
		}
	}
	return nil
}

func LinkToJobDistributor(ctx context.Context, input *LinkDonsToJDInput) (*cldf.Environment, []*DON, error) {
	if input == nil {
		return nil, nil, errors.New("input is nil")
	}
	if err := input.Validate(); err != nil {
		return nil, nil, errors.Wrap(err, "input validation failed")
	}

	jdConfig := jd.JDConfig{
		GRPC:  input.JdOutput.ExternalGRPCUrl,
		WSRPC: input.JdOutput.InternalWSRPCUrl,
		Creds: insecure.NewCredentials(),
	}

	jdClient, jdErr := jd.NewJDClient(jdConfig)
	if jdErr != nil {
		return nil, nil, errors.Wrap(jdErr, "failed to create JD client")
	}

	donJDClient := &JobDistributor{
		JobDistributor: jdClient,
	}

	dons := make([]*DON, len(input.NodeSetOutput))
	// var allNodesInfo []NodeInfo

	for idx, nodeOutput := range input.NodeSetOutput {
		// a maximum of 1 bootstrap is supported due to environment constraints
		// bootstrapNodeCount := 0
		// if input.Topology.DonsMetadata.List()[idx].ContainsBootstrapNode() {
		// 	bootstrapNodeCount = 1
		// }

		// nodeInfo, err := node.GetNodeInfo(nodeOutput.Output, nodeOutput.NodeSetName, input.Topology.DonsMetadata.List()[idx].ID, bootstrapNodeCount)
		// if err != nil {
		// 	return nil, nil, errors.Wrap(err, "failed to get node info")
		// }
		// allNodesInfo = append(allNodesInfo, nodeInfo...)

		supportedChains, schErr := findSupportedChainsForDON(input.Topology.DonsMetadata.List()[idx], input.BlockchainOutputs)
		if schErr != nil {
			return nil, nil, errors.Wrap(schErr, "failed to find supported chains for DON")
		}

		don, regErr := NewDON(ctx, input.Topology.DonsMetadata.List()[idx], nodeOutput.CLNodes, supportedChains, donJDClient)
		if regErr != nil {
			return nil, nil, fmt.Errorf("failed to create registered DON: %w", regErr)
		}

		dons[idx] = don

		// var regErr error
		// dons[idx], regErr = configureJDForDON(ctx, nodeInfo, supportedChains, input.JdOutput)
		// if regErr != nil {
		// 	return nil, nil, fmt.Errorf("failed to configure JD for DON: %w", regErr)
		// }
	}

	var nodeIDs []string
	for _, don := range dons {
		nodeIDs = append(nodeIDs, don.JDNodeIDs()...)
	}

	// dons = addOCRKeyLabelsToNodeMetadata(dons, input.Topology)

	// ctxWithTimeout, cancel := context.WithTimeout(ctx, 2*time.Minute)
	// defer cancel()

	// jd, jdErr := NewJDClient(ctxWithTimeout, JDConfig{
	// 	GRPC:     input.JdOutput.ExternalGRPCUrl,
	// 	WSRPC:    input.JdOutput.InternalWSRPCUrl,
	// 	Creds:    insecure.NewCredentials(),
	// 	NodeInfo: allNodesInfo,
	// })

	// if jdErr != nil {
	// 	return nil, nil, errors.Wrap(jdErr, "failed to create JD client")
	// }

	input.CldfEnvironment.Offchain = donJDClient
	input.CldfEnvironment.NodeIDs = nodeIDs

	return input.CldfEnvironment, dons, nil
}

// copied from flags package to avoid circular dependency
func HasFlagForAnyChain(values []string, capability string) bool {
	if slices.Contains(values, capability) {
		return true
	}

	for _, value := range values {
		if strings.HasPrefix(value, capability+"-") {
			return true
		}
	}

	return false
}

func findSupportedChainsForDON(donMetadata *DonMetadata, blockchainOutputs []*WrappedBlockchainOutput) ([]ChainConfig, error) {
	chains := make([]ChainConfig, 0)

	for chainSelector, bcOut := range blockchainOutputs {
		hasEVMChainEnabled := slices.Contains(donMetadata.EVMChains(), bcOut.ChainID)
		hasSolanaWriteCapability := HasFlagForAnyChain(donMetadata.Flags, WriteSolanaCapability)
		chainIsSolana := strings.EqualFold(bcOut.BlockchainOutput.Family, chainselectors.FamilySolana)
		if !hasEVMChainEnabled && (!hasSolanaWriteCapability || !chainIsSolana) {
			continue
		}

		cfg, cfgErr := ChainConfigFromWrapped(bcOut)
		if cfgErr != nil {
			return nil, errors.Wrapf(cfgErr, "failed to build chain config for chain selector %d", chainSelector)
		}

		chains = append(chains, cfg)
	}

	return chains, nil
}

type JDConfig struct {
	GRPC  string
	WSRPC string
	Creds credentials.TransportCredentials
	Auth  oauth2.TokenSource
	// NodeInfo []NodeInfo
}

// JobDistributor implements the OffchainClient interface in CLDF and wraps the CLDF JD client and add DON functionality.
// The CLDF JD client does not have the DON functionality, so we wrap it here.
type JobDistributor struct {
	*jd.JobDistributor
	don *DON
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
