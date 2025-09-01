/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"testing"

	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityADFSAuthServer(t *testing.T) {
	t.Parallel()

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateADFSAuthServer API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/adfs_auth_server/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/adfs_auth_server/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.ADFSAuthServer{
		ID:                             123,
		ResourceURI:                    "/api/admin/configuration/v1/adfs_auth_server/123/",
		Name:                           "adfs_auth_server-test",
		Description:                    "Test ADFSAuthServer",
		ClientID:                       "test-value",
		FederationServiceName:          "adfs_auth_server-test",
		FederationServiceIdentifier:    "test-value",
		RelyingPartyTrustIdentifierURL: "https://example.com",
	}

	// Mock the GetADFSAuthServer API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/adfs_auth_server/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		adfsauthserver := args.Get(2).(*config.ADFSAuthServer)
		*adfsauthserver = *mockState
	}).Maybe()

	// Mock the UpdateADFSAuthServer API call
	client.On("PutJSON", mock.Anything, "configuration/v1/adfs_auth_server/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.ADFSAuthServerUpdateRequest)
		adfsauthserver := args.Get(3).(*config.ADFSAuthServer)

		// Update mock state with all fields from the update request
		mockState.Name = updateRequest.Name
		mockState.Description = updateRequest.Description
		mockState.ClientID = updateRequest.ClientID
		mockState.FederationServiceName = updateRequest.FederationServiceName
		mockState.FederationServiceIdentifier = updateRequest.FederationServiceIdentifier
		mockState.RelyingPartyTrustIdentifierURL = updateRequest.RelyingPartyTrustIdentifierURL

		// Return updated state
		*adfsauthserver = *mockState
	}).Maybe()

	// Mock the DeleteADFSAuthServer API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/adfs_auth_server/123/"
	}), mock.Anything).Return(nil)

	testInfinityADFSAuthServer(t, client)
}

func testInfinityADFSAuthServer(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_adfs_auth_server_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_adfs_auth_server.adfs_auth_server-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_adfs_auth_server.adfs_auth_server-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.adfs_auth_server-test", "name", "adfs_auth_server-test"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.adfs_auth_server-test", "description", "Test ADFSAuthServer"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.adfs_auth_server-test", "federation_service_name", "adfs_auth_server-test"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.adfs_auth_server-test", "client_id", "test-value"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.adfs_auth_server-test", "federation_service_identifier", "test-value"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.adfs_auth_server-test", "relying_party_trust_identifier_url", "https://example.com"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_adfs_auth_server_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_adfs_auth_server.adfs_auth_server-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_adfs_auth_server.adfs_auth_server-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.adfs_auth_server-test", "name", "adfs_auth_server-test"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.adfs_auth_server-test", "description", "Updated Test ADFSAuthServer"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.adfs_auth_server-test", "federation_service_name", "adfs_auth_server-test"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.adfs_auth_server-test", "client_id", "updated-value"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.adfs_auth_server-test", "federation_service_identifier", "updated-value"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.adfs_auth_server-test", "relying_party_trust_identifier_url", "https://updated.example.com"),
				),
			},
		},
	})
}
