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
		ID:                           1,
		ResourceURI:                  "/api/admin/configuration/v1/global/1/",
		AWSAccessKey:                 test.StringPtr("non-default-aws-access-key"),
		AWSSecretKey:                 test.StringPtr("non-default-aws-secret-key"),
		AzureClientID:                test.StringPtr("non-default-azure-client-id"),
		AzureSecret:                  test.StringPtr("non-default-azure-secret"),
		AzureSubscriptionID:          test.StringPtr("non-default-azure-sub-id"),
		AzureTenant:                  test.StringPtr("non-default-azure-tenant"),
		BdpmMaxPinFailuresPerWindow:  99,               // default: 20
		BdpmMaxScanAttemptsPerWindow: 88,               // default: 20
		BdpmPinChecksEnabled:         false,            // default: true
		BdpmScanQuarantineEnabled:    false,            // default: true
		BurstingEnabled:              true,             // default: false
		BurstingMinLifetime:          test.IntPtr(123), // default: 50
		BurstingThreshold:            test.IntPtr(77),  // default: 5
		CloudProvider:                "GCP",            // default: "AWS"
		ContactEmailAddress:          "someone@notdefault.com",
		ContentSecurityPolicyHeader:  "custom-csp-header",
		ContentSecurityPolicyState:   false, // default: true
		CryptoMode:                   "on",  // default: "besteffort"
		DefaultTheme:                 test.StringPtr("dark"),
		DefaultToNewWebapp:           false,    // default: true
		DefaultWebapp:                "legacy", // default: "latest"
		DefaultWebappAlias:           test.StringPtr("custom-alias"),
		DeploymentUUID:               "custom-uuid",
		DisabledCodecs: []config.CodecValue{
			{Value: "H264_H_1"},
			{Value: "VP9"},
		},
		EjectLastParticipantBackstopTimeout: 10,    // default: 0
		EnableAnalytics:                     true,  // default: false
		EnableApplicationAPI:                false, // default: true
		EnableBreakoutRooms:                 true,  // default: false
		EnableChat:                          false, // default: true
		EnableDenoise:                       false, // default: true
		EnableDialout:                       false, // default: true
		EnableDirectory:                     false, // default: true
		EnableEdgeNonMesh:                   false, // default: true
		EnableFecc:                          false, // default: true
		EnableH323:                          false, // default: true
		EnableLegacyDialoutAPI:              true,  // default: false
		EnableLyncAutoEscalate:              true,  // default: false
		EnableLyncVbss:                      true,  // default: false
		EnableMlvad:                         true,  // default: false
		EnableMultiscreen:                   false, // default: true
		EnablePushNotifications:             true,  // default: false
		EnableRTMP:                          false, // default: true
		EnableSIP:                           false, // default: true
		EnableSIPUDP:                        true,  // default: false
		EnableSoftmute:                      false, // default: true
		EnableSSH:                           false, // default: true
		EnableTurn443:                       true,  // default: false
		EnableWebRTC:                        false, // default: true
		ErrorReportingEnabled:               true,  // default: false
		ErrorReportingURL:                   "https://custom-error-reporting.com",
		EsConnectionTimeout:                 99,    // default: 7
		EsInitialRetryBackoff:               22,    // default: 1
		EsMaximumDeferredPosts:              333,   // default: 1000
		EsMaximumRetryBackoff:               444,   // default: 1800
		EsMediaStreamsWait:                  55,    // default: 1
		EsMetricsUpdateInterval:             66,    // default: 60
		EsShortTermMemoryExpiration:         77,    // default: 2
		ExternalParticipantAvatarLookup:     false, // default: true
		GcpClientEmail:                      test.StringPtr("notdefault@gcp.com"),
		GcpPrivateKey:                       test.StringPtr("notdefaultkey"),
		GcpProjectID:                        test.StringPtr("notdefaultproject"),
		GuestsOnlyTimeout:                   999,  // default: 60
		LegacyAPIHTTP:                       true, // default: false
		LegacyAPIUsername:                   "notdefaultuser",
		LegacyAPIPassword:                   "notdefaultpass",
		LiveCaptionsAPIGateway:              "notdefaultgateway",
		LiveCaptionsAppID:                   "notdefaultappid",
		LiveCaptionsEnabled:                 false, // default: true
		LiveCaptionsPublicJWTKey:            "notdefaultjwtkey",
		LiveCaptionsVMRDefault:              true,  // default: false
		LiveviewShowConferences:             false, // default: true
		LocalMssipDomain:                    "notdefaultdomain",
		LogonBanner:                         "notdefaultbanner",
		LogsMaxAge:                          123,                  // default: 0
		ManagementQos:                       test.IntPtr(77),      // default: 0
		ManagementSessionTimeout:            99,                   // default: 30
		ManagementStartPage:                 "/custom/start/page", // default: "/admin/conferencingstatus/deploymentgraph/deployment_graph/"
		MaxCallrateIn:                       test.IntPtr(888),     // default: nil
		MaxCallrateOut:                      test.IntPtr(999),     // default: nil
		MaxPixelsPerSecond:                  "fullhd",             // default: "hd"
		MaxPresentationBandwidthRatio:       33,                   // default: 75
		MediaPortsEnd:                       50123,                // default: 49999
		MediaPortsStart:                     40123,                // default: 40000
		OcspResponderURL:                    "notdefaultocsp",
		OcspState:                           "ON", // default: "OFF"
		PinEntryTimeout:                     321,  // default: 120
		// PssCustomerID:                       "notdefaultpsscustomer",
		// PssEnabled:                          true, // default: false
		// PssGateway:                          "notdefaultpssgateway",
		// PssToken:                            "notdefaultpsstoken",
		SessionTimeoutEnabled:       false, // default: true
		SignallingPortsEnd:          39998, // default: 39999
		SignallingPortsStart:        34001, // default: 34000
		SipTLSCertVerifyMode:        "ON",  // default: "OFF"
		SiteBanner:                  "notdefaultsitebanner",
		SiteBannerBg:                "#ffffff", // default: "#c0c0c0"
		SiteBannerFg:                "#ff0000", // default: "#000000"
		TeamsEnablePowerpointRender: true,      // default: false
		// WaitingForChairTimeout:              901,       // default: 900
	}

	// Mock the GetGlobalconfiguration API call for Read operations
	client.On(
		"GetJSON",
		mock.Anything,                // context.Context
		"configuration/v1/global/1/", // string
		mock.Anything,                // *url.Values or nil
		mock.AnythingOfType("*config.GlobalConfiguration"), // pointer to config.GlobalConfiguration
	).Return(nil).Run(func(args mock.Arguments) {
		global_configuration := args.Get(3).(*config.GlobalConfiguration)
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
		mockState.BurstingMinLifetime = updateRequest.BurstingMinLifetime
		mockState.BurstingThreshold = updateRequest.BurstingThreshold
		mockState.CryptoMode = updateRequest.CryptoMode
		mockState.MaxPixelsPerSecond = updateRequest.MaxPixelsPerSecond
		mockState.CloudProvider = updateRequest.CloudProvider
		mockState.AWSAccessKey = updateRequest.AWSAccessKey
		mockState.AWSSecretKey = updateRequest.AWSSecretKey
		mockState.AzureClientID = updateRequest.AzureClientID
		mockState.AzureSecret = updateRequest.AzureSecret
		mockState.AzureSubscriptionID = updateRequest.AzureSubscriptionID
		mockState.AzureTenant = updateRequest.AzureTenant
		mockState.BdpmMaxPinFailuresPerWindow = updateRequest.BdpmMaxPinFailuresPerWindow
		mockState.BdpmMaxScanAttemptsPerWindow = updateRequest.BdpmMaxScanAttemptsPerWindow
		mockState.BdpmPinChecksEnabled = updateRequest.BdpmPinChecksEnabled
		mockState.BdpmScanQuarantineEnabled = updateRequest.BdpmScanQuarantineEnabled
		mockState.ContactEmailAddress = updateRequest.ContactEmailAddress
		mockState.ContentSecurityPolicyHeader = updateRequest.ContentSecurityPolicyHeader
		mockState.ContentSecurityPolicyState = updateRequest.ContentSecurityPolicyState
		mockState.DefaultTheme = updateRequest.DefaultTheme
		mockState.DefaultToNewWebapp = updateRequest.DefaultToNewWebapp
		mockState.DefaultWebapp = updateRequest.DefaultWebapp
		mockState.DefaultWebappAlias = updateRequest.DefaultWebappAlias
		mockState.DeploymentUUID = updateRequest.DeploymentUUID
		mockState.DisabledCodecs = updateRequest.DisabledCodecs
		mockState.EjectLastParticipantBackstopTimeout = updateRequest.EjectLastParticipantBackstopTimeout
		mockState.EnableAnalytics = updateRequest.EnableAnalytics
		mockState.EnableApplicationAPI = updateRequest.EnableApplicationAPI
		mockState.EnableBreakoutRooms = updateRequest.EnableBreakoutRooms
		mockState.EnableChat = updateRequest.EnableChat
		mockState.EnableDenoise = updateRequest.EnableDenoise
		mockState.EnableDialout = updateRequest.EnableDialout
		mockState.EnableDirectory = updateRequest.EnableDirectory
		mockState.EnableEdgeNonMesh = updateRequest.EnableEdgeNonMesh
		mockState.EnableFecc = updateRequest.EnableFecc
		mockState.EnableLegacyDialoutAPI = updateRequest.EnableLegacyDialoutAPI
		mockState.EnableLyncAutoEscalate = updateRequest.EnableLyncAutoEscalate
		mockState.EnableLyncVbss = updateRequest.EnableLyncVbss
		mockState.EnableMlvad = updateRequest.EnableMlvad
		mockState.EnableMultiscreen = updateRequest.EnableMultiscreen
		mockState.EnablePushNotifications = updateRequest.EnablePushNotifications
		mockState.EnableSIPUDP = updateRequest.EnableSIPUDP
		mockState.EnableSoftmute = updateRequest.EnableSoftmute
		mockState.EnableSSH = updateRequest.EnableSSH
		mockState.EnableTurn443 = updateRequest.EnableTurn443
		mockState.ErrorReportingEnabled = updateRequest.ErrorReportingEnabled
		mockState.ErrorReportingURL = updateRequest.ErrorReportingURL
		mockState.EsConnectionTimeout = updateRequest.EsConnectionTimeout
		mockState.EsInitialRetryBackoff = updateRequest.EsInitialRetryBackoff
		mockState.EsMaximumDeferredPosts = updateRequest.EsMaximumDeferredPosts
		mockState.EsMaximumRetryBackoff = updateRequest.EsMaximumRetryBackoff
		mockState.EsMediaStreamsWait = updateRequest.EsMediaStreamsWait
		mockState.EsMetricsUpdateInterval = updateRequest.EsMetricsUpdateInterval
		mockState.EsShortTermMemoryExpiration = updateRequest.EsShortTermMemoryExpiration
		mockState.ExternalParticipantAvatarLookup = updateRequest.ExternalParticipantAvatarLookup
		mockState.GcpClientEmail = updateRequest.GcpClientEmail
		mockState.GcpPrivateKey = updateRequest.GcpPrivateKey
		mockState.GcpProjectID = updateRequest.GcpProjectID
		mockState.GuestsOnlyTimeout = updateRequest.GuestsOnlyTimeout
		mockState.LegacyAPIHTTP = updateRequest.LegacyAPIHTTP
		mockState.LegacyAPIUsername = updateRequest.LegacyAPIUsername
		mockState.LegacyAPIPassword = updateRequest.LegacyAPIPassword
		mockState.LiveCaptionsAPIGateway = updateRequest.LiveCaptionsAPIGateway
		mockState.LiveCaptionsAppID = updateRequest.LiveCaptionsAppID
		mockState.LiveCaptionsEnabled = updateRequest.LiveCaptionsEnabled
		mockState.LiveCaptionsPublicJWTKey = updateRequest.LiveCaptionsPublicJWTKey
		mockState.LiveCaptionsVMRDefault = updateRequest.LiveCaptionsVMRDefault
		mockState.LiveviewShowConferences = updateRequest.LiveviewShowConferences
		mockState.LocalMssipDomain = updateRequest.LocalMssipDomain
		mockState.LogonBanner = updateRequest.LogonBanner
		mockState.LogsMaxAge = updateRequest.LogsMaxAge
		mockState.ManagementQos = updateRequest.ManagementQos
		mockState.ManagementSessionTimeout = updateRequest.ManagementSessionTimeout
		mockState.ManagementStartPage = updateRequest.ManagementStartPage
		mockState.MaxCallrateIn = updateRequest.MaxCallrateIn
		mockState.MaxCallrateOut = updateRequest.MaxCallrateOut
		mockState.MaxPresentationBandwidthRatio = updateRequest.MaxPresentationBandwidthRatio
		mockState.MediaPortsEnd = updateRequest.MediaPortsEnd
		mockState.MediaPortsStart = updateRequest.MediaPortsStart
		mockState.OcspResponderURL = updateRequest.OcspResponderURL
		mockState.OcspState = updateRequest.OcspState
		mockState.PinEntryTimeout = updateRequest.PinEntryTimeout
		mockState.PssCustomerID = updateRequest.PssCustomerID
		mockState.PssEnabled = updateRequest.PssEnabled
		mockState.PssGateway = updateRequest.PssGateway
		mockState.PssToken = updateRequest.PssToken
		mockState.SessionTimeoutEnabled = updateRequest.SessionTimeoutEnabled
		mockState.SignallingPortsEnd = updateRequest.SignallingPortsEnd
		mockState.SignallingPortsStart = updateRequest.SignallingPortsStart
		mockState.SipTLSCertVerifyMode = updateRequest.SipTLSCertVerifyMode
		mockState.SiteBanner = updateRequest.SiteBanner
		mockState.SiteBannerBg = updateRequest.SiteBannerBg
		mockState.SiteBannerFg = updateRequest.SiteBannerFg
		mockState.TeamsEnablePowerpointRender = updateRequest.TeamsEnablePowerpointRender
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
				Check:  resource.ComposeTestCheckFunc(
				// resource.TestCheckResourceAttrSet("pexip_infinity_global_configuration.global_configuration-test", "id"),
				// resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_webrtc", "true"),
				// resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_sip", "true"),
				// resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_h323", "true"),
				// resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_rtmp", "true"),
				// resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "crypto_mode", "besteffort"),
				// resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "max_pixels_per_second", "hd"),
				// resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "bursting_enabled", "true"),
				// resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "cloud_provider", "aws"),
				// resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "conference_creation_mode", "per_cluster"),
				// resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_analytics", "true"),
				// resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "media_ports_start", "40000"),
				// resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "media_ports_end", "40100"),
				// resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "signalling_ports_start", "5060"),
				// resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "signalling_ports_end", "5070"),
				// resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "guests_only_timeout", "300"),
				// resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "waiting_for_chair_timeout", "600"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_global_configuration_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_global_configuration.global_configuration-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_analytics", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_webrtc", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_sip", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_h323", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_rtmp", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "crypto_mode", "on"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "max_pixels_per_second", "fullhd"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "bursting_enabled", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "cloud_provider", "GCP"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "media_ports_start", "40123"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "media_ports_end", "50123"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "signalling_ports_start", "34001"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "signalling_ports_end", "39998"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "guests_only_timeout", "999"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "waiting_for_chair_timeout", "900"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "max_callrate_in", "888"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "max_callrate_out", "999"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_sip_udp", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "default_webapp_alias", "custom-alias"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "default_theme", "dark"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "disabled_codecs.#", "2"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_global_configuration.global_configuration-test", "disabled_codecs.*", "H264_H_1"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_global_configuration.global_configuration-test", "disabled_codecs.*", "VP9"),
					// resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "maximum_deferred_es_posts", "500"),
				),
			},
		},
	})
}
