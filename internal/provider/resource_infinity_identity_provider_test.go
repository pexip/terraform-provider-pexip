/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityIdentityProvider(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking
	mockState := &config.IdentityProvider{
		ID:                                  123,
		ResourceURI:                         "/api/admin/configuration/v1/identity_provider/123/",
		Name:                                "identity_provider-test",
		Description:                         "Test IdentityProvider",
		IdpType:                             "saml",
		UUID:                                "12345678-1234-1234-1234-123456789abc",
		SSOUrl:                              "https://example.com/sso",
		IdpEntityID:                         "https://example.com/entity",
		ServiceEntityID:                     "https://pexip.example.com/entity",
		ServicePrivateKey:                   "",
		SignatureAlgorithm:                  "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256",
		DigestAlgorithm:                     "http://www.w3.org/2001/04/xmlenc#sha256",
		DisplayNameAttributeName:            "displayName",
		RegistrationAliasAttributeName:      "email",
		AssertionConsumerServiceURL:         "https://localhost/samlconsumer/12345678-1234-1234-1234-123456789abc",
		WorkerFQDNACSURLs:                   true,
		DisablePopupFlow:                    true,
		OidcFlow:                            "code",
		OidcClientSecret:                    "",
		OidcTokenEndpointAuthScheme:         "client_secret_basic",
		OidcTokenSignatureScheme:            "rs256",
		OidcFranceConnectRequiredEidasLevel: "eidas1",
	}

	// Mock the CreateIdentityprovider API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/identity_provider/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/identity_provider/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.IdentityProviderCreateRequest)
		// Update mockState with values from create request
		mockState.Name = createReq.Name
		mockState.Description = createReq.Description
		mockState.IdpType = createReq.IdpType
		mockState.UUID = createReq.UUID
		mockState.SSOUrl = createReq.SSOUrl
		mockState.IdpEntityID = createReq.IdpEntityID
		mockState.IdpPublicKey = createReq.IdpPublicKey
		mockState.ServiceEntityID = createReq.ServiceEntityID
		mockState.ServicePublicKey = createReq.ServicePublicKey
		mockState.ServicePrivateKey = createReq.ServicePrivateKey
		mockState.SignatureAlgorithm = createReq.SignatureAlgorithm
		mockState.DigestAlgorithm = createReq.DigestAlgorithm
		mockState.DisplayNameAttributeName = createReq.DisplayNameAttributeName
		mockState.RegistrationAliasAttributeName = createReq.RegistrationAliasAttributeName
		mockState.AssertionConsumerServiceURL = createReq.AssertionConsumerServiceURL
		mockState.AssertionConsumerServiceURL2 = createReq.AssertionConsumerServiceURL2
		mockState.AssertionConsumerServiceURL3 = createReq.AssertionConsumerServiceURL3
		mockState.AssertionConsumerServiceURL4 = createReq.AssertionConsumerServiceURL4
		mockState.AssertionConsumerServiceURL5 = createReq.AssertionConsumerServiceURL5
		mockState.AssertionConsumerServiceURL6 = createReq.AssertionConsumerServiceURL6
		mockState.AssertionConsumerServiceURL7 = createReq.AssertionConsumerServiceURL7
		mockState.AssertionConsumerServiceURL8 = createReq.AssertionConsumerServiceURL8
		mockState.AssertionConsumerServiceURL9 = createReq.AssertionConsumerServiceURL9
		mockState.AssertionConsumerServiceURL10 = createReq.AssertionConsumerServiceURL10
		mockState.WorkerFQDNACSURLs = createReq.WorkerFQDNACSURLs
		mockState.DisablePopupFlow = createReq.DisablePopupFlow
		mockState.OidcFlow = createReq.OidcFlow
		mockState.OidcClientID = createReq.OidcClientID
		mockState.OidcClientSecret = createReq.OidcClientSecret
		mockState.OidcTokenURL = createReq.OidcTokenURL
		mockState.OidcUserInfoURL = createReq.OidcUserInfoURL
		mockState.OidcJWKSURL = createReq.OidcJWKSURL
		mockState.OidcTokenEndpointAuthScheme = createReq.OidcTokenEndpointAuthScheme
		mockState.OidcTokenSignatureScheme = createReq.OidcTokenSignatureScheme
		mockState.OidcDisplayNameClaimName = createReq.OidcDisplayNameClaimName
		mockState.OidcRegistrationAliasClaimName = createReq.OidcRegistrationAliasClaimName
		mockState.OidcAdditionalScopes = createReq.OidcAdditionalScopes
		mockState.OidcFranceConnectRequiredEidasLevel = createReq.OidcFranceConnectRequiredEidasLevel

		// Convert attributes from URIs to IdentityProviderAttribute structs
		if createReq.Attributes != nil {
			attrs := make([]config.IdentityProviderAttribute, 0, len(*createReq.Attributes))
			for _, uri := range *createReq.Attributes {
				// Extract ID from URI (e.g., "/api/admin/configuration/v1/identity_provider_attribute/1/")
				var id int
				if _, err := fmt.Sscanf(uri, "/api/admin/configuration/v1/identity_provider_attribute/%d/", &id); err == nil {
					attrs = append(attrs, config.IdentityProviderAttribute{
						ID:          id,
						ResourceURI: uri,
					})
				}
			}
			mockState.Attributes = &attrs
		} else {
			mockState.Attributes = nil
		}
	})

	// Mock the CreateIdentityProviderAttribute API calls for the attributes
	attrCreateResponse1 := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/identity_provider_attribute/1/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/identity_provider_attribute/", mock.MatchedBy(func(req *config.IdentityProviderAttributeCreateRequest) bool {
		return req.Name == "tf-test-displayName"
	}), mock.Anything).Return(attrCreateResponse1, nil).Maybe()

	attrCreateResponse2 := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/identity_provider_attribute/2/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/identity_provider_attribute/", mock.MatchedBy(func(req *config.IdentityProviderAttributeCreateRequest) bool {
		return req.Name == "tf-test-email"
	}), mock.Anything).Return(attrCreateResponse2, nil).Maybe()

	// Mock GetIdentityProviderAttribute for reads
	client.On("GetJSON", mock.Anything, "configuration/v1/identity_provider_attribute/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		attr := args.Get(3).(*config.IdentityProviderAttribute)
		*attr = config.IdentityProviderAttribute{
			ID:          1,
			Name:        "tf-test-displayName",
			Description: "Test attribute for display name",
			ResourceURI: "/api/admin/configuration/v1/identity_provider_attribute/1/",
		}
	}).Maybe()

	client.On("GetJSON", mock.Anything, "configuration/v1/identity_provider_attribute/2/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		attr := args.Get(3).(*config.IdentityProviderAttribute)
		*attr = config.IdentityProviderAttribute{
			ID:          2,
			Name:        "tf-test-email",
			Description: "Test attribute for email",
			ResourceURI: "/api/admin/configuration/v1/identity_provider_attribute/2/",
		}
	}).Maybe()

	// Mock UpdateIdentityProviderAttribute
	client.On("PutJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/identity_provider_attribute/1/" || path == "configuration/v1/identity_provider_attribute/2/"
	}), mock.Anything, mock.Anything).Return(nil).Maybe()

	// Mock DeleteIdentityProviderAttribute
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/identity_provider_attribute/1/" || path == "configuration/v1/identity_provider_attribute/2/"
	}), mock.Anything).Return(nil).Maybe()

	// Mock the GetIdentityprovider API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/identity_provider/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		identity_provider := args.Get(3).(*config.IdentityProvider)
		// Copy mockState but clear sensitive fields (API doesn't return them)
		*identity_provider = *mockState
		identity_provider.ServicePrivateKey = ""
		identity_provider.OidcClientSecret = ""
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
		mockState.AssertionConsumerServiceURL2 = updateRequest.AssertionConsumerServiceURL2
		mockState.AssertionConsumerServiceURL3 = updateRequest.AssertionConsumerServiceURL3
		mockState.AssertionConsumerServiceURL4 = updateRequest.AssertionConsumerServiceURL4
		mockState.AssertionConsumerServiceURL5 = updateRequest.AssertionConsumerServiceURL5
		mockState.AssertionConsumerServiceURL6 = updateRequest.AssertionConsumerServiceURL6
		mockState.AssertionConsumerServiceURL7 = updateRequest.AssertionConsumerServiceURL7
		mockState.AssertionConsumerServiceURL8 = updateRequest.AssertionConsumerServiceURL8
		mockState.AssertionConsumerServiceURL9 = updateRequest.AssertionConsumerServiceURL9
		mockState.AssertionConsumerServiceURL10 = updateRequest.AssertionConsumerServiceURL10
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

		// Convert attributes from URIs to IdentityProviderAttribute structs
		if updateRequest.Attributes != nil {
			attrs := make([]config.IdentityProviderAttribute, 0, len(*updateRequest.Attributes))
			for _, uri := range *updateRequest.Attributes {
				// Extract ID from URI (e.g., "/api/admin/configuration/v1/identity_provider_attribute/1/")
				var id int
				if _, err := fmt.Sscanf(uri, "/api/admin/configuration/v1/identity_provider_attribute/%d/", &id); err == nil {
					attrs = append(attrs, config.IdentityProviderAttribute{
						ID:          id,
						ResourceURI: uri,
					})
				}
			}
			if len(attrs) > 0 {
				mockState.Attributes = &attrs
			} else {
				mockState.Attributes = nil
			}
		} else {
			mockState.Attributes = nil
		}

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
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source: "hashicorp/random",
			},
		},
		Steps: []resource.TestStep{
			{
				// Step 1: Create with full config
				Config: test.LoadTestFolder(t, "resource_infinity_identity_provider_full"),
				Check: resource.ComposeTestCheckFunc(
					// IDs and required fields
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "resource_id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "uuid"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "name", "tf-test Identity Provider full"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "assertion_consumer_service_url"),

					// Optional basic fields
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "description", "Full test Identity Provider with all fields"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "idp_type", "oidc"),

					// SAML specific fields
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "sso_url", "https://idp.example.com/sso"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "idp_entity_id", "https://idp.example.com/entity"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "idp_public_key"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "service_entity_id", "https://pexip.example.com/entity"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "service_public_key"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "service_private_key"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "signature_algorithm", "http://www.w3.org/2001/04/xmldsig-more#rsa-sha384"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "digest_algorithm", "http://www.w3.org/2001/04/xmldsig-more#sha384"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "display_name_attribute_name", "displayName"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "registration_alias_attribute_name", "userPrincipalName"),

					// Additional assertion consumer service URLs
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "assertion_consumer_service_url2"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "assertion_consumer_service_url3"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "assertion_consumer_service_url4"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "assertion_consumer_service_url5"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "assertion_consumer_service_url6"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "assertion_consumer_service_url7"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "assertion_consumer_service_url8"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "assertion_consumer_service_url9"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "assertion_consumer_service_url10"),

					// Worker and popup settings
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "worker_fqdn_acs_urls", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "disable_popup_flow", "true"),

					// OIDC specific fields
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_flow", "implicit"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_client_id", "test-client-id-12345"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_client_secret", "test-client-secret-67890"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_token_url", "https://idp.example.com/oauth2/token"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_user_info_url", "https://idp.example.com/oauth2/userinfo"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_jwks_url", "https://idp.example.com/.well-known/jwks.json"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_token_endpoint_auth_scheme", "client_secret_basic"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_token_signature_scheme", "hs256"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_display_name_claim_name", "full_name"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_registration_alias_claim_name", "preferred_username"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_additional_scopes", "profile email phone address"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_france_connect_required_eidas_level", "eidas3"),

					// Attributes
					resource.TestCheckTypeSetElemAttrPair("pexip_infinity_identity_provider.test", "attributes.*", "pexip_infinity_identity_provider_attribute.attr1", "id"),
					resource.TestCheckTypeSetElemAttrPair("pexip_infinity_identity_provider.test", "attributes.*", "pexip_infinity_identity_provider_attribute.attr2", "id"),
				),
			},
			{
				// Step 2: Update to min config
				Config: test.LoadTestFolder(t, "resource_infinity_identity_provider_min"),
				Check: resource.ComposeTestCheckFunc(
					// IDs and required fields
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "resource_id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "uuid"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "name", "tf-test Identity Provider min"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "assertion_consumer_service_url"),

					// Optional fields cleared - verify defaults
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "idp_type", "saml"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "sso_url", ""),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "idp_entity_id", ""),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "idp_public_key", ""),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "service_entity_id", ""),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "service_public_key", ""),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "service_private_key", ""),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "signature_algorithm", "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "digest_algorithm", "http://www.w3.org/2001/04/xmlenc#sha256"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "display_name_attribute_name", "NameId"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "registration_alias_attribute_name", "NameId"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "assertion_consumer_service_url2", ""),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "assertion_consumer_service_url3", ""),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "assertion_consumer_service_url4", ""),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "assertion_consumer_service_url5", ""),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "assertion_consumer_service_url6", ""),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "assertion_consumer_service_url7", ""),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "assertion_consumer_service_url8", ""),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "assertion_consumer_service_url9", ""),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "assertion_consumer_service_url10", ""),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "worker_fqdn_acs_urls", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "disable_popup_flow", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_flow", "code"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_client_id", ""),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_client_secret", ""),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_token_url", ""),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_user_info_url", ""),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_jwks_url", ""),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_token_endpoint_auth_scheme", "client_secret_post"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_token_signature_scheme", "rs256"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_display_name_claim_name", "name"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_registration_alias_claim_name", "sub"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_additional_scopes", ""),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_france_connect_required_eidas_level", "disabled"),
				),
			},
			{
				// Step 3: Destroy and recreate with minimal config
				Config:  test.LoadTestFolder(t, "resource_infinity_identity_provider_min"),
				Destroy: true,
			},
			{
				// Step 4: Recreate with minimal config (after destroy)
				Config: test.LoadTestFolder(t, "resource_infinity_identity_provider_min"),
				Check: resource.ComposeTestCheckFunc(
					// IDs and required fields
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "resource_id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "uuid"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "name", "tf-test Identity Provider min"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "assertion_consumer_service_url"),

					// Optional fields - verify defaults on create
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "idp_type", "saml"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "signature_algorithm", "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "digest_algorithm", "http://www.w3.org/2001/04/xmlenc#sha256"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "display_name_attribute_name", "NameId"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "registration_alias_attribute_name", "NameId"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "worker_fqdn_acs_urls", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "disable_popup_flow", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_flow", "code"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_token_endpoint_auth_scheme", "client_secret_post"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_token_signature_scheme", "rs256"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_display_name_claim_name", "name"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_registration_alias_claim_name", "sub"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_france_connect_required_eidas_level", "disabled"),
				),
			},
			{
				// Step 5: Update to full config
				Config: test.LoadTestFolder(t, "resource_infinity_identity_provider_full"),
				Check: resource.ComposeTestCheckFunc(
					// IDs and required fields
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "resource_id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "uuid"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "name", "tf-test Identity Provider full"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "assertion_consumer_service_url"),

					// Optional basic fields
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "description", "Full test Identity Provider with all fields"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "idp_type", "oidc"),

					// SAML specific fields
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "sso_url", "https://idp.example.com/sso"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "idp_entity_id", "https://idp.example.com/entity"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "idp_public_key"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "service_entity_id", "https://pexip.example.com/entity"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "service_public_key"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "service_private_key"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "signature_algorithm", "http://www.w3.org/2001/04/xmldsig-more#rsa-sha384"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "digest_algorithm", "http://www.w3.org/2001/04/xmldsig-more#sha384"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "display_name_attribute_name", "displayName"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "registration_alias_attribute_name", "userPrincipalName"),

					// Additional assertion consumer service URLs
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "assertion_consumer_service_url2"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "assertion_consumer_service_url3"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "assertion_consumer_service_url4"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "assertion_consumer_service_url5"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "assertion_consumer_service_url6"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "assertion_consumer_service_url7"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "assertion_consumer_service_url8"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "assertion_consumer_service_url9"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider.test", "assertion_consumer_service_url10"),

					// Worker and popup settings
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "worker_fqdn_acs_urls", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "disable_popup_flow", "true"),

					// OIDC specific fields
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_flow", "implicit"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_client_id", "test-client-id-12345"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_client_secret", "test-client-secret-67890"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_token_url", "https://idp.example.com/oauth2/token"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_user_info_url", "https://idp.example.com/oauth2/userinfo"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_jwks_url", "https://idp.example.com/.well-known/jwks.json"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_token_endpoint_auth_scheme", "client_secret_basic"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_token_signature_scheme", "hs256"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_display_name_claim_name", "full_name"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_registration_alias_claim_name", "preferred_username"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_additional_scopes", "profile email phone address"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider.test", "oidc_france_connect_required_eidas_level", "eidas3"),

					// Attributes
					resource.TestCheckTypeSetElemAttrPair("pexip_infinity_identity_provider.test", "attributes.*", "pexip_infinity_identity_provider_attribute.attr1", "id"),
					resource.TestCheckTypeSetElemAttrPair("pexip_infinity_identity_provider.test", "attributes.*", "pexip_infinity_identity_provider_attribute.attr2", "id"),
				),
			},
		},
	})
}
