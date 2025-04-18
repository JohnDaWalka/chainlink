package metadata

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func Test_SimpleContract_Clone(t *testing.T) {
	original := SimpleContract{
		DeployedAt:  time.Now(),
		TxHash:      common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
		BlockNumber: 123456,
	}

	// Clone the original metadata
	cloned := original.Clone()

	// Assert that the cloned metadata matches the original
	require.Equal(t, original, cloned)

	// Modify the cloned metadata and ensure the original is unaffected
	cloned.BlockNumber = 654321
	require.NotEqual(t, original.BlockNumber, cloned.BlockNumber)
}
