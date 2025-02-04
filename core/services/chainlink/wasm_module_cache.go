package chainlink

import (
	"crypto/sha256"
	"encoding/hex"
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

type cachedWasmModuleFactory struct {
	lggr     logger.Logger
	cacheDir string
	stats    wasmModuleCacheStats

	wasmLocks sync.Map
}

func NewCachedWasmModuleFactory(lggr logger.Logger, cacheDir string, stats wasmModuleCacheStats) (*cachedWasmModuleFactory, error) {
	err := os.MkdirAll(cacheDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create cache directory: %w", err)
	}

	return &cachedWasmModuleFactory{lggr: lggr, cacheDir: cacheDir, stats: stats}, nil
}

type wasmLock struct {
	mux *sync.Mutex
}

func (f *cachedWasmModuleFactory) NewModule(engine *wasmtime.Engine, wasm []byte, isUncompressed bool, maxCompressedBinarySize uint64, maxDecompressedBinarySize uint64) (*wasmtime.Module, error) {
	sha := sha256.Sum256(wasm)
	wasmBytesID := hex.EncodeToString(sha[:])

	lock := wasmLock{mux: &sync.Mutex{}}
	actual, _ := f.wasmLocks.LoadOrStore(wasmBytesID, &lock)
	lock = *actual.(*wasmLock)
	lock.mux.Lock()
	defer lock.mux.Unlock()

	cacheFile := f.cacheDir + "/" + wasmBytesID[:] + ".serialized"

	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		wasm, err = host.ValidateAndDecompressBinary(wasm, isUncompressed, maxCompressedBinarySize, maxDecompressedBinarySize)
		if err != nil {
			return nil, fmt.Errorf("failed to validate and decompress binary: %w", err)
		}

		module, err := f.createAndCacheModule(engine, wasm, cacheFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create and cache module: %w", err)
		}
		return module, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to check if cache file exists: %w", err)
	} else {
		f.stats.OnCacheHit()
		mod, err := wasmtime.NewModuleDeserializeFile(engine, cacheFile)
		if err != nil {
			f.lggr.Warn("Failed to deserialize module from cache, recreating and caching module", "cacheFile", cacheFile, "error", err)
			// If the module fails to deserialize from the cache create it from new and cache it again
			// This is to handle the case where the engined configuration has changed or the cache is corrupted
			mod, err = f.createAndCacheModule(engine, wasm, cacheFile)
			if err != nil {
				return nil, fmt.Errorf("failed to create and cache module: %w", err)
			}
		}
		return mod, nil
	}
}

func (f *cachedWasmModuleFactory) createAndCacheModule(engine *wasmtime.Engine, wasm []byte, cacheFile string) (*wasmtime.Module, error) {
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
