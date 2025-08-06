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

func TestInfinityIdentityProvider(t *testing.T) {
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateIdentityprovider API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/identity_provider/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/identity_provider/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.IdentityProvider{
		ID:                                  123,
		ResourceURI:                         "/api/admin/configuration/v1/identity_provider/123/",
		Name:                                "identity_provider-test",
		Description:                         "Test IdentityProvider",
		IdpType:                             "saml",
		SSOUrl:                              "https://example.com",
		IdpEntityID:                         "test-value",
		IdpPublicKey:                        "test-value",
		ServiceEntityID:                     "test-value",
		ServicePublicKey:                    "test-value",
		ServicePrivateKey:                   "test-value",
		SignatureAlgorithm:                  "rsa-sha256",
		DigestAlgorithm:                     "sha256",
		DisplayNameAttributeName:            "identity_provider-test",
		RegistrationAliasAttributeName:      "identity_provider-test",
		AssertionConsumerServiceURL:         "https://example.com",
		WorkerFQDNACSURLs:                   true,
		DisablePopupFlow:                    true,
		OidcFlow:                            "authorization_code",
		OidcClientID:                        "test-value",
		OidcClientSecret:                    "test-value",
		OidcTokenURL:                        "https://example.com",
		OidcUserInfoURL:                     "https://example.com",
		OidcJWKSURL:                         "https://example.com",
		OidcTokenEndpointAuthScheme:         "client_secret_basic",
		OidcTokenSignatureScheme:            "rs256",
		OidcDisplayNameClaimName:            "identity_provider-test",
		OidcRegistrationAliasClaimName:      "identity_provider-test",
		OidcAdditionalScopes:                "test-value",
		OidcFranceConnectRequiredEidasLevel: "eidas1",
		Attributes:                          "test-value",
	}

	// Mock the GetIdentityprovider API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/identity_provider/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		identity_provider := args.Get(2).(*config.IdentityProvider)
		*identity_provider = *mockState
	}).Maybe()

	// Mock the UpdateIdentityprovider API call
	client.On("PutJSON", mock.Anything, "configuration/v1/identity_provider/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.IdentityProviderUpdateRequest)
		identity_provider := args.Get(3).(*config.IdentityProvider)

		// Update mock state based on request
		mockState.Name = updateRequest.Name
		mockState.Description = updateRequest.Description
		mockState.IdpType = updateRequest.IdpType
		mockState.SSOUrl = updateRequest.SSOUrl
		mockState.IdpEntityID = updateRequest.IdpEntityID
		mockState.IdpPublicKey = updateRequest.IdpPublicKey
		mockState.ServiceEntityID = updateRequest.ServiceEntityID
		mockState.ServicePublicKey = updateRequest.ServicePublicKey
		mockState.ServicePrivateKey = updateRequest.ServicePrivateKey
		mockState.SignatureAlgorithm = updateRequest.SignatureAlgorithm
		mockState.DigestAlgorithm = updateRequest.DigestAlgorithm
		mockState.DisplayNameAttributeName = updateRequest.DisplayNameAttributeName
		mockState.RegistrationAliasAttributeName = updateRequest.RegistrationAliasAttributeName
		mockState.AssertionConsumerServiceURL = updateRequest.AssertionConsumerServiceURL
		if updateRequest.WorkerFQDNACSURLs != nil {
			mockState.WorkerFQDNACSURLs = *updateRequest.WorkerFQDNACSURLs
		}
		if updateRequest.DisablePopupFlow != nil {
			mockState.DisablePopupFlow = *updateRequest.DisablePopupFlow
		}
		mockState.OidcFlow = updateRequest.OidcFlow
		mockState.OidcClientID = updateRequest.OidcClientID
		mockState.OidcClientSecret = updateRequest.OidcClientSecret
		mockState.OidcTokenURL = updateRequest.OidcTokenURL
		mockState.OidcUserInfoURL = updateRequest.OidcUserInfoURL
		mockState.OidcJWKSURL = updateRequest.OidcJWKSURL
		mockState.OidcTokenEndpointAuthScheme = updateRequest.OidcTokenEndpointAuthScheme
		mockState.OidcTokenSignatureScheme = updateRequest.OidcTokenSignatureScheme
		mockState.OidcDisplayNameClaimName = updateRequest.OidcDisplayNameClaimName
		mockState.OidcRegistrationAliasClaimName = updateRequest.OidcRegistrationAliasClaimName
		mockState.OidcAdditionalScopes = updateRequest.OidcAdditionalScopes
		mockState.OidcFranceConnectRequiredEidasLevel = updateRequest.OidcFranceConnectRequiredEidasLevel
		mockState.Attributes = updateRequest.Attributes

		// Return updated state
		*identity_provider = *mockState
	}).Maybe()

	// Mock the DeleteIdentityprovider API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/identity_provider/123/"
	}), mock.Anything).Return(nil)

	testInfinityIdentityProvider(t, client)
}

func testInfinityIdentityProvider(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_identity_provider_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.identity_provider-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.identity_provider-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.identity_provider-test", "name", "identity_provider-test"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.identity_provider-test", "description", "Test IdentityProvider"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.identity_provider-test", "display_name_attribute_name", "identity_provider-test"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.identity_provider-test", "registration_alias_attribute_name", "identity_provider-test"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.identity_provider-test", "worker_fqdn_acs_urls", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.identity_provider-test", "disable_popup_flow", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.identity_provider-test", "oidc_display_name_claim_name", "identity_provider-test"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.identity_provider-test", "oidc_registration_alias_claim_name", "identity_provider-test"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_identity_provider_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.identity_provider-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.identity_provider-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.identity_provider-test", "name", "identity_provider-test"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.identity_provider-test", "description", "Updated Test IdentityProvider"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.identity_provider-test", "idp_type", "oidc"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.identity_provider-test", "idp_entity_id", "updated-value"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.identity_provider-test", "idp_public_key", "updated-value"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.identity_provider-test", "oidc_token_url", "https://updated.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.identity_provider-test", "oidc_jwks_url", "https://updated.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.identity_provider-test", "worker_fqdn_acs_urls", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.identity_provider-test", "disable_popup_flow", "false"),
				),
			},
		},
	})
}
