package registrysyncer_test

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	p2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	kcr_v1 "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	kcr_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	evmclient "github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads/headstest"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmtestutils "github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
	syncerMocks "github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"

	captestutils "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/capabilities/testutils"
)

var writeChainCapability = kcr_v1.CapabilitiesRegistryCapability{
	LabelledName:   "write-chain",
	Version:        "1.0.1",
	CapabilityType: uint8(3),
}

func startNewChainWithRegistry(t *testing.T) (*kcr_v1.CapabilitiesRegistry, common.Address, *bind.TransactOpts, *simulated.Backend) {
	owner := evmtestutils.MustNewSimTransactor(t)

	i := &big.Int{}
	oneEth, _ := i.SetString("100000000000000000000", 10)
	gasLimit := ethconfig.Defaults.Miner.GasCeil * 2 // 60 M blocks

	simulatedBackend := simulated.NewBackend(gethtypes.GenesisAlloc{owner.From: {
		Balance: oneEth,
	}}, simulated.WithBlockGasLimit(gasLimit))
	simulatedBackend.Commit()

	CapabilitiesRegistryAddress, _, CapabilitiesRegistry, err := kcr_v1.DeployCapabilitiesRegistry(owner, simulatedBackend.Client())
	require.NoError(t, err, "DeployCapabilitiesRegistry failed")

	fmt.Println("Deployed CapabilitiesRegistry at", CapabilitiesRegistryAddress.Hex())
	simulatedBackend.Commit()

	return CapabilitiesRegistry, CapabilitiesRegistryAddress, owner, simulatedBackend
}

type crFactory struct {
	lggr      logger.Logger
	ht        logpoller.HeadTracker
	logPoller logpoller.LogPoller
	client    evmclient.Client
}

func (c *crFactory) NewContractReader(ctx context.Context, cfg []byte) (types.ContractReader, error) {
	crCfg := &evmrelaytypes.ChainReaderConfig{}
	if err := json.Unmarshal(cfg, crCfg); err != nil {
		return nil, err
	}

	svc, err := evm.NewChainReaderService(ctx, c.lggr, c.logPoller, c.ht, c.client, *crCfg)
	if err != nil {
		return nil, err
	}

	return svc, nil
}

func newContractReaderFactory(t *testing.T, simulatedBackend *simulated.Backend) *crFactory {
	lggr := logger.TestLogger(t)
	client := evmclient.NewSimulatedBackendClient(
		t,
		simulatedBackend,
		testutils.SimulatedChainID,
	)
	db := pgtest.NewSqlxDB(t)
	const finalityDepth = 2
	ht := headstest.NewSimulatedHeadTracker(client, false, finalityDepth)
	lp := logpoller.NewLogPoller(
		logpoller.NewORM(testutils.SimulatedChainID, db, lggr),
		client,
		lggr,
		ht,
		logpoller.Opts{
			PollPeriod:               100 * time.Millisecond,
			FinalityDepth:            finalityDepth,
			BackfillBatchSize:        3,
			RPCBatchSize:             2,
			KeepFinalizedBlocksDepth: 1000,
		},
	)
	return &crFactory{
		lggr:      lggr,
		client:    client,
		ht:        ht,
		logPoller: lp,
	}
}

func randomWord() [32]byte {
	word := make([]byte, 32)
	_, err := rand.Read(word)
	if err != nil {
		panic(err)
	}
	return [32]byte(word)
}

type launcher struct {
	localRegistry *registrysyncer.LocalRegistry
	mu            sync.RWMutex
}

func (l *launcher) OnNewRegistry(_ context.Context, localRegistry *registrysyncer.LocalRegistry) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.localRegistry = localRegistry
	return nil
}

type orm struct {
	ormMock               *syncerMocks.ORM
	mu                    sync.RWMutex
	latestLocalRegistryCh chan struct{}
	addLocalRegistryCh    chan struct{}
}

func newORM(t *testing.T) *orm {
	t.Helper()

	return &orm{
		ormMock:               syncerMocks.NewORM(t),
		latestLocalRegistryCh: make(chan struct{}, 1),
		addLocalRegistryCh:    make(chan struct{}, 1),
	}
}

func (o *orm) Cleanup() {
	o.mu.Lock()
	defer o.mu.Unlock()
	close(o.latestLocalRegistryCh)
	close(o.addLocalRegistryCh)
}

func (o *orm) AddLocalRegistry(ctx context.Context, localRegistry registrysyncer.LocalRegistry) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.addLocalRegistryCh <- struct{}{}
	err := o.ormMock.AddLocalRegistry(ctx, localRegistry)
	return err
}

func (o *orm) LatestLocalRegistry(ctx context.Context) (*registrysyncer.LocalRegistry, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.latestLocalRegistryCh <- struct{}{}
	return o.ormMock.LatestLocalRegistry(ctx)
}

func toPeerIDs(ids [][32]byte) []p2ptypes.PeerID {
	var pids []p2ptypes.PeerID
	for _, id := range ids {
		pids = append(pids, id)
	}
	return pids
}

func TestReader_Integration(t *testing.T) {
	ctx := testutils.Context(t)
	reg, regAddress, owner, sim := startNewChainWithRegistry(t)

	_, err := reg.AddCapabilities(owner, []kcr_v1.CapabilitiesRegistryCapability{writeChainCapability})
	require.NoError(t, err, "AddCapability failed for %s", writeChainCapability.LabelledName)
	sim.Commit()

	cid := fmt.Sprintf("%s@%s", writeChainCapability.LabelledName, writeChainCapability.Version)

	hid, err := reg.GetHashedCapabilityId(&bind.CallOpts{}, writeChainCapability.LabelledName, writeChainCapability.Version)
	require.NoError(t, err)

	_, err = reg.AddNodeOperators(owner, []kcr_v1.CapabilitiesRegistryNodeOperator{
		{
			Admin: owner.From,
			Name:  "TEST_NOP",
		},
	})
	require.NoError(t, err)
	sim.Commit()

	nodeSet := [][32]byte{
		randomWord(),
		randomWord(),
		randomWord(),
	}

	signersSet := [][32]byte{
		randomWord(),
		randomWord(),
		randomWord(),
	}

	encPubKey1 := randomWord()
	encPubKey2 := randomWord()
	encPubKey3 := randomWord()

	nodes := []kcr_v1.CapabilitiesRegistryNodeParams{
		{
			// The first NodeOperatorId has id 1 since the id is auto-incrementing.
			NodeOperatorId:      uint32(1),
			Signer:              signersSet[0],
			P2pId:               nodeSet[0],
			EncryptionPublicKey: encPubKey1,
			HashedCapabilityIds: [][32]byte{hid},
		},
		{
			// The first NodeOperatorId has id 1 since the id is auto-incrementing.
			NodeOperatorId:      uint32(1),
			Signer:              signersSet[1],
			P2pId:               nodeSet[1],
			EncryptionPublicKey: encPubKey2,
			HashedCapabilityIds: [][32]byte{hid},
		},
		{
			// The first NodeOperatorId has id 1 since the id is auto-incrementing.
			NodeOperatorId:      uint32(1),
			Signer:              signersSet[2],
			P2pId:               nodeSet[2],
			EncryptionPublicKey: encPubKey3,
			HashedCapabilityIds: [][32]byte{hid},
		},
	}
	_, err = reg.AddNodes(owner, nodes)
	require.NoError(t, err)
	sim.Commit()

	config := &capabilitiespb.CapabilityConfig{
		DefaultConfig: values.Proto(values.EmptyMap()).GetMapValue(),
		RemoteConfig: &capabilitiespb.CapabilityConfig_RemoteTriggerConfig{
			RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
				RegistrationRefresh: durationpb.New(20 * time.Second),
				RegistrationExpiry:  durationpb.New(60 * time.Second),
				// F + 1
				MinResponsesToAggregate: uint32(1) + 1,
				MessageExpiry:           durationpb.New(120 * time.Second),
			},
		},
	}
	configb, err := proto.Marshal(config)
	if err != nil {
		panic(err)
	}

	cfgs := []kcr_v1.CapabilitiesRegistryCapabilityConfiguration{
		{
			CapabilityId: hid,
			Config:       configb,
		},
	}
	_, err = reg.AddDON(
		owner,
		nodeSet,
		cfgs,
		true,
		true,
		1,
	)
	sim.Commit()

	require.NoError(t, err)

	db := pgtest.NewSqlxDB(t)
	factory := newContractReaderFactory(t, sim)
	syncerORM := registrysyncer.NewORM(db, logger.TestLogger(t))
	syncer, err := registrysyncer.New(logger.TestLogger(t), func() (p2ptypes.PeerID, error) { return p2ptypes.PeerID{}, nil }, factory, regAddress.Hex(), syncerORM)
	require.NoError(t, err)

	l := &launcher{}
	syncer.AddListener(l)

	err = syncer.Sync(ctx, false) // not looking to load from the DB in this specific test.
	s := l.localRegistry
	require.NoError(t, err)
	assert.Len(t, s.IDsToCapabilities, 1)

	gotCap := s.IDsToCapabilities[cid]
	assert.Equal(t, registrysyncer.Capability{
		CapabilityType: capabilities.CapabilityTypeTarget,
		ID:             "write-chain@1.0.1",
	}, gotCap)

	assert.Len(t, s.IDsToDONs, 1)
	expectedDON := capabilities.DON{
		ID:               1,
		ConfigVersion:    1,
		IsPublic:         true,
		AcceptsWorkflows: true,
		F:                1,
		Members:          toPeerIDs(nodeSet),
	}
	gotDon := s.IDsToDONs[1]
	assert.Equal(t, expectedDON, gotDon.DON)
	assert.Equal(t, configb, gotDon.CapabilityConfigurations[cid].Config)

	nodesInfo := []registrysyncer.NodeInfo{
		{
			NodeOperatorID:      uint32(1),
			ConfigCount:         1,
			WorkflowDONId:       1,
			Signer:              signersSet[0],
			P2PID:               p2ptypes.PeerID(nodeSet[0]),
			EncryptionPublicKey: encPubKey1,
			HashedCapabilityIDs: &[][32]byte{hid},
			CapabilitiesDONIds:  []*big.Int{},
			Version:             "v1",
		},
		{
			NodeOperatorID:      uint32(1),
			ConfigCount:         1,
			WorkflowDONId:       1,
			Signer:              signersSet[1],
			P2PID:               p2ptypes.PeerID(nodeSet[1]),
			EncryptionPublicKey: encPubKey2,
			HashedCapabilityIDs: &[][32]byte{hid},
			CapabilitiesDONIds:  []*big.Int{},
			Version:             "v1",
		},
		{
			NodeOperatorID:      uint32(1),
			ConfigCount:         1,
			WorkflowDONId:       1,
			Signer:              signersSet[2],
			P2PID:               p2ptypes.PeerID(nodeSet[2]),
			EncryptionPublicKey: encPubKey3,
			HashedCapabilityIDs: &[][32]byte{hid},
			CapabilitiesDONIds:  []*big.Int{},
			Version:             "v1",
		},
	}

	assert.Len(t, s.IDsToNodes, 3)
	assert.Equal(t, map[p2ptypes.PeerID]registrysyncer.NodeInfo{
		nodeSet[0]: nodesInfo[0],
		nodeSet[1]: nodesInfo[1],
		nodeSet[2]: nodesInfo[2],
	}, s.IDsToNodes)
}

func TestSyncer_DBIntegration(t *testing.T) {
	ctx := testutils.Context(t)
	reg, regAddress, owner, sim := startNewChainWithRegistry(t)

	_, err := reg.AddCapabilities(owner, []kcr_v1.CapabilitiesRegistryCapability{writeChainCapability})
	require.NoError(t, err, "AddCapability failed for %s", writeChainCapability.LabelledName)
	sim.Commit()

	cid, err := reg.GetHashedCapabilityId(&bind.CallOpts{}, writeChainCapability.LabelledName, writeChainCapability.Version)
	require.NoError(t, err)

	_, err = reg.AddNodeOperators(owner, []kcr_v1.CapabilitiesRegistryNodeOperator{
		{
			Admin: owner.From,
			Name:  "TEST_NOP",
		},
	})
	require.NoError(t, err)
	sim.Commit()

	nodeSet := [][32]byte{
		randomWord(),
		randomWord(),
		randomWord(),
	}

	signersSet := [][32]byte{
		randomWord(),
		randomWord(),
		randomWord(),
	}

	nodes := []kcr_v1.CapabilitiesRegistryNodeParams{
		{
			// The first NodeOperatorId has id 1 since the id is auto-incrementing.
			NodeOperatorId:      uint32(1),
			Signer:              signersSet[0],
			P2pId:               nodeSet[0],
			EncryptionPublicKey: randomWord(),
			HashedCapabilityIds: [][32]byte{cid},
		},
		{
			// The first NodeOperatorId has id 1 since the id is auto-incrementing.
			NodeOperatorId:      uint32(1),
			Signer:              signersSet[1],
			P2pId:               nodeSet[1],
			EncryptionPublicKey: randomWord(),
			HashedCapabilityIds: [][32]byte{cid},
		},
		{
			// The first NodeOperatorId has id 1 since the id is auto-incrementing.
			NodeOperatorId:      uint32(1),
			Signer:              signersSet[2],
			P2pId:               nodeSet[2],
			EncryptionPublicKey: randomWord(),
			HashedCapabilityIds: [][32]byte{cid},
		},
	}
	_, err = reg.AddNodes(owner, nodes)
	require.NoError(t, err)
	sim.Commit()

	config := &capabilitiespb.CapabilityConfig{
		DefaultConfig: values.Proto(values.EmptyMap()).GetMapValue(),
		RemoteConfig: &capabilitiespb.CapabilityConfig_RemoteTriggerConfig{
			RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
				RegistrationRefresh: durationpb.New(20 * time.Second),
				RegistrationExpiry:  durationpb.New(60 * time.Second),
				// F + 1
				MinResponsesToAggregate: uint32(1) + 1,
			},
		},
	}
	configb, err := proto.Marshal(config)
	require.NoError(t, err)

	cfgs := []kcr_v1.CapabilitiesRegistryCapabilityConfiguration{
		{
			CapabilityId: cid,
			Config:       configb,
		},
	}
	_, err = reg.AddDON(
		owner,
		nodeSet,
		cfgs,
		true,
		true,
		1,
	)
	require.NoError(t, err)
	sim.Commit()

	factory := newContractReaderFactory(t, sim)
	syncerORM := newORM(t)
	syncerORM.ormMock.On("LatestLocalRegistry", mock.Anything).Return(nil, errors.New("no state found"))
	syncerORM.ormMock.On("AddLocalRegistry", mock.Anything, mock.Anything).Return(nil)
	syncer, err := newTestSyncer(logger.TestLogger(t), func() (p2ptypes.PeerID, error) { return p2ptypes.PeerID{}, nil }, factory, regAddress.Hex(), syncerORM)
	require.NoError(t, err)
	require.NoError(t, syncer.Start(ctx))
	t.Cleanup(func() {
		syncerORM.Cleanup()
		require.NoError(t, syncer.Close())
	})

	l := &launcher{}
	syncer.AddListener(l)

	var latestLocalRegistryCalled, addLocalRegistryCalled bool
	timeout := time.After(testutils.WaitTimeout(t))

	for !latestLocalRegistryCalled || !addLocalRegistryCalled {
		select {
		case val := <-syncerORM.latestLocalRegistryCh:
			assert.Equal(t, struct{}{}, val)
			latestLocalRegistryCalled = true
		case val := <-syncerORM.addLocalRegistryCh:
			assert.Equal(t, struct{}{}, val)
			addLocalRegistryCalled = true
		case <-timeout:
			t.Fatal("test timed out; channels did not received data")
		}
	}
}

func TestSyncer_LocalNode(t *testing.T) {
	ctx := t.Context()
	lggr := logger.TestLogger(t)

	var pid p2ptypes.PeerID
	err := pid.UnmarshalText([]byte("12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N"))
	require.NoError(t, err)

	workflowDonNodes := []p2ptypes.PeerID{
		pid,
		randomWord(),
		randomWord(),
		randomWord(),
	}

	dID := uint32(1)
	// The below state describes a Workflow DON (AcceptsWorkflows = true),
	// which exposes the streams-trigger and write_chain capabilities.
	// We expect receivers to be wired up and both capabilities to be added to the registry.
	localRegistry := registrysyncer.NewLocalRegistry(
		lggr,
		func() (p2ptypes.PeerID, error) { return pid, nil },
		map[registrysyncer.DonID]registrysyncer.DON{
			registrysyncer.DonID(dID): {
				DON: capabilities.DON{
					ID:               dID,
					ConfigVersion:    uint32(2),
					F:                uint8(1),
					IsPublic:         true,
					AcceptsWorkflows: true,
					Members:          workflowDonNodes,
				},
			},
		},
		map[p2ptypes.PeerID]registrysyncer.NodeInfo{
			workflowDonNodes[0]: {
				NodeOperatorID:      1,
				Signer:              randomWord(),
				P2PID:               workflowDonNodes[0],
				EncryptionPublicKey: randomWord(),
				Version:             "v1",
			},
			workflowDonNodes[1]: {
				NodeOperatorID:      1,
				Signer:              randomWord(),
				P2PID:               workflowDonNodes[1],
				EncryptionPublicKey: randomWord(),
				Version:             "v1",
			},
			workflowDonNodes[2]: {
				NodeOperatorID:      1,
				Signer:              randomWord(),
				P2PID:               workflowDonNodes[2],
				EncryptionPublicKey: randomWord(),
				Version:             "v1",
			},
			workflowDonNodes[3]: {
				NodeOperatorID:      1,
				Signer:              randomWord(),
				P2PID:               workflowDonNodes[3],
				EncryptionPublicKey: randomWord(),
				Version:             "v1",
			},
		},
		map[string]registrysyncer.Capability{
			"test-target@1.0.0": {
				CapabilityType: capabilities.CapabilityTypeTarget,
				ID:             "write-chain@1.0.1",
			},
		},
	)

	node, err := localRegistry.LocalNode(ctx)
	require.NoError(t, err)

	don := capabilities.DON{
		ID:               dID,
		ConfigVersion:    2,
		Members:          workflowDonNodes,
		F:                1,
		IsPublic:         true,
		AcceptsWorkflows: true,
	}
	expectedNode := capabilities.Node{
		PeerID:              &pid,
		NodeOperatorID:      1,
		Signer:              localRegistry.IDsToNodes[pid].Signer,
		EncryptionPublicKey: localRegistry.IDsToNodes[pid].EncryptionPublicKey,
		WorkflowDON:         don,
		CapabilityDONs:      []capabilities.DON{don},
	}
	assert.Equal(t, expectedNode, node)
}

func TestReader_V2_Integration(t *testing.T) {
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	
	// Create a simulated backend similar to V1 tests
	owner := evmtestutils.MustNewSimTransactor(t)
	i := &big.Int{}
	oneEth, _ := i.SetString("100000000000000000000", 10)
	gasLimit := ethconfig.Defaults.Miner.GasCeil * 2 // 60 M blocks

	simulatedBackend := simulated.NewBackend(gethtypes.GenesisAlloc{owner.From: {
		Balance: oneEth,
	}}, simulated.WithBlockGasLimit(gasLimit))
	simulatedBackend.Commit()

	// Deploy a V2 capabilities registry
	regAddress, _, reg, err := kcr_v2.DeployCapabilitiesRegistry(owner, simulatedBackend.Client(), kcr_v2.CapabilitiesRegistryConstructorParams{})
	require.NoError(t, err, "DeployCapabilitiesRegistry failed")
	simulatedBackend.Commit()

	// Add a V2 capability with string ID and metadata
	writeChainCapabilityV2 := kcr_v2.CapabilitiesRegistryCapability{
		CapabilityId:          "write-chain@1.0.1",
		ConfigurationContract: common.Address{},
		Metadata:              []byte(`{"capabilityType": 3, "responseType": 1}`), // 3 = target capability
	}

	// Add capability
	_, err = reg.AddCapabilities(owner, []kcr_v2.CapabilitiesRegistryCapability{writeChainCapabilityV2})
	require.NoError(t, err, "AddCapability failed for %s", writeChainCapabilityV2.CapabilityId)
	simulatedBackend.Commit()

	// V2 uses string capability IDs directly
	cid := writeChainCapabilityV2.CapabilityId

	// Add node operator
	_, err = reg.AddNodeOperators(owner, []kcr_v2.CapabilitiesRegistryNodeOperator{
		{
			Admin: owner.From,
			Name:  "TEST_NOP_V2",
		},
	})
	require.NoError(t, err, "Failed to add node operator")
	simulatedBackend.Commit()

	// Create test nodes
	nodeSet := [][32]byte{
		randomWord(),
		randomWord(),
		randomWord(),
	}

	signersSet := [][32]byte{
		randomWord(),
		randomWord(),
		randomWord(),
	}

	encPubKey1 := randomWord()
	encPubKey2 := randomWord()
	encPubKey3 := randomWord()

	// V2 nodes use string capability IDs and require CsaKey
	nodes := []kcr_v2.CapabilitiesRegistryNodeParams{
		{
			NodeOperatorId:      uint32(1),
			Signer:              signersSet[0],
			P2pId:               nodeSet[0],
			EncryptionPublicKey: encPubKey1,
			CsaKey:              randomWord(),
			CapabilityIds:       []string{cid},
		},
		{
			NodeOperatorId:      uint32(1),
			Signer:              signersSet[1],
			P2pId:               nodeSet[1],
			EncryptionPublicKey: encPubKey2,
			CsaKey:              randomWord(),
			CapabilityIds:       []string{cid},
		},
		{
			NodeOperatorId:      uint32(1),
			Signer:              signersSet[2],
			P2pId:               nodeSet[2],
			EncryptionPublicKey: encPubKey3,
			CsaKey:              randomWord(),
			CapabilityIds:       []string{cid},
		},
	}
	_, err = reg.AddNodes(owner, nodes)
	require.NoError(t, err, "Failed to add nodes")
	simulatedBackend.Commit()

	// Create capability configuration
	config := &capabilitiespb.CapabilityConfig{
		DefaultConfig: values.Proto(values.EmptyMap()).GetMapValue(),
		RemoteConfig: &capabilitiespb.CapabilityConfig_RemoteTriggerConfig{
			RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
				RegistrationRefresh: durationpb.New(20 * time.Second),
				RegistrationExpiry:  durationpb.New(60 * time.Second),
				MinResponsesToAggregate: uint32(1) + 1,
				MessageExpiry:           durationpb.New(120 * time.Second),
			},
		},
	}
	configb, err := proto.Marshal(config)
	require.NoError(t, err)

	// V2 DON configuration uses string capability IDs
	cfgs := []kcr_v2.CapabilitiesRegistryCapabilityConfiguration{
		{
			CapabilityId: cid,
			Config:       configb,
		},
	}

	// Add DON using AddDONs with DON family (V2 feature)
	newDONs := []kcr_v2.CapabilitiesRegistryNewDONParams{
		{
			Name:                     "test-don-v2",
			DonFamilies:              []string{"workflow-don-family"},
			Config:                   []byte{},
			CapabilityConfigurations: cfgs,
			Nodes:                    nodeSet,
			F:                        1,
			IsPublic:                 true,
			AcceptsWorkflows:         true,
		},
	}
	_, err = reg.AddDONs(owner, newDONs)
	require.NoError(t, err)
	simulatedBackend.Commit()

	// Test the syncer using the same setup as V1 tests
	db := pgtest.NewSqlxDB(t)
	factory := newContractReaderFactory(t, simulatedBackend)
	syncerORM := registrysyncer.NewORM(db, lggr)
	syncer, err := registrysyncer.New(lggr, func() (p2ptypes.PeerID, error) { return p2ptypes.PeerID{}, nil }, factory, regAddress.Hex(), syncerORM)
	require.NoError(t, err)

	l := &launcher{}
	syncer.AddListener(l)

	err = syncer.Sync(ctx, false)
	require.NoError(t, err)

	s := l.localRegistry
	require.NotNil(t, s)

	// Test V2 capabilities with string IDs
	assert.Len(t, s.IDsToCapabilities, 1)
	gotCap := s.IDsToCapabilities[cid]
	assert.Equal(t, registrysyncer.Capability{
		CapabilityType: capabilities.CapabilityTypeTarget,
		ID:             "write-chain@1.0.1",
	}, gotCap)

	// Test V2 DON with family
	assert.Len(t, s.IDsToDONs, 1)
	expectedDON := capabilities.DON{
		ID:               1,
		ConfigVersion:    1,
		IsPublic:         true,
		AcceptsWorkflows: true,
		F:                1,
		Members:          toPeerIDs(nodeSet),
	}
	gotDon := s.IDsToDONs[1]
	assert.Equal(t, expectedDON, gotDon.DON)
	assert.Equal(t, configb, gotDon.CapabilityConfigurations[cid].Config)

	// Test V2 node info with string capability IDs
	expectedNodesInfo := []registrysyncer.NodeInfo{
		{
			NodeOperatorID:      uint32(1),
			ConfigCount:         1,
			WorkflowDONId:       1,
			Signer:              signersSet[0],
			P2PID:               p2ptypes.PeerID(nodeSet[0]),
			EncryptionPublicKey: encPubKey1,
			CapabilityIDs:       &[]string{cid}, // V2 uses string IDs
			CapabilitiesDONIds:  []*big.Int{},
			Version:             "v2",
		},
		{
			NodeOperatorID:      uint32(1),
			ConfigCount:         1,
			WorkflowDONId:       1,
			Signer:              signersSet[1],
			P2PID:               p2ptypes.PeerID(nodeSet[1]),
			EncryptionPublicKey: encPubKey2,
			CapabilityIDs:       &[]string{cid}, // V2 uses string IDs
			CapabilitiesDONIds:  []*big.Int{},
			Version:             "v2",
		},
		{
			NodeOperatorID:      uint32(1),
			ConfigCount:         1,
			WorkflowDONId:       1,
			Signer:              signersSet[2],
			P2PID:               p2ptypes.PeerID(nodeSet[2]),
			EncryptionPublicKey: encPubKey3,
			CapabilityIDs:       &[]string{cid}, // V2 uses string IDs
			CapabilitiesDONIds:  []*big.Int{},
			Version:             "v2",
		},
	}

	assert.Len(t, s.IDsToNodes, 3)
	assert.Equal(t, map[p2ptypes.PeerID]registrysyncer.NodeInfo{
		nodeSet[0]: expectedNodesInfo[0],
		nodeSet[1]: expectedNodesInfo[1],
		nodeSet[2]: expectedNodesInfo[2],
	}, s.IDsToNodes)
}

func TestSyncer_V2_DBIntegration(t *testing.T) {
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	
	// Create a simulated backend similar to V1 tests
	owner := evmtestutils.MustNewSimTransactor(t)
	i := &big.Int{}
	oneEth, _ := i.SetString("100000000000000000000", 10)
	gasLimit := ethconfig.Defaults.Miner.GasCeil * 2 // 60 M blocks

	simulatedBackend := simulated.NewBackend(gethtypes.GenesisAlloc{owner.From: {
		Balance: oneEth,
	}}, simulated.WithBlockGasLimit(gasLimit))
	simulatedBackend.Commit()

	// Deploy a V2 capabilities registry
	regAddress, _, reg, err := kcr_v2.DeployCapabilitiesRegistry(owner, simulatedBackend.Client(), kcr_v2.CapabilitiesRegistryConstructorParams{})
	require.NoError(t, err, "DeployCapabilitiesRegistry failed")
	simulatedBackend.Commit()

	// Add a V2 capability
	writeChainCapabilityV2 := kcr_v2.CapabilitiesRegistryCapability{
		CapabilityId:          "write-chain@1.0.1",
		ConfigurationContract: common.Address{},
		Metadata:              []byte(`{"capabilityType": 3, "responseType": 1}`), // 3 = target capability
	}

	_, err = reg.AddCapabilities(owner, []kcr_v2.CapabilitiesRegistryCapability{writeChainCapabilityV2})
	require.NoError(t, err, "AddCapability failed for %s", writeChainCapabilityV2.CapabilityId)
	simulatedBackend.Commit()

	cid := writeChainCapabilityV2.CapabilityId

	// Add node operator
	_, err = reg.AddNodeOperators(owner, []kcr_v2.CapabilitiesRegistryNodeOperator{
		{
			Admin: owner.From,
			Name:  "TEST_NOP_V2",
		},
	})
	require.NoError(t, err)
	simulatedBackend.Commit()

	// Create test nodes
	nodeSet := [][32]byte{
		randomWord(),
		randomWord(),
		randomWord(),
	}

	signersSet := [][32]byte{
		randomWord(),
		randomWord(),
		randomWord(),
	}

	nodes := []kcr_v2.CapabilitiesRegistryNodeParams{
		{
			NodeOperatorId:      uint32(1),
			Signer:              signersSet[0],
			P2pId:               nodeSet[0],
			EncryptionPublicKey: randomWord(),
			CsaKey:              randomWord(),
			CapabilityIds:       []string{cid},
		},
		{
			NodeOperatorId:      uint32(1),
			Signer:              signersSet[1],
			P2pId:               nodeSet[1],
			EncryptionPublicKey: randomWord(),
			CsaKey:              randomWord(),
			CapabilityIds:       []string{cid},
		},
		{
			NodeOperatorId:      uint32(1),
			Signer:              signersSet[2],
			P2pId:               nodeSet[2],
			EncryptionPublicKey: randomWord(),
			CsaKey:              randomWord(),
			CapabilityIds:       []string{cid},
		},
	}
	_, err = reg.AddNodes(owner, nodes)
	require.NoError(t, err)
	simulatedBackend.Commit()

	// Create capability configuration
	config := &capabilitiespb.CapabilityConfig{
		DefaultConfig: values.Proto(values.EmptyMap()).GetMapValue(),
		RemoteConfig: &capabilitiespb.CapabilityConfig_RemoteTriggerConfig{
			RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
				RegistrationRefresh: durationpb.New(20 * time.Second),
				RegistrationExpiry:  durationpb.New(60 * time.Second),
				MinResponsesToAggregate: uint32(1) + 1,
			},
		},
	}
	configb, err := proto.Marshal(config)
	require.NoError(t, err)

	cfgs := []kcr_v2.CapabilitiesRegistryCapabilityConfiguration{
		{
			CapabilityId: cid,
			Config:       configb,
		},
	}

	// Add DON using AddDONs with DON family (V2 feature)
	newDONs := []kcr_v2.CapabilitiesRegistryNewDONParams{
		{
			Name:                     "test-don-v2-db",
			DonFamilies:              []string{"workflow-don-family-v2"},
			Config:                   []byte{},
			CapabilityConfigurations: cfgs,
			Nodes:                    nodeSet,
			F:                        1,
			IsPublic:                 true,
			AcceptsWorkflows:         true,
		},
	}
	_, err = reg.AddDONs(owner, newDONs)
	require.NoError(t, err)
	simulatedBackend.Commit()

	// Test database integration
	syncerORM := newORM(t)
	syncerORM.ormMock.On("LatestLocalRegistry", mock.Anything).Return(nil, errors.New("no state found"))
	syncerORM.ormMock.On("AddLocalRegistry", mock.Anything, mock.Anything).Return(nil)

	factory := newContractReaderFactory(t, simulatedBackend)

	syncer, err := newTestSyncer(
		lggr,
		func() (p2ptypes.PeerID, error) { return p2ptypes.PeerID{}, nil },
		factory,
		regAddress.Hex(),
		syncerORM,
	)
	require.NoError(t, err)
	require.NoError(t, syncer.Start(ctx))
	
	t.Cleanup(func() {
		syncerORM.Cleanup()
		require.NoError(t, syncer.Close())
	})

	l := &launcher{}
	syncer.AddListener(l)

	// Test that the syncer calls the ORM methods
	var latestLocalRegistryCalled, addLocalRegistryCalled bool
	timeout := time.After(testutils.WaitTimeout(t))

	for !latestLocalRegistryCalled || !addLocalRegistryCalled {
		select {
		case val := <-syncerORM.latestLocalRegistryCh:
			assert.Equal(t, struct{}{}, val)
			latestLocalRegistryCalled = true
		case val := <-syncerORM.addLocalRegistryCh:
			assert.Equal(t, struct{}{}, val)
			addLocalRegistryCalled = true
		case <-timeout:
			t.Fatal("test timed out; channels did not received data")
		}
	}
}

func TestSyncer_V2_LocalNode(t *testing.T) {
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)

	var pid p2ptypes.PeerID
	err := pid.UnmarshalText([]byte("12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N"))
	require.NoError(t, err)

	workflowDonNodes := []p2ptypes.PeerID{
		pid,
		randomWord(),
		randomWord(),
		randomWord(),
	}

	dID := uint32(1)
	// Test V2 local registry with string capability IDs
	localRegistry := registrysyncer.NewLocalRegistry(
		lggr,
		func() (p2ptypes.PeerID, error) { return pid, nil },
		map[registrysyncer.DonID]registrysyncer.DON{
			registrysyncer.DonID(dID): {
				DON: capabilities.DON{
					ID:               dID,
					ConfigVersion:    uint32(2),
					F:                uint8(1),
					IsPublic:         true,
					AcceptsWorkflows: true,
					Members:          workflowDonNodes,
				},
			},
		},
		map[p2ptypes.PeerID]registrysyncer.NodeInfo{
			workflowDonNodes[0]: {
				NodeOperatorID:      1,
				Signer:              randomWord(),
				P2PID:               workflowDonNodes[0],
				EncryptionPublicKey: randomWord(),
				CapabilityIDs:       &[]string{"write-chain@1.0.1", "trigger@1.0.0"}, // V2 uses string IDs
				Version:             "v2",
			},
			workflowDonNodes[1]: {
				NodeOperatorID:      1,
				Signer:              randomWord(),
				P2PID:               workflowDonNodes[1],
				EncryptionPublicKey: randomWord(),
				CapabilityIDs:       &[]string{"write-chain@1.0.1", "trigger@1.0.0"}, // V2 uses string IDs
				Version:             "v2",
			},
			workflowDonNodes[2]: {
				NodeOperatorID:      1,
				Signer:              randomWord(),
				P2PID:               workflowDonNodes[2],
				EncryptionPublicKey: randomWord(),
				CapabilityIDs:       &[]string{"write-chain@1.0.1"}, // V2 uses string IDs
				Version:             "v2",
			},
			workflowDonNodes[3]: {
				NodeOperatorID:      1,
				Signer:              randomWord(),
				P2PID:               workflowDonNodes[3],
				EncryptionPublicKey: randomWord(),
				CapabilityIDs:       &[]string{"write-chain@1.0.1"}, // V2 uses string IDs
				Version:             "v2",
			},
		},
		map[string]registrysyncer.Capability{
			"write-chain@1.0.1": {
				CapabilityType: capabilities.CapabilityTypeTarget,
				ID:             "write-chain@1.0.1",
			},
			"trigger@1.0.0": {
				CapabilityType: capabilities.CapabilityTypeTrigger,
				ID:             "trigger@1.0.0",
			},
		},
	)

	node, err := localRegistry.LocalNode(ctx)
	require.NoError(t, err)

	don := capabilities.DON{
		ID:               dID,
		ConfigVersion:    2,
		Members:          workflowDonNodes,
		F:                1,
		IsPublic:         true,
		AcceptsWorkflows: true,
	}
	expectedNode := capabilities.Node{
		PeerID:              &pid,
		NodeOperatorID:      1,
		Signer:              localRegistry.IDsToNodes[pid].Signer,
		EncryptionPublicKey: localRegistry.IDsToNodes[pid].EncryptionPublicKey,
		WorkflowDON:         don,
		CapabilityDONs:      []capabilities.DON{don},
	}
	assert.Equal(t, expectedNode, node)

	// Test that V2 capabilities are properly handled
	assert.Len(t, localRegistry.IDsToCapabilities, 2)
	assert.Contains(t, localRegistry.IDsToCapabilities, "write-chain@1.0.1")
	assert.Contains(t, localRegistry.IDsToCapabilities, "trigger@1.0.0")

	// Test that V2 node info has string capability IDs
	nodeInfo := localRegistry.IDsToNodes[pid]
	assert.Equal(t, "v2", nodeInfo.Version)
	assert.NotNil(t, nodeInfo.CapabilityIDs)
	assert.Nil(t, nodeInfo.HashedCapabilityIDs)
	assert.Equal(t, []string{"write-chain@1.0.1", "trigger@1.0.0"}, *nodeInfo.CapabilityIDs)
}

// Add this helper struct to implement the ContractReaderFactory interface
type testContractReaderFactory struct {
	backendTH *captestutils.EVMBackendTH
	t         *testing.T
}

func (f *testContractReaderFactory) NewContractReader(ctx context.Context, bytes []byte) (types.ContractReader, error) {
	return f.backendTH.NewContractReader(ctx, f.t, bytes)
}

func newTestSyncer(
	lggr logger.Logger,
	getPeerID func() (p2ptypes.PeerID, error),
	relayer registrysyncer.ContractReaderFactory,
	registryAddress string,
	orm *orm,
) (registrysyncer.RegistrySyncer, error) {
	rs, err := registrysyncer.New(lggr, getPeerID, relayer, registryAddress, orm)
	if err != nil {
		return nil, err
	}
	return rs, nil
}
