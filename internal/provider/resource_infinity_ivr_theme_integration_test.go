//go:build integration

/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"crypto/tls"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/pexip/terraform-provider-pexip/internal/log"
	"github.com/stretchr/testify/require"

	"github.com/pexip/terraform-provider-pexip/internal/test"

	"github.com/pexip/go-infinity-sdk/v38"
)

func TestInfinityIvrThemeIntegration(t *testing.T) {
	_ = os.Setenv("TF_ACC", "1")

	// Base transport configuration
	baseTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // We need this because default certificate is not trusted
			MinVersion:         tls.VersionTLS12,
		},
		MaxIdleConns:        30,
		MaxIdleConnsPerHost: 5,
		IdleConnTimeout:     60 * time.Second,
	}

	// Conditionally wrap with logging transport if DEBUG_HTTP is set
	var transport http.RoundTripper = baseTransport
	if os.Getenv("DEBUG_HTTP") != "" {
		transport = &log.LoggingTransport{
			Base: baseTransport,
		}
	}

	client, err := infinity.New(
		infinity.WithBaseURL(test.INFINITY_BASE_URL),
		infinity.WithBasicAuth(test.INFINITY_USERNAME, test.INFINITY_PASSWORD),
		infinity.WithMaxRetries(2),
		infinity.WithTransport(transport),
	)
	require.NoError(t, err)

	testInfinityIvrTheme(t, client)
}
