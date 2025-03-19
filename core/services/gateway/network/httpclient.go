package network

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// HTTPClient interfaces defines a method to send HTTP requests
type HTTPClient interface {
	Send(ctx context.Context, req HTTPRequest) (*HTTPResponse, error)
}

type HTTPClientConfig struct {
	MaxResponseBytes uint32
	DefaultTimeout   time.Duration
}

var (
	defaultAllowedPorts     = []int{80, 443}
	defaultAllowedSchemes   = []string{"http", "https"}
	defaultMaxResponseBytes = uint32(26.4 * utils.KB)
	defaultTimeout          = 5 * time.Second
)

type HTTPRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
	Timeout time.Duration

	// Maximum number of bytes to read from the response body.  If 0, the default value is used.
	// Does not override a request specific value gte 0.
	MaxResponseBytes uint32
}

type HTTPResponse struct {
	StatusCode int               // HTTP status code
	Headers    map[string]string // HTTP headers
	Body       []byte            // HTTP response body
}

type httpClient struct {
	client *http.Client
	config HTTPClientConfig
	lggr   logger.Logger
}

// NewHTTPClient creates a new NewHTTPClient
// As of now, the client does not support TLS configuration but may be extended in the future
func NewHTTPClient(config HTTPClientConfig, lggr logger.Logger) (HTTPClient, error) {
	return &httpClient{
		config: config,
		client: &http.Client{
			Timeout:   config.DefaultTimeout,
			Transport: http.DefaultTransport,
		},
		lggr: lggr,
	}, nil
}

func (c *httpClient) Send(ctx context.Context, req HTTPRequest) (*HTTPResponse, error) {
	to := req.Timeout
	if to == 0 {
		to = c.config.DefaultTimeout
	}

	c.lggr.Debugw("sending HTTP request with timeout", "url", req.URL, "request timeout", to)
	timeoutCtx, cancel := context.WithTimeout(ctx, to)
	defer cancel()
	r, err := http.NewRequestWithContext(timeoutCtx, req.Method, req.URL, bytes.NewBuffer(req.Body))
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		c.lggr.Errorw("failed to send HTTP request", "url", req.URL, "err", err)
		return nil, err
	}
	defer resp.Body.Close()

	n := maxReadBytes(readSize{defaultSize: c.config.MaxResponseBytes, requestSize: req.MaxResponseBytes})
	c.lggr.Debugw("max bytes to read from HTTP response", "bytes", n)

	reader := http.MaxBytesReader(nil, resp.Body, int64(n))
	body, err := io.ReadAll(reader)
	if err != nil {
		c.lggr.Errorw("failed to read HTTP response body", "url", req.URL, "err", err)
		return nil, err
	}
	headers := make(map[string]string)
	for k, v := range resp.Header {
		// header values are usually an array of size 1
		// joining them to a single string in case array size is greater than 1
		headers[k] = strings.Join(v, ",")
	}
	c.lggr.Debugw("received HTTP response", "statusCode", resp.StatusCode, "url", req.URL, "headers", headers)

	return &HTTPResponse{
		Headers:    headers,
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}

type readSize struct {
	defaultSize uint32
	requestSize uint32
}

func maxReadBytes(sizes readSize) uint32 {
	if sizes.requestSize == 0 {
		return sizes.defaultSize
	}
	return minUint32(sizes.defaultSize, sizes.requestSize)
}

func minUint32(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}
