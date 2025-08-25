/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"testing"

	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityGlobalConfiguration(t *testing.T) {
	t.Parallel()

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking - include all fields that are set in config
	mockState := &config.GlobalConfiguration{
		ID:                          1,
		ResourceURI:                 "/api/admin/configuration/v1/global/1/",
		EnableWebRTC:                true,
		EnableSIP:                   true,
		EnableH323:                  true,
		EnableRTMP:                  true,
		CryptoMode:                  "besteffort",
		MaxPixelsPerSecond:          "720000",
		BurstingEnabled:             true,
		CloudProvider:               "aws",
		AWSAccessKey:                test.StringPtr("test-key"),
		AWSSecretKey:                test.StringPtr("test-secret"),
		AzureClientID:               test.StringPtr("test-client"),
		AzureSecret:                 test.StringPtr("test-secret"),
		ConferenceCreatePermissions: "user_admin",
		ConferenceCreationMode:      "per_cluster",
		EnableAnalytics:             true,
		EnableErrorReporting:        true,
		BandwidthRestrictions:       "restricted",
		AdministratorEmail:          "test@example.com",
		MediaPortsStart:             40000,
		MediaPortsEnd:               40100,
		SignallingPortsStart:        5060,
		SignallingPortsEnd:          5070,
		GuestsOnlyTimeout:           300,
		WaitingForChairTimeout:      600,
	}

	// Mock the GetGlobalconfiguration API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/global/1/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		global_configuration := args.Get(2).(*config.GlobalConfiguration)
		*global_configuration = *mockState
	}).Maybe()

	// Mock the UpdateGlobalconfiguration API call
	client.On("PutJSON", mock.Anything, "configuration/v1/global/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.GlobalConfigurationUpdateRequest)
		global_configuration := args.Get(3).(*config.GlobalConfiguration)

		// Update mock state based on request
		if updateRequest.EnableWebRTC != nil {
			mockState.EnableWebRTC = *updateRequest.EnableWebRTC
		}
		if updateRequest.EnableSIP != nil {
			mockState.EnableSIP = *updateRequest.EnableSIP
		}
		if updateRequest.EnableH323 != nil {
			mockState.EnableH323 = *updateRequest.EnableH323
		}
		if updateRequest.EnableRTMP != nil {
			mockState.EnableRTMP = *updateRequest.EnableRTMP
		}
		if updateRequest.CryptoMode != "" {
			mockState.CryptoMode = updateRequest.CryptoMode
		}
		if updateRequest.MaxPixelsPerSecond != "" {
			mockState.MaxPixelsPerSecond = updateRequest.MaxPixelsPerSecond
		}
		if updateRequest.BurstingEnabled != nil {
			mockState.BurstingEnabled = *updateRequest.BurstingEnabled
		}
		if updateRequest.CloudProvider != "" {
			mockState.CloudProvider = updateRequest.CloudProvider
		}
		if updateRequest.AWSAccessKey != nil {
			mockState.AWSAccessKey = updateRequest.AWSAccessKey
		}
		if updateRequest.AWSSecretKey != nil {
			mockState.AWSSecretKey = updateRequest.AWSSecretKey
		}
		if updateRequest.AzureClientID != nil {
			mockState.AzureClientID = updateRequest.AzureClientID
		}
		if updateRequest.AzureSecret != nil {
			mockState.AzureSecret = updateRequest.AzureSecret
		}
		if updateRequest.ConferenceCreatePermissions != "" {
			mockState.ConferenceCreatePermissions = updateRequest.ConferenceCreatePermissions
		}
		if updateRequest.ConferenceCreationMode != "" {
			mockState.ConferenceCreationMode = updateRequest.ConferenceCreationMode
		}
		if updateRequest.EnableAnalytics != nil {
			mockState.EnableAnalytics = *updateRequest.EnableAnalytics
		}
		if updateRequest.EnableErrorReporting != nil {
			mockState.EnableErrorReporting = *updateRequest.EnableErrorReporting
		}
		if updateRequest.BandwidthRestrictions != "" {
			mockState.BandwidthRestrictions = updateRequest.BandwidthRestrictions
		}
		if updateRequest.AdministratorEmail != "" {
			mockState.AdministratorEmail = updateRequest.AdministratorEmail
		}
		if updateRequest.MediaPortsStart != nil {
			mockState.MediaPortsStart = *updateRequest.MediaPortsStart
		}
		if updateRequest.MediaPortsEnd != nil {
			mockState.MediaPortsEnd = *updateRequest.MediaPortsEnd
		}
		if updateRequest.SignallingPortsStart != nil {
			mockState.SignallingPortsStart = *updateRequest.SignallingPortsStart
		}
		if updateRequest.SignallingPortsEnd != nil {
			mockState.SignallingPortsEnd = *updateRequest.SignallingPortsEnd
		}
		if updateRequest.GuestsOnlyTimeout != nil {
			mockState.GuestsOnlyTimeout = *updateRequest.GuestsOnlyTimeout
		}
		if updateRequest.WaitingForChairTimeout != nil {
			mockState.WaitingForChairTimeout = *updateRequest.WaitingForChairTimeout
		}

		// Return updated state
		*global_configuration = *mockState
	}).Maybe()

	testInfinityGlobalConfiguration(t, client)
}

func testInfinityGlobalConfiguration(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_global_configuration_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_global_configuration.global_configuration-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_webrtc", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_sip", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_h323", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_rtmp", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "crypto_mode", "besteffort"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "max_pixels_per_second", "720000"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "bursting_enabled", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "cloud_provider", "aws"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "conference_create_permissions", "user_admin"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "conference_creation_mode", "per_cluster"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_analytics", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_error_reporting", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "bandwidth_restrictions", "restricted"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "administrator_email", "test@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "media_ports_start", "40000"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "media_ports_end", "40100"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "signalling_ports_start", "5060"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "signalling_ports_end", "5070"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "guests_only_timeout", "300"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "waiting_for_chair_timeout", "600"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_global_configuration_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_global_configuration.global_configuration-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_webrtc", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_sip", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_h323", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_rtmp", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "crypto_mode", "required"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "max_pixels_per_second", "1080000"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "bursting_enabled", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "cloud_provider", "azure"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "conference_create_permissions", "admin_only"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "conference_creation_mode", "per_node"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_analytics", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_error_reporting", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "bandwidth_restrictions", "none"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "administrator_email", "updated@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "media_ports_start", "50000"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "media_ports_end", "50100"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "signalling_ports_start", "5080"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "signalling_ports_end", "5090"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "guests_only_timeout", "600"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "waiting_for_chair_timeout", "900"),
				),
			},
		},
	})
}
