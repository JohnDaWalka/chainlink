package feeds_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pb "github.com/smartcontractkit/chainlink-protos/orchestrator/feedsmanager"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/feeds"
)

func Test_RPCHandlers_TransferJob(t *testing.T) {
	var (
		ctx          = testutils.Context(t)
		jobID        = uuid.New()
		sourcePubKey = "source-manager-public-key"
	)

	h := setupTestHandlers(t)

	h.svc.
		On("TransferJob", ctx, &feeds.TransferJobArgs{
			RemoteUUID:          jobID,
			SourceManagerPubKey: sourcePubKey,
			TargetManagerID:     h.feedsManagerID,
		}).
		Return(nil)

	_, err := h.TransferJob(ctx, &pb.TransferJobRequest{
		Id:                  jobID.String(),
		SourceManagerPubKey: sourcePubKey,
	})
	require.NoError(t, err)
}

func Test_RPCHandlers_TransferJob_InvalidUUID(t *testing.T) {
	ctx := testutils.Context(t)

	h := setupTestHandlers(t)

	_, err := h.TransferJob(ctx, &pb.TransferJobRequest{
		Id:                  "invalid-uuid",
		SourceManagerPubKey: "source-key",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid UUID")
}
