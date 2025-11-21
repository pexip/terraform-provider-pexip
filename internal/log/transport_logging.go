package log

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LoggingTransport struct {
	Base http.RoundTripper
}

func (t *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	// Read and log request body if it exists
	var reqBody []byte
	if req.Body != nil {
		reqBody, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(reqBody)) // reattach
	}
	fmt.Printf("\n--> %s %s\n%s\n", req.Method, req.URL, string(reqBody))

	// Send the request
	resp, err := t.Base.RoundTrip(req)
	if err != nil {
		fmt.Printf("<-- ERROR: %v (%v)\n", err, time.Since(start))
		return nil, err
	}

	// Read and log response body (non-destructive)
	var respBody []byte
	if resp.Body != nil {
		respBody, _ = io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewReader(respBody)) // reattach
	}
	fmt.Printf("<-- %s %s (%v)\n%s\n", resp.Status, req.URL, time.Since(start), string(respBody))

	return resp, nil
}
