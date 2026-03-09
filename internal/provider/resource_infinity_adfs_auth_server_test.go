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

func TestInfinityADFSAuthServer(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

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
		Name:                           "tf-test-adfs-auth-server",
		Description:                    "Full test configuration for ADFS Auth Server",
		ClientID:                       "test-client-id-full",
		FederationServiceName:          "adfs-full.example.com",
		FederationServiceIdentifier:    "https://adfs-full.example.com/adfs/services/trust",
		RelyingPartyTrustIdentifierURL: "https://full.example.com",
	}

	// Mock the GetADFSAuthServer API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/adfs_auth_server/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		adfsauthserver := args.Get(3).(*config.ADFSAuthServer)
		*adfsauthserver = *mockState
	}).Maybe()

	// Mock the UpdateADFSAuthServer API call
	client.On("PutJSON", mock.Anything, "configuration/v1/adfs_auth_server/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.ADFSAuthServerUpdateRequest)
		adfsauthserver := args.Get(3).(*config.ADFSAuthServer)

		// Update mock state with all fields from the update request
		if updateRequest.Name != "" {
			mockState.Name = updateRequest.Name
		}
		if updateRequest.Description != "" {
			mockState.Description = updateRequest.Description
		}
		if updateRequest.ClientID != "" {
			mockState.ClientID = updateRequest.ClientID
		}
		if updateRequest.FederationServiceName != "" {
			mockState.FederationServiceName = updateRequest.FederationServiceName
		}
		if updateRequest.FederationServiceIdentifier != "" {
			mockState.FederationServiceIdentifier = updateRequest.FederationServiceIdentifier
		}
		if updateRequest.RelyingPartyTrustIdentifierURL != "" {
			mockState.RelyingPartyTrustIdentifierURL = updateRequest.RelyingPartyTrustIdentifierURL
		}

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
			// Test 1: Create with full configuration
			{
				Config: test.LoadTestFolder(t, "resource_infinity_adfs_auth_server_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "name", "tf-test-adfs-auth-server"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "description", "Full test configuration for ADFS Auth Server"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "client_id", "test-client-id-full"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "federation_service_name", "adfs-full.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "federation_service_identifier", "https://adfs-full.example.com/adfs/services/trust"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "relying_party_trust_identifier_url", "https://full.example.com"),
				),
			},
			// Test 2: Update to min configuration, then delete
			{
				Config: test.LoadTestFolder(t, "resource_infinity_adfs_auth_server_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "name", "tf-test-adfs-auth-server"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "description", "Full test configuration for ADFS Auth Server"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "client_id", "test-client-id-min"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "federation_service_name", "adfs-min.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "federation_service_identifier", "https://adfs-min.example.com/adfs/services/trust"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "relying_party_trust_identifier_url", "https://min.example.com"),
				),
			},
			// Test 3: Create with min configuration
			{
				Config: test.LoadTestFolder(t, "resource_infinity_adfs_auth_server_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "name", "tf-test-adfs-auth-server"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "client_id", "test-client-id-min"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "federation_service_name", "adfs-min.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "federation_service_identifier", "https://adfs-min.example.com/adfs/services/trust"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "relying_party_trust_identifier_url", "https://min.example.com"),
				),
			},
			// Test 4: Update to full configuration
			{
				Config: test.LoadTestFolder(t, "resource_infinity_adfs_auth_server_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "name", "tf-test-adfs-auth-server"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "description", "Full test configuration for ADFS Auth Server"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "client_id", "test-client-id-full"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "federation_service_name", "adfs-full.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "federation_service_identifier", "https://adfs-full.example.com/adfs/services/trust"),
					resource.TestCheckResourceAttr("pexip_infinity_adfs_auth_server.tf-test-adfs-auth-server", "relying_party_trust_identifier_url", "https://full.example.com"),
				),
			},
		},
	})
}
