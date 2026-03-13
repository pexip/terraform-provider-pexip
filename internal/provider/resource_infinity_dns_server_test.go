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

func TestInfinityDNSServer(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared mock state for the DNS server
	mockDNS := &config.DNSServer{
		ID:          1,
		ResourceURI: "/api/admin/configuration/v1/dns_server/1/",
		Address:     "4.2.2.2",
		Description: "",
	}

	// Step 1: Create DNS server with full config (4.2.2.2 with description)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/dns_server/", mock.Anything, mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/dns_server/1/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.DNSServerCreateRequest)
		mockDNS.Address = req.Address
		mockDNS.Description = req.Description
	}).Once()

	// Step 2: Update to min config (4.2.2.1, clear description)
	client.On("PutJSON", mock.Anything, "configuration/v1/dns_server/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.DNSServerUpdateRequest)
		// Update the shared mock state
		mockDNS.Address = req.Address
		mockDNS.Description = req.Description
		// Also update the response object if provided
		if args.Get(3) != nil {
			dns := args.Get(3).(*config.DNSServer)
			*dns = *mockDNS
		}
	}).Once()

	// Step 3: Delete DNS server (and final cleanup after test)
	client.On("DeleteJSON", mock.Anything, "configuration/v1/dns_server/1/", mock.Anything).Return(nil).Maybe()

	// Step 4: Recreate DNS server with min config (4.2.2.1, no description)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/dns_server/", mock.Anything, mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/dns_server/1/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.DNSServerCreateRequest)
		mockDNS.Address = req.Address
		mockDNS.Description = req.Description
	}).Once()

	// Step 5: Update to full config (4.2.2.2, add description)
	client.On("PutJSON", mock.Anything, "configuration/v1/dns_server/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.DNSServerUpdateRequest)
		// Update the shared mock state
		mockDNS.Address = req.Address
		mockDNS.Description = req.Description
		// Also update the response object if provided
		if args.Get(3) != nil {
			dns := args.Get(3).(*config.DNSServer)
			*dns = *mockDNS
		}
	}).Once()

	// Mock Read operations (GetJSON) - used throughout all steps
	client.On("GetJSON", mock.Anything, "configuration/v1/dns_server/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dns := args.Get(3).(*config.DNSServer)
		*dns = *mockDNS
	}).Maybe()

	testInfinityDNSServer(t, client)
}

func testInfinityDNSServer(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				// Step 1: Create with full config
				Config: test.LoadTestFolder(t, "resource_infinity_dns_server_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_dns_server.tf-test-dns", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_dns_server.tf-test-dns", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_dns_server.tf-test-dns", "address", "4.2.2.2"),
					resource.TestCheckResourceAttr("pexip_infinity_dns_server.tf-test-dns", "description", "tf-test Level 3 DNS Server"),
				),
			},
			{
				// Step 2: Update to min config (clear description)
				Config: test.LoadTestFolder(t, "resource_infinity_dns_server_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_dns_server.tf-test-dns", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_dns_server.tf-test-dns", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_dns_server.tf-test-dns", "address", "4.2.2.1"),
					resource.TestCheckResourceAttr("pexip_infinity_dns_server.tf-test-dns", "description", ""),
				),
			},
			{
				// Step 3: Destroy
				Config:  test.LoadTestFolder(t, "resource_infinity_dns_server_min"),
				Destroy: true,
			},
			{
				// Step 4: Create with min config (after destroy)
				Config: test.LoadTestFolder(t, "resource_infinity_dns_server_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_dns_server.tf-test-dns", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_dns_server.tf-test-dns", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_dns_server.tf-test-dns", "address", "4.2.2.1"),
					resource.TestCheckResourceAttr("pexip_infinity_dns_server.tf-test-dns", "description", ""),
				),
			},
			{
				// Step 5: Update to full config
				Config: test.LoadTestFolder(t, "resource_infinity_dns_server_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_dns_server.tf-test-dns", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_dns_server.tf-test-dns", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_dns_server.tf-test-dns", "address", "4.2.2.2"),
					resource.TestCheckResourceAttr("pexip_infinity_dns_server.tf-test-dns", "description", "tf-test Level 3 DNS Server"),
				),
			},
		},
	})
}
