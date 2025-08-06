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

func TestInfinityDNSServer(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateDNSServer API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/dns_server/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/dns_server/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.DNSServer{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/dns_server/123/",
		Address:     "192.168.1.50",
		Description: "Cloudflare DNS",
	}

	// Mock the GetDNSServer API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/dns_server/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dns_server := args.Get(2).(*config.DNSServer)
		*dns_server = *mockState
	}).Maybe()

	// Mock the UpdateDNSServer API call
	client.On("PutJSON", mock.Anything, "configuration/v1/dns_server/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.DNSServerUpdateRequest)
		dns_server := args.Get(3).(*config.DNSServer)

		// Update mock state based on request
		if updateRequest.Address != "" {
			mockState.Address = updateRequest.Address
		}
		if updateRequest.Description != "" {
			mockState.Description = updateRequest.Description
		}

		// Return updated state
		*dns_server = *mockState
	}).Maybe()

	// Mock the DeleteDNSServer API call - use mock.MatchedBy to match dynamic ID
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/dns_server/123/"
	}), mock.Anything).Return(nil)

	testInfinityDNSServer(t, client)
}

func testInfinityDNSServer(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_dns_server_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_dns_server.cloudflare-dns", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_dns_server.cloudflare-dns", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_dns_server.cloudflare-dns", "address", "192.168.1.50"),
					resource.TestCheckResourceAttr("pexip_infinity_dns_server.cloudflare-dns", "description", "Cloudflare DNS"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_dns_server_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_dns_server.cloudflare-dns", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_dns_server.cloudflare-dns", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_dns_server.cloudflare-dns", "address", "192.168.1.50"),
					resource.TestCheckResourceAttr("pexip_infinity_dns_server.cloudflare-dns", "description", "Cloudflare DNS - updated"),
				),
			},
		},
	})
}
