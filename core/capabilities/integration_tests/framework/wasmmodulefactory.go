package framework

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/bytecodealliance/wasmtime-go/v23"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
)

type WasmModuleFactory interface {
	NewWasmModuleFactoryFnForPeer(peerID string) host.WasmModuleFactoryFn
}

type inMemoryWasmModuleFactory struct {
	mux sync.Mutex
}

func NewInMemoryWasmModuleFactory() WasmModuleFactory {
	return &inMemoryWasmModuleFactory{}
}

func (f *inMemoryWasmModuleFactory) NewWasmModuleFactoryFnForPeer(peerID string) host.WasmModuleFactoryFn {
	return func(engine *wasmtime.Engine, wasm []byte) (*wasmtime.Module, error) {
		f.mux.Lock()
		defer f.mux.Unlock()
		return wasmtime.NewModule(engine, wasm)
	}
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

func (f *cachedWasmModuleFactory) NewWasmModuleFactoryFnForPeer(peerID string) host.WasmModuleFactoryFn {
	return func(engine *wasmtime.Engine, wasm []byte) (*wasmtime.Module, error) {

		sha := sha256.Sum256(wasm)
		wasmBytesID := hex.EncodeToString(sha[:])
		cacheFile := f.cacheDir + "/" + peerID + "/" + wasmBytesID[:] + ".serialized"

		if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
			err = os.MkdirAll(f.cacheDir+"/"+peerID, 0755)
			if err != nil {
				return nil, fmt.Errorf("failed to create cache directory: %w", err)
			}

			f.mux.Lock()
			module, err := wasmtime.NewModule(engine, wasm)
			f.mux.Unlock()
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
			f.mux.Lock()
			start := time.Now()
			fmt.Printf("NEWMODULE START\n")
			mod, err := wasmtime.NewModuleDeserializeFile(engine, cacheFile)
			fmt.Printf("NEWMODULE DONE: %s\n", time.Since(start))
			f.mux.Unlock()
			if err != nil {
				return nil, fmt.Errorf("failed to deserialize module: %w", err)
			}
			return mod, nil
		}
	}
}
