package log

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggingTransport_DefaultMaxBodyLength(t *testing.T) {
	// Create a test server that returns a large response
	largeBody := strings.Repeat("A", 2000)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(largeBody))
	}))
	defer server.Close()

	// Create transport with default max body length
	transport := &LoggingTransport{
		Base: http.DefaultTransport,
	}

	client := &http.Client{Transport: transport}
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Verify we got the full response
	body, _ := io.ReadAll(resp.Body)
	if len(body) != 2000 {
		t.Errorf("Expected body length 2000, got %d", len(body))
	}
}

func TestLoggingTransport_CustomMaxBodyLength(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("short response"))
	}))
	defer server.Close()

	// Create transport with custom max body length
	transport := &LoggingTransport{
		Base:          http.DefaultTransport,
		MaxBodyLength: 500,
	}

	client := &http.Client{Transport: transport}
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()
}

func TestLoggingTransport_UnlimitedBodyLength(t *testing.T) {
	// Create a test server
	largeBody := strings.Repeat("B", 5000)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(largeBody))
	}))
	defer server.Close()

	// Create transport with unlimited body length (negative value)
	transport := &LoggingTransport{
		Base:          http.DefaultTransport,
		MaxBodyLength: -1,
	}

	client := &http.Client{Transport: transport}
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Verify we got the full response
	body, _ := io.ReadAll(resp.Body)
	if len(body) != 5000 {
		t.Errorf("Expected body length 5000, got %d", len(body))
	}
}

func TestLoggingTransport_TruncateBody(t *testing.T) {
	tests := []struct {
		name          string
		maxBodyLength int
		bodySize      int
		expectTrunc   bool
	}{
		{"Default with small body", 0, 500, false},
		{"Default with large body", 0, 2000, true},
		{"Custom 200 with small body", 200, 100, false},
		{"Custom 200 with large body", 200, 500, true},
		{"Unlimited with large body", -1, 5000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport := &LoggingTransport{
				Base:          http.DefaultTransport,
				MaxBodyLength: tt.maxBodyLength,
			}

			body := bytes.Repeat([]byte("X"), tt.bodySize)
			truncated := transport.truncateBody(body)

			if tt.expectTrunc {
				if !strings.Contains(truncated, "truncated") {
					t.Errorf("Expected truncation message, but got: %s", truncated[:50])
				}
			} else {
				if strings.Contains(truncated, "truncated") {
					t.Errorf("Did not expect truncation, but got: %s", truncated[:50])
				}
			}
		})
	}
}
