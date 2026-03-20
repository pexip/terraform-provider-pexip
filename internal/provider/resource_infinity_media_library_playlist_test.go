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

func TestInfinityMediaLibraryPlaylist(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking - starts with full config
	mockState := &config.MediaLibraryPlaylist{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/media_library_playlist/123/",
		Name:        "tf-test-media-library-playlist",
		Description: "tf-test media library playlist description",
		Loop:        true,
		Shuffle:     true,
	}

	// Mock the CreateMedialibraryplaylist API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/media_library_playlist/123/",
	}
	// Step 1: Create with full config
	client.On("PostWithResponse", mock.Anything, "configuration/v1/media_library_playlist/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.MediaLibraryPlaylistCreateRequest)
		mockState.Name = createReq.Name
		mockState.Description = createReq.Description
		mockState.Loop = createReq.Loop
		mockState.Shuffle = createReq.Shuffle
	}).Once()

	// Step 4: Create with min config (after delete)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/media_library_playlist/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.MediaLibraryPlaylistCreateRequest)
		mockState.Name = createReq.Name
		mockState.Description = createReq.Description
		mockState.Loop = createReq.Loop
		mockState.Shuffle = createReq.Shuffle
	}).Once()

	// Mock the GetMedialibraryplaylist API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/media_library_playlist/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		media_library_playlist := args.Get(3).(*config.MediaLibraryPlaylist)
		*media_library_playlist = *mockState
	}).Maybe()

	// Mock the UpdateMedialibraryplaylist API call
	client.On("PutJSON", mock.Anything, "configuration/v1/media_library_playlist/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.MediaLibraryPlaylistUpdateRequest)
		media_library_playlist := args.Get(3).(*config.MediaLibraryPlaylist)

		// Update mock state based on request
		if updateRequest.Name != "" {
			mockState.Name = updateRequest.Name
		}
		// Description can be cleared by sending empty string
		mockState.Description = updateRequest.Description
		if updateRequest.Loop != nil {
			mockState.Loop = *updateRequest.Loop
		}
		if updateRequest.Shuffle != nil {
			mockState.Shuffle = *updateRequest.Shuffle
		}

		// Return updated state
		*media_library_playlist = *mockState
	}).Times(2) // Step 2: Update to min, Step 5: Update to full

	// Mock the DeleteMedialibraryplaylist API call
	// Step 3: Delete, and final cleanup at end of test
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/media_library_playlist/123/"
	}), mock.Anything).Return(nil).Times(2)

	testInfinityMediaLibraryPlaylist(t, client)
}

func testInfinityMediaLibraryPlaylist(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_playlist_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist.media_library_playlist-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist.media_library_playlist-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "name", "tf-test-media-library-playlist"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "description", "tf-test media library playlist description"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "loop", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "shuffle", "true"),
				),
			},
			// Step 2: Update to min config (description cleared, loop/shuffle reset to defaults)
			{
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_playlist_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist.media_library_playlist-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist.media_library_playlist-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "name", "tf-test-media-library-playlist"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "loop", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "shuffle", "false"),
				),
			},
			// Step 3: Delete the resource
			{
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_playlist_delete"),
			},
			// Step 4: Create with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_playlist_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist.media_library_playlist-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist.media_library_playlist-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "name", "tf-test-media-library-playlist"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "loop", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "shuffle", "false"),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_playlist_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist.media_library_playlist-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist.media_library_playlist-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "name", "tf-test-media-library-playlist"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "description", "tf-test media library playlist description"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "loop", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "shuffle", "true"),
				),
			},
		},
	})
}
