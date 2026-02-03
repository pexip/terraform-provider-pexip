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

func TestInfinityScheduledAlias(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock MS Exchange connector creation
	exchangeConnectorCreateResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/ms_exchange_connector/1/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/ms_exchange_connector/", mock.Anything, mock.Anything).Return(exchangeConnectorCreateResponse, nil)

	exchangeConnectorState := &config.MsExchangeConnector{
		ID:                             1,
		ResourceURI:                    "/api/admin/configuration/v1/ms_exchange_connector/1/",
		Name:                           "test-exchange-connector",
		Description:                    "Test Exchange Connector",
		AddinServerDomain:              "test-domain",
		UUID:                           "test-uuid-value",
		NonIdpParticipants:             "disallow_all",
		AuthenticationMethod:           "BASIC",
		AuthProvider:                   "ADFS",
		Username:                       "",
		Password:                       "",
		RoomMailboxName:                "",
		URL:                            "",
		ScheduledAliasDomain:           "",
		ScheduledAliasSuffixLength:     6,
		MeetingBufferBefore:            30,
		MeetingBufferAfter:             60,
		EnableDynamicVmrs:              true,
		EnablePersonalVmrs:             false,
		AllowNewUsers:                  true,
		DisableProxy:                   false,
		UseCustomAddInSources:          false,
		EnableAddinDebugLogs:           false,
		OauthClientSecret:              "",
		OauthAuthEndpoint:              "",
		OauthTokenEndpoint:             "",
		OauthRedirectURI:               "",
		OauthRefreshToken:              "",
		KerberosRealm:                  "",
		KerberosKdc:                    "",
		KerberosKdcHttpsProxy:          "",
		KerberosExchangeSpn:            "",
		KerberosEnableTls:              true,
		KerberosAuthEveryRequest:       false,
		KerberosVerifyTlsUsingCustomCa: false,
		AddinDisplayName:               "Pexip Scheduling Service",
		AddinDescription:               "Turns meetings into Pexip meetings",
		AddinProviderName:              "Pexip",
		AddinButtonLabel:               "Create a Pexip meeting",
		AddinGroupLabel:                "Pexip meeting",
		AddinSupertipTitle:             "Makes this a Pexip meeting",
		AddinSupertipDescription:       "Turns this meeting into an audio or video conference hosted in a Pexip VMR. The meeting is not scheduled until you select Send.",
		AddinAuthorityURL:              "",
		AddinOidcMetadataURL:           "",
		AddinAuthenticationMethod:      "EXCHANGE_USER_ID_TOKEN",
		PersonalVmrOauthClientID:       test.StringPtr("4189c2b4-92ca-416c-b7ea-bc3cfab3d0f0"),
		PersonalVmrOauthClientSecret:   "",
		PersonalVmrOauthAuthEndpoint:   "",
		PersonalVmrOauthTokenEndpoint:  "",
		PersonalVmrAdfsRelyingPartyTrustIdentifier: "",
		OfficeJsURL:                  "https://appsforoffice.microsoft.com/lib/1/hosted/office.js",
		MicrosoftFabricURL:           "https://appsforoffice.microsoft.com/fabric/1.0/fabric.min.css",
		MicrosoftFabricComponentsURL: "https://appsforoffice.microsoft.com/fabric/1.0/fabric.components.min.css",
		AdditionalAddInScriptSources: "",
		Domains:                      nil,
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

	client.On("GetJSON", mock.Anything, "configuration/v1/ms_exchange_connector/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		connector := args.Get(3).(*config.MsExchangeConnector)
		*connector = *exchangeConnectorState
	})

	client.On("DeleteJSON", mock.Anything, "configuration/v1/ms_exchange_connector/1/", mock.Anything).Return(nil).Maybe()

	// Mock the CreateScheduledalias API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/scheduled_alias/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/scheduled_alias/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.ScheduledAlias{
		ID:                123,
		ResourceURI:       "/api/admin/configuration/v1/scheduled_alias/123/",
		Alias:             "test-scheduled-alias",
		AliasNumber:       1234567890,
		NumericAlias:      "123456",
		UUID:              "11111111-1111-1111-1111-111111111111",
		ExchangeConnector: "/api/admin/configuration/v1/ms_exchange_connector/1/",
		IsUsed:            true,
		EWSItemUID:        test.StringPtr("test-ews-uid"),
	}

	// Mock the GetScheduledalias API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/scheduled_alias/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		scheduled_alias := args.Get(3).(*config.ScheduledAlias)
		*scheduled_alias = *mockState
	}).Maybe()

	// Mock the UpdateScheduledalias API call
	client.On("PutJSON", mock.Anything, "configuration/v1/scheduled_alias/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.ScheduledAliasUpdateRequest)
		scheduled_alias := args.Get(3).(*config.ScheduledAlias)

		// Update mock state based on request
		if updateRequest.Alias != "" {
			mockState.Alias = updateRequest.Alias
		}
		if updateRequest.AliasNumber != nil {
			mockState.AliasNumber = *updateRequest.AliasNumber
		}
		if updateRequest.NumericAlias != "" {
			mockState.NumericAlias = updateRequest.NumericAlias
		}
		if updateRequest.UUID != "" {
			mockState.UUID = updateRequest.UUID
		}
		if updateRequest.ExchangeConnector != "" {
			mockState.ExchangeConnector = updateRequest.ExchangeConnector
		}
		if updateRequest.IsUsed != nil {
			mockState.IsUsed = *updateRequest.IsUsed
		}
		if updateRequest.EWSItemUID != nil {
			mockState.EWSItemUID = updateRequest.EWSItemUID
		}

		// Return updated state
		*scheduled_alias = *mockState
	}).Maybe()

	// Mock the DeleteScheduledalias API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/scheduled_alias/123/"
	}), mock.Anything).Return(nil)

	testInfinityScheduledAlias(t, client)
}

func testInfinityScheduledAlias(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_scheduled_alias_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_scheduled_alias.scheduled_alias-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_scheduled_alias.scheduled_alias-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_alias.scheduled_alias-test", "is_used", "true"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_scheduled_alias_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_scheduled_alias.scheduled_alias-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_scheduled_alias.scheduled_alias-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_alias.scheduled_alias-test", "is_used", "false"),
				),
			},
		},
	})
}
