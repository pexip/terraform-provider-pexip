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

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"
	"github.com/pexip/terraform-provider-pexip/internal/test"
	"github.com/stretchr/testify/require"
)

func TestInfinityGMSGatewayTokenIntegration(t *testing.T) {
	t.Skip("Skipping: Requires Pexip signed GMS certificate and private key")
	_ = os.Setenv("TF_ACC", "1")

	// Verify required environment variables are set
	certificate := os.Getenv("TF_VAR_gms_certificate")
	privateKey := os.Getenv("TF_VAR_gms_private_key")

	require.NotEmpty(t, certificate, "TF_VAR_gms_certificate environment variable must be set")
	require.NotEmpty(t, privateKey, "TF_VAR_gms_private_key environment variable must be set")

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

	testInfinityGMSGatewayTokenIntegration(t, client)
}

func testInfinityGMSGatewayTokenIntegration(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_gms_gateway_token_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_gms_gateway_token.gms_gateway_token_integration_test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_gms_gateway_token.gms_gateway_token_integration_test", "certificate"),
					resource.TestCheckResourceAttrSet("pexip_infinity_gms_gateway_token.gms_gateway_token_integration_test", "private_key"),
					// API should return extracted certificates
					resource.TestCheckResourceAttrSet("pexip_infinity_gms_gateway_token.gms_gateway_token_integration_test", "intermediate_certificate"),
					resource.TestCheckResourceAttrSet("pexip_infinity_gms_gateway_token.gms_gateway_token_integration_test", "leaf_certificate"),
					resource.TestCheckResourceAttrSet("pexip_infinity_gms_gateway_token.gms_gateway_token_integration_test", "supports_direct_guest_join"),
				),
			},
		},
	})
}
