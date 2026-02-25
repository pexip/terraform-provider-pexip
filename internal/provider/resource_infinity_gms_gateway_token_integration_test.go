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

	"github.com/pexip/go-infinity-sdk/v38"
	"github.com/pexip/terraform-provider-pexip/internal/test"
	"github.com/stretchr/testify/require"
)

func TestInfinityGMSGatewayTokenIntegration(t *testing.T) {
	_ = os.Setenv("TF_ACC", "1")

	// Verify required environment variables are set
	certificate := os.Getenv("TF_VAR_infinity_gms_gw_token_cert")
	privateKey := os.Getenv("TF_VAR_infinity_gms_gw_token_key")
	certificateUpdated := os.Getenv("TF_VAR_infinity_gms_gw_token_cert2")
	privateKeyUpdated := os.Getenv("TF_VAR_infinity_gms_gw_token_key2")

	require.NotEmpty(t, certificate, "TF_VAR_infinity_gms_gw_token_cert environment variable must be set")
	require.NotEmpty(t, privateKey, "TF_VAR_infinity_gms_gw_token_key environment variable must be set")
	require.NotEmpty(t, certificateUpdated, "TF_VAR_infinity_gms_gw_token_cert2 environment variable must be set")
	require.NotEmpty(t, privateKeyUpdated, "TF_VAR_infinity_gms_gw_token_key2 environment variable must be set")

	client, err := infinity.New(
		infinity.WithBaseURL(test.INFINITY_BASE_URL),
		infinity.WithBasicAuth(test.INFINITY_USERNAME, test.INFINITY_PASSWORD),
		infinity.WithMaxRetries(2),
		infinity.WithTransport(&http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // We need this because default certificate is not trusted
				MinVersion:         tls.VersionTLS12,
			},
			MaxIdleConns:        30,
			MaxIdleConnsPerHost: 5,
			IdleConnTimeout:     60 * time.Second,
		}),
	)
	require.NoError(t, err)

	testInfinityGMSGatewayToken(t, client)
}
