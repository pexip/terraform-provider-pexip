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

func TestInfinityMediaLibraryPlaylist(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateMedialibraryplaylist API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/media_library_playlist/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/media_library_playlist/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.MediaLibraryPlaylist{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/media_library_playlist/123/",
		Name:        "media_library_playlist-test",
		Description: "Test MediaLibraryPlaylist",
		Loop:        true,
		Shuffle:     true,
	}

	// Mock the GetMedialibraryplaylist API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/media_library_playlist/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		media_library_playlist := args.Get(2).(*config.MediaLibraryPlaylist)
		*media_library_playlist = *mockState
	}).Maybe()

	// Mock the UpdateMedialibraryplaylist API call
	client.On("PutJSON", mock.Anything, "configuration/v1/media_library_playlist/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.MediaLibraryPlaylistUpdateRequest)
		media_library_playlist := args.Get(3).(*config.MediaLibraryPlaylist)

		// Update mock state based on request
		if updateRequest.Description != "" {
			mockState.Description = updateRequest.Description
		}
		if updateRequest.Loop != nil {
			mockState.Loop = *updateRequest.Loop
		}
		if updateRequest.Shuffle != nil {
			mockState.Shuffle = *updateRequest.Shuffle
		}

		// Return updated state
		*media_library_playlist = *mockState
	}).Maybe()

	// Mock the DeleteMedialibraryplaylist API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/media_library_playlist/123/"
	}), mock.Anything).Return(nil)

	testInfinityMediaLibraryPlaylist(t, client)
}

func testInfinityMediaLibraryPlaylist(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_playlist_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist.media_library_playlist-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist.media_library_playlist-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "name", "media_library_playlist-test"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "description", "Test MediaLibraryPlaylist"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "loop", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "shuffle", "true"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_playlist_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist.media_library_playlist-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist.media_library_playlist-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "name", "media_library_playlist-test"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "description", "Updated Test MediaLibraryPlaylist"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "loop", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "shuffle", "false"),
				),
			},
		},
	})
}
