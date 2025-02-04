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

type WasmModuleFactory interface {
	NewWasmModuleFactoryFnForPeer(peerID string) (host.WasmtimeModuleFactoryFn, error)
}

type noCacheModuleFactory struct {
}

func NewNoCacheWasmModuleFactory() WasmModuleFactory {
	return &noCacheModuleFactory{}
}

func (f *noCacheModuleFactory) NewWasmModuleFactoryFnForPeer(peerID string) (host.WasmtimeModuleFactoryFn, error) {
	return NoCacheWasmModuleFactory, nil
}

func NoCacheWasmModuleFactory(engine *wasmtime.Engine, wasm []byte, isUncompressed bool, maxCompressedBinarySize uint64, maxDecompressedBinarySize uint64) (*wasmtime.Module, error) {
	wasm, err := host.ValidateAndDecompressBinary(wasm, isUncompressed, maxCompressedBinarySize, maxDecompressedBinarySize)
	if err != nil {
		return nil, fmt.Errorf("failed to validate and decompress binary: %w", err)
	}
	return wasmtime.NewModule(engine, wasm)
}

type cachedWasmModuleFactory struct {
	lggr     logger.Logger
	mux      sync.Mutex
	cacheDir string
}

func NewCachedWasmModuleFactory(lggr logger.Logger, cacheDir string) (WasmModuleFactory, error) {
	cacheDir, err := filepath.Abs(cacheDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for cache directory: %w", err)
	}
	return &cachedWasmModuleFactory{lggr: lggr, cacheDir: cacheDir}, nil
}

func (f *cachedWasmModuleFactory) NewWasmModuleFactoryFnForPeer(peerID string) (host.WasmtimeModuleFactoryFn, error) {

	factory, err := chainlink.NewCachedWasmModuleFactory(f.lggr, filepath.Join(f.cacheDir, peerID))
	if err != nil {
		return nil, fmt.Errorf("failed to create cached wasm module factory: %w", err)
	}

	return factory.NewModule, nil
}
