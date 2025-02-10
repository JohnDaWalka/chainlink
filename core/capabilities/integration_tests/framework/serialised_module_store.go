package framework

import (
	"path/filepath"
	"sync"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/artifacts"
)

type SerialisedModuleStoreFactory interface {
	GetModuleStoreForPeer(peerID string) (artifacts.SerialisedModuleStore, error)
}

type FileBasedSerialisedModuleStoreFactory struct {
	storagePath string
}

func NewFileBasedSerialisedModuleStoreFactory(storagePath string) *FileBasedSerialisedModuleStoreFactory {
	return &FileBasedSerialisedModuleStoreFactory{
		storagePath: storagePath,
	}
}

func (f *FileBasedSerialisedModuleStoreFactory) GetModuleStoreForPeer(peerID string) (artifacts.SerialisedModuleStore, error) {
	return artifacts.NewFileBasedModuleStore(filepath.Join(f.storagePath, peerID))
}

type NoStoreModuleStoreFactory struct {
}

func NewNoStoreModuleStoreFactory() *NoStoreModuleStoreFactory {
	return &NoStoreModuleStoreFactory{}
}

func (n *NoStoreModuleStoreFactory) GetModuleStoreForPeer(peerID string) (artifacts.SerialisedModuleStore, error) {
	return NewNoStoreModuleStore()
}

type NoStoreModuleStore struct {
	mux                  sync.Mutex
	workflowIDToBinaryID map[string]string
}

func NewNoStoreModuleStore() (*NoStoreModuleStore, error) {
	return &NoStoreModuleStore{workflowIDToBinaryID: map[string]string{}}, nil
}

func (s *NoStoreModuleStore) StoreModule(workflowID string, binaryID string, module []byte) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.workflowIDToBinaryID[workflowID] = binaryID
	return nil
}

func (s *NoStoreModuleStore) GetModulePath(workflowID string) (string, bool, error) {
	return "", false, nil
}

// GetBinaryID retrieves the wasm binary ID for the given workflow ID. Returns the binary ID, a boolean indicating
// whether the binary ID for the module was found, and an error if one occurred.
func (s *NoStoreModuleStore) GetBinaryID(workflowID string) (string, bool, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	binaryID, ok := s.workflowIDToBinaryID[workflowID]
	return binaryID, ok, nil
}

func (s *NoStoreModuleStore) DeleteModule(workflowID string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	delete(s.workflowIDToBinaryID, workflowID)
	return nil
}
