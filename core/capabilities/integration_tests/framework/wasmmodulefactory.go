package framework

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/bytecodealliance/wasmtime-go/v28"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
)

type WasmModuleFactory interface {
	NewWasmModuleFactoryFnForPeer(peerID string) host.WasmtimeModuleFactoryFn
}

type noCacheModuleFactory struct {
}

func NewNoCacheWasmModuleFactory() WasmModuleFactory {
	return &noCacheModuleFactory{}
}

func (f *noCacheModuleFactory) NewWasmModuleFactoryFnForPeer(peerID string) host.WasmtimeModuleFactoryFn {
	return NoCacheWasmModuleFactory
}

func NoCacheWasmModuleFactory(engine *wasmtime.Engine, wasm []byte, isUncompressed bool, maxCompressedBinarySize uint64, maxDecompressedBinarySize uint64) (*wasmtime.Module, error) {
	wasm, err := host.ValidateAndDecompressBinary(wasm, isUncompressed, maxCompressedBinarySize, maxDecompressedBinarySize)
	if err != nil {
		return nil, fmt.Errorf("failed to validate and decompress binary: %w", err)
	}
	return wasmtime.NewModule(engine, wasm)
}

type cachedWasmModuleFactory struct {
	mux      sync.Mutex
	cacheDir string
}

func NewCachedWasmModuleFactory(cacheDir string) (WasmModuleFactory, error) {
	cacheDir, err := filepath.Abs(cacheDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for cache directory: %w", err)
	}
	return &cachedWasmModuleFactory{cacheDir: cacheDir}, nil
}

func (f *cachedWasmModuleFactory) NewWasmModuleFactoryFnForPeer(peerID string) host.WasmtimeModuleFactoryFn {
	return func(engine *wasmtime.Engine, wasm []byte, isUncompressed bool, maxCompressedBinarySize uint64, maxDecompressedBinarySize uint64) (*wasmtime.Module, error) {

		sha := sha256.Sum256(wasm)
		wasmBytesID := hex.EncodeToString(sha[:])
		cacheFile := f.cacheDir + "/" + peerID + "/" + wasmBytesID[:] + ".serialized"

		if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
			err = os.MkdirAll(f.cacheDir+"/"+peerID, 0755)
			if err != nil {
				return nil, fmt.Errorf("failed to create cache directory: %w", err)
			}

			wasm, err = host.ValidateAndDecompressBinary(wasm, isUncompressed, maxCompressedBinarySize, maxDecompressedBinarySize)
			if err != nil {
				return nil, fmt.Errorf("failed to validate and decompress binary: %w", err)
			}

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
		} else if err != nil {
			return nil, fmt.Errorf("failed to check if cache file exists: %w", err)
		} else {
			mod, err := wasmtime.NewModuleDeserializeFile(engine, cacheFile)
			if err != nil {
				return nil, fmt.Errorf("failed to deserialize module: %w", err)
			}
			return mod, nil
		}
	}
}
