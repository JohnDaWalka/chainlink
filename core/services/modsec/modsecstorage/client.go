package modsecstorage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type stdClient struct {
	endpoint string
	client   *http.Client
}

// Get implements Storage.
func (s *stdClient) Get(ctx context.Context, key string) ([]byte, error) {
	// the standard endpoint is GET /get/<key> where <key> is a message ID.
	url := fmt.Sprintf("%s/get/%s", s.endpoint, key)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get key %s: %s", key, resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// Set implements Storage.
func (s *stdClient) Set(ctx context.Context, key string, value []byte) error {
	// the standard endpoint is POST /set/<key> where <key> is a message ID.
	url := fmt.Sprintf("%s/set/%s", s.endpoint, key)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(value))
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to set key %s: %s", key, resp.Status)
	}

	return nil
}

// GetMany implements Storage.
func (s *stdClient) GetMany(ctx context.Context, keys []string) (map[string][]byte, error) {
	// the standard endpoint is POST /getmany with a JSON body of keys.
	url := fmt.Sprintf("%s/getmany", s.endpoint)
	body, err := json.Marshal(keys)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get keys: %s", resp.Status)
	}

	var results map[string][]byte
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetAll implements Storage.
func (s *stdClient) GetAll(ctx context.Context) (map[string][]byte, error) {
	// the standard endpoint is GET /getall
	url := fmt.Sprintf("%s/getall", s.endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get all keys: %s", resp.Status)
	}

	var results map[string][]byte
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}
	return results, nil
}

func NewStdClient(endpoint string) Storage {
	return &stdClient{
		endpoint: endpoint,
		client:   &http.Client{},
	}
}
