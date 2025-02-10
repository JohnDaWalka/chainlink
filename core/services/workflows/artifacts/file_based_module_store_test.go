package artifacts

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSerialisedModuleStore(t *testing.T) {
	storeDir := t.TempDir()
	store, err := NewFileBasedModuleStore(storeDir)
	storeDir = filepath.Join(storeDir, "serialised_modules")

	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	workflowID := "test-workflow"
	dataString := "test module data"

	moduleData := []byte(dataString)

	// Test StoreModule
	binaryID := "binaryID"
	err = store.StoreModule(workflowID, binaryID, moduleData)
	require.NoError(t, err)

	// Verify the file was created
	filePath := filepath.Join(storeDir, workflowID, binaryID)
	_, err = os.Stat(filePath)
	require.NoError(t, err)

	// Test GetModule
	path, exists, err := store.GetModulePath(workflowID)
	require.NoError(t, err)
	assert.True(t, exists, "expected module to exist")

	// read in the file at path and check it matches original data
	bytes, err := os.ReadFile(path)
	require.NoError(t, err)

	assert.Equal(t, moduleData, bytes, "expected module data to match")

	// Test getting the binary ID for a given workflow ID
	binaryID, exists, err = store.GetBinaryID(workflowID)
	require.NoError(t, err)
	assert.True(t, exists, "expected binary ID to exist")
	assert.Equal(t, "binaryID", binaryID, "expected binary ID to match")

	// Test DeleteModule
	err = store.DeleteModule(workflowID)
	if err != nil {
		t.Fatalf("failed to delete module: %v", err)
	}

	// Verify the file was deleted
	if _, err = os.Stat(filePath); !os.IsNotExist(err) {
		t.Fatalf("expected file %s to be deleted", filePath)
	}

	// Test GetModule for non-existent module
	_, exists, err = store.GetModulePath(workflowID)
	require.NoError(t, err)
	assert.False(t, exists, "expected module to not exist")
}
