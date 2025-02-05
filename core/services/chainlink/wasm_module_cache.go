package chainlink

import (
	"fmt"
	"os"
	"sync"

	"github.com/bytecodealliance/wasmtime-go/v28"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
)

type wasmModuleCacheStats interface {
	OnCacheHit()
	OnCacheMiss()
	OnCacheEviction(int)
	OnCacheAddition()
}

type fileBasedModuleCache struct {
	lggr     logger.Logger
	cacheDir string
	stats    wasmModuleCacheStats

	wasmLocks sync.Map
}

type WasmModuleCache interface {
	GetModuleFromBinaryID(binaryID host.BinaryID, initialFuel uint64) (*host.WasmTimeModule, error)
	GetModuleFromBinary(binary []byte, initialFuel uint64, isUncompressed bool,
		maxCompressedBinarySize uint64, maxDecompressedBinarySize uint64) (*host.WasmTimeModule, error)
}

func NewFileBasedModuleCache(lggr logger.Logger, cacheDir string, stats wasmModuleCacheStats) (*fileBasedModuleCache, error) {
	err := os.MkdirAll(cacheDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create cache directory: %w", err)
	}

	return &fileBasedModuleCache{lggr: lggr, cacheDir: cacheDir, stats: stats}, nil
}

type wasmLock struct {
	mux *sync.Mutex
}

var ErrSerialisedModuleNotFound = fmt.Errorf("serialised module not found")

// GetModuleFromBinaryID creates a new module if a serialised module is found in the cache, else it returns an
// ErrSerialisedModuleNotFound error.
func (f *fileBasedModuleCache) GetModuleFromBinaryID(binaryID host.BinaryID, initialFuel uint64) (*host.WasmTimeModule, error) {
	return f.getModule(binaryID, nil, initialFuel, false, 0, 0)
}

// GetModuleFromBinary creates a new module from the given binary, caching a serialised version of the module.  The returned
// module must be Closed.
func (f *fileBasedModuleCache) GetModuleFromBinary(binary []byte, initialFuel uint64, isUncompressed bool,
	maxCompressedBinarySize uint64, maxDecompressedBinarySize uint64) (*host.WasmTimeModule, error) {
	binaryID := host.FromBinary(binary)

	return f.getModule(binaryID, binary, initialFuel, isUncompressed, maxCompressedBinarySize, maxDecompressedBinarySize)
}

// getModule will return a module from the cache if it exists, otherwise if the binary for the module is provided it
// will create a new module and cache it.  If the binary is not provided, it will return an ErrSerialisedModuleNotFound error.
func (f *fileBasedModuleCache) getModule(binaryID host.BinaryID, binary []byte, initialFuel uint64,
	isUncompressed bool, maxCompressedBinarySize uint64, maxDecompressedBinarySize uint64) (*host.WasmTimeModule, error) {

	lock := f.getLockForBinaryID(binaryID)
	lock.mux.Lock()
	defer lock.mux.Unlock()

	cacheFile := f.getCacheFilePath(binaryID)

	engineCfg, err := host.GetEngineConfiguration(initialFuel)
	if err != nil {
		return nil, fmt.Errorf("failed to get engine configuration: %w", err)
	}
	engine := wasmtime.NewEngineWithConfig(engineCfg)

	var module *wasmtime.Module
	if f.isCacheFileMissing(cacheFile) {
		if binary == nil {
			return nil, ErrSerialisedModuleNotFound
		}
		module, err = f.createAndCacheNewModule(engine, binary, cacheFile, isUncompressed, maxCompressedBinarySize, maxDecompressedBinarySize)
		if err != nil {
			return nil, fmt.Errorf("failed to create and cache module: %w", err)
		}
	} else {
		module, err = f.loadModuleFromCache(engine, cacheFile, binary)
		if err != nil {
			return nil, fmt.Errorf("failed to load module from cache: %w", err)
		}
	}

	return host.NewWasmTimeModule(engine, module, engineCfg), nil
}

func (f *fileBasedModuleCache) getLockForBinaryID(binaryID host.BinaryID) *wasmLock {
	lock := wasmLock{mux: &sync.Mutex{}}
	actual, _ := f.wasmLocks.LoadOrStore(binaryID, &lock)
	return actual.(*wasmLock)
}

func (f *fileBasedModuleCache) getCacheFilePath(binaryID host.BinaryID) string {
	return f.cacheDir + "/" + string(binaryID) + ".serialized"
}

func (f *fileBasedModuleCache) isCacheFileMissing(cacheFile string) bool {
	_, err := os.Stat(cacheFile)
	return os.IsNotExist(err)
}

func (f *fileBasedModuleCache) createAndCacheNewModule(engine *wasmtime.Engine, binary []byte, cacheFile string,
	isUncompressed bool, maxCompressedBinarySize uint64, maxDecompressedBinarySize uint64) (*wasmtime.Module, error) {

	binary, err := host.ValidateAndDecompressBinary(binary, isUncompressed, maxCompressedBinarySize, maxDecompressedBinarySize)
	if err != nil {
		return nil, fmt.Errorf("failed to validate and decompress binary: %w", err)
	}

	module, err := f.createAndCacheModule(engine, binary, cacheFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create and cache module: %w", err)
	}
	return module, nil
}

func (f *fileBasedModuleCache) loadModuleFromCache(engine *wasmtime.Engine, cacheFile string, binary []byte) (*wasmtime.Module, error) {
	f.stats.OnCacheHit()
	mod, err := wasmtime.NewModuleDeserializeFile(engine, cacheFile)
	if err != nil {
		f.lggr.Warn("Failed to deserialize module from cache, recreating and caching module", "cacheFile", cacheFile, "error", err)
		mod, err = f.createAndCacheModule(engine, binary, cacheFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create and cache module: %w", err)
		}
	}
	return mod, nil
}

func (f *fileBasedModuleCache) createAndCacheModule(engine *wasmtime.Engine, wasm []byte, cacheFile string) (*wasmtime.Module, error) {
	f.stats.OnCacheMiss()
	module, err := wasmtime.NewModule(engine, wasm)
	if err != nil {
		return nil, fmt.Errorf("failed to create module: %w", err)
	}

	serialisedBytes, err := module.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialise module: %w", err)
	}

	err = os.WriteFile(cacheFile, serialisedBytes, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to write serialised bytes to cache: %w", err)
	}
	f.stats.OnCacheAddition()

	return module, nil
}

type NoopWasmModuleCacheStats struct{}

func (n *NoopWasmModuleCacheStats) OnCacheHit() {}

func (n *NoopWasmModuleCacheStats) OnCacheMiss() {}

func (n *NoopWasmModuleCacheStats) OnCacheEviction(int) {}

func (n *NoopWasmModuleCacheStats) OnCacheAddition() {}
