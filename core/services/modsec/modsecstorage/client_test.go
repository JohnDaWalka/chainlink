package modsecstorage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStdClient(t *testing.T) {
	t.Parallel()

	// Create a test server.
	server, cleanup := NewTestServer()
	t.Cleanup(cleanup)

	// Create a client that points to the test server.
	client := NewStdClient(server.URL)

	// Run test cases.
	t.Run("SetAndGet", func(t *testing.T) {
		key := "my-key"
		value := []byte("my-value")

		// Set a value.
		err := client.Set(t.Context(), key, value)
		require.NoError(t, err)

		// Get it back.
		retrieved, err := client.Get(t.Context(), key)
		require.NoError(t, err)
		require.Equal(t, value, retrieved)
	})

	t.Run("GetNotFound", func(t *testing.T) {
		_, err := client.Get(t.Context(), "non-existent-key")
		require.Error(t, err)
		require.Contains(t, err.Error(), "404 Not Found")
	})
}
