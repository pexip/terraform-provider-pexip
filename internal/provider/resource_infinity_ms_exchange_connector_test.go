package provider

import (
	"os"
	"testing"

	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityMsExchangeConnector(t *testing.T) {
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
		AuthenticationMethod:           "oauth2",
		AuthProvider:                   "azure",
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
		AddinAuthenticationMethod:      "web_api",
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
		Domains:                      stringPtr("test-value"),
		HostIdentityProviderGroup:    stringPtr("test-server.example.com"),
		IvrTheme:                     stringPtr("test-value"),
		NonIdpParticipants:           "test-value",
	}

	// Mock the GetMsexchangeconnector API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/ms_exchange_connector/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		ms_exchange_connector := args.Get(2).(*config.MsExchangeConnector)
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
		if updateRequest.Domains != nil {
			mockState.Domains = updateRequest.Domains
		}
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
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "authentication_method", "basic"),
					resource.TestCheckResourceAttr("pexip_infinity_ms_exchange_connector.ms_exchange_connector-test", "auth_provider", "exchange"),
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
