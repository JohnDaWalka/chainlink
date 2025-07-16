package solana

import (
	"github.com/gagliardetto/solana-go"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	solanaCs "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// TransferOwnershipCacheRequest wraps the generic request for cache contracts
type TransferOwnershipCacheRequest struct {
	ChainSel                    uint64
	CurrentOwner, ProposedOwner solana.PublicKey
	Version                     string
	Qualifier                   string
	MCMSCfg                     proposalutils.TimelockConfig
}

// TransferOwnershipCache implementation
var _ cldf.ChangeSetV2[*TransferOwnershipCacheRequest] = TransferOwnershipCache{}

type TransferOwnershipCache struct{}

func (cs TransferOwnershipCache) VerifyPreconditions(env cldf.Environment, req *TransferOwnershipCacheRequest) error {
	return solanaCs.GenericVerifyPreconditions(env, req.ChainSel, req.Version, req.Qualifier, "CacheContract")
}

func (cs TransferOwnershipCache) Apply(env cldf.Environment, req *TransferOwnershipCacheRequest) (cldf.ChangesetOutput, error) {
	genericReq := &solanaCs.TransferOwnershipRequest{
		ChainSel:      req.ChainSel,
		CurrentOwner:  req.CurrentOwner,
		ProposedOwner: req.ProposedOwner,
		Version:       req.Version,
		Qualifier:     req.Qualifier,
		MCMSCfg:       req.MCMSCfg,
		ContractConfig: solanaCs.ContractConfig{
			ContractType: "CacheContract",
			StateType:    "CacheState",
			OperationID:  "transfer-ownership-cache",
			Description:  "transfers ownership of cache to mcms",
		},
	}
	return solanaCs.GenericTransferOwnership(env, genericReq)
}