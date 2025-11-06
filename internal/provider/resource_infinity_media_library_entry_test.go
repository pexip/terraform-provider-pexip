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

func TestInfinityMediaLibraryEntry(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateMedialibraryentry API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/media_library_entry/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/media_library_entry/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.MediaLibraryEntry{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/media_library_entry/123/",
		Name:        "media_library_entry-test",
		Description: "Test MediaLibraryEntry",
		UUID:        "test-value",
		FileName:    "media_library_entry-test",
	}

	// Mock the GetMedialibraryentry API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/media_library_entry/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		media_library_entry := args.Get(3).(*config.MediaLibraryEntry)
		*media_library_entry = *mockState
	}).Maybe()

	// Mock the UpdateMedialibraryentry API call
	client.On("PutJSON", mock.Anything, "configuration/v1/media_library_entry/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.MediaLibraryEntryUpdateRequest)
		media_library_entry := args.Get(3).(*config.MediaLibraryEntry)

		// Update mock state based on request
		if updateRequest.Description != "" {
			mockState.Description = updateRequest.Description
		}
		if updateRequest.UUID != "" {
			mockState.UUID = updateRequest.UUID
		}

		// Return updated state
		*media_library_entry = *mockState
	}).Maybe()

	// Mock the DeleteMedialibraryentry API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/media_library_entry/123/"
	}), mock.Anything).Return(nil)

	testInfinityMediaLibraryEntry(t, client)
}

func testInfinityMediaLibraryEntry(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_entry_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_entry.media_library_entry-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_entry.media_library_entry-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "name", "media_library_entry-test"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "description", "Test MediaLibraryEntry"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "file_name", "media_library_entry-test"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_entry_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_entry.media_library_entry-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_entry.media_library_entry-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "name", "media_library_entry-test"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "description", "Updated Test MediaLibraryEntry"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "file_name", "media_library_entry-test"),
				),
			},
		},
	})
}
