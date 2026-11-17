package ovh

import (
	"io"
	"log"
	"net/http"
	"time"
)

// retryDelays are the waits before each retry of a request the OVHcloud API answered with a server
// error or a rate limit. Kept short on purpose: a retry lengthens the scan, and steampipe scans can
// run in parallel workers that may be cut short past a fixed timeout. A variable so tests can shorten
// it.
var retryDelays = []time.Duration{time.Second, 2 * time.Second}

// retryTransport retries the GET and HEAD requests the OVHcloud API answers with a server error (5xx)
// or a rate limit (429), both of which it returns transiently, in bursts. Wrapping the HTTP transport
// covers every call of every table, and retries single requests, never a whole listing, so no row
// can be streamed twice.
type retryTransport struct {
	base   http.RoundTripper
	delays []time.Duration
}

func newRetryTransport(base http.RoundTripper) http.RoundTripper {
	return &retryTransport{base: base, delays: retryDelays}
}

func isRetryableStatus(code int) bool {
	return code >= http.StatusInternalServerError || code == http.StatusTooManyRequests
}

func (t *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.base.RoundTrip(req)
	// Only idempotent requests without a body are safe to send again as they are
	if req.Method != http.MethodGet && req.Method != http.MethodHead {
		return resp, err
	}
	for _, delay := range t.delays {
		if err != nil || !isRetryableStatus(resp.StatusCode) {
			return resp, err
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
		// The standard logger, not plugin.Logger: many calls carry no query context to take it from
		log.Printf("[WARN] ovh: %s %s answered %d, retrying in %s", req.Method, req.URL.Path, resp.StatusCode, delay)
		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		case <-time.After(delay):
		}
		resp, err = t.base.RoundTrip(req)
	}
	return resp, err
}
