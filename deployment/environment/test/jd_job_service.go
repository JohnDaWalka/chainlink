package test

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/grpc"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink/deployment"
)

type UnimplementedJobServiceClient struct{}

func (s *UnimplementedJobServiceClient) BatchProposeJob(ctx context.Context, in *jobv1.BatchProposeJobRequest, opts ...grpc.CallOption) (*jobv1.BatchProposeJobResponse, error) {
	// TODO CCIP-3108  implement me
	panic("implement me")
}

func (s *UnimplementedJobServiceClient) DeleteJob(ctx context.Context, in *jobv1.DeleteJobRequest, opts ...grpc.CallOption) (*jobv1.DeleteJobResponse, error) {
	panic("unimplemented")
}

func (s *UnimplementedJobServiceClient) UpdateJob(ctx context.Context, in *jobv1.UpdateJobRequest, opts ...grpc.CallOption) (*jobv1.UpdateJobResponse, error) {
	panic("unimplemented")
}

// GetJob implements job.JobServiceClient.
func (s *UnimplementedJobServiceClient) GetJob(ctx context.Context, in *jobv1.GetJobRequest, opts ...grpc.CallOption) (*jobv1.GetJobResponse, error) {
	panic("unimplemented")
}

// GetProposal implements job.JobServiceClient.
func (s *UnimplementedJobServiceClient) GetProposal(ctx context.Context, in *jobv1.GetProposalRequest, opts ...grpc.CallOption) (*jobv1.GetProposalResponse, error) {
	panic("unimplemented")
}

// ListJobs implements job.JobServiceClient.
func (s *UnimplementedJobServiceClient) ListJobs(ctx context.Context, in *jobv1.ListJobsRequest, opts ...grpc.CallOption) (*jobv1.ListJobsResponse, error) {
	panic("unimplemented")
}

// ListProposals implements job.JobServiceClient.
func (s *UnimplementedJobServiceClient) ListProposals(ctx context.Context, in *jobv1.ListProposalsRequest, opts ...grpc.CallOption) (*jobv1.ListProposalsResponse, error) {
	panic("unimplemented")
}

// ProposeJob implements job.JobServiceClient.
func (s *UnimplementedJobServiceClient) ProposeJob(ctx context.Context, in *jobv1.ProposeJobRequest, opts ...grpc.CallOption) (*jobv1.ProposeJobResponse, error) {
	panic("unimplemented")
}

// RevokeJob implements job.JobServiceClient.
func (s *UnimplementedJobServiceClient) RevokeJob(ctx context.Context, in *jobv1.RevokeJobRequest, opts ...grpc.CallOption) (*jobv1.RevokeJobResponse, error) {
	panic("unimplemented")
}

type jobState struct {
	jobs      []*jobv1.Job
	proposals []*jobv1.Proposal
}

// jobStore is a thread-safe jobStore for wrappedNode
// it is indexed by both p2p key and csa key
type jobStore struct {
	mu  sync.RWMutex
	db2 map[string]*wrappedNode

	p2pToID map[p2pKey]string
	csaToID map[csaKey]string
}

func newJobStore(node []deployment.Node) *jobStore {
	s := &jobStore{
		db2:     make(map[string]*wrappedNode),
		csaToID: make(map[csaKey]string),
		p2pToID: make(map[p2pKey]string),
	}
	for _, v := range node {
		w := newWrapper(v)
		s.db2[v.NodeID] = w
		s.p2pToID[p2pKey(w.Node.PeerID.String())] = v.NodeID
		s.csaToID[w.Node.CSAKey] = v.NodeID
	}
	return s
}

func (s *jobStore) getNode(id string) (*wrappedNode, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.db2[id]
	if !ok {
		return nil, fmt.Errorf("node not found for id %s", id)
	}
	return n, nil
}

func (s *jobStore) getNodeByP2P(p2p p2pKey) (*wrappedNode, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.p2pToID[p2p]
	if !ok {
		return nil, fmt.Errorf("node not found for p2p %s", p2p)
	}
	return s.getNode(id)
}

func (s *jobStore) getNodeByCSA(csa csaKey) (*wrappedNode, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.csaToID[csa]
	if !ok {
		return nil, fmt.Errorf("node not found for csa key %s", csa)
	}
	return s.getNode(id)
}

func (s *jobStore) list() []*wrappedNode {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []*wrappedNode
	for _, v := range s.db2 {
		out = append(out, v)
	}
	return out
}

func (s *jobStore) put(n *wrappedNode) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.db2[n.Node.NodeID] = n
	s.csaToID[n.Node.CSAKey] = n.NodeID
	s.p2pToID[p2pKey(n.Node.PeerID.String())] = n.NodeID
}
