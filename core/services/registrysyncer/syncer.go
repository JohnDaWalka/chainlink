package registrysyncer

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	p2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

type Listener interface {
	OnNewRegistry(ctx context.Context, registry *LocalRegistry) error
}

type Syncer interface {
	services.Service
	AddListener(h ...Listener)
}

type ContractReaderFactory interface {
	NewContractReader(context.Context, []byte) (types.ContractReader, error)
}

type RegistrySyncer interface {
	Sync(ctx context.Context, isInitialSync bool) error
	AddListener(listeners ...Listener)
	Start(ctx context.Context) error
	Close() error
	Ready() error
	HealthReport() map[string]error
	Name() string
}

type registrySyncer struct {
	services.StateMachine
	metrics              *syncerMetricLabeler
	stopCh               services.StopChan
	listeners            []Listener
	relayer              ContractReaderFactory
	capabilitiesContract types.BoundContract
	getPeerID            func() (p2ptypes.PeerID, error)

	orm ORM

	updateChan chan *LocalRegistry

	capabilitiesRegistryReader CapabilitiesRegistryReader
	readerFactory              CapabilitiesRegistryReaderFactory

	wg   sync.WaitGroup
	lggr logger.Logger
	mu   sync.RWMutex
}

const capabilitiesRegistryContractName = "CapabilitiesRegistry"

var _ services.Service = &registrySyncer{}

var (
	defaultTickInterval = 12 * time.Second
)

// New instantiates a new RegistrySyncer
func New(
	lggr logger.Logger,
	getPeerID func() (p2ptypes.PeerID, error),
	relayer ContractReaderFactory,
	registryAddress string,
	orm ORM,
) (RegistrySyncer, error) {
	metricLabeler, err := newSyncerMetricLabeler()
	if err != nil {
		return nil, fmt.Errorf("failed to create syncer metric labeler: %w", err)
	}

	return &registrySyncer{
		metrics:    metricLabeler,
		stopCh:     make(services.StopChan),
		updateChan: make(chan *LocalRegistry),
		lggr:       logger.Named(lggr, "RegistrySyncer"),
		relayer:    relayer,
		capabilitiesContract: types.BoundContract{
			Address: registryAddress,
			Name:    capabilitiesRegistryContractName,
		},
		orm:           orm,
		getPeerID:     getPeerID,
		readerFactory: NewCapabilitiesRegistryReaderFactory(),
	}, nil
}

func (s *registrySyncer) Start(ctx context.Context) error {
	return s.StartOnce("RegistrySyncer", func() error {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.syncLoop()
		}()
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.updateStateLoop()
		}()
		return nil
	})
}

func (s *registrySyncer) syncLoop() {
	ctx, cancel := s.stopCh.NewCtx()
	defer cancel()

	ticker := time.NewTicker(defaultTickInterval)
	defer ticker.Stop()

	// Sync for a first time outside the loop; this means we'll start a remote
	// sync immediately once spinning up syncLoop, as by default a ticker will
	// fire for the first time at T+N, where N is the interval. We do not
	// increment RemoteRegistryFailureCounter the first time
	s.lggr.Debug("starting initial sync with remote registry")
	err := s.Sync(ctx, true)
	if err != nil {
		s.lggr.Errorw("failed to sync with remote registry", "error", err)
	}

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.lggr.Debug("starting regular sync with the remote registry")
			err := s.Sync(ctx, false)
			if err != nil {
				s.lggr.Errorw("failed to sync with remote registry", "error", err)
				s.metrics.incrementRemoteRegistryFailureCounter(ctx)
			}
		}
	}
}

func (s *registrySyncer) updateStateLoop() {
	ctx, cancel := s.stopCh.NewCtx()
	defer cancel()

	for {
		select {
		case <-s.stopCh:
			return
		case localRegistry, ok := <-s.updateChan:
			if !ok {
				// channel has been closed, terminating.
				return
			}
			if err := s.orm.AddLocalRegistry(ctx, *localRegistry); err != nil {
				s.lggr.Errorw("failed to save state to local registry", "error", err)
			}
		}
	}
}

func (s *registrySyncer) importOnchainRegistry(ctx context.Context) (*LocalRegistry, error) {
	// Create versioned reader if not already created
	// NOTE: This can't be called during syncer initialization because Bind() makes an onchain call
	// to verify that the contract address exists, and if called during initialization,
	// this results in a "no live nodes" error.
	if s.capabilitiesRegistryReader == nil {
		reader, err := s.readerFactory.NewCapabilitiesRegistryReader(
			ctx,
			s.relayer,
			s.capabilitiesContract.Address,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create capabilities registry reader: %w", err)
		}
		s.capabilitiesRegistryReader = reader
	}

	// Use versioned reader to get capabilities
	capabilityInfos, err := s.capabilitiesRegistryReader.GetCapabilities(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get capabilities: %w", err)
	}

	idsToCapabilities := map[string]Capability{}
	hashedIDsToCapabilityIDs := map[[32]byte]string{}
	for _, c := range capabilityInfos {
		idsToCapabilities[c.ID] = Capability{
			ID:             c.ID,
			CapabilityType: toCapabilityType(c.CapabilityType),
		}

		// V1-specific: build hash mapping from capabilities
		if c.HashedId != nil {
			hashedIDsToCapabilityIDs[*c.HashedId] = c.ID
		}
	}

	// Use versioned reader to get DONs
	donInfos, err := s.capabilitiesRegistryReader.GetDONs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get DONs: %w", err)
	}

	idsToDONs := map[DonID]DON{}
	for _, d := range donInfos {
		cc := map[string]CapabilityConfiguration{}
		for _, dc := range d.CapabilityConfigurations {
			// Handle V1 case where CapabilityId is a pointer to [32]byte
			var capabilityID string

			if dc.CapabilityId != nil {
				// This is V1 - we have the hash directly
				hash := *dc.CapabilityId

				cid, ok := hashedIDsToCapabilityIDs[hash]
				if !ok {
					return nil, fmt.Errorf("invariant violation: could not find capability ID for hashed ID %x", hash)
				}
				capabilityID = cid
			} else {
				// V2 case: capability ID would be a string field (not implemented yet)
				// For now, skip this case
				continue
			}

			cc[capabilityID] = CapabilityConfiguration{
				Config: dc.Config,
			}
		}

		idsToDONs[DonID(d.ID)] = DON{
			DON:                      *toDONInfoFromVersioned(d),
			CapabilityConfigurations: cc,
		}
	}

	// Use versioned reader to get nodes
	nodeInfos, err := s.capabilitiesRegistryReader.GetNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %w", err)
	}

	idsToNodes := map[p2ptypes.PeerID]kcr.INodeInfoProviderNodeInfo{}
	for _, node := range nodeInfos {
		// Convert versioned NodeInfo back to V1 format for compatibility
		v1Node := kcr.INodeInfoProviderNodeInfo{
			NodeOperatorId:      node.NodeOperatorID,
			ConfigCount:         node.ConfigCount,
			WorkflowDONId:       node.WorkflowDONId,
			Signer:              node.Signer,
			P2pId:               [32]byte(node.P2PID), // Direct conversion from PeerID to [32]byte
			EncryptionPublicKey: node.EncryptionPublicKey,
			HashedCapabilityIds: *node.HashedCapabilityIds, // Dereference the pointer
			CapabilitiesDONIds:  make([]*big.Int, len(node.CapabilitiesDONIds)),
		}

		// Convert uint32 slice to big.Int slice
		for i, donID := range node.CapabilitiesDONIds {
			v1Node.CapabilitiesDONIds[i] = big.NewInt(int64(donID))
		}

		idsToNodes[node.P2PID] = v1Node
	}

	return &LocalRegistry{
		lggr:              s.lggr,
		getPeerID:         s.getPeerID,
		IDsToDONs:         idsToDONs,
		IDsToCapabilities: idsToCapabilities,
		IDsToNodes:        idsToNodes,
	}, nil
}

func (s *registrySyncer) Sync(ctx context.Context, isInitialSync bool) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.listeners) == 0 {
		s.lggr.Warn("sync called, but no listeners are registered; nooping")
		return nil
	}

	var latestRegistry *LocalRegistry
	var err error

	if isInitialSync {
		s.lggr.Debug("syncing with local registry")
		latestRegistry, err = s.orm.LatestLocalRegistry(ctx)
		if err != nil {
			s.lggr.Warnw("failed to sync with local registry, using remote registry instead", "error", err)
		} else {
			latestRegistry.lggr = s.lggr
			latestRegistry.getPeerID = s.getPeerID
		}
	}

	if latestRegistry == nil {
		s.lggr.Debug("syncing with remote registry")
		importedRegistry, err := s.importOnchainRegistry(ctx)
		if err != nil {
			return fmt.Errorf("failed to sync with remote registry: %w", err)
		}
		latestRegistry = importedRegistry
		// Attempt to send local registry to the update channel without blocking
		// This is to prevent the tests from hanging if they are not calling `Start()` on the syncer
		select {
		case <-s.stopCh:
			s.lggr.Debug("sync cancelled, stopping")
		case s.updateChan <- latestRegistry:
			// Successfully sent state
			s.lggr.Debug("remote registry update triggered successfully")
		default:
			// No one is ready to receive the state, handle accordingly
			s.lggr.Debug("no listeners on update channel, remote registry update skipped")
		}
	}

	for _, listener := range s.listeners {
		lrCopy := deepCopyLocalRegistry(latestRegistry)
		if err := listener.OnNewRegistry(ctx, &lrCopy); err != nil {
			s.lggr.Errorf("error calling launcher: %s", err)
			s.metrics.incrementLauncherFailureCounter(ctx)
		}
	}

	return nil
}

func deepCopyLocalRegistry(lr *LocalRegistry) LocalRegistry {
	var lrCopy LocalRegistry
	lrCopy.lggr = lr.lggr
	lrCopy.getPeerID = lr.getPeerID
	lrCopy.IDsToDONs = make(map[DonID]DON, len(lr.IDsToDONs))
	for id, don := range lr.IDsToDONs {
		d := capabilities.DON{
			ID:               don.ID,
			ConfigVersion:    don.ConfigVersion,
			Members:          make([]p2ptypes.PeerID, len(don.Members)),
			F:                don.F,
			IsPublic:         don.IsPublic,
			AcceptsWorkflows: don.AcceptsWorkflows,
		}
		copy(d.Members, don.Members)
		capCfgs := make(map[string]CapabilityConfiguration, len(don.CapabilityConfigurations))
		for capID, capCfg := range don.CapabilityConfigurations {
			capCfgs[capID] = CapabilityConfiguration{
				Config: capCfg.Config,
			}
		}
		lrCopy.IDsToDONs[id] = DON{
			DON:                      d,
			CapabilityConfigurations: capCfgs,
		}
	}

	lrCopy.IDsToCapabilities = make(map[string]Capability, len(lr.IDsToCapabilities))
	for id, capability := range lr.IDsToCapabilities {
		cp := capability
		lrCopy.IDsToCapabilities[id] = cp
	}

	lrCopy.IDsToNodes = make(map[p2ptypes.PeerID]kcr.INodeInfoProviderNodeInfo, len(lr.IDsToNodes))
	for id, node := range lr.IDsToNodes {
		nodeInfo := kcr.INodeInfoProviderNodeInfo{
			NodeOperatorId:      node.NodeOperatorId,
			ConfigCount:         node.ConfigCount,
			WorkflowDONId:       node.WorkflowDONId,
			Signer:              node.Signer,
			P2pId:               node.P2pId,
			EncryptionPublicKey: node.EncryptionPublicKey,
			HashedCapabilityIds: make([][32]byte, len(node.HashedCapabilityIds)),
			CapabilitiesDONIds:  make([]*big.Int, len(node.CapabilitiesDONIds)),
		}
		copy(nodeInfo.HashedCapabilityIds, node.HashedCapabilityIds)
		copy(nodeInfo.CapabilitiesDONIds, node.CapabilitiesDONIds)
		lrCopy.IDsToNodes[id] = nodeInfo
	}

	return lrCopy
}

type ContractCapabilityType uint8

const (
	ContractCapabilityTypeTrigger ContractCapabilityType = iota
	ContractCapabilityTypeAction
	ContractCapabilityTypeConsensus
	ContractCapabilityTypeTarget
)

func toCapabilityType(capabilityType uint8) capabilities.CapabilityType {
	switch ContractCapabilityType(capabilityType) {
	case ContractCapabilityTypeTrigger:
		return capabilities.CapabilityTypeTrigger
	case ContractCapabilityTypeAction:
		return capabilities.CapabilityTypeAction
	case ContractCapabilityTypeConsensus:
		return capabilities.CapabilityTypeConsensus
	case ContractCapabilityTypeTarget:
		return capabilities.CapabilityTypeTarget
	default:
		return capabilities.CapabilityTypeUnknown
	}
}

func toDONInfo(don kcr.CapabilitiesRegistryDONInfo) *capabilities.DON {
	peerIDs := []p2ptypes.PeerID{}
	for _, p := range don.NodeP2PIds {
		peerIDs = append(peerIDs, p)
	}

	return &capabilities.DON{
		ID:               don.Id,
		ConfigVersion:    don.ConfigCount,
		Members:          peerIDs,
		F:                don.F,
		IsPublic:         don.IsPublic,
		AcceptsWorkflows: don.AcceptsWorkflows,
	}
}

func toDONInfoFromVersioned(don DONInfo) *capabilities.DON {
	peerIDs := []p2ptypes.PeerID{}
	for _, p := range don.NodeP2PIds {
		peerIDs = append(peerIDs, p)
	}

	return &capabilities.DON{
		ID:               don.ID,
		ConfigVersion:    don.ConfigCount,
		Members:          peerIDs,
		F:                don.F,
		IsPublic:         don.IsPublic,
		AcceptsWorkflows: don.AcceptsWorkflows,
	}
}

func (s *registrySyncer) AddListener(listeners ...Listener) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.listeners = append(s.listeners, listeners...)
}

func (s *registrySyncer) Close() error {
	return s.StopOnce("RegistrySyncer", func() error {
		close(s.stopCh)
		s.mu.Lock()
		defer s.mu.Unlock()
		close(s.updateChan)
		s.wg.Wait()
		return nil
	})
}

func (s *registrySyncer) HealthReport() map[string]error {
	return map[string]error{s.Name(): s.Healthy()}
}

func (s *registrySyncer) Name() string {
	return s.lggr.Name()
}
