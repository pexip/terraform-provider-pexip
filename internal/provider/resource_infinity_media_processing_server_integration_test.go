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

func TestInfinityMediaProcessingServerIntegration(t *testing.T) {
	_ = os.Setenv("TF_ACC", "1")

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

	testInfinityMediaProcessingServerIntegration(t, client)
}

func testInfinityMediaProcessingServerIntegration(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_media_processing_server_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_processing_server.media_processing_server-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_processing_server.media_processing_server-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_processing_server.media_processing_server-test", "fqdn", "tf-test-mps-full.test.local"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_processing_server.media_processing_server-test", "app_id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_processing_server.media_processing_server-test", "public_jwt_key"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_media_processing_server_full_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_processing_server.media_processing_server-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_processing_server.media_processing_server-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_processing_server.media_processing_server-test", "fqdn", "tf-test.updated.test.local"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_processing_server.media_processing_server-test", "app_id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_processing_server.media_processing_server-test", "public_jwt_key"),
				),
			},
		},
	})
}
