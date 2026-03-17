/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

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

func TestInfinityNTPServer(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking - initialize with defaults
	mockState := &config.NTPServer{
		ID:          1,
		ResourceURI: "/api/admin/configuration/v1/ntp_server/1/",
		Address:     "0.europe.pool.ntp.org",
		Description: "",
	}

	// Step 1: Create with full config
	client.On("PostWithResponse", mock.Anything, "configuration/v1/ntp_server/", mock.Anything, mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/ntp_server/1/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.NTPServerCreateRequest)
		mockState.Address = req.Address
		mockState.Description = req.Description
	}).Once()

	// Step 2: Update to min config
	client.On("PutJSON", mock.Anything, "configuration/v1/ntp_server/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.NTPServerUpdateRequest)
		mockState.Address = req.Address
		mockState.Description = req.Description
		if args.Get(3) != nil {
			ntp := args.Get(3).(*config.NTPServer)
			*ntp = *mockState
		}
	}).Once()

	// Step 3: Delete
	client.On("DeleteJSON", mock.Anything, "configuration/v1/ntp_server/1/", mock.Anything).Return(nil).Maybe()

	// Step 4: Recreate with min config
	client.On("PostWithResponse", mock.Anything, "configuration/v1/ntp_server/", mock.Anything, mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/ntp_server/1/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.NTPServerCreateRequest)
		mockState.Address = req.Address
		mockState.Description = req.Description
	}).Once()

	// Step 5: Update to full config
	client.On("PutJSON", mock.Anything, "configuration/v1/ntp_server/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.NTPServerUpdateRequest)
		mockState.Address = req.Address
		mockState.Description = req.Description
		if args.Get(3) != nil {
			ntp := args.Get(3).(*config.NTPServer)
			*ntp = *mockState
		}
	}).Once()

	// Mock Read operations (GetJSON) - used throughout all steps
	client.On("GetJSON", mock.Anything, "configuration/v1/ntp_server/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		ntp := args.Get(3).(*config.NTPServer)
		*ntp = *mockState
	}).Maybe()

	testInfinityNTPServer(t, client)
}

func testInfinityNTPServer(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				// Step 1: Create with full config
				Config: test.LoadTestFolder(t, "resource_infinity_ntp_server_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ntp_server.tf-test-ntp", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ntp_server.tf-test-ntp", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ntp_server.tf-test-ntp", "address", "1.europe.pool.ntp.org"),
					resource.TestCheckResourceAttr("pexip_infinity_ntp_server.tf-test-ntp", "description", "tf-test NTP Server Description"),
				),
			},
			{
				// Step 2: Update to min config (clear optional fields, reset to defaults)
				Config: test.LoadTestFolder(t, "resource_infinity_ntp_server_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ntp_server.tf-test-ntp", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ntp_server.tf-test-ntp", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ntp_server.tf-test-ntp", "address", "0.europe.pool.ntp.org"),
					resource.TestCheckResourceAttr("pexip_infinity_ntp_server.tf-test-ntp", "description", ""),
				),
			},
			{
				// Step 3: Destroy
				Config:  test.LoadTestFolder(t, "resource_infinity_ntp_server_min"),
				Destroy: true,
			},
			{
				// Step 4: Create with min config (after destroy)
				Config: test.LoadTestFolder(t, "resource_infinity_ntp_server_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ntp_server.tf-test-ntp", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ntp_server.tf-test-ntp", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ntp_server.tf-test-ntp", "address", "0.europe.pool.ntp.org"),
					resource.TestCheckResourceAttr("pexip_infinity_ntp_server.tf-test-ntp", "description", ""),
				),
			},
			{
				// Step 5: Update to full config
				Config: test.LoadTestFolder(t, "resource_infinity_ntp_server_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ntp_server.tf-test-ntp", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ntp_server.tf-test-ntp", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ntp_server.tf-test-ntp", "address", "1.europe.pool.ntp.org"),
					resource.TestCheckResourceAttr("pexip_infinity_ntp_server.tf-test-ntp", "description", "tf-test NTP Server Description"),
				),
			},
		},
	})
}
