package provider

import (
	"context"
	"os"
	"testing"

	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityAuthentication_URLPattern(t *testing.T) {
	t.Parallel()
	// This test specifically verifies that the URL pattern is correct
	client := infinity.NewClientMock()

	// Test that GetAuthentication calls the correct URL
	client.On("GetJSON", mock.Anything, "configuration/v1/authentication/1/", mock.Anything).Return(nil).Once()

	// Call the method to verify URL pattern
	_, _ = client.Config().GetAuthentication(context.TODO())

	// Verify the expectations were met
	client.AssertExpectations(t)
}

func TestInfinityAuthentication_UpdateURLPattern(t *testing.T) {
	t.Parallel()
	// This test specifically verifies that the UpdateAuthentication URL pattern is correct
	client := infinity.NewClientMock()

	// Test that UpdateAuthentication calls the correct URL
	client.On("PutJSON", mock.Anything, "configuration/v1/authentication/1/", mock.Anything, mock.Anything).Return(nil).Once()

	// Call the method to verify URL pattern
	_, _ = client.Config().UpdateAuthentication(context.TODO(), &config.AuthenticationUpdateRequest{})

	// Verify the expectations were met
	client.AssertExpectations(t)
}

func TestInfinityAuthentication(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	client := infinity.NewClientMock()

	mockState := &config.Authentication{
		Source:            "local",
		ClientCertificate: "disabled",
		ResourceURI:       "/configuration/v1/authentication/1/",
		// Don't set ApiOauth2Expiration - let it default to 0 but handle it properly in the resource
	}

	client.On("GetJSON", mock.Anything, "configuration/v1/authentication/1/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		auth := args.Get(2).(*config.Authentication)
		*auth = *mockState
	}).Maybe()

	client.On("PutJSON", mock.Anything, "configuration/v1/authentication/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.AuthenticationUpdateRequest)
		auth := args.Get(3).(*config.Authentication)

		// Update mockState with all the values from the update request
		if updateRequest.Source != "" {
			mockState.Source = updateRequest.Source
		}
		if updateRequest.ClientCertificate != "" {
			mockState.ClientCertificate = updateRequest.ClientCertificate
		}
		if updateRequest.ApiOauth2DisableBasic != nil {
			mockState.ApiOauth2DisableBasic = *updateRequest.ApiOauth2DisableBasic
		}
		if updateRequest.ApiOauth2AllowAllPerms != nil {
			mockState.ApiOauth2AllowAllPerms = *updateRequest.ApiOauth2AllowAllPerms
		}
		if updateRequest.ApiOauth2Expiration != nil {
			mockState.ApiOauth2Expiration = *updateRequest.ApiOauth2Expiration
		}
		if updateRequest.LdapServer != "" {
			mockState.LdapServer = updateRequest.LdapServer
		}
		if updateRequest.LdapBaseDN != "" {
			mockState.LdapBaseDN = updateRequest.LdapBaseDN
		}
		if updateRequest.LdapBindUsername != "" {
			mockState.LdapBindUsername = updateRequest.LdapBindUsername
		}
		if updateRequest.LdapBindPassword != "" {
			mockState.LdapBindPassword = updateRequest.LdapBindPassword
		}
		if updateRequest.LdapUserSearchDN != "" {
			mockState.LdapUserSearchDN = updateRequest.LdapUserSearchDN
		}
		if updateRequest.LdapUserFilter != "" {
			mockState.LdapUserFilter = updateRequest.LdapUserFilter
		}
		if updateRequest.LdapUserSearchFilter != "" {
			mockState.LdapUserSearchFilter = updateRequest.LdapUserSearchFilter
		}
		if updateRequest.LdapUserGroupAttributes != "" {
			mockState.LdapUserGroupAttributes = updateRequest.LdapUserGroupAttributes
		}
		if updateRequest.LdapGroupSearchDN != "" {
			mockState.LdapGroupSearchDN = updateRequest.LdapGroupSearchDN
		}
		if updateRequest.LdapGroupFilter != "" {
			mockState.LdapGroupFilter = updateRequest.LdapGroupFilter
		}
		if updateRequest.LdapGroupMembershipFilter != "" {
			mockState.LdapGroupMembershipFilter = updateRequest.LdapGroupMembershipFilter
		}
		if updateRequest.LdapUseGlobalCatalog != nil {
			mockState.LdapUseGlobalCatalog = *updateRequest.LdapUseGlobalCatalog
		}
		if updateRequest.LdapPermitNoTLS != nil {
			mockState.LdapPermitNoTLS = *updateRequest.LdapPermitNoTLS
		}
		if updateRequest.OidcMetadataURL != "" {
			mockState.OidcMetadataURL = updateRequest.OidcMetadataURL
		}
		if updateRequest.OidcMetadata != "" {
			mockState.OidcMetadata = updateRequest.OidcMetadata
		}
		if updateRequest.OidcClientID != "" {
			mockState.OidcClientID = updateRequest.OidcClientID
		}
		if updateRequest.OidcClientSecret != "" {
			mockState.OidcClientSecret = updateRequest.OidcClientSecret
		}
		if updateRequest.OidcPrivateKey != "" {
			mockState.OidcPrivateKey = updateRequest.OidcPrivateKey
		}
		if updateRequest.OidcAuthMethod != "" {
			mockState.OidcAuthMethod = updateRequest.OidcAuthMethod
		}
		if updateRequest.OidcScope != "" {
			mockState.OidcScope = updateRequest.OidcScope
		}
		if updateRequest.OidcAuthorizeURL != "" {
			mockState.OidcAuthorizeURL = updateRequest.OidcAuthorizeURL
		}
		if updateRequest.OidcTokenEndpointURL != "" {
			mockState.OidcTokenEndpointURL = updateRequest.OidcTokenEndpointURL
		}
		if updateRequest.OidcUsernameField != "" {
			mockState.OidcUsernameField = updateRequest.OidcUsernameField
		}
		if updateRequest.OidcGroupsField != "" {
			mockState.OidcGroupsField = updateRequest.OidcGroupsField
		}
		if updateRequest.OidcRequiredKey != "" {
			mockState.OidcRequiredKey = updateRequest.OidcRequiredKey
		}
		if updateRequest.OidcRequiredValue != "" {
			mockState.OidcRequiredValue = updateRequest.OidcRequiredValue
		}
		if updateRequest.OidcDomainHint != "" {
			mockState.OidcDomainHint = updateRequest.OidcDomainHint
		}
		if updateRequest.OidcLoginButton != "" {
			mockState.OidcLoginButton = updateRequest.OidcLoginButton
		}

		*auth = *mockState
	}).Maybe()

	testInfinityAuthentication(t, client)
}

func testInfinityAuthentication(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_authentication_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_authentication.authentication-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "source", "local"),
				),
			},
		},
	})
}
