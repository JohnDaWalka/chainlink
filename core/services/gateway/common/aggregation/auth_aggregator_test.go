package aggregation

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func TestAuthAggregator_StartStop(t *testing.T) {
	lggr := logger.Test(t)
	agg := NewAuthAggregator(lggr, 2, 100*time.Millisecond)

	ctx := testutils.Context(t)

	err := agg.Start(ctx)
	require.NoError(t, err)
	require.Equal(t, "Started", agg.State())

	// Test that starting again returns error
	err = agg.Start(ctx)
	require.Error(t, err)

	err = agg.Close()
	require.NoError(t, err)
	require.Equal(t, "Stopped", agg.State())

	// Test that closing again returns Error
	err = agg.Close()
	require.Error(t, err)
}

func getRandomECDSAPublicKey(t *testing.T) string {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	return crypto.PubkeyToAddress(key.PublicKey).Hex()
}

func TestWorkflowAuthObservation_Digest(t *testing.T) {
	publicKey1 := getRandomECDSAPublicKey(t)
	publicKey2 := getRandomECDSAPublicKey(t)
	observation1 := WorkflowAuthObservation{
		WorkflowID: "workflow-1",
		AuthorizedKey: gateway_common.AuthorizedKey{
			KeyType:   gateway_common.KeyTypeECDSA,
			PublicKey: publicKey1,
		},
	}

	observation2 := WorkflowAuthObservation{
		WorkflowID: "workflow-1",
		AuthorizedKey: gateway_common.AuthorizedKey{
			KeyType:   gateway_common.KeyTypeECDSA,
			PublicKey: publicKey1,
		},
	}

	observation3 := WorkflowAuthObservation{
		WorkflowID: "workflow-2",
		AuthorizedKey: gateway_common.AuthorizedKey{
			KeyType:   gateway_common.KeyTypeECDSA,
			PublicKey: publicKey1,
		},
	}

	observation4 := WorkflowAuthObservation{
		WorkflowID: "workflow-1",
		AuthorizedKey: gateway_common.AuthorizedKey{
			KeyType:   "secp256k1",
			PublicKey: publicKey1,
		},
	}

	observation5 := WorkflowAuthObservation{
		WorkflowID: "workflow-1",
		AuthorizedKey: gateway_common.AuthorizedKey{
			KeyType:   gateway_common.KeyTypeECDSA,
			PublicKey: publicKey2,
		},
	}

	require.Equal(t, observation1.Digest(), observation2.Digest())
	require.NotEqual(t, observation1.Digest(), observation3.Digest())
	require.NotEqual(t, observation1.Digest(), observation4.Digest())
	require.NotEqual(t, observation1.Digest(), observation5.Digest())
}

func TestAuthAggregator_Collect(t *testing.T) {
	lggr := logger.Test(t)
	agg := NewAuthAggregator(lggr, 2, 10*time.Second)

	observation := WorkflowAuthObservation{
		WorkflowID: "workflow-1",
		AuthorizedKey: gateway_common.AuthorizedKey{
			KeyType:   gateway_common.KeyTypeECDSA,
			PublicKey: getRandomECDSAPublicKey(t),
		},
	}

	err := agg.Collect(observation, "node1")
	require.NoError(t, err)
	require.Len(t, agg.observations, 1)
	require.Len(t, agg.observedAt, 1)
	require.Len(t, agg.observedAt["node1"], 1)

	digest := observation.Digest()
	nodeObs, exists := agg.observations[digest]
	require.True(t, exists)
	require.Equal(t, observation, nodeObs.observation)
	require.True(t, nodeObs.nodes.Contains("node1"))
	require.Len(t, nodeObs.nodes, 1)
	timestamp1, ok := agg.observedAt["node1"][digest]
	require.True(t, ok)

	// Test collecting from second node with same observation
	err = agg.Collect(observation, "node2")
	require.NoError(t, err)
	require.Len(t, agg.observations, 1)
	require.Len(t, agg.observedAt, 2)

	nodeObs, exists = agg.observations[digest]
	require.True(t, exists)
	require.True(t, nodeObs.nodes.Contains("node1"))
	require.True(t, nodeObs.nodes.Contains("node2"))
	require.Len(t, nodeObs.nodes, 2)

	// Test collecting from same node again (should update timestamp)
	time.Sleep(10 * time.Millisecond) // Small delay to ensure timestamp difference
	err = agg.Collect(observation, "node1")
	require.NoError(t, err)
	require.Len(t, agg.observations, 1)
	require.Len(t, agg.observedAt, 2)

	nodeObs, exists = agg.observations[digest]
	require.True(t, exists)
	require.Len(t, nodeObs.nodes, 2)

	digests, ok := agg.observedAt["node1"]
	require.True(t, ok)
	timestamp2, ok := digests[digest]
	require.True(t, ok)
	require.NotEqual(t, timestamp1, timestamp2)
}

func TestAuthAggregator_CollectDifferentObservations(t *testing.T) {
	lggr := logger.Test(t)
	agg := NewAuthAggregator(lggr, 2, 10*time.Second)

	observation1 := WorkflowAuthObservation{
		WorkflowID: "workflow-1",
		AuthorizedKey: gateway_common.AuthorizedKey{
			KeyType:   gateway_common.KeyTypeECDSA,
			PublicKey: getRandomECDSAPublicKey(t),
		},
	}

	observation2 := WorkflowAuthObservation{
		WorkflowID: "workflow-1",
		AuthorizedKey: gateway_common.AuthorizedKey{
			KeyType:   gateway_common.KeyTypeECDSA,
			PublicKey: getRandomECDSAPublicKey(t),
		},
	}

	// Collect different observations
	err := agg.Collect(observation1, "node1")
	require.NoError(t, err)
	err = agg.Collect(observation2, "node2")
	require.NoError(t, err)

	// Should have 2 different observations
	require.Len(t, agg.observations, 2)
	require.Len(t, agg.observedAt, 2)

	digest1 := observation1.Digest()
	digest2 := observation2.Digest()

	nodeObs1, exists := agg.observations[digest1]
	require.True(t, exists)
	require.Equal(t, observation1, nodeObs1.observation)
	require.True(t, nodeObs1.nodes.Contains("node1"))

	nodeObs2, exists := agg.observations[digest2]
	require.True(t, exists)
	require.Equal(t, observation2, nodeObs2.observation)
	require.True(t, nodeObs2.nodes.Contains("node2"))
}

func TestAuthAggregator_Aggregate(t *testing.T) {
	lggr := logger.Test(t)
	threshold := 2
	agg := NewAuthAggregator(lggr, threshold, 10*time.Second)

	publicKey1 := getRandomECDSAPublicKey(t)
	publicKey2 := getRandomECDSAPublicKey(t)
	publicKey3 := getRandomECDSAPublicKey(t)

	observation1 := WorkflowAuthObservation{
		WorkflowID: "workflow-1",
		AuthorizedKey: gateway_common.AuthorizedKey{
			KeyType:   gateway_common.KeyTypeECDSA,
			PublicKey: publicKey1,
		},
	}

	observation2 := WorkflowAuthObservation{
		WorkflowID: "workflow-1",
		AuthorizedKey: gateway_common.AuthorizedKey{
			KeyType:   gateway_common.KeyTypeECDSA,
			PublicKey: publicKey2,
		},
	}

	observation3 := WorkflowAuthObservation{
		WorkflowID: "workflow-2",
		AuthorizedKey: gateway_common.AuthorizedKey{
			KeyType:   gateway_common.KeyTypeECDSA,
			PublicKey: publicKey3,
		},
	}

	// Test aggregation with no observations
	result, err := agg.Aggregate()
	require.NoError(t, err)
	require.Empty(t, result)

	// Add observations below threshold
	err = agg.Collect(observation1, "node1")
	require.NoError(t, err)
	err = agg.Collect(observation2, "node1")
	require.NoError(t, err)
	err = agg.Collect(observation3, "node1")
	require.NoError(t, err)

	result, err = agg.Aggregate()
	require.NoError(t, err)
	require.Empty(t, result)

	// Add observations to reach threshold for workflow-1
	err = agg.Collect(observation1, "node2")
	require.NoError(t, err)
	err = agg.Collect(observation2, "node2")
	require.NoError(t, err)

	result, err = agg.Aggregate()
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, "workflow-1", result[0].WorkflowID)
	require.Len(t, result[0].AuthorizedKeys, 2)

	// Add observation to reach threshold for workflow-2
	err = agg.Collect(observation3, "node2")
	require.NoError(t, err)

	result, err = agg.Aggregate()
	require.NoError(t, err)
	require.Len(t, result, 2)
}

func TestAuthAggregator_ReapObservations(t *testing.T) {
	lggr := logger.Test(t)
	cleanupInterval := 1 * time.Second
	agg := NewAuthAggregator(lggr, 2, cleanupInterval)
	observation1 := WorkflowAuthObservation{
		WorkflowID: "workflow-1",
		AuthorizedKey: gateway_common.AuthorizedKey{
			KeyType:   gateway_common.KeyTypeECDSA,
			PublicKey: "test-public-key-1",
		},
	}

	observation2 := WorkflowAuthObservation{
		WorkflowID: "workflow-1",
		AuthorizedKey: gateway_common.AuthorizedKey{
			KeyType:   gateway_common.KeyTypeECDSA,
			PublicKey: "test-public-key-2",
		},
	}
	err := agg.Collect(observation1, "node1")
	require.NoError(t, err)
	err = agg.Collect(observation1, "node2")
	require.NoError(t, err)
	err = agg.Collect(observation2, "node1")
	require.NoError(t, err)
	require.Len(t, agg.observations, 2)
	require.Len(t, agg.observedAt, 2)

	err = agg.Start(testutils.Context(t))
	require.NoError(t, err)
	// Wait for cleanup interval to pass
	time.Sleep(cleanupInterval + 100*time.Millisecond)
	require.Empty(t, agg.observations)
}

func TestAuthAggregator_ReapObservations_UnexpiredObservation(t *testing.T) {
	lggr := logger.Test(t)
	cleanupInterval := 1 * time.Second
	agg := NewAuthAggregator(lggr, 2, cleanupInterval)
	observation := WorkflowAuthObservation{
		WorkflowID: "workflow-1",
		AuthorizedKey: gateway_common.AuthorizedKey{
			KeyType:   gateway_common.KeyTypeECDSA,
			PublicKey: "test-public-key-1",
		},
	}

	err := agg.Collect(observation, "node1")
	require.NoError(t, err)
	err = agg.Collect(observation, "node2")
	require.NoError(t, err)
	require.Len(t, agg.observations, 1)
	require.Len(t, agg.observedAt, 2)

	// Wait for cleanup interval to pass
	time.Sleep(cleanupInterval + 100*time.Millisecond)
	// Add observation from node3 (fresh)
	err = agg.Collect(observation, "node3")
	require.NoError(t, err)
	// Manually trigger cleanup
	agg.reapObservations()

	require.Len(t, agg.observations, 1)
	o, ok := agg.observations[observation.Digest()]
	require.True(t, ok)
	require.Equal(t, observation, o.observation)
	require.True(t, o.nodes.Contains("node3"))
}

func TestAuthAggregator_Collect_EdgeCases(t *testing.T) {
	lggr := logger.Test(t)

	t.Run("empty workflow ID", func(t *testing.T) {
		agg := NewAuthAggregator(lggr, 1, 10*time.Second)

		observation := WorkflowAuthObservation{
			WorkflowID: "",
			AuthorizedKey: gateway_common.AuthorizedKey{
				KeyType:   gateway_common.KeyTypeECDSA,
				PublicKey: getRandomECDSAPublicKey(t),
			},
		}

		err := agg.Collect(observation, "node1")
		require.Error(t, err)
	})

	t.Run("empty public key", func(t *testing.T) {
		agg := NewAuthAggregator(lggr, 1, 10*time.Second)

		observation := WorkflowAuthObservation{
			WorkflowID: "workflow-1",
			AuthorizedKey: gateway_common.AuthorizedKey{
				KeyType:   gateway_common.KeyTypeECDSA,
				PublicKey: "",
			},
		}

		err := agg.Collect(observation, "node1")
		require.Error(t, err)
	})

	t.Run("empty node address", func(t *testing.T) {
		agg := NewAuthAggregator(lggr, 1, 10*time.Second)

		observation := WorkflowAuthObservation{
			WorkflowID: "workflow-1",
			AuthorizedKey: gateway_common.AuthorizedKey{
				KeyType:   gateway_common.KeyTypeECDSA,
				PublicKey: getRandomECDSAPublicKey(t),
			},
		}

		err := agg.Collect(observation, "")
		require.Error(t, err)
	})
}
