package log

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// DefaultMaxBodyLength is the default maximum number of bytes to log from request/response bodies
	DefaultMaxBodyLength = 1000
)

type LoggingTransport struct {
	Base          http.RoundTripper
	MaxBodyLength int // Maximum number of bytes to log from request/response bodies. 0 = use default (1000), negative = unlimited.
}

// truncateBody truncates the body to the maximum length and adds a truncation message if needed
func (t *LoggingTransport) truncateBody(body []byte) string {
	// Determine the max length to use
	maxLen := t.MaxBodyLength
	if maxLen == 0 {
		// Not set, use default
		maxLen = DefaultMaxBodyLength
	} else if maxLen < 0 {
		// Negative value means unlimited
		return string(body)
	}

	// Check if truncation is needed
	if len(body) <= maxLen {
		return string(body)
	}

	return fmt.Sprintf("%s\n... [truncated %d bytes, showing first %d bytes]", string(body[:maxLen]), len(body)-maxLen, maxLen)
}

func (t *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	// Read and log request body if it exists
	var reqBody []byte
	if req.Body != nil {
		reqBody, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(reqBody)) // reattach
	}
	fmt.Printf("\n--> %s %s\n%s\n", req.Method, req.URL, t.truncateBody(reqBody))

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
	fmt.Printf("<-- %s %s (%v)\n%s\n", resp.Status, req.URL, time.Since(start), t.truncateBody(respBody))

	return resp, nil
}
