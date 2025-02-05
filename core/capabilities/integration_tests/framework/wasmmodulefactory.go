package framework

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/bytecodealliance/wasmtime-go/v28"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

type WasmModuleCacheFactory interface {
	NewWasmModuleCacheForPeer(peerID string) (chainlink.WasmModuleCache, error)
}

type inMemoryWasmModuleCacheFactory struct {
}

func NewInMemoryWasmModuleCacheFactory() WasmModuleCacheFactory {
	return &inMemoryWasmModuleCacheFactory{}
}

func (f *inMemoryWasmModuleCacheFactory) NewWasmModuleCacheForPeer(peerID string) (chainlink.WasmModuleCache, error) {
	return NewInMemoryWasmModuleCache(), nil
}

type InMemoryWasmModuleCache struct {
	binaryIDToModule map[host.BinaryID]*host.WasmTimeModule
	mux              sync.Mutex
}

func NewInMemoryWasmModuleCache() *InMemoryWasmModuleCache {
	return &InMemoryWasmModuleCache{
		binaryIDToModule: make(map[host.BinaryID]*host.WasmTimeModule),
	}
}

func (c *InMemoryWasmModuleCache) GetModuleFromBinaryID(binaryID host.BinaryID, initialFuel uint64) (*host.WasmTimeModule, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	module, ok := c.binaryIDToModule[binaryID]
	if !ok {
		return nil, chainlink.ErrSerialisedModuleNotFound
	}

	return module, nil
}

func (c *InMemoryWasmModuleCache) GetModuleFromBinary(binary []byte, initialFuel uint64, isUncompressed bool,
	maxCompressedBinarySize uint64, maxDecompressedBinarySize uint64) (*host.WasmTimeModule, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	engineCfg, err := host.GetEngineConfiguration(initialFuel)
	if err != nil {
		return nil, fmt.Errorf("failed to get engine configuration: %w", err)
	}

	engine := wasmtime.NewEngineWithConfig(engineCfg)

	wasm, err := host.ValidateAndDecompressBinary(binary, isUncompressed, maxCompressedBinarySize, maxDecompressedBinarySize)
	if err != nil {
		return nil, fmt.Errorf("failed to validate and decompress binary: %w", err)
	}

	module, err := wasmtime.NewModule(engine, wasm)
	if err != nil {
		return nil, fmt.Errorf("failed to create module: %w", err)
	}

	wasmModule := host.NewWasmTimeModule(engine, module, engineCfg)
	binaryID := host.FromBinary(wasm)

	c.binaryIDToModule[binaryID] = wasmModule

	return wasmModule, nil
}

type fileBasedWasmModuleCachedFactory struct {
	lggr     logger.Logger
	cacheDir string
}

func NewFileBasedWasmModuleCachedFactory(lggr logger.Logger, cacheDir string) (WasmModuleCacheFactory, error) {
	cacheDir, err := filepath.Abs(cacheDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for cache directory: %w", err)
	}
	return &fileBasedWasmModuleCachedFactory{lggr: lggr, cacheDir: cacheDir}, nil
}

func (f *fileBasedWasmModuleCachedFactory) NewWasmModuleCacheForPeer(peerID string) (chainlink.WasmModuleCache, error) {

	factory, err := chainlink.NewFileBasedModuleCache(f.lggr, filepath.Join(f.cacheDir, peerID), &chainlink.NoopWasmModuleCacheStats{})
	if err != nil {
		return nil, fmt.Errorf("failed to create cached wasm module factory: %w", err)
	}

	return factory, nil
}
