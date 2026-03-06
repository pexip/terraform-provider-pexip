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

func TestInfinityOAuth2ClientIntegration(t *testing.T) {
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

	testInfinityOAuth2ClientIntegration(t, client)
}

func testInfinityOAuth2ClientIntegration(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_oauth2_client_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_oauth2_client.oauth2_client-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_oauth2_client.oauth2_client-test", "client_id"),
					resource.TestCheckResourceAttr("pexip_infinity_oauth2_client.oauth2_client-test", "client_name", "tf-test oauth2_client RW"),
					resource.TestCheckResourceAttr("pexip_infinity_oauth2_client.oauth2_client-test", "role", "/api/admin/configuration/v1/role/1/"),
					resource.TestCheckResourceAttrSet("pexip_infinity_oauth2_client.oauth2_client-test", "private_key_jwt"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_oauth2_client_full_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_oauth2_client.oauth2_client-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_oauth2_client.oauth2_client-test", "client_id"),
					resource.TestCheckResourceAttr("pexip_infinity_oauth2_client.oauth2_client-test", "client_name", "tf-test oauth2_client RO"),
					resource.TestCheckResourceAttr("pexip_infinity_oauth2_client.oauth2_client-test", "role", "/api/admin/configuration/v1/role/2/"),
					resource.TestCheckResourceAttrSet("pexip_infinity_oauth2_client.oauth2_client-test", "private_key_jwt"),
				),
			},
		},
	})
}
