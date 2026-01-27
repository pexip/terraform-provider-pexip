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

func TestInfinityMediaLibraryEntry(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Register cleanup function to remove test artifacts
	t.Cleanup(func() {
		_ = os.Remove("earth.mp4")
		_ = os.Remove("rain.mp4")
	})

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking
	mockState := &config.MediaLibraryEntry{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/media_library_entry/123/",
		Name:        "media_library_entry-test",
		Description: "Test MediaLibraryEntry",
		UUID:        "test-value",
		FileName:    "earth.mp4",
	}

	// Mock the CreateMediaLibraryEntry API call (now using multipart form)
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/media_library_entry/123/",
	}
	client.On("PostMultipartFormWithFieldsAndResponse", mock.Anything, "configuration/v1/media_library_entry/",
		mock.MatchedBy(func(fields map[string]string) bool {
			return fields["name"] == "media_library_entry-test" && fields["description"] == "Test MediaLibraryEntry"
		}), "media_file", mock.Anything, mock.Anything, mock.Anything).Return(createResponse, nil)

	// Mock the GetMediaLibraryEntry API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/media_library_entry/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		media_library_entry := args.Get(3).(*config.MediaLibraryEntry)
		*media_library_entry = *mockState
	}).Maybe()

	// Mock the UpdateMediaLibraryEntry API call (now using multipart form with PATCH)
	client.On("PatchMultipartFormWithFieldsAndResponse", mock.Anything, "configuration/v1/media_library_entry/123/",
		mock.MatchedBy(func(fields map[string]string) bool {
			return fields["description"] == "Updated Test MediaLibraryEntry"
		}), "media_file", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Run(func(args mock.Arguments) {
		fields := args.Get(2).(map[string]string)
		media_library_entry := args.Get(6).(*config.MediaLibraryEntry)

		// Update mock state based on fields
		if desc, ok := fields["description"]; ok && desc != "" {
			mockState.Description = desc
		}
		if uuid, ok := fields["uuid"]; ok && uuid != "" {
			mockState.UUID = uuid
		}
		// Update filename from the file parameter
		if filename, ok := args.Get(4).(string); ok {
			mockState.FileName = filename
		}

		// Return updated state
		*media_library_entry = *mockState
	}).Maybe()

	// Mock the DeleteMediaLibraryEntry API call
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
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "file_name", "earth.mp4"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_entry_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_entry.media_library_entry-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_entry.media_library_entry-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "name", "media_library_entry-test"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "description", "Updated Test MediaLibraryEntry"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "file_name", "rain.mp4"),
				),
			},
		},
	})
}
