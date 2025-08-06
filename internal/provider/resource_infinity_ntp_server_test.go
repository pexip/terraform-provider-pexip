/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"
	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityNTPServer(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/ntp_server/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/ntp_server/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Track state to return different values before and after update
	updated := false

	// Mock the API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/ntp_server/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		ntp := args.Get(2).(*config.NTPServer)
		if updated {
			*ntp = config.NTPServer{
				ID:          123,
				Address:     "2.pool.ntp.org",
				Description: "NTP server 2 - updated",
				ResourceURI: "/api/admin/configuration/v1/ntp_server/123/",
			}
		} else {
			*ntp = config.NTPServer{
				ID:          123,
				Address:     "2.pool.ntp.org",
				Description: "NTP server 2",
				ResourceURI: "/api/admin/configuration/v1/ntp_server/123/",
			}
		}
	}).Maybe() // Called multiple times for reads

	// Mock the API call - use mock.MatchedBy to match dynamic ID
	client.On("PutJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/ntp_server/123/"
	}), mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updated = true // Mark as updated for subsequent reads
		ntp := args.Get(3).(*config.NTPServer)
		*ntp = config.NTPServer{
			ID:          123,
			Address:     "2.pool.ntp.org",
			Description: "NTP server 2 - updated",
			ResourceURI: "/api/admin/configuration/v1/ntp_server/123/",
		}
	}).Maybe()

	// Mock the API call - use mock.MatchedBy to match dynamic ID
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/ntp_server/123/"
	}), mock.Anything).Return(nil)

	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ntp_server_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ntp_server.ntp-2", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ntp_server.ntp-2", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ntp_server.ntp-2", "address", "2.pool.ntp.org"),
					resource.TestCheckResourceAttr("pexip_infinity_ntp_server.ntp-2", "description", "NTP server 2"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ntp_server_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ntp_server.ntp-2", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ntp_server.ntp-2", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ntp_server.ntp-2", "address", "2.pool.ntp.org"),
					resource.TestCheckResourceAttr("pexip_infinity_ntp_server.ntp-2", "description", "NTP server 2 - updated"),
				),
			},
		},
	})
}
