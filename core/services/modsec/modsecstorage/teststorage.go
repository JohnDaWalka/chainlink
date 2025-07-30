package modsecstorage

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
)

// ErrNotFound is returned when a key is not found in the test storage.
// This is similar to an HTTP 404 Not Found error.
var ErrNotFound = errors.New("not found")

type testStorage struct {
	storage map[string][]byte
}

// Get implements Storage.
func (t *testStorage) Get(ctx context.Context, key string) ([]byte, error) {
	obj, ok := t.storage[key]
	if !ok {
		return nil, ErrNotFound
	}
	return obj, nil
}

// Set implements Storage.
func (t *testStorage) Set(ctx context.Context, key string, value []byte) error {
	t.storage[key] = value
	return nil
}

// NewTestStorage creates a new Storage implementation, backed by a map.
// This is intended to be used for testing only.
func NewTestStorage() Storage {
	return &testStorage{
		storage: make(map[string][]byte),
	}
}

// NewTestServer creates a new httptest.Server backed by a test storage.
// This is intended to be used for testing only.
func NewTestServer() (*httptest.Server, func()) {
	testStore := NewTestStorage()

	handler := http.NewServeMux()
	handler.HandleFunc("/get/", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path[len("/get/"):]
		val, err := testStore.Get(r.Context(), key)
		if errors.Is(err, ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, "get error", http.StatusInternalServerError)
			return
		}
		if val == nil {
			http.NotFound(w, r)
			return
		}
		w.Write(val)
	})
	handler.HandleFunc("/set/", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path[len("/set/"):]
		val, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "read error", http.StatusInternalServerError)
			return
		}
		err = testStore.Set(r.Context(), key, val)
		if err != nil {
			http.Error(w, "set error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(handler)
	return server, server.Close
}
