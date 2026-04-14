/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"context"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityAuthentication_URLPattern(t *testing.T) {
	t.Parallel()
	client := infinity.NewClientMock()

	client.On("GetJSON", mock.Anything, "configuration/v1/authentication/1/", mock.Anything, mock.Anything).Return(nil).Once()

	_, _ = client.Config().GetAuthentication(context.TODO())

	client.AssertExpectations(t)
}

func TestInfinityAuthentication_UpdateURLPattern(t *testing.T) {
	t.Parallel()
	client := infinity.NewClientMock()

	client.On("PatchJSON", mock.Anything, "configuration/v1/authentication/1/", mock.Anything, mock.Anything).Return(nil).Once()

	_, _ = client.Config().UpdateAuthentication(context.TODO(), &config.AuthenticationUpdateRequest{})

	client.AssertExpectations(t)
}

func TestInfinityAuthentication(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	client := infinity.NewClientMock()

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
		OidcScope:                 "openid email profile",
		OidcAuthorizeURL:          "",
		OidcTokenEndpointURL:      "",
		OidcUsernameField:         "preferred_username",
		OidcGroupsField:           "groups",
		OidcRequiredKey:           "",
		OidcRequiredValue:         "",
		OidcDomainHint:            "",
		OidcLoginButton:           "",
	}

	// Delete mock — registered first so it takes priority over the general mock.
	// Fingerprinted by LdapServer == "" && ApiOauth2Expiration == 3600, which the full
	// config does not match (it sets ldap_server and api_oauth2_expiration=7200).
	client.On("PatchJSON", mock.Anything, "configuration/v1/authentication/1/",
		mock.MatchedBy(func(req *config.AuthenticationUpdateRequest) bool {
			return req.LdapServer == "" && req.ApiOauth2Expiration != nil && *req.ApiOauth2Expiration == 3600
		}), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.AuthenticationUpdateRequest)

		assert.Equal(t, "LOCAL", req.Source)
		assert.Equal(t, "NO", req.ClientCertificate)
		assert.NotNil(t, req.ApiOauth2DisableBasic)
		assert.False(t, *req.ApiOauth2DisableBasic)
		assert.NotNil(t, req.ApiOauth2AllowAllPerms)
		assert.False(t, *req.ApiOauth2AllowAllPerms)
		assert.NotNil(t, req.ApiOauth2Expiration)
		assert.Equal(t, 3600, *req.ApiOauth2Expiration)
		assert.Equal(t, "", req.LdapServer)
		assert.Equal(t, "", req.LdapBaseDN)
		assert.Equal(t, "", req.LdapBindUsername)
		assert.Equal(t, "", req.LdapBindPassword)
		assert.Equal(t, "", req.LdapUserSearchDN)
		assert.Equal(t, "(&(objectclass=person)(!(objectclass=computer)))", req.LdapUserFilter)
		assert.Equal(t, "(|(uid={username})(sAMAccountName={username}))", req.LdapUserSearchFilter)
		assert.Equal(t, "memberOf", req.LdapUserGroupAttributes)
		assert.Equal(t, "", req.LdapGroupSearchDN)
		assert.Equal(t, "(|(objectclass=group)(objectclass=groupOfNames)(objectclass=groupOfUniqueNames)(objectclass=posixGroup))", req.LdapGroupFilter)
		assert.Equal(t, "(|(member={userdn})(uniquemember={userdn})(memberuid={useruid}))", req.LdapGroupMembershipFilter)
		assert.NotNil(t, req.LdapUseGlobalCatalog)
		assert.False(t, *req.LdapUseGlobalCatalog)
		assert.NotNil(t, req.LdapPermitNoTLS)
		assert.False(t, *req.LdapPermitNoTLS)
		assert.Equal(t, "", req.OidcMetadataURL)
		assert.Equal(t, "", req.OidcClientID)
		assert.Equal(t, "", req.OidcClientSecret)
		assert.Equal(t, "", req.OidcPrivateKey)
		assert.Equal(t, "client_secret", req.OidcAuthMethod)
		assert.Equal(t, "openid email profile", req.OidcScope)
		assert.Equal(t, "", req.OidcAuthorizeURL)
		assert.Equal(t, "", req.OidcTokenEndpointURL)
		assert.Equal(t, "preferred_username", req.OidcUsernameField)
		assert.Equal(t, "groups", req.OidcGroupsField)
		assert.Equal(t, "", req.OidcRequiredKey)
		assert.Equal(t, "", req.OidcRequiredValue)
		assert.Equal(t, "", req.OidcDomainHint)
		assert.Equal(t, "", req.OidcLoginButton)

		// Reset mockState to defaults
		mockState.Source = "LOCAL"
		mockState.ClientCertificate = "NO"
		mockState.ApiOauth2DisableBasic = false
		mockState.ApiOauth2AllowAllPerms = false
		mockState.ApiOauth2Expiration = 3600
		mockState.LdapServer = ""
		mockState.LdapBaseDN = ""
		mockState.LdapBindUsername = ""
		mockState.LdapBindPassword = ""
		mockState.LdapUserSearchDN = ""
		mockState.LdapUserFilter = "(&(objectclass=person)(!(objectclass=computer)))"
		mockState.LdapUserSearchFilter = "(|(uid={username})(sAMAccountName={username}))"
		mockState.LdapUserGroupAttributes = "memberOf"
		mockState.LdapGroupSearchDN = ""
		mockState.LdapGroupFilter = "(|(objectclass=group)(objectclass=groupOfNames)(objectclass=groupOfUniqueNames)(objectclass=posixGroup))"
		mockState.LdapGroupMembershipFilter = "(|(member={userdn})(uniquemember={userdn})(memberuid={useruid}))"
		mockState.LdapUseGlobalCatalog = false
		mockState.LdapPermitNoTLS = false
		mockState.OidcMetadataURL = ""
		mockState.OidcMetadata = ""
		mockState.OidcClientID = ""
		mockState.OidcClientSecret = ""
		mockState.OidcPrivateKey = ""
		mockState.OidcAuthMethod = "client_secret"
		mockState.OidcScope = "openid email profile"
		mockState.OidcAuthorizeURL = ""
		mockState.OidcTokenEndpointURL = ""
		mockState.OidcUsernameField = "preferred_username"
		mockState.OidcGroupsField = "groups"
		mockState.OidcRequiredKey = ""
		mockState.OidcRequiredValue = ""
		mockState.OidcDomainHint = ""
		mockState.OidcLoginButton = ""
	}).Once()

	// General PatchJSON mock — handles all create and update calls.
	client.On("PatchJSON", mock.Anything, "configuration/v1/authentication/1/",
		mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.AuthenticationUpdateRequest)
		auth := args.Get(3).(*config.Authentication)

		mockState.Source = req.Source
		mockState.ClientCertificate = req.ClientCertificate
		mockState.ApiOauth2DisableBasic = *req.ApiOauth2DisableBasic
		mockState.ApiOauth2AllowAllPerms = *req.ApiOauth2AllowAllPerms
		mockState.ApiOauth2Expiration = *req.ApiOauth2Expiration
		mockState.LdapServer = req.LdapServer
		mockState.LdapBaseDN = req.LdapBaseDN
		mockState.LdapBindUsername = req.LdapBindUsername
		mockState.LdapUserSearchDN = req.LdapUserSearchDN
		mockState.LdapUserFilter = req.LdapUserFilter
		mockState.LdapUserSearchFilter = req.LdapUserSearchFilter
		mockState.LdapUserGroupAttributes = req.LdapUserGroupAttributes
		mockState.LdapGroupSearchDN = req.LdapGroupSearchDN
		mockState.LdapGroupFilter = req.LdapGroupFilter
		mockState.LdapGroupMembershipFilter = req.LdapGroupMembershipFilter
		mockState.LdapUseGlobalCatalog = *req.LdapUseGlobalCatalog
		mockState.LdapPermitNoTLS = *req.LdapPermitNoTLS
		mockState.OidcMetadataURL = req.OidcMetadataURL
		mockState.OidcMetadata = req.OidcMetadata
		mockState.OidcClientID = req.OidcClientID
		mockState.OidcAuthMethod = req.OidcAuthMethod
		mockState.OidcScope = req.OidcScope
		mockState.OidcAuthorizeURL = req.OidcAuthorizeURL
		mockState.OidcTokenEndpointURL = req.OidcTokenEndpointURL
		mockState.OidcUsernameField = req.OidcUsernameField
		mockState.OidcGroupsField = req.OidcGroupsField
		mockState.OidcRequiredKey = req.OidcRequiredKey
		mockState.OidcRequiredValue = req.OidcRequiredValue
		mockState.OidcDomainHint = req.OidcDomainHint
		mockState.OidcLoginButton = req.OidcLoginButton

		*auth = *mockState
	}).Maybe()

	// GetJSON mock — returns current mockState for all Read operations.
	client.On("GetJSON", mock.Anything, "configuration/v1/authentication/1/",
		mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		auth := args.Get(3).(*config.Authentication)
		*auth = *mockState
	}).Maybe()

	testInfinityAuthentication(t, client)
}

func testInfinityAuthentication(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		ExternalProviders: map[string]resource.ExternalProvider{
			"tls": {
				Source: "hashicorp/tls",
			},
		},
		Steps: []resource.TestStep{
			// Step 1: Apply full configuration with all fields set to non-default values.
			{
				Config: test.LoadTestFolder(t, "resource_infinity_authentication_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_authentication.authentication-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "source", "LOCAL"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "client_certificate", "NO"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "api_oauth2_disable_basic", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "api_oauth2_allow_all_perms", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "api_oauth2_expiration", "7200"),
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
			// Step 2: Destroy — triggers Delete which must reset all fields to API defaults.
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_authentication_full"),
				Destroy: true,
			},
			// Step 3: Re-apply min config (empty body, all defaults) and verify
			// the API returned all fields to their defaults after the destroy.
			{
				Config: test.LoadTestFolder(t, "resource_infinity_authentication_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_authentication.authentication-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "source", "LOCAL"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "client_certificate", "NO"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "api_oauth2_disable_basic", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "api_oauth2_allow_all_perms", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "api_oauth2_expiration", "3600"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_server", ""),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_base_dn", ""),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_bind_username", ""),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_user_filter", "(&(objectclass=person)(!(objectclass=computer)))"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_user_search_filter", "(|(uid={username})(sAMAccountName={username}))"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_user_group_attributes", "memberOf"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_group_filter", "(|(objectclass=group)(objectclass=groupOfNames)(objectclass=groupOfUniqueNames)(objectclass=posixGroup))"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_group_membership_filter", "(|(member={userdn})(uniquemember={userdn})(memberuid={useruid}))"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_use_global_catalog", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "ldap_permit_no_tls", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_metadata_url", ""),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_client_id", ""),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_auth_method", "client_secret"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_scope", "openid email profile"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_authorize_url", ""),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_token_endpoint_url", ""),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_username_field", "preferred_username"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_groups_field", "groups"),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_required_key", ""),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_required_value", ""),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_domain_hint", ""),
					resource.TestCheckResourceAttr("pexip_infinity_authentication.authentication-test", "oidc_login_button", ""),
				),
			},
		},
	})
}

func TestInfinityAuthenticationValidation(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	client := infinity.NewClientMock()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: `
resource "pexip_infinity_authentication" "authentication-test" {
  source = "INVALID"
}
`,
				ExpectError: regexp.MustCompile(`value must be one of`),
			},
			{
				Config: `
resource "pexip_infinity_authentication" "authentication-test" {
  client_certificate = "INVALID"
}
`,
				ExpectError: regexp.MustCompile(`value must be one of`),
			},
			{
				Config: `
resource "pexip_infinity_authentication" "authentication-test" {
  oidc_auth_method = "INVALID"
}
`,
				ExpectError: regexp.MustCompile(`value must be one of`),
			},
			{
				Config: `
resource "pexip_infinity_authentication" "authentication-test" {
  ldap_server = "` + strings.Repeat("a", 256) + `"
}
`,
				ExpectError: regexp.MustCompile(`string length must be at most 255`),
			},
			{
				Config: `
resource "pexip_infinity_authentication" "authentication-test" {
  ldap_base_dn = "` + strings.Repeat("a", 256) + `"
}
`,
				ExpectError: regexp.MustCompile(`string length must be at most 255`),
			},
			{
				Config: `
resource "pexip_infinity_authentication" "authentication-test" {
  ldap_bind_username = "` + strings.Repeat("a", 256) + `"
}
`,
				ExpectError: regexp.MustCompile(`string length must be at most 255`),
			},
			{
				Config: `
resource "pexip_infinity_authentication" "authentication-test" {
  ldap_bind_password = "` + strings.Repeat("a", 101) + `"
}
`,
				ExpectError: regexp.MustCompile(`string length must be at most 100`),
			},
			{
				Config: `
resource "pexip_infinity_authentication" "authentication-test" {
  ldap_user_search_dn = "` + strings.Repeat("a", 256) + `"
}
`,
				ExpectError: regexp.MustCompile(`string length must be at most 255`),
			},
			{
				Config: `
resource "pexip_infinity_authentication" "authentication-test" {
  ldap_user_filter = "` + strings.Repeat("a", 1025) + `"
}
`,
				ExpectError: regexp.MustCompile(`string length must be at most 1024`),
			},
			{
				Config: `
resource "pexip_infinity_authentication" "authentication-test" {
  ldap_user_search_filter = "` + strings.Repeat("a", 1025) + `"
}
`,
				ExpectError: regexp.MustCompile(`string length must be at most 1024`),
			},
			{
				Config: `
resource "pexip_infinity_authentication" "authentication-test" {
  ldap_user_group_attributes = "` + strings.Repeat("a", 101) + `"
}
`,
				ExpectError: regexp.MustCompile(`string length must be at most 100`),
			},
			{
				Config: `
resource "pexip_infinity_authentication" "authentication-test" {
  ldap_group_search_dn = "` + strings.Repeat("a", 256) + `"
}
`,
				ExpectError: regexp.MustCompile(`string length must be at most 255`),
			},
			{
				Config: `
resource "pexip_infinity_authentication" "authentication-test" {
  ldap_group_filter = "` + strings.Repeat("a", 1025) + `"
}
`,
				ExpectError: regexp.MustCompile(`string length must be at most 1024`),
			},
			{
				Config: `
resource "pexip_infinity_authentication" "authentication-test" {
  ldap_group_membership_filter = "` + strings.Repeat("a", 1025) + `"
}
`,
				ExpectError: regexp.MustCompile(`string length must be at most 1024`),
			},
			{
				Config: `
resource "pexip_infinity_authentication" "authentication-test" {
  oidc_domain_hint = "` + strings.Repeat("a", 256) + `"
}
`,
				ExpectError: regexp.MustCompile(`string length must be at most 255`),
			},
			{
				Config: `
resource "pexip_infinity_authentication" "authentication-test" {
  oidc_login_button = "` + strings.Repeat("a", 129) + `"
}
`,
				ExpectError: regexp.MustCompile(`string length must be at most 128`),
			},
		},
	})
}
