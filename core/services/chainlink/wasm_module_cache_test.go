package chainlink

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/bytecodealliance/wasmtime-go/v28"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func TestNewCachedWasmModuleFactory(t *testing.T) {
	lggr := logger.Test(t)
	cacheDir := t.TempDir()

	factory, err := NewCachedWasmModuleFactory(lggr, cacheDir)
	require.NoError(t, err)
	require.NotNil(t, factory)
}

func TestNewModule(t *testing.T) {
	lggr := logger.Test(t)
	cacheDir := t.TempDir()

	factory, err := NewCachedWasmModuleFactory(lggr, cacheDir)
	require.NoError(t, err)

	engine := wasmtime.NewEngine()
	wasm := []byte(`(module (func (export "run") (nop)))`)
	isUncompressed := true
	maxCompressedBinarySize := uint64(1024)
	maxDecompressedBinarySize := uint64(1024)

	module, err := factory.NewModule(engine, wasm, isUncompressed, maxCompressedBinarySize, maxDecompressedBinarySize)
	require.NoError(t, err)
	require.NotNil(t, module)

	sha := sha256.Sum256(wasm)
	wasmBytesID := hex.EncodeToString(sha[:])
	cacheFile := filepath.Join(cacheDir, wasmBytesID+".serialized")

	_, err = os.Stat(cacheFile)
	require.NoError(t, err)
}
