//go:build integration

/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"crypto/tls"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityMediaLibraryPlaylistEntryIntegration(t *testing.T) {
	_ = os.Setenv("TF_ACC", "1")

	client, err := infinity.New(
		infinity.WithBaseURL(test.INFINITY_BASE_URL),
		infinity.WithBasicAuth(test.INFINITY_USERNAME, test.INFINITY_PASSWORD),
		infinity.WithMaxRetries(2),
		infinity.WithTransport(&http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // We need this because default certificate is not trusted
				MinVersion:         tls.VersionTLS12,
			},
			MaxIdleConns:        30,
			MaxIdleConnsPerHost: 5,
			IdleConnTimeout:     60 * time.Second,
		}),
	)
	require.NoError(t, err)

	testAccInfinityMediaLibraryPlaylistEntry(t, client)
}

func testAccInfinityMediaLibraryPlaylistEntry(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				// Step 1: Create with full config
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_playlist_entry_full_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist_entry.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist_entry.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist_entry.test", "entry_type", "MEDIA"),
					resource.TestCheckResourceAttrPair(
						"pexip_infinity_media_library_playlist_entry.test", "media",
						"pexip_infinity_media_library_entry.tf-test-media-entry", "id",
					),
					resource.TestCheckResourceAttrPair(
						"pexip_infinity_media_library_playlist_entry.test", "playlist",
						"pexip_infinity_media_library_playlist.tf-test-playlist", "id",
					),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist_entry.test", "position", "5"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist_entry.test", "playcount", "3"),
				),
			},
			{
				// Step 2: Update to min config
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_playlist_entry_min_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist_entry.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist_entry.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist_entry.test", "entry_type", "MEDIA"),
					resource.TestCheckResourceAttrPair(
						"pexip_infinity_media_library_playlist_entry.test", "playlist",
						"pexip_infinity_media_library_playlist.tf-test-playlist", "id",
					),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist_entry.test", "position", "1"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist_entry.test", "playcount", "1"),
				),
			},
			{
				// Step 3: Destroy
				Config:  test.LoadTestFolder(t, "resource_infinity_media_library_playlist_entry_min_integration"),
				Destroy: true,
			},
			{
				// Step 4: Recreate with min config
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_playlist_entry_min_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist_entry.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist_entry.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist_entry.test", "entry_type", "MEDIA"),
					resource.TestCheckResourceAttrPair(
						"pexip_infinity_media_library_playlist_entry.test", "playlist",
						"pexip_infinity_media_library_playlist.tf-test-playlist", "id",
					),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist_entry.test", "position", "1"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist_entry.test", "playcount", "1"),
				),
			},
			{
				// Step 5: Update to full config
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_playlist_entry_full_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist_entry.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist_entry.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist_entry.test", "entry_type", "MEDIA"),
					resource.TestCheckResourceAttrPair(
						"pexip_infinity_media_library_playlist_entry.test", "media",
						"pexip_infinity_media_library_entry.tf-test-media-entry", "id",
					),
					resource.TestCheckResourceAttrPair(
						"pexip_infinity_media_library_playlist_entry.test", "playlist",
						"pexip_infinity_media_library_playlist.tf-test-playlist", "id",
					),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist_entry.test", "position", "5"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist_entry.test", "playcount", "3"),
				),
			},
		},
	})
}
