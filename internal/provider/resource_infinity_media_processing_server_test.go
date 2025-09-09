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

func TestInfinityMediaProcessingServer(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateMediaprocessingserver API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/media_processing_server/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/media_processing_server/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.MediaProcessingServer{
		ID:           123,
		ResourceURI:  "/api/admin/configuration/v1/media_processing_server/123/",
		FQDN:         "test-value",
		AppID:        test.StringPtr("test-value"),
		PublicJWTKey: "test-value",
	}

	// Mock the GetMediaprocessingserver API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/media_processing_server/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		media_processing_server := args.Get(3).(*config.MediaProcessingServer)
		*media_processing_server = *mockState
	}).Maybe()

	// Mock the UpdateMediaprocessingserver API call
	client.On("PutJSON", mock.Anything, "configuration/v1/media_processing_server/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.MediaProcessingServerUpdateRequest)
		media_processing_server := args.Get(3).(*config.MediaProcessingServer)

		// Update mock state based on request
		if updateRequest.FQDN != "" {
			mockState.FQDN = updateRequest.FQDN
		}
		if updateRequest.AppID != nil {
			mockState.AppID = updateRequest.AppID
		}
		if updateRequest.PublicJWTKey != "" {
			mockState.PublicJWTKey = updateRequest.PublicJWTKey
		}

		// Return updated state
		*media_processing_server = *mockState
	}).Maybe()

	// Mock the DeleteMediaprocessingserver API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/media_processing_server/123/"
	}), mock.Anything).Return(nil)

	testInfinityMediaProcessingServer(t, client)
}

func testInfinityMediaProcessingServer(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_media_processing_server_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_processing_server.media_processing_server-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_processing_server.media_processing_server-test", "resource_id"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_media_processing_server_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_processing_server.media_processing_server-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_processing_server.media_processing_server-test", "resource_id"),
				),
			},
		},
	})
}
