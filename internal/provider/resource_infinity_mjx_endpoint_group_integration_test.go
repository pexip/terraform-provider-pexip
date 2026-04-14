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

func TestInfinityMjxEndpointGroupIntegration(t *testing.T) {
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

	testInfinityMjxEndpointGroupIntegration(t, client)
}

func testInfinityMjxEndpointGroupIntegration(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_endpoint_group_full_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint_group.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint_group.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "name", "tf-test mjx-endpoint-integration-group full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "description", "Test MJX endpoint integration group"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint_group.test", "system_location"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint_group.test", "mjx_integration"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "disable_proxy", "true"),
				),
			},
			// Step 2: Update with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_endpoint_group_min_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint_group.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint_group.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "name", "tf-test mjx-endpoint-integration-group min"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "description", ""),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint_group.test", "system_location"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "disable_proxy", "false"),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_mjx_endpoint_group_min_integration"),
				Destroy: true,
			},
			// Step 4: Recreate with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_endpoint_group_min_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint_group.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint_group.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "name", "tf-test mjx-endpoint-integration-group min"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint_group.test", "system_location"),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_endpoint_group_full_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint_group.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint_group.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "name", "tf-test mjx-endpoint-integration-group full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "description", "Test MJX endpoint integration group"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint_group.test", "system_location"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint_group.test", "mjx_integration"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "disable_proxy", "true"),
				),
			},
		},
	})
}
