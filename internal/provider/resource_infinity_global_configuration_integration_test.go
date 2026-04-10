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
		ExternalProviders: map[string]resource.ExternalProvider{
			"tls": {
				Source: "hashicorp/tls",
			},
		},
		Steps: []resource.TestStep{
			// Step 1: Apply full configuration with all related resources.
			{
				Config: test.LoadTestFolder(t, "resource_infinity_global_configuration_full_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_global_configuration.global_configuration-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_global_configuration.global_configuration-test", "default_theme"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "bdpm_max_pin_failures_per_window", "99"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "bdpm_max_scan_attempts_per_window", "88"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "bdpm_pin_checks_enabled", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "bdpm_scan_quarantine_enabled", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "bursting_enabled", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "bursting_min_lifetime", "123"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "bursting_threshold", "77"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "cloud_provider", "GCP"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "contact_email_address", "tf-test@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "content_security_policy_state", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "crypto_mode", "on"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "disabled_codecs.#", "2"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_global_configuration.global_configuration-test", "disabled_codecs.*", "H264_H_1"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_global_configuration.global_configuration-test", "disabled_codecs.*", "VP9"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "eject_last_participant_backstop_timeout", "60"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_analytics", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_application_api", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_breakout_rooms", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_chat", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_clock", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_denoise", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_dialout", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_directory", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_edge_non_mesh", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_fecc", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_h323", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_legacy_dialout_api", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_lync_auto_escalate", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_lync_vbss", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_mlvad", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_rtmp", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_sip", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_sip_udp", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_softmute", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_ssh", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_turn_443", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_webrtc", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "error_reporting_enabled", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "error_reporting_url", "https://custom-error-reporting.com"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "es_connection_timeout", "99"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "es_initial_retry_backoff", "22"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "es_maximum_retry_backoff", "444"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "es_media_streams_wait", "55"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "es_metrics_update_interval", "66"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "es_short_term_memory_expiration", "77"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "external_participant_avatar_lookup", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "gcp_client_email", "tf-test@gcp.com"),
					resource.TestCheckResourceAttrSet("pexip_infinity_global_configuration.global_configuration-test", "gcp_private_key"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "gcp_project_id", "tf-test-project"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "guests_only_timeout", "999"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "legacy_api_username", "tf-test-user"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "live_captions_vmr_default", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "liveview_show_conferences", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "local_mssip_domain", "tf-test.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "logon_banner", "tf-test-logon-banner"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "logs_max_age", "123"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "management_qos", "46"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "management_session_timeout", "99"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "management_start_page", "/admin/conferencingstatus/conference/"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "max_callrate_in", "888"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "max_callrate_out", "999"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "max_pixels_per_second", "fullhd"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "max_presentation_bandwidth_ratio", "33"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "media_ports_end", "49999"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "media_ports_start", "40123"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "ocsp_responder_url", "https://ocsp.tf-test.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "ocsp_state", "ON"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "pin_entry_timeout", "321"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "session_timeout_enabled", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "signalling_ports_end", "39998"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "signalling_ports_start", "33001"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "sip_tls_cert_verify_mode", "ON"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "site_banner", "tf-test-site-banner"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "site_banner_bg", "#123456"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "site_banner_fg", "#ffffff"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "teams_enable_powerpoint_render", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "waiting_for_chair_timeout", "901"),
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
