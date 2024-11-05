package txm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/types"
)

type StuckTxDetectorConfig struct {
	BlockTime             time.Duration
	StuckTxBlockThreshold uint16
	DetectionApiUrl       string
}

type stuckTxDetector struct {
	lggr      logger.Logger
	chainType chaintype.ChainType
	config    StuckTxDetectorConfig
}

func NewStuckTxDetector(lggr logger.Logger, chaintype chaintype.ChainType, config StuckTxDetectorConfig) *stuckTxDetector {
	return &stuckTxDetector{
		lggr:      lggr,
		chainType: chaintype,
		config:    config,
	}
}

func (s *stuckTxDetector) DetectStuckTransaction(tx *types.Transaction) (bool, error) {
	switch s.chainType {
	// TODO: add correct chaintype
	case chaintype.ChainArbitrum:
		result, err := s.apiBasedDetection(tx)
		if result || err != nil {
			return result, err
		}
		return s.timeBasedDetection(tx), nil
	default:
		return s.timeBasedDetection(tx), nil
	}
}

func (s *stuckTxDetector) timeBasedDetection(tx *types.Transaction) bool {
	threshold := (s.config.BlockTime * time.Duration(s.config.StuckTxBlockThreshold))
	if time.Since(tx.LastBroadcastAt) > threshold && !tx.LastBroadcastAt.IsZero() {
		s.lggr.Debugf("TxID: %v last broadcast was: %v which is more than the configured threshold: %v. Transaction is now considered stuck and will be purged.",
			tx.ID, tx.LastBroadcastAt, threshold)
		return true
	}
	return false
}

type ApiResponse struct {
	Status string      `json:"status,omitempty"`
	Hash   common.Hash `json:"hash,omitempty"`
}

const (
	ApiStatusPending   = "PENDING"
	ApiStatusIncluded  = "INCLUDED"
	ApiStatusFailed    = "FAILED"
	ApiStatusCancelled = "CANCELLED"
	ApiStatusUnknown   = "UNKNOWN"
)

func (s *stuckTxDetector) apiBasedDetection(tx *types.Transaction) (bool, error) {
	for _, attempt := range tx.Attempts {
		resp, err := http.Get(s.config.DetectionApiUrl + attempt.Hash.String())
		if err != nil {
			return false, fmt.Errorf("failed to get transaction status for txID: %v, attemptHash: %v - %w", tx.ID, attempt.Hash, err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return false, err
		}
		var apiResponse ApiResponse
		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			return false, fmt.Errorf("failed to unmarshal response for txID: %v, attemptHash: %v - %w", tx.ID, attempt.Hash, err)
		}
		switch apiResponse.Status {
		case ApiStatusPending, ApiStatusIncluded:
			return false, nil
		case ApiStatusFailed, ApiStatusCancelled:
			return true, nil
		case ApiStatusUnknown:
			continue
		default:
			continue
		}
	}
	return false, nil
}
