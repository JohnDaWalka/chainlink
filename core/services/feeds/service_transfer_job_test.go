package feeds_test

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	proto "github.com/smartcontractkit/chainlink-protos/orchestrator/feedsmanager"

	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/feeds"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/utils/crypto"
)

func Test_Service_TransferJob(t *testing.T) {
	t.Parallel()

	var (
		sourceManagerID = int64(1)
		targetManagerID = int64(2)
		remoteUUID      = uuid.New()
		proposalID      = int64(123)

		sourcePubKey = "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
		targetPubKey = "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"

		sourcePubKeyBytes, _ = hex.DecodeString(sourcePubKey)
		targetPubKeyBytes, _ = hex.DecodeString(targetPubKey)

		csaKey, _ = csakey.NewV2()

		proposal = &feeds.JobProposal{
			ID:             proposalID,
			FeedsManagerID: sourceManagerID,
			RemoteUUID:     remoteUUID,
			Status:         feeds.JobProposalStatusPending,
		}

		sourceManager = &feeds.FeedsManager{
			ID:        sourceManagerID,
			Name:      "Source FMS",
			PublicKey: crypto.PublicKey(sourcePubKeyBytes),
		}

		targetManager = &feeds.FeedsManager{
			ID:        targetManagerID,
			Name:      "Target FMS",
			PublicKey: crypto.PublicKey(targetPubKeyBytes),
		}

		specs = []feeds.JobProposalSpec{
			{
				ID:            456,
				JobProposalID: proposalID,
				Definition:    "test spec definition",
				Version:       1,
				Status:        feeds.SpecStatusPending,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		}

		transferCompleteResponse = &proto.TransferedJobResponse{}

		args = &feeds.TransferJobArgs{
			RemoteUUID:          remoteUUID,
			TargetManagerID:     targetManagerID,
			SourceManagerPubKey: sourcePubKey,
		}
	)

	testCases := []struct {
		name    string
		before  func(svc *TestService)
		args    *feeds.TransferJobArgs
		wantErr string
	}{
		{
			name: "success",
			before: func(svc *TestService) {
				svc.orm.On("GetManager", mock.Anything, targetManagerID).Return(targetManager, nil)
				svc.orm.On("GetManagerByPublicKey", mock.Anything, crypto.PublicKey(sourcePubKeyBytes)).Return(sourceManager, nil)
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, remoteUUID).Return(proposal, nil)
				svc.connMgr.On("GetClient", targetManagerID).Return(svc.fmsClient, nil)
				svc.csaKeystore.On("GetAll").Return([]csakey.KeyV2{csaKey}, nil)
				svc.orm.On("Transact", mock.Anything, mock.AnythingOfType("func(feeds.ORM) error")).
					Return(nil).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(feeds.ORM) error)

						svc.orm.On("ListSpecsByJobProposalIDs", mock.Anything, []int64{proposalID}).Return(specs, nil)
						svc.orm.On("TransferJobProposal", mock.Anything, proposalID, targetManagerID).Return(nil)
						svc.fmsClient.On("TransferedJob", mock.Anything, mock.AnythingOfType("*feedsmanager.TransferedJobRequest")).
							Return(transferCompleteResponse, nil)

						fn(svc.orm)
					})
			},
			args: args,
		},
		{
			name: "CSA key not found",
			before: func(svc *TestService) {
				svc.orm.On("GetManager", mock.Anything, targetManagerID).Return(targetManager, nil)
				svc.orm.On("GetManagerByPublicKey", mock.Anything, crypto.PublicKey(sourcePubKeyBytes)).Return(sourceManager, nil)
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, remoteUUID).Return(proposal, nil)
				svc.connMgr.On("GetClient", targetManagerID).Return(svc.fmsClient, nil)
				svc.csaKeystore.On("GetAll").Return([]csakey.KeyV2{}, nil)
				svc.orm.On("Transact", mock.Anything, mock.AnythingOfType("func(feeds.ORM) error")).
					Return(errors.New("failed to build transfer request: no CSA key found for node")).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(feeds.ORM) error)

						svc.orm.On("ListSpecsByJobProposalIDs", mock.Anything, []int64{proposalID}).Return(specs, nil)

						fn(svc.orm)
					})
			},
			args:    args,
			wantErr: "failed to build transfer request",
		},
		{
			name: "success with correct CSA key in transfer request",
			before: func(svc *TestService) {
				svc.orm.On("GetManager", mock.Anything, targetManagerID).Return(targetManager, nil)
				svc.orm.On("GetManagerByPublicKey", mock.Anything, crypto.PublicKey(sourcePubKeyBytes)).Return(sourceManager, nil)
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, remoteUUID).Return(proposal, nil)
				svc.connMgr.On("GetClient", targetManagerID).Return(svc.fmsClient, nil)
				svc.csaKeystore.On("GetAll").Return([]csakey.KeyV2{csaKey}, nil)
				svc.orm.On("Transact", mock.Anything, mock.AnythingOfType("func(feeds.ORM) error")).
					Return(nil).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(feeds.ORM) error)

						svc.orm.On("ListSpecsByJobProposalIDs", mock.Anything, []int64{proposalID}).Return(specs, nil)
						svc.orm.On("TransferJobProposal", mock.Anything, proposalID, targetManagerID).Return(nil)
						svc.fmsClient.On("TransferedJob", mock.Anything, mock.MatchedBy(func(req *proto.TransferedJobRequest) bool {
							return true
						})).Return(transferCompleteResponse, nil)

						fn(svc.orm)
					})
			},
			args: args,
		},
		{
			name: "job proposal not found",
			before: func(svc *TestService) {
				svc.orm.On("GetManager", mock.Anything, targetManagerID).Return(targetManager, nil)
				svc.orm.On("GetManagerByPublicKey", mock.Anything, crypto.PublicKey(sourcePubKeyBytes)).Return(sourceManager, nil)
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, remoteUUID).Return(nil, sql.ErrNoRows)
			},
			args:    args,
			wantErr: "failed to get job proposal",
		},
		{
			name: "source manager mismatch",
			before: func(svc *TestService) {
				wrongSourceProposal := &feeds.JobProposal{
					ID:             proposalID,
					FeedsManagerID: 999, // Different from sourceManagerID
					RemoteUUID:     remoteUUID,
					Status:         feeds.JobProposalStatusPending,
				}

				svc.orm.On("GetManager", mock.Anything, targetManagerID).Return(targetManager, nil)
				svc.orm.On("GetManagerByPublicKey", mock.Anything, crypto.PublicKey(sourcePubKeyBytes)).Return(sourceManager, nil)
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, remoteUUID).Return(wrongSourceProposal, nil)
			},
			args:    args,
			wantErr: "job proposal does not belong to the specified source manager",
		},
		{
			name: "source and target managers are the same",
			before: func(svc *TestService) {
				svc.orm.On("GetManager", mock.Anything, sourceManagerID).Return(sourceManager, nil)
				svc.orm.On("GetManagerByPublicKey", mock.Anything, crypto.PublicKey(sourcePubKeyBytes)).Return(sourceManager, nil) // Same manager for both calls
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, remoteUUID).Return(proposal, nil)
			},
			args: &feeds.TransferJobArgs{
				RemoteUUID:          remoteUUID,
				TargetManagerID:     sourceManagerID,
				SourceManagerPubKey: sourcePubKey, // Same key to test error
			},
			wantErr: "source and target managers cannot be the same",
		},
		{
			name: "target manager not found",
			before: func(svc *TestService) {
				svc.orm.On("GetManager", mock.Anything, targetManagerID).Return(targetManager, nil)
				svc.orm.On("GetManagerByPublicKey", mock.Anything, crypto.PublicKey(sourcePubKeyBytes)).Return(nil, sql.ErrNoRows)
			},
			args:    args,
			wantErr: "failed to get source manager by public key",
		},
		{
			name: "cannot get FMS client",
			before: func(svc *TestService) {
				svc.orm.On("GetManager", mock.Anything, targetManagerID).Return(targetManager, nil)
				svc.orm.On("GetManagerByPublicKey", mock.Anything, crypto.PublicKey(sourcePubKeyBytes)).Return(sourceManager, nil)
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, remoteUUID).Return(proposal, nil)
				svc.connMgr.On("GetClient", targetManagerID).Return(nil, assert.AnError)
			},
			args:    args,
			wantErr: "failed to get target feeds manager client",
		},
		{
			name: "FMS transfer call fails",
			before: func(svc *TestService) {
				svc.orm.On("GetManager", mock.Anything, targetManagerID).Return(targetManager, nil)
				svc.orm.On("GetManagerByPublicKey", mock.Anything, crypto.PublicKey(sourcePubKeyBytes)).Return(sourceManager, nil)
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, remoteUUID).Return(proposal, nil)
				svc.connMgr.On("GetClient", targetManagerID).Return(svc.fmsClient, nil)
				svc.csaKeystore.On("GetAll").Return([]csakey.KeyV2{csaKey}, nil)
				svc.orm.On("Transact", mock.Anything, mock.AnythingOfType("func(feeds.ORM) error")).
					Return(assert.AnError).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(feeds.ORM) error)

						svc.orm.On("ListSpecsByJobProposalIDs", mock.Anything, []int64{proposalID}).Return(specs, nil)
						svc.orm.On("TransferJobProposal", mock.Anything, proposalID, targetManagerID).Return(nil)
						svc.fmsClient.On("TransferedJob", mock.Anything, mock.AnythingOfType("*feedsmanager.TransferedJobRequest")).
							Return(nil, assert.AnError)

						fn(svc.orm)
					})
			},
			args:    args,
			wantErr: assert.AnError.Error(),
		},
		{
			name: "database transfer fails",
			before: func(svc *TestService) {
				svc.orm.On("GetManager", mock.Anything, targetManagerID).Return(targetManager, nil)
				svc.orm.On("GetManagerByPublicKey", mock.Anything, crypto.PublicKey(sourcePubKeyBytes)).Return(sourceManager, nil)
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, remoteUUID).Return(proposal, nil)
				svc.connMgr.On("GetClient", targetManagerID).Return(svc.fmsClient, nil)
				svc.csaKeystore.On("GetAll").Return([]csakey.KeyV2{csaKey}, nil)
				svc.orm.On("Transact", mock.Anything, mock.AnythingOfType("func(feeds.ORM) error")).
					Return(assert.AnError).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(feeds.ORM) error)

						svc.orm.On("ListSpecsByJobProposalIDs", mock.Anything, []int64{proposalID}).Return(specs, nil)
						svc.orm.On("TransferJobProposal", mock.Anything, proposalID, targetManagerID).Return(assert.AnError)

						fn(svc.orm)
					})
			},
			args:    args,
			wantErr: assert.AnError.Error(),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := setupTestServiceCfg(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.JobDistributor.DisplayName = testutils.Ptr("Test-Node")
				c.JobDistributor.AllowedJobTransfers = []toml.JobTransferRule{
					{
						From: sourcePubKey,
						To:   targetPubKey,
					},
				}
			})
			tc.before(svc)

			err := svc.TransferJob(testutils.Context(t), tc.args)

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_Service_TransferJob_AllowedPartnersValidation(t *testing.T) {
	t.Parallel()

	var (
		sourceManagerID = int64(1)
		targetManagerID = int64(2)
		remoteUUID      = uuid.New()
		proposalID      = int64(123)

		sourcePubKey = "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
		targetPubKey = "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"

		sourcePubKeyBytes, _ = hex.DecodeString(sourcePubKey)
		targetPubKeyBytes, _ = hex.DecodeString(targetPubKey)

		proposal = &feeds.JobProposal{
			ID:             proposalID,
			FeedsManagerID: sourceManagerID,
			RemoteUUID:     remoteUUID,
			Status:         feeds.JobProposalStatusPending,
		}

		sourceManager = &feeds.FeedsManager{
			ID:        sourceManagerID,
			Name:      "Source-Manager",
			PublicKey: crypto.PublicKey(sourcePubKeyBytes),
		}

		targetManager = &feeds.FeedsManager{
			ID:        targetManagerID,
			Name:      "Target-Manager",
			PublicKey: crypto.PublicKey(targetPubKeyBytes),
		}

		args = &feeds.TransferJobArgs{
			RemoteUUID:          remoteUUID,
			TargetManagerID:     targetManagerID,
			SourceManagerPubKey: sourcePubKey,
		}
	)

	testCases := []struct {
		name                string
		allowedJobTransfers []toml.JobTransferRule
		wantErr             string
	}{
		{
			name: "transfer allowed - matching from/to rule",
			allowedJobTransfers: []toml.JobTransferRule{
				{
					From: sourcePubKey,
					To:   targetPubKey,
				},
				{
					From: targetPubKey,
					To:   "other-manager-key",
				},
			},
			wantErr: "",
		},
		{
			name: "transfer denied - no matching rule",
			allowedJobTransfers: []toml.JobTransferRule{
				{
					From: "other-source-key",
					To:   targetPubKey,
				},
				{
					From: sourcePubKey,
					To:   "other-target-key",
				},
			},
			wantErr: "job transfer not allowed: no transfer rule found from source manager '" + sourcePubKey + "' to target manager '" + targetPubKey + "'",
		},
		{
			name:                "transfer denied - empty allowed transfers list blocks all",
			allowedJobTransfers: []toml.JobTransferRule{},
			wantErr:             "job transfer not allowed: no transfer rule found from source manager '" + sourcePubKey + "' to target manager '" + targetPubKey + "'",
		},
		{
			name:                "transfer denied - nil allowed transfers list blocks all",
			allowedJobTransfers: nil,
			wantErr:             "job transfer not allowed: no transfer rule found from source manager '" + sourcePubKey + "' to target manager '" + targetPubKey + "'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := setupTestServiceCfg(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.JobDistributor.DisplayName = testutils.Ptr("Test-Node")
				c.JobDistributor.AllowedJobTransfers = tc.allowedJobTransfers
			})

			svc.orm.On("GetManager", mock.Anything, targetManagerID).Return(targetManager, nil)
			svc.orm.On("GetManagerByPublicKey", mock.Anything, crypto.PublicKey(sourcePubKeyBytes)).Return(sourceManager, nil)
			svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, remoteUUID).Return(proposal, nil)

			if tc.wantErr == "" {
				svc.connMgr.On("GetClient", targetManagerID).Return(svc.fmsClient, nil)
				csaKey, _ := csakey.NewV2()
				svc.csaKeystore.On("GetAll").Return([]csakey.KeyV2{csaKey}, nil)

				specs := []feeds.JobProposalSpec{
					{
						ID:            456,
						JobProposalID: proposalID,
						Definition:    "test spec definition",
						Version:       1,
						Status:        feeds.SpecStatusPending,
						CreatedAt:     time.Now(),
						UpdatedAt:     time.Now(),
					},
				}

				svc.orm.On("Transact", mock.Anything, mock.AnythingOfType("func(feeds.ORM) error")).
					Return(nil).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(feeds.ORM) error)

						svc.orm.On("ListSpecsByJobProposalIDs", mock.Anything, []int64{proposalID}).Return(specs, nil)
						svc.orm.On("TransferJobProposal", mock.Anything, proposalID, targetManagerID).Return(nil)
						svc.fmsClient.On("TransferedJob", mock.Anything, mock.AnythingOfType("*feedsmanager.TransferedJobRequest")).
							Return(&proto.TransferedJobResponse{}, nil)

						fn(svc.orm)
					})
			}

			err := svc.TransferJob(testutils.Context(t), args)

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
