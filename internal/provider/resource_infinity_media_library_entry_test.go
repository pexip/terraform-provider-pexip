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

	// Shared state for mocking - starts with full config
	mockState := &config.MediaLibraryEntry{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/media_library_entry/123/",
		Name:        "tf-test-media-library-entry",
		Description: "tf-test media library entry description",
		UUID:        "api-generated-uuid",
		FileName:    "rain.mp4",
	}

	// Mock the CreateMediaLibraryEntry API call (now using multipart form)
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/media_library_entry/123/",
	}
	// Step 1: Create with full config
	client.On("PostMultipartFormWithFieldsAndResponse", mock.Anything, "configuration/v1/media_library_entry/",
		mock.Anything, "media_file", mock.Anything, mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		fields := args.Get(2).(map[string]string)
		// Update mock state based on create request fields
		if name, ok := fields["name"]; ok {
			mockState.Name = name
		}
		if desc, ok := fields["description"]; ok {
			mockState.Description = desc
		} else {
			mockState.Description = ""
		}
		// UUID is computed only (assigned by API)
		mockState.UUID = "api-generated-uuid"
		// Update filename from the file parameter
		if filename, ok := args.Get(4).(string); ok {
			mockState.FileName = filename
		}
	}).Once()

	// Step 4: Create with min config (after delete)
	client.On("PostMultipartFormWithFieldsAndResponse", mock.Anything, "configuration/v1/media_library_entry/",
		mock.Anything, "media_file", mock.Anything, mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		fields := args.Get(2).(map[string]string)
		// Update mock state based on create request fields
		mockState.Name = fields["name"]
		if desc, ok := fields["description"]; ok {
			mockState.Description = desc
		} else {
			mockState.Description = ""
		}
		// UUID is computed only (assigned by API)
		mockState.UUID = "api-generated-uuid"
		// Update filename from the file parameter
		if filename, ok := args.Get(4).(string); ok {
			mockState.FileName = filename
		}
	}).Once()

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
			// Description can be cleared by sending empty string
			mockState.Description = desc
		}
		// UUID is computed only (not updated by user)
		// Update filename from the file parameter
		if filename, ok := args.Get(4).(string); ok {
			mockState.FileName = filename
		}

		// Return updated state
		*media_library_entry = *mockState
	}).Times(2) // Step 2: Update to min, Step 5: Update to full

	// Mock the DeleteMediaLibraryEntry API call
	// Step 3: Delete, and final cleanup at end of test
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/media_library_entry/123/"
	}), mock.Anything).Return(nil).Times(2)

	testInfinityMediaLibraryEntry(t, client)
}

func testInfinityMediaLibraryEntry(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_entry_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_entry.media_library_entry-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_entry.media_library_entry-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "name", "tf-test-media-library-entry"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "description", "tf-test media library entry description"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_entry.media_library_entry-test", "uuid"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "file_name", "rain.mp4"),
				),
			},
			// Step 2: Update to min config (description is cleared, uuid persists)
			{
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_entry_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_entry.media_library_entry-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_entry.media_library_entry-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "name", "tf-test-media-library-entry"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "description", ""),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_entry.media_library_entry-test", "uuid"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "file_name", "earth.mp4"),
				),
			},
			// Step 3: Delete the resource
			{
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_entry_delete"),
			},
			// Step 4: Create with min config - no description set, uuid is computed
			{
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_entry_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_entry.media_library_entry-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_entry.media_library_entry-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "name", "tf-test-media-library-entry"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_entry.media_library_entry-test", "uuid"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "file_name", "earth.mp4"),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_entry_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_entry.media_library_entry-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_entry.media_library_entry-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "name", "tf-test-media-library-entry"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "description", "tf-test media library entry description"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_entry.media_library_entry-test", "uuid"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "file_name", "rain.mp4"),
				),
			},
		},
	})
}
