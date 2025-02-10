package artifacts

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileBasedModuleStore struct {
	storeDir string
	mux      sync.Mutex
}

func NewFileBasedModuleStore(storeDir string) (*FileBasedModuleStore, error) {
	storeDir = filepath.Join(storeDir, "serialised_modules")
	err := os.MkdirAll(storeDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create store directory: %w", err)
	}
	return &FileBasedModuleStore{storeDir: storeDir}, nil
}

func (s *FileBasedModuleStore) StoreModule(workflowID string, binaryID string, module []byte) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	workflowDir := filepath.Join(s.storeDir, workflowID)
	err := os.MkdirAll(workflowDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create workflow directory: %w", err)
	}

	serialisedModuleFile := filepath.Join(workflowDir, binaryID)
	err = os.WriteFile(serialisedModuleFile, module, 0600)
	if err != nil {
		return fmt.Errorf("failed to store module: %w", err)
	}
	return nil
}

// GetModulePath retrieves the path to the serialised module for the given workflow ID. Returns the path, a boolean indicating
// whether the module was found, and an error if one occurred.
func (s *FileBasedModuleStore) GetModulePath(workflowID string) (string, bool, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	workflowDir := filepath.Join(s.storeDir, workflowID)
	if _, err := os.Stat(workflowDir); os.IsNotExist(err) {
		return "", false, nil
	}

	// Get the first file in the directory
	files, err := os.ReadDir(workflowDir)
	if err != nil {
		return "", false, fmt.Errorf("failed to read the workflow directory: %w", err)
	}
	// Should only be one file in the directory
	if len(files) == 0 {
		return "", false, errors.New("no serialised binary found")
	}

	if len(files) > 1 {
		return "", false, fmt.Errorf("unexpected number of files (%d) in the workflow directory", len(files))
	}

	binaryDir := filepath.Join(workflowDir, files[0].Name())
	return binaryDir, true, nil
}

// GetBinaryID retrieves the wasm binary ID for the given workflow ID. Returns the binary ID, a boolean indicating
// whether the binary ID for the module was found, and an error if one occurred.
func (s *FileBasedModuleStore) GetBinaryID(workflowID string) (string, bool, error) {
	modulePath, exists, err := s.GetModulePath(workflowID)
	if err != nil {
		return "", false, fmt.Errorf("failed to get module path: %w", err)
	}

	if !exists {
		return "", false, nil
	}

	// get the last part of the path, this is the binary ID
	_, binaryID := filepath.Split(modulePath)

	return binaryID, true, nil
}

// DeleteModule deletes the serialised module for the given workflow ID. Returns a boolean indicating whether the module was
// deleted and an error if one occurred.
func (s *FileBasedModuleStore) DeleteModule(workflowID string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	workflowDir := filepath.Join(s.storeDir, workflowID)

	if _, err := os.Stat(workflowDir); os.IsNotExist(err) {
		return nil
	}

	if err := os.RemoveAll(workflowDir); err != nil {
		return fmt.Errorf("failed to delete workflow dir: %w", err)
	}
	return nil
}
