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

func TestInfinityMediaLibraryPlaylistIntegration(t *testing.T) {
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

	testInfinityMediaLibraryPlaylistIntegration(t, client)
}

func testInfinityMediaLibraryPlaylistIntegration(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_playlist_full_integration"),
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
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_playlist_min_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist.media_library_playlist-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_playlist.media_library_playlist-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "name", "tf-test-media-library-playlist"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "loop", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_playlist.media_library_playlist-test", "shuffle", "false"),
				),
			},
			// Step 3: Destroy resources before recreate-from-scratch test
			{
				Config:       test.LoadTestFolder(t, "resource_infinity_media_library_playlist_min_integration"),
				ResourceName: "pexip_infinity_media_library_playlist.media_library_playlist-test",
				Destroy:      true,
			},
			// Step 4: Create with min config (after destroy)
			{
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_playlist_min_integration"),
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
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_playlist_full_integration"),
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
