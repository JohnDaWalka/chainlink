package metadata

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_EnvMetadata_Clone(t *testing.T) {
	original := SimpleEnv{
		DeployCounts: map[uint64]int64{
			1: 10,
			2: 20,
		},
	}

	clone := original.Clone()

	// Verify the clone is equal to the original
	require.Equal(t, original, clone)

	// Modify the clone and check that the original is unaffected
	clone.DeployCounts[1] = 30
	require.NotEqual(t, original.DeployCounts[1], clone.DeployCounts[1])
}
