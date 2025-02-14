package por

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	libcrecli "github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
)

func CreateConfigFile(t *testing.T, feedsConsumerAddress common.Address, feedID, dataURL string) *os.File {
	configFile, err := os.CreateTemp("", "config.json")
	require.NoError(t, err, "failed to create workflow config file")

	cleanFeedID := strings.TrimPrefix(feedID, "0x")
	feedLength := len(cleanFeedID)

	require.GreaterOrEqual(t, feedLength, 32, "feed ID must be at least 32 characters long")

	if feedLength > 32 {
		cleanFeedID = cleanFeedID[:32]
	}

	feedIDToUse := "0x" + cleanFeedID

	workflowConfig := libcrecli.PoRWorkflowConfig{
		FeedID:          feedIDToUse,
		URL:             dataURL,
		ConsumerAddress: feedsConsumerAddress.Hex(),
	}

	configMarshalled, err := json.Marshal(workflowConfig)
	require.NoError(t, err, "failed to marshal workflow config")

	_, err = configFile.Write(configMarshalled)
	require.NoError(t, err, "failed to write workflow config file")

	return configFile
}
