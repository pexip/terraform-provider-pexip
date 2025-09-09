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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityGoogleAuthServer(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateGoogleauthserver API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/google_auth_server/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/google_auth_server/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.GoogleAuthServer{
		ID:              123,
		ResourceURI:     "/api/admin/configuration/v1/google_auth_server/123/",
		Name:            "google_auth_server-test",
		Description:     "Test GoogleAuthServer",
		ApplicationType: "web",
		ClientID:        test.StringPtr("test-value"),
		ClientSecret:    "test-value",
	}

	// Mock the GetGoogleauthserver API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/google_auth_server/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		google_auth_server := args.Get(3).(*config.GoogleAuthServer)
		*google_auth_server = *mockState
	}).Maybe()

	// Mock the UpdateGoogleauthserver API call
	client.On("PutJSON", mock.Anything, "configuration/v1/google_auth_server/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.GoogleAuthServerUpdateRequest)
		google_auth_server := args.Get(3).(*config.GoogleAuthServer)

		// Update mock state based on request
		if updateRequest.Name != "" {
			mockState.Name = updateRequest.Name
		}
		if updateRequest.Description != "" {
			mockState.Description = updateRequest.Description
		}
		if updateRequest.ApplicationType != "" {
			mockState.ApplicationType = updateRequest.ApplicationType
		}
		if updateRequest.ClientID != nil && *updateRequest.ClientID != "" {
			mockState.ClientID = updateRequest.ClientID
		}
		if updateRequest.ClientSecret != "" {
			mockState.ClientSecret = updateRequest.ClientSecret
		}

		// Return updated state
		*google_auth_server = *mockState
	}).Maybe()

	// Mock the DeleteGoogleauthserver API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/google_auth_server/123/"
	}), mock.Anything).Return(nil)

	testInfinityGoogleAuthServer(t, client)
}

func testInfinityGoogleAuthServer(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_google_auth_server_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_google_auth_server.google_auth_server-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_google_auth_server.google_auth_server-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_google_auth_server.google_auth_server-test", "name", "google_auth_server-test"),
					resource.TestCheckResourceAttr("pexip_infinity_google_auth_server.google_auth_server-test", "description", "Test GoogleAuthServer"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_google_auth_server_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_google_auth_server.google_auth_server-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_google_auth_server.google_auth_server-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_google_auth_server.google_auth_server-test", "name", "google_auth_server-test"),
					resource.TestCheckResourceAttr("pexip_infinity_google_auth_server.google_auth_server-test", "description", "Updated Test GoogleAuthServer"),
				),
			},
		},
	})
}
