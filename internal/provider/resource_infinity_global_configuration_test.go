/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"os"
	"testing"

	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityGlobalConfiguration(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

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
		EnableAnalytics:             true,
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
	client.On("PatchJSON", mock.Anything, "configuration/v1/global/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.GlobalConfigurationUpdateRequest)
		global_configuration := args.Get(3).(*config.GlobalConfiguration)

		// Update mock state based on request
		mockState.EnableWebRTC = updateRequest.EnableWebRTC
		mockState.EnableSIP = updateRequest.EnableSIP
		mockState.EnableH323 = updateRequest.EnableH323
		mockState.EnableRTMP = updateRequest.EnableRTMP
		mockState.BurstingEnabled = updateRequest.BurstingEnabled
		mockState.CryptoMode = updateRequest.CryptoMode
		mockState.MaxPixelsPerSecond = updateRequest.MaxPixelsPerSecond
		mockState.CloudProvider = updateRequest.CloudProvider
		mockState.AWSAccessKey = updateRequest.AWSAccessKey
		mockState.AWSSecretKey = updateRequest.AWSSecretKey
		mockState.AzureClientID = updateRequest.AzureClientID
		mockState.AzureSecret = updateRequest.AzureSecret
		mockState.EnableAnalytics = updateRequest.EnableAnalytics
		mockState.MediaPortsStart = updateRequest.MediaPortsStart
		mockState.MediaPortsEnd = updateRequest.MediaPortsEnd
		mockState.SignallingPortsStart = updateRequest.SignallingPortsStart
		mockState.SignallingPortsEnd = updateRequest.SignallingPortsEnd
		mockState.GuestsOnlyTimeout = updateRequest.GuestsOnlyTimeout
		mockState.WaitingForChairTimeout = updateRequest.WaitingForChairTimeout

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
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "conference_creation_mode", "per_cluster"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_analytics", "true"),
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
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "conference_creation_mode", "per_node"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_analytics", "false"),
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
