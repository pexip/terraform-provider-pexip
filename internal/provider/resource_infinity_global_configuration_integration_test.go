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

func TestInfinityGlobalConfigurationIntegration(t *testing.T) {
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

	testInfinityGlobalConfigurationIntegration(t, client)
}

func testInfinityGlobalConfigurationIntegration(t *testing.T, client InfinityClient) {
	// pexip_infinity_global_configuration is a singleton — its Delete operation sends a PATCH to reset
	// all fields to API defaults, clearing any references to related resources. Once those
	// references are cleared, the related resources can be deleted normally in the destroy step.
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Apply full configuration with all related resources.
			{
				Config: test.LoadTestFolder(t, "resource_infinity_global_configuration_full_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_global_configuration.global_configuration-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_global_configuration.global_configuration-test", "default_theme"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "logon_banner", "tf-test-logon-banner"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "site_banner", "tf-test-site-banner"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "site_banner_bg", "#123456"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "site_banner_fg", "#ffffff"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "contact_email_address", "tf-test@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_analytics", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_breakout_rooms", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_clock", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_legacy_dialout_api", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "guests_only_timeout", "120"),
				),
			},
			// Step 2: Destroy — triggers Delete which resets all fields to API defaults and
			// clears the default_theme reference, allowing the IVR theme to be deleted cleanly.
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_global_configuration_full_integration"),
				Destroy: true,
			},
		},
	})
}
