/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

//go:build integration

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestAccInfinityMediaLibraryPlaylistEntry(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
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
