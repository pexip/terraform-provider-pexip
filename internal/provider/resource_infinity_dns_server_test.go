package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"
	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/pexip/terraform-provider-pexip/internal/test"
	"github.com/stretchr/testify/mock"
)

func TestInfinityDNSServer(t *testing.T) {
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateDNSServer API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/dns_server/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/dns_server/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Track state to return different values before and after update
	updated := false

	// Mock the GetDNSServer API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/dns_server/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dns := args.Get(2).(*config.DNSServer)
		if updated {
			*dns = config.DNSServer{
				ID:          123,
				Address:     "1.1.1.1",
				Description: "Cloudflare DNS - updated",
				ResourceURI: "/api/admin/configuration/v1/dns_server/123/",
			}
		} else {
			*dns = config.DNSServer{
				ID:          123,
				Address:     "1.1.1.1",
				Description: "Cloudflare DNS",
				ResourceURI: "/api/admin/configuration/v1/dns_server/123/",
			}
		}
	}).Maybe() // Called multiple times for reads

	// Mock the UpdateDNSServer API call - use mock.MatchedBy to match dynamic ID
	client.On("PutJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/dns_server/123/"
	}), mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updated = true // Mark as updated for subsequent reads
		dns := args.Get(3).(*config.DNSServer)
		*dns = config.DNSServer{
			ID:          123,
			Address:     "1.1.1.1",
			Description: "Cloudflare DNS - updated",
			ResourceURI: "/api/admin/configuration/v1/dns_server/123/",
		}
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
				Config: test.LoadTestData(t, "resource_infinity_dns_server_basic.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_dns_server.cloudflare-dns", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_dns_server.cloudflare-dns", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_dns_server.cloudflare-dns", "address", "1.1.1.1"),
					resource.TestCheckResourceAttr("pexip_infinity_dns_server.cloudflare-dns", "description", "Cloudflare DNS"),
				),
			},
			{
				Config: test.LoadTestData(t, "resource_infinity_dns_server_basic_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_dns_server.cloudflare-dns", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_dns_server.cloudflare-dns", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_dns_server.cloudflare-dns", "address", "1.1.1.1"),
					resource.TestCheckResourceAttr("pexip_infinity_dns_server.cloudflare-dns", "description", "Cloudflare DNS - updated"),
				),
			},
		},
	})
}
