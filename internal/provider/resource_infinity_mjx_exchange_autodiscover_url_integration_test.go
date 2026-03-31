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
	"github.com/stretchr/testify/require"

	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityMjxExchangeAutodiscoverURLIntegration(t *testing.T) {
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

	testInfinityMjxExchangeAutodiscoverURLIntegration(t, client)
}

func testInfinityMjxExchangeAutodiscoverURLIntegration(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_exchange_autodiscover_url_full_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_autodiscover_url.test_full", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_autodiscover_url.test_full", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test_full", "name", "tf-test mjx-exchange-autodiscover-url full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test_full", "description", "Test Exchange Autodiscover URL description"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test_full", "url", "https://autodiscover-full.example.com/autodiscover/autodiscover.xml"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_autodiscover_url.test_full", "exchange_deployment"),
				),
			},
			// Step 2: Update with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_exchange_autodiscover_url_min_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_autodiscover_url.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_autodiscover_url.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test", "name", "tf-test mjx-exchange-autodiscover-url min"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test", "url", "https://autodiscover.example.com/autodiscover/autodiscover.xml"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_autodiscover_url.test", "exchange_deployment"),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_mjx_exchange_autodiscover_url_min_integration"),
				Destroy: true,
			},
			// Step 4: Recreate with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_exchange_autodiscover_url_min_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_autodiscover_url.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_autodiscover_url.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test", "name", "tf-test mjx-exchange-autodiscover-url min"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test", "url", "https://autodiscover.example.com/autodiscover/autodiscover.xml"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_autodiscover_url.test", "exchange_deployment"),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_exchange_autodiscover_url_full_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_autodiscover_url.test_full", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_autodiscover_url.test_full", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test_full", "name", "tf-test mjx-exchange-autodiscover-url full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test_full", "description", "Test Exchange Autodiscover URL description"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test_full", "url", "https://autodiscover-full.example.com/autodiscover/autodiscover.xml"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_autodiscover_url.test_full", "exchange_deployment"),
				),
			},
		},
	})
}
