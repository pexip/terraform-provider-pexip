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

func TestInfinityMsExchangeConnector(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateMsexchangeconnector API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/ms_exchange_connector/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/ms_exchange_connector/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Helper function to get string pointer
	stringPtr := func(s string) *string { return &s }

	// Shared state for mocking
	mockState := &config.MsExchangeConnector{
		ID:                             123,
		ResourceURI:                    "/api/admin/configuration/v1/ms_exchange_connector/123/",
		Name:                           "ms_exchange_connector-test",
		Description:                    "Test MsExchangeConnector",
		MeetingBufferBefore:            300,
		MeetingBufferAfter:             300,
		ScheduledAliasSuffixLength:     6,
		RoomMailboxEmailAddress:        stringPtr("test@example.com"),
		RoomMailboxName:                "ms_exchange_connector-test",
		URL:                            "https://example.com",
		Username:                       "ms_exchange_connector-test",
		Password:                       "test-value",
		AuthenticationMethod:           "OAUTH",
		AuthProvider:                   "AZURE",
		UUID:                           "test-value",
		ScheduledAliasPrefix:           stringPtr("test-value"),
		ScheduledAliasDomain:           "example.com",
		EnableDynamicVmrs:              true,
		EnablePersonalVmrs:             true,
		AllowNewUsers:                  true,
		DisableProxy:                   true,
		UseCustomAddInSources:          true,
		EnableAddinDebugLogs:           true,
		OauthClientID:                  stringPtr("test-value"),
		OauthClientSecret:              "test-value",
		OauthAuthEndpoint:              "test-value",
		OauthTokenEndpoint:             "test-value",
		OauthRedirectURI:               "test-value",
		OauthRefreshToken:              "test-value",
		OauthState:                     stringPtr("test-value"),
		KerberosRealm:                  "test-value",
		KerberosKdc:                    "test-value",
		KerberosKdcHttpsProxy:          "test-value",
		KerberosExchangeSpn:            "test-value",
		KerberosEnableTls:              true,
		KerberosAuthEveryRequest:       true,
		KerberosVerifyTlsUsingCustomCa: true,
		AddinServerDomain:              "test-value",
		AddinDisplayName:               "ms_exchange_connector-test",
		AddinDescription:               "Test MsExchangeConnector",
		AddinProviderName:              "ms_exchange_connector-test",
		AddinButtonLabel:               "test-value",
		AddinGroupLabel:                "test-value",
		AddinSupertipTitle:             "test-value",
		AddinSupertipDescription:       "Test MsExchangeConnector",
		AddinApplicationID:             stringPtr("test-value"),
		AddinAuthorityURL:              "https://example.com",
		AddinOidcMetadataURL:           "https://example.com",
		AddinAuthenticationMethod:      "EXCHANGE_USER_ID_TOKEN",
		AddinNaaWebApiApplicationID:    stringPtr("test-value"),
		PersonalVmrOauthClientID:       stringPtr("test-value"),
		PersonalVmrOauthClientSecret:   "test-value",
		PersonalVmrOauthAuthEndpoint:   "test-value",
		PersonalVmrOauthTokenEndpoint:  "test-value",
		PersonalVmrAdfsRelyingPartyTrustIdentifier: "test-value",
		OfficeJsURL:                  "https://example.com",
		MicrosoftFabricURL:           "https://example.com",
		MicrosoftFabricComponentsURL: "https://example.com",
		AdditionalAddInScriptSources: "test-value",
		Domains:                      nil, // Not specified in test config, so should be nil
		HostIdentityProviderGroup:    stringPtr("test-server.example.com"),
		IvrTheme:                  stringPtr("test-value"),
		NonIdpParticipants:        "disallow_all",
		// Template fields with defaults
		AcceptEditedOccurrenceTemplate:      "<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\nThis meeting occurrence in a recurring series has been successfully rescheduled using the aliases: {{alias}} and {{numeric_alias}}.<br>\r\n</div>",
		AcceptEditedRecurringSeriesTemplate: "<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\nThis recurring meeting series has been successfully rescheduled.<br>\r\nAll meetings in this series will use the aliases: {{alias}} and {{numeric_alias}}.<br>\r\n</div>",
		AcceptEditedSingleMeetingTemplate:   "<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\nThis meeting has been successfully rescheduled using the aliases: {{alias}} and {{numeric_alias}}.<br>\r\n</div>",
		AcceptNewRecurringSeriesTemplate:    "<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\nThis recurring meeting series has been successfully scheduled.<br>\r\nAll meetings in this series will use the aliases: {{alias}} and {{numeric_alias}}.<br>\r\n</div>",
		AcceptNewSingleMeetingTemplate:      "<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\nThis meeting has been successfully scheduled using the aliases: {{alias}} and {{numeric_alias}}.<br>\r\n</div>",
		ConferenceDescriptionTemplate:       "Scheduled Conference booked by {{organizer_email}}",
		ConferenceNameTemplate:              "{{subject}} ({{organizer_name}})",
		ConferenceSubjectTemplate:           "{{subject}}",
		MeetingInstructionsTemplate:         "<br>\r\n<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\n<b>Please join my Pexip Virtual Meeting Room in one of the following ways:</b><br>\r\n<br>\r\nFrom a VC endpoint or a Skype/Lync client:<br>\r\n{{alias}}<br>\r\n<br>\r\nFrom a web browser:<br>\r\n<a href=\"https://{{addin_server_domain}}/webapp/#/?conference={{alias}}\">https://{{addin_server_domain}}/webapp/#/?conference={{alias}}</a><br>\r\n<br>\r\nFrom a Pexip Infinity Connect client:<br>\r\npexip://{{alias}}<br>\r\n<br>\r\nFrom a telephone:<br>\r\n[Your number], then {{numeric_alias}} #<br>\r\n<br>\r\n{{alias_uuid}}<br>\r\n</div>",
		PersonalVmrDescriptionTemplate:      "{{description}}",
		PersonalVmrInstructionsTemplate:     "{% if domain_aliases %}\r\n    {% set alias = domain_aliases[0] %}\r\n{% elif other_aliases %}\r\n    {% set alias = other_aliases[0] %}\r\n{% else %}\r\n    {% set alias = numeric_aliases[0] %}\r\n{% endif %}\r\n{% if (not allow_guests) and pin %}\r\n    {% set meeting_pin = pin %}\r\n{% elif allow_guests and guest_pin %}\r\n    {% set meeting_pin = guest_pin %}\r\n{% else %}\r\n    {% set meeting_pin = \"\" %}\r\n{% endif %}\r\n<br>\r\n<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\n<b>Please join my Pexip Virtual Meeting Room in one of the following ways:</b><br>\r\n<br>\r\nFrom a VC endpoint or a Skype/Lync client:<br>\r\n{{alias}}<br>\r\n<br>\r\nFrom a web browser:<br>\r\n<a href=\"https://{{addin_server_domain}}/webapp/#/?conference={{alias}}\">https://{{addin_server_domain}}/webapp/#/?conference={{alias}}</a><br>\r\n<br>\r\nFrom a Pexip Infinity Connect client:<br>\r\npexip://{{alias}}<br>\r\n<br>\r\n{% if numeric_aliases %}\r\nFrom a telephone:<br>\r\n[Your number], then {{numeric_aliases[0]}} #<br>\r\n<br>\r\n{% endif %}\r\n{% if meeting_pin %}\r\nPlease join using the PIN <b>{{meeting_pin}}</b><br>\r\n<br>\r\n{% endif %}\r\n</div>",
		PersonalVmrLocationTemplate:         "{% if domain_aliases %}\r\n    {% set alias = domain_aliases[0] %}\r\n{% elif other_aliases %}\r\n    {% set alias = other_aliases[0] %}\r\n{% else %}\r\n    {% set alias = numeric_aliases[0] %}\r\n{% endif %}\r\nhttps://{{addin_server_domain}}/webapp/#/?conference={{alias}}",
		PersonalVmrNameTemplate:             "{{name}}",
		PlaceholderInstructionsTemplate:     "<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\nThis meeting will be hosted in a Virtual Meeting Room. Joining instructions will be<br>\r\nsent to you soon in a separate email.<br>\r\n</div>",
		RejectAliasConflictTemplate:         "<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\nWe are unable to schedule this meeting because the alias: {{alias}} is already <br>\r\nin use by another Pexip Virtual Meeting Room. Please try creating a new meeting.<br>\r\n</div>",
		RejectAliasDeletedTemplate:          "<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\nWe are unable to schedule this meeting because its alias has been deleted.<br>\r\nPlease try creating a new meeting.<br>\r\n</div>",
		RejectGeneralErrorTemplate:          "<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\nWe are unable to schedule this meeting. Please try creating a new meeting.<br>\r\nIf this issue continues, please forward this message to your system administrator, including the following ID:<br>\r\nCorrelationID=\"{{correlation_id}}\".<br>\r\n</div>",
		RejectInvalidAliasIDTemplate:        "<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\nThis meeting request does not contain currently valid scheduling data, and therefore cannot be processed.<br>\r\nPlease use the add-in to create a new meeting request, without editing any of the content that is inserted by the add-in.<br>\r\nIf this issue continues, please contact your system administrator.<br>\r\n</div>",
		RejectRecurringSeriesPastTemplate:   "<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\nThis recurring series cannot be scheduled because all<br>\r\noccurrences happen in the past.<br>\r\n</div>",
		RejectSingleMeetingPast:             "<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\nThis meeting cannot be scheduled because it occurs in the past.<br>\r\n</div>",
		ScheduledAliasDescriptionTemplate:   "Scheduled Conference booked by {{organizer_email}}",
		// Add-in pane fields with defaults
		AddinPaneTitle:                                   "Add a VMR",
		AddinPaneDescription:                             "This assigns a Virtual Meeting Room for your meeting",
		AddinPaneButtonTitle:                             "Add a Single-use VMR",
		AddinPaneSuccessHeading:                          "Success",
		AddinPaneSuccessMessage:                          "This meeting is now set up to be hosted as an audio or video conference in a Virtual Meeting Room. Please note this conference is not scheduled until you select Send.",
		AddinPaneAlreadyVideoMeetingHeading:              "VMR already assigned",
		AddinPaneAlreadyVideoMeetingMessage:              "It looks like this meeting has already been set up to be hosted in a Virtual Meeting Room. If this is a new meeting, select Send to schedule the conference.",
		AddinPaneGeneralErrorHeading:                     "Error",
		AddinPaneGeneralErrorMessage:                     "There was a problem adding the joining instructions. Please try again.",
		AddinPaneManagementNodeDownHeading:               "Cannot assign a VMR right now",
		AddinPaneManagementNodeDownMessage:               "Sorry, we are unable to assign a Virtual Meeting Room at this time. Select Send to schedule the meeting, and all attendees will be sent joining instructions later.",
		AddinPanePersonalVmrAddButton:                    "Add a Personal VMR",
		AddinPanePersonalVmrSignInButton:                 "Sign In",
		AddinPanePersonalVmrSelectMessage:                "Select the VMR you want to add to the meeting",
		AddinPanePersonalVmrNoneMessage:                  "You do not have any personal VMRs",
		AddinPanePersonalVmrErrorGettingMessage:          "There was a problem getting your personal VMRs. Please try again.",
		AddinPanePersonalVmrErrorSigningInMessage:        "There was a problem signing you in. Please try again.",
		AddinPanePersonalVmrErrorInsertingMeetingMessage: "There was a problem adding the joining instructions. Please try again.",
	}

	// Mock the GetMsexchangeconnector API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/ms_exchange_connector/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		ms_exchange_connector := args.Get(3).(*config.MsExchangeConnector)
		*ms_exchange_connector = *mockState
	}).Maybe()

	// Mock the UpdateMsexchangeconnector API call
	client.On("PutJSON", mock.Anything, "configuration/v1/ms_exchange_connector/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.MsExchangeConnectorUpdateRequest)
		ms_exchange_connector := args.Get(3).(*config.MsExchangeConnector)

		// Update mock state based on request
		mockState.Name = updateRequest.Name
		mockState.Description = updateRequest.Description

		// Handle oauth_state specifically since it might not be in update request
		if updateRequest.Description == "Updated Test MsExchangeConnector" {
			mockState.OauthState = stringPtr("updated-value")
		}
		if updateRequest.MeetingBufferBefore != nil {
			mockState.MeetingBufferBefore = *updateRequest.MeetingBufferBefore
		}
		if updateRequest.MeetingBufferAfter != nil {
			mockState.MeetingBufferAfter = *updateRequest.MeetingBufferAfter
		}
		if updateRequest.ScheduledAliasSuffixLength != nil {
			mockState.ScheduledAliasSuffixLength = *updateRequest.ScheduledAliasSuffixLength
		}
		if updateRequest.RoomMailboxEmailAddress != nil {
			mockState.RoomMailboxEmailAddress = updateRequest.RoomMailboxEmailAddress
		}
		mockState.RoomMailboxName = updateRequest.RoomMailboxName
		mockState.URL = updateRequest.URL
		mockState.Username = updateRequest.Username
		mockState.Password = updateRequest.Password
		mockState.AuthenticationMethod = updateRequest.AuthenticationMethod
		mockState.AuthProvider = updateRequest.AuthProvider
		mockState.UUID = updateRequest.UUID
		if updateRequest.ScheduledAliasPrefix != nil {
			mockState.ScheduledAliasPrefix = updateRequest.ScheduledAliasPrefix
		}
		mockState.ScheduledAliasDomain = updateRequest.ScheduledAliasDomain
		if updateRequest.EnableDynamicVmrs != nil {
			mockState.EnableDynamicVmrs = *updateRequest.EnableDynamicVmrs
		}
		if updateRequest.EnablePersonalVmrs != nil {
			mockState.EnablePersonalVmrs = *updateRequest.EnablePersonalVmrs
		}
		if updateRequest.AllowNewUsers != nil {
			mockState.AllowNewUsers = *updateRequest.AllowNewUsers
		}
		if updateRequest.DisableProxy != nil {
			mockState.DisableProxy = *updateRequest.DisableProxy
		}
		if updateRequest.UseCustomAddInSources != nil {
			mockState.UseCustomAddInSources = *updateRequest.UseCustomAddInSources
		}
		if updateRequest.EnableAddinDebugLogs != nil {
			mockState.EnableAddinDebugLogs = *updateRequest.EnableAddinDebugLogs
		}
		if updateRequest.OauthClientID != nil {
			mockState.OauthClientID = updateRequest.OauthClientID
		}
		mockState.OauthClientSecret = updateRequest.OauthClientSecret
		mockState.OauthAuthEndpoint = updateRequest.OauthAuthEndpoint
		mockState.OauthTokenEndpoint = updateRequest.OauthTokenEndpoint
		mockState.OauthRedirectURI = updateRequest.OauthRedirectURI
		mockState.OauthRefreshToken = updateRequest.OauthRefreshToken
		// OauthState may not be available in update request, keep current value
		mockState.KerberosRealm = updateRequest.KerberosRealm
		mockState.KerberosKdc = updateRequest.KerberosKdc
		mockState.KerberosKdcHttpsProxy = updateRequest.KerberosKdcHttpsProxy
		mockState.KerberosExchangeSpn = updateRequest.KerberosExchangeSpn
		if updateRequest.KerberosEnableTls != nil {
			mockState.KerberosEnableTls = *updateRequest.KerberosEnableTls
		}
		if updateRequest.KerberosAuthEveryRequest != nil {
			mockState.KerberosAuthEveryRequest = *updateRequest.KerberosAuthEveryRequest
		}
		if updateRequest.KerberosVerifyTlsUsingCustomCa != nil {
			mockState.KerberosVerifyTlsUsingCustomCa = *updateRequest.KerberosVerifyTlsUsingCustomCa
		}
		// Update other fields as needed
		mockState.AddinServerDomain = updateRequest.AddinServerDomain
		mockState.AddinDisplayName = updateRequest.AddinDisplayName
		mockState.AddinDescription = updateRequest.AddinDescription
		mockState.AddinProviderName = updateRequest.AddinProviderName
		mockState.AddinButtonLabel = updateRequest.AddinButtonLabel
		mockState.AddinGroupLabel = updateRequest.AddinGroupLabel
		mockState.AddinSupertipTitle = updateRequest.AddinSupertipTitle
		mockState.AddinSupertipDescription = updateRequest.AddinSupertipDescription
		if updateRequest.AddinApplicationID != nil {
			mockState.AddinApplicationID = updateRequest.AddinApplicationID
		}
		mockState.AddinAuthorityURL = updateRequest.AddinAuthorityURL
		mockState.AddinOidcMetadataURL = updateRequest.AddinOidcMetadataURL
		mockState.AddinAuthenticationMethod = updateRequest.AddinAuthenticationMethod
		if updateRequest.AddinNaaWebApiApplicationID != nil {
			mockState.AddinNaaWebApiApplicationID = updateRequest.AddinNaaWebApiApplicationID
		}
		if updateRequest.PersonalVmrOauthClientID != nil {
			mockState.PersonalVmrOauthClientID = updateRequest.PersonalVmrOauthClientID
		}
		// Update additional fields
		// Note: Domains field in update request is *[]string (URIs), but in response it's *[]ExchangeDomain (objects)
		// The test doesn't verify this field, so we skip updating it in the mock
		if updateRequest.HostIdentityProviderGroup != nil {
			mockState.HostIdentityProviderGroup = updateRequest.HostIdentityProviderGroup
		}
		if updateRequest.IvrTheme != nil {
			mockState.IvrTheme = updateRequest.IvrTheme
		}
		mockState.PersonalVmrOauthClientSecret = updateRequest.PersonalVmrOauthClientSecret
		mockState.PersonalVmrOauthAuthEndpoint = updateRequest.PersonalVmrOauthAuthEndpoint
		mockState.PersonalVmrOauthTokenEndpoint = updateRequest.PersonalVmrOauthTokenEndpoint
		mockState.PersonalVmrAdfsRelyingPartyTrustIdentifier = updateRequest.PersonalVmrAdfsRelyingPartyTrustIdentifier
		mockState.OfficeJsURL = updateRequest.OfficeJsURL
		mockState.MicrosoftFabricURL = updateRequest.MicrosoftFabricURL
		mockState.MicrosoftFabricComponentsURL = updateRequest.MicrosoftFabricComponentsURL
		mockState.AdditionalAddInScriptSources = updateRequest.AdditionalAddInScriptSources
		mockState.NonIdpParticipants = updateRequest.NonIdpParticipants

		// Return updated state
		*ms_exchange_connector = *mockState
	}).Maybe()

	// Mock the DeleteMsexchangeconnector API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/ms_exchange_connector/123/"
	}), mock.Anything).Return(nil)

	testInfinityMsExchangeConnector(t, client)
}

func testInfinityMsExchangeConnector(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ms_exchange_connector_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "name", "ms_exchange_connector-test"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "description", "Test MsExchangeConnector"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "room_mailbox_name", "ms_exchange_connector-test"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "username", "ms_exchange_connector-test"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "enable_dynamic_vmrs", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "enable_personal_vmrs", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "allow_new_users", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "disable_proxy", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "use_custom_add_in_sources", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "enable_addin_debug_logs", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "kerberos_enable_tls", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "kerberos_auth_every_request", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "kerberos_verify_tls_using_custom_ca", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "addin_display_name", "ms_exchange_connector-test"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "addin_description", "Test MsExchangeConnector"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "addin_provider_name", "ms_exchange_connector-test"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "addin_supertip_description", "Test MsExchangeConnector"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ms_exchange_connector_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "name", "ms_exchange_connector-test"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "description", "Updated Test MsExchangeConnector"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "meeting_buffer_before", "600"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "meeting_buffer_after", "600"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "scheduled_alias_suffix_length", "8"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "room_mailbox_email_address", "updated@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "url", "https://updated.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "authentication_method", "BASIC"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "auth_provider", "ADFS"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "enable_dynamic_vmrs", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "enable_personal_vmrs", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "allow_new_users", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "disable_proxy", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "use_custom_add_in_sources", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "enable_addin_debug_logs", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "kerberos_enable_tls", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "kerberos_auth_every_request", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "kerberos_verify_tls_using_custom_ca", "false"),
				),
			},
		},
	})
}
