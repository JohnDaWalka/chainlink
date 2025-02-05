package chainlink

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
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
	factory, err := NewFileBasedModuleCache(lggr, cacheDir, stats)
	require.NoError(t, err)

	// generate 10 unique wasmBytes and get and store the module as a pair
	var modules []*host.WasmTimeModule
	for i := 0; i < 10; i++ {
		wasmBytes := getWasmBytes(t, fmt.Sprintf("Hello, WebAssembly! %d", i), wasmDir)
		module, err := getModule(factory, wasmBytes)
		require.NoError(t, err)
		modules = append(modules, module)
	}

	assert.Equal(t, 10, stats.GetMisses())
	assert.Equal(t, 10, stats.GetAdditions())

	// Now retrieve the module from the cache and confirm that the cache stats are updated correctly
	// and that the retrieved module is correct
	for i := 0; i < 10; i++ {
		wasmBytes := getWasmBytes(t, fmt.Sprintf("Hello, WebAssembly! %d", i), wasmDir)
		module := modules[i]

		retrievedModule, err := getModule(factory, wasmBytes)
		require.NoError(t, err)

		equal, err := host.ModuleEquals(module, retrievedModule)
		require.NoError(t, err)
		assert.True(t, equal)
	}

	assert.Equal(t, 10, stats.GetHits())
	assert.Equal(t, 10, stats.GetMisses())
}

func TestNewModuleMultithreaded_RetrievalUsingBinary(t *testing.T) {
	lggr := logger.Test(t)

	wasmDir := t.TempDir()
	cacheDir := t.TempDir()
	stats := &testCacheStats{}
	factory, err := NewFileBasedModuleCache(lggr, cacheDir, stats)
	require.NoError(t, err)

	// generate 10 unique wasmBytes and get and store the module as a pair
	var modules []*host.WasmTimeModule
	for i := 0; i < 10; i++ {
		wasmBytes := getWasmBytes(t, fmt.Sprintf("Hello, WebAssembly! %d", i), wasmDir)
		module, err := getModule(factory, wasmBytes)
		require.NoError(t, err)
		modules = append(modules, module)
	}

	assert.Equal(t, 10, stats.GetMisses())
	assert.Equal(t, 10, stats.GetAdditions())

	// Now retrieve the module from the cache and confirm that the cache stats are updated correctly
	// and that the retrieved module is correct
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			wasmBytes := getWasmBytes(t, fmt.Sprintf("Hello, WebAssembly! %d", i), wasmDir)
			module := modules[i]

			retrievedModule, err := getModule(factory, wasmBytes)
			require.NoError(t, err)

			equal, err := host.ModuleEquals(module, retrievedModule)
			require.NoError(t, err)
			assert.True(t, equal)
		}(i)
	}
	wg.Wait()

	assert.Equal(t, 10, stats.GetHits())
	assert.Equal(t, 10, stats.GetMisses())
}

func TestNewModuleMultithreaded_RetrievalUsingBinaryID(t *testing.T) {
	lggr := logger.Test(t)

	wasmDir := t.TempDir()
	cacheDir := t.TempDir()
	stats := &testCacheStats{}
	factory, err := NewFileBasedModuleCache(lggr, cacheDir, stats)
	require.NoError(t, err)

	// generate 10 unique wasmBytes and get and store the module as a pair
	var modules []*host.WasmTimeModule
	for i := 0; i < 10; i++ {
		wasmBytes := getWasmBytes(t, fmt.Sprintf("Hello, WebAssembly! %d", i), wasmDir)
		module, err := getModule(factory, wasmBytes)
		require.NoError(t, err)
		modules = append(modules, module)
	}

	assert.Equal(t, 10, stats.GetMisses())
	assert.Equal(t, 10, stats.GetAdditions())

	// Now retrieve the module from the cache and confirm that the cache stats are updated correctly
	// and that the retrieved module is correct
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			wasmBytes := getWasmBytes(t, fmt.Sprintf("Hello, WebAssembly! %d", i), wasmDir)
			module := modules[i]

			binaryID := host.FromBinary(wasmBytes)
			retrievedModule, err := factory.GetModuleFromBinaryID(binaryID, 0)
			require.NoError(t, err)

			equal, err := host.ModuleEquals(module, retrievedModule)
			require.NoError(t, err)
			assert.True(t, equal)
		}(i)
	}
	wg.Wait()

	assert.Equal(t, 10, stats.GetHits())
	assert.Equal(t, 10, stats.GetMisses())
}

func TestNewModuleMultithreaded_Adding(t *testing.T) {
	lggr := logger.Test(t)

	wasmDir := t.TempDir()
	cacheDir := t.TempDir()
	stats := &testCacheStats{}
	factory, err := NewFileBasedModuleCache(lggr, cacheDir, stats)
	require.NoError(t, err)

	// generate 10 unique wasmBytes and get and store the module as a pair
	var modules []*host.WasmTimeModule
	for i := 0; i < 10; i++ {
		wasmBytes := getWasmBytes(t, fmt.Sprintf("Hello, WebAssembly! %d", i), wasmDir)
		module, err := getModule(factory, wasmBytes)
		require.NoError(t, err)
		modules = append(modules, module)
	}

	// Now using a new cache, add/retrieve the module from the cache in a multithreaded environment
	// and confirm that the cache stats are updated correctly and the retrieved module is correct
	cacheDir = t.TempDir()
	stats = &testCacheStats{}
	factory, err = NewFileBasedModuleCache(lggr, cacheDir, stats)
	require.NoError(t, err)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			wasmBytes := getWasmBytes(t, fmt.Sprintf("Hello, WebAssembly! %d", i), wasmDir)
			module := modules[i]

			retrievedModule, err := getModule(factory, wasmBytes)
			require.NoError(t, err)
			equal, err := host.ModuleEquals(module, retrievedModule)
			require.NoError(t, err)
			assert.True(t, equal)
		}(i)

		go func(i int) {
			defer wg.Done()
			wasmBytes := getWasmBytes(t, fmt.Sprintf("Hello, WebAssembly! %d", i), wasmDir)
			module := modules[i]

			retrievedModule, err := getModule(factory, wasmBytes)
			require.NoError(t, err)
			equal, err := host.ModuleEquals(module, retrievedModule)
			require.NoError(t, err)
			assert.True(t, equal)
		}(i)
	}

	wg.Wait()

	assert.Equal(t, 10, stats.GetHits())
	assert.Equal(t, 10, stats.GetMisses())
	assert.Equal(t, 10, stats.GetAdditions())
}

func getModule(factory *fileBasedModuleCache, wasmBytes []byte) (*host.WasmTimeModule, error) {
	maxCompressedBinarySize := uint64(math.MaxUint64)
	maxDecompressedBinarySize := uint64(math.MaxUint64)

	return factory.GetModuleFromBinary(wasmBytes, 0, true, maxCompressedBinarySize, maxDecompressedBinarySize)
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
