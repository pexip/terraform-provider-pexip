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

func TestInfinityMediaLibraryEntryIntegration(t *testing.T) {
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

	testInfinityMediaLibraryEntryIntegration(t, client)
}

func testInfinityMediaLibraryEntryIntegration(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_entry_full_integration"),
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
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_entry_min_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_entry.media_library_entry-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_entry.media_library_entry-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "name", "tf-test-media-library-entry"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "description", ""),
					resource.TestCheckResourceAttrSet("pexip_infinity_media_library_entry.media_library_entry-test", "uuid"),
					resource.TestCheckResourceAttr("pexip_infinity_media_library_entry.media_library_entry-test", "file_name", "earth.mp4"),
				),
			},
			// Step 3: Destroy resources before recreate-from-scratch test
			{
				Config:       test.LoadTestFolder(t, "resource_infinity_media_library_entry_min_integration"),
				ResourceName: "pexip_infinity_media_library_entry.media_library_entry-test",
				Destroy:      true,
			},
			// Step 4: Create with min config (after destroy) - no description set, uuid is computed
			{
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_entry_min_integration"),
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
				Config: test.LoadTestFolder(t, "resource_infinity_media_library_entry_full_integration"),
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
