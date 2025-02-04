package chainlink

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"

	"github.com/bytecodealliance/wasmtime-go/v28"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

type bytePair struct {
	wasmBytes   []byte
	moduleBytes []byte
}

func TestNewModule(t *testing.T) {
	lggr := logger.Test(t)

	wasmDir := t.TempDir()
	cacheDir := t.TempDir()
	stats := &testCacheStats{}
	factory, err := NewCachedWasmModuleFactory(lggr, cacheDir, stats)
	require.NoError(t, err)

	// generate 10 unique wasmBytes and get and store the moduleBytes as a pair
	var bytePairs []bytePair
	for i := 0; i < 10; i++ {
		wasmBytes := getWasmBytes(t, fmt.Sprintf("Hello, WebAssembly! %d", i), wasmDir)
		moduleBytes, err := getSerializedModule(factory, wasmBytes)
		require.NoError(t, err)
		bytePairs = append(bytePairs, bytePair{wasmBytes: wasmBytes, moduleBytes: moduleBytes})
	}

	assert.Equal(t, 10, stats.GetMisses())
	assert.Equal(t, 10, stats.GetAdditions())

	// Now retrieve the moduleBytes from the cache  and confirm that the cache stats are updated correctly
	// and that the retrieved module bytes are correct
	for i := 0; i < 10; i++ {
		wasmBytes := bytePairs[i].wasmBytes
		moduleBytes := bytePairs[i].moduleBytes

		moduleBytes2, err := getSerializedModule(factory, wasmBytes)
		require.NoError(t, err)

		assert.Equal(t, moduleBytes, moduleBytes2)
	}

	assert.Equal(t, 10, stats.GetHits())
	assert.Equal(t, 10, stats.GetMisses())
}

// Now test the cache retrieving in a multithreaded environment
func TestNewModuleMultithreaded_Retrieval(t *testing.T) {
	lggr := logger.Test(t)

	wasmDir := t.TempDir()
	cacheDir := t.TempDir()
	stats := &testCacheStats{}
	factory, err := NewCachedWasmModuleFactory(lggr, cacheDir, stats)
	require.NoError(t, err)

	// generate 10 unique wasmBytes and get and store the moduleBytes as a pair
	var bytePairs []bytePair
	for i := 0; i < 10; i++ {
		wasmBytes := getWasmBytes(t, fmt.Sprintf("Hello, WebAssembly! %d", i), wasmDir)
		moduleBytes, err := getSerializedModule(factory, wasmBytes)
		require.NoError(t, err)
		bytePairs = append(bytePairs, bytePair{wasmBytes: wasmBytes, moduleBytes: moduleBytes})
	}

	assert.Equal(t, 10, stats.GetMisses())
	assert.Equal(t, 10, stats.GetAdditions())

	// Now retrieve the moduleBytes from the cache  and confirm that the cache stats are updated correctly
	// and that the retrieved module bytes are correct
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			wasmBytes := bytePairs[i].wasmBytes
			moduleBytes := bytePairs[i].moduleBytes

			retrievedModuleBytes, err := getSerializedModule(factory, wasmBytes)
			require.NoError(t, err)

			assert.Equal(t, moduleBytes, retrievedModuleBytes)
		}(i)
	}
	wg.Wait()

	assert.Equal(t, 10, stats.GetHits())
	assert.Equal(t, 10, stats.GetMisses())
}

// Now test the cache adding and retrieving in a multithreaded environment
func TestNewModuleMultithreaded_Adding(t *testing.T) {
	lggr := logger.Test(t)

	wasmDir := t.TempDir()
	cacheDir := t.TempDir()
	stats := &testCacheStats{}
	factory, err := NewCachedWasmModuleFactory(lggr, cacheDir, stats)
	require.NoError(t, err)

	// generate 10 unique wasmBytes and get and store the moduleBytes as a pair
	var bytePairs []bytePair
	for i := 0; i < 10; i++ {
		wasmBytes := getWasmBytes(t, fmt.Sprintf("Hello, WebAssembly! %d", i), wasmDir)
		moduleBytes, err := getSerializedModule(factory, wasmBytes)
		require.NoError(t, err)
		bytePairs = append(bytePairs, bytePair{wasmBytes: wasmBytes, moduleBytes: moduleBytes})
	}

	// Now using a new cache, add/retrieve the moduleBytes from the cache in a multithreaded environment
	// and confirm that the cache stats are updated correctly and the retrieved module bytes are correct
	cacheDir = t.TempDir()
	stats = &testCacheStats{}
	factory, err = NewCachedWasmModuleFactory(lggr, cacheDir, stats)
	require.NoError(t, err)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			bytes := bytePairs[i].wasmBytes
			moduleBytes, err := getSerializedModule(factory, bytes)
			require.NoError(t, err)
			assert.Equal(t, bytePairs[i].moduleBytes, moduleBytes)
		}(i)

		go func(i int) {
			defer wg.Done()
			bytes := bytePairs[i].wasmBytes
			moduleBytes, err := getSerializedModule(factory, bytes)
			require.NoError(t, err)
			assert.Equal(t, bytePairs[i].moduleBytes, moduleBytes)
		}(i)
	}

	wg.Wait()

	assert.Equal(t, 10, stats.GetHits())
	assert.Equal(t, 10, stats.GetMisses())
	assert.Equal(t, 10, stats.GetAdditions())

}

func getSerializedModule(factory *cachedWasmModuleFactory, wasmBytes []byte) ([]byte, error) {
	engine := wasmtime.NewEngine()
	isUncompressed := true
	maxCompressedBinarySize := uint64(math.MaxUint64)
	maxDecompressedBinarySize := uint64(math.MaxUint64)

	module, err := factory.NewModule(engine, wasmBytes, isUncompressed, maxCompressedBinarySize, maxDecompressedBinarySize)
	if err != nil {
		return nil, err
	}

	moduleBytes, err := module.Serialize()
	if err != nil {
		return nil, err
	}

	return moduleBytes, nil
}

func getWasmBytes(t *testing.T, message string, wasmDir string) []byte {
	mainGoContent := fmt.Sprintf(`package main
import "fmt"

func main() {
    fmt.Println("%s")
}`, message)

	mainGoPath := filepath.Join(wasmDir, "main.go")
	err := os.WriteFile(mainGoPath, []byte(mainGoContent), 0644)
	require.NoError(t, err)

	cmd := exec.Command("go", "build", "-o", filepath.Join(wasmDir, "main.wasm"), mainGoPath)
	cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")
	err = cmd.Run()
	require.NoError(t, err)

	wasmBytes, err := os.ReadFile(filepath.Join(wasmDir, "main.wasm"))
	require.NoError(t, err)
	return wasmBytes
}

type testCacheStats struct {
	hits      int
	misses    int
	evictions int
	additions int
	mu        sync.Mutex
}

func (s *testCacheStats) OnCacheHit() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hits++
}

func (s *testCacheStats) OnCacheMiss() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.misses++
}

func (s *testCacheStats) OnCacheEviction(count int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evictions += count
}

func (s *testCacheStats) OnCacheAddition() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.additions++
}

func (s *testCacheStats) GetHits() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.hits
}

func (s *testCacheStats) GetMisses() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.misses
}

func (s *testCacheStats) GetEvictions() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.evictions
}

func (s *testCacheStats) GetAdditions() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.additions
}
