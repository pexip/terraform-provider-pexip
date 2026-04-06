/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"os"
	"regexp"
	"testing"

	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityGlobalConfiguration(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking - include all fields that are set in config
	defaultTheme := &config.IVRTheme{Name: "Pexip theme (English_UK)"}
	mockState := &config.GlobalConfiguration{
		ID:                           1,
		ResourceURI:                  "/api/admin/configuration/v1/global/1/",
		AWSAccessKey:                 test.StringPtr("non-default-aws-access-key"),
		AWSSecretKey:                 test.StringPtr("non-default-aws-secret-key"),
		AzureClientID:                test.StringPtr("11111111-2222-3333-4444-555555555555"),
		AzureSecret:                  test.StringPtr("non-default-azure-secret"),
		AzureSubscriptionID:          test.StringPtr("22222222-3333-4444-5555-666666666666"),
		AzureTenant:                  test.StringPtr("33333333-4444-5555-6666-777777777777"),
		BdpmMaxPinFailuresPerWindow:  99,               // default: 20
		BdpmMaxScanAttemptsPerWindow: 88,               // default: 20
		BdpmPinChecksEnabled:         false,            // default: true
		BdpmScanQuarantineEnabled:    false,            // default: true
		BurstingEnabled:              true,             // default: false
		BurstingMinLifetime:          test.IntPtr(123), // default: 50
		BurstingThreshold:            test.IntPtr(77),  // default: 5
		CloudProvider:                "AWS",            // default: "AWS"
		ContactEmailAddress:          "someone@notdefault.com",
		ContentSecurityPolicyState:   false, // default: true
		CryptoMode:                   "on",  // default: "besteffort"
		DefaultTheme:                 defaultTheme,
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
		EnableClock:                         true,  // default: false
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
		EsMaximumRetryBackoff:               444,   // default: 1800
		EsMediaStreamsWait:                  55,    // default: 1
		EsMetricsUpdateInterval:             66,    // default: 60
		EsShortTermMemoryExpiration:         77,    // default: 2
		ExternalParticipantAvatarLookup:     false, // default: true
		GcpClientEmail:                      test.StringPtr("notdefault@gcp.com"),
		GcpPrivateKey:                       test.StringPtr("notdefaultkey"),
		GcpProjectID:                        test.StringPtr("notdefaultproject"),
		GuestsOnlyTimeout:                   999, // default: 60
		LegacyAPIUsername:                   "notdefaultuser",
		LegacyAPIPassword:                   "notdefaultpass",
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
		MediaPortsEnd:                       49999,                // default: 49999
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
		SignallingPortsStart:        33001, // default: 33000
		SipTLSCertVerifyMode:        "ON",  // default: "OFF"
		SiteBanner:                  "notdefaultsitebanner",
		SiteBannerBg:                "#ffffff", // default: "#c0c0c0"
		SiteBannerFg:                "#ff0000", // default: "#000000"
		TeamsEnablePowerpointRender: true,      // default: false
		WaitingForChairTimeout: 901, // default: 900
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

	// Mock the Delete (resets all fields to schema defaults).
	// Identified by BurstingEnabled==false && GuestsOnlyTimeout==60, which only Delete sends.
	client.On("PatchJSON", mock.Anything, "configuration/v1/global/1/",
		mock.MatchedBy(func(req *config.GlobalConfigurationUpdateRequest) bool {
			return !req.BurstingEnabled && req.GuestsOnlyTimeout == 60
		}), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.GlobalConfigurationUpdateRequest)

		// Assert every field is set to its schema default.
		assert.Nil(t, req.AWSAccessKey)
		assert.Nil(t, req.AWSSecretKey)
		assert.Nil(t, req.AzureClientID)
		assert.Nil(t, req.AzureSecret)
		assert.Nil(t, req.AzureSubscriptionID)
		assert.Nil(t, req.AzureTenant)
		assert.Equal(t, 50, *req.BurstingMinLifetime)
		assert.Equal(t, 5, *req.BurstingThreshold)
		assert.Nil(t, req.DefaultTheme)
		assert.Nil(t, req.DefaultWebappAlias)
		assert.Nil(t, req.GcpClientEmail)
		assert.Equal(t, "", *req.GcpPrivateKey)
		assert.Nil(t, req.GcpProjectID)
		assert.Equal(t, 0, *req.ManagementQos)
		assert.Nil(t, req.MaxCallrateIn)
		assert.Nil(t, req.MaxCallrateOut)
		assert.Equal(t, 20, req.BdpmMaxPinFailuresPerWindow)
		assert.Equal(t, 20, req.BdpmMaxScanAttemptsPerWindow)
		assert.True(t, req.BdpmPinChecksEnabled)
		assert.True(t, req.BdpmScanQuarantineEnabled)
		assert.False(t, req.BurstingEnabled)
		assert.Equal(t, "AWS", req.CloudProvider)
		assert.Equal(t, "", req.ContactEmailAddress)
		assert.True(t, req.ContentSecurityPolicyState)
		assert.Equal(t, "besteffort", req.CryptoMode)
		assert.Equal(t, 0, req.EjectLastParticipantBackstopTimeout)
		assert.False(t, req.EnableAnalytics)
		assert.True(t, req.EnableApplicationAPI)
		assert.False(t, req.EnableBreakoutRooms)
		assert.True(t, req.EnableChat)
		assert.False(t, req.EnableClock)
		assert.True(t, req.EnableDenoise)
		assert.True(t, req.EnableDialout)
		assert.True(t, req.EnableDirectory)
		assert.True(t, req.EnableEdgeNonMesh)
		assert.True(t, req.EnableFecc)
		assert.True(t, req.EnableH323)
		assert.False(t, req.EnableLegacyDialoutAPI)
		assert.False(t, req.EnableLyncAutoEscalate)
		assert.False(t, req.EnableLyncVbss)
		assert.False(t, req.EnableMlvad)
		assert.True(t, req.EnableRTMP)
		assert.True(t, req.EnableSIP)
		assert.False(t, req.EnableSIPUDP)
		assert.True(t, req.EnableSoftmute)
		assert.True(t, req.EnableSSH)
		assert.False(t, req.EnableTurn443)
		assert.True(t, req.EnableWebRTC)
		assert.False(t, req.ErrorReportingEnabled)
		assert.Equal(t, "https://acr.pexip.com", req.ErrorReportingURL)
		assert.Equal(t, 7, req.EsConnectionTimeout)
		assert.Equal(t, 1, req.EsInitialRetryBackoff)
		assert.Equal(t, 1000, req.EsMaximumDeferredPosts)
		assert.Equal(t, 1800, req.EsMaximumRetryBackoff)
		assert.Equal(t, 1, req.EsMediaStreamsWait)
		assert.Equal(t, 60, req.EsMetricsUpdateInterval)
		assert.Equal(t, 2, req.EsShortTermMemoryExpiration)
		assert.True(t, req.ExternalParticipantAvatarLookup)
		assert.Equal(t, 60, req.GuestsOnlyTimeout)
		assert.Equal(t, "", req.LegacyAPIUsername)
		assert.Equal(t, "", req.LegacyAPIPassword)
		assert.False(t, req.LiveCaptionsVMRDefault)
		assert.True(t, req.LiveviewShowConferences)
		assert.Equal(t, "", req.LocalMssipDomain)
		assert.Equal(t, "", req.LogonBanner)
		assert.Equal(t, 0, req.LogsMaxAge)
		assert.Equal(t, 30, req.ManagementSessionTimeout)
		assert.Equal(t, "/admin/conferencingstatus/deploymentgraph/deployment_graph/", req.ManagementStartPage)
		assert.Equal(t, "hd", req.MaxPixelsPerSecond)
		assert.Equal(t, 75, req.MaxPresentationBandwidthRatio)
		assert.Equal(t, 49999, req.MediaPortsEnd)
		assert.Equal(t, 40000, req.MediaPortsStart)
		assert.Equal(t, "", req.OcspResponderURL)
		assert.Equal(t, "OFF", req.OcspState)
		assert.Equal(t, 120, req.PinEntryTimeout)
		assert.True(t, req.SessionTimeoutEnabled)
		assert.Equal(t, 39999, req.SignallingPortsEnd)
		assert.Equal(t, 33000, req.SignallingPortsStart)
		assert.Equal(t, "OFF", req.SipTLSCertVerifyMode)
		assert.Equal(t, "", req.SiteBanner)
		assert.Equal(t, "#c0c0c0", req.SiteBannerBg)
		assert.Equal(t, "#000000", req.SiteBannerFg)
		assert.False(t, req.TeamsEnablePowerpointRender)
		assert.Equal(t, 900, req.WaitingForChairTimeout)

		// Update mockState to reflect the reset defaults so subsequent reads are correct.
		mockState.AWSAccessKey = nil
		mockState.AWSSecretKey = nil
		mockState.AzureClientID = nil
		mockState.AzureSecret = nil
		mockState.AzureSubscriptionID = nil
		mockState.AzureTenant = nil
		mockState.BdpmMaxPinFailuresPerWindow = 20
		mockState.BdpmMaxScanAttemptsPerWindow = 20
		mockState.BdpmPinChecksEnabled = true
		mockState.BdpmScanQuarantineEnabled = true
		mockState.BurstingEnabled = false
		mockState.BurstingMinLifetime = test.IntPtr(50)
		mockState.BurstingThreshold = test.IntPtr(5)
		mockState.CloudProvider = "AWS"
		mockState.ContactEmailAddress = ""
		mockState.ContentSecurityPolicyState = true
		mockState.CryptoMode = "besteffort"
		mockState.DefaultTheme = nil
		mockState.DefaultWebappAlias = nil
		mockState.EjectLastParticipantBackstopTimeout = 0
		mockState.EnableAnalytics = false
		mockState.EnableApplicationAPI = true
		mockState.EnableBreakoutRooms = false
		mockState.EnableChat = true
		mockState.EnableClock = false
		mockState.EnableDenoise = true
		mockState.EnableDialout = true
		mockState.EnableDirectory = true
		mockState.EnableEdgeNonMesh = true
		mockState.EnableFecc = true
		mockState.EnableH323 = true
		mockState.EnableLegacyDialoutAPI = false
		mockState.EnableLyncAutoEscalate = false
		mockState.EnableLyncVbss = false
		mockState.EnableMlvad = false
		mockState.EnableRTMP = true
		mockState.EnableSIP = true
		mockState.EnableSIPUDP = false
		mockState.EnableSoftmute = true
		mockState.EnableSSH = true
		mockState.EnableTurn443 = false
		mockState.EnableWebRTC = true
		mockState.ErrorReportingEnabled = false
		mockState.ErrorReportingURL = "https://acr.pexip.com"
		mockState.EsConnectionTimeout = 7
		mockState.EsInitialRetryBackoff = 1
		mockState.EsMaximumDeferredPosts = 1000
		mockState.EsMaximumRetryBackoff = 1800
		mockState.EsMediaStreamsWait = 1
		mockState.EsMetricsUpdateInterval = 60
		mockState.EsShortTermMemoryExpiration = 2
		mockState.ExternalParticipantAvatarLookup = true
		mockState.GcpClientEmail = nil
		mockState.GcpPrivateKey = test.StringPtr("")
		mockState.GcpProjectID = nil
		mockState.GuestsOnlyTimeout = 60
		mockState.LegacyAPIUsername = ""
		mockState.LegacyAPIPassword = ""
		mockState.LiveCaptionsVMRDefault = false
		mockState.LiveviewShowConferences = true
		mockState.LocalMssipDomain = ""
		mockState.LogonBanner = ""
		mockState.LogsMaxAge = 0
		mockState.ManagementQos = test.IntPtr(0)
		mockState.ManagementSessionTimeout = 30
		mockState.ManagementStartPage = "/admin/conferencingstatus/deploymentgraph/deployment_graph/"
		mockState.MaxCallrateIn = nil
		mockState.MaxCallrateOut = nil
		mockState.MaxPixelsPerSecond = "hd"
		mockState.MaxPresentationBandwidthRatio = 75
		mockState.MediaPortsEnd = 49999
		mockState.MediaPortsStart = 40000
		mockState.OcspResponderURL = ""
		mockState.OcspState = "OFF"
		mockState.PinEntryTimeout = 120
		mockState.SessionTimeoutEnabled = true
		mockState.SignallingPortsEnd = 39999
		mockState.SignallingPortsStart = 33000
		mockState.SipTLSCertVerifyMode = "OFF"
		mockState.SiteBanner = ""
		mockState.SiteBannerBg = "#c0c0c0"
		mockState.SiteBannerFg = "#000000"
		mockState.TeamsEnablePowerpointRender = false
		mockState.WaitingForChairTimeout = 900
	}).Once()

	// General PatchJSON mock — handles all create and update calls.
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
		mockState.DefaultWebappAlias = updateRequest.DefaultWebappAlias
		mockState.DeploymentUUID = updateRequest.DeploymentUUID
		mockState.DisabledCodecs = updateRequest.DisabledCodecs
		mockState.EjectLastParticipantBackstopTimeout = updateRequest.EjectLastParticipantBackstopTimeout
		mockState.EnableAnalytics = updateRequest.EnableAnalytics
		mockState.EnableApplicationAPI = updateRequest.EnableApplicationAPI
		mockState.EnableBreakoutRooms = updateRequest.EnableBreakoutRooms
		mockState.EnableChat = updateRequest.EnableChat
		mockState.EnableClock = updateRequest.EnableClock
		mockState.EnableDenoise = updateRequest.EnableDenoise
		mockState.EnableDialout = updateRequest.EnableDialout
		mockState.EnableDirectory = updateRequest.EnableDirectory
		mockState.EnableEdgeNonMesh = updateRequest.EnableEdgeNonMesh
		mockState.EnableFecc = updateRequest.EnableFecc
		mockState.EnableLegacyDialoutAPI = updateRequest.EnableLegacyDialoutAPI
		mockState.EnableLyncAutoEscalate = updateRequest.EnableLyncAutoEscalate
		mockState.EnableLyncVbss = updateRequest.EnableLyncVbss
		mockState.EnableMlvad = updateRequest.EnableMlvad
		mockState.EnableSIPUDP = updateRequest.EnableSIPUDP
		mockState.EnableSoftmute = updateRequest.EnableSoftmute
		mockState.EnableSSH = updateRequest.EnableSSH
		mockState.EnableTurn443 = updateRequest.EnableTurn443
		mockState.ErrorReportingEnabled = updateRequest.ErrorReportingEnabled
		mockState.ErrorReportingURL = updateRequest.ErrorReportingURL
		mockState.EsConnectionTimeout = updateRequest.EsConnectionTimeout
		mockState.EsInitialRetryBackoff = updateRequest.EsInitialRetryBackoff
		mockState.EsMaximumRetryBackoff = updateRequest.EsMaximumRetryBackoff
		mockState.EsMediaStreamsWait = updateRequest.EsMediaStreamsWait
		mockState.EsMetricsUpdateInterval = updateRequest.EsMetricsUpdateInterval
		mockState.EsShortTermMemoryExpiration = updateRequest.EsShortTermMemoryExpiration
		mockState.ExternalParticipantAvatarLookup = updateRequest.ExternalParticipantAvatarLookup
		mockState.GcpClientEmail = updateRequest.GcpClientEmail
		mockState.GcpPrivateKey = updateRequest.GcpPrivateKey
		mockState.GcpProjectID = updateRequest.GcpProjectID
		mockState.GuestsOnlyTimeout = updateRequest.GuestsOnlyTimeout
		mockState.LegacyAPIUsername = updateRequest.LegacyAPIUsername
		mockState.LegacyAPIPassword = updateRequest.LegacyAPIPassword
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
			// Step 1: Apply full configuration with all fields set to non-default values.
			{
				Config: test.LoadTestFolder(t, "resource_infinity_global_configuration_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_global_configuration.global_configuration-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "bursting_enabled", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "bursting_min_lifetime", "123"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "bursting_threshold", "77"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "crypto_mode", "on"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "max_pixels_per_second", "fullhd"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_analytics", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_webrtc", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_sip", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_h323", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_rtmp", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_sip_udp", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_application_api", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_breakout_rooms", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_chat", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_ssh", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "session_timeout_enabled", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "bdpm_pin_checks_enabled", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "bdpm_scan_quarantine_enabled", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "external_participant_avatar_lookup", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "liveview_show_conferences", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "guests_only_timeout", "999"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "max_callrate_in", "888"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "max_callrate_out", "999"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "media_ports_start", "40123"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "media_ports_end", "49999"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "signalling_ports_start", "33001"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "signalling_ports_end", "39998"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "ocsp_state", "ON"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "sip_tls_cert_verify_mode", "ON"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "site_banner", "notdefaultsitebanner"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "site_banner_bg", "#ffffff"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "site_banner_fg", "#ff0000"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "logon_banner", "notdefaultbanner"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "default_theme", "dark"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "disabled_codecs.#", "2"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_global_configuration.global_configuration-test", "disabled_codecs.*", "H264_H_1"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_global_configuration.global_configuration-test", "disabled_codecs.*", "VP9"),
				),
			},
			// Step 2: Destroy — triggers Delete which must reset all fields to schema defaults.
			// The delete-specific PatchJSON mock above asserts every expected value.
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_global_configuration_full"),
				Destroy: true,
			},
			// Step 3: Re-apply min config (just logon_banner) and verify the real API
			// returned all other fields to their defaults after the destroy.
			{
				Config: test.LoadTestFolder(t, "resource_infinity_global_configuration_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_global_configuration.global_configuration-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "logon_banner", "test-logon-banner"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "bursting_enabled", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "crypto_mode", "besteffort"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "max_pixels_per_second", "hd"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_analytics", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_webrtc", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_sip", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_h323", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_rtmp", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_sip_udp", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_application_api", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_breakout_rooms", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_chat", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "enable_ssh", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "session_timeout_enabled", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "bdpm_pin_checks_enabled", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "bdpm_scan_quarantine_enabled", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "external_participant_avatar_lookup", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "liveview_show_conferences", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "guests_only_timeout", "60"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "media_ports_start", "40000"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "media_ports_end", "49999"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "signalling_ports_start", "33000"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "signalling_ports_end", "39999"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "ocsp_state", "OFF"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "sip_tls_cert_verify_mode", "OFF"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "site_banner", ""),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "site_banner_bg", "#c0c0c0"),
					resource.TestCheckResourceAttr("pexip_infinity_global_configuration.global_configuration-test", "site_banner_fg", "#000000"),
				),
			},
		},
	})
}

func TestInfinityGlobalConfigurationValidation(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	client := infinity.NewClientMock()

	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				// bursting_enabled=true with no cloud_provider (defaults to AWS) but no AWS keys — expect errors
				Config: `
resource "pexip_infinity_global_configuration" "global_configuration-test" {
  bursting_enabled = true
}
`,
				ExpectError: regexp.MustCompile(`aws_access_key must be configured|aws_secret_key must be configured`),
			},
			{
				// bursting_enabled=true with cloud_provider=AWS but no AWS keys — expect errors for both missing keys
				Config: `
resource "pexip_infinity_global_configuration" "global_configuration-test" {
  bursting_enabled = true
  cloud_provider   = "AWS"
}
`,
				ExpectError: regexp.MustCompile(`aws_access_key must be configured|aws_secret_key must be configured`),
			},
			{
				// bursting_enabled=true with cloud_provider=AWS and only the access key — expect error for missing secret key
				Config: `
resource "pexip_infinity_global_configuration" "global_configuration-test" {
  bursting_enabled = true
  cloud_provider   = "AWS"
  aws_access_key   = "test-access-key"
}
`,
				ExpectError: regexp.MustCompile(`aws_secret_key must be configured`),
			},
		},
	})
}
