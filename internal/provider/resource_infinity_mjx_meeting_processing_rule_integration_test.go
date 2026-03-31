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

func TestInfinityMjxMeetingProcessingRuleIntegration(t *testing.T) {
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

	testInfinityMjxMeetingProcessingRuleIntegration(t, client)
}

func testInfinityMjxMeetingProcessingRuleIntegration(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_meeting_processing_rule_full_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_meeting_processing_rule.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_meeting_processing_rule.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "name", "tf-test mjx-meeting-processing-rule full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "description", "Test MJX meeting processing rule"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "priority", "10"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "enabled", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "meeting_type", "teams"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_meeting_processing_rule.test", "mjx_integration"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "transform_rule", "{{ domain }}"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "domain", "example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "company_id", "test-company-id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "include_pin", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "default_processing_enabled", "false"),
				),
			},
			// Step 2: Update with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_meeting_processing_rule_min_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_meeting_processing_rule.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_meeting_processing_rule.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "name", "tf-test mjx-meeting-processing-rule min"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "priority", "1"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "enabled", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "meeting_type", "pexipinfinity"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_meeting_processing_rule.test", "mjx_integration"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "match_string", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "replace_string", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "domain", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "company_id", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "include_pin", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "default_processing_enabled", "true"),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_mjx_meeting_processing_rule_min_integration"),
				Destroy: true,
			},
			// Step 4: Recreate with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_meeting_processing_rule_min_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_meeting_processing_rule.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_meeting_processing_rule.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "name", "tf-test mjx-meeting-processing-rule min"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "priority", "1"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "meeting_type", "pexipinfinity"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_meeting_processing_rule.test", "mjx_integration"),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_meeting_processing_rule_full_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_meeting_processing_rule.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_meeting_processing_rule.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "name", "tf-test mjx-meeting-processing-rule full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "description", "Test MJX meeting processing rule"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "priority", "10"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "enabled", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "meeting_type", "teams"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_meeting_processing_rule.test", "mjx_integration"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "transform_rule", "{{ domain }}"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "domain", "example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "company_id", "test-company-id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "include_pin", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "default_processing_enabled", "false"),
				),
			},
		},
	})
}
