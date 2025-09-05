package pipeline

import (
	"context"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

type HTTPRequest struct {
	Method      StringParam
	URL         URLParam
	ReqHeaders  []string
	RequestData MapParam
}

type HTTPResponse struct {
	ID         string
	StatusCode int
	RespHeaders  http.Header
	Body       []byte
	Error      error
}

type PendingRequest struct {
	RequestId  string
	Request    HTTPRequest
	ResponseCh chan HTTPResponse
	Timeout    time.Time
}

type BatchMiddleware struct {
	mu           sync.RWMutex
	batches      map[string]*RequestBatch // keyed by API endpoint
	batchSize    int
	batchTimeout time.Duration
	logger       logger.Logger
}

type RequestBatch struct {
	mu       sync.Mutex
	requests []PendingRequest
	timer    *time.Timer
	endpoint string
}

func NewBatchMiddleware(batchSize int, batchTimeout time.Duration, lggr logger.Logger) *BatchMiddleware {
	return &BatchMiddleware{
		batches:      make(map[string]*RequestBatch),
		batchSize:    batchSize,
		batchTimeout: batchTimeout,
		logger:       lggr,
	}
}

func (bm *BatchMiddleware) AddRequest(ctx context.Context, req HTTPRequest) (HTTPResponse, error) {
	endpoint := req.URL.String() // Use request URL as batch key

	// Find or create batch for this endpoint
	bm.mu.Lock()
	batch, exists := bm.batches[endpoint]
	if !exists {
		batch = &RequestBatch{
			endpoint: endpoint,
			requests: make([]PendingRequest, 0, bm.batchSize),
		}
		bm.batches[endpoint] = batch
	}
	bm.mu.Unlock()

	// Add new pending request to batch
	responseCh := make(chan HTTPResponse, 1)
	pendingReq := PendingRequest{
		RequestId: generateRequestID(),
		Request:    req,
		ResponseCh: responseCh,
		Timeout:    time.Now().Add(bm.batchTimeout * 2), // Allow extra time for batching
	}

	batch.mu.Lock()
	batch.requests = append(batch.requests, pendingReq)
	requestCount := len(batch.requests)

	// Set timer for first request in batch
	if requestCount == 1 {
		batch.timer = time.AfterFunc(bm.batchTimeout, func() {
			bm.executeBatch(endpoint)
		})
	}

	// Execute immediately if batch is full
	if requestCount >= bm.batchSize {
		if batch.timer != nil {
			batch.timer.Stop()
		}
		batch.mu.Unlock()
		go bm.executeBatch(endpoint)
	} else {
		batch.mu.Unlock()
	}

	// Wait for response
	select {
	case response := <-responseCh:
		return response, response.Error
	case <-ctx.Done():
		return HTTPResponse{}, ctx.Err()
	}
}

func (bm *BatchMiddleware) executeBatch(endpoint string) {
	bm.mu.Lock()
	batch, exists := bm.batches[endpoint]
	if !exists {
		bm.mu.Unlock()
		return
	}
	delete(bm.batches, endpoint)
	bm.mu.Unlock()

	batch.mu.Lock()
	requests := batch.requests
	batch.mu.Unlock()

	if len(requests) == 0 {
		return
	}

	// TODO: Convert individual requests to batch payload format
	// for i, req := range requests {
	//     // do something
	// }

	// responseBytes, statusCode, respHeaders, start, finish, err := makeHTTPRequest(req)
	// if err != nil {
	//     // handle err
	// }

	// TODO: Parse http response into separate responses
	responses := make([]HTTPResponse, len(requests))

	// Send responses back to waiting goroutines
	// TODO: match by request ID
	for i, req := range requests {
		if i < len(responses) {
			select {
			case req.ResponseCh <- responses[i]:
			default:
				bm.logger.Warnw("Failed to send response to channel", "requestID", req.RequestId)
			}
		}
	}
}

func generateRequestID() string {
    // todo: simple uuid - can revisit this if we want something else
    return uuid.New().String()
}
