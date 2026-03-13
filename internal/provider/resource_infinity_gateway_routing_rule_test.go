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
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityGatewayRoutingRule(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking - initialize with defaults
	mockState := &config.GatewayRoutingRule{
		ID:                              1,
		ResourceURI:                     "/api/admin/configuration/v1/gateway_routing_rule/1/",
		Name:                            "tf-test-gateway-routing-rule",
		Description:                     "",
		Priority:                        100,
		Enable:                          true,
		MatchString:                     ".*@example.com",
		ReplaceString:                   "",
		CalledDeviceType:                "external",
		OutgoingProtocol:                "sip",
		CallType:                        "video",
		DenoiseAudio:                    true,
		ExternalParticipantAvatarLookup: test.StringPtr("default"),
		LiveCaptionsEnabled:             "default",
		MatchIncomingCalls:              true,
		MatchIncomingH323:               true,
		MatchIncomingMSSIP:              true,
		MatchIncomingOnlyIfRegistered:   false,
		MatchIncomingSIP:                true,
		MatchIncomingTeams:              false,
		MatchIncomingWebRTC:             true,
		MatchOutgoingCalls:              false,
		MatchStringFull:                 false,
		Tag:                             "",
		TreatAsTrusted:                  false,
		CryptoMode:                      nil,
		MaxPixelsPerSecond:              nil,
		MaxCallrateIn:                   nil,
		MaxCallrateOut:                  nil,
	}

	// Step 1: Create with full config
	client.On("PostWithResponse", mock.Anything, "configuration/v1/gateway_routing_rule/", mock.Anything, mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/gateway_routing_rule/1/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.GatewayRoutingRuleCreateRequest)
		mockState.Name = req.Name
		mockState.Description = req.Description
		mockState.MatchString = req.MatchString
		mockState.Priority = req.Priority
		mockState.Tag = req.Tag
		mockState.ReplaceString = req.ReplaceString
		mockState.Enable = req.Enable
		mockState.MatchStringFull = req.MatchStringFull
		mockState.CalledDeviceType = req.CalledDeviceType
		mockState.OutgoingProtocol = req.OutgoingProtocol
		mockState.CallType = req.CallType
		mockState.CryptoMode = req.CryptoMode
		mockState.DenoiseAudio = req.DenoiseAudio
		mockState.MaxPixelsPerSecond = req.MaxPixelsPerSecond
		mockState.MaxCallrateIn = req.MaxCallrateIn
		mockState.MaxCallrateOut = req.MaxCallrateOut
		mockState.MatchIncomingCalls = req.MatchIncomingCalls
		mockState.MatchOutgoingCalls = req.MatchOutgoingCalls
		mockState.MatchIncomingSIP = req.MatchIncomingSIP
		mockState.MatchIncomingH323 = req.MatchIncomingH323
		mockState.MatchIncomingMSSIP = req.MatchIncomingMSSIP
		mockState.MatchIncomingWebRTC = req.MatchIncomingWebRTC
		mockState.MatchIncomingTeams = req.MatchIncomingTeams
		mockState.MatchIncomingOnlyIfRegistered = req.MatchIncomingOnlyIfRegistered
		mockState.ExternalParticipantAvatarLookup = req.ExternalParticipantAvatarLookup
		mockState.LiveCaptionsEnabled = req.LiveCaptionsEnabled
		mockState.TreatAsTrusted = req.TreatAsTrusted
	}).Once()

	// Step 2: Update to min config
	client.On("PutJSON", mock.Anything, "configuration/v1/gateway_routing_rule/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.GatewayRoutingRuleUpdateRequest)
		mockState.Name = req.Name
		mockState.Description = req.Description
		mockState.MatchString = req.MatchString
		mockState.Priority = req.Priority
		mockState.Tag = req.Tag
		mockState.ReplaceString = req.ReplaceString
		mockState.Enable = req.Enable
		mockState.MatchStringFull = req.MatchStringFull
		mockState.CalledDeviceType = req.CalledDeviceType
		mockState.OutgoingProtocol = req.OutgoingProtocol
		mockState.CallType = req.CallType
		mockState.CryptoMode = req.CryptoMode
		mockState.DenoiseAudio = req.DenoiseAudio
		mockState.MaxPixelsPerSecond = req.MaxPixelsPerSecond
		mockState.MaxCallrateIn = req.MaxCallrateIn
		mockState.MaxCallrateOut = req.MaxCallrateOut
		mockState.MatchIncomingCalls = req.MatchIncomingCalls
		mockState.MatchOutgoingCalls = req.MatchOutgoingCalls
		mockState.MatchIncomingSIP = req.MatchIncomingSIP
		mockState.MatchIncomingH323 = req.MatchIncomingH323
		mockState.MatchIncomingMSSIP = req.MatchIncomingMSSIP
		mockState.MatchIncomingWebRTC = req.MatchIncomingWebRTC
		mockState.MatchIncomingTeams = req.MatchIncomingTeams
		mockState.MatchIncomingOnlyIfRegistered = req.MatchIncomingOnlyIfRegistered
		mockState.ExternalParticipantAvatarLookup = req.ExternalParticipantAvatarLookup
		mockState.LiveCaptionsEnabled = req.LiveCaptionsEnabled
		mockState.TreatAsTrusted = req.TreatAsTrusted
		if args.Get(3) != nil {
			rule := args.Get(3).(*config.GatewayRoutingRule)
			*rule = *mockState
		}
	}).Once()

	// Step 3: Delete
	client.On("DeleteJSON", mock.Anything, "configuration/v1/gateway_routing_rule/1/", mock.Anything).Return(nil).Maybe()

	// Step 4: Recreate with min config
	client.On("PostWithResponse", mock.Anything, "configuration/v1/gateway_routing_rule/", mock.Anything, mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/gateway_routing_rule/1/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.GatewayRoutingRuleCreateRequest)
		mockState.Name = req.Name
		mockState.Description = req.Description
		mockState.MatchString = req.MatchString
		mockState.Priority = req.Priority
		mockState.Tag = req.Tag
		mockState.ReplaceString = req.ReplaceString
		mockState.Enable = req.Enable
		mockState.MatchStringFull = req.MatchStringFull
		mockState.CalledDeviceType = req.CalledDeviceType
		mockState.OutgoingProtocol = req.OutgoingProtocol
		mockState.CallType = req.CallType
		mockState.CryptoMode = req.CryptoMode
		mockState.DenoiseAudio = req.DenoiseAudio
		mockState.MaxPixelsPerSecond = req.MaxPixelsPerSecond
		mockState.MaxCallrateIn = req.MaxCallrateIn
		mockState.MaxCallrateOut = req.MaxCallrateOut
		mockState.MatchIncomingCalls = req.MatchIncomingCalls
		mockState.MatchOutgoingCalls = req.MatchOutgoingCalls
		mockState.MatchIncomingSIP = req.MatchIncomingSIP
		mockState.MatchIncomingH323 = req.MatchIncomingH323
		mockState.MatchIncomingMSSIP = req.MatchIncomingMSSIP
		mockState.MatchIncomingWebRTC = req.MatchIncomingWebRTC
		mockState.MatchIncomingTeams = req.MatchIncomingTeams
		mockState.MatchIncomingOnlyIfRegistered = req.MatchIncomingOnlyIfRegistered
		mockState.ExternalParticipantAvatarLookup = req.ExternalParticipantAvatarLookup
		mockState.LiveCaptionsEnabled = req.LiveCaptionsEnabled
		mockState.TreatAsTrusted = req.TreatAsTrusted
	}).Once()

	// Step 5: Update to full config
	client.On("PutJSON", mock.Anything, "configuration/v1/gateway_routing_rule/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.GatewayRoutingRuleUpdateRequest)
		mockState.Name = req.Name
		mockState.Description = req.Description
		mockState.MatchString = req.MatchString
		mockState.Priority = req.Priority
		mockState.Tag = req.Tag
		mockState.ReplaceString = req.ReplaceString
		mockState.Enable = req.Enable
		mockState.MatchStringFull = req.MatchStringFull
		mockState.CalledDeviceType = req.CalledDeviceType
		mockState.OutgoingProtocol = req.OutgoingProtocol
		mockState.CallType = req.CallType
		mockState.CryptoMode = req.CryptoMode
		mockState.DenoiseAudio = req.DenoiseAudio
		mockState.MaxPixelsPerSecond = req.MaxPixelsPerSecond
		mockState.MaxCallrateIn = req.MaxCallrateIn
		mockState.MaxCallrateOut = req.MaxCallrateOut
		mockState.MatchIncomingCalls = req.MatchIncomingCalls
		mockState.MatchOutgoingCalls = req.MatchOutgoingCalls
		mockState.MatchIncomingSIP = req.MatchIncomingSIP
		mockState.MatchIncomingH323 = req.MatchIncomingH323
		mockState.MatchIncomingMSSIP = req.MatchIncomingMSSIP
		mockState.MatchIncomingWebRTC = req.MatchIncomingWebRTC
		mockState.MatchIncomingTeams = req.MatchIncomingTeams
		mockState.MatchIncomingOnlyIfRegistered = req.MatchIncomingOnlyIfRegistered
		mockState.ExternalParticipantAvatarLookup = req.ExternalParticipantAvatarLookup
		mockState.LiveCaptionsEnabled = req.LiveCaptionsEnabled
		mockState.TreatAsTrusted = req.TreatAsTrusted
		if args.Get(3) != nil {
			rule := args.Get(3).(*config.GatewayRoutingRule)
			*rule = *mockState
		}
	}).Once()

	// Mock Read operations (GetJSON) - used throughout all steps
	client.On("GetJSON", mock.Anything, "configuration/v1/gateway_routing_rule/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		rule := args.Get(3).(*config.GatewayRoutingRule)
		*rule = *mockState
	}).Maybe()

	testInfinityGatewayRoutingRule(t, client)
}

func testInfinityGatewayRoutingRule(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				// Step 1: Create with full config
				Config: test.LoadTestFolder(t, "resource_infinity_gateway_routing_rule_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "name", "tf-test-gateway-routing-rule"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "description", "tf-test Gateway Routing Rule Description"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_string", ".*@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "priority", "100"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "enable", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_string_full", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "replace_string", "replaced@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "tag", "tf-test-tag"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "called_device_type", "registration"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "outgoing_protocol", "h323"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "call_type", "audio"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "crypto_mode", "on"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "denoise_audio", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "max_pixels_per_second", "fullhd"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "max_callrate_in", "2048"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "max_callrate_out", "4096"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_incoming_calls", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_outgoing_calls", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_incoming_sip", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_incoming_h323", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_incoming_mssip", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_incoming_webrtc", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_incoming_teams", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_incoming_only_if_registered", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "enable_participant_avatar_lookup", "yes"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "live_captions_enabled", "yes"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "treat_as_trusted", "true"),
				),
			},
			{
				// Step 2: Update to min config (clear optional fields, reset to defaults)
				Config: test.LoadTestFolder(t, "resource_infinity_gateway_routing_rule_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "name", "tf-test-gateway-routing-rule"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_string", ".*@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "priority", "100"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "enable", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_string_full", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "replace_string", ""),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "tag", ""),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "called_device_type", "external"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "outgoing_protocol", "sip"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "call_type", "video"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "denoise_audio", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_incoming_calls", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_outgoing_calls", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_incoming_sip", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_incoming_h323", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_incoming_mssip", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_incoming_webrtc", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_incoming_teams", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_incoming_only_if_registered", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "enable_participant_avatar_lookup", "default"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "live_captions_enabled", "default"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "treat_as_trusted", "false"),
				),
			},
			{
				// Step 3: Destroy
				Config:  test.LoadTestFolder(t, "resource_infinity_gateway_routing_rule_min"),
				Destroy: true,
			},
			{
				// Step 4: Create with min config (after destroy)
				Config: test.LoadTestFolder(t, "resource_infinity_gateway_routing_rule_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "name", "tf-test-gateway-routing-rule"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_string", ".*@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "priority", "100"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "tag", ""),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "replace_string", ""),
				),
			},
			{
				// Step 5: Update to full config
				Config: test.LoadTestFolder(t, "resource_infinity_gateway_routing_rule_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "name", "tf-test-gateway-routing-rule"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "description", "tf-test Gateway Routing Rule Description"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_string", ".*@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "priority", "100"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "enable", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_string_full", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "replace_string", "replaced@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "tag", "tf-test-tag"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "called_device_type", "registration"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "outgoing_protocol", "h323"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "call_type", "audio"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "crypto_mode", "on"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "denoise_audio", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "max_pixels_per_second", "fullhd"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "max_callrate_in", "2048"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "max_callrate_out", "4096"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_incoming_calls", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_outgoing_calls", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_incoming_sip", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_incoming_h323", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_incoming_mssip", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_incoming_webrtc", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_incoming_teams", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "match_incoming_only_if_registered", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "enable_participant_avatar_lookup", "yes"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "live_captions_enabled", "yes"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.tf-test-gateway-routing-rule", "treat_as_trusted", "true"),
				),
			},
		},
	})
}
