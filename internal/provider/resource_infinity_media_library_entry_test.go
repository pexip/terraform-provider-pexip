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

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking - starts with min config
	mockState := &config.MediaLibraryEntry{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/media_library_entry/123/",
		Name:        "tf-test-media-library-entry",
		Description: "",
		UUID:        "test-value",
		FileName:    "earth.mp4",
	}

	// Mock the CreateMediaLibraryEntry API call (now using multipart form)
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/media_library_entry/123/",
	}
	client.On("PostMultipartFormWithFieldsAndResponse", mock.Anything, "configuration/v1/media_library_entry/",
		mock.Anything, "media_file", mock.Anything, mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		fields := args.Get(2).(map[string]string)
		// Update mock state based on create request fields
		mockState.Name = fields["name"]
		if desc, ok := fields["description"]; ok {
			mockState.Description = desc
		}
		// Update filename from the file parameter
		if filename, ok := args.Get(4).(string); ok {
			mockState.FileName = filename
		}
	}).Maybe()

	// Mock the GetMediaLibraryEntry API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/media_library_entry/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		media_library_entry := args.Get(3).(*config.MediaLibraryEntry)
		*media_library_entry = *mockState
	}).Maybe()

	// Mock the UpdateMediaLibraryEntry API call (now using multipart form with PATCH)
	client.On("PatchMultipartFormWithFieldsAndResponse", mock.Anything, "configuration/v1/media_library_entry/123/",
		mock.Anything, "media_file", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Run(func(args mock.Arguments) {
		fields := args.Get(2).(map[string]string)
		media_library_entry := args.Get(6).(*config.MediaLibraryEntry)

		// Update mock state based on fields
		if name, ok := fields["name"]; ok && name != "" {
			mockState.Name = name
		}
		if desc, ok := fields["description"]; ok {
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
			// Step 1: Create with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_entry_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_entry.media_library_entry-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_entry.media_library_entry-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "name", "tf-test-media-library-entry"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "file_name", "earth.mp4"),
				),
			},
			// Step 2: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_entry_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_entry.media_library_entry-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_entry.media_library_entry-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "name", "tf-test-media-library-entry-full"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "description", "tf-test media library entry description"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "file_name", "rain.mp4"),
				),
			},
		},
	})
}
