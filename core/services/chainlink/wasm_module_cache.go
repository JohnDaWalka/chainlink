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

type cachedWasmModuleFactory struct {
	lggr     logger.Logger
	cacheDir string
	mux      sync.Mutex
}

func NewCachedWasmModuleFactory(lggr logger.Logger, cacheDir string) (*cachedWasmModuleFactory, error) {
	err := os.MkdirAll(cacheDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create cache directory: %w", err)
	}

	return &cachedWasmModuleFactory{lggr: lggr, cacheDir: cacheDir}, nil
}

func (f *cachedWasmModuleFactory) NewModule(engine *wasmtime.Engine, wasm []byte, isUncompressed bool, maxCompressedBinarySize uint64, maxDecompressedBinarySize uint64) (*wasmtime.Module, error) {
	f.mux.Lock()
	defer f.mux.Unlock()

	sha := sha256.Sum256(wasm)
	wasmBytesID := hex.EncodeToString(sha[:])
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
	module, err := wasmtime.NewModule(engine, wasm)

	if err != nil {
		return nil, fmt.Errorf("failed to create module: %w", err)
	}

	serialisedBytes, err := module.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialise module: %w", err)
	}

	err = os.WriteFile(cacheFile, serialisedBytes, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write serialised bytes to cache: %w", err)
	}
	return module, nil
}
