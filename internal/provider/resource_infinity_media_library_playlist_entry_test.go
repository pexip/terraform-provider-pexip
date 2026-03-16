/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

//go:build unit

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"
	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityMediaLibraryPlaylistEntry(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	client := infinity.NewClientMock()

	// Mock Media Library Entry (for "media" field reference)
	mockMediaEntry := &config.MediaLibraryEntry{}
	mediaEntryCreateResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/media_library_entry/1/",
	}
	client.On("PostMultipartFormWithFieldsAndResponse", mock.Anything, "configuration/v1/media_library_entry/", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mediaEntryCreateResponse, nil).Run(func(args mock.Arguments) {
		fields := args.Get(2).(map[string]string)
		*mockMediaEntry = config.MediaLibraryEntry{
			ID:          1,
			ResourceURI: "/api/admin/configuration/v1/media_library_entry/1/",
			Name:        fields["name"],
			Description: fields["description"],
			UUID:        fields["uuid"],
			FileName:    "test-file.mp4",
			MediaType:   "video",
			MediaFormat: "mp4",
			MediaSize:   1024,
		}
	})
	client.On("GetJSON", mock.Anything, "configuration/v1/media_library_entry/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		entry := args.Get(3).(*config.MediaLibraryEntry)
		*entry = *mockMediaEntry
	}).Maybe()
	client.On("PatchMultipartFormWithFieldsAndResponse", mock.Anything, "configuration/v1/media_library_entry/1/", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&types.PostResponse{}, nil).Run(func(args mock.Arguments) {
		fields := args.Get(2).(map[string]string)
		entry := args.Get(6).(*config.MediaLibraryEntry)
		mockMediaEntry.Name = fields["name"]
		mockMediaEntry.Description = fields["description"]
		mockMediaEntry.UUID = fields["uuid"]
		*entry = *mockMediaEntry
	}).Maybe()
	client.On("DeleteJSON", mock.Anything, "configuration/v1/media_library_entry/1/", mock.Anything).Return(nil)

	// Mock Media Library Playlist (for "playlist" field reference)
	mockPlaylist := &config.MediaLibraryPlaylist{}
	playlistCreateResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/media_library_playlist/2/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/media_library_playlist/", mock.Anything, mock.Anything).Return(playlistCreateResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.MediaLibraryPlaylistCreateRequest)
		*mockPlaylist = config.MediaLibraryPlaylist{
			ID:          2,
			ResourceURI: "/api/admin/configuration/v1/media_library_playlist/2/",
			Name:        createReq.Name,
			Description: createReq.Description,
			Loop:        createReq.Loop,
			Shuffle:     createReq.Shuffle,
		}
	})
	client.On("GetJSON", mock.Anything, "configuration/v1/media_library_playlist/2/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		playlist := args.Get(3).(*config.MediaLibraryPlaylist)
		*playlist = *mockPlaylist
	}).Maybe()
	client.On("PutJSON", mock.Anything, "configuration/v1/media_library_playlist/2/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.MediaLibraryPlaylistUpdateRequest)
		playlist := args.Get(3).(*config.MediaLibraryPlaylist)
		mockPlaylist.Name = updateReq.Name
		mockPlaylist.Description = updateReq.Description
		if updateReq.Loop != nil {
			mockPlaylist.Loop = *updateReq.Loop
		}
		if updateReq.Shuffle != nil {
			mockPlaylist.Shuffle = *updateReq.Shuffle
		}
		*playlist = *mockPlaylist
	}).Maybe()
	client.On("DeleteJSON", mock.Anything, "configuration/v1/media_library_playlist/2/", mock.Anything).Return(nil)

	// Shared mock state for the playlist entry
	mockState := &config.MediaLibraryPlaylistEntry{}

	// Mock CreateMediaLibraryPlaylistEntry
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/media_library_playlist_entry/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/media_library_playlist_entry/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.MediaLibraryPlaylistEntryCreateRequest)
		*mockState = config.MediaLibraryPlaylistEntry{
			ID:          123,
			ResourceURI: "/api/admin/configuration/v1/media_library_playlist_entry/123/",
			EntryType:   createReq.EntryType,
			Media:       createReq.Media,
			Playlist:    createReq.Playlist,
			Position:    createReq.Position,
			Playcount:   createReq.Playcount,
		}
	})

	// Mock GetMediaLibraryPlaylistEntry
	client.On("GetJSON", mock.Anything, "configuration/v1/media_library_playlist_entry/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		entry := args.Get(3).(*config.MediaLibraryPlaylistEntry)
		*entry = *mockState
	}).Maybe()

	// Mock UpdateMediaLibraryPlaylistEntry
	client.On("PutJSON", mock.Anything, "configuration/v1/media_library_playlist_entry/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.MediaLibraryPlaylistEntryUpdateRequest)
		entry := args.Get(3).(*config.MediaLibraryPlaylistEntry)

		// Update mock state based on request
		mockState.EntryType = updateReq.EntryType
		// Always update Media field (could be set to nil)
		mockState.Media = updateReq.Media
		// Always update Playlist field (could be set to nil)
		mockState.Playlist = updateReq.Playlist
		if updateReq.Position != nil {
			mockState.Position = *updateReq.Position
		}
		if updateReq.Playcount != nil {
			mockState.Playcount = *updateReq.Playcount
		}

		*entry = *mockState
	}).Maybe()

	// Mock DeleteMediaLibraryPlaylistEntry
	client.On("DeleteJSON", mock.Anything, "configuration/v1/media_library_playlist_entry/123/", mock.Anything).Return(nil)

	testInfinityMediaLibraryPlaylistEntry(t, client)
}

func testInfinityMediaLibraryPlaylistEntry(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				// Step 1: Create with full config
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_playlist_entry_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist_entry.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist_entry.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist_entry.test", "entry_type", "MEDIA"),
					resource.TestCheckResourceAttrPair(
						"pexip_infinity_media_library_playlist_entry.test", "media",
						"pexip_infinity_media_library_entry.test", "id",
					),
					resource.TestCheckResourceAttrPair(
						"pexip_infinity_media_library_playlist_entry.test", "playlist",
						"pexip_infinity_media_library_playlist.test", "id",
					),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist_entry.test", "position", "5"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist_entry.test", "playcount", "3"),
				),
			},
			{
				// Step 2: Update to min config
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_playlist_entry_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist_entry.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist_entry.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist_entry.test", "entry_type", "MEDIA"),
					resource.TestCheckResourceAttrPair(
						"pexip_infinity_media_library_playlist_entry.test", "playlist",
						"pexip_infinity_media_library_playlist.test", "id",
					),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist_entry.test", "position", "1"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist_entry.test", "playcount", "1"),
				),
			},
			{
				// Step 3: Destroy
				Config:  test.LoadTestFolder(t, "resource_infinity_media_library_playlist_entry_min"),
				Destroy: true,
			},
			{
				// Step 4: Recreate with min config
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_playlist_entry_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist_entry.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist_entry.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist_entry.test", "entry_type", "MEDIA"),
					resource.TestCheckResourceAttrPair(
						"pexip_infinity_media_library_playlist_entry.test", "playlist",
						"pexip_infinity_media_library_playlist.test", "id",
					),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist_entry.test", "position", "1"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist_entry.test", "playcount", "1"),
				),
			},
			{
				// Step 5: Update to full config
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_playlist_entry_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist_entry.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist_entry.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist_entry.test", "entry_type", "MEDIA"),
					resource.TestCheckResourceAttrPair(
						"pexip_infinity_media_library_playlist_entry.test", "media",
						"pexip_infinity_media_library_entry.test", "id",
					),
					resource.TestCheckResourceAttrPair(
						"pexip_infinity_media_library_playlist_entry.test", "playlist",
						"pexip_infinity_media_library_playlist.test", "id",
					),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist_entry.test", "position", "5"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist_entry.test", "playcount", "3"),
				),
			},
		},
	})
}
