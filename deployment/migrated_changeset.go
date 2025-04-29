package deployment

import (
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// Aliases for changeset types
type (
	ChangeSet[C any]           = deployment.ChangeSet[C]
	ChangeLogic[C any]         = deployment.ChangeLogic[C]
	PreconditionVerifier[C any] = deployment.PreconditionVerifier[C]
	ChangeSetV2[C any]         = deployment.ChangeSetV2[C]
	ProposedJob                = deployment.ProposedJob
	ChangesetOutput            = deployment.ChangesetOutput
	ViewState                  = deployment.ViewState
	ViewStateV2                = deployment.ViewStateV2
)

// Aliases for functions
var (
	CreateChangeSet       = deployment.CreateChangeSet
	CreateLegacyChangeSet = deployment.CreateLegacyChangeSet
	MergeChangesetOutput  = deployment.MergeChangesetOutput
)

// Aliases for errors
var (
	ErrInvalidConfig      = deployment.ErrInvalidConfig
	ErrInvalidEnvironment = deployment.ErrInvalidEnvironment
)
