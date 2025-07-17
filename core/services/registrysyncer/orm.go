package registrysyncer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

type capabilitiesRegistryNodeInfo struct {
	NodeOperatorID      uint32       `json:"nodeOperatorId"`
	ConfigCount         uint32       `json:"configCount"`
	WorkflowDONId       uint32       `json:"workflowDONId"`
	Signer              types.PeerID `json:"signer"`
	P2PID               types.PeerID `json:"p2pId"`
	EncryptionPublicKey [32]byte     `json:"encryptionPublicKey"`
	CapabilitiesDONIds  []string     `json:"capabilitiesDONIds"`

	// Version-specific fields
	HashedCapabilityIDs *[]types.PeerID `json:"hashedCapabilityIds,omitempty"`
	CapabilityIDs       *[]string       `json:"capabilityIds,omitempty"`
	Version             string          `json:"version"`
}

func (l *LocalRegistry) MarshalJSON() ([]byte, error) {
	idsToNodes := make(map[types.PeerID]capabilitiesRegistryNodeInfo)
	for k, v := range l.IDsToNodes {
		var hashedCapabilityIDs *[]types.PeerID
		if v.HashedCapabilityIDs != nil {
			hashedIDs := make([]types.PeerID, len(*v.HashedCapabilityIDs))
			for i, id := range *v.HashedCapabilityIDs {
				hashedIDs[i] = types.PeerID(id[:])
			}
			hashedCapabilityIDs = &hashedIDs
		}

		var capabilityIDs *[]string
		if v.CapabilityIDs != nil {
			capabilityIDs = v.CapabilityIDs
		}

		capabilitiesDONIds := make([]string, len(v.CapabilitiesDONIds))
		for i, id := range v.CapabilitiesDONIds {
			capabilitiesDONIds[i] = id.String()
		}

		idsToNodes[k] = capabilitiesRegistryNodeInfo{
			NodeOperatorID:      v.NodeOperatorID,
			ConfigCount:         v.ConfigCount,
			WorkflowDONId:       v.WorkflowDONId,
			Signer:              types.PeerID(v.Signer[:]),
			P2PID:               v.P2PID,
			EncryptionPublicKey: v.EncryptionPublicKey,
			HashedCapabilityIDs: hashedCapabilityIDs,
			CapabilityIDs:       capabilityIDs,
			CapabilitiesDONIds:  capabilitiesDONIds,
			Version:             v.Version,
		}
	}

	b, err := json.Marshal(&struct {
		IDsToDONs         map[DonID]DON
		IDsToNodes        map[types.PeerID]capabilitiesRegistryNodeInfo
		IDsToCapabilities map[string]Capability
	}{
		IDsToDONs:         l.IDsToDONs,
		IDsToNodes:        idsToNodes,
		IDsToCapabilities: l.IDsToCapabilities,
	})
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

func (l *LocalRegistry) UnmarshalJSON(data []byte) error {
	temp := struct {
		IDsToDONs         map[DonID]DON
		IDsToNodes        map[types.PeerID]capabilitiesRegistryNodeInfo
		IDsToCapabilities map[string]Capability
	}{
		IDsToDONs:         make(map[DonID]DON),
		IDsToNodes:        make(map[types.PeerID]capabilitiesRegistryNodeInfo),
		IDsToCapabilities: make(map[string]Capability),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("failed to unmarshal state: %w", err)
	}

	l.IDsToDONs = temp.IDsToDONs

	l.IDsToNodes = make(map[types.PeerID]NodeInfo)
	for peerID, v := range temp.IDsToNodes {
		var hashedCapabilityIDs *[][32]byte
		if v.HashedCapabilityIDs != nil {
			hashedIDs := make([][32]byte, len(*v.HashedCapabilityIDs))
			for i, id := range *v.HashedCapabilityIDs {
				hashedIDs[i] = [32]byte(id)
			}
			hashedCapabilityIDs = &hashedIDs
		}

		var capabilityIDs *[]string
		if v.CapabilityIDs != nil {
			capabilityIDs = v.CapabilityIDs
		}

		capabilitiesDONIds := make([]*big.Int, len(v.CapabilitiesDONIds))
		for i, id := range v.CapabilitiesDONIds {
			bigInt := new(big.Int)
			bigInt.SetString(id, 10)
			capabilitiesDONIds[i] = bigInt
		}

		nodeInfo := NodeInfo{
			NodeOperatorID:      v.NodeOperatorID,
			ConfigCount:         v.ConfigCount,
			WorkflowDONId:       v.WorkflowDONId,
			Signer:              [32]byte(v.Signer),
			P2PID:               v.P2PID,
			EncryptionPublicKey: v.EncryptionPublicKey,
			HashedCapabilityIDs: hashedCapabilityIDs,
			CapabilityIDs:       capabilityIDs,
			CapabilitiesDONIds:  capabilitiesDONIds,
			Version:             v.Version,
		}

		l.IDsToNodes[peerID] = nodeInfo
	}

	l.IDsToCapabilities = temp.IDsToCapabilities

	return nil
}

type ORM interface {
	AddLocalRegistry(ctx context.Context, localRegistry LocalRegistry) error
	LatestLocalRegistry(ctx context.Context) (*LocalRegistry, error)
}

type orm struct {
	ds   sqlutil.DataSource
	lggr logger.Logger
}

var _ ORM = (*orm)(nil)

func NewORM(ds sqlutil.DataSource, lggr logger.Logger) orm {
	namedLogger := logger.Named(lggr, "RegistrySyncerORM")
	return orm{
		ds:   ds,
		lggr: namedLogger,
	}
}

func (orm orm) AddLocalRegistry(ctx context.Context, localRegistry LocalRegistry) error {
	return sqlutil.TransactDataSource(ctx, orm.ds, nil, func(tx sqlutil.DataSource) error {
		localRegistryJSON, err := localRegistry.MarshalJSON()
		if err != nil {
			return err
		}
		hash := sha256.Sum256(localRegistryJSON)
		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO registry_syncer_states (data, data_hash) VALUES ($1, $2) ON CONFLICT (data_hash) DO NOTHING`,
			localRegistryJSON, hex.EncodeToString(hash[:]),
		)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `DELETE FROM registry_syncer_states
WHERE data_hash NOT IN (
    SELECT data_hash FROM registry_syncer_states
    ORDER BY id DESC
    LIMIT 10
);`)
		return err
	})
}

func (orm orm) LatestLocalRegistry(ctx context.Context) (*LocalRegistry, error) {
	var localRegistry LocalRegistry
	var localRegistryJSON string
	err := orm.ds.GetContext(ctx, &localRegistryJSON, `SELECT data FROM registry_syncer_states ORDER BY id DESC LIMIT 1`)
	if err != nil {
		return nil, err
	}
	err = localRegistry.UnmarshalJSON([]byte(localRegistryJSON))
	if err != nil {
		return nil, err
	}
	return &localRegistry, nil
}
