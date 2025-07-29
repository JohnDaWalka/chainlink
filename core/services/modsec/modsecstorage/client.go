package modsecstorage

import (
	"bytes"
	"context"
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

func NewStdClient(endpoint string) Storage {
	return &stdClient{
		endpoint: endpoint,
		client:   &http.Client{},
	}
}
