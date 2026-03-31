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

func TestInfinityMjxIntegrationIntegration(t *testing.T) {
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

	testInfinityMjxIntegrationIntegration(t, client)
}

func testInfinityMjxIntegrationIntegration(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_integration_full_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_integration.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_integration.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "name", "tf-test mjx-integration full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "description", "Test MJX integration description"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "display_upcoming_meetings", "14"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "enable_non_video_meetings", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "enable_private_meetings", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "end_buffer", "10"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "start_buffer", "10"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "ep_username", "ep-user@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "ep_use_https", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "ep_verify_certificate", "true"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_integration.test", "graph_deployment"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "process_alias_private_meetings", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "replace_empty_subject", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "replace_subject_type", "ALL"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "replace_subject_template", "Meeting: {{ subject }}"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "use_webex", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "webex_api_domain", "custom.webexapis.com"),
				),
			},
			// Step 2: Update with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_integration_min_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_integration.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_integration.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "name", "tf-test mjx-integration min"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "display_upcoming_meetings", "7"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "enable_non_video_meetings", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "enable_private_meetings", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "end_buffer", "0"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "start_buffer", "5"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "ep_username", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "ep_use_https", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "ep_verify_certificate", "false"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_integration.test", "graph_deployment"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "process_alias_private_meetings", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "replace_empty_subject", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "replace_subject_type", "PRIVATE"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "replace_subject_template", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "use_webex", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "webex_api_domain", "webexapis.com"),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_mjx_integration_min_integration"),
				Destroy: true,
			},
			// Step 4: Recreate with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_integration_min_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_integration.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_integration.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "name", "tf-test mjx-integration min"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "display_upcoming_meetings", "7"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_integration.test", "graph_deployment"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "use_webex", "false"),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_integration_full_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_integration.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_integration.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "name", "tf-test mjx-integration full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "description", "Test MJX integration description"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "display_upcoming_meetings", "14"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "enable_non_video_meetings", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "enable_private_meetings", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "end_buffer", "10"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "start_buffer", "10"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "ep_username", "ep-user@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "ep_use_https", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "ep_verify_certificate", "true"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_integration.test", "graph_deployment"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "process_alias_private_meetings", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "replace_empty_subject", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "replace_subject_type", "ALL"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "replace_subject_template", "Meeting: {{ subject }}"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "use_webex", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "webex_api_domain", "custom.webexapis.com"),
				),
			},
		},
	})
}
