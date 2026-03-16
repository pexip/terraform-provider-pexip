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

func TestInfinityConferenceIntegration(t *testing.T) {
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

	testInfinityConferenceIntegration(t, client)
}

func testInfinityConferenceIntegration(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			// Test 1: Create with full configuration
			{
				Config: test.LoadTestFolder(t, "resource_infinity_conference_full_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_conference.tf-test-conference", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_conference.tf-test-conference", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "name", "tf-test-conference"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "breakout_rooms", "true"),
				),
			},
			// Test 2: Update to min configuration
			{
				Config: test.LoadTestFolder(t, "resource_infinity_conference_min_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_conference.tf-test-conference", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_conference.tf-test-conference", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "name", "tf-test-conference"),
				),
			},
			// Test 3: Create with min configuration (after destroy)
			{
				Config: test.LoadTestFolder(t, "resource_infinity_conference_min_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_conference.tf-test-conference", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_conference.tf-test-conference", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "name", "tf-test-conference"),
				),
			},
			// Test 4: Update to full configuration
			{
				Config: test.LoadTestFolder(t, "resource_infinity_conference_full_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_conference.tf-test-conference", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_conference.tf-test-conference", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "name", "tf-test-conference"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "breakout_rooms", "true"),
				),
			},
		},
	})
}
