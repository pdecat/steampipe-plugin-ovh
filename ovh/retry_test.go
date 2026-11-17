package ovh

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// statusServer answers each request with the next status of codes, then 200 once they are used up.
func statusServer(t *testing.T, codes ...int) (*httptest.Server, *int) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		if calls <= len(codes) {
			w.WriteHeader(codes[calls-1])
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)
	return server, &calls
}

func testRetryClient() *http.Client {
	return &http.Client{Transport: &retryTransport{base: http.DefaultTransport, delays: []time.Duration{0, 0}}}
}

func TestRetryTransportRecovers(t *testing.T) {
	server, calls := statusServer(t, http.StatusInternalServerError, http.StatusTooManyRequests)
	resp, err := testRetryClient().Get(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK || *calls != 3 {
		t.Fatalf("got status %d after %d calls, want 200 after 3", resp.StatusCode, *calls)
	}
}

func TestRetryTransportGivesUp(t *testing.T) {
	server, calls := statusServer(t, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout, http.StatusInternalServerError)
	resp, err := testRetryClient().Get(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusGatewayTimeout || *calls != 3 {
		t.Fatalf("got status %d after %d calls, want the third attempt's 504 after 3", resp.StatusCode, *calls)
	}
}

func TestRetryTransportDoesNotRetryClientErrors(t *testing.T) {
	for _, code := range []int{http.StatusBadRequest, http.StatusForbidden, http.StatusNotFound} {
		server, calls := statusServer(t, code)
		resp, err := testRetryClient().Get(server.URL)
		if err != nil {
			t.Fatalf("%d: unexpected error: %v", code, err)
		}
		_ = resp.Body.Close()
		if resp.StatusCode != code || *calls != 1 {
			t.Fatalf("%d: got status %d after %d calls, want it back after 1", code, resp.StatusCode, *calls)
		}
	}
}

func TestRetryTransportDoesNotRetryNonIdempotentRequests(t *testing.T) {
	server, calls := statusServer(t, http.StatusInternalServerError)
	resp, err := testRetryClient().Post(server.URL, "application/json", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusInternalServerError || *calls != 1 {
		t.Fatalf("got status %d after %d calls, want the 500 back after 1", resp.StatusCode, *calls)
	}
}

func TestRetryTransportStopsOnCancel(t *testing.T) {
	server, calls := statusServer(t, http.StatusInternalServerError, http.StatusInternalServerError)
	client := &http.Client{Transport: &retryTransport{base: http.DefaultTransport, delays: []time.Duration{time.Hour}}}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
	_, err := client.Do(req)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("got %v, want the context deadline", err)
	}
	if *calls != 1 {
		t.Fatalf("got %d calls, want 1", *calls)
	}
}
