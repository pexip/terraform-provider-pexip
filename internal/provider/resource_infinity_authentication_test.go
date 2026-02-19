/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"context"
	"os"
	"testing"

	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityAuthentication_URLPattern(t *testing.T) {
	t.Parallel()
	// This test specifically verifies that the URL pattern is correct
	client := infinity.NewClientMock()

	// Test that GetAuthentication calls the correct URL
	client.On("GetJSON", mock.Anything, "configuration/v1/authentication/1/", mock.Anything, mock.Anything).Return(nil).Once()

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
	client.On("PatchJSON", mock.Anything, "configuration/v1/authentication/1/", mock.Anything, mock.Anything).Return(nil).Once()

	// Call the method to verify URL pattern
	_, _ = client.Config().UpdateAuthentication(context.TODO(), &config.AuthenticationUpdateRequest{})

	// Verify the expectations were met
	client.AssertExpectations(t)
}

func TestInfinityAuthentication(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	client := infinity.NewClientMock()

	// Initialize mockState with default values (min config)
	mockState := &config.Authentication{
		ResourceURI:               "/api/admin/configuration/v1/authentication/1/",
		Source:                    "LOCAL",
		ClientCertificate:         "NO",
		ApiOauth2DisableBasic:     false,
		ApiOauth2AllowAllPerms:    false,
		ApiOauth2Expiration:       3600,
		LdapServer:                "",
		LdapBaseDN:                "",
		LdapBindUsername:          "",
		LdapBindPassword:          "",
		LdapUserSearchDN:          "",
		LdapUserFilter:            "(&(objectclass=person)(!(objectclass=computer)))",
		LdapUserSearchFilter:      "(|(uid={username})(sAMAccountName={username}))",
		LdapUserGroupAttributes:   "memberOf",
		LdapGroupSearchDN:         "",
		LdapGroupFilter:           "(|(objectclass=group)(objectclass=groupOfNames)(objectclass=groupOfUniqueNames)(objectclass=posixGroup))",
		LdapGroupMembershipFilter: "(|(member={userdn})(uniquemember={userdn})(memberuid={useruid}))",
		LdapUseGlobalCatalog:      false,
		LdapPermitNoTLS:           false,
		OidcMetadataURL:           "",
		OidcMetadata:              "",
		OidcClientID:              "",
		OidcClientSecret:          "",
		OidcPrivateKey:            "",
		OidcAuthMethod:            "client_secret",
		OidcScope:                 "openid profile email",
		OidcAuthorizeURL:          "",
		OidcTokenEndpointURL:      "",
		OidcUsernameField:         "preferred_username",
		OidcGroupsField:           "groups",
		OidcRequiredKey:           "",
		OidcRequiredValue:         "",
		OidcDomainHint:            "",
		OidcLoginButton:           "",
	}

	// Mock GetJSON to return current mockState
	client.On("GetJSON", mock.Anything, "configuration/v1/authentication/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		auth := args.Get(3).(*config.Authentication)
		*auth = *mockState
	}).Maybe()

	// Mock PatchJSON to update mockState
	client.On("PatchJSON", mock.Anything, "configuration/v1/authentication/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.AuthenticationUpdateRequest)
		auth := args.Get(3).(*config.Authentication)

		// Update mockState with values from the update request
		mockState.Source = updateRequest.Source
		mockState.ClientCertificate = updateRequest.ClientCertificate

		if updateRequest.ApiOauth2DisableBasic != nil {
			mockState.ApiOauth2DisableBasic = *updateRequest.ApiOauth2DisableBasic
		}
		if updateRequest.ApiOauth2AllowAllPerms != nil {
			mockState.ApiOauth2AllowAllPerms = *updateRequest.ApiOauth2AllowAllPerms
		}
		if updateRequest.ApiOauth2Expiration != nil {
			mockState.ApiOauth2Expiration = *updateRequest.ApiOauth2Expiration
		}

		// LDAP fields
		mockState.LdapServer = updateRequest.LdapServer
		mockState.LdapBaseDN = updateRequest.LdapBaseDN
		mockState.LdapBindUsername = updateRequest.LdapBindUsername
		// Don't update password from request - keep empty to simulate hashing
		mockState.LdapUserSearchDN = updateRequest.LdapUserSearchDN
		mockState.LdapUserFilter = updateRequest.LdapUserFilter
		mockState.LdapUserSearchFilter = updateRequest.LdapUserSearchFilter
		mockState.LdapUserGroupAttributes = updateRequest.LdapUserGroupAttributes
		mockState.LdapGroupSearchDN = updateRequest.LdapGroupSearchDN
		mockState.LdapGroupFilter = updateRequest.LdapGroupFilter
		mockState.LdapGroupMembershipFilter = updateRequest.LdapGroupMembershipFilter

		if updateRequest.LdapUseGlobalCatalog != nil {
			mockState.LdapUseGlobalCatalog = *updateRequest.LdapUseGlobalCatalog
		}
		if updateRequest.LdapPermitNoTLS != nil {
			mockState.LdapPermitNoTLS = *updateRequest.LdapPermitNoTLS
		}

		// OIDC fields
		mockState.OidcMetadataURL = updateRequest.OidcMetadataURL
		mockState.OidcMetadata = updateRequest.OidcMetadata
		mockState.OidcClientID = updateRequest.OidcClientID
		// Don't update secret from request - keep empty to simulate hashing
		// Hash the private key to simulate the API behavior
		if updateRequest.OidcPrivateKey != "" {
			mockState.OidcPrivateKey = "hashed_private_key_value"
		} else {
			mockState.OidcPrivateKey = ""
		}
		mockState.OidcAuthMethod = updateRequest.OidcAuthMethod
		mockState.OidcScope = updateRequest.OidcScope
		mockState.OidcAuthorizeURL = updateRequest.OidcAuthorizeURL
		mockState.OidcTokenEndpointURL = updateRequest.OidcTokenEndpointURL
		mockState.OidcUsernameField = updateRequest.OidcUsernameField
		mockState.OidcGroupsField = updateRequest.OidcGroupsField
		mockState.OidcRequiredKey = updateRequest.OidcRequiredKey
		mockState.OidcRequiredValue = updateRequest.OidcRequiredValue
		mockState.OidcDomainHint = updateRequest.OidcDomainHint
		mockState.OidcLoginButton = updateRequest.OidcLoginButton

		*auth = *mockState
	}).Maybe()

	testInfinityAuthentication(t, client)
}

func testInfinityAuthentication(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				// Step 1: Create with full config
				Config: test.LoadTestFolder(t, "resource_infinity_authentication_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_authentication.authentication-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "source", "LOCAL"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "client_certificate", "NO"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "api_oauth2_allow_all_perms", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "api_oauth2_expiration", "7200"),
					// LDAP configuration
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_server", "ldap.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_base_dn", "dc=example,dc=com"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_bind_username", "CN=Service Account,OU=Service Accounts,DC=example,DC=com"),
					resource.TestCheckResourceAttrSet("pexip_infinity_authentication.authentication-test", "ldap_bind_password"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_user_search_dn", "OU=Users,DC=example,DC=com"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_user_filter", "(&(objectclass=person))"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_user_search_filter", "(|(uid={username}))"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_user_group_attributes", "memberOftest"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_group_search_dn", "OU=Groups,DC=example,DC=com"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_group_filter", "(|(objectclass=group)(objectclass=groupOfNames)(objectclass=groupOfUniqueNames)(objectclass=posixGroup)(objectclass=testGroup))"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_group_membership_filter", "(|(member={userdn})(uniquemember={userdn}))"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_use_global_catalog", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_permit_no_tls", "true"),
					// OIDC configuration
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_metadata_url", "https://auth.example.com/.well-known/openid-configuration"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_client_id", "pexip-infinity-client-id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_authentication.authentication-test", "oidc_client_secret"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_auth_method", "private_key"),
					resource.TestCheckResourceAttrSet("pexip_infinity_authentication.authentication-test", "oidc_private_key"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_scope", "openid profile email groups"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_authorize_url", "https://auth.example.com/oauth2/authorize"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_token_endpoint_url", "https://auth.example.com/oauth2/token"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_username_field", "preferred_username_test"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_groups_field", "testgroups"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_required_key", "department"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_required_value", "IT"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_domain_hint", "example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_login_button", "Sign in with Corporate SSO"),
				),
			},
			{
				// Step 2: Update to min config
				Config: test.LoadTestFolder(t, "resource_infinity_authentication_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_authentication.authentication-test", "id"),
					// Verify defaults are restored
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "source", "LOCAL"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "client_certificate", "NO"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "api_oauth2_disable_basic", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "api_oauth2_allow_all_perms", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "api_oauth2_expiration", "3600"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_server", ""),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_base_dn", ""),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_bind_username", ""),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_metadata_url", ""),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_client_id", ""),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_auth_method", "client_secret"),
				),
			},
			{
				// Step 3: Destroy (no-op for singleton, but included for consistency)
				Config:  test.LoadTestFolder(t, "resource_infinity_authentication_min"),
				Destroy: true,
			},
			{
				// Step 4: Recreate with min config (actually just updates since it's a singleton)
				Config: test.LoadTestFolder(t, "resource_infinity_authentication_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_authentication.authentication-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "source", "LOCAL"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "client_certificate", "NO"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "api_oauth2_expiration", "3600"),
				),
			},
			{
				// Step 5: Update to full config
				Config: test.LoadTestFolder(t, "resource_infinity_authentication_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_authentication.authentication-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "source", "LOCAL"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "client_certificate", "NO"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "api_oauth2_allow_all_perms", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "api_oauth2_expiration", "7200"),
					// LDAP configuration
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_server", "ldap.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_base_dn", "dc=example,dc=com"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_bind_username", "CN=Service Account,OU=Service Accounts,DC=example,DC=com"),
					resource.TestCheckResourceAttrSet("pexip_infinity_authentication.authentication-test", "ldap_bind_password"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_user_search_dn", "OU=Users,DC=example,DC=com"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_user_filter", "(&(objectclass=person))"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_user_search_filter", "(|(uid={username}))"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_user_group_attributes", "memberOftest"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_group_search_dn", "OU=Groups,DC=example,DC=com"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_group_filter", "(|(objectclass=group)(objectclass=groupOfNames)(objectclass=groupOfUniqueNames)(objectclass=posixGroup)(objectclass=testGroup))"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_group_membership_filter", "(|(member={userdn})(uniquemember={userdn}))"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_use_global_catalog", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_permit_no_tls", "true"),
					// OIDC configuration
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_metadata_url", "https://auth.example.com/.well-known/openid-configuration"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_client_id", "pexip-infinity-client-id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_authentication.authentication-test", "oidc_client_secret"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_auth_method", "private_key"),
					resource.TestCheckResourceAttrSet("pexip_infinity_authentication.authentication-test", "oidc_private_key"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_scope", "openid profile email groups"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_authorize_url", "https://auth.example.com/oauth2/authorize"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_token_endpoint_url", "https://auth.example.com/oauth2/token"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_username_field", "preferred_username_test"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_groups_field", "testgroups"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_required_key", "department"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_required_value", "IT"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_domain_hint", "example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_login_button", "Sign in with Corporate SSO"),
				),
			},
		},
	})
}
