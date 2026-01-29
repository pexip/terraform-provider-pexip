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

	stringPtr := func(s string) *string { return &s }

	exchangeConnectorState := &config.MsExchangeConnector{
		ID:                             1,
		ResourceURI:                    "/api/admin/configuration/v1/ms_exchange_connector/1/",
		Name:                           "test-exchange-connector",
		Description:                    "Test Exchange Connector",
		MeetingBufferBefore:            300,
		MeetingBufferAfter:             300,
		ScheduledAliasSuffixLength:     6,
		RoomMailboxEmailAddress:        stringPtr("test@example.com"),
		RoomMailboxName:                "test-exchange-connector",
		URL:                            "https://exchange.test.local",
		Username:                       "testuser",
		Password:                       "testpass",
		AuthenticationMethod:           "oauth2",
		AuthProvider:                   "azure",
		UUID:                           "test-uuid-value",
		ScheduledAliasPrefix:           stringPtr("test"),
		ScheduledAliasDomain:           "example.com",
		EnableDynamicVmrs:              true,
		EnablePersonalVmrs:             true,
		AllowNewUsers:                  true,
		DisableProxy:                   true,
		UseCustomAddInSources:          true,
		EnableAddinDebugLogs:           true,
		OauthClientID:                  stringPtr("test-client-id"),
		OauthClientSecret:              "test-secret",
		OauthAuthEndpoint:              "test-auth-endpoint",
		OauthTokenEndpoint:             "test-token-endpoint",
		OauthRedirectURI:               "test-redirect-uri",
		OauthRefreshToken:              "test-refresh-token",
		OauthState:                     stringPtr("test-state"),
		KerberosRealm:                  "test-realm",
		KerberosKdc:                    "test-kdc",
		KerberosKdcHttpsProxy:          "test-proxy",
		KerberosExchangeSpn:            "test-spn",
		KerberosEnableTls:              true,
		KerberosAuthEveryRequest:       true,
		KerberosVerifyTlsUsingCustomCa: true,
		AddinServerDomain:              "test-domain",
		AddinDisplayName:               "test-exchange-connector",
		AddinDescription:               "Test Exchange Connector",
		AddinProviderName:              "test-exchange-connector",
		AddinButtonLabel:               "test-button",
		AddinGroupLabel:                "test-group",
		AddinSupertipTitle:             "test-title",
		AddinSupertipDescription:       "Test Exchange Connector",
		AddinApplicationID:             stringPtr("test-app-id"),
		AddinAuthorityURL:              "https://example.com",
		AddinOidcMetadataURL:           "https://example.com",
		AddinAuthenticationMethod:      "web_api",
		AddinNaaWebApiApplicationID:    stringPtr("test-naa-app-id"),
		PersonalVmrOauthClientID:       stringPtr("test-vmr-client-id"),
		PersonalVmrOauthClientSecret:   "test-vmr-secret",
		PersonalVmrOauthAuthEndpoint:   "test-vmr-auth",
		PersonalVmrOauthTokenEndpoint:  "test-vmr-token",
		PersonalVmrAdfsRelyingPartyTrustIdentifier: "test-adfs",
		OfficeJsURL:                  "https://example.com",
		MicrosoftFabricURL:           "https://example.com",
		MicrosoftFabricComponentsURL: "https://example.com",
		AdditionalAddInScriptSources: "test-sources",
		Domains: &[]config.ExchangeDomain{
			{
				ID:                1,
				Domain:            "test-domain",
				ExchangeConnector: "/api/admin/configuration/v1/ms_exchange_connector/1/",
				ResourceURI:       "/api/admin/configuration/v1/exchange_domain/1/",
			},
		},
		HostIdentityProviderGroup: nil,
		IvrTheme:                  nil,
		NonIdpParticipants:        "test-participants",
	}

	client.On("GetJSON", mock.Anything, "configuration/v1/ms_exchange_connector/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		connector := args.Get(3).(*config.MsExchangeConnector)
		*connector = *exchangeConnectorState
	}).Maybe()

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
